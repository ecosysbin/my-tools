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

package framework

import (
	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/go-resty/resty/v2"
)

type Interface interface {
	CasdoorInterface
	MiddlewareInterface

	HttpClient() *resty.Client
	ComponentConfig() ComponentConfigInterface
}

type ComponentConfigInterface interface {
	GetServerIP() string
	GetServerPort() string
	GetTokenKey() string
	GetRespCacheKey() string
	GetHttpKubeApiserver() string
	GetProducts() map[string]interface{}
	GetHttpKubeConfig() string
	GetHttpVncServer() string
	GetLoggerLogLevel() string
	GetMysqlPoolMax() int
	GetMysqlUrl() string
	GetCasdoorEndpoint() string
	GetCasdoorClientId() string
	GetCasdoorClientSecret() string
	GetCasdoorOrganizationName() string
	GetCasdoorApplicationName() string
	GetCasdoorCertificate() string
	GetPlatformDomain() string
	GetStorageMin() int64
	GetStorageMax() int64
	GetStorageDefault() int64
	GetImagesFilePath() string
	GetProductsFilePath() string
}

type MiddlewareInterfaceOptions func(*gcpctx.GCPContext) bool

type MiddlewareInterface interface {
	ParseUserToken(filterOpts ...MiddlewareInterfaceOptions) gcpctx.GCPContextHandlerFunc
}

type CasdoorInterface interface {
	ParseJwtToken(token string) (*casdoorsdk.Claims, error)
}

// Handlers Interface
type WebHandlerInterface interface {
	ListProfiles(c *gcpctx.GCPContext)
}

type ReverseProxyHandlerInterface interface{}

// 描述信息接口
type ProfilesHandlerInterface interface {
	ListProductProfiles(c *gcpctx.GCPContext)
	ListImageProfiles(c *gcpctx.GCPContext)
	ListStorageProfiles(c *gcpctx.GCPContext)
}

// 虚拟机操作接口
type VirtualServerHandlerInterface interface {
	// 创建虚拟机 osm场景
	CreateOSMVirtualServer(c *gcpctx.GCPContext)
	// 创建虚拟机 bsm场景
	CreateBSMVirtualServer(c *gcpctx.GCPContext)
	// 查询虚拟机列表
	ListVirtualServers(c *gcpctx.GCPContext)
	// 删除虚拟机
	DeleteVirtualServer(c *gcpctx.GCPContext)
	// 重启虚拟机
	RestartVirtualServer(c *gcpctx.GCPContext)
	// 暂停虚拟机
	StopVirtualServer(c *gcpctx.GCPContext)
	// 启动虚拟机
	StartVirtualServer(c *gcpctx.GCPContext)
	// 虚拟机绑定存储卷
	AddVolume(c *gcpctx.GCPContext)
	// 虚拟机卸载存储卷
	RemoveVolume(c *gcpctx.GCPContext)
}
