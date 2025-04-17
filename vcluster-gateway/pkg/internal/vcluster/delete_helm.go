package forkvcluster

import (
	"context"
	"fmt"
	"time"

	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/utils"

	"github.com/loft-sh/log"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd/app/localkubernetes"
	"github.com/loft-sh/vcluster/pkg/helm"
	"github.com/loft-sh/vcluster/pkg/util/translate"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/vcluster/find"
)

// DeleteHelmOption 使用 Option模式创建 CreateHelm
type DeleteHelmOption func(*DeleteHelm)

type DeleteHelm struct {
	*GlobalFlags
	*DeleteOptions
	rawConfig        *clientcmdapi.Config
	restConfig       *rest.Config
	kubeClient       *kubernetes.Clientset
	kubeClientConfig *clientcmd.ClientConfig
	log              log.Logger
}

func NewDeleteHelm(options ...DeleteHelmOption) *DeleteHelm {
	dh := &DeleteHelm{}

	for _, option := range options {
		option(dh)
	}

	return dh
}

func WithDeleteGlobalFlags(globalFlags *GlobalFlags) DeleteHelmOption {
	return func(ch *DeleteHelm) {
		ch.GlobalFlags = globalFlags
	}
}

func WithDeleteLogger(logger log.Logger) DeleteHelmOption {
	return func(ch *DeleteHelm) {
		ch.log = logger
	}
}

func WithDeleteK8sConfig(config *v1.RootK8sConfig) DeleteHelmOption {
	return func(dh *DeleteHelm) {
		dh.kubeClientConfig = config.RootClientConfig
		dh.kubeClient = config.RootKubeClientSet
		dh.rawConfig = config.RootRawConfig
	}
}

func WithDeleteOptions(options *DeleteOptions) DeleteHelmOption {
	return func(dh *DeleteHelm) {
		dh.DeleteOptions = options
	}
}

// DeleteOptions holds the delete cmd options
type DeleteOptions struct {
	Project             string
	Wait                bool
	KeepPVC             bool
	DeleteNamespace     bool
	DeleteConfigMap     bool
	AutoDeleteNamespace bool
}

func NewDefaultDeleteOptions() *DeleteOptions {
	return &DeleteOptions{
		Project:             "",
		Wait:                false,
		KeepPVC:             false,
		DeleteNamespace:     true,
		DeleteConfigMap:     false,
		AutoDeleteNamespace: true,
	}
}

// DeleteHelm executes the functionality
func (cmd *DeleteHelm) DeleteHelm(ctx context.Context, vClusterId string) error {
	vClusterList, err := find.FindInContext(ctx,
		cmd.Context,
		vClusterId,
		utils.GetVClusterNamespaceName(vClusterId),
		time.Second*10,
		false,
		cmd.kubeClientConfig)
	if err != nil {
		cmd.log.Errorf("failed to find vcluster list, err: %v", err)
		return err
	}

	cmd.log.Infof("DeleteHelm, successfully find vcluster in context, vclusterList: %v", vClusterList)

	if len(vClusterList) == 0 {
		cmd.log.Warn("DeleteHelm, not found vcluster")
		return errors.New("not found vcluster")
	}

	err = cmd.prepare(&vClusterList[0])
	if err != nil {
		return err
	}

	cmd.log.Infof("DeleteHelm, scuccessfully delete vclsuter")

	return cmd.deleteResources(ctx, vClusterId)
}

// deleteResources performs the deletion of resources
func (cmd *DeleteHelm) deleteResources(ctx context.Context, vClusterId string) error {
	// 获取 helm 默认二进制文件路径
	helmBinaryPath := getHelmBinaryPath()

	cmd.log.Infof("deleteResources, cmd.Namespace: %s, cmd.AutoDeleteNamespace: %v", cmd.Namespace, cmd.AutoDeleteNamespace)

	// 检查命名空间
	if cmd.AutoDeleteNamespace {
		namespace, err := cmd.kubeClient.CoreV1().Namespaces().Get(ctx, cmd.Namespace, metav1.GetOptions{})
		if err != nil {
			cmd.log.Errorf("deleteResources, cmd.AutoDeleteNamespace is true, but invoke clientSet get k8s namespace failed, cmd.NameSpace: %s", cmd.Namespace)
		} else if namespace != nil && namespace.Annotations != nil && namespace.Annotations[createdByVClusterAnnotation] == "true" {
			cmd.log.Infof("deleteResources, cmd.AutoDeleteNamespace is true, will set cmd.DeleteNamespace true")
			cmd.DeleteNamespace = true
		}
	}

	cmd.log.Infof("deleteResources, will be deleted vcluster using helm command")

	// 删除 Helm chart
	err := helm.NewClient(cmd.rawConfig, cmd.log, helmBinaryPath).Delete(vClusterId, cmd.Namespace)
	if err != nil {
		return err
	}

	cmd.log.Infof("deleteResources, successfully deleted vcluster in namespace %s; now start clean vcluster resource, like: pvc, configmap...", cmd.Namespace)

	// 尝试删除 PVC
	if !cmd.KeepPVC && !cmd.DeleteNamespace {
		pvcName := fmt.Sprintf("data-%s-0", vClusterId)
		pvcNameForK8sAndEks := fmt.Sprintf("data-%s-etcd-0", vClusterId)

		client, err := kubernetes.NewForConfig(cmd.restConfig)
		if err != nil {
			return err
		}

		err = client.CoreV1().PersistentVolumeClaims(cmd.Namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return errors.Wrap(err, "delete pvc")
			}
		} else {
			cmd.log.Infof("deleteResources, successfully deleted virtual cluster pvc %s in namespace %s", pvcName, cmd.Namespace)
		}

		// 删除 K8s 和 EKS 发行版的 PVC
		err = client.CoreV1().PersistentVolumeClaims(cmd.Namespace).Delete(ctx, pvcNameForK8sAndEks, metav1.DeleteOptions{})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return errors.Wrap(err, "delete pvc")
			}
		} else {
			cmd.log.Infof("deleteResources, successfully deleted virtual cluster pvc %s in namespace %s", pvcName, cmd.Namespace)
		}
	}

	// 尝试删除 ConfigMap
	if cmd.DeleteConfigMap {
		client := cmd.kubeClient
		configMapName := fmt.Sprintf("configmap-%s", vClusterId)

		// 尝试删除 ConfigMap
		err = client.CoreV1().ConfigMaps(cmd.Namespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return errors.Wrap(err, "delete configmap")
			}
		} else {
			cmd.log.Infof("deleteResources, successfully deleted ConfigMap %s in namespace %s", configMapName, cmd.Namespace)
		}
	}

	// 尝试删除命名空间
	if cmd.DeleteNamespace {
		client := cmd.kubeClient

		// 删除命名空间
		err = client.CoreV1().Namespaces().Delete(ctx, cmd.Namespace, metav1.DeleteOptions{})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return errors.Wrap(err, "delete namespace")
			}
		} else {
			cmd.log.Infof("deleteResources, successfully deleted virtual cluster namespace %s", cmd.Namespace)
		}

		// 删除多命名空间模式下的命名空间
		namespaces, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
			LabelSelector: translate.MarkerLabel + "=" + translate.SafeConcatName(cmd.Namespace, "x", vClusterId),
		})
		if err != nil && !kerrors.IsForbidden(err) {
			return errors.Wrap(err, "list namespaces")
		}

		// 删除所有命名空间
		if namespaces != nil && len(namespaces.Items) > 0 {
			for _, namespace := range namespaces.Items {
				err = client.CoreV1().Namespaces().Delete(ctx, namespace.Name, metav1.DeleteOptions{})
				if err != nil {
					if !kerrors.IsNotFound(err) {
						return errors.Wrap(err, "delete namespace")
					}
				} else {
					cmd.log.Infof("deleteResources, successfully deleted virtual cluster namespace %s", namespace.Name)
				}
			}
		}

		// 等待命名空间删除
		time.Sleep(time.Second * 3)
		nameSpace, err := client.CoreV1().Namespaces().Get(ctx, cmd.Namespace, metav1.GetOptions{})
		if err != nil {
			// 设置删除传播策略为 Foreground，确保级联删除所有资源
			deletePropagation := metav1.DeletePropagationForeground
			err = client.CoreV1().Namespaces().Delete(context.TODO(), nameSpace.Name, metav1.DeleteOptions{
				PropagationPolicy: &deletePropagation,
			})

			// 清空命名空间对象的 finalizers，防止删除卡住
			nameSpace.ObjectMeta.Finalizers = nil

			_, _ = client.CoreV1().Namespaces().Update(ctx, nameSpace, metav1.UpdateOptions{})
			if err != nil {
				cmd.log.Warnf("deleteResources, error deleting Namespace: %v\n", err)
			}
		}

		// 等待 vcluster 删除
		if cmd.Wait {
			cmd.log.Infof("deleteResources, waiting for vcluster to be deleted...")

			for {
				_, err = client.CoreV1().Namespaces().Get(ctx, cmd.Namespace, metav1.GetOptions{})
				if err != nil {
					break
				}

				nameSpace, err := client.CoreV1().Namespaces().Get(ctx, cmd.Namespace, metav1.GetOptions{})
				if err != nil {
					nameSpace.Spec.Finalizers = []corev1.FinalizerName{}
					_, _ = client.CoreV1().Namespaces().Update(ctx, nameSpace, metav1.UpdateOptions{})
				}
			}
			cmd.log.Infof("deleteResources, vcluser resource was deleted")
		}
	}

	return nil
}

func (cmd *DeleteHelm) prepare(vCluster *find.VCluster) error {
	// 加载 raw 配置
	rawConfig, err := vCluster.ClientFactory.RawConfig()
	if err != nil {
		return fmt.Errorf("there is an error loading your current kube config (%w), please make sure you have access to a kubernetes cluster and the command `kubectl get namespaces` is working", err)
	}

	err = cmd.deleteContext(&rawConfig, find.VClusterContextName(vCluster.Name, vCluster.Namespace, vCluster.Context), vCluster.Context)
	if err != nil {
		return errors.Wrap(err, "delete kube context")
	}

	rawConfig.CurrentContext = vCluster.Context
	restConfig, err := vCluster.ClientFactory.ClientConfig()
	if err != nil {
		return err
	}

	err = localkubernetes.CleanupLocal(vCluster.Name, vCluster.Namespace, &rawConfig, cmd.log)
	if err != nil {
		cmd.log.Warnf("error cleaning up: %v", err)
	}

	// 构建代理名称
	proxyName := find.VClusterConnectBackgroundProxyName(vCluster.Name, vCluster.Namespace, rawConfig.CurrentContext)
	_ = localkubernetes.CleanupBackgroundProxy(proxyName, cmd.log)

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	cmd.Namespace = vCluster.Namespace
	cmd.rawConfig = &rawConfig
	cmd.restConfig = restConfig
	cmd.kubeClient = kubeClient
	return nil
}

func (cmd *DeleteHelm) deleteContext(kubeConfig *clientcmdapi.Config, kubeContext string, otherContext string) error {
	// 获取上下文
	contextRaw, ok := kubeConfig.Contexts[kubeContext]
	if !ok {
		return nil
	}

	// 删除上下文
	delete(kubeConfig.Contexts, kubeContext)

	removeAuthInfo := true
	removeCluster := true

	// 检查 AuthInfo 或 Cluster 是否被其他上下文使用
	for name, ctx := range kubeConfig.Contexts {
		if name != kubeContext && ctx.AuthInfo == contextRaw.AuthInfo {
			removeAuthInfo = false
		}

		if name != kubeContext && ctx.Cluster == contextRaw.Cluster {
			removeCluster = false
		}
	}

	// 如果 AuthInfo 没有被其他上下文使用，则删除它
	if removeAuthInfo {
		delete(kubeConfig.AuthInfos, contextRaw.AuthInfo)
	}

	// 如果 Cluster 没有被其他上下文使用，则删除它
	if removeCluster {
		delete(kubeConfig.Clusters, contextRaw.Cluster)
	}

	if kubeConfig.CurrentContext == kubeContext {
		kubeConfig.CurrentContext = ""

		if otherContext != "" {
			kubeConfig.CurrentContext = otherContext
		} else if len(kubeConfig.Contexts) > 0 {
			for contextName, contextObj := range kubeConfig.Contexts {
				if contextObj != nil {
					kubeConfig.CurrentContext = contextName
					break
				}
			}
		}
	}

	return clientcmd.ModifyConfig(clientcmd.NewDefaultClientConfigLoadingRules(), *kubeConfig, false)
}
