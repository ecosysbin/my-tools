package v1

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	log "github.com/sirupsen/logrus"
)

type UpdateVClusterParams struct {
	AppId string

	Orders []*Order

	Logger *log.Logger

	Name             string
	Desc             string
	UserName         string
	TenantId         string
	RawName          string
	RawDesc          string
	IsInit           bool
	EnableHA         bool
	FallbackDns      string
	CustomHelmConfig CustomHelmConfig
}

type OrderDetail = CreateVClusterParams

type CreateVClusterParams struct {
	Orders []*Order

	Logger *log.Logger

	Name             string
	Desc             string
	UserName         string
	TenantId         string
	RawName          string
	RawDesc          string
	IsInit           bool
	EnableHA         bool
	FallbackDns      string
	CustomHelmConfig CustomHelmConfig
}

type Order struct {
	ProductID        int
	CycleCount       int
	Amount           int
	ActualAmount     int
	OrderType        int
	ResourceTypeID   int
	ProductCode      string
	ResourceTypeCode string

	InstanceID         string
	NodePoolInstanceId string
	InstanceSpecs      []*InstanceSpec
}

type InstanceSpec struct {
	ResourceSpecId        int    `json:"resourceSpecId"`
	ResourceSpecCode      string `json:"resourceSpecCode"`
	ResourceSpecParamId   int    `json:"resourceSpecParamId"`
	ResourceSpecParamCode string `json:"resourceSpecParamCode"`
	ParamName             string `json:"paramName"`
	ParamValue            string `json:"paramValue"`
	ParamUnit             int    `json:"paramUnit,omitempty"`
	ParamType             int    `json:"paramType,omitempty"`
}

type CustomHelmConfig struct {
	EnableCustomization bool
	Repo                string
	ValuesContent       string
}

type RootK8sConfig struct {
	RootClientConfig  *clientcmd.ClientConfig
	RootKubeClientSet *kubernetes.Clientset
	RootRawConfig     *api.Config
}

type VClusterInfo struct {
	Username             string
	TenantId             string
	Id                   string
	VClusterId           string
	Context              string
	Name                 string
	Comment              string
	Product              string
	StorageClass         string
	DefaultStorageClass  string
	ChartRepo            string
	DefaultImageRegistry string
	Desc                 string
	InstanceId           string
	NodePoolInstanceId   string
	ManagerBy            string
	Upgrade              bool
	IsInit               bool
	EnableHA             bool
	CephClusterId        string
	StorageManagerHost   string

	Logger       *log.Logger
	OrderDetails *OrderDetail

	CustomHelmConfig CustomHelmConfig
	FallbackDns      string
}

type UpdateVClusterResponse = CreateVClusterResponse

type CreateVClusterResponse struct {
	VClusterId string
}

type DeleteVClusterParams struct {
	Logger *log.Logger

	Id         string
	Username   string
	TenantId   string
	TenantType string
}

type DeleteVClusterResponse struct {
	Message string
}

type PauseVClusterParams struct {
	Logger *log.Logger

	Username   string
	Id         string
	TenantId   string
	TenantType string
}

type PauseVClusterResponse struct {
	Message string
}

type (
	ResumeVClusterParams   = PauseVClusterParams
	ResumeVClusterResponse = PauseVClusterResponse
)

type GetKubeConfigParams struct {
	Logger *log.Logger

	VClusterId   string
	Username     string
	TenantId     string
	TenantType   string
	KubeConnHost string
}

type GetVClusterTokenResponse = KubeConfig

type KubeConfig struct {
	APIVersion     string      `json:"apiVersion" yaml:"apiVersion" `
	Clusters       []Clusters  `json:"clusters" yaml:"clusters"`
	Contexts       []Contexts  `json:"contexts" yaml:"contexts"`
	CurrentContext string      `json:"current-context" yaml:"current-context"`
	Kind           string      `json:"kind" yaml:"kind"`
	Preferences    Preferences `json:"preferences" yaml:"preferences"`
	Users          []Users     `json:"users" yaml:"users"`
}

type Cluster struct {
	InsecureSkipTLSVerify bool   `json:"insecure-skip-tls-verify" yaml:"insecure-skip-tls-verify"`
	Server                string `json:"server" yaml:"server"`
}

type Clusters struct {
	Cluster Cluster `json:"cluster" yaml:"cluster"`
	Name    string  `json:"name" yaml:"name"`
}

type Context struct {
	Cluster string `json:"cluster" yaml:"cluster"`
	User    string `json:"user" yaml:"user"`
}

type Contexts struct {
	Context Context `json:"context" yaml:"context"`
	Name    string  `json:"name" yaml:"name"`
}

type Preferences struct{}

type User struct {
	Token string `json:"token" yaml:"token"`
}

type Users struct {
	Name string `json:"name" yaml:"name"`
	User User   `json:"user" yaml:"user"`
}

type QueryOperateStatusRequest struct {
	Logger *log.Logger

	Username   string
	TenantId   string
	TenantType string

	AppId  string
	Action string
}

type GetVClusterStatusRequest struct {
	Logger *log.Logger

	Username   string
	TenantId   string
	TenantType string

	AppId string
}

type GetVClusterResourceDetailsRequest struct {
	Logger *log.Logger

	Username   string
	TenantId   string
	TenantType string

	AppId string
}

type GetVClusterContainerIDRequest struct {
	Logger *log.Logger

	Username   string
	TenantId   string
	TenantType string

	VClusterId    string
	Namespace     string
	PodName       string
	ContainerName string
}
