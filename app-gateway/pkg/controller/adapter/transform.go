package adapter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"emperror.dev/errors"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
)

func TransformCreateAppParamsToAppRequestData(params *v1.CreateAppParams) (*v1.AppRequestData, error) {
	if len(params.OrderList) == 0 {
		return nil, errors.Errorf("orders is empty")
	}
	var mainOrder v1.Order
	var storageOrders = []v1.Order{}
	for _, order := range params.OrderList {
		if order.OrderInfo.ResourceTypeCode == "CVC" || order.OrderInfo.ResourceTypeCode == "VC" || order.OrderInfo.ResourceTypeCode == "APP" {
			mainOrder = order
		}
		if order.OrderInfo.ResourceTypeCode == "GPU-NODE" || order.OrderInfo.ResourceTypeCode == "GPU-NODE-SUB" || order.OrderInfo.ResourceTypeCode == "CPU-NODE-SUB" {
			mainOrder = order
		}
		if order.OrderInfo.ResourceTypeCode == "FS" || order.OrderInfo.ResourceTypeCode == "FS-NODE-SUB" {
			storageOrders = append(storageOrders, order)
		}
	}
	// instanceId
	instanceId := mainOrder.InstanceId
	// storageList
	storage := []v1.StorageLimit{}
	for _, storageOrder := range storageOrders {
		for _, spec := range storageOrder.InstanceSpec {
			// 兼容1.2.0之前的版本，ParamValue可能是数字，也可能是字符串
			ParamNum, _ := strconv.Atoi(spec.ParamValue)
			// if err != nil {
			// return nil, errors.Errorf("invalid storage limit: %s", spec.ParamValue)
			// }
			storage = append(storage, v1.StorageLimit{
				Limit:  int64(ParamNum),
				FsName: spec.ResourceSpecParamCode,
			})
		}
	}
	// spec
	var cpus float64
	var memory int64
	var gpus []v1.Gpu
	var domain, applicationName string
	for _, spec := range mainOrder.InstanceSpec {
		if spec.ResourceSpecParamCode == "cpu-cores" {
			// 将cpu从string转换为float32
			ParamNum, err := strconv.ParseFloat(spec.ParamValue, 32)
			// ParamNum, err := strconv.Atoi(spec.ParamValue)
			if err != nil {
				return nil, errors.Errorf("invalid cpus limits: %s", spec.ParamValue)
			}
			cpus = ParamNum
		}
		if spec.ResourceSpecParamCode == "mem-size" {
			ParamNum, err := strconv.Atoi(spec.ParamValue)
			if err != nil {
				return nil, errors.Errorf("invalid memory limits: %s", spec.ParamValue)
			}
			memory = int64(ParamNum)
		}
		// 待确认后修改正确字段
		if strings.HasPrefix(spec.ResourceSpecParamCode, "nvidia") {
			ParamNum, err := strconv.Atoi(spec.ParamValue)
			if err != nil {
				return nil, errors.Errorf("invalid gpu limits: %s", spec.ParamValue)
			}
			gpus = append(gpus, v1.Gpu{
				Type:         spec.ParamName,
				Count:        int64(ParamNum),
				ResourceName: spec.ResourceSpecParamCode,
			})
		}
		if spec.ResourceSpecParamCode == "managed-by" {
			domain = spec.ParamValue
		}
		if spec.ResourceSpecParamCode == "application-name" {
			applicationName = spec.ParamValue
		}
	}
	nodeSelector := map[string]string{}
	if mainOrder.NodePoolInstanceId == "" {
		nodeSelector["dc.com/osm.nodepool.type"] = "share"
	} else {
		nodeSelector["dc.com/osm.nodepool.type"] = "exclusive"
		nodeSelector["dc.com/osm.nodepool.tenantId"] = params.TenantId
	}
	return &v1.AppRequestData{
		GcpTenantID:     params.TenantId,
		Name:            params.Name,
		AppId:           params.AppId,
		AidcId:          params.AidcId,
		Desc:            params.Desc,
		Domain:          domain,
		ApplicationName: applicationName,
		User:            params.User,
		InstanceId:      instanceId,
		NodeSelector:    nodeSelector,
		StorageList:     storage,
		Spec: v1.Spec{
			Cpus: cpus,
			Mem:  memory,
			Gpu:  gpus,
		},
	}, nil
}

func checkHttpStatus(status int) error {
	if status != http.StatusOK && status != http.StatusCreated && status != http.StatusAccepted {
		return errors.Errorf("Unexpected HTTP status code: %d", status)
	}
	return nil
}

func TransformConditionsToString(conditions *v1.Conditions) string {
	encode, _ := json.Marshal(conditions)
	return string(encode)
}

func TransformStringToConditions(conditions string) (*v1.Conditions, error) {
	decode := []byte(conditions)
	var result v1.Conditions
	err := json.Unmarshal(decode, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal conditions")
	}
	return &result, nil
}
