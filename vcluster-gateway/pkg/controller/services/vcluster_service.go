package services

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"connectrpc.com/connect"

	v1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
	"vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1/vclusterv1connect"
	"vcluster-gateway/pkg/controller/framework"
	"vcluster-gateway/version"
)

type VClusterServer struct {
	controller framework.Interface

	vclusterv1.UnimplementedVClusterGatewayServiceServer
	vclusterv1connect.UnimplementedVClusterGatewayServiceHandler
}

func NewVClusterServer(controller framework.Interface) *VClusterServer {
	return &VClusterServer{
		controller: controller,
	}
}

// CheckHealth 健康检查
// HTTP:
// GET /v1/health/check
func (vcs *VClusterServer) CheckHealth(ctx context.Context, req *connect.Request[vclusterv1.CheckHealthRequest]) (*connect.Response[vclusterv1.CheckHealthResponse], error) {
	return connect.NewResponse(&vclusterv1.CheckHealthResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
	}), nil
}

// VersionInformation vCluster-gateway 服务的版本信息
// HTTP:
// GET /v1/version
func (vcs *VClusterServer) VersionInformation(ctx context.Context, req *connect.Request[vclusterv1.VersionInformationRequest]) (*connect.Response[vclusterv1.VersionInformationResponse], error) {
	info := &vclusterv1.VersionInformationResponse_Info{
		Version:   version.Version,
		GitCommit: version.GitCommit,
		BuildAt:   version.BuildAt,
	}

	return connect.NewResponse(&vclusterv1.VersionInformationResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: info,
	}), nil
}

// CreateVCluster 创建 vCluster
// HTTP:
// POST /v1/app
func (vcs *VClusterServer) CreateVCluster(ctx context.Context, req *connect.Request[vclusterv1.CreateVClusterRequest]) (*connect.Response[vclusterv1.CreateVClusterResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	formatInstanceSpecs := func(before []*vclusterv1.CreateVClusterRequest_Orderlist_InstanceSpec) (after []*v1.InstanceSpec) {
		after = make([]*v1.InstanceSpec, 0, len(before))
		for _, instanceSpec := range before {
			after = append(after, &v1.InstanceSpec{
				ResourceSpecId:        int(instanceSpec.ResourceSpecId),
				ResourceSpecCode:      instanceSpec.ResourceSpecCode,
				ResourceSpecParamId:   int(instanceSpec.ResourceSpecParamId),
				ResourceSpecParamCode: instanceSpec.ResourceSpecParamCode,
				ParamName:             instanceSpec.ParamName,
				ParamValue:            instanceSpec.ParamValue,
			})
		}
		return after
	}
	formatOrderList := func(before []*vclusterv1.CreateVClusterRequest_Orderlist) (after []*v1.Order) {
		after = make([]*v1.Order, 0, len(before))
		for _, order := range before {
			after = append(after, &v1.Order{
				ProductID:          int(order.OrderInfo.ProductId),
				CycleCount:         int(order.OrderInfo.CycleCount),
				Amount:             int(order.OrderInfo.Amount),
				ActualAmount:       int(order.OrderInfo.ActualAmount),
				OrderType:          int(order.OrderInfo.OrderType),
				ProductCode:        order.OrderInfo.ProductCode,
				ResourceTypeID:     int(order.OrderInfo.ResourceTypeId),
				ResourceTypeCode:   order.OrderInfo.ResourceTypeCode,
				InstanceID:         order.InstanceId,
				NodePoolInstanceId: order.NodePoolInstanceId,
				InstanceSpecs:      formatInstanceSpecs(order.InstanceSpec),
			})
		}
		return after
	}

	var getCustomHelmConfig = func() v1.CustomHelmConfig {
		var customHelmConfig v1.CustomHelmConfig
		if req.Msg.GetCustomHelmConfig() != nil {
			customHelmConfig.EnableCustomization = req.Msg.GetCustomHelmConfig().EnableCustomization
			customHelmConfig.Repo = req.Msg.GetCustomHelmConfig().Repo
			customHelmConfig.ValuesContent = req.Msg.GetCustomHelmConfig().ValuesContent
		}
		return customHelmConfig
	}

	var getFallbackDns = func() string {
		return os.Getenv("FALLBACK_DNS")
	}

	resp, _ := vcs.controller.VClusterController().
		CreateVCluster(&v1.CreateVClusterParams{
			Orders:           formatOrderList(req.Msg.GetOrderlist()),
			Logger:           gcpLogger,
			Name:             req.Msg.GetName(),
			Desc:             req.Msg.GetDesc(),
			UserName:         username,
			TenantId:         tenantId,
			RawName:          req.Msg.GetVclusterName(),
			RawDesc:          req.Msg.GetVclusterDesc(),
			IsInit:           req.Msg.GetIsInit(),
			EnableHA:         req.Msg.EnableHa,
			FallbackDns:      getFallbackDns(),
			CustomHelmConfig: getCustomHelmConfig(),
		})

	return connect.NewResponse(resp), nil
}

// UpdateVCluster 更新 vCluster
// HTTP:
// PUT /v1/app/{app_id}
func (vcs *VClusterServer) UpdateVCluster(ctx context.Context, req *connect.Request[vclusterv1.UpdateVClusterRequest]) (*connect.Response[vclusterv1.UpdateVClusterResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()
	appId := req.Msg.AppId

	formatInstanceSpecs := func(before []*vclusterv1.UpdateVClusterRequest_Orderlist_InstanceSpec) (after []*v1.InstanceSpec) {
		after = make([]*v1.InstanceSpec, 0, len(before))
		for _, instanceSpec := range before {
			after = append(after, &v1.InstanceSpec{
				ResourceSpecId:        int(instanceSpec.ResourceSpecId),
				ResourceSpecCode:      instanceSpec.ResourceSpecCode,
				ResourceSpecParamId:   int(instanceSpec.ResourceSpecParamId),
				ResourceSpecParamCode: instanceSpec.ResourceSpecParamCode,
				ParamName:             instanceSpec.ParamName,
				ParamValue:            instanceSpec.ParamValue,
			})
		}
		return after
	}
	formatOrderList := func(before []*vclusterv1.UpdateVClusterRequest_Orderlist) (after []*v1.Order) {
		after = make([]*v1.Order, 0, len(before))
		for _, order := range before {
			after = append(after, &v1.Order{
				ProductID:          int(order.OrderInfo.ProductId),
				CycleCount:         int(order.OrderInfo.CycleCount),
				Amount:             int(order.OrderInfo.Amount),
				ActualAmount:       int(order.OrderInfo.ActualAmount),
				OrderType:          int(order.OrderInfo.OrderType),
				ProductCode:        order.OrderInfo.ProductCode,
				ResourceTypeID:     int(order.OrderInfo.ResourceTypeId),
				ResourceTypeCode:   order.OrderInfo.ResourceTypeCode,
				InstanceID:         order.InstanceId,
				NodePoolInstanceId: order.NodePoolInstanceId,
				InstanceSpecs:      formatInstanceSpecs(order.InstanceSpec),
			})
		}
		return after
	}

	gcpLogger.Infof("UpdateVCluster, recieve request: %v, appId: %s", req.Msg, appId)

	var getCustomHelmConfig = func() v1.CustomHelmConfig {
		var customHelmConfig v1.CustomHelmConfig
		if req.Msg.GetCustomHelmConfig() != nil {
			customHelmConfig.EnableCustomization = req.Msg.GetCustomHelmConfig().EnableCustomization
			customHelmConfig.Repo = req.Msg.GetCustomHelmConfig().Repo
			customHelmConfig.ValuesContent = req.Msg.GetCustomHelmConfig().ValuesContent
		}
		return customHelmConfig
	}

	var getFallbackDns = func() string {
		return os.Getenv("FALLBACK_DNS")
	}

	resp, _ := vcs.controller.VClusterController().
		UpdateVCluster(&v1.UpdateVClusterParams{
			AppId:            appId,
			Orders:           formatOrderList(req.Msg.GetOrderlist()),
			Logger:           gcpLogger,
			Name:             req.Msg.GetName(),
			Desc:             req.Msg.GetDesc(),
			UserName:         username,
			TenantId:         tenantId,
			RawName:          req.Msg.GetVclusterName(),
			RawDesc:          req.Msg.GetVclusterDesc(),
			IsInit:           req.Msg.GetIsInit(),
			EnableHA:         req.Msg.EnableHa,
			FallbackDns:      getFallbackDns(),
			CustomHelmConfig: getCustomHelmConfig(),
		})

	return connect.NewResponse(resp), nil
}

// DeleteVCluster 删除 vCluster
// HTTP:
// DELETE /v1/app/{app_id}
func (vcs *VClusterServer) DeleteVCluster(ctx context.Context, req *connect.Request[vclusterv1.DeleteVClusterRequest]) (*connect.Response[vclusterv1.DeleteVClusterResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId := req.Msg.AppId
	resp, _ := vcs.controller.VClusterController().
		DeleteVCluster(&v1.DeleteVClusterParams{
			Logger:     gcpLogger,
			Id:         appId,
			Username:   username,
			TenantId:   tenantId,
			TenantType: "",
		})

	return connect.NewResponse(resp), nil
}

// GetKubeConfig 获取集群的 kubeconfig
// HTTP:
// GET /v1/kubeconfig/{app_id}
func (vcs *VClusterServer) GetKubeConfig(ctx context.Context, req *connect.Request[vclusterv1.GetKubeConfigRequest]) (*connect.Response[vclusterv1.GetKubeConfigResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId := req.Msg.AppId

	resp, _ := vcs.controller.VClusterController().
		GetKubeConfig(&v1.GetKubeConfigParams{
			Logger:     gcpLogger,
			VClusterId: appId,
			Username:   username,
			TenantId:   tenantId,
			TenantType: "",
		})

	return connect.NewResponse(resp), nil
}

// GetKubeConfigBase64 获取集群的 kubeconfig
// HTTP:
// GET /v1/kubeconfig/base64/{app_id}
func (vcs *VClusterServer) GetKubeConfigBase64(ctx context.Context, req *connect.Request[vclusterv1.GetKubeConfigRequest]) (*connect.Response[vclusterv1.GetKubeConfigBase64Response], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId := req.Msg.AppId

	resp, _ := vcs.controller.VClusterController().
		GetKubeConfigBase64(&v1.GetKubeConfigParams{
			Logger:     gcpLogger,
			VClusterId: appId,
			Username:   username,
			TenantId:   tenantId,
			TenantType: "",
		})

	return connect.NewResponse(resp), nil
}

// PauseVCluster 暂停 vCluster
// HTTP:
// POST /v1/app/{app_id}/pause
func (vcs *VClusterServer) PauseVCluster(ctx context.Context, req *connect.Request[vclusterv1.PauseVClusterRequest]) (*connect.Response[vclusterv1.PauseVClusterResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()
	appId := req.Msg.AppId

	resp, _ := vcs.controller.VClusterController().
		PauseVCluster(&v1.PauseVClusterParams{
			Logger:     gcpLogger,
			Username:   username,
			Id:         appId,
			TenantId:   tenantId,
			TenantType: "",
		})

	return connect.NewResponse(resp), nil
}

// ResumeVCluster 恢复 vCluster
// HTTP:
// POST /v1/app/{app_id}/recover
func (vcs *VClusterServer) ResumeVCluster(ctx context.Context, req *connect.Request[vclusterv1.ResumeVClusterRequest]) (*connect.Response[vclusterv1.ResumeVClusterResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId := req.Msg.AppId

	resp, _ := vcs.controller.VClusterController().
		ResumeVCluster(&v1.ResumeVClusterParams{
			Logger:     gcpLogger,
			Username:   username,
			Id:         appId,
			TenantId:   tenantId,
			TenantType: "",
		})

	return connect.NewResponse(resp), nil
}

// QueryOperateStatus 查询 操作 vCluster 的状态
// HTTP:
// GET /v1/app/{app_id}?type={create|update|delete|stop|recover}
// NOTE: app-gateway 通过该接口查询状态
func (vcs *VClusterServer) QueryOperateStatus(ctx context.Context, req *connect.Request[vclusterv1.QueryOperateStatusRequest]) (*connect.Response[vclusterv1.QueryOperateStatusResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId, operateType := req.Msg.AppId, req.Msg.Action

	resp, _ := vcs.controller.VClusterController().
		QueryOperateStatus(&v1.QueryOperateStatusRequest{
			Logger:     gcpLogger,
			Username:   username,
			TenantId:   tenantId,
			TenantType: "",

			AppId:  appId,
			Action: operateType,
		})

	return connect.NewResponse(resp), nil
}

// GetVClusterStatus 获取 vCluster 的状态
// HTTP:
// GET /v1/app/{app_id}/status
func (vcs *VClusterServer) GetVClusterStatus(ctx context.Context, req *connect.Request[vclusterv1.GetVClusterStatusRequest]) (*connect.Response[vclusterv1.GetVClusterStatusResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId := req.Msg.AppId

	resp, _ := vcs.controller.VClusterController().
		GetVClusterStatus(&v1.GetVClusterStatusRequest{
			Logger:     gcpLogger,
			Username:   username,
			TenantId:   tenantId,
			TenantType: "",

			AppId: appId,
		})

	return connect.NewResponse(resp), nil
}

// GetVClusterResourceDetails 获取 vCluster 集群资源详细信息，包括 ResourceQuota 和 Configurations（一些配置信息，例如是否开启 service export）
// HTTP:
// GET /v1/app/{app_id}/resourcedetails
func (vcs *VClusterServer) GetVClusterResourceDetails(ctx context.Context, req *connect.Request[vclusterv1.GetVClusterResourceDetailsRequest]) (*connect.Response[vclusterv1.GetVClusterResourceDetailsResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	appId := req.Msg.AppId

	resp, _ := vcs.controller.VClusterController().
		GetVClusterResourceDetails(&v1.GetVClusterResourceDetailsRequest{
			Logger:     gcpLogger,
			Username:   username,
			TenantId:   tenantId,
			TenantType: "",

			AppId: appId,
		})

	return connect.NewResponse(resp), nil
}

// GetVClusterContainerID 获取 vCluster 集群中的容器 ID
// HTTP:
// GET /v1/vclusters/{vcluster_id}/namespaces/{namespace}/pods/{pod_name}/containers/{container_name}/id
func (vcs *VClusterServer) GetVClusterContainerID(ctx context.Context, req *connect.Request[vclusterv1.GetVClusterContainerIDRequest]) (*connect.Response[vclusterv1.GetVClusterContainerIDResponse], error) {
	username, tenantId := "", ""
	gcpLogger := log.New()

	resp, _ := vcs.controller.VClusterController().
		GetVClusterContainerID(&v1.GetVClusterContainerIDRequest{
			Logger:        gcpLogger,
			Username:      username,
			TenantId:      tenantId,
			TenantType:    "",
			VClusterId:    req.Msg.VclusterId,
			Namespace:     req.Msg.Namespace,
			PodName:       req.Msg.PodName,
			ContainerName: req.Msg.ContainerName,
		})

	return connect.NewResponse(resp), nil
}
