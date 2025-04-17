package utils

func GetVClusterNamespaceName(vclusterId string) string {
	return "vcluster-" + vclusterId
}

func GetInfraVClusterNamespaceName(vclusterId string) string {
	return "infra-vc-" + vclusterId
}

func GetVClusterSecretName(vclusterId string) string {
	return "vc-" + vclusterId
}

func GetVClusterServiceName(vclusterId string) string {
	return vclusterId
}

func GetVClusterResourceQuotaName(vclusterId string) string {
	return vclusterId + "-quota"
}
