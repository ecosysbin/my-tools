package kube

import (
	clientset "k8s.io/client-go/kubernetes"
)

func KubeConfig(master, kubeconfig string) (KubeConfiguration, error) {
	config := DefaultKubeOptions
	config.Master = master
	config.Kubeconfig = kubeconfig
	restConfig, err := config.RestConfig()
	if err != nil {
		return config, err
	}
	config.KubeRestConfig = restConfig
	return config, nil
}

func CreateClients(config *KubeConfiguration) (
	clientSet clientset.Interface,
	err error) {

	kubeConfig, err := config.RestConfig()
	if err != nil {
		return
	}

	clientSet, err = clientset.NewForConfig(kubeConfig)
	if err != nil {
		return
	}
	return
}
