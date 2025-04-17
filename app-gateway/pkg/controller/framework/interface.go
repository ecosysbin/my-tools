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
	"context"

	"github.com/allegro/bigcache/v3"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
	"github.com/go-resty/resty/v2"
	gcpauthz "gitlab.datacanvas.com/AlayaNeW/OSM/gokit/authz"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/connectrpc/interceptors"
	gcpctx "gitlab.datacanvas.com/AlayaNeW/OSM/gokit/gin/context"
	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"
	"gorm.io/gorm"
	"k8s.io/client-go/informers"

	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
)

type Interface interface {
	CasdoorInterface
	GCPAuthzInterface
	// MiddlewareInterface

	interceptors.Interceptors

	HttpClient() *resty.Client
	ComponentConfig() *v1.VclusterGatewayConfiguration
	AppController() AppControllerInterface
	AppRepo() AppRepoInterface
}

type AppControllerInterface interface {
	GetAppList(params *v1.ListAppParams) (*appv1.ListAppResponse, error)
	GetApp(params *v1.GetAppParams) (*appv1.GetAppResponse, error)
	CreateApp(params *v1.WorkflowAppParams) (*appv1.CreateAppResponse, error)
	UpdateApp(params *v1.WorkflowAppParams) (*appv1.UpdateAppResponse, error)
	DeleteApp(params *v1.WorkflowAppParams) (*appv1.DeleteAppResponse, error)
	PauseApp(params *v1.WorkflowAppParams) (*appv1.PauseAppResponse, error)
	ResumeApp(params *v1.WorkflowAppParams) (*appv1.ResumeAppResponse, error)
	ListAppConfig() ([]repo.AppConfig, error)
	AddConfig(config repo.AppConfig) error
	DeleteConfig(appType string) error
}

type ComponentConfigInterface interface {
	GetServerIP() string
	GetServerPort() string
	GetSwaggerHost() string
	GetSwaggerBasePath() string
	GetTokenKey() string
	GetRespCacheKey() string
	GetVclusterGatewayDefaultCluster() string
	GetVclusterGatewayKubeDaemonHost() string
	GetVclusterGatewayKubeConfig() string
	GetVclusterGatewayDsn() string
	GetVclusterGatewayStorageClass() string
	GetMetricsEndpoint() string
	GetIamAgentEndpoint() string
	GetPlatform() string
	GetRegion() string
	GetCasdoorEndpoint() string
	GetCasdoorClientId() string
	GetCasdoorClientSecret() string
	GetVclusterGatewayConductorEndPoint() string
	GetCasdoorOrganizationName() string
	GetCasdoorApplicationName() string
	GetCasdoorCertificate() string
	GetVClusterAllRootCluster() *v1.AllClusterStruct
	GetVclusterGatewayDefaultImageRegistry() string
	GetVClusterRootCluster(string) *v1.RootClusterStruct
	GetClusterDB() *gorm.DB
	GetRootFactor(string) *informers.SharedInformerFactory
	GetVClusterTable() v1.VCluster
	GetVStorageTable() v1.VStorage
	GetVclusterProductSpec() v1.ProductSpec
	GetVclusterGatewayChatRepo() string
	GetVclusterGatewayConductorApiClient() *client.APIClient
	GetApsProvisioningUrl() string
	GetApsApiKey() string
	GetCephClusterId() string
	GetApsTenantStatusURL() string
	GetApsServingUrl() string
	GetApsTrainingUrl() string
	GetCacheInstance() *bigcache.BigCache
	//GetApsOperationURL() string
	GetVclusterGatewayConductorMetadataClient() *client.MetadataResourceApiService
	GetVclusterGatewayConductorWorkflowExecutor() *executor.WorkflowExecutor
	GetVclusterGatewayConductorTaskRunner() *worker.TaskRunner
}

type MiddlewareInterfaceOptions func(*gcpctx.GCPContext) bool

// type MiddlewareInterface interface {
// 	ParseUserToken(filterOpts ...MiddlewareInterfaceOptions) gcpctx.GCPContextHandlerFunc
// }

type CasdoorInterface interface {
	ParseJwtToken(ctx context.Context, tokenString string) (*casdoorsdk.Claims, error)
	ResetJwkSet(ctx context.Context) error
}

type GCPAuthzInterface interface {
	IsAllowed(policies []gcpauthz.Policy, ctx gcpauthz.Authzable) gcpauthz.AuthResult
	ParsePoliciesFromJSON(data []byte) ([]gcpauthz.Policy, error)
}

type (
	AppRepoInterface interface {
		GetConfig(appType string) (repo.AppConfig, error)
		ListAppConfig() ([]repo.AppConfig, error)
		AddConfig(config repo.AppConfig) error
		DeleteConfig(appType string) error
		ListAll(tenantId string) ([]repo.AppRecord, error)
		ListPageAll(options repo.ListOptions) ([]repo.AppRecord, int64, error)
		// ListDeletedPageAll(options repo.ListOptions) ([]repo.AppRecord, int64, error)
		// ListDeletedAll(tenantId string) ([]repo.AppRecord, error)
		Store(app repo.AppRecord) error
		Update(app repo.AppRecord) error
		GetByAppId(appId string) (repo.AppRecord, error)
	}
)
