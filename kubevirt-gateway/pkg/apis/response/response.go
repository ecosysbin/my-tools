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

	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"
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
		msg:  "",
	}

	ErrAccessDeny = &gcpResponse{
		code: 1001,
		msg:  "Access Deny",
	}
	ErrJwtVerify = &gcpResponse{
		code: 1002,
		msg:  "Jwt Token Verify Error",
	}
	ErrAuthorization = &gcpResponse{
		code: 1003,
		msg:  "Authorization Error",
	}
	ErrJsonUnmarshal = &gcpResponse{
		code: 1004,
		msg:  "JSON Unmarshal Error",
	}
	ErrBindParams = &gcpResponse{
		code: 1006,
		msg:  "Invalid Request Parameters Error",
	}
	ErrStoreVirtualServer = &gcpResponse{
		code: 3002,
		msg:  "store virtualserver to db Error",
	}
	ErrListVirtualServer = &gcpResponse{
		code: 3003,
		msg:  "list virtualserver from db Error",
	}
	ErrCheckRequestParam = &gcpResponse{
		code: 3004,
		msg:  "invalid request param",
	}
	ErrCreateStorageVolume = &gcpResponse{
		code: 3005,
		msg:  "create storage volume Error",
	}
	ErrCreateVirtualServer = &gcpResponse{
		code: 3006,
		msg:  "create VirtualServer Error",
	}
	ErrCreateVirtualServerSshPort = &gcpResponse{
		code: 3006,
		msg:  "create VirtualServer Ssh Port Error",
	}
)

func ErrCreateResponse(msg string) *gcpResponse {
	return &gcpResponse{
		code: 3001,
		msg:  msg,
	}
}

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
