package datasource

// 代表 vcluster 集群的状态
type vclusterStatus string

const (
	vclusterStatusRunning  vclusterStatus = "Running"
	vclusterStatusStarting vclusterStatus = "Starting"
	vclusterStatusPaused   vclusterStatus = "Paused"
	vclusterStatusDeleted  vclusterStatus = "Deleted"
	vclusterStatusUnknown  vclusterStatus = "Unknown"
	vclusterStatusFailed   vclusterStatus = "Failed"
)

const (
	appName                  = "vcluster"
	vclusterPrefix           = "vcluster-"
	serviceAccountAnnotation = "vcluster.loft.sh/service-account-name"
	corednsServiceAccount    = "coredns"
	deleteFlag               = 1 // 0 代表未删除，1 代表删除
)

const (
	DefaultClientQPS   = 2000
	DefaultClientBurst = 2000
)
