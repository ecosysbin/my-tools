package usecase

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v3"

	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	processor "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/processor"
)

func TestProcessVClusterInfo(t *testing.T) {
	// 构造测试数据
	info := &v1.VClusterInfo{
		Username:             "aps001",
		TenantId:             "52a8bad92ae4848c6e06b7b92233e8a9",
		Id:                   "001",
		VClusterId:           "lmlab06201856",
		Context:              "defaultK8sContext",
		Name:                 "lmlab06201856",
		Comment:              "some comment",
		Product:              "some product",
		StorageClass:         "default",
		ChartRepo:            "https://charts.bitnami.com/bitnami",
		DefaultImageRegistry: "registry.datacanvas.com",
		Desc:                 "some desc",
		ManagerBy:            "aps-training",
		Upgrade:              false,
		CephClusterId:        "001",
		InstanceId:           "ff0bca82-6ccf-45b7-b944-79fa5f4a402b",
		OrderDetails: &v1.CreateVClusterParams{
			Orders: []*v1.Order{
				{
					InstanceID: "ff0bca82-6ccf-45b7-b944-79fa5f4a402b",
					InstanceSpecs: []*v1.InstanceSpec{
						{
							ResourceSpecId:        39,
							ResourceSpecCode:      "VC-MIX",
							ResourceSpecParamId:   126,
							ResourceSpecParamCode: "nvidia/mig-h800-pcie-3g.40gb",
							ParamName:             "NVIDIA H800",
							ParamValue:            "1",
						},
						{
							ResourceSpecId:        39,
							ResourceSpecCode:      "VC-MIX",
							ResourceSpecParamId:   125,
							ResourceSpecParamCode: "nvidia/mig-h800-pcie-2g.20gb",
							ParamName:             "NVIDIA H800",
							ParamValue:            "1",
						},
						{
							ResourceSpecId:        39,
							ResourceSpecCode:      "VC-MIX",
							ResourceSpecParamId:   23,
							ResourceSpecParamCode: "managed-by",
							ParamName:             "纳管方",
							ParamValue:            "aps-training",
						},
					},
					ProductID:        11,
					CycleCount:       1,
					OrderType:        1,
					ProductCode:      "APS-MIGtest02",
					ResourceTypeCode: "VC",
					ResourceTypeID:   3,
				},
				{
					InstanceID: "a3ed03df-d8e5-43e8-ada2-581890b60ff2",
					InstanceSpecs: []*v1.InstanceSpec{
						{
							ResourceSpecId:        35,
							ResourceSpecCode:      "FS-HDD",
							ResourceSpecParamId:   19,
							ResourceSpecParamCode: "cephfs-mix",
							ParamName:             "文件存储配额",
							ParamValue:            "64",
						},
					},
					ProductID:        6,
					CycleCount:       1,
					OrderType:        1,
					ResourceTypeCode: "FS",
					ResourceTypeID:   2,
				},
			},
		},
	}
	// 初始化处理器
	p := processor.NewHelmValuesProcessor()

	// 设置默认值
	err := p.ApplySetters(info)
	if err != nil {
		t.Fatalf("Failed to set default values: %v", err)
	}

	// 将 Values 结构体序列化为 YAML 格式
	data, err := yaml.Marshal(p.Values)
	if err != nil {
		t.Fatalf("Failed to marshal values to YAML: %v", err)
	}

	// 打印出 YAML 来检查结果
	// fmt.Println(string(data))

	// 将 YAML 数据保存到文件
	err = ioutil.WriteFile("output.yaml", data, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write YAML to file: %v", err)
	}
}
