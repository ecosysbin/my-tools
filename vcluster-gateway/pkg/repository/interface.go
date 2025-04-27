package repository

import (
	"context"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	v1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
	"vcluster-gateway/pkg/internal/model"
)

type VClusterRepository interface {
	GetVClusterIdByInstanceId(instanceId string) (string, error)
	GetVClusterById(id string) (*model.VCluster, error)
	UpdateVCluster(vc *model.VCluster) error
	CheckVClusterNameExistByTenantId(tenantId string, vClusterName string) bool
	GenerateUniqueVClusterId(ctx context.Context) (string, error)
	CheckVClusterExistByTenantId(vClusterId string, tenantId string) bool
	CheckVClusterExistById(vClusterId string) bool
	CheckVClusterNameExist(vclusterName string) bool
	CheckVClusterNameExistAndDeleted(vclusterName string, isDeleted int) bool
	DeleteVClusterDBResources(id string) error
	CreateVClusterRecord(vclusterId string, serverStatus string, params *v1.VClusterInfo) error
	GetRootK8sConfig() (*v1.RootK8sConfig, error)
	GetVClusterClientSet() (*kubernetes.Clientset, error)
	GetVClusterInformerFactory() (informers.SharedInformerFactory, error)
	GetKubeConfig(ctx context.Context, vClusterId string, kubeConnHost string) (*v1.GetVClusterTokenResponse, error)

	CheckInstanceIdExist(instanceId string) bool

	GetVClusterResourceDetails(ctx context.Context, clusterId string) (*vclusterv1.GetVClusterResourceDetailsResponse_Data, error)
	GetVClusterContainerID(ctx context.Context, params *v1.GetVClusterContainerIDRequest) (containerId string, err error)

	DeleteVClusterNamespace(ctx context.Context, clusterId string) error
}
