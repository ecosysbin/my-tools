package consts

// globalAnnotations
const (
	AnnotationVClusterSpec       = "vcluster-spec"
	AnnotationVClusterOwner      = "vcluster-owner"
	AnnotationVClusterTenantID   = "vcluster-tenantid"
	AnnotationVClusterDescribe   = "vcluster-describe"
	AnnotationVClusterName       = "vcluster-name"
	AnnotationVClusterInstanceID = "gcp-order-instanceId"
	AnnotationZetyunGPUSpec      = "zetyun-gpu-desc"
	AnnotationGCPManagerBy       = "GCPManagerBy"
)

// isolation.resourceQuota
const (
	ResourceQuotaRequests            = "requests"
	ResourceQuotaLimits              = "limits"
	ResourceQuotaLimit               = "limit"
	ResourceQuotaNVIDIAPrefix        = "nvidia/"
	ResourceQuotaNVIDIAGPUResource   = "nvidia.com/"
	ResourceQuotaManagedBy           = "managed-by"
	ResourceQuotaCpuCores            = "cpu-cores"
	ResourceQuotaDefaultCpuCores     = 6
	ResourceQuotaMemorySize          = "mem-size"
	ResourceQuotaDefaultMemorySize   = 3
	ResourceQuotaStorageClassRequest = ".storageclass.storage.k8s.io/requests.storage"
	ResourceQuotaRequestsCPU         = ResourceQuotaRequests + ".cpu"
	ResourceQuotaLimitsCPU           = ResourceQuotaLimits + ".cpu"
	ResourceQuotaRequestsMemory      = ResourceQuotaRequests + ".memory"
	ResourceQuotaLimitsMemory        = ResourceQuotaLimits + ".memory"
	ResourceQuotaRequestsStorage     = ResourceQuotaRequests + ".storage"
)

// mapServices
const (
	MapServiceApsPrefix  = "aps-"
	MapServiceAimService = "aim-svc"
)

const (
	ApsServerStorageSize    = 50
	VClusterETCDStorageSize = 5
)

const (
	DefaultResourceType            = "vcluster"
	LabelsGCPCollector             = "gcp.com"
	LabelsGCPCollectorInstanceID   = LabelsGCPCollector + "/instance-id"
	LabelsGCPCollectorResourceType = LabelsGCPCollector + "/resource-type"

	PodLabelsGCPCollector             = "gcp.com"
	PodLabelsGCPCollectorInstanceID   = PodLabelsGCPCollector + "/instance-id"
	PodLabelsGCPCollectorResourceType = PodLabelsGCPCollector + "/resource-type"
)

const (
	ManageByRAW          = "raw"
	ManageByVcluster     = "vcluster"
	ManageByAPSTraininig = "aps-training"
	ManageByAPSServing   = "aps-serving"
	ManageByAlayaStudio  = "alaya-studio"
)

const (
	ZetyunStorageClassTypeGCP         = "gcp"
	ZetyunStorageClassTypeAPS         = "aps"
	ZetyunStorageClassTypeAlayaStudio = "alaya-studio"
)

var StorageClassTypeMap = map[string]string{
	ManageByAPSServing:   ZetyunStorageClassTypeAPS,
	ManageByAPSTraininig: ZetyunStorageClassTypeAPS,
	ManageByAlayaStudio:  ZetyunStorageClassTypeAlayaStudio,
}

const (
	DomainServing  = "serving"
	DomainTraining = "training"
)

const (
	TaskCreateApsVCluster = "task_vcluster_provision"

	TaskCreateAps      = "task_aps_provision"
	TaskQueryApsStatus = "task_aps_provision_status"

	TaskCreateAlayaStudio      = "task_create_alaya_studio"
	TaskQueryAlayaStudioStatus = "task_query_alaya_studio_status"
)
