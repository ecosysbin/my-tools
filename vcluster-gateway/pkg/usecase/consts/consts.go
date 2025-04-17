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
	ResourceQuotaDefaultCpuCores     = 16
	ResourceQuotaDefaultMemorySize   = 16
	ResourceQuotaMemorySize          = "mem-size"
	ResourceQuotaDefaultStorage      = "storageClass"
	ResourceEnableServiceExporter    = "sync-serviceexporter"
	ResourceQuotaStorageClassRequest = ".storageclass.storage.k8s.io/requests.storage"
	ResourceQuotaRequestsCPU         = ResourceQuotaRequests + ".cpu"
	ResourceQuotaLimitsCPU           = ResourceQuotaLimits + ".cpu"
	ResourceQuotaRequestsMemory      = ResourceQuotaRequests + ".memory"
	ResourceQuotaLimitsMemory        = ResourceQuotaLimits + ".memory"
	ResourceQuotaRequestsStorage     = ResourceQuotaRequests + ".storage"

	ResourceQuotaServicesLoadBalancers = "services.loadbalancers"
	ResourceQuotaServicesNodeports     = "services.nodeports"

	ResourceQuotaIngressesEnabled = "ingresses.enabled"
)

var ResourceTypeCodeSingleStorageMap = map[string]bool{
	"FS": true,
}

type MultiStorageRequesterArgs string

const (
	StorageTypes       MultiStorageRequesterArgs = "--storage-types"
	OrganizationId     MultiStorageRequesterArgs = "--organization-id"
	VClusterId         MultiStorageRequesterArgs = "--vcluster-id"
	StorageManagerHost MultiStorageRequesterArgs = "--storagemanager-host"
	ApiKey             MultiStorageRequesterArgs = "--apikey"
)

func (m MultiStorageRequesterArgs) ToArgs(value string) string {
	return string(m) + "=" + value
}

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
	VClusterStatusCreating = "Creating"
	VClusterStatusUpdating = "Updating"
	VClusterStatusRunning  = "Running"
	VClusterStatusDeleting = "Deleting"
	VClusterStatusResuming = "Resuming"
	VClusterStatusPausing  = "Pausing"
	VClusterStatusPaused   = "Paused"
	VClusterStatusFailed   = "Failed"
	VClusterStatusDeleted  = "Deleted"
)

const (
	TenantTypePlatform      = "TENANT_TYPE_PLATFORM"
	TenantTypePlatformAlias = "3"

	StatusProcessing = "processing"
	StatusSuccess    = "success"

	ActionCreate = "create"
	ActionUpdate = "update"
	ActionResume = "resume"
	ActionDelete = "delete"
	ActionPause  = "pause"
)

type HelmValueFilePath string

const (
	ControlPlaneAffinityTolerations HelmValueFilePath = "hack/helm_manifests/controlplane-affinity-tolerations.yaml"
	EnableHA                        HelmValueFilePath = "hack/helm_manifests/enableHA.yaml"
	EnableStoragePlugin             HelmValueFilePath = "hack/helm_manifests/enable-storage-plugin.yaml"
	EnableServiceExporter           HelmValueFilePath = "hack/helm_manifests/enable-serviceexporter.yaml"
	Extensions                      HelmValueFilePath = "hack/helm_manifests/extensions.yaml"
)

func (h HelmValueFilePath) String() string {
	return string(h)
}
