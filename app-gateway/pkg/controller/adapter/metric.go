package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
)

type MetricClient struct {
	Header map[string]string
	Server string
	Path   string
}

func NewMetricClient(Header map[string]string, server, path string) *MetricClient {
	return &MetricClient{
		Header: Header,
		Server: server,
		Path:   path,
	}
}

func (client *MetricClient) AllTypeGpuUsageMetrics(instanceId string) (map[string]float32, error) {
	// query=sum by(gpu_display_name)(dc_pod_resource_limits{resource=~"^gpu.*",pod=~".*",tenant_id="9202acbc-a07b-4624-aab7-3316369648e5"})
	// query=sum by(gpu_display_name)(dc_pod_resource_limits{resource=~"^gpu.*",pod=~".*",instance_id="05e9d8ff-c31e-4a6c-8853-dea8b896a830",is_base_model="1"})
	// query=sum%20by(gpu_display_name)(dc_pod_resource_limits%7Bresource=~%22%5Egpu.*%22,pod=~%22.*%22,instance_id=%2205e9d8ff-c31e-4a6c-8853-dea8b896a830%22,is_base_model=%221%22%7D)
	// sum by(gpu_display_name)(dc_pod_resource_limits{resource=~"^gpu.*",pod=~".*",instance_id="2a36a409-8eee-448d-a4d4-22d701ad00fa", is_base_model=""} or dc_pod_resource_limits{resource=~"^gpu.*",pod=~".*",instance_id="2a36a409-8eee-448d-a4d4-22d701ad00fa", is_base_model="1"})
	// sum%20by(gpu_display_name)(dc_pod_resource_limits%7Bresource=~%22%5Egpu.*%22,pod=~%22.*%22,instance_id=%222a36a409-8eee-448d-a4d4-22d701ad00fa%22,%20is_base_model=%22%22%7D%20or%20dc_pod_resource_limits%7Bresource=~%22%5Egpu.*%22,pod=~%22.*%22,instance_id=%222a36a409-8eee-448d-a4d4-22d701ad00fa%22,%20is_base_model=%221%22%7D)
	query := "query=sum+by%28gpu_display_name%29%28dc_pod_resource_limits%7Bresource%3D%7E%22%5Egpu.*%22%2Cpod%3D%7E%22.*%22%2Cinstance_id%3D%22" + instanceId + "%22%2C+is_base_model%3D%22%22%7D+or+dc_pod_resource_limits%7Bresource%3D%7E%22%5Egpu.*%22%2Cpod%3D%7E%22.*%22%2Cinstance_id%3D%22" + instanceId + "%22%2C+is_base_model%3D%221%22%7D%29"
	data, err := QueryMetrics(client.Header, client.Server, client.Path, query)
	if err != nil {
		return nil, err
	}
	log.Infof("app %s gpu metric data: %v", instanceId, data.Result)
	result := make(map[string]float32)
	for _, item := range data.Result {
		if len(item.Value) != 2 || len(item.Metric) != 1 {
			log.Warnf("invalid gpu metric data: %v", data.Result)
			continue
		}
		metric := item.Metric["gpu_display_name"]
		value := item.Value[1]
		if value == nil {
			continue
		}
		if strValue, ok := value.(string); ok {
			num, err := strconv.ParseFloat(strValue, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse gpu metric value, value %s", strValue)
			}
			result[metric] = float32(num)
		}
	}
	return result, nil
}

func TestCheckoutGpuMetric(input string) (map[string]int64, error) {
	type respBody struct {
		Status string          `json:"code"`
		Data   v1.MetricResult `json:"data"`
	}
	var createResp respBody
	if err := json.Unmarshal([]byte(input), &createResp); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal response body, err %v", err)
	}
	result := make(map[string]int64)
	for _, item := range createResp.Data.Result {
		if len(item.Value) != 2 || len(item.Metric) != 1 {
			log.Warnf("invalid gpu metric data: %v", createResp.Data.Result)
			continue
		}
		metric := item.Metric["gpu_display_name"]
		value := item.Value[1]
		if value == nil {
			continue
		}
		if strValue, ok := value.(string); ok {
			num, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse gpu metric value, value %s", strValue)
			}
			result[metric] = num
		}
	}
	return result, nil
}

func (client *MetricClient) UsageMetrics() (v1.UsageMetrics, error) {
	useageMetrics := v1.UsageMetrics{}
	memUsage, err := client.MemMetrics()
	if err != nil {
		return useageMetrics, err
	}
	useageMetrics.MemUsage = memUsage
	cpuUsage, err := client.CpuMetrics()
	if err != nil {
		return useageMetrics, err
	}
	useageMetrics.CpuUsage = cpuUsage
	storageUsage, err := client.StorageMetrics()
	if err != nil {
		return useageMetrics, err
	}
	useageMetrics.Storage = storageUsage
	gpuUsage, err := client.GpuMetrics("nvidia.com/gpu")
	if err != nil {
		return useageMetrics, err
	}
	useageMetrics.GpuUsage = gpuUsage
	return useageMetrics, nil
}

func (client *MetricClient) MemMetrics() (map[string]float32, error) {
	// query=sum by(tenant_id) (dc_pod_resource_limits{resource="memory",pod=~".*"})
	// query=sum by(instance_id) (dc_pod_resource_limits{resource="memory",pod=~".*"})
	// {
	// 	"status": "success",
	// 	"data": {
	// 		"resultType": "vector",
	// 		"result": [
	// 			{
	// 				"metric": {
	// 					"tenant_id": "02fd3e42-5bae-4411-b4e8-1f6024bc0f18"
	// 				},
	// 				"value": [
	// 					1726730962.945,
	// 					"152471339008"
	// 				]
	// 			}
	// 		]
	// 	}
	// }
	query := "query=sum%20by(instance_id)%20(dc_pod_resource_limits%7Bresource%3D%22memory%22%2Cpod%3D~%22.*%22%7D)"
	data, err := QueryMetrics(client.Header, client.Server, client.Path, query)
	if err != nil {
		return nil, err
	}
	result := make(map[string]float32)
	for _, item := range data.Result {
		if len(item.Value) == 0 {
			continue
		}
		value := item.Value[1]
		if value == nil {
			continue
		}
		// if value, ok := value.(float64); ok {
		// 	result[item.Metric["instance_id"]] = int64(value / (1024 * 1024))
		// }

		if strValue, ok := value.(string); ok {
			num, err := strconv.ParseFloat(strValue, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse gpu metric value, value %s", strValue)
			}
			numStr := fmt.Sprintf("%.3f", num/(1024*1024*1024))
			num, _ = strconv.ParseFloat(numStr, 32)
			result[item.Metric["instance_id"]] = float32(num)
		}
	}
	return result, nil
}

func (client *MetricClient) StorageMetrics() (map[string]float32, error) {
	// query=sum by (tenant_id)( dc_project_dir_filestore_usage_current_total )
	// query=sum by (instance_id)(dc_project_dir_filestore_usage_current_total{project_id="all"})
	query := "query=sum%20by%20(instance_id)(dc_project_dir_filestore_usage_current_total%7Bproject_id%3D%22all%22%7D)"
	data, err := QueryMetrics(client.Header, client.Server, client.Path, query)
	if err != nil {
		return nil, err
	}
	result := make(map[string]float32)
	for _, item := range data.Result {
		if len(item.Value) == 0 {
			continue
		}
		value := item.Value[1]
		if value == nil {
			continue
		}
		if strValue, ok := value.(string); ok {
			num, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse gpu metric value, value %s", strValue)
			}
			numStr := fmt.Sprintf("%.3f", num/1024)
			num, _ = strconv.ParseFloat(numStr, 32)
			result[item.Metric["instance_id"]] = float32(num)
		}
	}
	return result, nil
}

func (client *MetricClient) CpuMetrics() (map[string]float32, error) {
	// query=sum by(tenant_id) (dc_pod_resource_limits{resource="cpu",pod=~".*"})
	// query=sum by(instance_id) (dc_pod_resource_limits{resource="cpu",pod=~".*"})
	// 	{
	//         "status": "success",
	//         "data": {
	//                 "resultType": "vector",
	//                 "result": [{
	//                         "metric": {
	//                                 "tenant_id": "cb142bac-406c-4d56-9c6f-8e15d0b8750d"
	//                         },
	//                         "value": [1724398824.334, "461975552"]
	//                 }, {
	//                         "metric": {
	//                                 "tenant_id": "563f17a7-4838-4436-a4c0-f3b900e723f1"
	//                         },
	//                         "value": [1724398824.334, "654336000"]
	//                 }]
	//         }
	// }
	query := "query=sum%20by(instance_id)%20(dc_pod_resource_limits%7Bresource%3D%22cpu%22%2Cpod%3D~%22.*%22%7D)"
	data, err := QueryMetrics(client.Header, client.Server, client.Path, query)
	if err != nil {
		return nil, err
	}
	result := make(map[string]float32)
	for _, item := range data.Result {
		if len(item.Value) == 0 {
			continue
		}
		value := item.Value[1]
		if value == nil {
			continue
		}
		// if int64Value, ok := value.(float64); ok {
		// 	result[item.Metric["instance_id"]] = int64(int64Value)
		// }

		if strValue, ok := value.(string); ok {
			num, err := strconv.ParseFloat(strValue, 32)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse gpu metric value, value %s", strValue)
			}
			// numStr := fmt.Sprintf("%.3f", num/1024)
			// num, _ = strconv.ParseFloat(numStr, 32)
			result[item.Metric["instance_id"]] = float32(num)
		}
	}
	return result, nil
}

func gpuQUery(gpuType string) string {
	// query=sum by(tenant_id) (dc_pod_resource_limits{gpu_display_name="",resource=~"^gpu.*",pod=~".*"})
	return "query=sum%20by(instance_id)%20(dc_pod_resource_limits%7Bgpu_display_name%3D%22" + gpuType + "%22%2Cresource%3D~%22%5Egpu.*%22%2Cpod%3D~%22.*%22%7D)"
}

func (client *MetricClient) GpuMetrics(gputype string) (map[string]int64, error) {
	query := gpuQUery(gputype)
	data, err := QueryMetrics(client.Header, client.Server, client.Path, query)
	if err != nil {
		return nil, err
	}
	result := make(map[string]int64)
	for _, item := range data.Result {
		if len(item.Value) == 0 {
			continue
		}
		value := item.Value[1]
		if value == nil {
			continue
		}
		if int64Value, ok := value.(float64); ok {
			result[item.Metric["instance_id"]] = int64(int64Value)
		}
	}
	return result, nil
}

func metricUrl(server, path, query string) string {
	if strings.Contains(server, "http") {
		return fmt.Sprintf("%s%s?%s", server, path, query)
	} else {
		return fmt.Sprintf("http://%s%s?%s", server, path, query)
	}
}

func QueryMetrics(header map[string]string, server, path, query string) (*v1.MetricResult, error) {
	url := metricUrl(server, path, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create new HTTP request, URL: %s", url)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to perform HTTP request, URL: %s", url)
	}
	if err := checkHttpStatus(resp.StatusCode); err != nil {
		return nil, errors.Wrapf(err, "Unexpected HTTP status code, URL: %s", url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read response body, URL: %s", url)
	}
	type respBody struct {
		Status string          `json:"code"`
		Data   v1.MetricResult `json:"data"`
	}
	var createResp respBody
	if err = json.Unmarshal(body, &createResp); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal response body, URL: %s", url)
	}
	// log.Infof("url %s response %s", url, string(body))
	return &createResp.Data, nil
}
