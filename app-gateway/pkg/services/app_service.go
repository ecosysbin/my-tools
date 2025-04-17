package services

import (
	"context"
	"fmt"
	"io"
	"strings"

	"encoding/json"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1/appv1connect"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/adapter"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/utils"
)

type AppServer struct {
	controller framework.Interface

	appv1connect.UnimplementedAppServiceHandler
}

func NewAppServer(controller framework.Interface) *AppServer {
	return &AppServer{
		controller: controller,
	}
}

func (vcs *AppServer) SetLogLevel(ctx context.Context, req *connect.Request[appv1.SetLoglevelRequest]) (*connect.Response[appv1.SetLoglevelResponse], error) {
	// utils.Logger.Level = req.Msg.Level
	return connect.NewResponse(&appv1.SetLoglevelResponse{
		Status: int32(appv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:    "alive",
	}), nil
}

func (vcs *AppServer) CheckHealth(ctx context.Context, req *connect.Request[appv1.CheckHealthRequest]) (*connect.Response[appv1.CheckHealthResponse], error) {
	return connect.NewResponse(&appv1.CheckHealthResponse{
		Status: int32(appv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Msg:    "alive",
	}), nil
}

func (as *AppServer) GetVClusterToken(ctx context.Context, req *connect.Request[appv1.GetVClusterTokenRequest]) (*connect.Response[appv1.GetVClusterTokenResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId := ""
	if user != nil {
		tenantId = user.Properties[v1.UserTenantId]
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	id := req.Msg.Id
	// token := req.Header().Get("X-Access-Token")
	// if token != "" {
	// instanceMap, err := as.Auth(token, AUTH_ACTION_GET)
	// if err != nil {
	// logger.Infof("auth failed, %v", err)
	// return connect.NewResponse(&appv1.GetVClusterTokenResponse{
	// Status: 170000,
	// Msg:    "auth failed: " + err.Error(),
	// }), nil
	// }
	// if _, ok := instanceMap[INSTANCE_MANAGER]; !ok {
	// if _, ok := instanceMap[id]; !ok {
	// return connect.NewResponse(&appv1.GetVClusterTokenResponse{
	// Status: 170000,
	// Msg:    "auth failed, no access to this app",
	// }), nil
	// }
	// }
	// }
	logger.Infof("get vcluster token, app id: %s", id)

	// auth
	// authToken := req.Header().Get("X-Access-Token")
	// if authToken != "" && as.
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.GetVClusterTokenResponse{
			Status: 170002,
			Msg:    "db list app config err: " + err.Error(),
		}), nil
	}
	// params := &v1.GetAppParams{
	// 	TenantId:     tenantId,
	// 	AppId:        id,
	// 	Logger:       logger,
	// 	AppConfigMap: appConfigMap,
	// }
	// resp, err := as.controller.AppController().
	// 	GetApp(params)
	resp, err := as.controller.AppRepo().GetByAppId(id)
	if err != nil {
		logger.Infof("get app %s err, %v", id, err)
		return connect.NewResponse(&appv1.GetVClusterTokenResponse{
			Status: 170002,
			Msg:    "get app err: " + err.Error(),
		}), nil
	}

	appConfig := appConfigMap["vcluster"]
	url := fmt.Sprintf("%s/v1/kubeconfig/base64/%s", appConfig.Doamin, resp.AppId)

	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"

	tokenResp, err := adapter.HttpGetRequest(httpHeader, url)
	if err != nil {
		return connect.NewResponse(&appv1.GetVClusterTokenResponse{
			Status: 170003,
			Msg:    "http get kubeconfig err: " + err.Error(),
		}), nil
	}
	defer tokenResp.Body.Close()
	body, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return connect.NewResponse(&appv1.GetVClusterTokenResponse{
			Status: 170003,
			Msg:    "http get kubeconfig err: " + err.Error(),
		}), nil
	}

	logger.Infof("Received response: %s", string(body))
	type RespBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
		Data string `json:"data,omitempty"`
	}
	var createResp RespBody
	if err = json.Unmarshal(body, &createResp); err != nil {
		logger.Errorf("Failed to unmarshal response body, URL: %s, err: %v, body: %s", url, err, string(body))
		return connect.NewResponse(&appv1.GetVClusterTokenResponse{
			Status: 170003,
			Msg:    "http get kubeconfig err: Unmarshal response body err: " + err.Error(),
		}), nil
	}
	if createResp.Code != 0 {
		return connect.NewResponse(&appv1.GetVClusterTokenResponse{
			Status: 170003,
			Msg:    "http get kubeconfig err: " + string(body),
		}), nil
	}
	return connect.NewResponse(&appv1.GetVClusterTokenResponse{
		Status: int32(appv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Data:   createResp.Data,
	}), nil
}

func (as *AppServer) UpdateInstanceStatus(ctx context.Context, req *connect.Request[appv1.UpdateInstanceStatusRequest]) (*connect.Response[appv1.UpdateInstanceStatusResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId := ""
	if user != nil {
		tenantId = user.Properties[v1.UserTenantId]
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// 转发服务不做授权判断
	// token := req.Header().Get("X-Access-Token")
	// auth
	// if token != "" {
	// 	if _, err := as.Auth(token, AUTH_ACTION_UPDATE); err != nil {
	// 		logger.Infof("auth failed, %v", err)
	// 		return connect.NewResponse(&appv1.UpdateInstanceStatusResponse{
	// 			Status: 170000,
	// 			Msg:    "auth failed",
	// 		}), nil
	// 	}
	// }
	// logger.Infof("tenantId: %s username: %s", tenantId, username)
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.UpdateInstanceStatusResponse{
			Status: 170002,
			Msg:    "list app config err:" + err.Error(),
		}), nil
	}
	config := appConfigMap["aps-serving"]
	url := fmt.Sprintf("%s/updateInstanceStatus", config.Url)
	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"
	httpHeader["apikey"] = as.controller.ComponentConfig().GetApsApiKey()
	// req body
	reqBody := v1.UpdateAppInstanceStatusData{
		TenantId:             req.Msg.TenantId,
		InstanceId:           req.Msg.InstanceId,
		AppInstanceId:        req.Msg.AppInstanceId,
		Valid:                req.Msg.Valid,
		ProductCategoryValue: req.Msg.ProductCategoryValue,
		Reason: v1.Reason{
			Zh: req.Msg.Reason.Zh,
			En: req.Msg.Reason.En,
		},
	}
	reqBodyMar, err := json.Marshal(reqBody)
	if err != nil {
		return connect.NewResponse(&appv1.UpdateInstanceStatusResponse{
			Status: 170002,
			Msg:    "update instance status err: " + err.Error(),
		}), nil
	}
	if _, err := adapter.AppPostRestRequest(httpHeader, url, reqBodyMar, v1.HttpLogPrint{}); err != nil {
		logger.Info("update instance status err, %v", err)
		if strings.Contains(err.Error(), "unmarshal response body") {
			return connect.NewResponse(&appv1.UpdateInstanceStatusResponse{
				Status: 200,
			}), nil
		}
		return connect.NewResponse(&appv1.UpdateInstanceStatusResponse{
			Status: 170002,
			Msg:    "update instance status err: " + err.Error(),
		}), nil
	}
	return connect.NewResponse(&appv1.UpdateInstanceStatusResponse{
		Status: 200,
	}), nil
}

func (as *AppServer) PauseApp(ctx context.Context, req *connect.Request[appv1.PauseAppRequest]) (*connect.Response[appv1.PauseAppResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		if tenantType == "TENANT_TYPE_PLATFORM" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// token := req.Header().Get("X-Access-Token")
	// // auth
	// if token != "" {
	// 	if _, err := as.Auth(token, AUTH_ACTION_PAUSE); err != nil {
	// 		logger.Infof("auth failed, %v", err)
	// 		return connect.NewResponse(&appv1.PauseAppResponse{
	// 			Status: 170000,
	// 			Msg:    "auth failed: " + err.Error(),
	// 		}), nil
	// 	}
	// }
	logger.Infof("Pause App: %v", req.Msg.Ids)
	// 查询所有app实例
	apps, err := as.controller.AppRepo().ListAll(tenantId)
	if err != nil {
		return connect.NewResponse(&appv1.PauseAppResponse{
			Status: 170002,
			Msg:    "db query err: " + err.Error(),
		}), nil
	}
	appMap := map[string]repo.AppRecord{}
	for _, app := range apps {
		appMap[app.Id] = app
	}
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.PauseAppResponse{
			Status: 170002,
			Msg:    "list app config err: " + err.Error(),
		}), nil
	}
	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"

	for _, id := range req.Msg.Ids {
		if appRecord, ok := appMap[id]; ok {
			appRecord.Status = v1.APP_STATUS_PAUSING
			appConfig := appConfigMap[appRecord.ManageBy]

			apiKey := ""
			if strings.HasPrefix(appRecord.ManageBy, "aps") {
				apiKey = as.controller.ComponentConfig().GetApsApiKey()
			} else {
				apiKey = as.controller.ComponentConfig().GetAlayaStudioApiKey()
			}
			httpHeader["apikey"] = apiKey

			params := &v1.WorkflowAppParams{
				Logger:     logger,
				Action:     v1.ACTION_PAUSE,
				HttpHeader: httpHeader,
				AppConfig:  appConfig,
				AppRecord:  appRecord,
				Conditions: v1.Conditions{
					Action:    v1.ACTION_PAUSE,
					PreStatus: v1.APP_STATUS_RUNNING,
					Status:    v1.APP_STATUS_PENDING, // 小状态，可以表示创建时的具体状态，做流程控制使用
					Events:    []string{utils.ParseTimeEvent("prepare to pause app.")},
				},
			}
			// 1. 修改资源实例状态 (pausing)
			params.AppRecord.Conditions = adapter.TransformConditionsToString(&params.Conditions)
			if err := as.controller.AppRepo().Update(params.AppRecord); err != nil {
				params.Logger.Errorf("store app err, %v", err)
				return connect.NewResponse(&appv1.PauseAppResponse{
					Status: 170002,
					Msg:    "pause app err: " + err.Error(),
				}), nil
			}
			log.Infof("pausing app name %s, id %s", params.AppRecord.Name, params.AppRecord.Id)
			// 2. pause实例
			if _, err := as.controller.AppController().PauseApp(params); err != nil {
				logger.Infof("pause app err, %v", err)
				return connect.NewResponse(&appv1.PauseAppResponse{
					Status: 170002,
					Msg:    "pause app err: " + err.Error(),
				}), nil
			}
		}
	}
	return connect.NewResponse(&appv1.PauseAppResponse{
		Status: 200,
	}), nil
}

func (as *AppServer) ResumeApp(ctx context.Context, req *connect.Request[appv1.ResumeAppRequest]) (*connect.Response[appv1.ResumeAppResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		if tenantType == "TENANT_TYPE_PLATFORM" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// token := req.Header().Get("X-Access-Token")
	// // auth
	// if token != "" {
	// 	if _, err := as.Auth(token, AUTH_ACTION_RESUME); err != nil {
	// 		logger.Infof("auth failed, %v", err)
	// 		return connect.NewResponse(&appv1.ResumeAppResponse{
	// 			Status: 170000,
	// 			Msg:    "auth failed: " + err.Error(),
	// 		}), nil
	// 	}
	// }
	logger.Infof("Resume App: %v", req.Msg.Ids)
	// logger.Infof("tenantId: %s username: %s", tenantId, username)

	// 查询所有app实例
	apps, err := as.controller.AppRepo().ListAll(tenantId)
	if err != nil {
		return connect.NewResponse(&appv1.ResumeAppResponse{
			Status: 170002,
			Msg:    "db query err: " + err.Error(),
		}), nil
	}
	appMap := map[string]repo.AppRecord{}
	for _, app := range apps {
		appMap[app.Id] = app
	}
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.ResumeAppResponse{
			Status: 170002,
			Msg:    "list app config err: " + err.Error(),
		}), nil
	}
	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"

	for _, id := range req.Msg.Ids {
		if appRecord, ok := appMap[id]; ok {
			appRecord.Status = v1.APP_STATUS_RESUMING
			appConfig := appConfigMap[appRecord.ManageBy]

			apiKey := ""
			if strings.HasPrefix(appRecord.ManageBy, "aps") {
				apiKey = as.controller.ComponentConfig().GetApsApiKey()
			} else {
				apiKey = as.controller.ComponentConfig().GetAlayaStudioApiKey()
			}
			httpHeader["apikey"] = apiKey

			params := &v1.WorkflowAppParams{
				Logger:     logger,
				Action:     v1.ACTION_RESUME,
				HttpHeader: httpHeader,
				AppConfig:  appConfig,
				AppRecord:  appRecord,
				Conditions: v1.Conditions{
					Action:    v1.ACTION_RESUME,
					PreStatus: v1.APP_STATUS_PAUSED,
					Status:    v1.APP_STATUS_PENDING, // 小状态，可以表示创建时的具体状态，做流程控制使用
					Events:    []string{utils.ParseTimeEvent("prepare to resume app.")},
				},
			}
			// 1. 修改资源实例状态 (resuming)
			params.AppRecord.Conditions = adapter.TransformConditionsToString(&params.Conditions)
			if err := as.controller.AppRepo().Update(params.AppRecord); err != nil {
				params.Logger.Errorf("store app err, %v", err)
				return connect.NewResponse(&appv1.ResumeAppResponse{
					Status: 170002,
					Msg:    "resume app err: " + err.Error(),
				}), nil
			}
			log.Infof("resuming app name %s, id %s", params.AppRecord.Name, params.AppRecord.Id)
			// 2. resume实例
			if _, err := as.controller.AppController().ResumeApp(params); err != nil {
				logger.Infof("resume app err, %v", err)
				return connect.NewResponse(&appv1.ResumeAppResponse{
					Status: 170002,
					Msg:    "resume app err: " + err.Error(),
				}), nil
			}
		}
	}
	return connect.NewResponse(&appv1.ResumeAppResponse{
		Status: 200,
	}), nil
}

func (as *AppServer) UpdateApp(ctx context.Context, req *connect.Request[appv1.UpdateAppRequest]) (*connect.Response[appv1.UpdateAppResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		if tenantType == "TENANT_TYPE_PLATFORM" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// auth
	// token := req.Header().Get("X-Access-Token")
	// if token != "" {
	// 	if _, err := as.Auth(token, AUTH_ACTION_UPDATE); err != nil {
	// 		logger.Infof("auth failed, %v", err)
	// 		return connect.NewResponse(&appv1.UpdateAppResponse{
	// 			Status: 170000,
	// 			Msg:    "auth failed: " + err.Error(),
	// 		}), nil
	// 	}
	// }
	logger.Infof("Update App: %v", req.Msg.AppId)
	// 查询app实例
	_, err := as.controller.AppRepo().GetByAppId(req.Msg.AppId)
	if err != nil {
		logger.Infof("repo get app err, %v", err)
		return connect.NewResponse(&appv1.UpdateAppResponse{
			Status: 170001,
			Msg:    "request param err: get app by id failed: " + err.Error(),
		}), nil
	}
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.UpdateAppResponse{
			Status: 170002,
			Msg:    "list app config err: " + err.Error(),
		}), nil
	}

	// 构造updateAppParam
	updateAppParam := v1.CreateAppParams{
		Logger: logger,
		Id:     req.Msg.AppId,
		// AppId:     req.Msg.AppId,
		AidcId:    int64(req.Msg.AidcId),
		Name:      req.Msg.GetName(),
		Desc:      req.Msg.GetDesc(),
		OrderList: formatOrderList(req.Msg.GetOrderlist()),
		UserName:  username,
		TenantId:  tenantId,
		User: v1.GCPUser{
			UserID:   user.Properties["userId"],
			UserName: username,
			NickName: user.Name,
			Email:    user.Email,
			PhoneNum: user.Phone,
		},
		// HttpHeader:   httpHeader,
	}
	// 构建AppRequestData对象
	appRequestData, err := adapter.TransformCreateAppParamsToAppRequestData(&updateAppParam)
	if err != nil {
		logger.Infof("TransformCreateAppParamsToAppRequestData err, %v req.msg: %s", err, toSerialize(req.Msg))
		return connect.NewResponse(&appv1.UpdateAppResponse{
			Status: 170002,
			Msg:    "check app parameters err: " + err.Error(),
		}), nil
	}

	// appconfig
	appConfig := appConfigMap[appRequestData.Domain]
	appRequestData.Spec.Kubeconfig = appConfig.KubeConfig

	// 查询当前的appRecord
	repoApp, err := as.controller.AppRepo().GetByAppId(req.Msg.AppId)
	if err != nil {
		logger.Info("get app record err, %v", err)
		return connect.NewResponse(&appv1.UpdateAppResponse{
			Status: 170001,
			Msg:    "request param err: get app by id failed: " + err.Error(),
		}), nil
	}
	repoApp.Name = appRequestData.Name
	repoApp.Desc = appRequestData.Desc
	repoApp.Status = v1.APP_STATUS_UPDATING
	repoApp.CreateTime = utils.TimeNow()
	repoApp.Message = toSerialize(appRequestData)
	repoApp.OriginMessage = toSerialize(req.Msg)

	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"
	apiKey := ""
	if strings.HasPrefix(appRequestData.Domain, "aps") {
		apiKey = as.controller.ComponentConfig().GetApsApiKey()
	} else {
		apiKey = as.controller.ComponentConfig().GetAlayaStudioApiKey()
	}
	httpHeader["apikey"] = apiKey

	params := v1.WorkflowAppParams{
		Logger:         logger,
		Action:         v1.ACTION_UPDATE,
		UpdateAppReq:   req.Msg,
		AppConfig:      appConfig,
		AppRequestData: *appRequestData,
		AppRecord:      repoApp,
		HttpHeader:     httpHeader,
		Conditions: v1.Conditions{
			Action:    v1.ACTION_UPDATE,
			PreStatus: v1.APP_STATUS_RUNNING,
			Status:    v1.APP_STATUS_PENDING, // 小状态，可以表示创建时的具体状态，做流程控制使用
			Events:    []string{utils.ParseTimeEvent("prepare to update app.")},
		},
	}
	// 1. 资源实例入库 (pending)
	params.AppRecord.Conditions = adapter.TransformConditionsToString(&params.Conditions)
	if err := as.controller.AppRepo().Update(params.AppRecord); err != nil {
		params.Logger.Errorf("store app err, %v", err)
		return connect.NewResponse(&appv1.UpdateAppResponse{
			Status: 170001,
			Msg:    "create app err: " + err.Error(),
		}), nil
	}
	// 2. 更新app实例
	resp, err := as.controller.AppController().
		UpdateApp(&params)
	if err != nil {
		logger.Infof("create app err, %v", err)
		return connect.NewResponse(&appv1.UpdateAppResponse{
			Status: 170002,
			Msg:    "create app err: " + err.Error(),
		}), nil
	}
	return connect.NewResponse(resp), nil
}

func (as *AppServer) GetApp(ctx context.Context, req *connect.Request[appv1.GetAppRequest]) (*connect.Response[appv1.GetAppResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		if tenantType == "TENANT_TYPE_PLATFORM" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// logger.Infof("user: %v", user)
	// auth
	token := req.Header().Get("X-Access-Token")

	logger.Infof("appId: %s", req.Msg.AppId)
	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"
	apiKey := ""
	httpHeader["apikey"] = apiKey

	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.GetAppResponse{
			Status: 170002,
			Msg:    "db list app config err: " + err.Error(),
		}), nil
	}

	// auth
	// authToken := req.Header().Get("X-Access-Token")
	// if authToken != "" && as.
	params := &v1.GetAppParams{
		TenantId:     tenantId,
		TenantType:   tenantType,
		ReqToken:     token,
		AppId:        req.Msg.AppId,
		AppConfigMap: appConfigMap,
		HttpHeader:   httpHeader,
		Logger:       logger,
	}
	resp, err := as.controller.AppController().
		GetApp(params)
	if err != nil {
		logger.Infof("get app %s err, %v", req.Msg.AppId, err)
		return connect.NewResponse(&appv1.GetAppResponse{
			Status: 170002,
			Msg:    "get app err: " + err.Error(),
		}), nil
	}
	// vcluster不需要鉴权
	if token != "" && resp.GetData().GetManageBy() != "vcluster" {
		instanceMap, err := as.Auth(token, AUTH_ACTION_GET)
		if err != nil {
			logger.Infof("auth failed, %v", err)
			return connect.NewResponse(&appv1.GetAppResponse{
				Status: 170000,
				Msg:    "auth failed: " + err.Error(),
			}), nil
		}
		if _, ok := instanceMap[INSTANCE_MANAGER]; !ok {
			if _, ok := instanceMap[req.Msg.AppId]; !ok {
				return connect.NewResponse(&appv1.GetAppResponse{
					Status: 170000,
					Msg:    "auth failed, no access to this app",
				}), nil
			}
		}
	}

	return connect.NewResponse(resp), nil
}

func (as *AppServer) ListApp(ctx context.Context, req *connect.Request[appv1.ListAppRequest]) (*connect.Response[appv1.ListAppResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	// tenantId := user.Properties[v1.UserTenantId]
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		// 针对的是管理员租户和apiKey场景
		if tenantType == "TENANT_TYPE_PLATFORM" || tenantId == "datacanvas.internal.server" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// logger.Infof("user: %v", user)
	logger.Infof("tenantId: %s tenantType: %s username: %s delete: %d", tenantId, tenantType, username, req.Msg.Delete)
	// auth
	token := req.Header().Get("X-Access-Token")
	appAuthMap := map[string]string{}
	var err error
	// var response *connect.Response
	if token != "" {
		appAuthMap, err = as.Auth(token, AUTH_ACTION_GET)
		if err != nil {
			logger.Infof("auth failed, %v", err)
			appAuthMap = map[string]string{}
			// return connect.NewResponse(&appv1.ListAppResponse{
			// 	Status: 170000,
			// 	Msg:    "auth failed" + err.Error(),
			// }), nil
			// response := connect.NewResponse(&appv1.ListAppResponse{
			// 	Status: 170000,
			// 	Msg:    "auth failed: " + err.Error(),
			// })
			// return response, nil
		}
	}
	// tenantId, tenantType := as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
	params := &v1.ListAppParams{
		Logger:     logger,
		TenantId:   tenantId,
		Username:   username,
		Delete:     req.Msg.Delete,
		Page:       int(req.Msg.PageIndex),
		PageSize:   int(req.Msg.PageNum),
		Id:         req.Msg.Id,
		InstanceId: req.Msg.InstanceId,
		Name:       req.Msg.Name,
		ManageBy:   req.Msg.ManageBy,
		CreateBy:   req.Msg.CreateUser,
		Status:     req.Msg.Status,
		TenantIds:  req.Msg.TenantIds,
		CreateTime: v1.PeriodTime{
			StartTime: req.Msg.CreateStartTime,
			EndTime:   req.Msg.CreateEndTime,
		},
		DeleteTime: v1.PeriodTime{
			StartTime: req.Msg.DeleteStartTime,
			EndTime:   req.Msg.DeleteEndTime,
		},
	}
	apps, err := as.controller.AppController().GetAppList(params)
	if err != nil {
		logger.Infof("list app err, %v", err)
		response := connect.NewResponse(&appv1.ListAppResponse{
			Status: 170002,
			Msg:    "db list app err: " + err.Error(),
		})
		return response, nil
	}
	// token为空或是管理员租户或apiKey场景
	if _, ok := appAuthMap[INSTANCE_MANAGER]; ok || token == "" {
		response := connect.NewResponse(apps)
		return response, nil
	}
	// 使用authMap过滤app
	responseApp := []*appv1.ListAppResponse_App{}
	for _, app := range apps.Data.Rows {
		_, ok := appAuthMap[app.AppId]
		if !ok && app.ManageBy != "vcluster" {
			// 默认是false，再dto中做了配置
			app.AccessDeny = "true"
		}
		responseApp = append(responseApp, app)
	}
	apps.Data.Rows = responseApp
	response := connect.NewResponse(apps)
	logger.Infof("list app response success")
	return response, nil
}

func (as *AppServer) DeleteApp(ctx context.Context, req *connect.Request[appv1.DeleteAppRequest]) (*connect.Response[appv1.DeleteAppResponse], error) {
	// accessToken := req.Header().Get("X-Access-Token")
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		if tenantType == "TENANT_TYPE_PLATFORM" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	logger.Infof("Delete App: %v", req.Msg.Ids)
	logger.Infof("Delete App deletestorage: %v", req.Msg.IsDeleteStorage)
	// logger.Infof("user: %v", user)
	// auth
	// token := req.Header().Get("X-Access-Token")
	// if token != "" {
	// 	if _, err := as.Auth(token, AUTH_ACTION_DELETE); err != nil {
	// 		logger.Infof("auth failed, %v", err)
	// 		return connect.NewResponse(&appv1.DeleteAppResponse{
	// 			Status: 170000,
	// 			Msg:    "auth failed: " + err.Error(),
	// 		}), nil
	// 	}
	// }

	// 查询所有app实例
	apps, err := as.controller.AppRepo().ListAll(tenantId)
	if err != nil {
		return connect.NewResponse(&appv1.DeleteAppResponse{
			Status: 170002,
			Msg:    "db query err: " + err.Error(),
		}), nil
	}
	appMap := map[string]repo.AppRecord{}
	for _, app := range apps {
		appMap[app.Id] = app
	}
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.DeleteAppResponse{
			Status: 170002,
			Msg:    "list app config err: " + err.Error(),
		}), nil
	}
	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"

	httpParams := map[string]string{}
	httpParams["isDeleteStorage"] = req.Msg.IsDeleteStorage
	logger.Info("applist map: %v", appMap)
	for _, id := range req.Msg.Ids {
		if appRecord, ok := appMap[id]; ok {
			// 记录释放之前的app状态，再释放失败时需要恢复到该状态
			preStatus := appRecord.Status
			appRecord.Status = v1.APP_STATUS_DELETING
			appConfig := appConfigMap[appRecord.ManageBy]

			apiKey := ""
			if strings.HasPrefix(appRecord.ManageBy, "aps") {
				apiKey = as.controller.ComponentConfig().GetApsApiKey()
			} else {
				apiKey = as.controller.ComponentConfig().GetAlayaStudioApiKey()
			}
			httpHeader["apikey"] = apiKey

			params := &v1.WorkflowAppParams{
				Logger:     logger,
				Action:     v1.ACTION_DELETE,
				HttpHeader: httpHeader,
				HttpParams: httpParams,
				AppConfig:  appConfig,
				AppRecord:  appRecord,
				Conditions: v1.Conditions{
					Action:    v1.ACTION_DELETE,
					PreStatus: preStatus,
					Status:    v1.APP_STATUS_PENDING, // 小状态，可以表示创建时的具体状态，做流程控制使用
					Events:    []string{utils.ParseTimeEvent("prepare to delete app.")},
				},
			}
			// 1. 修改资源实例状态 (deleting)
			params.AppRecord.Conditions = adapter.TransformConditionsToString(&params.Conditions)
			if err := as.controller.AppRepo().Update(params.AppRecord); err != nil {
				params.Logger.Errorf("store app err, %v", err)
				return connect.NewResponse(&appv1.DeleteAppResponse{
					Status: 170001,
					Msg:    "delete app err: " + err.Error(),
				}), nil
			}
			log.Infof("deleting app name %s, id %s", params.AppRecord.Name, params.AppRecord.Id)
			// 2. 工作流删除资源
			if _, err := as.controller.AppController().DeleteApp(params); err != nil {
				logger.Infof("delete app err, %v", err)
				return connect.NewResponse(&appv1.DeleteAppResponse{
					Status: 170002,
					Msg:    "delete app err: " + err.Error(),
				}), nil
			}
		}
	}
	return connect.NewResponse(&appv1.DeleteAppResponse{
		Status: 200,
	}), nil
}

var formatInstanceSpecs = func(before []*appv1.Orderlist_InstanceSpec) (after []v1.InstanceSpec) {
	after = make([]v1.InstanceSpec, 0, len(before))
	for _, instanceSpec := range before {
		after = append(after, v1.InstanceSpec{
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

var formatOrderList = func(before []*appv1.Orderlist) (after []v1.Order) {
	after = make([]v1.Order, 0, len(before))
	for _, order := range before {
		after = append(after, v1.Order{
			OrderInfo: v1.OrderInfo{
				ProductId:        int(order.OrderInfo.ProductId),
				ProductCode:      order.OrderInfo.ProductCode,
				CycleCount:       int(order.OrderInfo.CycleCount),
				Amount:           int(order.OrderInfo.Amount),
				OrderType:        int(order.OrderInfo.OrderType),
				ResourceTypeId:   int(order.OrderInfo.ResourceTypeId),
				ResourceTypeCode: order.OrderInfo.ResourceTypeCode,
				ActualAmount:     int(order.OrderInfo.ActualAmount),
			},
			InstanceSpec:       formatInstanceSpecs(order.InstanceSpec),
			InstanceId:         order.InstanceId,
			NodePoolInstanceId: order.NodePoolInstanceId,
		})
	}
	return after
}

func (as *AppServer) CreateApp(ctx context.Context, req *connect.Request[appv1.CreateAppRequest]) (*connect.Response[appv1.CreateAppResponse], error) {
	user, username := as.controller.GetUser(ctx), as.controller.GetUsername(ctx)
	tenantId, tenantType := "", ""
	if user != nil {
		tenantId, tenantType = as.controller.GetUser(ctx).Properties[v1.UserTenantId], as.controller.GetUser(ctx).Properties[v1.UserTenantType]
		if tenantType == "TENANT_TYPE_PLATFORM" {
			// repo 在tenantId = "" 查询所有租户的app
			tenantId = ""
		}
	}
	logger := as.controller.GetLogger(ctx).WithField("username", username).WithField("tenantId", tenantId)
	// 创建先不鉴权，后续根据iam通知再加上
	// logger.Info("header: %v", accessToken)
	// logger.Infof("user: %v", user)
	// auth
	// token := req.Header().Get("X-Access-Token")
	// if token != "" {
	// 	if _, err := as.Auth(token, AUTH_ACTION_CREATE); err != nil {
	// 		logger.Infof("auth failed, %v", err)
	// 		return connect.NewResponse(&appv1.CreateAppResponse{
	// 			Status: 170000,
	// 			Msg:    "auth failed: " + err.Error(),
	// 		}), nil
	// 	}
	// }
	// app config
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		logger.Info("list app config err, %v", err)
		return connect.NewResponse(&appv1.CreateAppResponse{
			Status: 170002,
			Msg:    "db list app config err: " + err.Error(),
		}), nil
	}

	// 构造CreateAppParams
	id := uuid.New().String()
	createAppParam := v1.CreateAppParams{
		Logger:    logger,
		Id:        id,
		AppId:     id,
		AidcId:    int64(req.Msg.GetAidcId()),
		Name:      req.Msg.GetName(),
		Desc:      req.Msg.GetDesc(),
		OrderList: formatOrderList(req.Msg.GetOrderlist()),
		UserName:  username,
		TenantId:  tenantId,
		User: v1.GCPUser{
			UserID:   user.Properties["userId"],
			UserName: username,
			NickName: user.Name,
			Email:    user.Email,
			PhoneNum: user.Phone,
		},
		// HttpHeader:   httpHeader,
	}
	// 构建AppRequestData对象
	appRequestData, err := adapter.TransformCreateAppParamsToAppRequestData(&createAppParam)
	if err != nil {
		logger.Infof("TransformCreateAppParamsToAppRequestData err, %v, req.msg: %s", err, toSerialize(req.Msg))
		return connect.NewResponse(&appv1.CreateAppResponse{
			Status: 170002,
			Msg:    "transform app request data err: " + err.Error(),
		}), nil
	}
	if err := appRequestData.IsValid(); err != nil {
		logger.Infof("check app request data err, %v msg.req: %s", err, toSerialize(req.Msg))
		return connect.NewResponse(&appv1.CreateAppResponse{
			Status: 170002,
			Msg:    "check parameters err, " + err.Error(),
		}), nil
	}
	// appconfig
	appConfig := appConfigMap[appRequestData.Domain]
	appRequestData.Spec.Kubeconfig = appConfig.KubeConfig

	// 构造appRecord record
	repoApp := repo.AppRecord{
		Id:            appRequestData.AppId,
		InstanceId:    appRequestData.InstanceId,
		Name:          appRequestData.Name,
		Desc:          appRequestData.Desc,
		TenantId:      appRequestData.GcpTenantID,
		CreateUser:    appRequestData.User.UserName,
		AppId:         appRequestData.AppId,
		CreateTime:    utils.TimeNow(),
		Status:        v1.APP_STATUS_CREATING,
		ManageBy:      appRequestData.Domain,
		OriginMessage: toSerialize(req.Msg),
		Message:       toSerialize(appRequestData),
	}

	// http header
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"
	apiKey := ""
	if strings.HasPrefix(appRequestData.Domain, "aps") {
		apiKey = as.controller.ComponentConfig().GetApsApiKey()
	} else {
		apiKey = as.controller.ComponentConfig().GetAlayaStudioApiKey()
	}
	httpHeader["apikey"] = apiKey

	params := v1.WorkflowAppParams{
		Logger: logger,
		Action: v1.ACTION_CREATE,
		// Repo:           as.controller.AppRepo(),
		CreateAppReq:   req.Msg,
		AppConfig:      appConfig,
		AppRequestData: *appRequestData,
		Conditions: v1.Conditions{
			Action:    v1.ACTION_CREATE,
			PreStatus: "",
			Status:    v1.APP_STATUS_PENDING, // 小状态，可以表示创建时的具体状态，做流程控制使用
			Events:    []string{utils.ParseTimeEvent("prepare to create app.")},
		},
		AppRecord:  repoApp,
		HttpHeader: httpHeader,
	}
	// 1. 资源实例入库 (Creating)
	params.AppRecord.Conditions = adapter.TransformConditionsToString(&params.Conditions)
	if err := as.controller.AppRepo().Store(params.AppRecord); err != nil {
		params.Logger.Errorf("store app err, %v", err)
		return connect.NewResponse(&appv1.CreateAppResponse{
			Status: 170001,
			Msg:    "create app err: " + err.Error(),
		}), nil
	}
	// 2. 创建app
	resp, err := as.controller.AppController().
		CreateApp(&params)
	if err != nil {
		logger.Infof("create app err, %v msg.req: %s", err, toSerialize(req.Msg))
		return connect.NewResponse(&appv1.CreateAppResponse{
			Status: 170002,
			Msg:    "create app err: " + err.Error(),
		}), nil
	}
	return connect.NewResponse(resp), nil
}

func toSerialize(input interface{}) string {
	b, _ := json.Marshal(input)
	return string(b)
}

func (as *AppServer) appConfigMap() (map[string]v1.AppConfig, error) {
	// 查询appConfig
	appConfigs, err := as.controller.AppRepo().ListAppConfig()
	if err != nil {
		return nil, err
	}
	// 构造创建app config
	appConfigMap := map[string]v1.AppConfig{}
	for _, config := range appConfigs {
		appConfigMap[config.AppType] = v1.AppConfig{
			Doamin:     config.Domain,
			AppType:    config.AppType,
			Url:        fmt.Sprintf("%s/%s/%s", config.Domain, config.Version, "app"),
			KubeConfig: config.KubeConfig,
		}
	}
	return appConfigMap, nil
}

func (as *AppServer) SyncAppStatus() error {
	apps, err := as.controller.AppRepo().ListAll("")
	if err != nil {
		return err
	}
	appConfigMap, err := as.appConfigMap()
	if err != nil {
		return err
	}
	// http header
	httpHeader := map[string]string{}
	// sync has no token
	// httpHeader["X-Access-Token"] = req.Header().Get("X-Access-Token")
	httpHeader["Content-Type"] = "application/json"

	for _, app := range apps {
		if v1.IsFailedStatus(app.Status) || app.Status == v1.APP_STATUS_RUNNING || app.Status == v1.APP_STATUS_PAUSED || app.Status == v1.App_STATUS_DELETED {
			continue
		}
		log.Infof("sync app %s status %s", app.Name, app.Status)
		var appRequestData v1.AppRequestData
		if err := json.Unmarshal([]byte(app.Message), &appRequestData); err != nil {
			log.Warnf("sync app checkout appRequestData err, %v", err)
			continue
		}
		var appOriginRequestData appv1.CreateAppRequest
		if err := json.Unmarshal([]byte(app.OriginMessage), &appOriginRequestData); err != nil {
			log.Warnf("sync app checkout appOriginRequestData err, %v", err)
			continue
		}
		apiKey := ""
		if strings.HasPrefix(appRequestData.Domain, "aps") {
			apiKey = as.controller.ComponentConfig().GetApsApiKey()
		} else {
			apiKey = as.controller.ComponentConfig().GetAlayaStudioApiKey()
		}
		httpHeader["apikey"] = apiKey
		params := v1.WorkflowAppParams{
			Logger:         log.GlobalLogger().WithField("tenantId", app.TenantId).WithField("app", app.Name),
			AppConfig:      appConfigMap[appRequestData.Domain],
			AppRequestData: appRequestData,
			CreateAppReq:   &appOriginRequestData,
			HttpHeader:     httpHeader,
			AppRecord:      app,
			Conditions: v1.Conditions{
				PreStatus: app.Status,
				Status:    v1.APP_STATUS_PENDING, // 小状态，可以表示创建时的具体状态，做流程控制使用
				Events:    []string{utils.ParseTimeEvent("sync app status.")},
			},
		}

		if app.Status == v1.APP_STATUS_CREATING {
			params.Action = v1.ACTION_CREATE
			params.Conditions.Action = v1.ACTION_CREATE
			if _, err := as.controller.AppController().CreateApp(&params); err != nil {
				log.Warnf("sync app %s create app err, %v", app.Name, err)
			}
		}
		if app.Status == v1.APP_STATUS_DELETING {
			params.Action = v1.ACTION_DELETE
			params.Conditions.Action = v1.ACTION_DELETE
			if _, err := as.controller.AppController().DeleteApp(&params); err != nil {
				log.Warnf("sync app %s create app err, %v", app.Name, err)
			}
		}
		if app.Status == v1.APP_STATUS_UPDATING {
			params.Action = v1.ACTION_UPDATE
			params.Conditions.Action = v1.ACTION_UPDATE
			if _, err := as.controller.AppController().UpdateApp(&params); err != nil {
				log.Warnf("sync app %s create app err, %v", app.Name, err)
			}
		}
		if app.Status == v1.APP_STATUS_PAUSING {
			params.Action = v1.ACTION_PAUSE
			params.Conditions.Action = v1.ACTION_PAUSE
			if _, err := as.controller.AppController().PauseApp(&params); err != nil {
				log.Warnf("sync app %s create app err, %v", app.Name, err)
			}
		}
		if app.Status == v1.APP_STATUS_RESUMING {
			params.Action = v1.ACTION_RESUME
			params.Conditions.Action = v1.ACTION_RESUME
			if _, err := as.controller.AppController().ResumeApp(&params); err != nil {
				log.Warnf("sync app %s create app err, %v", app.Name, err)
			}
		}
	}
	return nil
}
