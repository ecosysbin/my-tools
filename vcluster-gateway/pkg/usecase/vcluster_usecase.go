package usecase

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/datasource"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/utils"
	forkvcluster "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/vcluster"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/vcluster/lifecycle"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/repository"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/consts"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/processor"
)

type cachePrefix string

const (
	// e.g. kubeconfig:vClusterId, vc_name: matrix01, instance_id: 123456789
	kubeconfigCachePrefix cachePrefix = "kubeconfig:"
	vcNameCachePrefix     cachePrefix = "vc_name:"
	instanceIdCachePrefix cachePrefix = "instance_id:"
)

func (cp cachePrefix) combineKey(value string) string {
	return string(cp) + value
}

type VClusterUseCase struct {
	Repo repository.VClusterRepository

	cache            *ristretto.Cache // 主要用于缓存，类似于 redis
	idempotencyCache *ristretto.Cache // 主要用于处理幂等，例如重复 name、instance_id 的请求发起时，在操作数据库前返回
}

func NewVClusterUseCase(repo repository.VClusterRepository) *VClusterUseCase {
	vcuc := &VClusterUseCase{
		Repo: repo,
	}

	commonConfig := &ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	}

	cache, err := ristretto.NewCache(commonConfig)
	if err != nil {
		panic(err)
	}

	idempotencyCache, err := ristretto.NewCache(commonConfig)
	if err != nil {
		panic(err)
	}

	vcuc.cache = cache
	vcuc.idempotencyCache = idempotencyCache

	return vcuc
}

func (vcuc *VClusterUseCase) CreateVCluster(ctx context.Context, info *v1.VClusterInfo) (v1.CreateVClusterResponse, error) {
	var resp v1.CreateVClusterResponse
	var err error

	// info.Name = info.OrderDetails.Name
	// info.Desc = info.OrderDetails.Desc

	// 检查 instance_id 幂等性
	// idKey := string(instanceIdCachePrefix) + info.InstanceId
	idKey := instanceIdCachePrefix.combineKey(info.InstanceId)
	if _, found := vcuc.idempotencyCache.Get(idKey); found {
		// 缓存中存在
		return resp, errors.Wrapf(v1.ErrorInstanceIdAlreadyExists, "failed to CreateVCluster, instance_id: %s is already in process", info.InstanceId)
	} else {
		// 缓存中不存在，从数据库中查询
		if exist := vcuc.Repo.CheckInstanceIdExist(info.InstanceId); exist {
			vcuc.idempotencyCache.Set(idKey, struct{}{}, 1)
			vcuc.cache.Wait()
			return resp, errors.Wrapf(v1.ErrorInstanceIdAlreadyExists, "failed to CreateVCluster, instance_id: %s is already in process", info.InstanceId)
		}
	}

	// 检查 name 幂等性
	//nameKey := vcNameCachePrefix.combineKey(info.Name)
	//if _, found := vcuc.idempotencyCache.Get(nameKey); found {
	//	// 缓存中存在
	//	return resp, errors.Wrapf(v1.ErrorVClusterNameAlreadyExists, "failed to CreateVCluster, vcluster name: %s is already in process", info.Name)
	//} else {
	//	if exist := vcuc.Repo.CheckVClusterNameExist(info.Name); exist {
	//		// 缓存中不存在，从数据库中查询
	//		vcuc.idempotencyCache.Set(nameKey, struct{}{}, 1)
	//		vcuc.cache.Wait()
	//		return resp, errors.Wrapf(v1.ErrorVClusterNameAlreadyExists, "failed to CreateVCluster, vcluster name: %s is already in process", info.Name)
	//	}
	//}
	//vcuc.idempotencyCache.Set(nameKey, struct{}{}, 1)

	// 将 instance_id 添加到缓存中
	vcuc.idempotencyCache.Set(idKey, struct{}{}, 1)
	vcuc.idempotencyCache.Wait()

	if exist := vcuc.Repo.CheckVClusterNameExistByTenantId(info.TenantId, info.Name); exist {
		return resp, errors.Wrapf(v1.ErrorVClusterNameAlreadyExists, "failed to CreateVCluster, vcluster name: %s is already in process", info.Name)
	}

	ctx = context.Background()
	info.Id, err = vcuc.Repo.GenerateUniqueVClusterId(ctx)
	if err != nil {
		return resp, errors.Wrapf(err, "failed to generate vcluster id")
	}

	// 数据先入库
	err = vcuc.Repo.CreateVClusterRecord(info.Id, consts.VClusterStatusCreating, info)
	if err != nil {
		return resp, errors.Wrapf(err, "failed to createOrUpgrade vcluster db record")
	}

	go func() {
		var errInGoRoutine error

		defer func() {
			if errInGoRoutine != nil {
				vc, err := vcuc.Repo.GetVClusterById(info.Id)
				if err == nil && vc != nil {
					vc.ServerStatus = consts.VClusterStatusFailed
					vc.Reason = errInGoRoutine.Error()
					if err := vcuc.Repo.UpdateVCluster(vc); err != nil {
						log.Errorf("failed to update vcluster status to failed, err: %v", err)
					}
				}
			}
		}()

		// 真正创建集群
		resp.VClusterId, errInGoRoutine = vcuc.createOrUpgrade(ctx, info)
		if errInGoRoutine != nil {
			log.Errorf("failed to createOrUpgrade vcluster, err: %v", err)
			return
		}
	}()

	resp.VClusterId = info.Id
	return resp, nil
}

func (vcuc *VClusterUseCase) UpdateVCluster(ctx context.Context, info *v1.VClusterInfo) (v1.CreateVClusterResponse, error) {
	var resp v1.UpdateVClusterResponse
	var err error

	// info.Name = info.OrderDetails.Name
	// info.Desc = info.OrderDetails.Desc

	vCluster, err := vcuc.Repo.GetVClusterById(info.VClusterId)
	if err != nil {
		return resp, errors.Errorf("failed to update vcluster, vclusterId: %s is not exist", info.VClusterId)
	}

	if vCluster.InstanceId == "" {
		return resp, errors.Errorf("failed to update vcluster, vclusterId: %s is not exist", info.VClusterId)
	}
	if vCluster.IsDeleted == 1 {
		return resp, errors.Errorf("failed to update vcluster, vclusterId: %s is deleted", info.VClusterId)
	}
	if vCluster.Status != consts.VClusterStatusRunning {
		return resp, errors.Errorf("failed to update vcluster, vclusterId: %s is not running", info.VClusterId)
	}

	// 修改状态为 Updating
	vCluster.ServerStatus = consts.VClusterStatusUpdating
	err = vcuc.Repo.UpdateVCluster(vCluster)
	if err != nil {
		return resp, errors.Wrapf(err, "failed to update vcluster db record")
	}

	go func() {
		var errInGoRoutine error

		defer func() {
			if errInGoRoutine != nil {
				vCluster.Status = consts.VClusterStatusFailed
				vCluster.Reason = errInGoRoutine.Error()
				if err = vcuc.Repo.UpdateVCluster(vCluster); err != nil {
					log.Errorf("failed to update vcluster db record: %v", err)
				}
			}
		}()

		// 更新集群
		_, errInGoRoutine = vcuc.createOrUpgrade(ctx, info)
		if errInGoRoutine != nil {
			log.Errorf("failed to createOrUpgrade vcluster: %v", err)
			return
		}
	}()

	resp.VClusterId = info.VClusterId
	return resp, nil
}

func (vcuc *VClusterUseCase) createOrUpgrade(ctx context.Context, info *v1.VClusterInfo) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	vclusterLogger := log.New()

	// 创建一个用于处理 helm values.yaml 值的 processor
	process := processor.NewHelmValuesProcessor()
	process.SetProcessorLogger(vclusterLogger)

	// 注册 set 函数，每个函数都会对 values.yaml 中一个属性进行赋值
	process.RegisterDefaultSetters()

	// 注册预处理函数，获取 ManagerBy 参数
	process.RegisterPreprocessors(func(info *v1.VClusterInfo) {
		for _, spec := range info.OrderDetails.Orders[0].InstanceSpecs {
			if strings.HasPrefix(spec.ResourceSpecParamCode, consts.ResourceQuotaManagedBy) {
				info.ManagerBy = spec.ParamValue
				break
			}
		}
	})

	// 执行注册的 set 函数
	err := process.ApplySetters(info)
	if err != nil {
		return "", err
	}

	rootK8sConfig, err := vcuc.Repo.GetRootK8sConfig()
	if err != nil {
		return "", err
	}

	// 创建一个执行器，用于执行 helm createOrUpgrade/upgrade 命令
	globalFlags := forkvcluster.NewDefaultGlobalFlags()
	globalFlags.SetK8sContext(datasource.DefaultCluster)

	createHelm := forkvcluster.NewCreateHelm(
		forkvcluster.WithGlobalFlags(globalFlags),
		forkvcluster.WithCreateOptions(forkvcluster.NewDefaultCreateOptions()),
		forkvcluster.WithCreateChartRepo(info.ChartRepo),
		forkvcluster.WithCreateK8sConfig(rootK8sConfig),
	)

	vclusterLogger.Infof("createOrUpgrade vcluster, successfully NewCreateHelm, CreateHelm: %v, GlobalFlags: %v, CreateOptions: %v", createHelm, *createHelm.GlobalFlags, *createHelm.CreateOptions)

	nodeSelector := make(map[string]string)
	if !info.IsInit {
		if info.NodePoolInstanceId == "" {
			nodeSelector["dc.com/osm.nodepool.type"] = "share"
		} else {
			nodeSelector["dc.com/osm.nodepool.type"] = "exclusive"
			nodeSelector["dc.com/osm.nodepool.tenantId"] = info.TenantId
			// nodeSelector["dc.com/osm.nodepool.orderInstanceId"] = info.NodePoolInstanceId
		}
	}
	// 初始化执行器
	err = createHelm.ValidateAndInitialize(ctx, info.Id, nodeSelector, info.IsInit)
	if err != nil {
		return "", err
	}

	// 将拼接好的 values 输出到一个文件
	valuesFile, err := process.GenerateValuesFile(info.Id)
	if err != nil {
		return "", errors.Wrapf(err, "failed to GenerateValuesFile, vclusterId: %s", info.Id)
	}

	// 执行操作
	valuesFiles := []string{valuesFile}

	if info.IsInit {
		valuesFiles = append(valuesFiles, consts.ControlPlaneAffinityTolerations.String())
	} else {
		//valuesFiles = append(valuesFiles, consts.EnableStoragePlugin.String())
	}

	if info.EnableHA {
		valuesFiles = append(valuesFiles, consts.EnableHA.String())
	}

	var enableServiceExporterFn = func(info *v1.VClusterInfo) bool {
		for _, spec := range info.OrderDetails.Orders[0].InstanceSpecs {
			if spec.ResourceSpecParamCode == consts.ResourceEnableServiceExporter {
				if spec.ParamValue == "true" {
					return true
				}
			}
		}
		return false
	}

	if enableServiceExporterFn(info) {
		valuesFiles = append(valuesFiles, consts.EnableServiceExporter.String())
	}

	var enableStoragePluginFn = func(info *v1.VClusterInfo) bool {
		for _, order := range info.OrderDetails.Orders {
			if consts.ResourceTypeCodeSingleStorageMap[order.ResourceTypeCode] {
				return true
			}
		}
		return false
	}

	if enableStoragePluginFn(info) {
		valuesFiles = append(valuesFiles, consts.EnableStoragePlugin.String())
	}

	// 如果使用自定义的 helm 配置，那么清空通过请求参数拼接的 values 文件
	if info.CustomHelmConfig.EnableCustomization {
		if info.CustomHelmConfig.ValuesContent != "" {
			// 如果 values 内容不为空，则使用自定义的 values, 清空 valuesFiles 并且将用户传递的 values 写入临时文件
			valuesFiles = make([]string, 0)
			filename, err := processor.WriteToTempFile(&info.CustomHelmConfig.ValuesContent, info.Id, vclusterLogger)
			if err != nil {
				return "", errors.Wrapf(err, "failed to WriteToTempFile")
			}
			valuesFiles = append(valuesFiles, filename)

			err = os.Remove(valuesFile)
			if err != nil {
				vclusterLogger.Warnf("failed to remove values file in custom helm config, err: %v", err)
				err = nil
			}
			vclusterLogger.Infof("remove values file in custom helm config, values file: %s", valuesFile)
		}

		if info.CustomHelmConfig.Repo != "" {
			// 如果 chart repo 不为空，则使用自定义的 chart repo
			createHelm.ChartRepo = info.CustomHelmConfig.Repo
		}
	}

	valuesFiles = append(valuesFiles, consts.Extensions.String())

	vclusterLogger.Infof("valuesFiles: %v", valuesFiles)

	err = createHelm.Deploy(ctx, info.Id, valuesFiles)
	if err != nil {
		return "", err
	}

	return info.Id, err
}

// DeleteVCluster 删除 vCluster
// 1. 先校验权限
// 2. 将 ServerStatus 设置为 "Deleting"
// 3. 调用 APS 删除接口租户
// 4. 删除 vcluster
func (vcuc *VClusterUseCase) DeleteVCluster(ctx context.Context, params *v1.DeleteVClusterParams) (*v1.DeleteVClusterResponse, error) {
	params.Logger.Infof("DeleteVCluster, recieved request, params: %+v", params)

	// 1. 校验权限
	allowFlag := false

	// 检查 TenantType
	if params.TenantType == "3" || params.TenantType == "TENANT_TYPE_PLATFORM" {
		params.Logger.Infof("DeleteVCluster, allowFlag is true, TenantType: %s, TenantId: %s", params.TenantType, params.TenantId)
		allowFlag = true
	} else {
		// 获取 vCluster 并检查 TenantId
		vc, err := vcuc.Repo.GetVClusterById(params.Id)
		if err != nil {
			return &v1.DeleteVClusterResponse{Message: "failed to get vcluster by id"}, err
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return &v1.DeleteVClusterResponse{Message: "tenant id not match"}, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	// 检查 allowFlag 是否为 true
	if !allowFlag {
		params.Logger.Warnf("DeleteVCluster, allowFlag is false, TenantType: %s, TenantId: %s", params.TenantType, params.TenantId)
		return &v1.DeleteVClusterResponse{Message: "tenant type not match"}, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	vc, err := vcuc.Repo.GetVClusterById(params.Id)
	if err != nil {
		return &v1.DeleteVClusterResponse{Message: "failed to get vcluster by id"}, err
	}

	// 如果状态已经为 Deleted，直接返回
	if vc.IsDeleted == 1 || vc.ServerStatus == consts.VClusterStatusDeleted {
		params.Logger.Infof("DeleteVCluster, vcluster is already deleted, vclusterId: %s", params.Id)
		return &v1.DeleteVClusterResponse{Message: "vcluster is deleted"}, nil
	}

	//if vc.ServerStatus == consts.VClusterStatusDeleting {
	//	params.Logger.Infof("DeleteVCluster, vcluster is already deleting, vclusterId: %s", params.Id)
	//	return &v1.DeleteVClusterResponse{Message: "vcluster is deleting"}, nil
	//}

	// 先将状态设置为 Deleting，将 ServerStatus 变为 Deleted，由 K8s Informer 去监听删除事件并设置
	vc.ServerStatus = consts.VClusterStatusDeleting

	err = vcuc.Repo.UpdateVCluster(vc)
	if err != nil {
		return &v1.DeleteVClusterResponse{Message: "failed to update vcluster"}, err
	}

	params.Logger.Infof("DeleteVCluster, set ServerStatus to Deleting, vclusterId: %s", params.Id)

	go func() {
		var errInGoRoutine error

		defer func() {
			if errInGoRoutine != nil {
				vc, err := vcuc.Repo.GetVClusterById(params.Id)
				if err == nil && vc != nil {
					vc.ServerStatus = consts.VClusterStatusFailed
					vc.Reason = errInGoRoutine.Error()
					if err := vcuc.Repo.UpdateVCluster(vc); err != nil {
						params.Logger.Errorf("failed to update vcluster status, err: %v, vclusterId: %s", err, params.Id)
					}
				}
			}
		}()

		// 3. 真正删除 vcluster
		msg, errInGoRoutine := vcuc.delete(ctx, params)
		if errInGoRoutine != nil {
			params.Logger.Warnf("DeleteVCluster, failed to delete vcluster, msg: %s, error: %s", msg, err.Error())
			return
		}

		params.Logger.Infof("DeleteVCluster, delete vcluster successfully, id: %s", params.Id)
	}()

	return &v1.DeleteVClusterResponse{Message: "Deleting"}, nil
}

func (vcuc *VClusterUseCase) delete(ctx context.Context, params *v1.DeleteVClusterParams) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ctx = context.WithValue(ctx, "logger", params.Logger)

	globalFlags := forkvcluster.NewDefaultGlobalFlags()
	globalFlags.SetK8sContext(datasource.DefaultCluster)

	rootK8sConfig, err := vcuc.Repo.GetRootK8sConfig()
	if err != nil {
		return "", err
	}

	deleteHelm := forkvcluster.NewDeleteHelm(
		forkvcluster.WithDeleteGlobalFlags(globalFlags),
		// forkvcluster.WithDeleteLogger(log.New()),
		forkvcluster.WithDeleteK8sConfig(rootK8sConfig),
		forkvcluster.WithDeleteOptions(forkvcluster.NewDefaultDeleteOptions()),
	)

	err = deleteHelm.DeleteHelm(ctx, params.Id)
	if err != nil {
		// 调用 helm 删除 vcluster 失败，尝试直接清理 namespace
		err = vcuc.Repo.DeleteVClusterNamespace(ctx, params.Id)
		if err != nil {
			params.Logger.Warnf("failed to delete vcluster namespace, err: %v, vclusterId: %s", err, params.Id)
		} else {
			params.Logger.Info("delete vcluster successfully, but directly delete namespace")
		}

		// 更改数据库中的 vcluster 状态为 deleted
		vc, err := vcuc.Repo.GetVClusterById(params.Id)
		if err == nil && vc != nil {
			vc.ServerStatus = consts.VClusterStatusDeleted
			if err := vcuc.Repo.UpdateVCluster(vc); err != nil {
				params.Logger.Errorf("failed to update vcluster status, err: %v, vclusterId: %s", err, params.Id)
			}
		}
		return "Delete vcluster successfully", nil
	}

	params.Logger.Info("delete vcluster successfully")

	return "Delete vcluster successfully", nil
}

func (vcuc *VClusterUseCase) PauseVClusters(ctx context.Context, params *v1.PauseVClusterParams) (*v1.PauseVClusterResponse, error) {
	id := params.Id

	// 校验权限
	allowFlag := false

	if params.TenantType == "3" || params.TenantType == "TENANT_TYPE_PLATFORM" {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(id)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", id)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	vc, err := vcuc.Repo.GetVClusterById(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", id)
	}

	if vc.ServerStatus == consts.VClusterStatusPaused {
		params.Logger.Warnf("vcluster already paused, vclusterId: %s", id)
		return &v1.PauseVClusterResponse{Message: "Already paused"}, nil
	}

	if vc.ServerStatus == consts.VClusterStatusPausing {
		params.Logger.Warnf("vcluster is pausing, vclusterId: %s", id)
		return &v1.PauseVClusterResponse{Message: "Pausing"}, nil
	}

	// 否则更新成 Pausing
	vc.ServerStatus = consts.VClusterStatusPausing
	err = vcuc.Repo.UpdateVCluster(vc)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update vcluster, id: %s", id)
	}

	go func() {
		var errInGoRoutine error

		defer func() {
			if errInGoRoutine != nil {
				vc, err := vcuc.Repo.GetVClusterById(id)
				if err == nil && vc != nil {
					vc.ServerStatus = consts.VClusterStatusFailed
					vc.Reason = errInGoRoutine.Error()
					if err := vcuc.Repo.UpdateVCluster(vc); err != nil {
						params.Logger.Errorf("failed to update vcluster status, err: %v, vclusterId: %s", err, id)
					}
				}
			}
		}()

		ctx = context.Background()
		rootK8sConfig, errInGoRoutine := vcuc.Repo.GetRootK8sConfig()
		if errInGoRoutine != nil {
			params.Logger.Errorf("failed to get root k8s config, err: %v, vclusterId: %s", err, id)
			return
		}

		errInGoRoutine = vcuc.pause(ctx, id, rootK8sConfig)
		if errInGoRoutine != nil {
			params.Logger.Errorf("failed to pause vcluster, err: %v, vclusterId: %s", err, id)
			return
		}

		params.Logger.Infof("vcluster paused, vclusterId: %s", id)
	}()

	return &v1.PauseVClusterResponse{Message: "Pausing"}, nil
}

func (vcuc *VClusterUseCase) pause(ctx context.Context, vClusterId string, rootK8sConfig *v1.RootK8sConfig) error {
	exist := vcuc.Repo.CheckVClusterExistById(vClusterId)
	if !exist {
		return errors.Errorf("vcluster %s not exist", vClusterId)
	}

	err := lifecycle.PauseVClusterWithCleanup(ctx,
		rootK8sConfig.RootKubeClientSet,
		vClusterId,
		utils.GetVClusterNamespaceName(vClusterId),
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "vcluster: %s pause failed", vClusterId)
	}

	return nil
}

// ResumeVClusters 恢复集群
// 1. 校验权限
// 2. 设置状态为 Resuming
// 3. 真正去恢复 VCluster
// 4. 调用 APS 恢复接口
func (vcuc *VClusterUseCase) ResumeVClusters(ctx context.Context, params *v1.PauseVClusterParams) (*v1.ResumeVClusterResponse, error) {
	id := params.Id

	// 1. 先校验权限
	allowFlag := false

	if params.TenantType == "3" || params.TenantType == "TENANT_TYPE_PLATFORM" {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(params.Id)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", id)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	vc, err := vcuc.Repo.GetVClusterById(params.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", id)
	}

	if vc.ServerStatus == consts.VClusterStatusRunning {
		params.Logger.Infof("vcluster is already running, id: %s", params.Id)
		return &v1.ResumeVClusterResponse{Message: "VCluster is already resumed"}, nil
	}

	if vc.ServerStatus == consts.VClusterStatusResuming {
		params.Logger.Infof("vcluster is already resuming, id: %s", params.Id)
		return &v1.ResumeVClusterResponse{Message: "VCluster is already resuming"}, nil
	}

	// 否则设置为 Resuming
	vc.ServerStatus = consts.VClusterStatusResuming
	err = vcuc.Repo.UpdateVCluster(vc)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update vcluster, id: %s", params.Id)
	}

	params.Logger.Infof("success to update vcluster status to Resuming, id: %s", params.Id)

	go func() {
		var errInGoRoutine error

		defer func() {
			if errInGoRoutine != nil {
				vc, err := vcuc.Repo.GetVClusterById(id)
				if err != nil || vc == nil {
					params.Logger.Errorf("failed to get vcluster by id, err: %v, vclusterId: %s", err, id)
					return
				}
				vc.ServerStatus = consts.VClusterStatusFailed
				vc.Reason = errInGoRoutine.Error()
				if err := vcuc.Repo.UpdateVCluster(vc); err != nil {
					params.Logger.Errorf("failed to update vcluster status, err: %v, vclusterId: %s", err, id)
				}
			}
		}()

		ctx = context.WithValue(ctx, "logger", params.Logger)

		rootK8sConfig, errInGoRoutine := vcuc.Repo.GetRootK8sConfig()
		if err != nil {
			params.Logger.Errorf("failed to get root k8s config, err: %v, vclusterId: %s", err, params.Id)
			return
		}

		errInGoRoutine = vcuc.resume(ctx, params.Id, rootK8sConfig)
		if err != nil {
			params.Logger.Errorf("failed to resume vcluster, err: %v, vclusterId: %s", err, params.Id)
			return
		}
	}()

	return &v1.ResumeVClusterResponse{Message: "Resuming"}, nil
}

func (vcuc *VClusterUseCase) resume(ctx context.Context, vClusterId string, rootK8sConfig *v1.RootK8sConfig) error {
	exist := vcuc.Repo.CheckVClusterExistById(vClusterId)
	if !exist {
		return errors.Errorf("vcluster %s not exist", vClusterId)
	}

	err := lifecycle.ResumeVCluster(ctx,
		rootK8sConfig.RootKubeClientSet,
		vClusterId,
		utils.GetVClusterNamespaceName(vClusterId),
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "vcluster: %s resume failed", vClusterId)
	}

	return nil
}

func (vcuc *VClusterUseCase) GetKubeConfig(ctx context.Context, params *v1.GetKubeConfigParams) (*v1.GetVClusterTokenResponse, error) {
	// key := string(kubeconfigCachePrefix) + params.VClusterId
	key := kubeconfigCachePrefix.combineKey(params.VClusterId)

	obj, found := vcuc.cache.Get(key)
	if found {
		params.Logger.Infof("successfully got vCluster kubeconfig from cache for vclusterId: %s", params.VClusterId)
		return obj.(*v1.GetVClusterTokenResponse), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	allowFlag := false

	if params.TenantType == "3" || params.TenantType == "TENANT_TYPE_PLATFORM" {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(params.VClusterId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.VClusterId)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	ctx = context.WithValue(ctx, "logger", params.Logger)

	kubeconfig, err := vcuc.Repo.GetKubeConfig(ctx, params.VClusterId, params.KubeConnHost)
	if err != nil {
		params.Logger.Errorf("Failed to get vCluster token for tenantId: %s, vclusterId: %s, error: %v", params.TenantId, params.VClusterId, err)
		return nil, errors.Wrapf(err, "failed to get vcluster token, tenantId: %s, vclusterId: %s", params.TenantId, params.VClusterId)
	}

	vcuc.cache.Set(key, kubeconfig, 1)
	vcuc.cache.Wait()

	params.Logger.Infof("Successfully got vCluster token for tenantId: %s, vclusterId: %s", params.TenantId, params.VClusterId)

	return kubeconfig, nil
}

func (vcuc *VClusterUseCase) QueryOperateStatus(ctx context.Context, params *v1.QueryOperateStatusRequest) (*vclusterv1.QueryOperateStatusResponse_Data, error) {
	// 1. 校验权限
	var allowFlag bool

	if params.TenantType == consts.TenantTypePlatformAlias || params.TenantType == consts.TenantTypePlatform {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(params.AppId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.AppId)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	// 2. 调用数据库查询状态
	vc, err := vcuc.Repo.GetVClusterById(params.AppId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.AppId)
	}

	// 3. 判断操作类型返回状态
	retStatus := consts.StatusProcessing
	switch params.Action {
	case consts.ActionCreate, consts.ActionUpdate, consts.ActionResume:
		if vc.ServerStatus == consts.VClusterStatusRunning {
			retStatus = consts.StatusSuccess
		}
	case consts.ActionDelete:
		if vc.ServerStatus == consts.VClusterStatusDeleted {
			retStatus = consts.StatusSuccess
		}
	case consts.ActionPause:
		if vc.ServerStatus == consts.VClusterStatusPaused {
			retStatus = consts.StatusSuccess
		}
	}

	// 4. 返回响应数据
	return &vclusterv1.QueryOperateStatusResponse_Data{
		AppId:  params.AppId,
		Status: retStatus,
		Reason: vc.Reason,
	}, nil
}

func (vcuc *VClusterUseCase) GetVClusterStatus(ctx context.Context, params *v1.GetVClusterStatusRequest) (*vclusterv1.GetVClusterStatusResponse_Data, error) {
	var allowFlag bool

	if params.TenantType == consts.TenantTypePlatformAlias || params.TenantType == consts.TenantTypePlatform {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(params.AppId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.AppId)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	vc, err := vcuc.Repo.GetVClusterById(params.AppId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.AppId)
	}

	return &vclusterv1.GetVClusterStatusResponse_Data{
		Status: vc.ServerStatus,
	}, nil
}

func (vcuc *VClusterUseCase) GetVClusterResourceDetails(ctx context.Context, params *v1.GetVClusterResourceDetailsRequest) (*vclusterv1.GetVClusterResourceDetailsResponse_Data, error) {
	var allowFlag bool

	if params.TenantType == consts.TenantTypePlatformAlias || params.TenantType == consts.TenantTypePlatform {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(params.AppId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.AppId)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	data, err := vcuc.Repo.GetVClusterResourceDetails(ctx, params.AppId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster resource quotas by id: %s", params.AppId)
	}

	return &vclusterv1.GetVClusterResourceDetailsResponse_Data{
		UtilizationRate: data.UtilizationRate,
		Configurations:  data.Configurations,
	}, nil
}

func (vcuc *VClusterUseCase) GetVClusterContainerID(ctx context.Context, params *v1.GetVClusterContainerIDRequest) (*vclusterv1.GetVClusterContainerIDResponse_Data, error) {
	var allowFlag bool

	if params.TenantType == consts.TenantTypePlatformAlias || params.TenantType == consts.TenantTypePlatform {
		allowFlag = true
	} else {
		vc, err := vcuc.Repo.GetVClusterById(params.VClusterId)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get vcluster by id: %s", params.VClusterId)
		}

		if vc.TenantId == params.TenantId {
			allowFlag = true
		} else {
			return nil, errors.Errorf("tenant id not match, vcluster tenant id: %s, request tenant id: %s", vc.TenantId, params.TenantId)
		}
	}

	if !allowFlag {
		return nil, errors.Errorf("tenant type not match, request tenant type: %s", params.TenantType)
	}

	containerId, err := vcuc.Repo.GetVClusterContainerID(ctx, params)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster container id: %s", params.VClusterId)
	}

	return &vclusterv1.GetVClusterContainerIDResponse_Data{
		ContainerId: containerId,
	}, nil

}
