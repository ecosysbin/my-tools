package vcluster

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase"
)

type VClusterController struct {
	controller framework.Interface
	useCase    *usecase.VClusterUseCase
}

func NewVClusterController(controller framework.Interface) *VClusterController {
	vcc := &VClusterController{
		controller: controller,
	}

	err := vcc.controller.DIContainer().Invoke(func(useCase *usecase.VClusterUseCase) {
		vcc.useCase = useCase
	})
	if err != nil {
		panic(err)
	}

	return vcc
}

func (vcc *VClusterController) CreateVCluster(params *v1.CreateVClusterParams) (*vclusterv1.CreateVClusterResponse, error) {
	if len(params.Orders) == 0 {
		msg := "orderList is empty"
		return &vclusterv1.CreateVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_CREATE_CLUSTER_ERROR),
			Msg:  msg,
			Data: nil,
		}, errors.New(msg)
	}

	info := &v1.VClusterInfo{
		FallbackDns:          params.FallbackDns,
		CustomHelmConfig:     params.CustomHelmConfig,
		EnableHA:             params.EnableHA,
		IsInit:               params.IsInit,
		Username:             params.UserName,
		TenantId:             params.TenantId,
		Name:                 params.RawName,
		Comment:              params.RawDesc,
		InstanceId:           params.Orders[0].InstanceID,
		NodePoolInstanceId:   params.Orders[0].NodePoolInstanceId,
		Upgrade:              false,
		CephClusterId:        vcc.controller.ComponentConfig().GetCephClusterId(),
		Context:              "",
		StorageClass:         vcc.controller.ComponentConfig().GetVclusterGatewayStorageClass(),
		DefaultStorageClass:  vcc.controller.ComponentConfig().GetVclusterGatewayDefaultStorageClass(),
		ChartRepo:            vcc.controller.ComponentConfig().GetVclusterGatewayChatRepo(),
		DefaultImageRegistry: vcc.controller.ComponentConfig().GetVclusterGatewayDefaultImageRegistry(),
		StorageManagerHost:   vcc.controller.ComponentConfig().GetStorageManagerHost(),
		Logger:               params.Logger,
		OrderDetails:         params,
		VClusterId:           "",
		Product:              "",
		Desc:                 "",
		ManagerBy:            "",
	}

	if info.Name == "" {
		info.Name = params.Name
	}
	if info.Comment == "" {
		info.Comment = params.Desc
	}

	info.Logger = log.New()

	resp, err := vcc.useCase.CreateVCluster(context.Background(), info)
	if err != nil {
		params.Logger.Errorf("create vcluster error: %v", err)

		if errors.Is(err, v1.ErrorInstanceIdAlreadyExists) {
			return &vclusterv1.CreateVClusterResponse{
				Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
				Msg:  v1.Cause(err),
				Data: nil,
			}, nil
		}

		if errors.Is(err, v1.ErrorVClusterNameAlreadyExists) {
			return &vclusterv1.CreateVClusterResponse{
				Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_CREATE_CLUSTER_ERROR),
				Msg:  v1.Cause(err),
				Data: nil,
			}, nil
		}

		return &vclusterv1.CreateVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_CREATE_CLUSTER_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_CREATE_CLUSTER_ERROR.String(),
			Data: nil,
		}, err
	}

	return &vclusterv1.CreateVClusterResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.Number()),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: &vclusterv1.CreateVClusterResponse_Data{AppId: resp.VClusterId},
	}, nil
}

func (vcc *VClusterController) UpdateVCluster(params *v1.UpdateVClusterParams) (*vclusterv1.UpdateVClusterResponse, error) {
	if len(params.Orders) == 0 {
		msg := "orderList is empty"
		return &vclusterv1.UpdateVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_UPDATE_CLUSTER_ERROR),
			Msg:  msg,
			Data: nil,
		}, errors.New(msg)
	}

	orderDetails := &v1.OrderDetail{
		Orders:   params.Orders,
		Logger:   params.Logger,
		Name:     params.Name,
		Desc:     params.Desc,
		UserName: params.UserName,
		TenantId: params.TenantId,
		RawName:  params.RawName,
		RawDesc:  params.RawDesc,
	}
	info := &v1.VClusterInfo{
		FallbackDns:          params.FallbackDns,
		CustomHelmConfig:     params.CustomHelmConfig,
		EnableHA:             params.EnableHA,
		IsInit:               params.IsInit,
		Id:                   params.AppId,
		Username:             params.UserName,
		TenantId:             params.TenantId,
		Name:                 params.RawName,
		Comment:              params.RawDesc,
		InstanceId:           params.Orders[0].InstanceID,
		Upgrade:              false,
		CephClusterId:        vcc.controller.ComponentConfig().GetCephClusterId(),
		Context:              "",
		StorageClass:         vcc.controller.ComponentConfig().GetVclusterGatewayStorageClass(),
		DefaultStorageClass:  vcc.controller.ComponentConfig().GetVclusterGatewayDefaultStorageClass(),
		ChartRepo:            vcc.controller.ComponentConfig().GetVclusterGatewayChatRepo(),
		DefaultImageRegistry: vcc.controller.ComponentConfig().GetVclusterGatewayDefaultImageRegistry(),
		StorageManagerHost:   vcc.controller.ComponentConfig().GetStorageManagerHost(),
		Logger:               params.Logger,
		OrderDetails:         orderDetails,
		VClusterId:           params.AppId,
		Product:              "",
		Desc:                 "",
		ManagerBy:            "",
	}

	if info.Name == "" {
		info.Name = params.Name
	}
	if info.Comment == "" {
		info.Comment = params.Desc
	}

	info.Logger = log.New()

	resp, err := vcc.useCase.UpdateVCluster(context.Background(), info)
	if err != nil {
		params.Logger.Errorf("update vcluster error: %v", err)
		return &vclusterv1.UpdateVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_UPDATE_CLUSTER_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_UPDATE_CLUSTER_ERROR.String(),
			Data: nil,
		}, err
	}

	return &vclusterv1.UpdateVClusterResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.Number()),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: &vclusterv1.UpdateVClusterResponse_Data{AppId: resp.VClusterId},
	}, nil
}

func (vcc *VClusterController) DeleteVCluster(params *v1.DeleteVClusterParams) (*vclusterv1.DeleteVClusterResponse, error) {
	_, err := vcc.useCase.DeleteVCluster(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("delete vcluster error: %v", err)
		return &vclusterv1.DeleteVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_DELETE_CLUSTER_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_DELETE_CLUSTER_ERROR.String(),
			Data: nil,
		}, nil
	}

	return &vclusterv1.DeleteVClusterResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: nil,
	}, nil
}

func (vcc *VClusterController) PauseVCluster(params *v1.PauseVClusterParams) (*vclusterv1.PauseVClusterResponse, error) {
	_, err := vcc.useCase.PauseVClusters(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("pause vcluster error: %v", err)
		return &vclusterv1.PauseVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_PAUSE_VCLUSTER_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_PAUSE_VCLUSTER_ERROR.String(),
			Data: nil,
		}, err
	}

	return &vclusterv1.PauseVClusterResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: &vclusterv1.PauseVClusterResponse_Data{},
	}, nil
}

func (vcc *VClusterController) ResumeVCluster(params *v1.ResumeVClusterParams) (*vclusterv1.ResumeVClusterResponse, error) {
	_, err := vcc.useCase.ResumeVClusters(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("resume vcluster error: %v", err)

		return &vclusterv1.ResumeVClusterResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR.String(),
			Data: nil,
		}, err
	}

	return &vclusterv1.ResumeVClusterResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: &vclusterv1.ResumeVClusterResponse_Data{},
	}, nil
}

func (vcc *VClusterController) GetKubeConfig(params *v1.GetKubeConfigParams) (*vclusterv1.GetKubeConfigResponse, error) {
	params.KubeConnHost = vcc.controller.ComponentConfig().GetVclusterGatewayKubeDaemonHost()

	resp, err := vcc.useCase.GetKubeConfig(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("get vcluster token error: %v", err)
		return &vclusterv1.GetKubeConfigResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_KUBECONFIG_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_KUBECONFIG_ERROR.String(),
			Data: nil,
		}, err
	}

	clusters := make([]*vclusterv1.GetKubeConfigResponse_KubeConfig_NamedCluster, 0)
	for _, v := range resp.Clusters {
		namedCluster := &vclusterv1.GetKubeConfigResponse_KubeConfig_NamedCluster{
			Cluster: &vclusterv1.GetKubeConfigResponse_KubeConfig_NamedCluster_Cluster{
				InsecureSkipTlsVerify: v.Cluster.InsecureSkipTLSVerify,
				Server:                v.Cluster.Server,
			},
			Name: v.Name,
		}
		clusters = append(clusters, namedCluster)
	}

	contexts := make([]*vclusterv1.GetKubeConfigResponse_KubeConfig_NamedContext, 0)
	for _, v := range resp.Contexts {
		namedContext := &vclusterv1.GetKubeConfigResponse_KubeConfig_NamedContext{
			Context: &vclusterv1.GetKubeConfigResponse_KubeConfig_NamedContext_Context{
				Cluster: v.Context.Cluster,
				User:    v.Context.User,
			},
			Name: v.Name,
		}
		contexts = append(contexts, namedContext)
	}

	users := make([]*vclusterv1.GetKubeConfigResponse_KubeConfig_NamedUser, 0)
	for _, v := range resp.Users {
		namedUser := &vclusterv1.GetKubeConfigResponse_KubeConfig_NamedUser{
			User: &vclusterv1.GetKubeConfigResponse_KubeConfig_NamedUser_User{
				Token: v.User.Token,
			},
			Name: v.Name,
		}
		users = append(users, namedUser)
	}

	return &vclusterv1.GetKubeConfigResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.Number()),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: &vclusterv1.GetKubeConfigResponse_KubeConfig{
			ApiVersion:     resp.APIVersion,
			Clusters:       clusters,
			Contexts:       contexts,
			CurrentContext: resp.CurrentContext,
			Kind:           resp.Kind,
			Preferences:    &vclusterv1.GetKubeConfigResponse_KubeConfig_Preferences{},
			Users:          users,
		},
	}, nil
}

func (vcc *VClusterController) GetKubeConfigBase64(params *v1.GetKubeConfigParams) (*vclusterv1.GetKubeConfigBase64Response, error) {
	params.KubeConnHost = vcc.controller.ComponentConfig().GetVclusterGatewayKubeDaemonHost()

	resp, err := vcc.useCase.GetKubeConfig(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("get vcluster token error: %v", err)
		return &vclusterv1.GetKubeConfigBase64Response{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_KUBECONFIG_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_KUBECONFIG_ERROR.String(),
			Data: "",
		}, err
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		params.Logger.Errorf("get vcluster token error: %v", err)
		return &vclusterv1.GetKubeConfigBase64Response{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_KUBECONFIG_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_KUBECONFIG_ERROR.String(),
			Data: "",
		}, err
	}

	return &vclusterv1.GetKubeConfigBase64Response{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.Number()),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: base64.StdEncoding.EncodeToString(jsonBytes),
	}, nil
}

func (vcc *VClusterController) QueryOperateStatus(params *v1.QueryOperateStatusRequest) (*vclusterv1.QueryOperateStatusResponse, error) {
	if params.AppId == "" || params.Action == "" {
		msg := "app_id or action is empty"
		return &vclusterv1.QueryOperateStatusResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR),
			Msg:  msg,
			Data: nil,
		}, errors.New(msg)
	}

	data, err := vcc.useCase.QueryOperateStatus(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("query operate status error: %v", err)

		return &vclusterv1.QueryOperateStatusResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR.String(),
			Data: nil,
		}, err
	}

	return &vclusterv1.QueryOperateStatusResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: data,
	}, nil
}

func (vcc *VClusterController) GetVClusterStatus(params *v1.GetVClusterStatusRequest) (*vclusterv1.GetVClusterStatusResponse, error) {
	if params.AppId == "" {
		msg := "app_id is empty"
		return &vclusterv1.GetVClusterStatusResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_STATUS_ERROR),
			Msg:  msg,
			Data: nil,
		}, errors.New(msg)
	}

	data, err := vcc.useCase.GetVClusterStatus(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("get vcluster status error: %v", err)

		return &vclusterv1.GetVClusterStatusResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_STATUS_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_STATUS_ERROR.String(),
			Data: nil,
		}, err
	}

	return &vclusterv1.GetVClusterStatusResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: data,
	}, nil
}

func (vcc *VClusterController) GetVClusterResourceDetails(params *v1.GetVClusterResourceDetailsRequest) (*vclusterv1.GetVClusterResourceDetailsResponse, error) {
	if params.AppId == "" {
		msg := "app_id is empty"
		return &vclusterv1.GetVClusterResourceDetailsResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_RESOURCE_DETAILS_ERROR),
			Msg:  msg,
			Data: nil,
		}, errors.New(msg)
	}

	data, err := vcc.useCase.GetVClusterResourceDetails(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("get vcluster resource details error: %v", err)

		return &vclusterv1.GetVClusterResourceDetailsResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_RESOURCE_DETAILS_ERROR),
			Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_RESOURCE_DETAILS_ERROR.String(),
			Data: nil,
		}, err
	}

	params.Logger.Infof("GetVClusterResourceDetails, data: %+v", data)

	return &vclusterv1.GetVClusterResourceDetailsResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: data,
	}, nil
}

func (vcc *VClusterController) GetVClusterContainerID(params *v1.GetVClusterContainerIDRequest) (*vclusterv1.GetVClusterContainerIDResponse, error) {
	if params.VClusterId == "" || params.Namespace == "" || params.PodName == "" || params.ContainerName == "" {
		msg := "vcluster_id, namespace, pod_name or container_name is empty"
		return &vclusterv1.GetVClusterContainerIDResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_CONTAINERID_ERROR),
			Msg:  msg,
			Data: nil,
		}, errors.New(msg)
	}

	data, err := vcc.useCase.GetVClusterContainerID(context.Background(), params)
	if err != nil {
		params.Logger.Errorf("get vcluster container id error: %v", err)

		return &vclusterv1.GetVClusterContainerIDResponse{
			Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_GET_CONTAINERID_ERROR),
			Msg:  err.Error(),
			Data: nil,
		}, err
	}

	return &vclusterv1.GetVClusterContainerIDResponse{
		Code: int32(vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:  vclusterv1.VClusterStatus_VCLUSTER_STATUS_OK.String(),
		Data: data,
	}, nil
}
