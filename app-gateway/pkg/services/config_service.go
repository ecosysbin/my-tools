package services

import (
	"context"

	"connectrpc.com/connect"
	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
)

func (as *AppServer) ListConfig(ctx context.Context, req *connect.Request[appv1.ListConfigRequest]) (*connect.Response[appv1.ListConfigResponse], error) {
	// accessToken := req.Header().Get("X-Access-Token")
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId := user.Properties[v1.UserTenantId]
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	logger.Infof("user: %v", user)
	logger.Infof("tenantId: %s username: %s", tenantId, username)

	configs, err := as.controller.AppController().ListAppConfig()
	if err != nil {
		logger.Infof("list app err, %v", err)
		return connect.NewResponse(&appv1.ListConfigResponse{
			Status: 160005,
			Msg:    "list config err",
		}), nil
	}

	configList := make([]*appv1.ListConfigResponse_Config, 0)
	for _, c := range configs {
		cnfig := &appv1.ListConfigResponse_Config{
			AppType:    c.AppType,
			Domain:     c.Domain,
			Version:    c.Version,
			Kubeconfig: c.KubeConfig,
		}
		configList = append(configList, cnfig)
	}
	return connect.NewResponse(&appv1.ListConfigResponse{
		Status: 200,
		Data:   configList,
	}), nil
}

func (as *AppServer) AddConfig(ctx context.Context, req *connect.Request[appv1.AddConfigRequest]) (*connect.Response[appv1.AddConfigResponse], error) {
	// accessToken := req.Header().Get("X-Access-Token")
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId := user.Properties[v1.UserTenantId]
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	logger.Infof("user: %v", user)
	logger.Infof("tenantId: %s username: %s add or update config", tenantId, username)

	config := repo.AppConfig{
		AppType:    req.Msg.AppType,
		Domain:     req.Msg.Domain,
		Version:    req.Msg.Version,
		KubeConfig: req.Msg.Kubeconfig,
	}

	err := as.controller.AppController().AddConfig(config)
	if err != nil {
		logger.Infof("add app err, %v", err)
		return connect.NewResponse(&appv1.AddConfigResponse{
			Status: 160006,
			Msg:    "add config err",
		}), nil
	}

	return connect.NewResponse(&appv1.AddConfigResponse{
		Status: 200,
	}), nil
}

func (as *AppServer) DeleteConfig(ctx context.Context, req *connect.Request[appv1.DeleteConfigRequest]) (*connect.Response[appv1.DeleteConfigResponse], error) {
	// accessToken := req.Header().Get("X-Access-Token")
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId := user.Properties[v1.UserTenantId]
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	logger.Infof("user: %v", user)
	logger.Infof("tenantId: %s username: %s", tenantId, username)

	err := as.controller.AppController().DeleteConfig(req.Msg.AppType)
	if err != nil {
		logger.Infof("delete config err, %v", err)
		return connect.NewResponse(&appv1.DeleteConfigResponse{
			Status: 160007,
			Msg:    "delete config err",
		}), nil
	}

	return connect.NewResponse(&appv1.DeleteConfigResponse{
		Status: 200,
	}), nil
}
