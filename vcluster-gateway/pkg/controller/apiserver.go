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
	_ "vcluster-gateway/docs"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

// func (c *Controller) SetGlobalMiddleware() {
// 	c.Handler.Use(gin.Logger())

// 	// // These handlers are before c.Next().
// 	c.Handler.Use(gcpctx.GCPContextWrapper(sessionuuid.Middleware()))
// 	c.Handler.Use(gcpctx.GCPContextWrapper(c.ParseUserToken(
// 		middleware.FilterRoutersOptions(
// 			[]string{
// 				"/swagger/",
// 				"/favicon.ico",
// 				"/version",
// 			}))))

// 	// The logging middleware must be placed after the sessionuuid
// 	// and parseUserToken middleware, as it requires the previous
// 	// two middleware to set the username and sessionuuid.
// 	// c.Handler.Use(gcpctx.GCPContextWrapper(gcplog.Middleware()))

// 	// These handlers are after c.Next().
// 	// This method must be the last one.
// 	c.Handler.Use(gin.Recovery())
// }

// Register all Handlers
func (c *Controller) registerHandlers() {
	// Swagger API
	// c.Handler.GET("/kvm/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Version Handler
	// c.Handler.GET("/kvm/v1/version", gcpctx.GCPContextWrapper(version.VersionInformation))
	c.Router.GET("/kvm/v1/product-profiles", c.vClusterHandler.ListProductProfiles)
	// c.registerProfileHandler()
}

// // Register profile Handler
// func (c *Controller) registerProfileHandler() {
// 	// 查询概览信息

// }
