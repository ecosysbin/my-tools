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

package v1

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
	"github.com/fsnotify/fsnotify"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

//var _ framework.ComponentConfigInterface = &VclusterGatewayConfiguration{}

type VclusterGatewayConfiguration struct {
	ConfigFilePath string `json:"configFilePath" yaml:"configFilePath"`
	EnableWatch    bool   `json:"enableWatch" yaml:"enableWatch"`

	Server          Server           `json:"server" yaml:"server"`
	Mysql           Mysql            `json:"mysql" yaml:"mysql"`
	VclusterGateway VclusterGateway  `json:"vclusterGateway" yaml:"vclusterGateway"`
	Metrics         MetricsServer    `json:"metrics" yaml:"metrics"`
	IamAgent        IamAgent         `json:"iamagent" yaml:"iamagent"`
	Casdoor         Casdoor          `json:"casdoor" yaml:"casdoor"`
	AllCluster      AllClusterStruct `json:"-" yaml:"-"`
	ProductSpec     ProductSpec      `json:"productSpec" yaml:"-"`
	Iam             Iam              `json:"iam" yaml:"iam"`
	Conductor       Conductor        `json:"conductor" yaml:"conductor"`
	Aps             Aps              `json:"aps" yaml:"aps"`
	AlayaStudio     AlayaStudio      `json:"alayaStudio" yaml:"alayaStudio"`
	Ceph            Ceph             `json:"ceph" yaml:"ceph"`
	Cache           Cache            `json:"cache" yaml:"cache"`
}

type MetricResult struct {
	Result []MetricData `json:"result"`
}

type MetricData struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type UsageMetrics struct {
	MemUsage map[string]float32
	CpuUsage map[string]float32
	GpuUsage map[string]int64
	Storage  map[string]float32
}

type MetricsServer struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type Mysql struct {
	PoolMax int    `json:"poolMax" yaml:"poolMax"`
	Url     string `json:"url" yaml:"url"`
}

type Ceph struct {
	ClusterId string `json:"clusterId" yaml:"clusterId"`
}

type Server struct {
	Platform     string  `json:"platform" yaml:"platform" mapstructure:"platform"`
	ServerName   string  `json:"serverName" yaml:"serverName" mapstructure:"serverName"`
	Region       string  `json:"region" yaml:"region" mapstructure:"region"`
	IP           string  `json:"-" yaml:"-" mapstructure:"-"`
	Port         string  `json:"port" yaml:"port"`
	TokenKey     string  `json:"tokenKey" yaml:"tokenKey"`
	RespCacheKey string  `json:"respCacheKey" yaml:"respCacheKey"`
	Swagger      Swagger `json:"swagger" yaml:"swagger"`
}

type Swagger struct {
	Host     string `json:"host" yaml:"host"`
	BasePath string `json:"basePath" yaml:"basePath"`
}

type VclusterGateway struct {
	DefaultCluster       string `json:"defaultCluster" yaml:"defaultCluster"`
	KubeDaemonHost       string `json:"kubeDaemonHost" yaml:"kubeDaemonHost"`
	KubeConfig           string `json:"kubeConfig" yaml:"kubeConfig"`
	Dsn                  string `json:"dsn" yaml:"dsn"`
	StorageClass         string `json:"storageClass" yaml:"storageClass"`
	ChartRepo            string `json:"chartRepo" yaml:"chartRepo"`
	DefaultImageRegistry string `json:"defaultImageRegistry" yaml:"defaultImageRegistry"`
}

type Casdoor struct {
	Endpoint         string `json:"endpoint" yaml:"endpoint"`
	ClientId         string `json:"clientId" yaml:"clientId"`
	ClientSecret     string `json:"clientSecret" yaml:"clientSecret"`
	OrganizationName string `json:"organizationName" yaml:"organizationName"`
	ApplicationName  string `json:"applicationName" yaml:"applicationName"`
	Certificate      string `json:"certificate" yaml:"certificate"`
}

type AllClusterStruct struct {
	Cluster        map[string]RootClusterStruct
	CurrentContext string
	DB             *gorm.DB
	VClusterDB     ClusterDB
	Cache          *bigcache.BigCache
}

type Cache struct {
	LifeTime         int  `yaml:"lifeTime"`
	CleanWindow      int  `yaml:"cleanWindow"`
	Enable           bool `json:"enable" yaml:"enable"`
	HardMaxCacheSize int  `yaml:"hardMaxCacheSize"`
	Shards           int  `yaml:"shards"`
}

type RootClusterStruct struct {
	// 底层 K8s 配置
	RootRawConfig     *api.Config
	RootRestConfig    *rest.Config
	RootKubeClientSet *kubernetes.Clientset
	RootClientConfig  *clientcmd.ClientConfig
	Factory           *informers.SharedInformerFactory

	// VCluster K8s 配置
	VClusterMap map[string]VClusterStruct
}

type ClusterDB struct {
	VCluster VCluster
	VStorage VStorage
	VGpu     VGpu
}

type VCluster struct {
	gorm.Model
	UserName        string         `gorm:"column:user_name;size:60" json:"-" swaggerignore:"true"`
	TenantId        string         `gorm:"column:tenant_id;size:36" json:"-" swaggerignore:"true"`
	VClusterId      string         `gorm:"column:vcluster_id;size:12" json:"id"`
	VClusterName    string         `gorm:"column:vcluster_name;size:30" json:"name"`
	RootClusterName string         `gorm:"column:root_cluster_name;size:30" json:"context" `
	CreateTime      *time.Time     `gorm:"column:create_time" json:"createTime"`
	DeleteTime      *time.Time     `gorm:"column:delete_time" json:"deleteTime"`
	Status          Status         `gorm:"column:status" json:"status"`
	Comment         string         `gorm:"column:comment;size:100" json:"comment"`
	StartTime       *time.Time     `gorm:"column:started_time" json:"startedTime"`
	IsDeleted       int            `gorm:"column:is_deleted;size:1" json:"-" swaggerignore:"true"`
	Namespace       string         `gorm:"column:namespace;size:24" json:"namespace"`
	Product         string         `gorm:"column:product;size:20" json:"product"`
	InstanceSpec    string         `gorm:"column:instance_spec;type:text" json:"instanceSpec"`
	InstanceId      string         `gorm:"column:instance_id;size:36" json:"instanceId"`
	ManageBy        string         `gorm:"column:manage_by;size:16;default:'raw'" json:"manageBy"`
	UtilizationRate Resourcequotas `gorm:"-" json:"utilizationRate"`
	AppHost         string         `gorm:"cloumn:app_host" json:"appHost"`
	AppUrl          string         `gorm:"cloumn:app_url" json:"appURL"`
	ApsURL          string         `gorm:"column:aps_url;size:300" json:"apsURL"`
	AlayaStudioURL  string         `gorm:"column:alaya_studio_url;size:300" json:"alayaStudioURL"`
}

type VStorage struct {
	gorm.Model
	VClusterID       string `gorm:"column:vcluster_id;index;size:12; not null"`
	VStorageType     string `gorm:"column:vstorage_type;index;size:60; not null"`
	VStorageCapacity int    `gorm:"column:vstorage_capacity;size:60; not null"`
	IsDeleted        int    `gorm:"column:is_deleted;size:1;"`
	Name             string `gorm:"column:name;index;size:60;"`
}

type VGpu struct {
	gorm.Model
	ClusterID    string `gorm:"index;size:12;not null"`
	GpuType      string `gorm:"size:50;not null"`
	ResourceName string `gorm:"size:50;not null"`
}

const (
	ACTION_STATUS_SUCCESS = "success"
	ACTION_STATUS_FAILED  = "failed"

	ACTION_CREATE = "create"
	ACTION_UPDATE = "update"
	ACTION_DELETE = "delete"
	ACTION_PAUSE  = "pause"
	ACTION_RESUME = "resume"
)

const (
	APP_STATUS_PENDING       = "Pending"
	APP_STATUS_CREATING      = "Creating"
	APP_STATUS_UPDATING      = "Updating"
	APP_STATUS_RUNNING       = "Running"
	APP_STATUS_CREATE_FAILED = "Create_Failed"
	APP_STATUS_PAUSE_FAILED  = "Pause_Failed"
	APP_STATUS_RESUME_FAILED = "Resume_Failed"
	APP_STATUS_DELETE_FAILED = "Delete_Failed"
	APP_STATUS_UPDATE_FAILED = "Update_Failed"
	APP_STATUS_PAUSING       = "Pausing"
	APP_STATUS_PAUSED        = "Paused"
	APP_STATUS_RESUMING      = "Resuming"
	APP_STATUS_DELETING      = "Deleting"
	App_STATUS_DELETED       = "Deleted"
)

func ActionStatusFailed(action string) string {
	switch action {
	case ACTION_CREATE:
		return APP_STATUS_CREATE_FAILED
	case ACTION_UPDATE:
		return APP_STATUS_UPDATE_FAILED
	case ACTION_DELETE:
		return APP_STATUS_DELETE_FAILED
	case ACTION_PAUSE:
		return APP_STATUS_PAUSE_FAILED
	case ACTION_RESUME:
		return APP_STATUS_RESUME_FAILED
	default:
		return APP_STATUS_CREATE_FAILED
	}
}

func IsFailedStatus(status string) bool {
	return status == APP_STATUS_CREATE_FAILED || status == APP_STATUS_PAUSE_FAILED || status == APP_STATUS_RESUME_FAILED || status == APP_STATUS_DELETE_FAILED
}

type Status string

const (
	StatusRunning Status = "Running"

	StatusStarting Status = "Starting"
	StatusPaused   Status = "Paused"
	StatusDeleted  Status = "Deleted"
	StatusUnknown  Status = "Unknown"
)

type VQuota string

// const (
// 	VQuotaHard VQuota = "hard"
// 	VQuotaUsed VQuota = "used"
// )

type ServiceType string

const (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the cluster IP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"

	// ServiceTypeLoadBalancer means a service will be exposed via an
	// external load balancer (if the cloud provider supports it), in addition
	// to 'NodePort' type.
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"

	// ServiceTypeExternalName means a service consists of only a reference to
	// an external name that kubedns or equivalent will return as a CNAME
	// record, with no exposing or proxying of any pods involved.
	ServiceTypeExternalName ServiceType = "ExternalName"
)

// vcluster结构体，存储vcluster的配置信息，使用时判断name是否为空来进行连接，删除vcluster时清理对应map： map[instance_id]
type VClusterStruct struct {
	VClusterRawConfig     api.Config
	VClusterTokenConfig   api.Config
	VClusterRestConfig    rest.Config
	VClusterKubeClientSet kubernetes.Clientset
	VClusterClientConfig  clientcmd.ClientConfig
	Name                  string
	Namespace             string
	Created               metav1.Time
}

type Namespace struct {
	Name       string      `json:"name"`
	CreateTime time.Time   `json:"createTime"`
	Status     string      `json:"status"`
	Object     interface{} `json:"object"`
}

type Pod struct {
	Name       string      `json:"name"`
	CreateTime time.Time   `json:"createTime"`
	Namespace  string      `json:"namespace"`
	Status     Status      `json:"status"`
	Events     []Event     `json:"events"`
	Restarts   int         `json:"restarts"`
	Node       string      `json:"node"`
	Ready      string      `json:"ready"`
	Object     interface{} `json:"object"`
}

type Service struct {
	Name        string      `json:"name"`
	CreateTime  time.Time   `json:"createTime"`
	ClusterIP   string      `json:"clusterIP"`
	Namespace   string      `json:"namespace"`
	Type        string      `json:"type"`
	ExternalIPs []string    `json:"externalIPs"`
	Ports       string      `json:"ports"`
	Selector    string      `json:"selector"`
	Object      interface{} `json:"object"`
}

type Deployment struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreateTime time.Time   `json:"createTime"`
	Ready      string      `json:"ready"`
	Avaliable  int32       `json:"avaliable"`
	Uptodate   int32       `json:"updatedReplicas"`
	Containers string      `json:"containers"`
	Images     string      `json:"images"`
	Selector   string      `json:"selector"`
	Object     interface{} `json:"object"`
}

type StatefulSet struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreateTime time.Time   `json:"createTime"`
	Ready      string      `json:"ready"`
	Containers string      `json:"containers"`
	Images     string      `json:"images"`
	Selector   string      `json:"selector"`
	Object     interface{} `json:"object"`
}

type Secret struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreateTime time.Time   `json:"createTime"`
	Type       string      `json:"type"`
	Data       int         `json:"data"`
	Object     interface{} `json:"object"`
}

type Configmap struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreateTime time.Time   `json:"createTime"`
	Data       int         `json:"data"`
	Object     interface{} `json:"object"`
}

type Resourcequotas struct {
	Gpu     []Struct         `json:"gpu"`
	Memory  map[VQuota]int64 `json:"memory"`
	Cpu     map[VQuota]int64 `json:"cpu"`
	Storage []Struct         `json:"storage"`
}
type Struct struct {
	Name string `json:"name"`
	Hard int64  `json:"hard"`
	Used int64  `json:"used"`
}

type Pvc struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreateTime time.Time   `json:"createTime"`
	Data       int         `json:"data"`
	Object     interface{} `json:"object"`
}

type Ingress struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreateTime time.Time   `json:"createTime"`
	Class      string      `json:"class"`
	Hosts      string      `json:"hosts"`
	Object     interface{} `json:"object"`
}

type EventV1 struct {
	CreateTime time.Time   `json:"createTime"`
	LastTime   time.Time   `json:"lastTime"`
	Level      string      `json:"level"`
	Reason     string      `json:"reason"`
	Message    string      `json:"message"`
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	Object     interface{} `json:"object"`
}

type Event struct {
	Time      time.Time `json:"time"`
	Level     string    `json:"level"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
}

// 此处为kubeconfig的数据结构
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

type Preferences struct {
}

type User struct {
	Token string `json:"token" yaml:"token"`
}

type Users struct {
	Name string `json:"name" yaml:"name"`
	User User   `json:"user" yaml:"user"`
}

type ProductSpec struct {
	ProductCategories []ProductCategories `json:"product_categories,omitempty"`
	Products          []Products          `json:"products,omitempty"`
	Storage           []Storage           `json:"storage,omitempty"`
}

type ProductCategories struct {
	Code  string `json:"code,omitempty"`
	Value string `json:"value,omitempty"`
	Seq   string `json:"seq,omitempty"`
}
type Configs struct {
	ConfigKey   string `json:"configKey,omitempty"`
	ConfigValue string `json:"configValue,omitempty"`
}
type Products struct {
	Code     string    `json:"code,omitempty"`
	Name     string    `json:"name,omitempty"`
	Category string    `json:"category,omitempty"`
	Configs  []Configs `json:"configs,omitempty"`
}

type Storage struct {
	Min          int    `json:"min,omitempty"`
	Max          int    `json:"max,omitempty"`
	Default      int    `json:"default,omitempty"`
	Name         string `json:"name,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
}

type Iam struct {
	PolicyEndPoint string `json:"policyEndPoint" yaml:"policyEndPoint" mapstructure:"policyEndPoint"`
}

type Conductor struct {
	EndPoint         string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	BatchSize        int    `json:"batchSize" yaml:"batchSize" mapstructure:"batchSize"`
	PollInterval     int    `json:"pollInterval" yaml:"pollInterval" mapstructure:"pollInterval"`
	WorkflowExecutor *executor.WorkflowExecutor
	ApiClient        *client.APIClient
	MetadataClient   *client.MetadataResourceApiService
	TaskRunner       *worker.TaskRunner
}

type Aps struct {
	ProvisioningUrl    string `json:"provisioningUrl" yaml:"provisioningUrl" mapstructure:"provisioningUrl"`
	ApiKey             string `json:"apikey" yaml:"apikey" mapstructure:"apikey"`
	GetTenantStatusURL string `json:"getTenantURL" yaml:"getTenantURL" mapstructure:"getTenantURL"`
	TrainingUrl        string `json:"trainingUrl" yaml:"trainingUrl" mapstructure:"trainingUrl"`
	ServingUrl         string `json:"servingUrl" yaml:"servingUrl" mapstructure:"trainingUrl"`
}

type AlayaStudio struct {
	Host   string `json:"host" yaml:"host"`
	ApiKey string `json:"apiKey" yaml:"apiKey"`
}

type IamAgent struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type AppRequestData struct {
	GcpTenantID     string            `json:"gcpTenantId"`
	Name            string            `json:"name"`
	AppId           string            `json:"appId"`
	AidcId          int64             `json:"aidcId"`
	Desc            string            `json:"desc"`
	Domain          string            `json:"domain"`
	ApplicationName string            `json:"applicationName"`
	InstanceId      string            `json:"instanceId"`
	NodeSelector    map[string]string `json:"nodeSelector"`
	StorageList     []StorageLimit    `json:"storageList"`
	Spec            Spec              `json:"spec"`
	User            GCPUser           `json:"user"`
}

func (ad *AppRequestData) IsValid() error {
	if ad.Name == "" {
		return errors.New("app name is empty")
	}
	if ad.Domain == "" {
		return errors.New("app type is empty")
	}
	if ad.Domain != "vcluster" && ad.Domain != "aps-training" && ad.Domain != "aps-serving" && ad.Domain != "alaya-studio" {
		return fmt.Errorf("app type %s is not supported", ad.Domain)
	}
	if ad.InstanceId == "" {
		return errors.New("instance id is empty")
	}
	return nil
}

type StorageLimit struct {
	Origin string `json:"origin"`
	Limit  int64  `json:"limit"`
	FsName string `json:"fsName"`
}

type Gpu struct {
	Type         string `json:"type"`
	Count        int64  `json:"count"`
	ResourceName string `json:"k8sResourceName"`
}

type GCPUser struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	PhoneNum string `json:"phoneNum"`
}

type Spec struct {
	Cpus       float64 `json:"cpus"`
	Mem        int64   `json:"mem"`
	Gpu        []Gpu   `json:"gpu"`
	Kubeconfig string  `json:"kubeconfig"`
}

type APSTenant struct {
	Code int          `json:"code"`
	Data []APSRespone `json:"data"`
}

type APSRespone struct {
	GcpTenantID string `json:"gcpTenantId"`
	VclusterID  string `json:"vclusterId"`
	TenantID    string `json:"tenantId"`
	Status      string `json:"status"`
	TenantURL   string `json:"tenantURL"`
	Message     string `json:"message"`
	UpdateTime  string `json:"updateTime"`
	AppURL      string `json:"appURL"`
	AppPath     string `json:"appPath"`
	AppBaseURL  string `json:"appBaseURL"`
}

type APSOpen struct {
	Code int        `json:"code"`
	Data APSOpenRes `json:"data"`
}
type APSOpenRes struct {
	GcpTenantId string `json:"gcpTenantId"`
	TenantId    string `json:"tenantId"`
}

type AlayaStudioData struct {
	GcpTenantId     string      `json:"gcpTenantId"`
	VClusterId      string      `json:"vclusterId"`
	StorageList     interface{} `json:"storageList"`
	VClusterSpec    interface{} `json:"vclusterSpec"`
	ApplicationName string      `json:"applicationName"`
}

type RawData struct {
	GcpTenantId     string      `json:"gcpTenantId"`
	VClusterId      string      `json:"vclusterId"`
	StorageList     interface{} `json:"storageList"`
	VClusterSpec    interface{} `json:"vclusterSpec"`
	ApplicationName string      `json:"applicationName"`
}

type CreateAppParams struct {
	Logger       *log.Logger `json:"-"`
	Id           string      `json:"id"`
	AppId        string      `json:"appId"`
	AidcId       int64       `json:"aidcId"`
	Name         string      `json:"name"`
	Desc         string      `json:"desc"`
	ManagerBy    string      `json:"managerBy"`
	OrderList    []Order     `json:"orderlist"`
	UserName     string      `json:"userName"`
	TenantId     string      `json:"tenantId"`
	VClusterName string      `json:"vclusterName"`
	VClusterDesc string      `json:"vclusterDesc"`
	User         GCPUser     `json:"user"`
}

type WorkflowAppParams struct {
	Logger *log.Logger `json:"-"`
	Action string      `json:"action"`
	// Repo framework.RepoIn
	AppConfig      AppConfig         `json:"appConfig"`
	HttpHeader     map[string]string `json:"httpHeader"`
	HttpParams     `json:"params"`
	UpdateAppReq   *appv1.UpdateAppRequest `json:"updateAppReq"`
	CreateAppReq   *appv1.CreateAppRequest `json:"createAppReq"`
	AppRequestData AppRequestData          `json:"appRequestData"`
	AppRecord      repo.AppRecord          `json:"appRecord"`
	Conditions     `json:"conditions"`
}

type Conditions struct {
	Action    string   `json:"action"`
	PreStatus string   `json:"preStatus"`
	Status    string   `json:"status"`
	Reason    string   `json:"reason"`
	Events    []string `json:"events"`
}

type HttpParams map[string]string

func (h HttpParams) Encode() string {
	encodeStr := ""
	if len(h) == 0 {
		return ""
	}
	for k, v := range h {
		encodeStr += fmt.Sprintf("%s=%s&", k, v)
	}
	return encodeStr[0 : len(encodeStr)-1]
}

type HttpLogPrint struct {
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
	Resp   string            `json:"resp"`
}

type AppConfig struct {
	AppType    string `json:"appType"`
	Doamin     string `json:"domain"`
	Url        string `json:"url"`
	KubeConfig string `json:"kubeconfig"`
	VclusterId string `json:"vclusterId"`
}

type Order struct {
	OrderInfo          OrderInfo      `json:"orderInfo"`
	InstanceId         string         `json:"instanceId"`
	NodePoolInstanceId string         `json:"nodePoolInstanceId"`
	InstanceSpec       []InstanceSpec `json:"instanceSpec"`
}

type OrderInfo struct {
	ProductId        int    `json:"productId"`
	ProductCode      string `json:"productCode"`
	CycleCount       int    `json:"cycleCount"`
	Amount           int    `json:"amount"`
	OrderType        int    `json:"orderType"`
	ResourceTypeId   int    `json:"resourceTypeId"`
	ResourceTypeCode string `json:"resourceTypeCode"`
	ActualAmount     int    `json:"actualAmount"`
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

type ListAppParams struct {
	Delete     int32
	Context    string
	Username   string
	TenantType string
	TenantId   string
	Id         string
	Name       string
	InstanceId string
	Page       int
	PageSize   int
	ManageBy   string
	CreateBy   string
	Status     []string
	TenantIds  []string
	CreateTime PeriodTime
	DeleteTime PeriodTime
	Logger     *log.Logger
}

type PeriodTime struct {
	StartTime string
	EndTime   string
}

type GetAppParams struct {
	ReqToken     string
	TenantId     string
	TenantType   string
	AppId        string
	AppConfigMap map[string]AppConfig
	HttpHeader   map[string]string
	Logger       *log.Logger
}

type UpdateAppInstanceStatusData struct {
	TenantId             string `json:"tenantId"`
	InstanceId           string `json:"instanceId"`
	AppInstanceId        string `json:"appInstanceId"`
	Valid                bool   `json:"valid"`
	ProductCategoryValue string `json:"productCategoryValue"`
	Reason               Reason `json:"reason"`
}

type Reason struct {
	Zh string `json:"zh-CN"`
	En string `json:"en-US"`
}

func (lvp *ListAppParams) GetContext() string {
	if lvp.Context == "" {
		return DefaultKubernetesContext
	}

	return lvp.Context
}

func (lvp *ListAppParams) IsDeleted() bool {
	return lvp.Delete == 1
}

type (
	ResumeAppParams = PauseAppParams
)

type PauseAppParams struct {
	Logger *log.Logger

	Username   string
	Ids        []string
	TenantId   string
	TenantType string

	ApsInterfacePrefix string
	ApsApiKey          string
}

type DeleteAppParams struct {
	Logger    *log.Logger
	AppConfig AppConfig      `json:"appConfig"`
	AppRecord repo.AppRecord `json:"appRecord"`
}

func (u VCluster) TableName() string {
	return "vcluster"
}

func (u VStorage) TableName() string {
	return "vstorage"
}

func (u VGpu) TableName() string {
	return "vgpu"
}

func (c *VclusterGatewayConfiguration) GetMysqlPoolMax() int {
	return c.Mysql.PoolMax
}

func (c *VclusterGatewayConfiguration) GetMysqlUrl() string {
	return c.Mysql.Url
}

func (c *VclusterGatewayConfiguration) GetClusterDB() *gorm.DB {
	return c.AllCluster.DB
}

func (c *VclusterGatewayConfiguration) GetRootFactor(k8sContext string) *informers.SharedInformerFactory {
	return c.AllCluster.Cluster[k8sContext].Factory
}

func (c *VclusterGatewayConfiguration) GetVclusterProductSpec() ProductSpec {
	return c.ProductSpec
}

func (c *VclusterGatewayConfiguration) GetVClusterTable() VCluster {
	return c.AllCluster.VClusterDB.VCluster
}

func (c *VclusterGatewayConfiguration) GetVStorageTable() VStorage {
	return c.AllCluster.VClusterDB.VStorage
}

func (c *VclusterGatewayConfiguration) GetVClusterAllRootCluster() *AllClusterStruct {
	return &c.AllCluster
}

func (c *VclusterGatewayConfiguration) GetVClusterRootCluster(k8sContext string) RootClusterStruct {
	return c.AllCluster.Cluster[k8sContext]
}

func (c *VclusterGatewayConfiguration) GetSwaggerHost() string {
	return c.Server.Swagger.Host
}

func (c *VclusterGatewayConfiguration) GetSwaggerBasePath() string {
	return c.Server.Swagger.BasePath
}

func (c *VclusterGatewayConfiguration) GetServerIP() string {
	return c.Server.IP
}

func (c *VclusterGatewayConfiguration) GetServerPort() string {
	return c.Server.Port
}

func (c *VclusterGatewayConfiguration) GetPlatform() string {
	return c.Server.Platform
}

func (c *VclusterGatewayConfiguration) GetRegion() string {
	return c.Server.Region
}

func (c *VclusterGatewayConfiguration) GetTokenKey() string {
	return c.Server.TokenKey
}

func (c *VclusterGatewayConfiguration) GetRespCacheKey() string {
	return c.Server.RespCacheKey
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayDefaultCluster() string {
	return c.VclusterGateway.DefaultCluster
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayKubeDaemonHost() string {
	return c.VclusterGateway.KubeDaemonHost
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayKubeConfig() string {
	return c.VclusterGateway.KubeConfig
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayDsn() string {
	return c.VclusterGateway.Dsn
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayDefaultImageRegistry() string {
	return c.VclusterGateway.DefaultImageRegistry
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayStorageClass() string {
	return c.VclusterGateway.StorageClass
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayChatRepo() string {
	return c.VclusterGateway.ChartRepo
}

func (c *VclusterGatewayConfiguration) GetMetricsEndpoint() string {
	return c.Metrics.Endpoint
}

func (c *VclusterGatewayConfiguration) GetIamAgentEndpoint() string {
	return c.IamAgent.Endpoint
}

func (c *VclusterGatewayConfiguration) GetCasdoorEndpoint() string {
	return c.Casdoor.Endpoint
}

func (c *VclusterGatewayConfiguration) GetCasdoorClientId() string {
	return c.Casdoor.ClientId
}

func (c *VclusterGatewayConfiguration) GetCasdoorClientSecret() string {
	return c.Casdoor.ClientSecret
}

func (c *VclusterGatewayConfiguration) GetCasdoorOrganizationName() string {
	return c.Casdoor.OrganizationName
}

func (c *VclusterGatewayConfiguration) GetCasdoorApplicationName() string {
	return c.Casdoor.ApplicationName
}

func (c *VclusterGatewayConfiguration) GetCasdoorCertificate() string {
	return c.Casdoor.Certificate
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorEndPoint() string {
	return c.Conductor.EndPoint
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorBatchSize() int {
	return c.Conductor.BatchSize
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorPollInterval() int {
	return c.Conductor.PollInterval
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorApiClient() *client.APIClient {
	return c.Conductor.ApiClient
}

func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorMetadataClient() *client.MetadataResourceApiService {
	return c.Conductor.MetadataClient
}
func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorWorkflowExecutor() *executor.WorkflowExecutor {
	return c.Conductor.WorkflowExecutor
}
func (c *VclusterGatewayConfiguration) GetVclusterGatewayConductorTaskRunner() *worker.TaskRunner {
	return c.Conductor.TaskRunner
}

func (c *VclusterGatewayConfiguration) GetPolicyEndPoint() string {
	return c.Iam.PolicyEndPoint
}

func (c *VclusterGatewayConfiguration) GetApsProvisioningUrl() string {
	return c.Aps.ProvisioningUrl
}
func (c *VclusterGatewayConfiguration) GetApsApiKey() string {
	return c.Aps.ApiKey
}
func (c *VclusterGatewayConfiguration) GetApsServingUrl() string {
	return c.Aps.ServingUrl
}
func (c *VclusterGatewayConfiguration) GetApsTrainingUrl() string {
	return c.Aps.TrainingUrl
}
func (c *VclusterGatewayConfiguration) GetApsTenantStatusURL() string {
	return c.Aps.GetTenantStatusURL
}

func (c *VclusterGatewayConfiguration) GetAlayaStudioHost() string {
	return c.AlayaStudio.Host
}
func (c *VclusterGatewayConfiguration) GetAlayaStudioApiKey() string {
	return c.AlayaStudio.ApiKey
}

func (c *VclusterGatewayConfiguration) GetCephClusterId() string {
	return c.Ceph.ClusterId
}

func (c *VclusterGatewayConfiguration) GetCacheInstance() *bigcache.BigCache {
	return c.AllCluster.Cache
}

func (c *VclusterGatewayConfiguration) ReadConfFromFile() error {
	file, err := os.ReadFile(c.ConfigFilePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, c); err != nil {
		return err
	}

	return nil
}

func (c *VclusterGatewayConfiguration) Watch(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Warnf("Create new config file watcher failed with error: %v", err)
		return
	}

	defer watcher.Close()

	if err := watcher.Add(c.ConfigFilePath); err != nil {
		log.Warnf("Watching config file %s with error: %v", c.ConfigFilePath, err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Infof("Config file %s is changed, refresh the configuration.", c.ConfigFilePath)

				if err := c.ReadConfFromFile(); err != nil {
					log.Warnf("Refresh config file with error: %v", err)
				}

				log.Infof("refresh conf: %#v", c)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Warnf("Watch config file %s with error: %v", c.ConfigFilePath, err)

		case <-ctx.Done():
			log.Info("Get shutdown signal.")
			return
		}
	}
}

func Default(cfg *VclusterGatewayConfiguration) error {
	cfg.ConfigFilePath = "/etc/config.yaml"
	cfg.EnableWatch = false

	cfg.Server.IP = "0.0.0.0"
	cfg.Server.Port = "8083"
	cfg.Server.TokenKey = "X-token"
	cfg.Server.RespCacheKey = "respBody"
	cfg.Server.Swagger.Host = "localhost:8088"
	cfg.Server.Swagger.BasePath = ""
	return nil
}
