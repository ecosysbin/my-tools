//
// Copyright 2023 The Zetyun.GCP Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package response

import (
	"net/http"

	gcpctx "gitlab.datacanvas.com/AlayaNeW/OSM/gokit/gin/context"
)

type GCPResponse interface {
	Code() int
	Msg() string
}

type gcpResponse struct {
	code int
	msg  string
}

func (g *gcpResponse) Code() int {
	return g.code
}

func (g *gcpResponse) Msg() string {
	return g.msg
}

var (
	SuccessGCPResponse = &gcpResponse{
		code: 200,
		msg:  "Execution succeeded",
	}

	ErrAccessDeny = &gcpResponse{
		code: 142001,
		msg:  "Access Deny",
	}

	ErrJwtVerify = &gcpResponse{
		code: 142002,
		msg:  "Jwt Token Verify Error",
	}

	ErrCreateClusterError = &gcpResponse{
		code: 142003,
		msg:  "Create Cluster Error",
	}

	ErrVClusterSelect = &gcpResponse{
		code: 142004,
		msg:  "Select cluster error",
	}

	ErrClusterDelete = &gcpResponse{
		code: 142005,
		msg:  "Delete cluster error",
	}

	ErrVClusterToken = &gcpResponse{
		code: 142006,
		msg:  "Request token error",
	}

	ErrBodyVerify = &gcpResponse{
		code: 142007,
		msg:  "request body Verify error",
	}

	ErrVClusterEvent = &gcpResponse{
		code: 142008,
		msg:  "Get cluster event error",
	}

	ErrVClusterNS = &gcpResponse{
		code: 142009,
		msg:  "Get cluster namespace error",
	}

	ErrVClusterDeploy = &gcpResponse{
		code: 142010,
		msg:  "Get cluster deploy error",
	}

	ErrVClusterStatefulSet = &gcpResponse{
		code: 142011,
		msg:  "Get cluster statefulset error",
	}
	ErrVClusterIngress = &gcpResponse{
		code: 142012,
		msg:  "Get cluster ingress error",
	}
	ErrVClusterPod = &gcpResponse{
		code: 142013,
		msg:  "Get cluster pod error",
	}
	ErrVClusterSecret = &gcpResponse{
		code: 142014,
		msg:  "Get cluster secret error",
	}
	ErrVClusterService = &gcpResponse{
		code: 142015,
		msg:  "Get cluster service error",
	}
	ErrVClusterConfigmap = &gcpResponse{
		code: 142016,
		msg:  "Get cluster configmap error",
	}

	ErrVClusterList = &gcpResponse{
		code: 142017,
		msg:  "Get cluster resource list error",
	}

	ErrResourceNoExist = &gcpResponse{
		code: 142018,
		msg:  "Get cluster resource no exist ",
	}

	ErrVClusterStructure = &gcpResponse{
		code: 142019,
		msg:  "Cluster structure assemble error",
	}

	ErrCreateClusterWorkflowError = &gcpResponse{
		code: 142020,
		msg:  "Create Cluster Workflow Error",
	}

	ErrPauseCluster = &gcpResponse{
		code: 142021,
		msg:  "Pause cluster error",
	}

	ErrResumeCluster = &gcpResponse{
		code: 142022,
		msg:  "Resume cluster error",
	}
)

type BaseResponse struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func Response(c *gcpctx.GCPContext, gcpResp GCPResponse, data interface{}) {
	httpStatusCode := http.StatusOK

	switch gcpResp.Code() {
	case 1006:
		httpStatusCode = http.StatusBadRequest
	}

	c.JSON(httpStatusCode, BaseResponse{
		Status: gcpResp.Code(),
		Msg:    gcpResp.Msg(),
		Data:   data,
	})

}

func ResponseYaml(c *gcpctx.GCPContext, gcpResp GCPResponse, data interface{}) {
	httpStatusCode := http.StatusOK

	switch gcpResp.Code() {
	case 1006:
		httpStatusCode = http.StatusBadRequest
	}

	c.YAML(httpStatusCode, BaseResponse{
		Status: gcpResp.Code(),
		Msg:    gcpResp.Msg(),
		Data:   data,
	})

}
