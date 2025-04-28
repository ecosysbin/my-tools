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

	"github.com/gin-gonic/gin"
)

type RestResponse interface {
	Code() int
	Msg() string
}

type restResponse struct {
	code int
	msg  string
}

func (g *restResponse) Code() int {
	return g.code
}

func (g *restResponse) Msg() string {
	return g.msg
}

var (
	SuccessGCPResponse = &restResponse{
		code: 200,
		msg:  "",
	}

	ErrAccessDeny = &restResponse{
		code: 1001,
		msg:  "Access Deny",
	}
	ErrJwtVerify = &restResponse{
		code: 1002,
		msg:  "Jwt Token Verify Error",
	}
	ErrAuthorization = &restResponse{
		code: 1003,
		msg:  "Authorization Error",
	}
	ErrJsonUnmarshal = &restResponse{
		code: 1004,
		msg:  "JSON Unmarshal Error",
	}
	ErrBindParams = &restResponse{
		code: 1006,
		msg:  "Invalid Request Parameters Error",
	}
	ErrStoreVirtualServer = &restResponse{
		code: 3002,
		msg:  "store virtualserver to db Error",
	}
	ErrListVirtualServer = &restResponse{
		code: 3003,
		msg:  "list virtualserver from db Error",
	}
	ErrCheckRequestParam = &restResponse{
		code: 3004,
		msg:  "invalid request param",
	}
	ErrCreateStorageVolume = &restResponse{
		code: 3005,
		msg:  "create storage volume Error",
	}
	ErrCreateVirtualServer = &restResponse{
		code: 3006,
		msg:  "create VirtualServer Error",
	}
	ErrCreateVirtualServerSshPort = &restResponse{
		code: 3006,
		msg:  "create VirtualServer Ssh Port Error",
	}
)

func ErrCreateResponse(msg string) *restResponse {
	return &restResponse{
		code: 3001,
		msg:  msg,
	}
}

type BaseResponse struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func Response(c *gin.Context, restResp RestResponse, data interface{}) {
	httpStatusCode := http.StatusOK

	switch restResp.Code() {
	case 1006:
		httpStatusCode = http.StatusBadRequest
	}

	c.JSON(httpStatusCode, BaseResponse{
		Status: restResp.Code(),
		Msg:    restResp.Msg(),
		Data:   data,
	})

}
