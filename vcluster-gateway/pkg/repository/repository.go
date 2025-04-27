package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	v1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
	"vcluster-gateway/pkg/datasource"
	"vcluster-gateway/pkg/internal/model"
	"vcluster-gateway/pkg/internal/utils"
	"vcluster-gateway/pkg/usecase/consts"

	log "github.com/sirupsen/logrus"
)

var _ VClusterRepository = &VClusterRepositoryImpl{}

type VClusterRepositoryImpl struct {
	dbDataSource  datasource.VClusterDBDataSource
	k8sDataSource datasource.VClusterKubernetesDataSource
}

func NewVClusterRepository(dbDataSource datasource.VClusterDBDataSource,
	k8sDataSource datasource.VClusterKubernetesDataSource,
) VClusterRepository {
	vcRepoImpl := &VClusterRepositoryImpl{
		dbDataSource:  dbDataSource,
		k8sDataSource: k8sDataSource,
	}

	return vcRepoImpl
}

func (vcRepoImpl *VClusterRepositoryImpl) CreateVClusterRecord(vclusterId string, serverStatus string, params *v1.VClusterInfo) error {
	marshal := func(v interface{}) string {
		vv, _ := json.Marshal(v)
		return string(vv)
	}

	var getVClusterNamespaceByIsInit = func(vclusterId string, isInit bool) string {
		if isInit {
			return utils.GetInfraVClusterNamespaceName(vclusterId)
		}
		return utils.GetVClusterNamespaceName(vclusterId)
	}

	createTime := time.Now()

	vcModel := &model.VCluster{
		UserName:        params.Username,
		TenantId:        params.TenantId,
		VClusterId:      vclusterId,
		VClusterName:    params.Name,
		RootClusterName: datasource.DefaultCluster,
		CreateTime:      &createTime,
		DeleteTime:      nil,
		ServerStatus:    serverStatus,
		Comment:         params.Comment,
		StartTime:       nil,
		IsDeleted:       0,
		Namespace:       getVClusterNamespaceByIsInit(vclusterId, params.IsInit),
		InstanceSpec:    marshal(params.OrderDetails.Orders[0].InstanceSpecs),
		InstanceId:      params.OrderDetails.Orders[0].InstanceID,
		ManageBy:        params.ManagerBy,
	}
	err := vcRepoImpl.dbDataSource.CreateVCluster(vcModel)
	if err != nil {
		return errors.Wrapf(err, "create vcluster record failed, vclusterId: %s", vclusterId)
	}

	var vStorages []*model.VStorage
	for _, order := range params.OrderDetails.Orders[1:] {
		// orders 列表的第一个 order 是存储的 VCluster 相关配置的，跳过
		for _, spec := range order.InstanceSpecs {
			if spec.ResourceSpecCode != "" && spec.ResourceSpecParamCode != "" {
				key := spec.ResourceSpecCode + "/" + spec.ResourceSpecParamCode + "-" + vclusterId
				vStorageType, vStorageName := strings.Split(key, "/")[0], strings.Split(key, "/")[1]
				vStorages = append(vStorages, &model.VStorage{
					VClusterID:       vclusterId,
					VStorageType:     vStorageType,
					Name:             vStorageName,
					VStorageCapacity: 0,
				})
			}
		}
	}

	if len(vStorages) > 0 {
		err = vcRepoImpl.dbDataSource.CreateVStorages(vStorages)
		if err != nil {
			return errors.Wrapf(err, "create vstorage record failed, vclusterId: %s", vclusterId)
		}
	}

	instanceSpecs := params.OrderDetails.Orders[0].InstanceSpecs
	gpuDescAnnotationSet := make(map[string]struct{})
	combineAnnotations := make([]string, 0)

	for _, spec := range instanceSpecs {
		if strings.HasPrefix(spec.ResourceSpecParamCode, consts.ResourceQuotaNVIDIAPrefix) {
			if _, exist := gpuDescAnnotationSet[spec.ResourceSpecParamCode]; !exist {
				gpuDescAnnotationSet[spec.ResourceSpecParamCode] = struct{}{}
				combineAnnotations = append(combineAnnotations, spec.ResourceSpecParamCode)
			}
		}
	}
	gpusFiledTmp := strings.Join(combineAnnotations, "+")
	const (
		nvidiaGPUPrefix  = `nvidia`
		huaweiGPUPrefix  = "huawei"
		managedByPrefix  = "managed-by"
		resourceQuotaGPU = `isolation.resourceQuota.quota.requests\\.`
	)

	cutString := func(s, sep string, n int) string {
		parts := strings.Split(s, sep)
		if len(parts) < n {
			return s
		}
		return parts[n]
	}
	var vGpus []*model.VGpu
	if gpusFiledTmp != "" {
		for _, spec := range strings.Split(gpusFiledTmp, "+") {
			resourceName := ""
			if strings.HasPrefix(spec, nvidiaGPUPrefix) {
				resourceName = "nvidia.com/" + cutString(spec, "/", 1)
			}
			if strings.HasPrefix(spec, huaweiGPUPrefix) {
				resourceName = "huawei.com/" + cutString(spec, "/", 1)
			}

			vGpus = append(vGpus, &model.VGpu{
				ClusterID:    vclusterId,
				GpuType:      strings.Split(cutString(spec, "/", 1), ":")[0],
				ResourceName: resourceName,
			})
		}

		if len(vGpus) > 0 {
			err = vcRepoImpl.dbDataSource.CreateVGpus(vGpus)
			if err != nil {
				return errors.Wrapf(err, "create vgpu record failed, vclusterId: %s", vclusterId)
			}
		}
	}

	return nil
}

func extractDomain(inputURL string) string {
	// 解析输入的URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return inputURL
	}

	// 提取URL的Scheme（协议）和Host（域名和端口）
	domain := parsedURL.Scheme + "://" + parsedURL.Host
	return domain
}

func (vcRepoImpl *VClusterRepositoryImpl) getGPUResources(clusterID string, resourceQuota *corev1.ResourceQuota) ([]v1.Gpu, error) {
	vGPUs, err := vcRepoImpl.dbDataSource.FindVClusterGPUs(clusterID)
	if err != nil {
		return nil, err
	}

	var gpuResources []v1.Gpu
	for _, vGPU := range vGPUs {
		resourceName := corev1.ResourceName("requests." + strings.Split(vGPU.ResourceName, ":")[0])
		quota := resourceQuota.Status.Hard[resourceName]
		gpuResources = append(gpuResources, v1.Gpu{
			Type:         vGPU.GpuType,
			ResourceName: vGPU.ResourceName,
			Count:        quota.Value(),
		})
	}

	return gpuResources, nil
}

func (vcRepoImpl *VClusterRepositoryImpl) getStorageResources(clusterID string, resourceQuota *corev1.ResourceQuota) ([]v1.StorageList, error) {
	vStorages, err := vcRepoImpl.dbDataSource.FindVClusterStorages(clusterID)
	if err != nil {
		return nil, err
	}

	var storageResources []v1.StorageList
	for _, vStorage := range vStorages {
		if vStorage.Name == "" {
			continue
		}

		limit := resourceQuota.Status.Hard[corev1.ResourceName(vStorage.Name+".storageclass.storage.k8s.io/requests.storage")]
		storageResources = append(storageResources, v1.StorageList{
			Type:         "hdd",
			Limit:        limit.Value() / 1024 / 1024,
			StorageClass: vStorage.Name,
		})
	}

	return storageResources, nil
}

func (vcRepoImpl *VClusterRepositoryImpl) GetVClusterById(id string) (*model.VCluster, error) {
	return vcRepoImpl.dbDataSource.GetVClusterById(id)
}

func (vcRepoImpl *VClusterRepositoryImpl) UpdateVCluster(vc *model.VCluster) error {
	return vcRepoImpl.dbDataSource.UpdateVClusterSingle(vc)
}

func (vcRepoImpl *VClusterRepositoryImpl) CheckVClusterNameExistByTenantId(tenantId string, vClusterName string) bool {
	return vcRepoImpl.dbDataSource.CheckVClusterNameExistByTenantId(tenantId, vClusterName)
}

func (vcRepoImpl *VClusterRepositoryImpl) CheckVClusterNameExistAndDeleted(vclusterName string, isDeleted int) bool {
	return vcRepoImpl.dbDataSource.CheckVClusterNameExistAndDeleted(vclusterName, isDeleted)
}

const (
	charset         = "abcdefghijklmnopqrstuvwxyz0123456789"
	idLength        = 12
	generateIdTimes = 5
)

// GenerateUniqueVClusterId 生成一个 12 位的 vcluster id
// 格式：vc+uuid
// loft-sh/vcluster 要求 vcluster 的 name 要求必须由小写字母、数字、'-' 组成，并且以字母开头，字母数字结尾
func (vcRepoImpl *VClusterRepositoryImpl) GenerateUniqueVClusterId(ctx context.Context) (string, error) {
	gcpLogger := ctx.Value("logger").(*log.Logger)

	var times int

	for {
		id, err := generateIdWithPrefix("vc", charset, idLength-2)
		if err != nil {
			return "", errors.Wrap(err, "failed to generate unique vclusterId")
		}

		exist := vcRepoImpl.dbDataSource.CheckVClusterExistById(id)
		if !exist {
			gcpLogger.Infof("successfully generated a unique vclusterId: %s", id)
			return id, nil
		}

		times++
		if times >= generateIdTimes {
			break
		}
	}

	return "", errors.New("failed to generate a unique vclusterId")
}

func generateIdWithPrefix(prefix string, charset string, length int) (string, error) {
	id, err := gonanoid.Generate(charset, length)
	if err != nil {
		return "", err
	}
	return prefix + id, nil
}

func (vcRepoImpl *VClusterRepositoryImpl) CheckVClusterExistByTenantId(vClusterId string, tenantId string) bool {
	return vcRepoImpl.dbDataSource.CheckVClusterExistByTenantIdAndVClusterId(vClusterId, tenantId)
}

func (vcRepoImpl *VClusterRepositoryImpl) CheckVClusterExistById(vClusterId string) bool {
	return vcRepoImpl.dbDataSource.CheckVClusterExistById(vClusterId)
}

func (vcRepoImpl *VClusterRepositoryImpl) GetVClusterIdByInstanceId(instanceId string) (string, error) {
	vCluster, err := vcRepoImpl.dbDataSource.GetVClusterByInstanceId(instanceId)
	if err != nil {
		return "", errors.Wrap(err, "dbDataSource.GetVClusterByInstanceIdAndK8sCtx failed")
	}
	return vCluster.VClusterId, nil
}

func (vcRepoImpl *VClusterRepositoryImpl) DeleteVClusterDBResources(id string) error {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		err := vcRepoImpl.dbDataSource.DeleteVClusterById(id)
		if err != nil {
			return errors.Wrap(err, "failed to delete vcluster")
		}
		return nil
	})

	g.Go(func() error {
		err := vcRepoImpl.dbDataSource.DeleteVStorageById(id)
		if err != nil {
			return errors.Wrap(err, "failed to delete vstorage")
		}
		return nil
	})

	g.Go(func() error {
		err := vcRepoImpl.dbDataSource.DeleteVGpuById(id)
		if err != nil {
			return errors.Wrap(err, "failed to delete vgpu")
		}
		return nil
	})

	// Wait for all goroutines to complete and check for errors
	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "failed to delete resources")
	}

	return nil
}

func (vcRepoImpl *VClusterRepositoryImpl) GetRootK8sConfig() (*v1.RootK8sConfig, error) {
	clientSet := vcRepoImpl.k8sDataSource.GetClientSet()
	clientConfig, rowConfig := vcRepoImpl.k8sDataSource.GetConfigs()

	return &v1.RootK8sConfig{
		RootKubeClientSet: clientSet,
		RootClientConfig:  clientConfig,
		RootRawConfig:     rowConfig,
	}, nil
}

func (vcRepoImpl *VClusterRepositoryImpl) GetVClusterClientSet() (*kubernetes.Clientset, error) {
	return vcRepoImpl.k8sDataSource.GetClientSet(), nil
}

func (vcRepoImpl *VClusterRepositoryImpl) GetVClusterInformerFactory() (informers.SharedInformerFactory, error) {
	return vcRepoImpl.k8sDataSource.GetSharedInformerFactory(), nil
}

func (vcRepoImpl *VClusterRepositoryImpl) GetKubeConfig(ctx context.Context, vClusterId string, kubeConnHost string) (*v1.GetVClusterTokenResponse, error) {
	logger := ctx.Value("logger").(*log.Logger)

	// 当前 vcluster 所在的命名空间
	// 示例：
	// vClusterId: vozhafmqfme0
	// namespace: vcluster-vozhafmqfme0   svc: vozhafmqfme0   secret: vc-vozhafmqfme0
	vcNamespace := fmt.Sprintf("vcluster-%s", vClusterId)

	logger.Infof("GetKubeConfig, vcNamespace: %s", vcNamespace)

	// 获取可以操作 VCluster 的 clientset
	vcClientSet, err := vcRepoImpl.k8sDataSource.GenerateVClusterClientSet(vClusterId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vcluster client set, vcluster id: %s", vClusterId)
	}
	if vcClientSet == nil {
		return nil, errors.Errorf("failed to get vcluster client set, vcClientSet is nil, vcluster id: %s", vClusterId)
	}

	kubeSystemNs := "kube-system"
	clusterAdminSa := "cluster-admin-token-sa"

	// 在 VCluster 内部创建一个 ServiceAccount，如果不存在
	_, err = vcClientSet.CoreV1().ServiceAccounts(kubeSystemNs).Get(ctx, clusterAdminSa, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			_, err = vcClientSet.CoreV1().ServiceAccounts(kubeSystemNs).Create(ctx, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:        clusterAdminSa,
					Annotations: map[string]string{"kubernetes.io/service-account.name": vClusterId},
				},
			}, metav1.CreateOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "failed to create cluster-admin-token-sa")
			}
		} else {
			return nil, errors.Wrap(err, "failed to get cluster-admin-token-sa")
		}
	}

	// 在 VCluster 内部创建 ClusterRoleBinding，如果不存在
	_, err = vcClientSet.RbacV1().ClusterRoleBindings().Get(ctx, clusterAdminSa, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			_, err = vcClientSet.RbacV1().ClusterRoleBindings().Create(ctx, &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterAdminSa,
					// Namespace: kubeSystemNs,
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: rbacv1.SchemeGroupVersion.Group,
					Kind:     "ClusterRole",
					Name:     "cluster-admin",
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      "ServiceAccount",
						Name:      clusterAdminSa,
						Namespace: kubeSystemNs,
					},
				},
			}, metav1.CreateOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "failed to create cluster-admin-token-sa")
			}
		} else {
			return nil, errors.Wrap(err, "failed to get cluster-admin-token-sa")
		}
	}

	// 创建 ServiceAccount token
	expirationSeconds := kubeconfigExpirationSeconds
	request, err := vcClientSet.CoreV1().ServiceAccounts(kubeSystemNs).CreateToken(ctx, clusterAdminSa, &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences:         getAudiences(kubeConnHost),
			ExpirationSeconds: &expirationSeconds,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create service account token")
	}

	// 创建 TokenReview
	_, err = vcClientSet.AuthenticationV1().TokenReviews().Create(ctx, &authenticationv1.TokenReview{
		Spec: authenticationv1.TokenReviewSpec{
			Token: request.Status.Token,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token review")
	}

	// 获取可以操作 RootCluster 的 clientset
	rootClusterConfig, err := vcRepoImpl.GetRootK8sConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get root cluster config")
	}
	rootClientSet := rootClusterConfig.RootKubeClientSet

	// 如果 ingress 已经存在就删除
	ingress, err := rootClientSet.NetworkingV1().Ingresses(vcNamespace).Get(ctx, vClusterId, metav1.GetOptions{})
	if err != nil && !kerrors.IsNotFound(err) {
		return nil, errors.Wrap(err, "failed to get ingress")
	}

	err = rootClientSet.NetworkingV1().Ingresses(vcNamespace).Delete(ctx, vClusterId, metav1.DeleteOptions{})
	if err != nil && !kerrors.IsNotFound(err) {
		return nil, errors.Wrap(err, "failed to delete existing ingress")
	}

	// 新建 ingress
	ingress = vcRepoImpl.newIngress(vClusterId, vcNamespace, kubeConnHost)
	createdIngress, err := rootClientSet.NetworkingV1().Ingresses(vcNamespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		ingress = vcRepoImpl.newIngress(vClusterId, utils.GetInfraVClusterNamespaceName(vClusterId), kubeConnHost)
		createdIngress, err = rootClientSet.NetworkingV1().Ingresses(utils.GetInfraVClusterNamespaceName(vClusterId)).Create(ctx, ingress, metav1.CreateOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create ingress")
		}
	}

	logger.Infof("Created Ingress %s in namespace %s for vCluster %s with host %s", createdIngress.Name, vcNamespace, vClusterId, kubeConnHost)

	// 从数据库中获取 VCluster 信息
	vcluster, err := vcRepoImpl.dbDataSource.GetVClusterById(vClusterId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get vcluster id")
	}

	// 生成可以外部访问的 kubeconfig 文件
	kubeconfig, err := vcRepoImpl.newVCKubeconfig(request.Status.Token, kubeConnHost, vcluster.VClusterId, vcluster.VClusterName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to foo")
	}

	logger.Infof("Successfully created kubeconfig for vCluster %s in namespace %s", vClusterId, vcNamespace)

	return kubeconfig, nil
}

// newVCKubeconfig 根据提供的 token、kubeConnHost、vClusterId 和 vClusterName 生成一个新的 KubeConfig
func (vcRepoImpl *VClusterRepositoryImpl) newVCKubeconfig(token string, kubeConnHost string, vClusterId, vClusterName string) (*v1.KubeConfig, error) {
	cluster := api.NewCluster()
	cluster.Server = "https://" + kubeConnHost + "/" + datasource.DefaultCluster + "/" + vClusterId
	cluster.InsecureSkipTLSVerify = true

	authInfo := api.NewAuthInfo()
	authInfo.Token = token

	config := api.NewConfig()
	config.Clusters[vClusterName] = cluster
	config.AuthInfos[vClusterName] = authInfo

	kubeContext := api.NewContext()
	kubeContext.Cluster = vClusterName
	kubeContext.AuthInfo = vClusterName

	config.Contexts[vClusterName] = kubeContext
	config.CurrentContext = vClusterName

	config.APIVersion = "v1"
	config.Kind = "Config"

	kubeConfigBytes, err := clientcmd.Write(*config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write kube config")
	}

	var kubeconfig v1.KubeConfig
	err = yaml.Unmarshal(kubeConfigBytes, &kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal kube config")
	}

	return &kubeconfig, nil
}

// 创建一个 ingress struct
func (vcRepoImpl *VClusterRepositoryImpl) newIngress(name, namespace string, KubeConnectHost string) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                  "nginx",
				"nginx.ingress.kubernetes.io/proxy-buffering":  "off",
				"nginx.ingress.kubernetes.io/upstream-vhost":   strings.Split(KubeConnectHost, ":")[0],
				"nginx.ingress.kubernetes.io/server-alias":     strings.Split(KubeConnectHost, ":")[0],
				"nginx.ingress.kubernetes.io/proxy-body-size":  "5M",
				"ingress.kubernetes.io/force-ssl-redirect":     "false",
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/rewrite-target":   "/$2",
			},
		},
		Spec: networkingv1.IngressSpec{
			// IngressClassName: &ingressClass,
			Rules: []networkingv1.IngressRule{
				{
					Host: strings.Split(KubeConnectHost, ":")[0],
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/" + datasource.DefaultCluster + "/" + name + "(/|$)(.*)",
									// PathType: &pathType,
									PathType: func(s string) *networkingv1.PathType {
										pt := networkingv1.PathType(s)
										return &pt
									}("ImplementationSpecific"),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: name,
											Port: networkingv1.ServiceBackendPort{
												Number: 443,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getAudiences(host string) []string {
	return []string{
		"https://" + host,
		"https://kubernetes.default.svc.cluster.local",
		"https://kubernetes.default.svc",
		"https://kubernetes.default",
	}
}

func (vcRepoImpl *VClusterRepositoryImpl) CheckVClusterNameExist(vclusterName string) bool {
	return vcRepoImpl.dbDataSource.CheckVClusterNameExist(vclusterName)
}

func (vcRepoImpl *VClusterRepositoryImpl) CheckInstanceIdExist(instanceId string) bool {
	return vcRepoImpl.dbDataSource.CheckInstanceIdExist(instanceId)
}

func (vcRepoImpl *VClusterRepositoryImpl) GetVClusterResourceDetails(ctx context.Context, clusterId string) (*vclusterv1.GetVClusterResourceDetailsResponse_Data, error) {
	// 从底层 K8s 获取
	clientSet := vcRepoImpl.k8sDataSource.GetClientSet()
	resourcequota, err := clientSet.CoreV1().ResourceQuotas(utils.GetVClusterNamespaceName(clusterId)).Get(ctx, utils.GetVClusterResourceQuotaName(clusterId), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	//vGPUs, err := vcRepoImpl.dbDataSource.FindVClusterGPUs(clusterId)
	//if err != nil {
	//	return nil, err
	//}

	vStorages, err := vcRepoImpl.dbDataSource.FindVClusterStorages(clusterId)
	if err != nil {
		return nil, err
	}

	vclusterv1Resourcequotas := vcRepoImpl.convertResourceQuota(resourcequota, vStorages)

	var enableServiceExporterFn = func(clusterId string, clientSet *kubernetes.Clientset) bool {

		namespace := utils.GetVClusterNamespaceName(clusterId)
		deploymentName := clusterId

		deployment, err := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			log.Warnf("GetVClusterResourceDetails, failed to get deployment %s: %v", deploymentName, err)
			return false
		}

		// Get the Pod corresponding to the Deployment
		pods, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
		})
		if err != nil {
			log.Warnf("GetVClusterResourceDetails, failed to list pods: %v", err)
			return false
		}

		// Find the syncer container in the Pod
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if container.Name == "syncer" {
					// Extract the CONFIG environment variable
					for _, env := range container.Env {
						if env.Name == "CONFIG" {
							if strings.Contains(env.Value, "ServiceExporter") {

								return true
							}
						}
					}
				}
			}
		}

		return false
	}

	resp := &vclusterv1.GetVClusterResourceDetailsResponse_Data{
		Configurations:  map[string]string{"enableServiceExport": strconv.FormatBool(enableServiceExporterFn(clusterId, clientSet))},
		UtilizationRate: vclusterv1Resourcequotas,
	}

	return resp, nil
}

func (vcRepoImpl *VClusterRepositoryImpl) convertResourceQuota(resourcequota *corev1.ResourceQuota, vStorages []model.VStorage) *vclusterv1.GetVClusterResourceDetailsResponse_Resourcequotas {
	// 默认为每一个 vc 集群分配了 3G 内存和 6 核 cpu
	memHard := resourcequota.Status.Hard[corev1.ResourceLimitsMemory]
	memUsed := resourcequota.Status.Used[corev1.ResourceLimitsMemory]

	cpuHard := resourcequota.Status.Hard[corev1.ResourceLimitsCPU]
	cpuUsed := resourcequota.Status.Used[corev1.ResourceLimitsCPU]

	//memHardNum, memUsedNum := memHard.MilliValue()/1000/1024/1024/1024, memUsed.MilliValue()/1000/1024/1024/1024
	memHardNum, _ := memHard.MilliValue()/1000/1024/1024/1024, memUsed.MilliValue()/1000/1024/1024/1024

	//log.Infof("转换 ResourceQuota 资源：memHard.String(): %s, memUsed.String(): %s", memHard.String(), memUsed.String())
	//log.Infof("转换 ResourceQuota 资源：memHardNum: %d, memUsedNum: %d", memHardNum, memUsedNum)

	//if memHardNum-3 > 0 {
	//	memHardNum = memHardNum - 3
	//}

	//memUsedMB := int64(memUsed.MilliValue()/1000/1024/1024) - 217
	memUsedMB := int64(memUsed.MilliValue() / 1000 / 1024 / 1024)
	memUsedGB := memUsedMB / 1024
	//if memUsedGB-2 > 0 {
	//	memUsedGB = memUsedGB - 2
	//} else {
	//	memUsedGB = 0
	//}

	cpuHardNum, cpuUsedNum := cpuHard.MilliValue()/1000, cpuUsed.MilliValue()/1000
	//log.Infof("转换 ResourceQuota 资源：cpuHard.String(): %s, cpuUsed.String(): %s", cpuHard.String(), cpuUsed.String())
	//log.Infof("转换 ResourceQuota 资源：cpuHardNum: %d, cpuUsedNum: %d", cpuHardNum, cpuUsedNum)

	//if cpuHardNum-6 > 0 {
	//	cpuHardNum = cpuHardNum - 6
	//}
	//if cpuUsedNum-6 > 0 {
	//	cpuUsedNum = cpuUsedNum - 6
	//}

	vclusterv1Resourcequotas := &vclusterv1.GetVClusterResourceDetailsResponse_Resourcequotas{
		Gpu: []*vclusterv1.GetVClusterResourceDetailsResponse_Resourcequotas_Quota{},
		Memory: map[string]float32{
			quotaHard: float32(memHardNum),
			quotaUsed: float32(memUsedGB),
		},
		Cpu: map[string]float32{
			quotaHard: float32(cpuHardNum),
			quotaUsed: float32(cpuUsedNum),
		},
		Storage: []*vclusterv1.GetVClusterResourceDetailsResponse_Resourcequotas_Quota{},
	}

	// 提取 GPU 资源配额信息
	for resourceName, _ := range resourcequota.Status.Hard {
		// 只处理以 nvidia.com 开头的 GPU 资源
		if strings.HasPrefix(string(resourceName), "requests.nvidia.com") {
			gpuType := strings.TrimPrefix(string(resourceName), "requests.")
			gpuHard := resourcequota.Status.Hard[resourceName]
			gpuUsed := resourcequota.Status.Used[resourceName]

			vclusterv1Resourcequotas.Gpu = append(vclusterv1Resourcequotas.Gpu, &vclusterv1.GetVClusterResourceDetailsResponse_Resourcequotas_Quota{
				Name: gpuType, // GPU 类型（去掉前缀后的部分）
				Hard: float32(gpuHard.Value()),
				Used: float32(gpuUsed.Value()),
			})
		}
	}

	for _, storage := range vStorages {
		storageHard := resourcequota.Status.Hard[corev1.ResourceName(storage.Name+".storageclass.storage.k8s.io/requests.storage")]
		storageUsed := resourcequota.Status.Used[corev1.ResourceName(storage.Name+".storageclass.storage.k8s.io/requests.storage")]

		vclusterv1Resourcequotas.Storage = append(vclusterv1Resourcequotas.Storage, &vclusterv1.GetVClusterResourceDetailsResponse_Resourcequotas_Quota{
			Name: storage.VStorageType,
			Hard: float32(storageHard.MilliValue() / 1000 / 1024 / 1024 / 1024),
			Used: float32(storageUsed.MilliValue() / 1000 / 1024 / 1024 / 1024),
		})
	}

	return vclusterv1Resourcequotas
}

func (vcRepoImpl *VClusterRepositoryImpl) GetVClusterContainerID(ctx context.Context, params *v1.GetVClusterContainerIDRequest) (containerId string, err error) {
	clientSet := vcRepoImpl.k8sDataSource.GetClientSet()

	// 构建 LabelSelector
	labelSelector := fmt.Sprintf(
		"dc.com/osm-vcluster-id=%s,dc.com/osm-pod-namespace=%s,dc.com/osm-pod-name=%s",
		params.VClusterId, params.Namespace, params.PodName,
	)

	// 使用 LabelSelector 列出符合条件的 Pod
	podList, err := clientSet.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		params.Logger.Errorf("Error listing pods with label selector %s in namespace %s: %v", labelSelector, params.Namespace, err)
		return "", err
	}

	// 检查是否找到符合条件的 Pod
	if len(podList.Items) == 0 {
		params.Logger.Errorf("No pod found with label selector %s in namespace %s", labelSelector, params.Namespace)
		return "", err
	}

	targetPod := &podList.Items[0] // 只有一个符合条件的 Pod

	// 获取指定容器的 ID
	for _, containerStatus := range targetPod.Status.ContainerStatuses {
		if containerStatus.Name == params.ContainerName {
			return containerStatus.ContainerID, nil
		}
	}

	params.Logger.Errorf("Container %s not found in pod %s", params.ContainerName, params.PodName)
	return "", err
}

func (vcRepoImpl *VClusterRepositoryImpl) DeleteVClusterNamespace(ctx context.Context, clusterId string) error {
	clientSet := vcRepoImpl.k8sDataSource.GetClientSet()
	err := clientSet.CoreV1().Namespaces().Delete(ctx, utils.GetVClusterNamespaceName(clusterId), metav1.DeleteOptions{})
	if err != nil {
		// namespace 不存在，则尝试删除 infra-vcluster-id 的 namespace
		// 并且直接返回 nil
		if kerrors.IsNotFound(err) {
			_ = clientSet.CoreV1().Namespaces().Delete(ctx, utils.GetInfraVClusterNamespaceName(clusterId), metav1.DeleteOptions{})
			return nil
		}
		return errors.Wrapf(err, "failed to delete namespace %s", utils.GetVClusterNamespaceName(clusterId))
	}
	return nil
}
