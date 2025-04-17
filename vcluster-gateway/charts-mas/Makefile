.PHONY: package-chart
package-chart: ## Build package vcluster-k8s chart.
	helm package vcluster-mas --version 0.18.1
	helm repo index .

.PHONY: push-chart
push-chart: package-chart ## Push chart to minio.
	cd charts && mc cp ./index.yaml gcp-vcluster-chart/charts/vcluster/k8s/charts-mas/ --recursive
	cd charts && mc cp ./vcluster-k8s-0.18.1.tgz gcp-vcluster-chart/charts/vcluster/k8s/charts-mas/ --recursive