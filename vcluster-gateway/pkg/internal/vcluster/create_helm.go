// Package forkvcluster 包封装了使用 Helm 创建/升级 vCluster 的功能
// 大部分代码都 fork 自 github.com/loft-sh/vcluster/pkg/cli/create_helm.go
// 该包的主要功能包括：
// - 验证和初始化 vCluster 创建/升级所需的参数和配置
// - 使用 Helm 部署和升级 vCluster
// - 确保命名空间存在并创建必要的资源
// - 获取和设置 Kubernetes 版本、Service CIDR 等信息
// - 通过多种可选参数（Option 模式）配置 CreateHelm 结构体
package forkvcluster

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/loft-sh/log"
	"github.com/loft-sh/vcluster/cmd/vclusterctl/cmd/app/localkubernetes"
	"github.com/loft-sh/vcluster/pkg/embed"
	"github.com/loft-sh/vcluster/pkg/helm"
	"github.com/loft-sh/vcluster/pkg/telemetry"
	"github.com/loft-sh/vcluster/pkg/upgrade"
	"github.com/loft-sh/vcluster/pkg/util"
	"github.com/loft-sh/vcluster/pkg/util/cliconfig"
	"github.com/loft-sh/vcluster/pkg/util/servicecidr"
	"github.com/loft-sh/vcluster/pkg/values"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	v1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	"vcluster-gateway/pkg/internal/utils"
)

var allowedDistros = []string{"k3s", "k0s", "k8s", "eks"}

// CreateHelmOption 使用 Option模式创建 CreateHelm
type CreateHelmOption func(*CreateHelm)

// CreateHelm 封装了调用 helm create/upgrade 的参数
// 其中 GlobalFlags和 CreateOptions 的定义基本来自 loft-sh/vcluster 源码
type CreateHelm struct {
	*GlobalFlags
	*CreateOptions
	log log.Logger

	chartValues string

	localCluster     bool
	kubeClientConfig clientcmd.ClientConfig
	kubeClient       kubernetes.Clientset
	rawConfig        clientcmdapi.Config
}

func NewCreateHelm(options ...CreateHelmOption) *CreateHelm {
	ch := &CreateHelm{}

	for _, option := range options {
		option(ch)
	}

	return ch
}

func WithGlobalFlags(globalFlags *GlobalFlags) CreateHelmOption {
	return func(ch *CreateHelm) {
		ch.GlobalFlags = globalFlags
	}
}

func WithLogger(logger log.Logger) CreateHelmOption {
	return func(ch *CreateHelm) {
		ch.log = logger
	}
}

func WithCreateOptions(options *CreateOptions) CreateHelmOption {
	return func(ch *CreateHelm) {
		// options.KubeConfigContextName = ch.KubeConfigContextName
		ch.CreateOptions = options
	}
}

func WithCreateChartRepo(chartRepo string) CreateHelmOption {
	return func(ch *CreateHelm) {
		switch chartRepo {
		case "localhost", "127.0.0.1", "0.0.0.0", "":
			ch.LocalChartDir = "./charts/vcluster-k8s/"
		default:
			ch.ChartRepo = chartRepo
		}
	}
}

func WithCreateK8sConfig(config *v1.RootK8sConfig) CreateHelmOption {
	return func(ch *CreateHelm) {
		ch.kubeClientConfig = *config.RootClientConfig
		ch.kubeClient = *config.RootKubeClientSet
		ch.rawConfig = *config.RootRawConfig
	}
}

// ValidateAndInitialize 调用的方法大部都 fork 自 loft-sh/vcluster 源码
// 这里主要封装一下，暴露给使用者一个简单的方法
func (cmd *CreateHelm) ValidateAndInitialize(ctx context.Context, vClusterId string, nodeSelector map[string]string, infra bool) error {
	// 验证是否使用了不适用于 OSS 的 vCluster.Pro 参数
	if err := cmd.validateOSSFlags(); err != nil {
		return err
	}

	// 确保使用正确的 Chart 版本
	cmd.insureCorrectVersion()

	// 确保命名空间存在，如果不存在则创建一个新的命名空间
	if err := cmd.ensureNamespace(ctx, vClusterId, nodeSelector, infra); err != nil {
		return err
	}

	// 获取服务 CIDR，如果没有设置则获取默认值
	cmd.getServiceCIDR(ctx)

	// 获取 Kubernetes 版本信息
	kubernetesVersion, err := cmd.getKubernetesVersion()
	if err != nil {
		cmd.log.Warn("Failed to get Kubernetes version: ", err.Error())
	}

	// 将 CreateHelm.CreateOptions 转换为 ChartOptions
	chartOptions, err := cmd.toChartOptions(kubernetesVersion, cmd.log)
	if err != nil {
		cmd.log.Warn("Failed to get default values: ", err.Error())
	}

	// 创建新的日志记录器
	logger := logr.New(cmd.log.LogrLogSink())

	// 获取默认的 release 值
	cmd.chartValues, err = values.GetDefaultReleaseValues(chartOptions, logger)
	if err != nil {
		cmd.log.Error("Failed to get default values: ", err.Error())
	}

	// 将 CreateHelm.CreateOptions.Values 字段格式化
	if err = cmd.formatValues(); err != nil {
		return err
	}

	return nil
}

// Deploy 调用 helm client 创建 vCluster
func (cmd *CreateHelm) Deploy(ctx context.Context, vClusterId string, valueFilenames []string) (err error) {
	// check if there is a vcluster directory already
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get current work directory: %w", err)
	}
	if _, err := os.Stat(filepath.Join(workDir, cmd.ChartName)); err == nil {
		return fmt.Errorf("aborting vcluster creation. Current working directory contains a file or a directory with the name equal to the vcluster chart name - \"%s\". Please execute vcluster create command from a directory that doesn't contain a file or directory named \"%s\"", cmd.ChartName, cmd.ChartName)
	}

	if cmd.LocalChartDir == "" {
		chartEmbedded := false
		if cmd.ChartVersion == upgrade.GetVersion() { // use embedded chart if default version
			embeddedChartName := fmt.Sprintf("%s-%s.tgz", cmd.ChartName, upgrade.GetVersion())
			// not using filepath.Join because the embed.FS separator is not OS specific
			embeddedChartPath := fmt.Sprintf("charts/%s", embeddedChartName)
			embeddedChartFile, err := embed.Charts.ReadFile(embeddedChartPath)
			if err != nil && errors.Is(err, fs.ErrNotExist) {
				cmd.log.Infof("Chart not embedded: %q, pulling from helm repository.", err)
			} else if err != nil {
				cmd.log.Errorf("Unexpected error while accessing embedded file: %q", err)
			} else {
				temp, err := os.CreateTemp("", fmt.Sprintf("%s%s", embeddedChartName, "-"))
				if err != nil {
					cmd.log.Errorf("Error creating temp file: %v", err)
				} else {
					defer temp.Close()
					defer os.Remove(temp.Name())
					_, err = temp.Write(embeddedChartFile)
					if err != nil {
						cmd.log.Errorf("Error writing package file to temp: %v", err)
					}
					cmd.LocalChartDir = temp.Name()
					chartEmbedded = true
					cmd.log.Infof("Using embedded chart: %q", embeddedChartName)
				}
			}
		}

		// rewrite chart location, this is an optimization to avoid
		// downloading the whole index.yaml and parsing it
		if !chartEmbedded && cmd.ChartRepo == loftChartRepo && cmd.ChartVersion != "" { // specify versioned path to repo url
			cmd.LocalChartDir = loftChartRepo + "/charts/" + cmd.ChartName + "-" + strings.TrimPrefix(cmd.ChartVersion, "v") + ".tgz"
		}
	}

	cmd.log.Infof("Create/Upgrade vcluster %s...", vClusterId)

	// 强制更新 Deployment, 触发 informer 的 updateEvent
	deployment, err := cmd.kubeClient.AppsV1().Deployments("vcluster-"+vClusterId).Get(ctx, vClusterId, metav1.GetOptions{})
	if err == nil || deployment != nil {
		if deployment.Annotations == nil {
			deployment.Annotations = make(map[string]string)
		}

		deployment.Annotations["upgradeTime"] = time.Now().Format(time.RFC3339)

		_, _ = cmd.kubeClient.AppsV1().Deployments("vcluster-"+vClusterId).Update(ctx, deployment, metav1.UpdateOptions{})
	}

	// we have to upgrade / install the chart
	err = helm.NewClient(&cmd.rawConfig, cmd.log, defaultHelmBinaryPath).Upgrade(ctx, vClusterId, cmd.Namespace, helm.UpgradeOptions{
		Chart:       cmd.ChartName,
		Repo:        cmd.ChartRepo,
		Version:     cmd.ChartVersion,
		Path:        cmd.LocalChartDir,
		Values:      cmd.chartValues, // 调用 vcluster 源码中的 GetDefaultReleaseValues 生成的配置
		ValuesFiles: valueFilenames,  // 通过 -f 指定的 values.yaml 文件，需要添加到 ValuesFiles 里
		SetValues:   cmd.SetValues,   // 如果使用 key=value 的方式设置参数，需要添加到 SetValues 里
	})
	if err != nil {
		cmd.log.Error("Error while installing chart: ", err.Error())
		return err
	}

	cmd.log.Infof("Successfully created virtual cluster %s in namespace %s.", vClusterId, cmd.Namespace)
	return nil
}

// validateOSSFlags 检查是否使用了不适用于 OSS 的 vCluster.Pro 参数，并返回相应的错误信息
func (cmd *CreateHelm) validateOSSFlags() error {
	if cmd.Project != "" {
		return fmt.Errorf("cannot use --project as you are not connected to a vCluster.Pro instance." + loginText)
	}
	if cmd.Cluster != "" {
		return fmt.Errorf("cannot use --cluster as you are not connected to a vCluster.Pro instance." + loginText)
	}
	if cmd.Template != "" {
		return fmt.Errorf("cannot use --template as you are not connected to a vCluster.Pro instance." + loginText)
	}
	if cmd.TemplateVersion != "" {
		return fmt.Errorf("cannot use --template-version as you are not connected to a vCluster.Pro instance." + loginText)
	}
	if len(cmd.Links) > 0 {
		return fmt.Errorf("cannot use --link as you are not connected to a vCluster.Pro instance." + loginText)
	}
	if cmd.Params != "" {
		return fmt.Errorf("cannot use --params as you are not connected to a vCluster.Pro instance." + loginText)
	}
	if len(cmd.SetParams) > 0 {
		return fmt.Errorf("cannot use --set-params as you are not connected to a vCluster.Pro instance." + loginText)
	}

	return nil
}

// insureCorrectVersion 确保使用正确的 Chart 版本，如果版本是开发版本，则将其设置为空字符串
func (cmd *CreateHelm) insureCorrectVersion() {
	if cmd.ChartVersion == upgrade.DevelopmentVersion {
		cmd.ChartVersion = ""
	}
}

// ensureNamespace 确保命名空间存在，如果不存在则创建一个新的命名空间
func (cmd *CreateHelm) ensureNamespace(ctx context.Context, id string, nodeSelector map[string]string, infra bool) error {
	cmd.log.Infof("ensureNamespace, cmd.Namespace: %s", cmd.Namespace)

	var err error
	if cmd.Namespace == "" {
		cmd.log.Infof("ensureNamespace, cmd.Namespace is empty,")
		cmd.log.Infof("ensureNamespace will use namespace %s to create the vcluster...", cmd.Namespace)

		//if infra {
		//	cmd.Namespace = "infra-" + "vc-" + id
		//} else {
		//	cmd.Namespace = "vcluster-" + id
		//}
		if infra {
			cmd.Namespace = utils.GetInfraVClusterNamespaceName(id)
		} else {
			cmd.Namespace = utils.GetVClusterNamespaceName(id)
		}
	}

	// make sure namespace exists
	namespace, err := cmd.kubeClient.CoreV1().Namespaces().Get(ctx, cmd.Namespace, metav1.GetOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return cmd.createNamespace(ctx, nodeSelector)
		} else if !kerrors.IsForbidden(err) {
			return err
		}
	} else if namespace.DeletionTimestamp != nil {
		cmd.log.Infof("Waiting until namespace is terminated...")
		err := wait.PollUntilContextTimeout(ctx, time.Second, time.Minute*2, false, func(ctx context.Context) (bool, error) {
			namespace, err := cmd.kubeClient.CoreV1().Namespaces().Get(ctx, cmd.Namespace, metav1.GetOptions{})
			if err != nil {
				if kerrors.IsNotFound(err) {
					return true, nil
				}

				return false, err
			}

			return namespace.DeletionTimestamp == nil, nil
		})
		if err != nil {
			return err
		}

		// create namespace
		return cmd.createNamespace(ctx, nodeSelector)
	}

	return nil
}

// createNamespace 创建命名空间
func (cmd *CreateHelm) createNamespace(ctx context.Context, nodeSelector map[string]string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: cmd.Namespace,
			Annotations: map[string]string{
				createdByVClusterAnnotation: "true",
			},
		},
	}

	if len(nodeSelector) > 0 {
		nodeSelectorSlice := make([]string, 0, len(nodeSelector))
		for k, v := range nodeSelector {
			nodeSelectorSlice = append(nodeSelectorSlice, k+"="+v)
		}

		ns.Annotations["scheduler.alpha.kubernetes.io/node-selector"] = strings.Join(nodeSelectorSlice, ",")
	}

	// try to create the namespace
	cmd.log.Infof("Creating namespace %s", cmd.Namespace)
	_, err := cmd.kubeClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "create namespace")
	}
	return nil
}

// getServiceCIDR 获取 Service CIDR，如果没有设置则获取默认值
func (cmd *CreateHelm) getServiceCIDR(ctx context.Context) {
	if cmd.CIDR == "" {
		cidr, warning := servicecidr.GetServiceCIDR(ctx, &cmd.kubeClient, cmd.Namespace)
		if warning != "" {
			cmd.log.Warn(warning)
		}
		cmd.CIDR = cidr
	}
}

// getKubernetesVersion 获取 k8s 版本信息
func (cmd *CreateHelm) getKubernetesVersion() (*version.Info, error) {
	var (
		kubernetesVersion *version.Info
		err               error
	)
	if cmd.KubernetesVersion != "" {
		if cmd.KubernetesVersion[0] != 'v' {
			cmd.KubernetesVersion = "v" + cmd.KubernetesVersion
		}

		if !semver.IsValid(cmd.KubernetesVersion) {
			return nil, fmt.Errorf("please use valid semantic versioning format, e.g. vX.X")
		}

		majorMinorVer := semver.MajorMinor(cmd.KubernetesVersion)

		if splittedVersion := strings.Split(cmd.KubernetesVersion, "."); len(splittedVersion) > 2 {
			cmd.log.Warnf("currently we only support major.minor version (%s) and not the patch version (%s)", majorMinorVer, cmd.KubernetesVersion)
		}

		parsedVersion, err := values.ParseKubernetesVersionInfo(majorMinorVer)
		if err != nil {
			return nil, err
		}

		kubernetesVersion = &version.Info{
			Major: parsedVersion.Major,
			Minor: parsedVersion.Minor,
		}
	}

	if kubernetesVersion == nil {
		kubernetesVersion, err = cmd.kubeClient.DiscoveryClient.ServerVersion()
		if err != nil {
			return nil, err
		}
	}

	return kubernetesVersion, nil
}

// toChartOptions 将 CreateHelm.CreateOption 转为 ChartOptions
func (cmd *CreateHelm) toChartOptions(kubernetesVersion *version.Info, log log.Logger) (*values.ChartOptions, error) {
	if !util.Contains(cmd.Distro, allowedDistros) {
		return nil, fmt.Errorf("unsupported distro %s, please select one of: %s", cmd.Distro, strings.Join(allowedDistros, ", "))
	}

	if cmd.ChartName == "vcluster" && cmd.Distro != "k3s" {
		cmd.ChartName += "-" + cmd.Distro
	}

	// check if we're running in isolated mode
	if cmd.Isolate {
		// In this case, default the ExposeLocal variable to false
		// as it will always fail creating a vcluster in isolated mode
		cmd.ExposeLocal = false
	}

	// check if we should create with node port
	clusterType := localkubernetes.DetectClusterType(&cmd.rawConfig)
	if cmd.ExposeLocal && clusterType.LocalKubernetes() {
		log.Infof("Detected local kubernetes cluster %s. Will deploy vcluster with a NodePort & sync real nodes", clusterType)
		cmd.localCluster = true
	}

	return &values.ChartOptions{
		ChartName:          cmd.ChartName,
		ChartRepo:          cmd.ChartRepo,
		ChartVersion:       cmd.ChartVersion,
		CIDR:               cmd.CIDR,
		DisableIngressSync: cmd.DisableIngressSync,
		Expose:             cmd.Expose,
		SyncNodes:          cmd.localCluster,
		NodePort:           cmd.localCluster,
		Isolate:            cmd.Isolate,
		KubernetesVersion: values.Version{
			Major: kubernetesVersion.Major,
			Minor: kubernetesVersion.Minor,
		},
		DisableTelemetry:    cliconfig.GetConfig(log).TelemetryDisabled,
		InstanceCreatorType: "vclusterctl",
		MachineID:           telemetry.GetMachineID(log),
	}, nil
}

// generateValues 格式化 CreateHelm.CreateOption.Values 字段
// Notice: 目前我们创建 VCluster 时 Values 字段是空的
func (cmd *CreateHelm) formatValues() error {
	getBase64DecodedString := func(values string) (string, error) {
		strDecoded, err := base64.StdEncoding.DecodeString(values)
		if err != nil {
			return "", err
		}
		return string(strDecoded), nil
	}

	var newExtraValues []string
	for _, value := range cmd.Values {
		decodedString, err := getBase64DecodedString(value)
		// ignore decoding errors and treat it as non-base64 string
		if err != nil {
			newExtraValues = append(newExtraValues, value)
			continue
		}
		// write a temporary values file
		tempFile, err := os.CreateTemp("", "")
		tempValuesFile := tempFile.Name()
		if err != nil {
			return errors.Wrap(err, "create temp values file")
		}
		defer func(name string) {
			_ = os.Remove(name)
		}(tempValuesFile)
		_, err = tempFile.Write([]byte(decodedString))
		if err != nil {
			return errors.Wrap(err, "write values to temp values file")
		}
		err = tempFile.Close()
		if err != nil {
			return errors.Wrap(err, "close temp values file")
		}
		// setting new file to extraValues slice to process it further.
		newExtraValues = append(newExtraValues, tempValuesFile)
	}
	// resetting this as the base64 encoded strings should be removed and only valid file names should be kept.
	cmd.Values = newExtraValues

	return nil
}

// CreateOptions holds the create cmd options
type CreateOptions struct {
	KubeConfigContextName string
	ChartVersion          string
	ChartName             string
	ChartRepo             string
	LocalChartDir         string
	Distro                string
	CIDR                  string
	Values                []string
	SetValues             []string
	DeprecatedExtraValues []string

	KubernetesVersion string

	CreateNamespace    bool
	DisableIngressSync bool
	UpdateCurrent      bool
	Expose             bool
	ExposeLocal        bool

	Connect bool
	Upgrade bool
	Isolate bool

	// Pro
	Project         string
	Cluster         string
	Template        string
	TemplateVersion string
	Links           []string
	Annotations     []string
	Labels          []string
	Params          string
	SetParams       []string
	DisablePro      bool
}

func NewDefaultCreateOptions() *CreateOptions {
	return &CreateOptions{
		// KubeConfigContextName: ch.GlobalFlags.Context,
		ChartRepo:          defaultChartRepo,
		ChartVersion:       defaultChartVersion,
		ChartName:          defaultChartName,
		LocalChartDir:      defaultLocalChartDir,
		Distro:             defaultDistro,
		CIDR:               defaultCIDR,
		KubernetesVersion:  defaultKubernetesVersion,
		CreateNamespace:    defaultCreateNamespace,
		DisableIngressSync: defaultDisableIngressSync,
		UpdateCurrent:      defaultUpdateCurrent,
		Expose:             defaultExpose,
		ExposeLocal:        defaultExposeLocal,
		Connect:            defaultConnect,
		Upgrade:            defaultUpgrade,
		Isolate:            defaultIsolate,
		Project:            defaultProject,
		Cluster:            defaultCluster,
		Template:           defaultTemplate,
		TemplateVersion:    defaultTemplateVersion,
		Params:             defaultParams,
		DisablePro:         defaultDisablePro,
	}
}
