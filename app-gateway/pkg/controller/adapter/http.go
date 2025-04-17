package adapter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"emperror.dev/errors"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
)

type RespBody struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg,omitempty"`
	Data RespBodyData `json:"data,omitempty"`
}

type RespBodyData struct {
	AppId      string `json:"appId,omitempty"`
	Url        string `json:"url,omitempty"`
	MonitorUrl string `json:"monitorUrl,omitempty"`
	Status     string `json:"status,omitempty"`
}

func AppPostRestRequest(header map[string]string, url string, data []byte, printData v1.HttpLogPrint) (*RespBody, error) {
	log.Infof("Sending provisioning request to URL: %s,Method: %s header: %v,data: %s", printData.Url, printData.Method, printData.Header, printData.Body)
	// log.Infof("Sending provisioning request to URL: %s,Method: %s header: %v,data: %s", printData.Url, printData.Method, printData.Header, string(data))
	resp, err := HttpPostRequest(header, url, data)
	if err != nil {
		return nil, err
	}
	if err := checkHttpStatus(resp.StatusCode); err != nil {
		log.Errorf("Unexpected HTTP status code, URL: %s, data: %s err, %v", string(data))
		return nil, errors.Wrapf(err, "[Unexpected HTTP Status code, URL: %s Code: %d]", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to read response body, URL: %s]", url)
	}

	log.Infof("Received response: %s", string(body))

	var createResp RespBody
	if err = json.Unmarshal(body, &createResp); err != nil {
		return nil, errors.Wrapf(err, "[Failed to unmarshal response body, URL: %s]", url)
	}
	if createResp.Code == 0 && strings.ToLower(createResp.Data.Status) == v1.ACTION_STATUS_FAILED {
		// 兼容aps failed场景，字段在status里面，其他异常使用code非0表示
		createResp.Code = 1
		createResp.Msg = string(body)
	}
	if createResp.Code != 0 && createResp.Msg == "" {
		// 兼容aps failed场景，code不为空，但是错误信息在data里面
		createResp.Msg = string(body)
	}
	log.Infof("Response OpenAps code: %d, data: %v", createResp.Code, createResp.Data)
	return &createResp, nil
}

func HttpPostRequest(header map[string]string, url string, data []byte) (*http.Response, error) {
	httpReqParam := HttpRequstParam{
		method: http.MethodPost,
		url:    url,
		header: header,
		data:   data,
		retry:  3,
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// resp, err := client.Do(req)
	resp, err := httpDoRetryWithNewReq(client, httpReqParam)
	if err != nil {
		log.Errorf("Failed to perform HTTP request, URL: %s, data: %s err, %v", url, string(data), err)
		return nil, errors.Wrapf(err, "[Failed to perform HTTP request, URL: %s]", url)
	}
	if err := checkHttpStatus(resp.StatusCode); err != nil {
		return nil, errors.Wrapf(err, "[Unexpected HTTP Status code, URL: %s, Code: %d]", url, resp.StatusCode)
	}
	return resp, nil
}

func HttpPutRequest(header map[string]string, url string, data []byte) (*RespBody, error) {
	httpReqParam := HttpRequstParam{
		method: http.MethodPut,
		url:    url,
		header: header,
		data:   data,
		retry:  3,
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	log.Infof("Sending provisioning request to URL: %s, data: %s", url, string(data))
	resp, err := httpDoRetryWithNewReq(client, httpReqParam)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to perform HTTP request, URL: %s]", url)
	}
	if err := checkHttpStatus(resp.StatusCode); err != nil {
		return nil, errors.Wrapf(err, "[Unexpected HTTP Status code, URL: %s, Code: %d]", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to read response body, URL: %s]", url)
	}

	log.Infof("Received response: %s", string(body))

	var createResp RespBody
	if err = json.Unmarshal(body, &createResp); err != nil {
		return nil, errors.Wrapf(err, "[Failed to unmarshal response body, URL: %s]", url)
	}
	if createResp.Code == 0 && createResp.Data.Status == v1.ACTION_STATUS_FAILED {
		// 兼容aps failed场景，字段在status里面，其他异常使用code非0表示
		createResp.Code = 1
		createResp.Msg = string(body)
	}
	if createResp.Code != 0 && createResp.Msg == "" {
		// 兼容aps failed场景，code不为空，但是错误信息在data里面
		createResp.Msg = string(body)
	}
	log.Infof("Response OpenAps code: %d, data: %v", createResp.Code, createResp.Data)
	return &createResp, nil
}

func AppGetRequest(header map[string]string, url string) (*RespBody, error) {
	log.Infof("Sending provisioning request to URL: %s", url)
	resp, err := HttpGetRequest(header, url)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to perform HTTP request, URL: %s]", url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to read response body, URL: %s]", url)
	}

	log.Infof("Received response: %s", string(body))

	var createResp RespBody
	if err = json.Unmarshal(body, &createResp); err != nil {
		log.Errorf("Failed to unmarshal response body, URL: %s, err: %v, body: %s", url, err, string(body))
		return nil, errors.Wrapf(err, "[Failed to unmarshal response body, URL: %s]", url)
	}
	if createResp.Code == 0 && createResp.Data.Status == v1.ACTION_STATUS_FAILED {
		// 兼容aps failed场景，字段在status里面，其他异常使用code非0表示
		createResp.Code = 1
		createResp.Msg = string(body)
	}
	if createResp.Code != 0 && createResp.Msg == "" {
		// 兼容aps failed场景，code不为空，但是错误信息在data里面
		createResp.Msg = string(body)
	}
	return &createResp, nil
}

func HttpGetRequest(header map[string]string, url string) (*http.Response, error) {
	httpReqParam := HttpRequstParam{
		method: http.MethodGet,
		url:    url,
		header: header,
		retry:  3,
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	log.Infof("Sending provisioning request to URL: %s", url)
	resp, err := httpDoRetryWithNewReq(client, httpReqParam)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to perform HTTP request, URL: %s]", url)
	}
	if err := checkHttpStatus(resp.StatusCode); err != nil {
		return nil, errors.Wrapf(err, "[Unexpected HTTP Status code, URL: %s, Code: %d]", url, resp.StatusCode)
	}
	return resp, nil
}

func HttpDeleteRequest(header map[string]string, url string) (*RespBody, error) {
	httpReqParam := HttpRequstParam{
		method: http.MethodDelete,
		url:    url,
		header: header,
		data:   nil,
		retry:  3,
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	log.Infof("Sending provisioning request to URL: %s", url)
	resp, err := httpDoRetryWithNewReq(client, httpReqParam)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to perform HTTP request, URL: %s]", url)
	}
	if err := checkHttpStatus(resp.StatusCode); err != nil {
		return nil, errors.Wrapf(err, "[Unexpected HTTP Status code, URL: %s]", url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to read response body, URL: %s]", url)
	}

	log.Infof("Received response: %s", string(body))

	var respBody RespBody
	if err = json.Unmarshal(body, &respBody); err != nil {
		return nil, errors.Wrapf(err, "[Failed to unmarshal response body, URL: %s]", url)
	}
	if respBody.Code == 0 && respBody.Data.Status == v1.ACTION_STATUS_FAILED {
		// 兼容aps failed场景，字段在status里面，其他异常使用code非0表示
		respBody.Code = 1
		respBody.Msg = string(body)
	}
	if respBody.Code != 0 && respBody.Msg == "" {
		// 兼容aps failed场景，code不为空，但是错误信息在data里面
		respBody.Msg = string(body)
	}
	log.Infof("Response OpenAps code: %d, data: %v", respBody.Code, respBody.Data)
	return &respBody, nil
}

type HttpRequstParam struct {
	url    string
	header map[string]string
	data   []byte
	method string
	retry  int
}

func httpDoRetryWithNewReq(client *http.Client, requestParam HttpRequstParam) (*http.Response, error) {
	var req *http.Request
	var resp *http.Response
	var err error

	for i := 0; i < requestParam.retry; i++ {
		req, err = http.NewRequest(requestParam.method, requestParam.url, bytes.NewBuffer(requestParam.data))
		if err != nil {
			return nil, errors.Wrapf(err, "[Failed to create new HTTP request, URL: %s]", requestParam.url)
		}

		for k, v := range requestParam.header {
			req.Header.Set(k, v)
		}
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}
		log.Warnf("HTTP request failed, retrying... %d, err: %v", i, err)
		time.Sleep(1 * time.Second)
	}
	return resp, err
}
