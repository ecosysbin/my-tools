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
	"go.uber.org/dig"
	"gorm.io/gorm"
	"k8s.io/client-go/informers"

	gcpauthz "gitlab.datacanvas.com/AlayaNeW/OSM/gokit/authz"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/connectrpc/interceptors"
	gcpctx "gitlab.datacanvas.com/AlayaNeW/OSM/gokit/gin/context"

	configv1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
)

type Interface interface {
	CasdoorInterface
	GCPAuthzInterface

	interceptors.Interceptors

	DIContainer() DIContainerInterface

	HttpClient() *resty.Client
	ComponentConfig() *configv1.VclusterGatewayConfiguration

	VClusterController() VClusterControllerInterface
}

type VClusterControllerInterface interface {
	CreateVCluster(params *v1.CreateVClusterParams) (*vclusterv1.CreateVClusterResponse, error)
	UpdateVCluster(params *v1.UpdateVClusterParams) (*vclusterv1.UpdateVClusterResponse, error)
	DeleteVCluster(params *v1.DeleteVClusterParams) (*vclusterv1.DeleteVClusterResponse, error)
	PauseVCluster(params *v1.PauseVClusterParams) (*vclusterv1.PauseVClusterResponse, error)
	ResumeVCluster(params *v1.ResumeVClusterParams) (*vclusterv1.ResumeVClusterResponse, error)
	GetKubeConfig(params *v1.GetKubeConfigParams) (*vclusterv1.GetKubeConfigResponse, error)
	GetKubeConfigBase64(params *v1.GetKubeConfigParams) (*vclusterv1.GetKubeConfigBase64Response, error)
	QueryOperateStatus(params *v1.QueryOperateStatusRequest) (*vclusterv1.QueryOperateStatusResponse, error)
	GetVClusterStatus(params *v1.GetVClusterStatusRequest) (*vclusterv1.GetVClusterStatusResponse, error)
	GetVClusterResourceDetails(params *v1.GetVClusterResourceDetailsRequest) (*vclusterv1.GetVClusterResourceDetailsResponse, error)
	GetVClusterContainerID(params *v1.GetVClusterContainerIDRequest) (*vclusterv1.GetVClusterContainerIDResponse, error)
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
	GetCasdoorEndpoint() string
	GetCasdoorClientId() string
	GetCasdoorClientSecret() string
	GetVclusterGatewayConductorEndPoint() string
	GetCasdoorOrganizationName() string
	GetCasdoorApplicationName() string
	GetCasdoorCertificate() string
	GetVClusterAllRootCluster() *configv1.AllClusterStruct
	GetVclusterGatewayDefaultImageRegistry() string
	GetVClusterRootCluster(string) *configv1.RootClusterStruct
	GetClusterDB() *gorm.DB
	GetRootFactor(string) *informers.SharedInformerFactory
	GetVClusterTable() configv1.VCluster
	GetVStorageTable() configv1.VStorage
	GetVclusterProductSpec() configv1.ProductSpec
	GetVclusterGatewayChatRepo() string
	GetVclusterGatewayConductorApiClient() *client.APIClient
	GetCephClusterId() string
	GetCacheInstance() *bigcache.BigCache
	// GetApsOperationURL() string
	GetVclusterGatewayConductorMetadataClient() *client.MetadataResourceApiService
	GetVclusterGatewayConductorWorkflowExecutor() *executor.WorkflowExecutor
	GetVclusterGatewayConductorTaskRunner() *worker.TaskRunner
}

type MiddlewareInterfaceOptions func(*gcpctx.GCPContext) bool

type MiddlewareInterface interface {
	ParseUserToken(filterOpts ...MiddlewareInterfaceOptions) gcpctx.GCPContextHandlerFunc
}

type CasdoorInterface interface {
	ParseJwtToken(ctx context.Context, tokenString string) (*casdoorsdk.Claims, error)
	ResetJwkSet(ctx context.Context) error
}

type GCPAuthzInterface interface {
	IsAllowed(policies []gcpauthz.Policy, ctx gcpauthz.Authzable) gcpauthz.AuthResult
	ParsePoliciesFromJSON(data []byte) ([]gcpauthz.Policy, error)
}

// Handlers Interface
type WebHandlerInterface interface {
	ListProfiles(c *gcpctx.GCPContext)
}

type ReverseProxyHandlerInterface interface{}

type DIContainerInterface interface {
	Invoke(function interface{}, opts ...dig.InvokeOption) error
}
