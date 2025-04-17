package app

import (
	"encoding/json"
	"time"

	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/adapter"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/utils"
)

func TransformAppRecordToGetAppResponse(appRecored repo.AppRecord, usedMetrics *v1.UsageMetrics) *appv1.GetAppResponse_App {
	appResponse := appv1.GetAppResponse_App{
		Id:              appRecored.Id,
		TenantId:        appRecored.TenantId,
		InstanceId:      appRecored.InstanceId,
		AppId:           appRecored.AppId,
		Name:            appRecored.Name,
		InstanceSpec:    appRecored.Desc,
		Comment:         appRecored.Desc,
		ManageBy:        appRecored.ManageBy,
		AppURL:          appRecored.Url,
		AppBaseURL:      utils.ExtractDomain(appRecored.Url),
		AppPath:         utils.ExtractURLPath(appRecored.Url),
		MonitorUrl:      appRecored.MonitorUrl,
		Status:          appRecored.Status,
		Reason:          appRecored.Reason,
		CreateTime:      formatTime(appRecored.CreateTime),
		StartedTime:     formatTime(appRecored.StartedTime),
		DeleteTime:      formatTime(appRecored.DeleteTime),
		UtilizationRate: &appv1.Resourcequotas{},
		CreateUser:      appRecored.CreateUser,
	}

	metrics := &appv1.Resourcequotas{}
	appRequestData := checkoutAppRequestData(appRecored.Message)
	if appRequestData != nil {
		// cpu metric
		if usedMetrics != nil {
			metrics.Cpu = map[string]float32{
				"used": usedMetrics.CpuUsage[appRecored.AppId],
				"hard": float32(appRequestData.Spec.Cpus),
			}
			metrics.Memory = map[string]float32{
				"used": usedMetrics.MemUsage[appRecored.AppId],
				"hard": float32(appRequestData.Spec.Mem),
			}
		}
	}
	appResponse.UtilizationRate = metrics
	return &appResponse
}

func TransformAppRecordToAppResponse(appRecored repo.AppRecord, usedMetrics *v1.UsageMetrics, gpuAllTypeMetrics map[string]map[string]float32) *appv1.ListAppResponse_App {
	appResponse := &appv1.ListAppResponse_App{
		Id:              appRecored.Id,
		TenantId:        appRecored.TenantId,
		InstanceId:      appRecored.InstanceId,
		AppId:           appRecored.AppId,
		Name:            appRecored.Name,
		Comment:         appRecored.Desc,
		InstanceSpec:    appRecored.Desc,
		ManageBy:        appRecored.ManageBy,
		AppURL:          appRecored.Url,
		AppBaseURL:      utils.ExtractDomain(appRecored.Url),
		AppPath:         utils.ExtractURLPath(appRecored.Url),
		MonitorUrl:      appRecored.MonitorUrl,
		Status:          appRecored.Status,
		Reason:          appRecored.Reason,
		CreateTime:      formatTime(appRecored.CreateTime),
		StartedTime:     formatTime(appRecored.StartedTime),
		DeleteTime:      formatTime(appRecored.DeleteTime),
		UtilizationRate: &appv1.Resourcequotas{},
		CreateUser:      appRecored.CreateUser,
		AccessDeny:      "false",
	}

	metrics := &appv1.Resourcequotas{}
	storageMetrics := []*appv1.Resourcequotas_Quota{}
	gpuMetrics := []*appv1.Resourcequotas_Quota{}
	config := map[string]string{}
	appRequestData := checkoutAppRequestData(appRecored.Message)
	if appRequestData != nil {
		// cpu metric
		if usedMetrics != nil {
			metrics.Cpu = map[string]float32{
				"used": usedMetrics.CpuUsage[appRecored.InstanceId],
				// "used": usedMetrics.CpuUsage[appRecored.AppId],
				"hard": float32(appRequestData.Spec.Cpus),
			}
			metrics.Memory = map[string]float32{
				"used": usedMetrics.MemUsage[appRecored.InstanceId],
				// "used": usedMetrics.MemUsage[appRecored.AppId],
				"hard": float32(appRequestData.Spec.Mem),
			}
		}
	}
	// storage metric
	for _, storage := range appRequestData.StorageList {
		storageMetrics = append(storageMetrics, &appv1.Resourcequotas_Quota{
			// Name: storage.Name,
			Hard: float32(storage.Limit),
			// Used: usedMetrics.Storage[appRecored.AppId],
			Used: usedMetrics.Storage[appRecored.InstanceId],
		})
	}
	// gpu metric
	for _, gpu := range appRequestData.Spec.Gpu {
		// gpu 总量为0时，则不显示
		if gpu.Count == 0 {
			continue
		}
		gpuMetrics = append(gpuMetrics, &appv1.Resourcequotas_Quota{
			Name: gpu.Type,
			Hard: float32(gpu.Count),
			// Used: gpuAllTypeMetrics[appRecored.AppId][gpu.Type],
			Used: gpuAllTypeMetrics[appRecored.InstanceId][gpu.Type],
		})
	}
	if len(appRequestData.StorageList) == 1 {
		config["storageType"] = appRequestData.StorageList[0].FsName
	}
	if config["storageType"] == "" {
		// 兼容1.2.1到1.3.0升级
		log.Warnf("app %s storage type is empty", appRecored.Id)
		config["storageType"] = "capacity"
	}

	metrics.Storage = storageMetrics
	metrics.Gpu = gpuMetrics
	appResponse.UtilizationRate = metrics
	appResponse.Config = config
	// 获取conditions
	var conditions = &appv1.ListAppResponse_Conditions{}
	syncWorker := GetSyncWorker(appRecored.Id)
	if syncWorker != nil {
		conditions = &appv1.ListAppResponse_Conditions{
			Action: syncWorker.Conditions.Action,
			Status: syncWorker.Conditions.Status,
			Reason: syncWorker.Conditions.Reason,
			Events: syncWorker.Conditions.Events,
		}
	} else {
		log.Infof("sync worker not found for app %s", appRecored.Id)
		recordConditions, err := adapter.TransformStringToConditions(appRecored.Conditions)
		if err != nil {
			log.Infof("transform string %s to conditions error: %v", appRecored.Conditions, err)
		}
		if err == nil && recordConditions != nil {
			conditions = &appv1.ListAppResponse_Conditions{
				Action: recordConditions.Action,
				Status: recordConditions.Status,
				Reason: recordConditions.Reason,
				Events: recordConditions.Events,
			}
		}
	}
	appResponse.Conditions = conditions
	return appResponse
}

// func formatGpuType(gpuType string) string {
// 	if strings.Contains(gpuType, " ") {
// 		return strings.ReplaceAll(gpuType, " ", "-")
// 	}
// 	return gpuType
// }

func checkoutAppRequestData(msg string) *v1.AppRequestData {
	var appRequestData v1.AppRequestData
	json.Unmarshal([]byte(msg), &appRequestData)
	return &appRequestData
}

func TransformAppRecordListToAppResponse(appRecoredList []repo.AppRecord, usageMetrics *v1.UsageMetrics, gpuAllTypeMetrics map[string]map[string]float32) []*appv1.ListAppResponse_App {
	appResponseList := []*appv1.ListAppResponse_App{}
	for _, appRecored := range appRecoredList {
		appResponseList = append(appResponseList, TransformAppRecordToAppResponse(appRecored, usageMetrics, gpuAllTypeMetrics))
	}
	return appResponseList
}

const (
	TimeFormat = "2006-01-02 15:04:05"
)

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(v1.TimeFormat)
}
