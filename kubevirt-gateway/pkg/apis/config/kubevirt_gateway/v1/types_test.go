package v1

import (
	"fmt"
	"testing"
)

func TestReadConfFromFile(t *testing.T) {
	var kubevirtConfig = &KubevirtGatewayConfiguration{
		ConfigFilePath:     "../../../../../config.yaml",
		PlatformConfigPath: "../../../../../gcp.yaml",
	}
	if err := kubevirtConfig.ReadConfFromFile(); err != nil {
		t.Errorf("Failed to read config file: %v", err)
	}
	fmt.Printf("kubeconfig: %v", *kubevirtConfig)
}
