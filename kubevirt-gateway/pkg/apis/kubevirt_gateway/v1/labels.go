package v1

const (
	LabelGCPBindInstanceIdKey   = "gcp.com/bind-instance-id"
	LabelGCPBindInstanceNameKey = "gcp.com/bind-instance-name"
	LabelGCPMountPointKey       = "gcp.com/mount-point"
	LabelGCPPurposeKey          = "gcp.com/purpose"
	LabelGCPCreateUserKey       = "gcp.com/create-user"
	LabelGCPCreateAppKey        = "gcp.com/create-app"
	LabelGCPReleaseByInstance   = "gcp.com/release-by-instance"
	LabelGCPCreateByInstance    = "gcp.com/create-by-instance" // 随实例一起创建的磁盘会有这个标签, 出错回滚时使用

	LabelGCPCreateApp     = "gcp-kubevirt-gateway"
	LabelGCPPurposeSystem = "system"
	LabelGCPPurposeData   = "data"

	// 计量计费相关标签
	LabelGCPInstanceIdKey     = "gcp.com/instance-id"
	LabelGCPResourceTypeKey   = "gcp.com/resource-type"
	LabelGCPResourceTypeValue = "kvm"
)
