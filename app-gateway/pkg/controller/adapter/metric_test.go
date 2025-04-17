package adapter

import (
	"fmt"
	"testing"
)

func TestQueryMetrics(t *testing.T) {
	server := "10.220.11.136:31284"
	path := "/api/v1/query"
	query := "query=sum%20by(tenant_id)%20(dc_pod_resource_limits%7Bresource%3D%22memory%22%2Cpod%3D~%22.*%22%7D)"

	resp, err := QueryMetrics(nil, server, path, query)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("response: %v", resp)
}

func TestTestCheckoutGpuMetric(t *testing.T) {
	input := `{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": [
				{
					"metric": {
						"gpu_display_name": "NVIDIA-Tesla-P4"
					},
					"value": [
						1725520816.459,
						"1"
					]
				}
			]
		}
	}`
	result, err := TestCheckoutGpuMetric(input)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Printf("result: %v", result)
}
