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

package middleware

import (
	"strings"

	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"
	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/response"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"
)

var _ framework.MiddlewareInterface = &MiddlewareController{}

type MiddlewareController struct {
	controller framework.Interface
}

func New(controller framework.Interface) *MiddlewareController {
	return &MiddlewareController{
		controller: controller,
	}
}

func (mc *MiddlewareController) ParseUserToken(filterOpts ...framework.MiddlewareInterfaceOptions) gcpctx.GCPContextHandlerFunc {
	return func(c *gcpctx.GCPContext) {
		isSkip := false
		for _, opt := range filterOpts {
			if opt(c) {
				isSkip = true
			}
		}

		if !isSkip {
			token := c.GetHeader(mc.controller.ComponentConfig().GetTokenKey())
			if token == "" {
				log.Error("X-Access-Token is must")
				response.Response(c, response.ErrAccessDeny, nil)
				c.Abort()
				return
			}
			// 内部接口，不需要验证token
			if token == "admin" {
				c.SetUesrName("admin")
				c.Next()
				return
			}
			// If Casdoor is enabled, please uncomment the code below.
			// Alternatively, if you need to parse user tokens,
			// you can implement your own ParseJwtToken method.

			claims, err := mc.controller.ParseJwtToken(token)
			if err != nil {
				log.Infof("authentication failed")
				response.Response(c, response.ErrJwtVerify, nil)
				c.Abort()
				return
			}
			role, ok := claims.Properties["role"]
			if ok {
				log.Infof("User %v role is %v", claims.Name, role)
				if !strings.Contains(role, "bsm-admin") && !strings.Contains(role, "osm-admin") {
					response.Response(c, response.ErrAuthorization, "user role error")
					c.Abort()
					return
				}
			} else {
				log.Infof("User %v no role", claims.Name)
				response.Response(c, response.ErrAuthorization, "user role error")
				c.Abort()
				return
			}
			// // 进行授权判断
			// roleMap := map[string]interface{}{}
			// for _, role := range claims.User.Roles {
			// 	roleMap[role.Name] = nil
			// }
			// _, admin_exist := roleMap[ROLE_ADMIN]
			// _, kvmUser_exist := roleMap[ROLE_KVMUSER]
			// if !admin_exist && !kvmUser_exist {
			// 	// 授权失败
			// 	log.Infof("Authorization failed")
			// 	response.Response(c, response.ErrAuthorization, nil)
			// 	c.Abort()
			// 	return
			// }
			c.SetUesrName(claims.Name)
		}

		c.Next()
	}
}

// const (
// 	ROLE_ADMIN   = "console-admin"
// 	ROLE_KVMUSER = "kvm-user"
// )
