package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/adapter"
)

const (
	INSTANCE_MANAGER = "*"

	AUTH_ACTION_GET    = "get"
	AUTH_ACTION_LIST   = "list"
	AUTH_ACTION_CREATE = "create"
	AUTH_ACTION_DELETE = "delete"
	AUTH_ACTION_PAUSE  = "pause"
	AUTH_ACTION_RESUME = "resume"
	AUTH_ACTION_UPDATE = "update"
)

func (as *AppServer) Auth(token, action string) (map[string]string, error) {
	iamAgent := as.controller.ComponentConfig().GetIamAgentEndpoint()
	url := fmt.Sprintf("%s/proxy/api/user/allow/resources", iamAgent)
	httpHeader := map[string]string{}
	httpHeader["X-Access-Token"] = token
	httpHeader["Content-Type"] = "application/json"

	iamAuhReq := IamAuthRequest{
		Platform: as.controller.ComponentConfig().GetPlatform(),
		Service:  "vcluster",
		Region:   as.controller.ComponentConfig().GetRegion(),
		Action:   action,
	}
	reqBody, err := json.Marshal(iamAuhReq)
	if err != nil {
		return nil, errors.Wrapf(err, "[marshal iam auth request failed, url: %s]", url)
	}
	resp, err := AuthPostRestRequest(httpHeader, url, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "[auth failed]")
	}
	if !resp.Data.IsAllow {
		return nil, fmt.Errorf("auth failed: %v, auth url: %s", resp.Data, url)
	}
	if action == AUTH_ACTION_CREATE ||
		action == AUTH_ACTION_DELETE ||
		action == AUTH_ACTION_UPDATE ||
		action == AUTH_ACTION_RESUME ||
		action == AUTH_ACTION_PAUSE {
		return nil, nil
	}
	resourceMap := map[string]string{}
	for _, resource := range resp.Data.Resources {
		if !resource.Allow {
			continue
		}
		for _, r := range resource.Resources {
			if r == "*" {
				resourceMap[INSTANCE_MANAGER] = ""
				continue
			}
			resourceMap[checkoutResourceInstance(r)] = ""
		}
	}
	return resourceMap, nil
}

func checkoutResourceInstance(resource string) string {
	// resource: gcp:vcluster:cn-maanshan-a:*:instance/vc1p2iwcshyu
	array := strings.Split(resource, "instance/")
	if len(array) != 2 {
		return ""
	}
	return array[1]
}

func AuthPostRestRequest(httpHeader map[string]string, url string, reqBody []byte) (*RespBody, error) {
	resp, err := adapter.HttpPostRequest(httpHeader, url, reqBody)
	if err != nil {
		return nil, errors.Wrapf(err, "[HttpPostRequest failed, url: %s, body: %s]", url, string(reqBody))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "[Unexpected HTTP status code %d, URL: %s]", resp.StatusCode, url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[Failed to read response body, URL: %s, body: %s]", url, string(reqBody))
	}

	var createResp RespBody
	if err = json.Unmarshal(body, &createResp); err != nil {
		return nil, errors.Wrapf(err, "[Failed to unmarshal response body, URL: %s, data: %s]", url, string(body))
	}
	if createResp.Code != 0 {
		return nil, errors.Wrapf(err, "[Failed to auth, response code is %d data: %s]", createResp.Code, string(body))

	}
	log.Infof("[Auth success, url: %s, data: %s]", url, string(body))
	return &createResp, nil
}

type IamAuthRequest struct {
	Platform string `json:"platform"`
	Service  string `json:"service"`
	Region   string `json:"region"`
	Action   string `json:"action"`
}

type RespBody struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg,omitempty"`
	Data RespBodyData `json:"data,omitempty"`
}

type RespBodyData struct {
	IsAllow   bool       `json:"isAllow"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	UserId    string   `json:"userId"`
	PolicyId  string   `json:"policyId"`
	Allow     bool     `json:"allow"`
	Resources []string `json:"resources"`
}
