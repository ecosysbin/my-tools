package adapter

import (
	"testing"

	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
)

func TestTransform(t *testing.T) {
	// TODO: write test cases
	// 创建测试数据

	user := v1.GCPUser{
		UserID:   "userIdValue",
		UserName: "userNameValue",
		NickName: "nickNameValue",
		Email:    "emailValue",
		PhoneNum: "phoneNumValue",
	}
	// 构造 CreateAppParams
	appParams := &v1.CreateAppParams{
		Id:           "appIdValue",
		AppId:        "appIdValue",
		Name:         "appNameValue",
		Desc:         "appDescriptionValue",
		ManagerBy:    "managerValue",
		OrderList:    []v1.Order{mainOrder(), storageOrder()},
		UserName:     "creatorUserNameValue",
		TenantId:     "tenantIdValue",
		VClusterName: "vclusterNameValue",
		VClusterDesc: "vclusterDescriptionValue",
		User:         user,
	}
	_, err := TransformCreateAppParamsToAppRequestData(appParams)
	if err != nil {
		t.Errorf("test failed, err: %v", err)
	}
}

func mainOrder() v1.Order {
	orderInfo := v1.OrderInfo{
		ProductId:        1,
		ProductCode:      "product001",
		OrderType:        1,
		ResourceTypeId:   2,
		ResourceTypeCode: "CVC",
	}
	manageBySpec := v1.InstanceSpec{
		ResourceSpecId:        123,
		ResourceSpecCode:      "VC-CPU",
		ResourceSpecParamId:   456,
		ResourceSpecParamCode: "managed-by",
		ParamName:             "paramName",
		ParamValue:            "alaya-studio",
		ParamUnit:             789,
		ParamType:             10,
	}
	aplicationSpec := v1.InstanceSpec{
		ResourceSpecId:        123,
		ResourceSpecCode:      "VC-CPU",
		ResourceSpecParamId:   456,
		ResourceSpecParamCode: "cpu-cores",
		ParamName:             "paramName",
		ParamValue:            "12.5",
		ParamUnit:             789,
		ParamType:             10,
	}
	cpuSpec := v1.InstanceSpec{
		ResourceSpecId:        123,
		ResourceSpecCode:      "VC-CPU",
		ResourceSpecParamId:   456,
		ResourceSpecParamCode: "mem-size",
		ParamName:             "paramName",
		ParamValue:            "32",
		ParamUnit:             789,
		ParamType:             10,
	}
	memSpec := v1.InstanceSpec{
		ResourceSpecId:        123,
		ResourceSpecCode:      "specCode",
		ResourceSpecParamId:   456,
		ResourceSpecParamCode: "paramCode",
		ParamName:             "paramName",
		ParamValue:            "paramValue",
		ParamUnit:             789,
		ParamType:             10,
	}
	return v1.Order{
		OrderInfo:          orderInfo,
		InstanceId:         "instanceIdValue",
		NodePoolInstanceId: "nodePoolInstanceIdValue",
		InstanceSpec:       []v1.InstanceSpec{cpuSpec, memSpec, manageBySpec, aplicationSpec},
	}
}

func storageOrder() v1.Order {
	orderInfo := v1.OrderInfo{
		ProductId:        1,
		ProductCode:      "PRO-HDDalaya",
		OrderType:        1,
		ResourceTypeId:   2,
		ResourceTypeCode: "FS",
	}

	storage := v1.InstanceSpec{
		ResourceSpecId:        123,
		ResourceSpecCode:      "FS-HDD",
		ResourceSpecParamId:   456,
		ResourceSpecParamCode: "unite-cephfs-hdd",
		ParamName:             "paramName",
		ParamValue:            "50",
		ParamUnit:             789,
		ParamType:             10,
	}
	return v1.Order{
		OrderInfo:          orderInfo,
		InstanceId:         "instanceIdValue",
		NodePoolInstanceId: "nodePoolInstanceIdValue",
		InstanceSpec:       []v1.InstanceSpec{storage},
	}
}
