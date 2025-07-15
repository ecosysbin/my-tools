package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// _forceBytesImport 强制导入 bytes 包，避免条件编译错误
func _forceBytesImport() {
	_ = bytes.Buffer{}
}

// GetUserRequest 对应.api文件中的GetUserRequest类型
type GetUserRequest struct {
	Authorization string `header:"authorization"`
	Name          string `path:"name"`
	Delete        bool   `form:"delete,optional"`
}

// GetUserResponse 对应.api文件中的GetUserResponse类型
type GetUserResponse struct {
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
	Age        string `json:"age"`
}

// AddUserRequest 对应.api文件中的AddUserRequest类型
type AddUserRequest struct {
	Authorization string `header:"authorization"`
	Name          string `json:"name"`
	Age           string `json:"age"`
}

// AddUserResponse 对应.api文件中的AddUserResponse类型
type AddUserResponse struct {
	Message string `json:"message"`
}

// DeleteUserRequest 对应.api文件中的DeleteUserRequest类型
type DeleteUserRequest struct {
	Authorization string `header:"authorization"`
	Name          string `path:"name"`
}

// DeleteUserResponse 对应.api文件中的DeleteUserResponse类型
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// UpdateUserRequest 对应.api文件中的UpdateUserRequest类型
type UpdateUserRequest struct {
	Authorization string `header:"authorization"`
	Name          string `path:"name"`
	Age           string `json:"age"`
}

// UpdateUserResponse 对应.api文件中的UpdateUserResponse类型
type UpdateUserResponse struct {
	Message string `json:"message"`
}

// TestApiClient 是访问TestApi服务的客户端
type TestApiClient struct {
	domain string
	client *http.Client
}

// NewTestApi 创建一个新的TestApiClient实例
func NewTestApi(domain string) *TestApiClient {
	return &TestApiClient{
		domain: domain,
		client: &http.Client{},
	}
}

// AddUserHandler 对应.api文件中的AddUserHandler接口
func (c *TestApiClient) AddUserHandler(ctx context.Context, req AddUserRequest) (*AddUserResponse, error) {
	url := fmt.Sprintf("%s%s", c.domain, "/v1/user")

	// 创建HTTP POST请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reqObj, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// 设置header参数
	reqObj.Header.Set("authorization", req.Authorization)

	// 设置请求体类型
	reqObj.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(reqObj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var response AddUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteUserHandler 对应.api文件中的DeleteUserHandler接口
func (c *TestApiClient) DeleteUserHandler(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error) {
	url := fmt.Sprintf("%s%s", c.domain, "/v1/user/"+req.Name+"")

	// 创建HTTP DELETE请求
	reqObj, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置header参数
	reqObj.Header.Set("authorization", req.Authorization)

	// 发送请求
	resp, err := c.client.Do(reqObj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var response DeleteUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateUserHandler 对应.api文件中的UpdateUserHandler接口
func (c *TestApiClient) UpdateUserHandler(ctx context.Context, req UpdateUserRequest) (*UpdateUserResponse, error) {
	url := fmt.Sprintf("%s%s", c.domain, "/v1/user/"+req.Name+"")

	// 创建HTTP PUT请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reqObj, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// 设置header参数
	reqObj.Header.Set("authorization", req.Authorization)

	// 设置请求体类型
	reqObj.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(reqObj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var response UpdateUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUserHandler 对应.api文件中的GetUserHandler接口
func (c *TestApiClient) GetUserHandler(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	url := fmt.Sprintf("%s%s", c.domain, "/v1/user/"+req.Name+"")

	// 处理form参数（查询字符串）
	query := url.Values{}
	if req.Delete != false {
		query.Add("delete", fmt.Sprintf("%v", req.Delete))
	}
	if len(query) > 0 {
		url += "?" + query.Encode()
	}

	// 创建HTTP GET请求
	reqObj, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置header参数
	reqObj.Header.Set("authorization", req.Authorization)

	// 发送请求
	resp, err := c.client.Do(reqObj)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var response GetUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
