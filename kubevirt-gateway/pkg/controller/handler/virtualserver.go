package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"
	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	v1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kubevirt_gateway/v1"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/response"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/domain"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/repo"
	"k8s.io/client-go/tools/cache"
)

var _ framework.VirtualServerHandlerInterface = &VirtualServerHandler{}

type VirtualServerHandler struct {
	controller           framework.Interface
	virtualServerManager domain.VirtualServerManager
}

func NewVirtualServerHandler(controller framework.Interface, virtualServerManager domain.VirtualServerManager) (*VirtualServerHandler, error) {
	// watch 虚拟机状态变化，更新到数据库, 可以通过workers（int）传参，并发对workequeue进行消费
	// for i := 0; i < workers; i++ {
	// 	go virtualServerManager.Sync()
	// }
	go virtualServerManager.Sync()
	// 启动pvc informer
	go virtualServerManager.StartInformer()
	if ok := cache.WaitForCacheSync(virtualServerManager.StopCh, virtualServerManager.ResourcesSynced...); !ok {
		log.Infof("failed to wait for caches to sync")
		virtualServerManager.StopInformer()
		return nil, fmt.Errorf("failed to wait for caches to sync")
	}
	return &VirtualServerHandler{
		controller:           controller,
		virtualServerManager: virtualServerManager,
	}, nil
}

// CreateOSMVirtualServer godoc
//
//	@Summary		Create OSM VirtualServer
//	@Description	Create OSM VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string					true	"用户 JWT token"
//	@param			post			body		v1.OSMVirtualServerReq	true	"云服务器创建参数"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver [post]
func (w *VirtualServerHandler) CreateOSMVirtualServer(c *gcpctx.GCPContext) {
	username := c.GetUesrName()
	// logger: log with username and sessionuuid
	logger := c.Logger()

	logger.Infof("username: %s", username)
	// 0. 解析用户请求参数
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Infof("read request body err, %v", err)
		response.Response(c, response.ErrBindParams, nil)
		return
	}
	osmVirtualServerReq := v1.OSMVirtualServerReq{}
	if err := json.Unmarshal(body, &osmVirtualServerReq); err != nil {
		logger.Infof("read virtualserver %s err, %v", osmVirtualServerReq.Name, err)
		response.Response(c, response.ErrBindParams, nil)
		return
	}
	virtualserver := v1.TransformOSMVirtualServerReqToOSMVirtualServer(osmVirtualServerReq)
	// 1. 对下发服务器配置进行校验
	if err := virtualserver.CheckParam(); err != nil {
		logger.Infof("check virtualserver %s request param err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	// 2. 初始化数据
	if err := virtualserver.Init(username, w.virtualServerManager.Kubevirt, w.virtualServerManager.Kubecli, w.virtualServerManager.Repo); err != nil {
		logger.Infof("check virtualserver %s request param err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	// 3. 初始化资源
	products := w.controller.ComponentConfig().GetProducts()
	if err := virtualserver.InitResource(w.virtualServerManager.Repo, nil, products); err != nil {
		logger.Infof("check virtualserver %s request param err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	// 4. 将用户数据初始化到数据库，异步发起创建kubevirt虚拟机流程
	if err := w.virtualServerManager.AddVirtualServer(virtualserver); err != nil {
		logger.Infof("add virtualserver %s err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrStoreVirtualServer, nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, nil)
}

// CreateBSMVirtualServer godoc
//
//	@Summary		Create BSM VirtualServer
//	@Description	Create BSM VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string					true	"用户 JWT token"
//	@param			post			body		v1.BSMVirtualServerReq	true	"云服务器创建参数"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/bsm-virtualserver [post]
func (w *VirtualServerHandler) CreateBSMVirtualServer(c *gcpctx.GCPContext) {
	username := c.GetUesrName()
	// logger: log with username and sessionuuid
	logger := c.Logger()

	logger.Infof("username: %s", username)
	// 0. 解析用户请求参数
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Infof("read request body err, %v", err)
		response.Response(c, response.ErrBindParams, nil)
		return
	}
	bsmVirtualServerReq := v1.BSMVirtualServerReq{}
	if err := json.Unmarshal(body, &bsmVirtualServerReq); err != nil {
		logger.Infof("read virtualserver %s err, %v", bsmVirtualServerReq.Name, err)
		response.Response(c, response.ErrBindParams, nil)
		return
	}
	virtualserver := v1.TransformBSMVirtualServerReqToVirtualServer(bsmVirtualServerReq)
	if err := virtualserver.CheckParam(); err != nil {
		logger.Infof("check virtualserver %s request param err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	// 2. 初始化数据
	if err := virtualserver.Init(username, w.virtualServerManager.Kubevirt, w.virtualServerManager.Kubecli, w.virtualServerManager.Repo); err != nil {
		logger.Infof("check virtualserver %s request param err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	// 3. 初始化资源, bsm场景下，产品管理在产品服务，子产品接到请求即认为库存充足，使用产品配置创建实例
	// products := w.controller.ComponentConfig().GetProducts()
	if err := virtualserver.InitResource(w.virtualServerManager.Repo, bsmVirtualServerReq.ProductConfig, nil); err != nil {
		logger.Infof("check virtualserver %s request param err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	// 4. 将用户数据初始化到数据库，异步发起创建kubevirt虚拟机流程
	if err := w.virtualServerManager.AddVirtualServer(virtualserver); err != nil {
		logger.Infof("add virtualserver %s err, %v", virtualserver.Name, err)
		response.Response(c, response.ErrStoreVirtualServer, nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, nil)
}

// ListVirtualServers godoc
//
//	@Summary		List VirtualServers
//	@Description	List VirtualServers
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.VirtualServers}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualservers [get]
func (w *VirtualServerHandler) ListVirtualServers(c *gcpctx.GCPContext) {
	username := c.GetUesrName()
	deleted := c.Query("deleted")
	// logger: log with username and sessionuuid
	logger := c.Logger()

	logger.Infof("username: %s", username)
	var vmRepos []repo.VirtualServer
	var err error
	if deleted == v1.DELETED_YES {
		vmRepos, err = w.virtualServerManager.Repo.ListDeletedVms(username)
	} else {
		vmRepos, err = w.virtualServerManager.Repo.List(username)
	}
	if err != nil {
		logger.Infof("list vm from db err, %v", err)
		response.Response(c, response.ErrListVirtualServer, nil)
		return
	}
	virtualServers := v1.TransformRepoToVirtualServers(vmRepos)
	virtualServers = w.InitDomainName(virtualServers)
	response.Response(c, response.SuccessGCPResponse, virtualServers)
}

func (w *VirtualServerHandler) InitDomainName(virtualServers []v1.VirtualServer) []v1.VirtualServer {
	newVirtualserver := []v1.VirtualServer{}
	for _, virtualserver := range virtualServers {
		virtualserver.State.DomainName = w.controller.ComponentConfig().GetPlatformDomain()
		newVirtualserver = append(newVirtualserver, virtualserver)
	}
	return newVirtualserver
}

// GetVirtualServer godoc
//
//	@Summary		Get VirtualServer
//	@Description	Get VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.VirtualServer}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualservers/{virtualserver.instanceId} [get]
func (w *VirtualServerHandler) GetVirtualServer(c *gcpctx.GCPContext) {
	username := c.GetUesrName()
	// logger: log with username and sessionuuid
	logger := c.Logger()
	logger.Infof("username: %s", username)
	instanceId := c.Param("instanceId")
	if instanceId == "" {
		logger.Infof("instanceId not fount in param")
		response.Response(c, response.ErrCheckRequestParam, nil)
		return
	}
	virtualServer, err := w.virtualServerManager.Repo.GetVmById("", instanceId)
	if err != nil {
		log.Infof("repo get virtualServer %s, err", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, virtualServer)
}

// DeleteVirtualServer godoc
//
//	@Summary		Delete VirtualServer
//	@Description	Delete VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver/{virtualserver.instanceId} [delete]
func (w *VirtualServerHandler) DeleteVirtualServer(c *gcpctx.GCPContext) {
	username := c.GetUesrName()
	// logger: log with username and sessionuuid
	logger := c.Logger()
	logger.Infof("username: %s", username)
	instanceId := c.Param("instanceId")
	if instanceId == "" {
		logger.Infof("instanceId not fount in param")
		response.Response(c, response.ErrCheckRequestParam, nil)
		return
	}
	virtualServerRepo, err := w.virtualServerManager.Repo.GetVmById("", instanceId)
	if err != nil {
		// 删除时实例不存在直接返回成功
		if strings.Contains(err.Error(), "not found") {
			response.Response(c, response.SuccessGCPResponse, nil)
			return
		}
		log.Infof("repo get virtualServer %s, err", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	logger.Infof("delete vm %s", virtualServerRepo.Name)
	if err := w.virtualServerManager.DeleteVirtualServer(v1.TransformRepoToVirtualServer(virtualServerRepo)); err != nil {
		logger.Infof("delete virtualserver %s err, %v", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, nil)
}

// RestartVirtualServer godoc
//
//	@Summary		Restart VirtualServer
//	@Description	Restart VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver/{virtualserver.instanceId}/restart [get]
func (w *VirtualServerHandler) RestartVirtualServer(c *gcpctx.GCPContext) {
	// username := c.GetUesrName()

	// logger: log with username and sessionuuid
	logger := c.Logger()

	instanceId := c.Param("instanceId")
	if instanceId == "" {
		logger.Infof("instanceId not fount in param")
		response.Response(c, response.ErrCheckRequestParam, nil)
		return
	}
	virtualServerRepo, err := w.virtualServerManager.Repo.GetVmById("", instanceId)
	if err != nil {
		log.Infof("repo get virtualServer %s, err", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	f := w.virtualServerManager.RestartVirtualServer
	if virtualServerRepo.Status == v1.STATE_STOPPED {
		// vm为停止状态，kubevirt不支持重启，则使用启动指令
		f = w.virtualServerManager.StartKubevirtVm
	}
	if err := f(v1.TransformRepoToVirtualServer(virtualServerRepo)); err != nil {
		logger.Infof("delete virtualserver %s err, %v", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, nil)
}

// StopVirtualServer godoc
//
//	@Summary		StopVirtualServer
//	@Description	Stop VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver/{virtualserver.instanceId}/stop [get]
func (w *VirtualServerHandler) StopVirtualServer(c *gcpctx.GCPContext) {
	// username := c.GetUesrName()

	// logger: log with username and sessionuuid
	logger := c.Logger()

	instanceId := c.Param("instanceId")
	if instanceId == "" {
		logger.Infof("instanceId not fount in param")
		response.Response(c, response.ErrCheckRequestParam, nil)
		return
	}
	virtualServerRepo, err := w.virtualServerManager.Repo.GetVmById("", instanceId)
	if err != nil {
		log.Infof("repo get virtualServer %s, err", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}

	if err := w.virtualServerManager.StopVirtualServer(v1.TransformRepoToVirtualServer(virtualServerRepo)); err != nil {
		logger.Infof("delete virtualserver %s err, %v", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, nil)
}

// StartVirtualServer godoc
//
//	@Summary		Start VirtualServer
//	@Description	Start VirtualServer
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver/{virtualserver.instanceId}/start [get]
func (w *VirtualServerHandler) StartVirtualServer(c *gcpctx.GCPContext) {
	// username := c.GetUesrName()

	// logger: log with username and sessionuuid
	logger := c.Logger()

	instanceId := c.Param("instanceId")
	if instanceId == "" {
		logger.Infof("instanceId not fount in param")
		response.Response(c, response.ErrCheckRequestParam, nil)
		return
	}
	virtualServerRepo, err := w.virtualServerManager.Repo.GetVmById("", instanceId)
	if err != nil {
		log.Infof("repo get virtualServer %s, err", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}

	if err := w.virtualServerManager.StartVirtualServer(v1.TransformRepoToVirtualServer(virtualServerRepo)); err != nil {
		logger.Infof("delete virtualserver %s err, %v", instanceId, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, nil)
}

// AddVolume godoc
//
//	@Summary		Bind Volume
//	@Description	Bind Volume
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver/{virtualserver.instanceId}/bind/storage/{volume.name} [post]
func (w *VirtualServerHandler) AddVolume(c *gcpctx.GCPContext) {
	username := c.GetUesrName()

	// logger: log with username and sessionuuid
	logger := c.Logger()

	logger.Infof("username: %s", username)
	response.Response(c, response.SuccessGCPResponse, nil)
}

// RemoveVolume godoc
//
//	@Summary		Remove Volume
//	@Description	Remove Volume
//	@Tags			virtualserver
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/virtualserver/{virtualserver.instanceId}/unbind/storage/{volume.name} [post]
func (w *VirtualServerHandler) RemoveVolume(c *gcpctx.GCPContext) {
	username := c.GetUesrName()

	// logger: log with username and sessionuuid
	logger := c.Logger()

	logger.Infof("username: %s", username)
	response.Response(c, response.SuccessGCPResponse, nil)
}
