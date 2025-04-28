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

package options

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	log "github.com/sirupsen/logrus"

	configv1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
)

type VclusterGatewayOptions struct {
	*configv1.VclusterGatewayConfiguration
}

func (o *VclusterGatewayOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.ConfigFilePath, "config", o.ConfigFilePath, "Config file path.")
	fs.BoolVar(&o.EnableWatch, "enable-watch", o.EnableWatch, "If true, will auto watch config file and refresh configuration.")

	fs.StringVar(&o.Server.IP, "ip", o.Server.IP, "Http server listen address.")
	fs.StringVar(&o.Server.Port, "port", o.Server.Port, "Http server listen port.")
	fs.StringVar(&o.Server.TokenKey, "token-key", o.Server.TokenKey, "Http server token key.")
	fs.StringVar(&o.Server.RespCacheKey, "resp-cache-key", o.Server.RespCacheKey, "Http server resp cache key.")

	fs.StringVar(&o.Casdoor.Endpoint, "casdoor-endpoint", o.Casdoor.Endpoint, "Casdoor server endpoint.")
	fs.StringVar(&o.Casdoor.ClientId, "casdoor-client-id", o.Casdoor.ClientId, "Casdoor server client id.")
	fs.StringVar(&o.Casdoor.ClientSecret, "casdoor-client-secret", o.Casdoor.ClientSecret, "Casdoor server client secret.")
	fs.StringVar(&o.Casdoor.OrganizationName, "casdoor-organization-name", o.Casdoor.OrganizationName, "Casdoor server organization name.")
	fs.StringVar(&o.Casdoor.ApplicationName, "casdoor-application-name", o.Casdoor.ApplicationName, "Casdoor server application name.")
	fs.StringVar(&o.Casdoor.Certificate, "casdoor-certificate", o.Casdoor.Certificate, "Casdoor server certificate.")
	fs.StringVar(&o.Server.Swagger.Host, "swagger-host", o.Server.Swagger.Host, "Swagger host.")
	fs.StringVar(&o.Server.Swagger.BasePath, "swagger-basepath", o.Server.Swagger.BasePath, "Swagger base path.")
}

func (o *VclusterGatewayOptions) ApplyTo(cfg *configv1.VclusterGatewayConfiguration) error {
	if o.ConfigFilePath != "" {
		cfg.ConfigFilePath = o.ConfigFilePath

		if err := cfg.ReadConfFromFile(); err != nil {
			return err
		}
	}

	if o.EnableWatch {
		cfg.EnableWatch = o.EnableWatch
		// If watch config file changes is enabled,
		// the command line arguments will be overwritten.
		return nil
	}

	if o.Server.TokenKey != "" {
		cfg.Server.TokenKey = o.Server.TokenKey
	}
	if o.Server.RespCacheKey != "" {
		cfg.Server.RespCacheKey = o.Server.RespCacheKey
	}
	if o.Casdoor.Endpoint != "" {
		cfg.Casdoor.Endpoint = o.Casdoor.Endpoint
	}
	if o.Casdoor.ClientId != "" {
		cfg.Casdoor.ClientId = o.Casdoor.ClientId
	}
	if o.Casdoor.ClientSecret != "" {
		cfg.Casdoor.ClientSecret = o.Casdoor.ClientSecret
	}
	if o.Casdoor.OrganizationName != "" {
		cfg.Casdoor.OrganizationName = o.Casdoor.OrganizationName
	}
	if o.Casdoor.ApplicationName != "" {
		cfg.Casdoor.ApplicationName = o.Casdoor.ApplicationName
	}
	if o.Casdoor.Certificate != "" {
		cfg.Casdoor.Certificate = o.Casdoor.Certificate
	}
	if o.VclusterGateway.Dsn != "" {
		cfg.VclusterGateway.Dsn = o.VclusterGateway.Dsn
	}
	if o.VclusterGateway.DefaultCluster != "" {
		cfg.VclusterGateway.DefaultCluster = o.VclusterGateway.DefaultCluster
	}
	if o.VclusterGateway.StorageClass != "" {
		cfg.VclusterGateway.StorageClass = o.VclusterGateway.StorageClass
	}
	if o.VclusterGateway.DefaultStorageClass != "" {
		cfg.VclusterGateway.DefaultStorageClass = o.VclusterGateway.DefaultStorageClass
	}
	if o.VclusterGateway.ChartRepo != "" {
		cfg.VclusterGateway.ChartRepo = o.VclusterGateway.ChartRepo
	}
	if o.VclusterGateway.DefaultImageRegistry != "" {
		cfg.VclusterGateway.DefaultImageRegistry = o.VclusterGateway.DefaultImageRegistry
	}
	if o.Conductor.EndPoint != "" {
		cfg.Conductor.EndPoint = o.Conductor.EndPoint
	}

	if o.Conductor.BatchSize != 0 {
		cfg.Conductor.EndPoint = o.Conductor.EndPoint
	} else {
		cfg.Conductor.BatchSize = 5
	}

	if o.Conductor.PollInterval != 0 {
		cfg.Conductor.PollInterval = o.Conductor.PollInterval
	} else {
		cfg.Conductor.PollInterval = 5
	}

	if o.Aps.ProvisioningUrl != "" {
		cfg.Aps.ProvisioningUrl = o.Aps.ProvisioningUrl
	}
	if o.Aps.ApiKey != "" {
		cfg.Aps.ApiKey = o.Aps.ApiKey
	}
	if o.Aps.GetTenantStatusURL != "" {
		cfg.Aps.GetTenantStatusURL = o.Aps.GetTenantStatusURL
	}
	if o.Aps.ServingUrl != "" {
		cfg.Aps.ServingUrl = o.Aps.ServingUrl
	}
	if o.Aps.TrainingUrl != "" {
		cfg.Aps.TrainingUrl = o.Aps.TrainingUrl
	}
	if o.Ceph.ClusterId != "" {
		cfg.Ceph.ClusterId = o.Ceph.ClusterId
	}

	if o.StorageManager.Host != "" {
		cfg.StorageManager.Host = o.StorageManager.Host
	}

	if o.Cache.Enable {
		cfg.Cache.LifeTime = o.Cache.LifeTime
		cfg.Cache.CleanWindow = o.Cache.CleanWindow
		cfg.Cache.Enable = o.Cache.Enable
		cfg.Cache.HardMaxCacheSize = o.Cache.HardMaxCacheSize
		cfg.Cache.Shards = o.Cache.Shards
	} else {
		o.Cache.Enable = false
	}
	cfg.AllCluster = o.setRootClusterStructV1()
	return nil
}

func (o *VclusterGatewayOptions) setRootClusterStructV1() configv1.AllClusterStruct {
	if o.ConfigFilePath == "" {
		return configv1.AllClusterStruct{}
	}

	// 注释使用 kubeconfig 初始化底层 K8s 的代码
	//config, err := os.ReadFile(o.GetVclusterGatewayKubeConfig())
	//
	//if err != nil {
	//	log.Error("read kube config failed,", err)
	//	os.Exit(1)
	//}
	//clientConfig, err := clientcmd.NewClientConfigFromBytes(config)
	//if err != nil {
	//	log.Error("create root cluster struct failed,", err)
	//	os.Exit(1)
	//}

	o.AllCluster = configv1.AllClusterStruct{
		Cluster:        map[string]configv1.RootClusterStruct{},
		CurrentContext: o.VclusterGateway.DefaultCluster,
		DB:             nil,
	}

	// o.createRootClusterStructV1(clientConfig)
	o.connectDb(o.GetVclusterGatewayDsn())

	if o.Cache.Enable {
		o.createBigCache()
	}

	o.Conductor.ApiClient = client.NewAPIClient(
		nil,
		settings.NewHttpSettings(
			o.Conductor.EndPoint,
		))
	o.Conductor.WorkflowExecutor = executor.NewWorkflowExecutor(o.Conductor.ApiClient)
	o.Conductor.MetadataClient = &client.MetadataResourceApiService{APIClient: o.Conductor.ApiClient}
	o.Conductor.TaskRunner = worker.NewTaskRunnerWithApiClient(o.Conductor.ApiClient)
	return o.AllCluster
}

func (o *VclusterGatewayOptions) createBigCache() {
	var err error

	o.AllCluster.Cache, err = bigcache.New(context.Background(), bigcache.Config{
		Shards:               o.Cache.Shards,
		LifeWindow:           time.Duration(o.Cache.LifeTime) * time.Second,
		CleanWindow:          time.Duration(o.Cache.CleanWindow) * time.Second,
		MaxEntriesInWindow:   0,
		MaxEntrySize:         0,
		StatsEnabled:         true,
		Verbose:              false,
		Hasher:               nil,
		HardMaxCacheSize:     o.Cache.HardMaxCacheSize,
		OnRemove:             nil,
		OnRemoveWithMetadata: nil,
		OnRemoveWithReason:   nil,
		Logger:               nil,
	})
	if err != nil {
		panic(err)
	}
}

func (o *VclusterGatewayOptions) connectDb(dsn string) {
	// d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	DisableForeignKeyConstraintWhenMigrating: true,
	// 	Logger:                                   logger.Default.LogMode(logger.Warn),
	// 	AllowGlobalUpdate:                        true,
	// },
	// )
	// if err != nil {
	// 	log.Error("open mysql failed,", err)
	// 	os.Exit(1)
	// }

	// o.AllCluster.DB = d
	// sqlDB, e := o.AllCluster.DB.DB()
	// if e != nil {
	// 	log.Error("sql.db  failed,", e)
	// 	os.Exit(1)
	// }

	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// sqlDB.SetConnMaxLifetime(time.Minute * 30)
}

func (o *VclusterGatewayOptions) createVClusterToken() {
}

func (o *VclusterGatewayOptions) createRootClusterStructV1(clientConfig clientcmd.ClientConfig) {
	rootRawConfig, err := clientConfig.RawConfig()
	if err != nil {
		log.Errorf("there is an error loading your current kube config (%w), please make sure you have access to a kubernetes cluster and the command `kubectl get namespaces` is working", err)
	}

	rootRawConfig.CurrentContext = rootRawConfig.Contexts[rootRawConfig.CurrentContext].Cluster
	o.AllCluster.CurrentContext = rootRawConfig.CurrentContext

	for name, cluster := range rootRawConfig.Clusters {
		rootConfig := api.Config{
			Kind:           "Config",
			APIVersion:     "v1",
			Preferences:    api.Preferences{},
			Clusters:       make(map[string]*api.Cluster),
			AuthInfos:      make(map[string]*api.AuthInfo),
			Contexts:       make(map[string]*api.Context),
			CurrentContext: "",
			Extensions:     nil,
		}

		rootConfig.Clusters[name] = cluster
		rootConfig.AuthInfos[name] = rootRawConfig.AuthInfos[name]
		rootConfig.Contexts[name] = rootRawConfig.Contexts[name]

		rootConfig.CurrentContext = name

		ConfigBytes, e := clientcmd.Write(rootConfig)
		if e != nil {
			log.Errorf("write kube config %v", err)
		}

		ClientConfig, errClientConfig := clientcmd.NewClientConfigFromBytes(ConfigBytes)

		if errClientConfig != nil {
			log.Errorf("create clientConfig %v", err)
		}

		RestConfig, rerr := ClientConfig.ClientConfig()
		if rerr != nil {
			log.Errorf("load RestConfig %v", rerr)
		}
		RestConfig.Burst = 2000
		RestConfig.QPS = 2000

		KubeClientSet, errClientSet := kubernetes.NewForConfig(RestConfig)

		if errClientSet != nil {
			log.Errorf("create KubeClientSet %v", rerr)
		}

		factory := informers.NewSharedInformerFactoryWithOptions(KubeClientSet, time.Second*30, informers.WithNamespace(corev1.NamespaceAll))
		o.AllCluster.Cluster[name] = configv1.RootClusterStruct{
			RootRawConfig:     &rootConfig,
			RootRestConfig:    RestConfig,
			RootKubeClientSet: KubeClientSet,
			RootClientConfig:  &ClientConfig,
			Factory:           &factory,
			VClusterMap:       map[string]configv1.VClusterStruct{},
		}

		cs := o.AllCluster.Cluster[name].RootKubeClientSet
		_, gErr := cs.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if gErr != nil {
			log.Errorf("root cluster %s get pods succes: %v", name, gErr)
		} else {
			log.Infof("root cluster %s get pods succes", name)
		}

	}
}

func (o *VclusterGatewayOptions) Validate() []error {
	if o == nil {
		return nil
	}

	var errs []error

	if o.ConfigFilePath == "" {
		if o.Server.TokenKey == "" {
			errs = append(errs, fmt.Errorf("server TokenKey is must"))
		}

		if o.Server.RespCacheKey == "" {
			errs = append(errs, fmt.Errorf("server RespCacheKey is must"))
		}
		if o.Casdoor.Endpoint == "" {
			errs = append(errs, fmt.Errorf("casdoor Endpoint is must"))
		}

		if o.Casdoor.ClientId == "" {
			errs = append(errs, fmt.Errorf("casdoor ClientId is must"))
		}

		if o.Casdoor.ClientSecret == "" {
			errs = append(errs, fmt.Errorf("casdoor ClientSecret is must"))
		}

		if o.Casdoor.OrganizationName == "" {
			errs = append(errs, fmt.Errorf("casdoor OrganizationName is must"))
		}

		if o.Casdoor.ApplicationName == "" {
			errs = append(errs, fmt.Errorf("casdoor ApplicationName is must"))
		}

		if o.Casdoor.Certificate == "" {
			errs = append(errs, fmt.Errorf("casdoor Certificate is must"))
		}
	}

	return errs
}
