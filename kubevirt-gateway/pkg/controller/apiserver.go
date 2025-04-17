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

package controller

import (
	_ "gitlab.datacanvas.com/aidc/kubevirt-gateway/docs"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/middleware"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/version"

	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"
	gcplog "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/log"
	"gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/sessionuuid"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (c *Controller) SetGlobalMiddleware() {
	c.Handler.Use(gin.Logger())

	// // These handlers are before c.Next().
	c.Handler.Use(gcpctx.GCPContextWrapper(sessionuuid.Middleware()))
	c.Handler.Use(gcpctx.GCPContextWrapper(c.ParseUserToken(
		middleware.FilterRoutersOptions(
			[]string{
				"/swagger/",
				"/favicon.ico",
				"/version",
			}))))

	// The logging middleware must be placed after the sessionuuid
	// and parseUserToken middleware, as it requires the previous
	// two middleware to set the username and sessionuuid.
	c.Handler.Use(gcpctx.GCPContextWrapper(gcplog.Middleware()))

	// These handlers are after c.Next().
	// This method must be the last one.
	c.Handler.Use(gin.Recovery())
}

// Register all Handlers
func (c *Controller) registerHandlers() {
	c.SetGlobalMiddleware()

	// Swagger API
	c.Handler.GET("/kvm/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Version Handler
	c.Handler.GET("/kvm/v1/version", gcpctx.GCPContextWrapper(version.VersionInformation))
	c.registerProfileHandler()
	c.registerVirtualServerHandler()
	c.registerVncProxyHandler()
}

// Register profile Handler
func (c *Controller) registerProfileHandler() {
	// 查询概览信息
	c.Handler.GET("/kvm/v1/product-profiles", gcpctx.GCPContextWrapper(c.ListProductProfiles))
	c.Handler.GET("/kvm/v1/image-profiles", gcpctx.GCPContextWrapper(c.ListImageProfiles))
	c.Handler.GET("/kvm/v1/storage-profiles", gcpctx.GCPContextWrapper(c.ListStorageProfiles))
	// 查询产品可用数量
	c.Handler.GET("/kvm/v1/product/available", gcpctx.GCPContextWrapper(c.Available))
}

// Register virtualserver Handler
func (c *Controller) registerVirtualServerHandler() {
	// 创建虚拟机
	c.Handler.POST("/kvm/v1/virtualserver", gcpctx.GCPContextWrapper(c.CreateOSMVirtualServer))
	// 创建虚拟机
	c.Handler.POST("/kvm/v1/bsm-virtualserver", gcpctx.GCPContextWrapper(c.CreateBSMVirtualServer))
	// 获取虚拟机列表
	c.Handler.GET("/kvm/v1/virtualservers", gcpctx.GCPContextWrapper(c.ListVirtualServers))
	// 查询虚拟机
	c.Handler.GET("/kvm/v1/virtualservers/:instanceId", gcpctx.GCPContextWrapper(c.GetVirtualServer))
	// 删除虚拟机
	c.Handler.DELETE("/kvm/v1/virtualserver/:instanceId", gcpctx.GCPContextWrapper(c.DeleteVirtualServer))
	// 重启虚拟机
	c.Handler.GET("/kvm/v1/virtualserver/:instanceId/restart", gcpctx.GCPContextWrapper(c.RestartVirtualServer))
	// 停止虚拟机
	c.Handler.GET("/kvm/v1/virtualserver/:instanceId/stop", gcpctx.GCPContextWrapper(c.StopVirtualServer))
	// 启动虚拟机
	c.Handler.GET("/kvm/v1/virtualserver/:instanceId/start", gcpctx.GCPContextWrapper(c.StartVirtualServer))
}

func (c *Controller) registerVncProxyHandler() {
	c.Handler.Any("/vnc/v1/vnc_lite.html", gcpctx.GCPContextWrapper(c.VncProxy))
	c.Handler.Any("/vnc/v1/core/*matepath", gcpctx.GCPContextWrapper(c.VncProxy))
	c.Handler.Any("/vnc/v1/vendor/*matepath", gcpctx.GCPContextWrapper(c.VncProxy))
	c.Handler.Any("/vnc/v1/k8s/apis/subresources.kubevirt.io/v1alpha3/namespaces/virtualserver/virtualmachineinstances/:instanceId/vnc", gcpctx.GCPContextWrapper(c.VncProxyVerify))
}
