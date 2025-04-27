package datasource

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"vcluster-gateway/pkg/internal/model"
)

// VClusterDBDataSource define DB DataSource Functions
type VClusterDBDataSource interface {
	FindVlusters(tenantType, tenantId string, rootClusterName string, isDeleted int32, permission []string) ([]model.VCluster, error)
	FindVClusterGPUs(vclusterId string) ([]model.VGpu, error)
	FindVClusterStorages(vclusterId string) ([]model.VStorage, error)
	FindVStorageByVClusterIdAndNameType(vclusterId, name, storageType string) ([]*model.VStorage, error)
	GetVClusterById(vClusterId string) (*model.VCluster, error)
	GetVClusterByInstanceId(instanceId string) (*model.VCluster, error)
	GetVStorageByVClusterId(vclusterId string) (*model.VStorage, error)
	CheckVClusterNameExistByTenantId(tenantId string, name string) bool
	CheckVClusterExistById(vClusterId string) bool
	CheckVClusterExistByTenantIdAndVClusterId(vClusterId string, tenantId string) bool
	CheckVStorageExistByVClusterId(vclusterId string) bool
	CheckVGpuExistByVClusterId(vclusterId string) bool
	DeleteVClusterById(id string) error
	DeleteVStorageById(id string) error
	DeleteVGpuById(id string) error
	UpdateVClusterSingle(vcluster *model.VCluster) error
	UpdateVCluster(oldModel *model.VCluster, needUpdate *model.VCluster) error
	UpdateVStorage(oldModel *model.VStorage, needUpdate *model.VStorage) error
	CreateVCluster(vc *model.VCluster) error
	CreateVStorages(storages []*model.VStorage) error
	CreateVGpus(storages []*model.VGpu) error
	CheckVClusterNameExist(vclusterName string) bool
	CheckVClusterNameExistAndDeleted(vclusterName string, isDeleted int) bool
	CheckInstanceIdExist(instanceId string) bool
}

// VClusterKubernetesDataSource define Kubernetes DataSource Functions
type VClusterKubernetesDataSource interface {
	// ***** 操作底层 k8s 的接口 *****
	GetClientSet() *kubernetes.Clientset
	GetSharedInformerFactory() informers.SharedInformerFactory
	GetConfigs() (*clientcmd.ClientConfig, *api.Config)
	Shutdown()

	// ***** 操作 vCluster 的接口 *****
	GenerateVClusterClientSet(vClusterId string) (*kubernetes.Clientset, error)
}
