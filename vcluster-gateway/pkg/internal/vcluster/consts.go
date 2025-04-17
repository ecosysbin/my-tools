package forkvcluster

const (
	loftChartRepo             = "http://harbor.zetyun.cn:9000/charts/vcluster"
	defaultChartRepo          = loftChartRepo
	defaultChartVersion       = "0.18.1"
	defaultChartName          = "vcluster"
	defaultLocalChartDir      = ""
	defaultDistro             = "k8s"
	defaultCIDR               = ""
	defaultKubernetesVersion  = ""
	defaultCreateNamespace    = true
	defaultDisableIngressSync = false
	defaultUpdateCurrent      = true
	defaultExpose             = false
	defaultExposeLocal        = false
	defaultConnect            = true
	defaultUpgrade            = false
	defaultIsolate            = false
	defaultProject            = ""
	defaultCluster            = ""
	defaultTemplate           = ""
	defaultTemplateVersion    = ""
	defaultParams             = ""
	defaultDisablePro         = true

	defaultHelmBinaryPath = "/usr/local/bin/helm"

	loginText = "\nPlease run:\n * 'vcluster login' to connect to an existing vCluster.Pro instance\n * 'vcluster pro start' to deploy a new vCluster.Pro instance"

	createdByVClusterAnnotation = "vcluster.loft.sh/created"
)

func getHelmBinaryPath() string {
	return defaultHelmBinaryPath
}
