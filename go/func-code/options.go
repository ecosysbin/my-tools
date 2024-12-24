package func-code

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
)

func HttpPostRequest(header map[string]string, url string, data []byte) (*http.Response, error) {
	httpReqParam := NewHttpRequstParam(
		WithHeader(header),
		WithMethodUrl(url),
		WithMethod(http.MethodGet),
		WithMethodRetry(3),
		WithData(data),
	)

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

func HttpGetRequest(header map[string]string, url string) (*http.Response, error) {
	httpReqParam := NewHttpRequstParam(
		WithHeader(header),
		WithMethodUrl(url),
		WithMethod(http.MethodGet),
		WithMethodRetry(3),
	)
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

type HttpRequstParam struct {
	url    string
	header map[string]string
	data   []byte
	method string
	retry  int
}

func NewHttpRequstParam(opts ...HttpParamOption) HttpRequstParam {
	var param HttpRequstParam

	for _, opt := range opts {
		opt(&param)
	}
	return param
}

type HttpParamOption func(param *HttpRequstParam)

func WithHeader(header map[string]string) HttpParamOption {
	return func(param *HttpRequstParam) {
		param.header = header
	}
}

func WithData(data []byte) HttpParamOption {
	return func(param *HttpRequstParam) {
		param.data = data
	}
}

func WithMethod(method string) HttpParamOption {
	return func(param *HttpRequstParam) {
		param.method = method
	}
}

func WithMethodUrl(url string) HttpParamOption {
	return func(param *HttpRequstParam) {
		param.url = url
	}
}

func WithMethodRetry(retry int) HttpParamOption {
	return func(param *HttpRequstParam) {
		param.retry = retry
	}
}
