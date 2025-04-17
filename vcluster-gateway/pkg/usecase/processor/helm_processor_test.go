package processor_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/consts"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/models"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/processor"
)

func mockVClusterInfo() *v1.VClusterInfo {
	return &v1.VClusterInfo{
		Id:                   "aB3xYz9KlM7n",
		Username:             "aps001",
		InstanceId:           "ff0bca82-6ccf-45b7-b944-79fa5f4a402b",
		Desc:                 "this is a desc",
		Name:                 "lmlab06201856",
		TenantId:             "52a8bad92ae4848c6e06b7b92233e8a9",
		ManagerBy:            "this is a manager",
		CephClusterId:        "7e4c0eda-02e0-11ef-bedf-e4434b2ce110",
		StorageClass:         "local-path",
		DefaultImageRegistry: "harbor.zetyun.cn/aidc/vcluster/",
		OrderDetails: &v1.CreateVClusterParams{
			Orders: []*v1.Order{
				{
					InstanceID: "ff0bca82-6ccf-45b7-b944-79fa5f4a402b",
					InstanceSpecs: []*v1.InstanceSpec{
						{
							ParamName:             "NVIDIA H800",
							ParamValue:            "1",
							ResourceSpecCode:      "VC-MIX",
							ResourceSpecParamCode: "nvidia/mig-h800-pcie-3g.40gb",
						},
						{
							ParamName:             "NVIDIA H800",
							ParamValue:            "1",
							ResourceSpecCode:      "VC-MIX",
							ResourceSpecParamCode: "nvidia/mig-h800-pcie-2g.20gb",
						},
						{
							ParamName:             "纳管方",
							ParamValue:            "aps-training",
							ResourceSpecCode:      "VC-MIX",
							ResourceSpecParamCode: "managed-by",
						},
						{
							ParamName:             "xx",
							ParamValue:            "4",
							ResourceSpecCode:      "xx",
							ResourceSpecParamCode: "services.loadbalancers",
						},
						{
							ParamName:             "xx",
							ParamValue:            "8",
							ResourceSpecCode:      "xx",
							ResourceSpecParamCode: "services.nodeports",
						},
						{
							ParamName:             "xx",
							ParamValue:            "1",
							ResourceSpecCode:      "xx",
							ResourceSpecParamCode: "ingresses.enabled",
						},
					},
					OrderType:        1,
					ProductCode:      "APS-MIGtest02",
					ResourceTypeCode: "VC",
				},
				{
					InstanceID: "a3ed03df-d8e5-43e8-ada2-581890b60ff2",
					InstanceSpecs: []*v1.InstanceSpec{
						{
							ParamName:             "文件存储配额",
							ParamValue:            "64",
							ResourceSpecCode:      "FS-HDD",
							ResourceSpecParamCode: "cephfs-mix",
						},
					},
					OrderType:        1,
					ResourceTypeCode: "FS",
				},
			},
		},
	}
}

func mockHelmValuesProcessor() *processor.HelmValuesProcessor {
	process := processor.NewHelmValuesProcessor()

	return process
}

func mockHelmValuesProcessorForSetter() *processor.HelmValuesProcessor {
	process := processor.NewHelmValuesProcessor()

	obtainManagerBy := func(info *v1.VClusterInfo) {
		for _, spec := range info.OrderDetails.Orders[0].InstanceSpecs {
			if strings.HasPrefix(spec.ResourceSpecParamCode, consts.ResourceQuotaManagedBy) {
				info.ManagerBy = spec.ParamValue
				break
			}
		}
	}

	process.RegisterPreprocessors(obtainManagerBy)

	return process
}

func TestHelmValuesProcessor_RegisterDefaultSetters(t *testing.T) {
	tests := []struct {
		name                  string
		notExpectSetterNumber int
		expectError           bool
	}{
		{
			name:                  "success",
			notExpectSetterNumber: 0,
			expectError:           false,
		},
	}

	for _, tt := range tests {
		mockProcessor := mockHelmValuesProcessor()
		mockProcessor.RegisterDefaultSetters()

		assert.NotEqual(t, tt.notExpectSetterNumber, len(mockProcessor.Setters))
	}
}

func TestHelmValuesProcessor_RegisterSetters(t *testing.T) {
	type setterFunc = processor.SetterFunc

	tests := []struct {
		name               string
		expectSetterNumber int
		setterFunc         []setterFunc
	}{
		{
			name:               "success with no setters",
			expectSetterNumber: 0,
			setterFunc:         []setterFunc{},
		},
		{
			name:               "success with single setter",
			expectSetterNumber: 1,
			setterFunc: []setterFunc{
				func(*v1.VClusterInfo, *models.Values) {
					// mock1
				},
			},
		},
		{
			name:               "success with multiple setters",
			expectSetterNumber: 3,
			setterFunc: []setterFunc{
				func(*v1.VClusterInfo, *models.Values) {
					// mock1
				},
				func(*v1.VClusterInfo, *models.Values) {
					// mock2
				},
				func(*v1.VClusterInfo, *models.Values) {
					// mock3
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessor()
			mockProcessor.RegisterSetters(tt.setterFunc...)

			assert.Equal(t, tt.expectSetterNumber, len(mockProcessor.Setters))
		})
	}
}

func TestHelmValuesProcessor_RegisterPreprocessors(t *testing.T) {
	type preprocessFunc = processor.PreprocessFunc
	tests := []struct {
		name                    string
		expectPreprocessorCount int
		preprocessors           []preprocessFunc
	}{
		{
			name:                    "no preprocessors",
			expectPreprocessorCount: 0,
			preprocessors:           []preprocessFunc{},
		},
		{
			name:                    "single preprocessor",
			expectPreprocessorCount: 1,
			preprocessors: []preprocessFunc{
				func(info *v1.VClusterInfo) {
					// mock preprocessor 1
				},
			},
		},
		{
			name:                    "multiple preprocessors",
			expectPreprocessorCount: 3,
			preprocessors: []preprocessFunc{
				func(info *v1.VClusterInfo) {
					// mock preprocessor 1
				},
				func(info *v1.VClusterInfo) {
					// mock preprocessor 2
				},
				func(info *v1.VClusterInfo) {
					// mock preprocessor 3
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessor()
			mockProcessor.RegisterPreprocessors(tt.preprocessors...)

			assert.Equal(t, tt.expectPreprocessorCount, len(mockProcessor.Preprocessors))
		})
	}
}

func TestHelmValuesProcessor_ApplySetters(t *testing.T) {
	tests := []struct {
		name                       string
		info                       *v1.VClusterInfo
		preprocessors              []processor.PreprocessFunc
		setters                    []processor.SetterFunc
		expectedNamePrefix         string
		expectDefaultImageRegistry string
		expectError                bool
	}{
		{
			name: "success with preprocessors and setters",
			info: mockVClusterInfo(),
			preprocessors: []processor.PreprocessFunc{func(info *v1.VClusterInfo) {
				info.Name = "processed-" + info.Username
			}},
			setters: []processor.SetterFunc{func(info *v1.VClusterInfo, values *models.Values) {
				values.DefaultImageRegistry = info.DefaultImageRegistry
			}},
			expectedNamePrefix:         "processed-",
			expectDefaultImageRegistry: mockVClusterInfo().DefaultImageRegistry,
			expectError:                false,
		},
		{
			name:                       "success with no preprocessors and no setters",
			info:                       mockVClusterInfo(),
			preprocessors:              []processor.PreprocessFunc{},
			setters:                    []processor.SetterFunc{},
			expectedNamePrefix:         "",
			expectDefaultImageRegistry: "",
			expectError:                false,
		},
		{
			name:                       "error with nil info",
			info:                       nil,
			preprocessors:              []processor.PreprocessFunc{},
			setters:                    []processor.SetterFunc{},
			expectedNamePrefix:         "",
			expectDefaultImageRegistry: "",
			expectError:                true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessor()
			mockProcessor.RegisterSetters(tt.setters...)
			mockProcessor.RegisterPreprocessors(tt.preprocessors...)

			err := mockProcessor.ApplySetters(tt.info)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.info != nil {
					assert.True(t, strings.HasPrefix(tt.info.Name, tt.expectedNamePrefix))
					assert.Equal(t, tt.expectDefaultImageRegistry, mockProcessor.Values.DefaultImageRegistry)
				}
			}
		})
	}
}

func TestHelmValuesProcessor_SetGlobalAnnotations(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/globalAnnotations_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/globalAnnotations_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetGlobalAnnotations)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetResourceQuota(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/resourceQuota_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/resourceQuota_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetResourceQuota)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, strings.Trim(expected, "\n"), strings.Trim(*generateYaml, "\n"))
		})
	}
}

func TestHelmValuesProcessor_SetZetyun(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/zetyun_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/zetyun_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetZetyun)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()
			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetMapServicesVirtual(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/mapServices_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/mapServices_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetMapServicesVirtual)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()
			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetSyncer(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/syncer_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/syncer_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetSyncer)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())
			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetLabels(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/labels_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/labels_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetLabels)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetPodLabels(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/podLabels_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/podLabels_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetPodLabels)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetPlugin(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/plugin_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/plugin_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetPlugin)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetEtcd(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/etcd_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/etcd_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetEtcd)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetDefaultImageRegistry(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/defaultImageRegistry_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/defaultImageRegistry_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetDefaultImageRegistry)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()

			assert.NoError(t, err)
			assert.Equal(t, expected, *generateYaml)
		})
	}
}

func TestHelmValuesProcessor_SetSync(t *testing.T) {
	tests := []struct {
		name         string
		expectedFile string
		expectError  bool
	}{
		{
			name:         "success",
			expectedFile: "assets/sync_001.yaml",
			expectError:  false,
		},
		{
			name:         "error-file-not-found",
			expectedFile: "assets/sync_002.yaml",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := mockHelmValuesProcessorForSetter()
			mockProcessor.RegisterSetters(processor.SetSync)

			_ = mockProcessor.ApplySetters(mockVClusterInfo())

			expectedBytes, err := os.ReadFile(tt.expectedFile)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			expected := string(expectedBytes)
			generateYaml, err := mockProcessor.GenerateYaml()
			assert.NoError(t, err)
			assert.Equal(t, strings.Trim(expected, "\n"), strings.Trim(*generateYaml, "\n"))
		})
	}
}
