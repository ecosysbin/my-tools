package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/consts"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/models"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase/utils"
)

// PreprocessFunc 预处理 VClusterInfo 的函数类型
type PreprocessFunc func(*v1.VClusterInfo)

// SetterFunc 设置 HelmValuesProcessor.Values 的函数
type SetterFunc func(from *v1.VClusterInfo, to *models.Values)

type HelmValuesProcessor struct {
	Logger *log.Logger

	Values        *models.Values
	Preprocessors []PreprocessFunc
	Setters       []SetterFunc
}

func NewHelmValuesProcessor() *HelmValuesProcessor {
	return &HelmValuesProcessor{
		Values: models.NewHelmValues(),
	}
}

func (p *HelmValuesProcessor) SetProcessorLogger(logger *log.Logger) {
	p.Logger = logger
}

// RegisterDefaultSetters 注册默认的 SetterFunc
func (p *HelmValuesProcessor) RegisterDefaultSetters() {
	p.Setters = []SetterFunc{
		SetGlobalAnnotations,
		SetResourceQuota,
		SetZetyun,
		SetMapServicesVirtual,
		SetSyncer,
		SetLabels,
		SetPodLabels,
		SetPlugin,
		SetEtcd,
		SetDefaultImageRegistry,
		SetSync,
	}
}

// RegisterSetters 注册 SetterFunc
func (p *HelmValuesProcessor) RegisterSetters(setters ...SetterFunc) {
	p.Setters = append(p.Setters, setters...)
}

// RegisterPreprocessors 注册预处理函数，因为 set 函数只从 consts.go 中获取常量或者根据 VClusterInfo 变量获取值
// 但 set 函数执行前，需要处理一下 VClusterInfo，添加预处理函数，可以保证 set 执行时，不需要修改 VClusterInfo 入参
func (p *HelmValuesProcessor) RegisterPreprocessors(preprocessors ...PreprocessFunc) {
	p.Preprocessors = append(p.Preprocessors, preprocessors...)
}

// ApplySetters 设置默认值
func (p *HelmValuesProcessor) ApplySetters(info *v1.VClusterInfo) error {
	if info == nil {
		return errors.New("VClusterInfo cannot be nil")
	}

	for _, preprocess := range p.Preprocessors {
		preprocess(info)
	}

	for _, setter := range p.Setters {
		setter(info, p.Values)
	}

	return nil
}

// SetGlobalAnnotations 设置 globalAnnotations 相关参数
func SetGlobalAnnotations(info *v1.VClusterInfo, values *models.Values) {
	marshal := func(v any) string {
		if v == nil {
			return ""
		}
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}

	values.GlobalAnnotations[consts.AnnotationVClusterSpec] = marshal(info.OrderDetails.Orders[0].InstanceSpecs)
	values.GlobalAnnotations[consts.AnnotationVClusterOwner] = info.Username
	values.GlobalAnnotations[consts.AnnotationVClusterInstanceID] = info.InstanceId
	values.GlobalAnnotations[consts.AnnotationVClusterDescribe] = info.Desc
	values.GlobalAnnotations[consts.AnnotationVClusterName] = info.Name
	values.GlobalAnnotations[consts.AnnotationVClusterTenantID] = info.TenantId
	values.GlobalAnnotations[consts.AnnotationGCPManagerBy] = info.ManagerBy

	instanceSpecs := info.OrderDetails.Orders[0].InstanceSpecs
	gpuDescAnnotationSet := make(map[string]struct{})
	combineAnnotations := make([]string, 0)

	for _, spec := range instanceSpecs {
		if strings.HasPrefix(spec.ResourceSpecParamCode, consts.ResourceQuotaNVIDIAPrefix) {
			if _, exist := gpuDescAnnotationSet[spec.ResourceSpecParamCode]; !exist {
				gpuDescAnnotationSet[spec.ResourceSpecParamCode] = struct{}{}
				combineAnnotations = append(combineAnnotations, spec.ResourceSpecParamCode)
			}
		}
	}
	values.GlobalAnnotations[consts.AnnotationZetyunGPUSpec] = strings.Join(combineAnnotations, "+")

	for _, order := range info.OrderDetails.Orders[1:] {
		// orders 列表的第一个 order 是存储的 VCluster 相关配置的，跳过
		for _, spec := range order.InstanceSpecs {
			if spec.ResourceSpecCode != "" && spec.ResourceSpecParamCode != "" {
				key := spec.ResourceSpecCode + "/" + spec.ResourceSpecParamCode + "-" + info.Id
				values.GlobalAnnotations[key] = consts.ResourceQuotaLimit
			}
		}
	}

	return
}

// SetResourceQuota 设置 resourceQuota 相关参数
func SetResourceQuota(info *v1.VClusterInfo, values *models.Values) {
	cutString := func(s, sep string, n int) string {
		parts := strings.Split(s, sep)
		if len(parts) < n {
			return s
		}
		return parts[n]
	}

	values.Isolation.Enabled = true
	values.Isolation.ResourceQuota.Enabled = true

	instanceSpecs := info.OrderDetails.Orders[0].InstanceSpecs
	gpuNumber := 0
	cpuSet, memorySet := false, false
	for _, spec := range instanceSpecs {
		// 添加对 GPU 的资源限制
		if spec.ResourceSpecParamCode == "" {
			continue
		}
		if strings.HasPrefix(spec.ResourceSpecParamCode, consts.ResourceQuotaNVIDIAPrefix) {
			gpuRequestKey := consts.ResourceQuotaRequests + "." + consts.ResourceQuotaNVIDIAGPUResource +
				strings.Split(cutString(spec.ResourceSpecParamCode, "/", 1), ":")[0]

			num, _ := strconv.Atoi(spec.ParamValue)

			values.Isolation.ResourceQuota.Quota[gpuRequestKey] = num

			gpuNumber += num
		} else if strings.HasPrefix(spec.ResourceSpecParamCode, consts.ResourceQuotaNVIDIAGPUResource) {
			gpuRequestKey := consts.ResourceQuotaRequests + "." + spec.ResourceSpecParamCode
			if strings.Contains(spec.ResourceSpecParamCode, ":") {
				gpuRequestKey = consts.ResourceQuotaRequests + "." + strings.Split(spec.ResourceSpecParamCode, ":")[0]
			}

			num, _ := strconv.Atoi(spec.ParamValue)
			values.Isolation.ResourceQuota.Quota[gpuRequestKey] = num
		}

		if spec.ResourceSpecParamCode == consts.ResourceQuotaCpuCores {
			cpuCores, _ := strconv.Atoi(spec.ParamValue)
			totalCpuCores := consts.ResourceQuotaDefaultCpuCores + cpuCores

			values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaRequestsCPU] = totalCpuCores
			values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaLimitsCPU] = totalCpuCores

			cpuSet = true
		}

		if spec.ResourceSpecParamCode == consts.ResourceQuotaMemorySize {
			memorySize, _ := strconv.Atoi(spec.ParamValue)
			totalMemorySize := utils.ConvertGiga(consts.ResourceQuotaDefaultMemorySize + memorySize)

			values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaRequestsMemory] = totalMemorySize
			values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaLimitsMemory] = totalMemorySize

			memorySet = true
		}

		// 表示创建的 VCluster 使用默认存储
		if spec.ResourceSpecParamCode == consts.ResourceQuotaDefaultStorage {
			value, err := strconv.Atoi(spec.ParamValue)
			if err != nil {
				info.Logger.Warnf("Error converting default storage value to int: %v", err)
				continue
			}
			if value == 0 {
				// 默认存储容量为 0，跳过
				continue
			}

			// 配置默认存储的容量
			key := info.DefaultStorageClass + "-" + info.Id + consts.ResourceQuotaStorageClassRequest
			values.Isolation.ResourceQuota.Quota[key] = utils.ConvertGiga(spec.ParamValue)
		}

	}

	//for _, order := range info.OrderDetails.Orders[1:] {
	//	// orders 列表的第一个 order 是存储的 VCluster 相关配置的，跳过
	//	for _, spec := range order.InstanceSpecs {
	//		key := spec.ResourceSpecParamCode + "-" + info.Id + consts.ResourceQuotaStorageClassRequest
	//		values.Isolation.ResourceQuota.Quota[key] = utils.ConvertGiga(spec.ParamValue)
	//
	//	}
	//}

	if !cpuSet && !memorySet {
		// 通过 GPU 卡数来计算 CPU 和内存的资源限制
		values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaRequestsCPU] = consts.ResourceQuotaDefaultCpuCores + gpuNumber*8
		values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaLimitsCPU] = consts.ResourceQuotaDefaultCpuCores + gpuNumber*8
		values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaRequestsMemory] = utils.ConvertGiga(consts.ResourceQuotaDefaultMemorySize + gpuNumber*64)
		values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaLimitsMemory] = utils.ConvertGiga(consts.ResourceQuotaDefaultMemorySize + gpuNumber*64)
	}

	// 限制 services.nodeports 和 services.loadbalancers
	for _, spec := range instanceSpecs {
		if spec.ResourceSpecParamCode == consts.ResourceQuotaServicesLoadBalancers {
			num, _ := strconv.Atoi(spec.ParamValue)
			values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaServicesLoadBalancers] = num
		} else if spec.ResourceSpecParamCode == consts.ResourceQuotaServicesNodeports {
			num, _ := strconv.Atoi(spec.ParamValue)
			values.Isolation.ResourceQuota.Quota[consts.ResourceQuotaServicesNodeports] = num
		}
	}

	if info.FallbackDns != "" {
		values.Isolation.NetworkPolicy.FallbackDns = info.FallbackDns
	}

	return
}

// SetZetyun 设置 zetyun 相关参数
func SetZetyun(info *v1.VClusterInfo, values *models.Values) {
	values.Zetyun.StorageClass.ClusterId = info.CephClusterId

	// TODO 这里是否要设置为 true 呢
	// values.Zetyun.StorageClass.Enabled = true

	storageClassType, ok := consts.StorageClassTypeMap[info.ManagerBy]
	if !ok {
		storageClassType = consts.ZetyunStorageClassTypeGCP
	}

	instanceSpecs := info.OrderDetails.Orders[0].InstanceSpecs
	for _, spec := range instanceSpecs {
		if spec.ResourceSpecParamCode == consts.ResourceQuotaDefaultStorage {
			value, err := strconv.Atoi(spec.ParamValue)
			if err != nil {
				info.Logger.Warnf("Error converting default storage value to int: %v", err)
				continue
			}
			if value == 0 {
				// 默认存储容量为 0，跳过
				continue
			}

			// 创建 vcluster 默认存储的 storageClass
			values.Zetyun.StorageClass.List = append(values.Zetyun.StorageClass.List, models.NameTypeNs{
				Name: info.DefaultStorageClass,
				Type: storageClassType,
				Ns:   info.DefaultStorageClass,
			})
		}
	}
	//for _, order := range info.OrderDetails.Orders[1:] {
	//	for _, spec := range order.InstanceSpecs {
	//		values.Zetyun.StorageClass.List = append(values.Zetyun.StorageClass.List, models.NameTypeNs{
	//			Name: spec.ResourceSpecParamCode,
	//			Type: storageClassType,
	//			Ns:   spec.ResourceSpecParamCode,
	//		})
	//	}
	//}
	return
}

// SetMapServicesVirtual 设置 aps 相关参数
func SetMapServicesVirtual(info *v1.VClusterInfo, values *models.Values) {
	if !strings.HasPrefix(info.ManagerBy, consts.MapServiceApsPrefix) {
		return
	}
	from := "vcluster-" + info.Id + "/" + consts.MapServiceAimService
	to := consts.MapServiceAimService
	values.MapServices.FromVirtual = append(values.MapServices.FromVirtual, models.FromTo{From: from, To: to})
	return
}

// SetSyncer 设置 syncer 相关参数
func SetSyncer(info *v1.VClusterInfo, values *models.Values) {
	for i, order := range info.OrderDetails.Orders {
		if i == 0 {
			continue
		}

		for _, spec := range order.InstanceSpecs {
			values.Syncer.Env = append(values.Syncer.Env, models.EnvConfig{Name: "storage-" + strconv.Itoa(i), Value: spec.ResourceSpecParamCode + "-" + info.Id})
		}
	}

	values.Syncer.Env = append(values.Syncer.Env, models.EnvConfig{Name: "instance-id", Value: info.InstanceId})

	return
}

// SetLabels 设置 Labels 相关参数
func SetLabels(info *v1.VClusterInfo, values *models.Values) {
	values.Labels[consts.LabelsGCPCollectorInstanceID] = info.InstanceId
	values.Labels[consts.LabelsGCPCollectorResourceType] = consts.DefaultResourceType
	return
}

// SetPodLabels 设置 PodLabels 相关参数
func SetPodLabels(info *v1.VClusterInfo, values *models.Values) {
	values.PodLabels[consts.PodLabelsGCPCollectorInstanceID] = info.InstanceId
	values.PodLabels[consts.PodLabelsGCPCollectorResourceType] = consts.DefaultResourceType
	return
}

// SetPlugin 设置 Plugin 相关参数
func SetPlugin(info *v1.VClusterInfo, values *models.Values) {
	// TODO：将默认存储，通过 Plugin 同步到 VCluster 中
	//for i, order := range info.OrderDetails.Orders {
	//	if i == 0 {
	//		continue
	//	}
	//
	//	for _, spec := range order.InstanceSpecs {
	//		values.Plugin.Hooks.Env = append(values.Plugin.Hooks.Env, models.NameValue{Name: "storage-" + strconv.Itoa(i), Value: spec.ResourceSpecParamCode + "-" + info.Id})
	//	}
	//}

	// Hooks 的 Env 变量

	idx := 1
	instanceSpecs := info.OrderDetails.Orders[0].InstanceSpecs
	for _, spec := range instanceSpecs {
		if spec.ResourceSpecParamCode == consts.ResourceQuotaDefaultStorage {
			// 表示默认存储
			value, err := strconv.Atoi(spec.ParamValue)
			if err != nil {
				info.Logger.Warnf("Error converting default storage value to int: %v", err)
				continue
			}

			if value == 0 {
				// 默认存储容量为 0，跳过
				continue
			}

			values.Plugin.Hooks.Env = append(values.Plugin.Hooks.Env, models.NameValue{Name: "storage-" + strconv.Itoa(idx), Value: info.DefaultStorageClass + "-" + info.Id})

			idx++
		}
	}

	values.Plugin.Hooks.Env = append(values.Plugin.Hooks.Env, models.NameValue{Name: "instance-id", Value: info.InstanceId})

	if info.IsInit {
		values.Plugin.Hooks.Env = append(values.Plugin.Hooks.Env, models.NameValue{Name: "INITIALIZE", Value: "true"})
	} else {
		values.Plugin.Hooks.Env = append(values.Plugin.Hooks.Env, models.NameValue{Name: "INITIALIZE", Value: "false"})
	}

	var enableStoragePluginFn = func(info *v1.VClusterInfo) bool {
		for _, order := range info.OrderDetails.Orders {
			if consts.ResourceTypeCodeSingleStorageMap[order.ResourceTypeCode] {
				return true
			}
		}
		return false
	}

	if enableStoragePluginFn(info) {
		// multiStorage 插件的 Args 变量
		storageTypes := make([]string, 0)
		for _, order := range info.OrderDetails.Orders {
			if !consts.ResourceTypeCodeSingleStorageMap[order.ResourceTypeCode] {
				continue
			}

			for _, spec := range order.InstanceSpecs {
				storageTypes = append(storageTypes, spec.ResourceSpecParamCode+":"+spec.ParamValue+":"+spec.ResourceSpecCode)
			}
		}

		storageTypesArg := consts.StorageTypes.ToArgs(strings.Join(storageTypes, ","))
		organizationIdArg := consts.OrganizationId.ToArgs(info.TenantId)
		vclusterIdArg := consts.VClusterId.ToArgs(info.Id)
		apiKeyArg := consts.ApiKey.ToArgs(os.Getenv("GCP_AUTHZ_API_KEY"))

		values.Plugin.MultiStorageRequester.Args = []string{storageTypesArg, organizationIdArg, vclusterIdArg, apiKeyArg}

		if info.StorageManagerHost != "" {
			storageManagerHostArg := consts.StorageManagerHost.ToArgs(info.StorageManagerHost)
			values.Plugin.MultiStorageRequester.Args = append(values.Plugin.MultiStorageRequester.Args, storageManagerHostArg)
		}
	}

	return
}

// SetEtcd 设置 Etcd 相关参数
func SetEtcd(info *v1.VClusterInfo, values *models.Values) {
	if info.StorageClass != "" {
		values.Etcd.Storage.ClassName = info.StorageClass
	}
	return
}

// SetDefaultImageRegistry 设置 DefaultImageRegistry 相关参数
func SetDefaultImageRegistry(info *v1.VClusterInfo, values *models.Values) {
	values.DefaultImageRegistry = info.DefaultImageRegistry
	return
}

func SetSync(info *v1.VClusterInfo, values *models.Values) {
	instanceSpecs := info.OrderDetails.Orders[0].InstanceSpecs
	for _, spec := range instanceSpecs {
		if spec.ResourceSpecParamCode == consts.ResourceQuotaIngressesEnabled {
			enable, _ := strconv.Atoi(spec.ParamValue)
			if enable == 1 {
				values.Sync.Ingresses.Enabled = true
			} else {
				values.Sync.Ingresses.Enabled = false
			}
		}
	}
}

// GenerateYaml 将拼接的参数序列化到一个 yaml 字符串中
func (p *HelmValuesProcessor) GenerateYaml() (*string, error) {
	var b bytes.Buffer

	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	if err := encoder.Encode(&p.Values); err != nil {
		return nil, err
	}

	valuesYamlStr := b.String()
	return &valuesYamlStr, nil
}

const valuesFileTemplate = "vcluster-%s-values-file-*.yaml"

// GenerateValuesFile 生成 values.yaml 文件
// TODO: 这里生成的 values.yaml 文件，目前没有删除；后续可以考虑使用 pvc 持久化 或者将生成的文件发送到存储服务
func (p *HelmValuesProcessor) GenerateValuesFile(vclusterId string) (filename string, err error) {
	// 格式化成 yaml
	generateYaml, err := p.GenerateYaml()
	if err != nil {
		return filename, err
	}

	// 抽取写入临时文件的方法
	filename, err = WriteToTempFile(generateYaml, vclusterId, p.Logger)
	if err != nil {
		p.Logger.Errorf("Error writing formatted values YAML to temp file: %v", err)
		return filename, err
	}

	p.Logger.Infof("HelmValuesProcessor, Successfully GenerateValuesFile, filename: %s", filename)

	return filename, nil
}

// WriteToTempFile 将给定的 YAML 内容写入临时文件
func WriteToTempFile(content *string, vclusterId string, logger *log.Logger) (filename string, err error) {
	// 将 Values 写入一个临时文件
	tempFile, err := os.CreateTemp("", fmt.Sprintf(valuesFileTemplate, vclusterId))
	if err != nil {
		logger.Errorf("Error creating temp file for formatted values YAML: %v", err)
		return "", err
	}
	// defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(*content)
	if err != nil {
		logger.Errorf("Error writing formatted values YAML to temp file: %v", err)
		tempFile.Close() // 关闭文件
		return "", err
	}
	tempFile.Close()

	return tempFile.Name(), nil
}
