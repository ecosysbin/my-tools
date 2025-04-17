package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/model"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/utils"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/vcluster/find"
)

type ClusterName = string

const (
	DefaultCluster ClusterName = "inCluster"
	// ResyncPeriod is the resync period for the shared informer factory
	ResyncPeriod = time.Second * 30
)

type KubeClientManager struct {
	ClientSet       *kubernetes.Clientset
	InformerFactory informers.SharedInformerFactory
}

var _ VClusterKubernetesDataSource = &KubernetesDataSource{}

type KubernetesDataSource struct {
	clientSet       *kubernetes.Clientset
	informerFactory informers.SharedInformerFactory

	clientConfig *clientcmd.ClientConfig
	rawConfig    *api.Config

	stopCh chan struct{}
	logger *log.Logger

	DB VClusterDBDataSource
}

func NewKubernetesDataSource(db VClusterDBDataSource) (VClusterKubernetesDataSource, error) {
	logger := log.WithField("layer", "kubernetes_datasource")

	// 使用 ServiceAccount 的 token 生成 kubeconfig 文件
	if err := generateKubeconfig(); err != nil {
		return nil, errors.Errorf("error generating kubeconfig: %w", err)
	}
	defer func() {
		home := os.Getenv("HOME")
		if home != "" {
			path := filepath.Join(home, ".kube", "config")
			if err := os.RemoveAll(path); err != nil {
				logger.Errorf("Failed to remove kubeconfig file: %v", err)
			}
		}
	}()

	// 初始化底层集群的 clientSet
	clientSet, err := innerCluster()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize clientSet")
	}

	// 初始化 clientConfig 和 rawConfig
	clientConfig, rawConfig, err := initClientConfigAndRawConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize clientConfig and rawConfig")
	}

	var kds KubernetesDataSource

	kds = KubernetesDataSource{
		clientSet:       clientSet,
		informerFactory: kds.initSharedInformerFactory(clientSet),
		clientConfig:    clientConfig,
		rawConfig:       rawConfig,
		stopCh:          make(chan struct{}),
		DB:              db,
		logger:          logger,
	}

	kds.startInformerFactory()

	logger.Infof("NewKubernetesDataSource initialized successfully")

	return &kds, nil
}

func (kds *KubernetesDataSource) Shutdown() {
	// 关闭 底层 K8s 集群 Informer 使用的 stopCh
	close(kds.stopCh)

	log.Infof("KubernetesDataSource shut down successfully")
}

func (kds *KubernetesDataSource) GetClientSet() *kubernetes.Clientset {
	return kds.clientSet
}

func (kds *KubernetesDataSource) GetSharedInformerFactory() informers.SharedInformerFactory {
	return kds.informerFactory
}

func (kds *KubernetesDataSource) GetConfigs() (*clientcmd.ClientConfig, *api.Config) {
	return kds.clientConfig, kds.rawConfig
}

func (kds *KubernetesDataSource) GenerateVClusterClientSet(vClusterId string) (*kubernetes.Clientset, error) {
	generateOriginClientSet := func(secret *corev1.Secret, service *corev1.Service) (*kubernetes.Clientset, error) {
		vcKubeconfigBytes := secret.Data["config"]
		vcKubeconfig, err := clientcmd.Load(vcKubeconfigBytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load kubeconfig")
		}

		serviceIp := "https://" + service.Spec.ClusterIP
		vcKubeconfig.Clusters["my-vcluster"].Server = serviceIp

		vcClientConfig, err := clientcmd.NewDefaultClientConfig(*vcKubeconfig, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create client config")
		}

		vcClientConfig.Timeout = time.Minute * 5

		vcClientSet, err := kubernetes.NewForConfig(vcClientConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create clientset")
		}

		return vcClientSet, nil
	}

	logger := kds.logger.WithField("func", "generateVClusterClientSet")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rootClientSet := kds.clientSet

	secret, err := rootClientSet.CoreV1().Secrets(utils.GetVClusterNamespaceName(vClusterId)).Get(ctx, utils.GetVClusterSecretName(vClusterId), metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			secret, err = rootClientSet.CoreV1().Secrets(utils.GetInfraVClusterNamespaceName(vClusterId)).Get(ctx, utils.GetVClusterSecretName(vClusterId), metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get secret")
			}
		}
	}

	logger.Infof("success create secret, secret name: %s, vcluster id: %s", secret.Name, vClusterId)

	service, err := rootClientSet.CoreV1().Services(utils.GetVClusterNamespaceName(vClusterId)).Get(ctx, utils.GetVClusterServiceName(vClusterId), metav1.GetOptions{})
	if err != nil {
		logger.Warnf("failed to get service, err: %v, vcluster id :%s", err, vClusterId)
		if kerrors.IsNotFound(err) {
			service, err = rootClientSet.CoreV1().Services(utils.GetInfraVClusterNamespaceName(vClusterId)).Get(ctx, utils.GetVClusterServiceName(vClusterId), metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get service")
			}
		}
	}

	logger.Infof("success create service, service name: %s, vcluster id: %s", service.Name, vClusterId)

	vcClientSet, err := generateOriginClientSet(secret, service)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate vc origin clientset")
	}

	logger.Infof("successfully get vcluster clientset: %+v", vcClientSet)

	return vcClientSet, nil
}

//func (kds *KubernetesDataSource) WaitForAllVClusterPodsReady(oldPod, newPod *corev1.Pod) {
//	if oldPod.ObjectMeta.ResourceVersion == newPod.ObjectMeta.ResourceVersion {
//		return
//	}
//
//	namespace := newPod.Namespace
//	if !strings.HasPrefix(namespace, vclusterPrefix) {
//		// 没有匹配到 vclusterPrefix 直接返回
//		return
//	}
//
//	id := namespace[len(vclusterPrefix):]
//	if len(id) != 12 {
//		// 如果长度不等于12，直接返回
//		return
//	}
//
//	// 当前 Pod 是 vcluster-id 命名空间下的
//
//	if newPod.Annotations[serviceAccountAnnotation] == corednsServiceAccount {
//		// 当前 Pod 是 coreDNS
//		if newPod.Status.Phase == corev1.PodRunning && newPod.DeletionTimestamp == nil {
//			// 并且 Pod 的状态是 Running 且没有被标记删除
//		}
//	}
//}

func isVClusterDeployment(deployment *appsv1.Deployment) bool {
	id, exists := deployment.Labels["release"]
	if !exists || len(id) != 12 {
		// 如果 vcluster id 长度不为 12 位，直接返回 false
		return false
	}

	availableSet := map[string]bool{
		utils.GetVClusterNamespaceName(id):      true,
		utils.GetInfraVClusterNamespaceName(id): true,
	}

	if _, ok := availableSet[deployment.Namespace]; !ok {
		// 当前 Deployment 所在的 namespace 不是 vcluster-id 格式或者 infra-vcluster-id 格式，直接返回 false
		return false
	}

	if deployment.Labels["app"] != appName {
		// 如果 Deployment 的 app 标签不为 "vcluster"，直接返回 false
		return false
	}

	// 满足所有条件，返回 true
	return true
}

// DeleteDeploymentEventHandler Kubernetes Deployment Event Handler, will be update mysql record
func (kds *KubernetesDataSource) DeleteDeploymentEventHandler(deployment *appsv1.Deployment) {
	isVCluster := isVClusterDeployment(deployment)
	if !isVCluster {
		return
	}

	id, _ := deployment.Labels["release"]

	logger := kds.logger.WithField("func", "DeleteDeploymentEventHandler").
		WithField("vclusterId", id).
		WithField("deploymentName", deployment.Name).
		WithField("deploymentNamespace", deployment.Namespace)

	logger.Infof("kubernetes deployment informer received Delete Event")

	record, err := kds.DB.GetVClusterById(id)
	if err != nil {
		logger.Warnf("failed to GetVClusterById, err: %v, vcluster id :%s", err, id)
		return
	}

	deleteTime := time.Now()
	err = kds.DB.UpdateVCluster(record, &model.VCluster{
		DeleteTime:   &deleteTime,
		Status:       string(vclusterStatusDeleted),
		ServerStatus: string(vclusterStatusDeleted),
		IsDeleted:    deleteFlag,
	})
	if err != nil {
		logger.Errorf("failed to UpdateVCluster, err: %v, vcluster id :%s", err, id)
		return
	}

	vStorage, err := kds.DB.GetVStorageByVClusterId(id)
	if err != nil {
		logger.Errorf("failed to GetVStorageByVClusterId, err: %v, vcluster id :%s", err, id)
		return
	}

	err = kds.DB.UpdateVStorage(vStorage, &model.VStorage{
		IsDeleted: deleteFlag,
	})
	if err != nil {
		logger.Errorf("failed to UpdateVStorage, err: %v, vcluster id :%s", err, id)
		return
	}
}

func (kds *KubernetesDataSource) UpdateVClusterRecord(oldDeployment, newDeployment *appsv1.Deployment) {
	isVClster := isVClusterDeployment(newDeployment)
	if !isVClster {
		return
	}

	if oldDeployment.ResourceVersion == newDeployment.ResourceVersion {
		return
	}

	id, _ := newDeployment.Labels["release"]

	logger := kds.logger.WithField("func", "UpdateVClusterRecord").
		WithField("vclusterId", id).
		WithField("deploymentName", newDeployment.Name).
		WithField("deploymentNamespace", newDeployment.Namespace)

	logger.Infof("kubernetes deployment informer received Update Event")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)

	clusters, err := find.FindInContext(ctx, DefaultCluster, newDeployment.Name, newDeployment.Namespace, time.Second*10, false, kds.clientConfig)
	if err != nil || len(clusters) == 0 {
		logger.Errorf("failed to find vcluster, err: %v, vcluster id: %s, vcluster len: %d", err, id, len(clusters))
		return
	}

	vclusterModel, err := kds.DB.GetVClusterById(id)
	if err != nil {
		logger.Errorf("failed to GetVClusterById, err: %v, vcluster id: %s", err, id)
		return
	}

	var serverStatus vclusterStatus
	clusterStatus := clusters[0].Status

	switch clusterStatus {
	case find.StatusRunning:
		serverStatus = vclusterStatusRunning
		err = kds.updateVClusterWithRunning(vclusterModel, clusterStatus, serverStatus)

	case find.StatusPaused:
		serverStatus = vclusterStatusPaused
		err = kds.updateVClusterWithNotRunning(id, vclusterModel, clusterStatus, serverStatus)

	case find.StatusPending, find.StatusUnknown:
		go func() {
			ticker := time.NewTicker(10 * time.Second) // 每 10 秒检查一次状态
			defer ticker.Stop()

			timeout := time.After(5 * time.Minute) // 总超时时间为 5 分钟
			var clusterStatusTmp find.Status

			for {
				select {
				case <-timeout:
					// 超时逻辑
					logger.Warnf("Timeout reached while finding vcluster, proceeding to update database.")
					err := kds.updateVClusterWithNotRunning(id, vclusterModel, clusterStatusTmp, vclusterStatusFailed)
					if err != nil {
						logger.Errorf("failed to update vcluster, err: %v, vcluster id: %s", err, id)
					}
					return

				case <-ticker.C:
					// 定时检查逻辑
					clusters, err := find.FindInContext(ctx, DefaultCluster, newDeployment.Name, newDeployment.Namespace, time.Second*10, false, kds.clientConfig)
					if err != nil || len(clusters) == 0 {
						logger.Errorf("failed to find vcluster, err: %v, vcluster id: %s, vcluster len: %d", err, id, len(clusters))
						continue // 如果失败，继续下一轮检查
					}

					clusterStatusTmp = clusters[0].Status
					if clusterStatusTmp == find.StatusRunning {
						logger.Infof("Vcluster is running, id: %s", id)
						return // 状态为 running，提前结束
					}
				}
			}
		}()

		serverStatus = ""
		err = kds.updateVClusterWithNotRunning(id, vclusterModel, clusterStatus, serverStatus)

		//default:
		//	serverStatus = vclusterStatusFailed
		//	err = kds.updateVClusterWithNotRunning(id, vclusterModel, clusterStatus, serverStatus)
	}

	if err != nil {
		logger.Errorf("failed to update vcluster, err: %v, vcluster id: %s", err, id)
	}
}

func (kds *KubernetesDataSource) updateVClusterWithRunning(vcluster *model.VCluster, status find.Status, serverStatus vclusterStatus) error {
	startTime := time.Now().Add(5 * time.Second)

	needUpdates := &model.VCluster{Status: string(status), StartTime: &startTime}
	if serverStatus != "" {
		needUpdates.ServerStatus = string(serverStatus)
	}

	err := kds.DB.UpdateVCluster(vcluster, needUpdates)
	if err != nil {
		return errors.Wrapf(err, "failed to UpdateVCluster, err: %v, vcluster id :%s", err, vcluster.VClusterId)
	}

	return nil
}

func (kds *KubernetesDataSource) updateVClusterWithNotRunning(id string, vcluster *model.VCluster, status find.Status, serverStatus vclusterStatus) error {
	needUpdates := &model.VCluster{Status: string(status)}
	if serverStatus != "" {
		needUpdates.ServerStatus = string(serverStatus)
	}

	if serverStatus == vclusterStatusFailed {
		kds.logger.Warnf("updateVClusterWithNotRunning, serverStatus: %s, id: %s, now need get events via clientset", serverStatus, id)
		events, err := kds.GetVClusterNamespaceEvents(id)
		if err == nil && events != nil {
			needUpdates.Reason = *events
		} else if err != nil {
			kds.logger.Errorf("Failed to get events for vcluster %s: %v", id, err)
		} else {
			kds.logger.Warnf("No events found for vcluster %s", id)
		}
	}

	err := kds.DB.UpdateVCluster(vcluster, needUpdates)
	if err != nil {
		return errors.Wrapf(err, "failed to UpdateVCluster, err: %v, vcluster id :%s", err, id)
	}

	return nil
}

func (kds *KubernetesDataSource) GetVClusterNamespaceEvents(vclusterId string) (*string, error) {
	eventList, err := kds.clientSet.CoreV1().Events(utils.GetVClusterNamespaceName(vclusterId)).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	abnormalEvents := make([]string, 0, len(eventList.Items))
	for _, event := range eventList.Items {
		if event.Type == corev1.EventTypeWarning && event.Message != "" {
			abnormalEvents = append(abnormalEvents, event.Message)
		}
	}

	jsonBytes, err := json.Marshal(abnormalEvents)
	if err != nil {
		return nil, err
	}

	jsonStr := string(jsonBytes)
	return &jsonStr, nil
}

// generateKubeconfig executes the shell script to generate the kubeconfig file.
// This function uses the in-cluster method to create the kubeconfig file
// by reading the Kubernetes service account tokens and CA certificates
// mounted inside the pod at /var/run/secrets/kubernetes.io/serviceaccount.
func generateKubeconfig() error {
	cmd := exec.Command("sh", "-c", kubeconfigScript)
	var stderr strings.Builder
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to execute generateKubeconfig script, err: %v, stderr: %s", err, stderr.String())
	}

	return nil
}

// innerCluster creates an in-cluster Kubernetes clientSet.
func innerCluster() (*kubernetes.Clientset, error) {
	// Get the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get in-cluster config")
	}

	config.QPS = DefaultClientQPS
	config.Burst = DefaultClientBurst

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create in-cluster clientSet")
	}

	return clientSet, nil
}

func initClientConfigAndRawConfig() (*clientcmd.ClientConfig, *api.Config, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	bytes, err := os.ReadFile(kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
	}

	clientConfig, err := clientcmd.NewClientConfigFromBytes(bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client config from bytes: %w", err)
	}

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get raw config: %w", err)
	}

	return &clientConfig, &rawConfig, nil
}

// initSharedInformerFactory initializes root k8s SharedInformerFactory.
func (kds *KubernetesDataSource) initSharedInformerFactory(clientSet *kubernetes.Clientset) informers.SharedInformerFactory {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientSet, ResyncPeriod, informers.WithNamespace(corev1.NamespaceAll))

	// podInformer 主要作用是通过 update 事件来判断一个 vCluster 是否启动完成（所有 pod 都 running）
	podInformer := informerFactory.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			// oldPod := oldObj.(*corev1.Pod)
			// newPod := newObj.(*corev1.Pod)

			// kds.WaitForAllVClusterPodsReady(oldPod, newPod)
		},
	})

	// deploymentInformer 主要作用：修改 vCluster 在数据库中的状态，例如 Deleted、Paused
	deploymentInformer := informerFactory.Apps().V1().Deployments().Informer()

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldDeployment, newDeployment := oldObj.(*appsv1.Deployment), newObj.(*appsv1.Deployment)

			kds.UpdateVClusterRecord(oldDeployment, newDeployment)
		},
		DeleteFunc: func(obj interface{}) {
			deployment := obj.(*appsv1.Deployment)

			kds.DeleteDeploymentEventHandler(deployment)
		},
	})

	informerFactory.Start(kds.stopCh)

	return informerFactory
}

func (kds *KubernetesDataSource) startInformerFactory() {
	// 启动 informer
	kds.informerFactory.Start(kds.stopCh)

	// 等待缓存同步
	for informer, synced := range kds.informerFactory.WaitForCacheSync(kds.stopCh) {
		if !synced {
			log.Errorf("failed to sync cache for informer: %s", informer)
			return
		}
	}

	log.Infof("successfully sync informer cacche")
}
