package app

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"

	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/adapter"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/utils"
)

type AppController struct {
	controller framework.Interface
}

func NewAppController(controller framework.Interface) *AppController {
	ac := &AppController{
		controller: controller,
	}
	return ac
}

func (ac *AppController) GetApp(params *v1.GetAppParams) (*appv1.GetAppResponse, error) {
	app, err := ac.controller.AppRepo().GetByAppId(params.AppId)
	if err != nil {
		params.Logger.Errorf("get app %s err, %v", params.AppId, err)
		return nil, err
	}
	if params.ReqToken != "" && app.TenantId != params.TenantId && params.TenantType != "TENANT_TYPE_PLATFORM" {
		return &appv1.GetAppResponse{
			Status: 170000,
			Msg:    "auth failed, tenantId not match " + params.TenantId,
		}, nil
	}

	// params.Logger.Infof("getapp: %v", app)
	responseApp := TransformAppRecordToGetAppResponse(app, nil)
	// rawCluster get cluster metrics
	if app.ManageBy == "vcluster" {
		appConfig := params.AppConfigMap[app.ManageBy]
		metrics, config, err := ac.rawClusterConfigs(appConfig.Url, app.AppId, params.HttpHeader)
		if err != nil {
			params.Logger.Errorf("get app %s metrics err, %v", params.AppId, err)
		}
		responseApp.UtilizationRate = metrics
		responseApp.Config = config
	}
	return &appv1.GetAppResponse{
		Status: int32(appv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Data:   responseApp,
	}, nil
}

func (ac *AppController) rawClusterConfigs(url, appId string, header map[string]string) (*appv1.Resourcequotas, map[string]string, error) {
	// 根据appId从vluster-gateway获取指标数据
	getUrl := fmt.Sprintf("%s/%s/resourcedetails", url, appId)
	resp, err := adapter.HttpGetRequest(header, getUrl)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "[Failed to get app resourcequota, app: %s]", appId)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to read response body, URL: %s", getUrl)
	}
	// {"code":0, "msg":"VCLUSTER_STATUS_OK", "data":{"utilizationRate":{"gpu":[{"name":"nvidia.com/gpu-tesla-p4", "hard":2, "used":0}], "memory":{"hard":1600,
	// "used":8}, "cpu":{"hard":9994, "used":6}, "storage":[{"name":"FS-HDD", "hard":100, "used":0}]}}}
	log.Infof("Sending provisioning request to URL: %s response: %s", body, string(body))
	type detail struct {
		Rate   appv1.Resourcequotas `json:"utilizationRate"`
		Config map[string]string    `json:"configurations"`
	}

	type respBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
		Data detail `json:"data,omitempty"`
	}

	var getResp respBody
	if err = json.Unmarshal(body, &getResp); err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to unmarshal response body, URL: %s, body: %s", getUrl, string(body))
	}

	if getResp.Code != 0 {
		return nil, nil, fmt.Errorf("failed to get app resourcequota, app: %s, code: %d, msg: %s", appId, getResp.Code, getResp.Msg)
	}
	return &getResp.Data.Rate, getResp.Data.Config, nil
}

func (ac *AppController) GetAppList(params *v1.ListAppParams) (*appv1.ListAppResponse, error) {
	// 监控数据没有，得组装
	options := repo.ListOptions{
		TenantId:   params.TenantId,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Id:         params.Id,
		Name:       params.Name,
		InstanceId: params.InstanceId,
		Deleted:    params.Delete,
		ManageBy:   params.ManageBy,
		CreateBy:   params.CreateBy,
		Status:     params.Status,
		TenantIds:  params.TenantIds,
		CreateTime: &repo.PeriodTime{
			StartTime: params.CreateTime.StartTime,
			EndTime:   params.CreateTime.EndTime,
		},
		DeleteTime: &repo.PeriodTime{
			StartTime: params.DeleteTime.StartTime,
			EndTime:   params.DeleteTime.EndTime,
		},
	}
	var apps []repo.AppRecord
	var size int64
	var err error
	// if params.Delete == 1 {
	// 	apps, size, err = ac.controller.AppRepo().ListDeletedPageAll(options)
	// } else {
	apps, size, err = ac.controller.AppRepo().ListPageAll(options)
	// }
	if err != nil {
		params.Logger.Errorf("list apps err, %v", err)
		return nil, err
	}
	// sort
	// sort.Slice(apps, func(i, j int) bool {
	// 	return apps[i].CreateTime.After(*apps[j].CreateTime)
	// })
	usageMetrics, err := ac.UsageMetrics()
	if err != nil {
		log.Warnf("get usage metrics err, %v", err)
	}
	log.Infof("usageMetrics %v", usageMetrics)
	gpuMetrics := ac.GpuUsageMetrics(apps)
	// test ok
	// gpuMetrics := map[string]map[string]int64{}
	// gpuMetrics["0c9f5fe4-e20b-4840-a0d1-48f8f0acdf12"] = map[string]int64{"NVIDIA-Tesla-P4": 1}

	log.Infof("gpuMetrics %v", gpuMetrics)
	responseAppList := TransformAppRecordListToAppResponse(apps, &usageMetrics, gpuMetrics)
	data := &appv1.ListAppResponse_Data{
		Count: int32(size),
		Rows:  responseAppList,
	}
	// log.Infof("responseAppList appList %v", responseAppList)
	return &appv1.ListAppResponse{
		Status: int32(appv1.VClusterStatus_VCLUSTER_STATUS_OK),
		Data:   data,
	}, nil
}

func (ac *AppController) UpdateApp(params *v1.WorkflowAppParams) (*appv1.UpdateAppResponse, error) {
	updateAppWorkFlow := AppWorkFlow{
		Metadata: params,
		Works: []AppWork{
			{
				WorkName:     WORK_UPDATE_UPDATEINSTANCE,
				Work:         ac.updateAppInstance,
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_UPDATE_UPDATEUPDATINGSTATUS,
				Work:     ac.updateAppStatus,
			},
			{
				WorkName:     WORK_UPDATE_CHECKAPPUPDATESTATUS,
				Work:         WorkforTimeout(WORK_UPDATE_CHECKAPPUPDATESTATUS, 5, 3*60, ac.checkAppStatus),
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_UPDATE_UPDATERUNNIGSTATUS,
				Work:     ac.updateAppStatus,
			},
		},
	}
	go updateAppWorkFlow.Start()
	return &appv1.UpdateAppResponse{
		Status: 200,
		Data:   params.AppRecord.Id,
	}, nil
}

func (ac *AppController) CreateApp(params *v1.WorkflowAppParams) (*appv1.CreateAppResponse, error) {
	createAppWorkFlow := AppWorkFlow{
		Metadata: params,
		Works: []AppWork{
			{
				WorkName:     WORK_CREATE_CREATEINSTANCE,
				Work:         ac.createAppInstance,
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_CREATE_UPDATECREATESTATUS,
				Work:     ac.updateAppStatus,
			},
			{
				WorkName:     WORK_CREATE_CHECKAPPCREATESTATUS,
				Work:         WorkforTimeout(WORK_CREATE_CHECKAPPCREATESTATUS, 5, 3*60, ac.checkAppStatus),
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_CREATE_UPDATERUNNIGSTATUS,
				Work:     ac.updateAppStatus,
			},
		},
	}
	go createAppWorkFlow.Start()
	return &appv1.CreateAppResponse{
		Status: 200,
		Data:   params.AppRecord.Id,
	}, nil
}

func reasonDetail(reason string) string {
	return fmt.Sprintf("Time: %s Reason: [%s]", utils.FormartTimeNow(), reason)
}

func (ac *AppController) createAppInstance(params *v1.WorkflowAppParams) (err error) {
	defer func() {
		if err != nil {
			// params.AppRecord.Reason = err.Error()
			// params.AppRecord.Status = v1.APP_STATUS_FAILED
			params.Conditions.Status = v1.APP_STATUS_CREATE_FAILED
			params.Conditions.Reason = reasonDetail(err.Error())
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("create app instance failed, %v", err)))

			// app的状态需要恢复到之前的状态，因为创建过程之前的状态为空，则直接使用二级状态作为app的状态
			params.AppRecord.Status = params.Conditions.Status
			params.AppRecord.Reason = params.Conditions.Reason
			RefreshSyncWorker(params)
		}
	}()

	logHeader := deepcopyMap(params.HttpHeader)
	logHeader["X-Access-Token"] = utils.MaskToken(logHeader["X-Access-Token"])
	logHeader["apikey"] = utils.MaskToken(logHeader["apikey"])
	var reqBody []byte
	var logBody []byte
	if params.AppRecord.ManageBy == "vcluster" {
		reqBody, err = json.Marshal(params.CreateAppReq)
		if err != nil {
			return errors.Wrapf(err, "[Failed to marshal request data, params: %v]", params.CreateAppReq)
		}
		logBody = reqBody
	} else {
		formatReqBody := params.AppRequestData
		formatReqBody.Spec.Kubeconfig = utils.MaskToken(formatReqBody.Spec.Kubeconfig)
		logBody, _ = json.Marshal(formatReqBody)
		reqBody, err = json.Marshal(params.AppRequestData)
		if err != nil {
			return errors.Wrapf(err, "[Failed to marshal request data, params: %v]", params.AppRequestData)
		}
	}
	httpLogPrint := v1.HttpLogPrint{
		Url:    params.AppConfig.Url,
		Header: logHeader,
		Method: "POST",
		Body:   string(logBody),
	}
	params.Logger.Infof("http create app instance type %s, url %s", params.AppRecord.ManageBy, params.AppConfig.Url)
	resp, err := adapter.AppPostRestRequest(params.HttpHeader, params.AppConfig.Url, reqBody, httpLogPrint)
	if err != nil {
		return errors.Wrapf(err, "[Failed to create app: %s, err, %v]", params.AppRecord.Name, err)
	}
	if resp.Code != 0 {
		return errors.Errorf("create app: %s Http Post err, errCode: %d, msg: [%s]", params.AppRecord.Name, resp.Code, resp.Msg)
	}
	if resp.Data.AppId != "" {
		// 下层应用选择用自己的appId
		params.AppRecord.AppId = resp.Data.AppId
	}
	params.Conditions.Status = v1.APP_STATUS_CREATING
	params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("send create app request to %s", params.AppRecord.ManageBy)))
	return nil
}

func deepcopyMap(originalMap map[string]string) map[string]string {
	// 创建一个新的 map
	copiedMap := make(map[string]string)

	// 遍历原始 map，将每个键值对拷贝到新 map 中
	for key, value := range originalMap {
		copiedMap[key] = value
	}
	return copiedMap
}

func (ac *AppController) updateAppInstance(params *v1.WorkflowAppParams) (err error) {
	var reqBody []byte
	defer func() {
		if err != nil {
			params.Conditions.Status = v1.ActionStatusFailed(params.Action)
			params.Conditions.Reason = err.Error()
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("update app instance failed, %v", err)))

			params.AppRecord.Reason = params.Conditions.Reason
			// 失败时，app的状态需要恢复到之前的状态，若之前的状态为空，则直接使用二级状态作为app的状态
			if params.Conditions.PreStatus == "" {
				params.AppRecord.Status = params.Conditions.Status
			} else {
				params.AppRecord.Status = params.Conditions.PreStatus
			}
			RefreshSyncWorker(params)
		}
	}()
	if params.AppRecord.ManageBy == "raw" || params.AppRecord.ManageBy == "vcluster" {
		reqBody, err = json.Marshal(params.UpdateAppReq)
		if err != nil {
			return errors.Wrapf(err, "[Failed to marshal request data, params: %v]", params.CreateAppReq)
		}
	} else {
		reqBody, err = json.Marshal(params.AppRequestData)
		if err != nil {
			return errors.Wrapf(err, "[Failed to marshal request data, params: %v]", params.AppRequestData)
		}
	}
	params.Logger.Infof("http update app instance type %s, url %s", params.AppRecord.ManageBy, params.AppConfig.Url)
	url := fmt.Sprintf("%s/%s", params.AppConfig.Url, params.AppRecord.AppId)
	resp, err := adapter.HttpPutRequest(params.HttpHeader, url, reqBody)
	if err != nil {
		return errors.Wrapf(err, "Failed to create app: %s, err, %v", params.AppRecord.Name, err)
	}
	if resp.Code != 0 {
		return errors.Errorf("update app: %s Http Post err, errCode %d, msg: [%s]", params.AppRecord.Name, resp.Code, resp.Msg)
	}
	if resp.Data.AppId != "" {
		// 下层应用选择用自己的appId
		params.AppRecord.AppId = resp.Data.AppId
	}
	// params.AppRecord.Status = v1.APP_STATUS_UPDATING
	params.Conditions.Status = v1.APP_STATUS_UPDATING
	params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("send update app request to %s", params.AppRecord.ManageBy)))
	return nil
}

func (ac *AppController) deleteAppInstance(params *v1.WorkflowAppParams) (err error) {
	defer func() {
		if err != nil {
			params.Conditions.Status = v1.ActionStatusFailed(params.Action)
			params.Conditions.Reason = err.Error()
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("delete app instance failed, %v", err)))

			params.AppRecord.Reason = params.Conditions.Reason
			// 失败时，app的状态需要恢复到之前的状态，若之前的状态为空，则直接使用二级状态作为app的状态
			if params.Conditions.PreStatus == "" {
				params.AppRecord.Status = params.Conditions.Status
			} else {
				params.AppRecord.Status = params.Conditions.PreStatus
			}
			RefreshSyncWorker(params)
		}
	}()
	url := fmt.Sprintf("%s/%s", params.AppConfig.Url, params.AppRecord.AppId)
	if params.HttpParams != nil {
		url = fmt.Sprintf("%s?%s", url, params.HttpParams.Encode())
	}
	resp, err := adapter.HttpDeleteRequest(params.HttpHeader, url)
	if err != nil {
		return errors.Wrapf(err, "Failed to delete app: %s, err, %v", params.AppRecord.Name, err)
	}
	if resp.Code != 0 {
		return errors.Errorf("delete app: %s Http Post err, errCode %d, msg: [%s]", params.AppRecord.Name, resp.Code, resp.Msg)
	}
	params.Conditions.Status = v1.APP_STATUS_DELETING
	params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("send delete app request to %s", params.AppRecord.ManageBy)))
	return nil
}

func (ac *AppController) pauseAppInstance(params *v1.WorkflowAppParams) (err error) {
	defer func() {
		if err != nil {
			params.Conditions.Status = v1.ActionStatusFailed(params.Action)
			params.Conditions.Reason = err.Error()
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("update app instance failed, %v", err)))

			params.AppRecord.Reason = params.Conditions.Reason
			// 失败时，app的状态需要恢复到之前的状态，若之前的状态为空，则直接使用二级状态作为app的状态
			if params.Conditions.PreStatus == "" {
				params.AppRecord.Status = params.Conditions.Status
			} else {
				params.AppRecord.Status = params.Conditions.PreStatus
			}
			RefreshSyncWorker(params)
		}
	}()
	url := fmt.Sprintf("%s/%s/pause", params.AppConfig.Url, params.AppRecord.AppId)
	resp, err := adapter.AppPostRestRequest(params.HttpHeader, url, nil, v1.HttpLogPrint{})
	if err != nil {
		return errors.Wrapf(err, "[Failed to create app: %s, err, %v]", params.AppRecord.Name, err)
	}
	if resp.Code != 0 {
		return errors.Errorf("create app: %s Http Post err, errCode %d, msg [%s]", params.AppRecord.Name, resp.Code, resp.Msg)
	}
	params.Conditions.Status = v1.APP_STATUS_PAUSING
	params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("send pause app request to %s", params.AppRecord.ManageBy)))
	return nil
}

func (ac *AppController) resumeAppInstance(params *v1.WorkflowAppParams) (err error) {
	defer func() {
		if err != nil {
			params.Conditions.Status = v1.ActionStatusFailed(params.Action)
			params.Conditions.Reason = err.Error()
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("resume app instance failed, %v", err)))

			params.AppRecord.Reason = params.Conditions.Reason
			// 失败时，app的状态需要恢复到之前的状态，若之前的状态为空，则直接使用二级状态作为app的状态
			if params.Conditions.PreStatus == "" {
				params.AppRecord.Status = params.Conditions.Status
			} else {
				params.AppRecord.Status = params.Conditions.PreStatus
			}
			RefreshSyncWorker(params)
		}
	}()
	url := fmt.Sprintf("%s/%s/resume", params.AppConfig.Url, params.AppRecord.AppId)
	resp, err := adapter.AppPostRestRequest(params.HttpHeader, url, nil, v1.HttpLogPrint{})
	if err != nil {
		return errors.Wrapf(err, "Failed to create app: %s, err, %v", params.AppRecord.Name, err)
	}
	if resp.Code != 0 {
		return errors.Errorf("create app: %s Http Post err, errCode %d, msg [%s]", params.AppRecord.Name, resp.Code, resp.Msg)
	}
	params.Conditions.Status = v1.APP_STATUS_RESUMING
	params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("send resume app request to %s", params.AppRecord.ManageBy)))
	return nil
}

func (ac *AppController) checkAppStatus(params *v1.WorkflowAppParams) (err error) {
	defer func() {
		if err != nil {
			// params.AppRecord.Reason = err.Error()
			// params.AppRecord.Status = v1.APP_STATUS_FAILED
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(err.Error()))
			// 失败时，app的状态需要恢复到之前的状态，若之前的状态为空，则直接使用二级状态作为app的状态
			// need retry err 无需恢复app状态
			if !errors.Is(err, ErrNeedRetry) {
				params.Conditions.Status = v1.ActionStatusFailed(params.Action)
				params.Conditions.Reason = err.Error()
				if params.Conditions.PreStatus == "" {
					params.AppRecord.Status = params.Conditions.Status
				} else {
					params.AppRecord.Status = params.Conditions.PreStatus
				}
				params.AppRecord.Reason = params.Conditions.Reason
			}
		}
	}()
	getUrl := fmt.Sprintf("%s/%s?action=%s", params.AppConfig.Url, params.AppRecord.AppId, params.Action)
	resp, err := adapter.AppGetRequest(params.HttpHeader, getUrl)
	if err != nil {
		return errors.Wrapf(err, "[Failed to create app, app: %s]", params.AppRecord.Name)
	}
	if resp.Code != 0 {
		return errors.Errorf("http get app: %s status failed, code %d msg [%s]", params.AppRecord.Name, resp.Code, resp.Msg)
	}
	// 上面需要速错
	if resp.Data.Status == v1.ACTION_STATUS_SUCCESS {
		// 任务完成对conditions进行清空
		params.AppRecord.Conditions = ""
		// 修改资源实例状态为deleted
		if params.Action == v1.ACTION_DELETE {
			params.AppRecord.DeleteTime = utils.TimeNow()
			params.Conditions.Status = v1.App_STATUS_DELETED
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("delete app return %v", resp.Data.Status)))
			params.AppRecord.Status = v1.App_STATUS_DELETED
			params.AppRecord.Deleted = 1
		}
		if params.Action == v1.ACTION_CREATE {
			// 修改资源实例状态为running
			params.AppRecord.StartedTime = utils.TimeNow()
			params.Conditions.Status = v1.APP_STATUS_RUNNING
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("create app return %v", resp.Data.Status)))
			params.AppRecord.Status = v1.APP_STATUS_RUNNING
			params.AppRecord.Url = resp.Data.Url
			params.AppRecord.MonitorUrl = resp.Data.MonitorUrl
			params.AppRecord.Reason = ""
		}
		if params.Action == v1.ACTION_UPDATE {
			// 修改资源实例状态为running
			params.AppRecord.StartedTime = utils.TimeNow()
			params.Conditions.Status = v1.APP_STATUS_RUNNING
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("update app return %v", resp.Data.Status)))
			params.AppRecord.Status = v1.APP_STATUS_RUNNING
			params.AppRecord.Url = resp.Data.Url
			params.AppRecord.MonitorUrl = resp.Data.MonitorUrl
			params.AppRecord.Reason = ""
		}
		if params.Action == v1.ACTION_PAUSE {
			// 修改资源实例状态为paused
			params.Conditions.Status = v1.APP_STATUS_PAUSED
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("pause app return %v", resp.Data.Status)))
			params.AppRecord.Status = v1.APP_STATUS_PAUSED
		}
		if params.Action == v1.ACTION_RESUME {
			// 修改资源实例状态为running
			params.AppRecord.StartedTime = utils.TimeNow()
			params.Conditions.Status = v1.APP_STATUS_RUNNING
			params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("resume app return %v", resp.Data.Status)))
			params.AppRecord.Status = v1.APP_STATUS_RUNNING
		}
		// success后清理缓存中的事件记录
		StopSyncWorker(params)
	} else {
		// 这里需要继续循环watch
		// return errors.Errorf("app: %s status is %s", params.AppRecord.Name, resp.Data.Status)
		// 任务完成对conditions进行清空
		// params.Conditions.Events = append(params.Conditions.Events, utils.ParseTimeEvent(fmt.Sprintf("app: %s status is %s, wait for success", params.AppRecord.Name, resp.Data.Status)))
		return errors.Wrapf(ErrNeedRetry, "app: %s current status is %s", params.AppRecord.Name, params.Conditions.Status)
	}
	return nil
}

func (ac *AppController) updateAppStatus(params *v1.WorkflowAppParams) error {
	if params.AppRecord.Reason != "" {
		params.AppRecord.Reason = reasonDetail(params.AppRecord.Reason)
	}
	params.AppRecord.Conditions = adapter.TransformConditionsToString(&params.Conditions)
	if err := ac.controller.AppRepo().Update(params.AppRecord); err != nil {
		return errors.Errorf("update app %s status %s err, %v", params.AppRecord.Name, params.AppRecord.Status, err)
	}
	return nil
}

func (ac *AppController) DeleteApp(params *v1.WorkflowAppParams) (*appv1.DeleteAppResponse, error) {
	deleteAppWorkFlow := AppWorkFlow{
		Metadata: params,
		Works: []AppWork{
			{
				WorkName:     WORK_DELETE_DELETEINSTANCE,
				Work:         ac.deleteAppInstance,
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName:     WORK_DELETE_CHECKAPPDELETESTATUS,
				Work:         WorkforTimeout(WORK_DELETE_CHECKAPPDELETESTATUS, 5, 3*60, ac.checkAppStatus),
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_DELETE_UPDATEDELETEDTATUS,
				Work:     ac.updateAppStatus,
			},
		},
	}
	go deleteAppWorkFlow.Start()
	return &appv1.DeleteAppResponse{
		Status: 200,
	}, nil
}

func (ac *AppController) PauseApp(params *v1.WorkflowAppParams) (*appv1.PauseAppResponse, error) {
	pauseAppWorkFlow := AppWorkFlow{
		Metadata: params,
		Works: []AppWork{
			{
				WorkName:     WORK_PAUSE_PAUSEINSTANCE,
				Work:         ac.pauseAppInstance,
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName:     WORK_PAUSE_CHECKAPPPAUSESTATUS,
				Work:         WorkforTimeout(WORK_PAUSE_CHECKAPPPAUSESTATUS, 5, 3*60, ac.checkAppStatus),
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_PAUSE_UPDATEDPAUSEDTATUS,
				Work:     ac.updateAppStatus,
			},
		},
	}
	go pauseAppWorkFlow.Start()
	return &appv1.PauseAppResponse{
		Status: 200,
	}, nil
}

func (ac *AppController) ResumeApp(params *v1.WorkflowAppParams) (*appv1.ResumeAppResponse, error) {
	pauseAppWorkFlow := AppWorkFlow{
		Metadata: params,
		Works: []AppWork{
			{
				WorkName:     WORK_RESUME_RESUMEINSTANCE,
				Work:         ac.resumeAppInstance,
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName:     WORK_RESUME_CHECKAPPPSESUMESTATUS,
				Work:         WorkforTimeout(WORK_RESUME_CHECKAPPPSESUMESTATUS, 5, 3*60, ac.checkAppStatus),
				FailCallBack: ac.updateAppStatus,
			},
			{
				WorkName: WORK_PAUSE_UPDATEDPAUSEDTATUS,
				Work:     ac.updateAppStatus,
			},
		},
	}
	go pauseAppWorkFlow.Start()
	return &appv1.ResumeAppResponse{
		Status: 200,
	}, nil
}

func (ac *AppController) UsageMetrics() (v1.UsageMetrics, error) {
	metricClient := adapter.NewMetricClient(nil, ac.controller.ComponentConfig().GetMetricsEndpoint(), "/api/v1/query")
	return metricClient.UsageMetrics()
}

func (ac *AppController) GpuUsageMetrics(apps []repo.AppRecord) map[string]map[string]float32 {
	metricClient := adapter.NewMetricClient(nil, ac.controller.ComponentConfig().GetMetricsEndpoint(), "/api/v1/query")
	gpuMetrics := map[string]map[string]float32{}
	for _, app := range apps {
		appGpuMetric, err := metricClient.AllTypeGpuUsageMetrics(app.InstanceId)
		// appGpuMetric, err := metricClient.AllTypeGpuUsageMetrics(app.AppId)
		if err != nil {
			log.Warnf("get app %s gpu metrics err, %v", app.AppId, err)
			continue
		}
		gpuMetrics[app.InstanceId] = appGpuMetric
	}
	return gpuMetrics
}
