package v1

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/repo"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientset "k8s.io/client-go/kubernetes"
	kubevirtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
)

const (
	STATE_PREPAREING = "preparing"
	STATE_STARTING   = "starting"
	STATE_RUNNING    = "running"
	STATE_ABNORMITY  = "abnormity"
	STATE_STOPPING   = "stopping"
	STATE_RESTARTING = "restarting"
	STATE_STOPPED    = "stopped"
	STATE_DELETING   = "deleting"
	STATE_DELETED    = "deleted"
)

const (
	DELETED_YES     = "1"
	DELETED_YES_INT = 1
	DELETED_NO      = "0"
	DELETED_NO_INT  = 0
)

const (
	STORAGE_SYSTEM          = "system"
	STORAGE_DATA            = "data"
	IMAGE_STORAGE_PRE       = "gcp-"
	VIRTUALSERVER_NAMESPACE = "virtualserver"
	DEFAULT_SSH_PORT        = 22
	VOLUME_MOUNTPOINT_PRE   = "/dev/"
)

const (
	SYSTEM_VOLUME_NAME    = "system-volume"
	CLOUDINIT_VOLUME_NAME = "cloud-init-volume"
	DEFAULT_NAME          = "default"
	DEFAULT_MACHINE_TYPE  = "q35"
	PVC_PURPOSE_SYSTEM    = "system"
	PVC_PURPOSE_DATA      = "data"
)

const (
	SCHEDULE_TOLERATIONS_KEY   = "chedulerserver"
	SCHEDULE_TOLERATIONS_VALUE = "kubevirt"
)

const (
	SSHPORT_MIN = 31000
	SSHPORT_MAX = 32000
)

var (
	MinDiskCapacity int64
	MaxDiskCapacity int64
)

type VirtualServers []VirtualServer

// 模型设计, 参数设计，参考k8s进行处理
type VirtualServer struct {
	Name             string            `json:"name"       binding:"required"`
	Desc             string            `json:"desc"       binding:"required"`
	Image            string            `json:"image"       binding:"required"`
	ProductCode      string            `json:"productCode"       binding:"required"`
	Storage          *Storage          `json:"storage,omitempty"`
	CloudInit        *CloudInit        `json:"cloudinit,omitempty"`
	Labels           map[string]string `json:"-"` // 创建虚拟机时需要打的标签，与数据创建模型和返回模型无关，无需json解析
	CreateUser       string            `json:"-"` // 创建虚拟机时需要打的标签，与数据创建模型和返回模型无关，无需json解析
	StartRun         bool              `json:"-"` // 创建完虚拟机是否启动，默认为true
	KubevirtResource `json:"-"`
	CreateTime       *time.Time `json:"createTime,omitempty"`
	StartedTime      *time.Time `json:"startedTime,omitempty"`
	DeleteTime       *time.Time `json:"deleteTime,omitempty"`
	State            State      `json:"state"       binding:""` // 创建时无需下发，查询时需要返回
	Deleted          int32      `json:"-"`
}

type BSMVirtualServerReq struct {
	Name          string          `json:"name"       binding:"required"`
	Desc          string          `json:"desc"       binding:"required"`
	Image         string          `json:"image"       binding:"required"`
	Storage       *Storage        `json:"storage,omitempty"`
	CloudInit     *CloudInit      `json:"cloudinit,omitempty"`
	InstanceId    string          `json:"instanceId"       binding:"required"`
	ProductInfo   productInfo     `json:"productInfo"       binding:"required"`
	ProductConfig []productConfig `json:"productConfig"       binding:"required"`
}

type OSMVirtualServerReq struct {
	Name        string      `json:"name"       binding:"required"`
	Desc        string      `json:"desc"       binding:"required"`
	Image       string      `json:"image"       binding:"required"`
	Storage     *Storage    `json:"storage,omitempty"`
	CloudInit   *CloudInit  `json:"cloudinit,omitempty"`
	ProductInfo productInfo `json:"productInfo"       binding:"required"`
}

type productConfig struct {
	Key   string `json:"configKey"       binding:"required"`
	Value string `json:"configValue"       binding:"required"`
}

type productInfo struct {
	ProductCode string `json:"productCode"       binding:"required"`
	// ProductCategory string `json:"productCode"`
}

type KubevirtResource struct {
	ResourceList corev1.ResourceList
	Gpus         []kubevirtv1.GPU
}

type Storage struct {
	SystemStorage StorageEntity   `json:"systemStorage"       binding:"required"`
	DataStorage   []StorageEntity `json:"dataStorage"       binding:"required"`
}

func (s *Storage) Init(vm *VirtualServer) {
	// 初始化系统盘存储
	if s.SystemStorage.IsNew {
		s.SystemStorage.PvcName = s.SystemStorage.NewPvcName(vm, PVC_PURPOSE_SYSTEM, 0)
	}
	s.SystemStorage.CloneSource = vm.Image // getImageOriginStorage(vm.Image)
	s.SystemStorage.Labels = map[string]string{
		LabelGCPCreateUserKey:       vm.CreateUser,
		LabelGCPPurposeKey:          LabelGCPPurposeSystem,
		LabelGCPCreateAppKey:        LabelGCPCreateApp,
		LabelGCPBindInstanceIdKey:   vm.State.InstanceId,
		LabelGCPBindInstanceNameKey: vm.Name,
		LabelGCPCreateByInstance:    vm.State.InstanceId,
	}
	if !s.SystemStorage.IsNew {
		s.SystemStorage.Labels = map[string]string{
			// 非新磁盘，做增量标签
			LabelGCPBindInstanceIdKey:   vm.State.InstanceId,
			LabelGCPBindInstanceNameKey: vm.Name,
		}
	}
	// 是否随实例释放
	if s.SystemStorage.ReleaseWithInstance {
		s.SystemStorage.Labels = utils.MergeMaps(s.SystemStorage.Labels, map[string]string{
			LabelGCPReleaseByInstance: vm.State.InstanceId,
		})
	}
	// 初始化数据盘存储
	datas := []StorageEntity{}
	for i, ds := range s.DataStorage {
		if ds.IsNew {
			ds.PvcName = s.SystemStorage.NewPvcName(vm, PVC_PURPOSE_DATA, i)
		}
		ds.Labels = map[string]string{
			LabelGCPCreateUserKey:       vm.CreateUser,
			LabelGCPPurposeKey:          LabelGCPPurposeData,
			LabelGCPCreateAppKey:        LabelGCPCreateApp,
			LabelGCPBindInstanceIdKey:   vm.State.InstanceId,
			LabelGCPBindInstanceNameKey: vm.Name,
			LabelGCPCreateByInstance:    vm.State.InstanceId,
		}
		if !ds.IsNew {
			ds.Labels = map[string]string{
				// 非新磁盘，做增量标签
				LabelGCPBindInstanceIdKey:   vm.State.InstanceId,
				LabelGCPBindInstanceNameKey: vm.Name,
			}
		}
		// 是否随实例释放
		if ds.ReleaseWithInstance {
			ds.Labels = utils.MergeMaps(ds.Labels, map[string]string{
				LabelGCPReleaseByInstance: vm.State.InstanceId,
			})
		}
		datas = append(datas, ds)
	}
	s.DataStorage = datas
}

func (virtualServer *VirtualServer) Namespace() string {
	return VIRTUALSERVER_NAMESPACE
}

// 构造kubevirt虚拟服务器模型(后续考虑挪到apis/v1初始化模型时完成)
func (virtualServer *VirtualServer) ConstructKubevirtVm() *kubevirtv1.VirtualMachine {
	// 存储卷
	volumes := []kubevirtv1.Volume{}
	// 磁盘
	disks := []kubevirtv1.Disk{}
	// 系统存储卷
	systemVolume := kubevirtv1.Volume{
		Name: virtualServer.Storage.SystemStorage.PvcName,
		VolumeSource: kubevirtv1.VolumeSource{
			PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
				PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: virtualServer.Storage.SystemStorage.PvcName,
				},
			},
		},
	}
	systemDisk := kubevirtv1.Disk{
		Name: virtualServer.Storage.SystemStorage.PvcName,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: kubevirtv1.DiskBusSCSI,
			},
		},
	}
	volumes = append(volumes, systemVolume)
	disks = append(disks, systemDisk)
	// 数据盘存储卷
	for _, v := range virtualServer.Storage.DataStorage {
		dataVolume := kubevirtv1.Volume{
			Name: v.PvcName,
			VolumeSource: kubevirtv1.VolumeSource{
				PersistentVolumeClaim: &kubevirtv1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: v.PvcName,
					},
				},
			},
		}
		dataDisk := kubevirtv1.Disk{
			Name: v.PvcName,
			DiskDevice: kubevirtv1.DiskDevice{
				Disk: &kubevirtv1.DiskTarget{
					Bus: kubevirtv1.DiskBusSCSI,
				},
			},
		}
		volumes = append(volumes, dataVolume)
		disks = append(disks, dataDisk)
	}
	// cloudInit初始化存储卷
	cloudInitVolume := kubevirtv1.Volume{
		Name: CLOUDINIT_VOLUME_NAME,
		VolumeSource: kubevirtv1.VolumeSource{
			CloudInitNoCloud: &kubevirtv1.CloudInitNoCloudSource{
				// UserData: constructUserData(virtualServer.CloudInit.User, virtualServer.CloudInit.PassWorld, virtualServer.CloudInit.SshKey),
				UserDataBase64: constructUserData(virtualServer.CloudInit.User, virtualServer.CloudInit.PassWorld, virtualServer.CloudInit.SshKey),
			},
		},
	}
	cloudInitDisk := kubevirtv1.Disk{
		Name: CLOUDINIT_VOLUME_NAME,
		DiskDevice: kubevirtv1.DiskDevice{
			Disk: &kubevirtv1.DiskTarget{
				Bus: kubevirtv1.DiskBusVirtio,
			},
		},
	}
	volumes = append(volumes, cloudInitVolume)
	disks = append(disks, cloudInitDisk)
	// 网络配置
	inf := kubevirtv1.Interface{
		Name: DEFAULT_NAME,
		InterfaceBindingMethod: kubevirtv1.InterfaceBindingMethod{
			Bridge: &kubevirtv1.InterfaceBridge{},
		}}

	networkInterfaces := []kubevirtv1.Interface{inf}

	networks := []kubevirtv1.Network{}
	networks = append(networks, kubevirtv1.Network{
		Name: DEFAULT_NAME,
		NetworkSource: kubevirtv1.NetworkSource{
			Pod: &kubevirtv1.PodNetwork{},
		},
	})

	// kubevirt虚拟机模型
	vm := &kubevirtv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      virtualServer.KubevirtVmName(),
			Namespace: virtualServer.Namespace(),
			Labels:    virtualServer.Labels,
		},
		Spec: kubevirtv1.VirtualMachineSpec{
			Running: &virtualServer.StartRun,
			Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: virtualServer.Labels,
				},
				Spec: kubevirtv1.VirtualMachineInstanceSpec{
					Volumes: volumes,
					Domain: kubevirtv1.DomainSpec{
						Devices: kubevirtv1.Devices{
							GPUs:       virtualServer.KubevirtResource.Gpus,
							Disks:      disks,
							Interfaces: networkInterfaces,
						},
						Resources: kubevirtv1.ResourceRequirements{
							Requests: virtualServer.KubevirtResource.ResourceList,
							Limits:   virtualServer.KubevirtResource.ResourceList,
						},
						Machine: &kubevirtv1.Machine{
							Type: DEFAULT_MACHINE_TYPE,
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      SCHEDULE_TOLERATIONS_KEY,
							Operator: corev1.TolerationOpEqual,
							Value:    SCHEDULE_TOLERATIONS_VALUE,
							Effect:   corev1.TaintEffectNoSchedule,
						},
					},
					Networks: networks,
				},
			},
		},
	}
	return vm
}

//go:embed templates
var tmpl embed.FS

func constructUserData(username, pwd string, sshKey []string) string {
	// userData := fmt.Sprintf("#cloud-config\ndisable_root: false\nssh_pwauth: true\nuser: %s\npassword: %s\nkeyboard:\n  layout: us\nlocale: en_US", username, pwd)
	// return userData
	userData := constructCloudinit(username, pwd, sshKey)
	userDataBase64 := base64.StdEncoding.EncodeToString([]byte(userData))
	return userDataBase64
	// return constructCloudinit(username, pwd, sshKey)
}

func constructCloudinit(username, pwd string, sshKeys []string) string {
	temp, err := template.ParseFS(tmpl, "templates/cloudinit.yaml")
	if err != nil {
		log.Infof("parse cloudinit template err %v", err)
		return ""
	}
	type CloudInit struct {
		UserName string
		PassWord string
		SshKeys  []template.HTML
	}
	decSshKeys := []template.HTML{}
	for _, key := range sshKeys {
		// decSshKey, err := base64.StdEncoding.DecodeString(key)
		// if err == nil {
		// decSshKeys = append(decSshKeys, template.HTML(string(decSshKey)))
		decSshKeys = append(decSshKeys, template.HTML(string(key)))
		// }
	}

	// decPwd, err := base64.StdEncoding.DecodeString(pwd)
	// if err == nil {
	// 	pwd = string(decPwd)
	// }

	cfg := &CloudInit{
		UserName: username,
		PassWord: pwd,
		SshKeys:  decSshKeys,
	}
	var w bytes.Buffer
	if err := temp.ExecuteTemplate(&w, "cloudinit.yaml", cfg); err != nil {
		log.Infof("ExecuteTemplate cloudinit err %v", err)
		return ""
	}
	return w.String()
}

// 构造存储卷名称（区分不同用户下可以创建相同名称的存储卷）
func (s *StorageEntity) NewPvcName(vm *VirtualServer, purpose string, number int) string {
	if number == 1 {
		return fmt.Sprintf("%s-%s--%s-%s", strings.ToLower(vm.CreateUser), "kvm", strings.ToLower(vm.Name), purpose)
	} else {
		return fmt.Sprintf("%s-%s--%s-%s%d", strings.ToLower(vm.CreateUser), "kvm", strings.ToLower(vm.Name), purpose, number+1)
	}
}

// uuid作为kubevirt vm的名称，vnc访问时放在地址栏比较合适
func (vm *VirtualServer) KubevirtVmName() string {
	return vm.State.InstanceId
}

// 不用使用uuid作为service的名称
func (vm *VirtualServer) SshServiceName() string {
	return fmt.Sprintf("%s-%s", strings.ToLower(vm.CreateUser), strings.ToLower(vm.Name))
}

type StorageEntity struct {
	PvcName             string            `json:"pvcName"`
	Capacity            int               `json:"capacity"`
	StorageClass        string            `json:"storageClass"`
	IsNew               bool              `json:"isNew"       binding:"required"`
	ReleaseWithInstance bool              `json:"releaseWithInstance"       binding:"required"` // 默认为false
	CloneSource         string            `json:"-"`
	Labels              map[string]string `json:"-"` // 创建Pvc时打的标签，与数据创建模型和返回模型无关，无需json解析
}

type CloudInit struct {
	User      string   `json:"user"`
	PassWorld string   `json:"pwd"`
	SshKey    []string `json:"sshkeys"`
}

func (c *CloudInit) isValid() error {
	if (c.User == "" && c.PassWorld != "") || (c.PassWorld != "" && c.User == "") {
		return fmt.Errorf("invalid user init")
	}
	if c.User == "" && len(c.SshKey) == 0 {
		return fmt.Errorf("invalid user init")
	}
	if len(c.SshKey) > 10 {
		return fmt.Errorf("sshkey size should less than 10")
	}
	return nil
}

type State struct {
	InstanceId   string        `json:"instanceId"       binding:""`
	Status       string        `json:"status"       binding:""`
	Vnc          string        `json:"vnc"       binding:""`
	DomainName   string        `json:"domainName"       binding:""`
	SshPort      int32         `json:"sshPort"       binding:""`
	FailedStatus *FailedStatus `json:"failedStatus,omitempty"` // 出错时返回
}

type FailedStatus struct {
	Reason  string `json:"reason"       binding:""`
	Message string `json:"message"       binding:""`
}

func vnc(instanceId string) string {
	return fmt.Sprintf("/vnc/v1/vnc_lite.html?path=vnc/v1/k8s/apis/subresources.kubevirt.io/v1alpha3/namespaces/virtualserver/virtualmachineinstances/%s/vnc", instanceId)
}

func getValidSshPort(min, max int32, vms []repo.VirtualServer) int32 {
	usedPorts := map[int32]interface{}{}
	for _, v := range vms {
		usedPorts[v.SshPort] = nil
	}
	validSshPorts := []int32{}
	for i := min; i < max; i++ {
		if _, ok := usedPorts[i]; !ok {
			validSshPorts = append(validSshPorts, i)
		}
	}
	if len(validSshPorts) == 0 {
		return -1
	}
	index := rand.Intn(len(validSshPorts))
	return validSshPorts[index]
}

func (vm *VirtualServer) Init(createUser string, kubevirt kubecli.KubevirtClient,
	kubeclient clientset.Interface, repo pkg.VirtualServerRepo) error {
	// 0. 基础元数据初始化
	vm.CreateTime = utils.TimeNow()
	vm.CreateUser = createUser
	vm.StartRun = true // 创建完虚拟机是否启动，默认为true
	// 检查虚拟机名称是否冲突
	vms, err := repo.List(createUser)
	if err != nil {
		return fmt.Errorf("check vm err, %v", err)
	}
	for _, v := range vms {
		if v.Name == vm.Name {
			return fmt.Errorf("vm name %s already exist", vm.Name)
		}
	}
	// ssh映射端口 SSHPORT_MIN 预留给vnc服务
	sshPort := getValidSshPort(SSHPORT_MIN+1, SSHPORT_MAX, vms)
	if sshPort <= 0 {
		return fmt.Errorf("can not find valid service ssh port")
	}
	if vm.State.InstanceId == "" {
		// osm场景InstanceId由kvm自动生成
		vm.State.InstanceId = uuid.New().String()
	}
	vm.State.Status = STATE_PREPAREING
	vm.State.SshPort = sshPort
	vm.Labels = map[string]string{
		LabelGCPCreateUserKey: vm.CreateUser,
		// LabelGCPBindInstanceIdKey:   vm.State.InstanceId,
		LabelGCPInstanceIdKey:       vm.State.InstanceId,
		LabelGCPResourceTypeKey:     LabelGCPResourceTypeValue,
		LabelGCPBindInstanceNameKey: vm.Name,
	}
	// 初始化资源配额
	// vm.InitResource()

	// 初始化存储卷
	vm.Storage.Init(vm)
	// 检查存储卷
	// 获取全部存储卷
	userVolumeMap := map[string]corev1.PersistentVolumeClaim{}
	volumeList, err := kubeclient.CoreV1().PersistentVolumeClaims(vm.Namespace()).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("check volume err, %v", err)
	}
	for _, v := range volumeList.Items {
		userVolumeMap[v.Name] = v
	}
	// 检查系统存储卷是否冲突
	pvc, ok := userVolumeMap[vm.Storage.SystemStorage.PvcName]
	if vm.Storage.SystemStorage.IsNew && ok {
		return fmt.Errorf("volume %s already exist", vm.Storage.SystemStorage.PvcName)
	} else if !vm.Storage.SystemStorage.IsNew {
		if !ok {
			// 卷不存在
			return fmt.Errorf("volume %s not found", vm.Storage.SystemStorage.PvcName)
		}
		// 卷正在被使用
		instance, ok := pvc.Labels[LabelGCPBindInstanceIdKey]
		if ok {
			return fmt.Errorf("volume %s already in %s use", vm.Storage.SystemStorage.PvcName, instance)
		}
		// 卷的状态没有bound成功
		if pvc.Status.Phase != corev1.ClaimBound {
			return fmt.Errorf("pvc %s status %s not bound", pvc.Name, pvc.Status.Phase)
		}
	}

	// 检查数据存储卷是否冲突
	dataStorage := []StorageEntity{}
	for _, v := range vm.Storage.DataStorage {
		pvc, ok := userVolumeMap[v.PvcName]
		if v.IsNew && ok {
			return fmt.Errorf("volume %s already exist", v.PvcName)
		} else if !v.IsNew {
			if !ok {
				return fmt.Errorf("volume %s not found", v.PvcName)
			}
			// 卷正在被使用
			instance, ok := pvc.Labels[LabelGCPBindInstanceIdKey]
			if ok {
				return fmt.Errorf("volume %s already in %s use", v.PvcName, instance)
			}
		}
		dataStorage = append(dataStorage, v)
	}
	vm.Storage.DataStorage = dataStorage
	return nil
}

func (disk *StorageEntity) ConstructPvc(virtualServer *VirtualServer) *corev1.PersistentVolumeClaim {
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      disk.PvcName,
			Namespace: virtualServer.Namespace(),
			Labels:    disk.Labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &disk.StorageClass,
			// 后续文件存储时，考虑ReadOnlyMany
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(fmt.Sprintf("%vGi", disk.Capacity)),
				},
			},
		},
	}
	if disk.CloneSource != "" {
		// DataSource不为空时，走clone
		cloneSource := &corev1.TypedLocalObjectReference{
			Kind: "PersistentVolumeClaim", // 当前clone磁盘只支持pvc类型。后续考虑其他类型
			Name: disk.CloneSource,
		}
		pvc.Spec.DataSource = cloneSource
	}
	return &pvc
}

func (vm *VirtualServer) ConstructSshPortService() *corev1.Service {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vm.SshServiceName(),
			Namespace: vm.Namespace(),
			Labels:    vm.Labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: vm.Labels,
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Port:       vm.State.SshPort,
					TargetPort: intstr.FromInt(DEFAULT_SSH_PORT),
					NodePort:   vm.State.SshPort,
				},
			},
		},
	}
	return &svc
}

func (vm *VirtualServer) InitResource(repo pkg.VirtualServerRepo, productConfig []productConfig, products map[string]interface{}) (err error) {
	var productResource *ProductResource
	// productConfig 不为空，则为bsm场景
	if productConfig != nil {
		productResource, err = vm.getBSMProductResource(productConfig)
		if err != nil {
			return err
		}
	} else {
		// osm场景
		if localOSMProduct == nil {
			localOSMProduct, err = TransformLocalProductsToOSMProduct(products)
			if err != nil {
				return err
			}
		}
		productResource, err = vm.getOSMProductResource(localOSMProduct, repo)
		if err != nil {
			return err
		}
	}
	// 获取配置产品总量
	kubevirtResource := KubevirtResource{}
	// 配置GPU资源
	// productResource不为空，无需再做判空处理
	gpuNum := productResource.GpuNum
	for gpuNum > 0 && productResource.GpuK8sResource != "" {
		kubevirtResource.Gpus = append(kubevirtResource.Gpus, kubevirtv1.GPU{
			Name:       uuid.New().String(),
			DeviceName: productResource.GpuK8sResource,
		})
		gpuNum--
	}
	// 配置cpu和memory
	cpu, err := resource.ParseQuantity(productResource.Cpu)
	if err != nil {
		return fmt.Errorf("invalid cpu request %s", productResource.Cpu)
	}
	mempry, err := resource.ParseQuantity(productResource.Mem)
	if err != nil {
		return fmt.Errorf("invalid cpu request %s", productResource.Mem)
	}
	kubevirtResource.ResourceList = corev1.ResourceList{
		corev1.ResourceCPU:    cpu,
		corev1.ResourceMemory: mempry,
	}
	vm.KubevirtResource = kubevirtResource
	return nil
}

// todo 将bsm productConfig转化成ProductResource
func (vm *VirtualServer) getBSMProductResource(productConfig []productConfig) (*ProductResource, error) {
	pResource := ProductResource{}
	pConfigMap := map[string]string{}
	for _, pConfig := range productConfig {
		pConfigMap[pConfig.Key] = pConfig.Value
	}
	// 配置GPU
	if value, ok := pConfigMap[ProductK8sResourceKey]; ok {
		pResource.GpuK8sResource = value
	}
	if value, ok := pConfigMap[ProductGpuNumKey]; ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%s transform to int failed", value)
		}
		pResource.GpuNum = num
	}
	// 配置cpu
	if value, ok := pConfigMap[ProductCpuRequestKey]; ok {
		pResource.Cpu = value
	}
	// 配置memory
	if value, ok := pConfigMap[ProductMemRequestKey]; ok {
		pResource.Mem = value
	}
	return &pResource, nil
}

const (
	ProductGpuTypeKey     = "gpu_type"
	ProductK8sResourceKey = "gpu_k8s_resource"
	ProductGpuNumKey      = "gpu_num"
	ProductMemRequestKey  = "mem"
	ProductCpuRequestKey  = "cpu"
)

func (vm *VirtualServer) getOSMProductResource(osmProduct OSMProduct, repo pkg.VirtualServerRepo) (*ProductResource, error) {
	// 查询全量虚拟机列表
	vms, err := repo.List("")
	if err != nil {
		log.Infof("list vm from repo err, %v", err)
		return nil, err
	}
	// 查询各产品的使用数
	usedProductMount := 0
	for _, v := range vms {
		if v.Product == vm.ProductCode {
			usedProductMount++
		}
	}
	// 剩余产品库存量检查
	product := osmProduct[vm.ProductCode]
	if product.Num <= usedProductMount {
		return nil, fmt.Errorf("product %s number %d used %d", vm.ProductCode, product.Num, usedProductMount)
	}
	return product.ProductResource, nil
}

var localOSMProduct OSMProduct

type OSMProduct map[string]LocalProduct

type ProductResource struct {
	GpuK8sResource string
	GpuNum         int
	Mem            string
	Cpu            string
}

type LocalProduct struct {
	Num int
	*ProductResource
}

// func getImageOriginStorage(imageName string) string {
// 	return IMAGE_STORAGE_PRE + imageName
// }

func (vm *VirtualServer) CheckParam() error {
	// 校验虚拟机名称
	if !isValidName(vm.Name) {
		return fmt.Errorf("invalid virtualserver name %s, virtualserver name must match regex '^[a-z]([-a-z0-9]*[a-z0-9])?$'", vm.Name)
	}
	// 校验image
	if vm.Image == "" {
		return fmt.Errorf("invalid image %s", vm.Image)
	}
	// 校验产品（bsm场景产品code也是需要的，返回虚拟机详情后，前端需要使用产品code查询产品配置）
	if vm.ProductCode == "" {
		return fmt.Errorf("invalid product %s", vm.ProductCode)
	}
	// 校验初始化用户名密码
	if err := vm.CloudInit.isValid(); err != nil {
		return err
	}
	// 存储校验
	// if vm.Storage.SystemStorage.IsNew {
	// 	if vm.Storage.SystemStorage.Capacity < 10 || vm.Storage.SystemStorage.Capacity > 100 {
	// 		return fmt.Errorf("invalid disk %s Capacity %d, disk Capacity must match [10, 100]", vm.Storage.SystemStorage.PvcName, vm.Storage.SystemStorage.Capacity)
	// 	}
	// }
	// 检查数据盘
	if len(vm.Storage.DataStorage) > 16 {
		return fmt.Errorf("invalid disk numbers, disk numbers should less than 16")
	}
	// for _, disk := range vm.Storage.DataStorage {
	// 	if disk.IsNew {
	// 		if disk.Capacity < int(MinDiskCapacity) || disk.Capacity > int(MaxDiskCapacity) {
	// 			return fmt.Errorf("invalid disk %s Capacity %d, disk Capacity must match [%d, %d]", disk.PvcName, disk.Capacity, MinDiskCapacity, MaxDiskCapacity)
	// 		}
	// 	}
	// }
	return nil
}

func isValidName(name string) bool {
	pattern := "^[a-zA-Z][a-zA-Z0-9]{3,19}$"
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}
