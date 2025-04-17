package tests

import (
	appconfig "gitlab.datacanvas.com/aidc/kubevirt-gateway/cmd/app/config"
	configv1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/config/kubevirt_gateway/v1"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kube"
)

func MockConfig() *appconfig.Config {
	kubeConfig := kube.KubeConfiguration{}
	kubevirtConfig := &configv1.KubevirtGatewayConfiguration{
		Server: configv1.Server{
			Port: "8089",
		},
		Http: configv1.Http{
			KubeApiserver: "",
			KubeConfig:    "",
			VncServer:     "",
		},
	}
	appConfig := &appconfig.Config{
		KubeConfig:      kubeConfig,
		ComponentConfig: kubevirtConfig,
	}
	return appConfig
}
