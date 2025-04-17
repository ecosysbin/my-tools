package domain

import (
	"context"
	"fmt"
	"time"

	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg"
	v1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kubevirt_gateway/v1"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubevirtV1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"

	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type VirtualServerManager struct {
	// 当前实例所管理的虚拟服务器对象，回滚的时候从这里获取对象的详情进行回滚
	workqueue       workqueue.RateLimitingInterface
	InformerFactory informers.SharedInformerFactory
	// podInformer is the pods shared informer
	podInformer coreinformers.PodInformer
	// podLister can list/get pods from the shared informer's store
	podLister corelisters.PodLister
	// pvcInformer is the pvcs shared informer
	pvcInformer coreinformers.PersistentVolumeClaimInformer
	// pvcLister can list/get pvcs from the shared informer's store
	pvcLister       corelisters.PersistentVolumeClaimLister
	ResourcesSynced []cache.InformerSynced
	StopCh          chan struct{}
	Context         context.Context
	Repo            pkg.VirtualServerRepo
	Kubevirt        kubecli.KubevirtClient
	Kubecli         clientset.Interface
}

func NewVirtualServerManager(context context.Context, repo pkg.VirtualServerRepo, kubevirtCli kubecli.KubevirtClient, kubecli clientset.Interface) (VirtualServerManager, error) {
	factory := informers.NewSharedInformerFactoryWithOptions(kubecli, time.Second*10, informers.WithNamespace(v1.VIRTUALSERVER_NAMESPACE))
	podInformer := factory.Core().V1().Pods()
	pvcInformer := factory.Core().V1().PersistentVolumeClaims()
	manager := VirtualServerManager{
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kubevirt-gateway"),
		InformerFactory: factory,
		podInformer:     podInformer,
		podLister:       podInformer.Lister(),
		pvcInformer:     pvcInformer,
		pvcLister:       pvcInformer.Lister(),
		ResourcesSynced: []cache.InformerSynced{podInformer.Informer().HasSynced, pvcInformer.Informer().HasSynced},
		StopCh:          make(chan struct{}),
		Context:         context,
		Repo:            repo,
		Kubevirt:        kubevirtCli,
		Kubecli:         kubecli,
	}
	manager.podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    manager.AddEnventToWorkqueue,
		UpdateFunc: manager.UpdateEnventToWorkqueue,
		DeleteFunc: manager.AddEnventToWorkqueue,
	})
	return manager, nil
}

func (m *VirtualServerManager) StartInformer() {
	defer close(m.StopCh)
	go m.InformerFactory.Start(m.StopCh)
	<-m.StopCh
}

func (m *VirtualServerManager) StopInformer() {
	close(m.StopCh)
}

func (m *VirtualServerManager) Sync() {
	for m.processNextVMItem() {
	}
}

func (m *VirtualServerManager) processNextVMItem() bool {
	log.Infof("vm process item")
	obj, shutdown := m.workqueue.Get()

	if shutdown {
		return false
	}
	err := func(obj interface{}) error {
		defer m.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			m.workqueue.Forget(obj)
			log.Infof("expected string in workqueue but got %#v", obj)
			return nil
		}
		if err := m.sync(key); err != nil {
			if m.workqueue.NumRequeues(key) < 15 {
				// Put the item back on the workqueue to handle any transient errors.
				m.workqueue.AddRateLimited(key)
				return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
			}
		}
		return nil
	}(obj)
	if err != nil {
		log.Infof(err.Error())
		return true
	}
	return true
}

func (m *VirtualServerManager) sync(instanceId string) error {
	virtualServerRepo, err := m.Repo.GetVmById("", instanceId)
	if err != nil {
		log.Infof("repo get virtualServer %s, err", instanceId, err)
		return err
	}
	// kubevirt-gateway重启时，vm事件会触发，但是对象数据库状态为running，则不会做具体操作
	if virtualServerRepo.Status == v1.STATE_RUNNING {
		log.Infof("virtualServer %s status %s", virtualServerRepo.Name, v1.STATE_RUNNING)
		return nil
	}

	vmPod, err := m.getKubevirtVMPod(instanceId)
	if err != nil {
		return err
	}
	// 停止kubevirt-vm时，同步状态
	if virtualServerRepo.Status == v1.STATE_STOPPING && vmPod == nil {
		virtualServerRepo.Status = v1.STATE_STOPPED
		log.Infof("update vm %s status %s", virtualServerRepo.Name, virtualServerRepo.Status)
		if err := m.Repo.Update(&virtualServerRepo); err != nil {
			log.Infof("update virtualServer %s err, %v", virtualServerRepo.Name, err)
			return err
		}
		return nil
	}
	if vmPod == nil {
		log.Infof("vm %s pod not found", instanceId)
		return nil
	}
	// 创建kubevirt-vm时同步状态
	if virtualServerRepo.Status == v1.STATE_PREPAREING && vmPod != nil {
		virtualServerRepo.Status = v1.STATE_STARTING
		if err := m.Repo.Update(&virtualServerRepo); err != nil {
			log.Infof("update virtualserver %s err, %v", virtualServerRepo.Name, err)
			return err
		}
		return nil
	}
	// 删除vm时，完成全流程后，走主业务流程同步状态

	// 启动成功或重启成功kubevirt-vm同步状态
	if (virtualServerRepo.Status == v1.STATE_STARTING || virtualServerRepo.Status == v1.STATE_RESTARTING) && vmPod.Status.Phase == corev1.PodRunning {
		// 先检查vm状态，在重启时，pod状态不变化，vm状态会更新
		vm, err := m.Kubevirt.VirtualMachine(v1.VIRTUALSERVER_NAMESPACE).Get(context.Background(), instanceId, &metav1.GetOptions{})
		if err != nil {
			log.Infof("get virtualserver %s VirtualMach %s err, %v", virtualServerRepo.Name, instanceId, err)
			return err
		}
		if len(vm.Status.StateChangeRequests) != 0 {
			log.Infof("kubevirt vm %s is in %s state action", virtualServerRepo.Name, vm.Status.StateChangeRequests[0].Action)
			// 重新塞入队列里等待下次调度
			m.workqueue.Add(instanceId)
			return nil
		}
		log.Infof("kubevirt vm %s status %v", virtualServerRepo.Name, vm.Status.Ready)
		// 查询vmi是否卷挂载成功，没挂载成功则不处理，等待下次轮训
		// 查询虚拟服务器挂载点列表
		vmi, err := m.Kubevirt.VirtualMachineInstance(v1.VIRTUALSERVER_NAMESPACE).Get(context.Background(), instanceId, &metav1.GetOptions{})
		if err != nil {
			log.Infof("get virtualserver %s VirtualMachineInstance %s err, %v", virtualServerRepo.Name, instanceId, err)
			return err
		}
		volumePoints := map[string]string{}
		for _, volume := range vmi.Status.VolumeStatus {
			if volume.Target != "" {
				volumePoints[volume.Name] = v1.VOLUME_MOUNTPOINT_PRE + volume.Target
			}
		}
		if len(volumePoints) == 0 {
			log.Infof("virtualserver %s vmi already in starting", virtualServerRepo.Name)
			// 重新塞入队列里等待下次调度
			m.workqueue.Add(instanceId)
			return nil
		}
		// vmi启动完成,pod一定会有更新事件，接收这次事件即可
		// 更新存储卷label挂载点（因为挂载点格式/dev/vda，label的value不支持，所以挂载点放在annotion里）
		// 更新系统盘挂载点label，虚拟机中卷volume的名称即对应存储卷pvc的名称
		disks, err := m.virtualServerDisks(virtualServerRepo.Id)
		if err != nil {
			log.Infof("list disk err, %v", err)
			return err
		}
		for _, disk := range disks {
			mountPoint := volumePoints[disk.Name]
			disk.Annotations[v1.LabelGCPMountPointKey] = mountPoint
			err = updateVolume(context.Background(), disk, m.Kubecli)
			if err != nil {
				log.Infof("update virtualServer %s volume %s mountPoint err, %v", virtualServerRepo.Name, disk.Name, err)
			}
		}

		preStatus := virtualServerRepo.Status
		virtualServerRepo.Status = v1.STATE_RUNNING
		virtualServerRepo.StartedTime = utils.TimeNow()
		log.Infof("update vm %s from %s to %s status", virtualServerRepo.Name, preStatus, virtualServerRepo.Status)
		if err := m.Repo.Update(&virtualServerRepo); err != nil {
			log.Infof("update virtualServer %s err, %v", virtualServerRepo.Name, err)
			return err
		}
		return nil
	}
	// 虚拟机运行异常状态同步
	if vmPod != nil && vmPod.Status.Phase == corev1.PodFailed && vmPod.Status.Phase != corev1.PodFailed {
		virtualServerRepo.Status = v1.STATE_ABNORMITY
		virtualServerRepo.Reason = vmPod.Status.Conditions[0].Reason
		virtualServerRepo.Message = vmPod.Status.Conditions[0].Message
		log.Infof("update vm %s status %s", virtualServerRepo.Name, virtualServerRepo.Status)
		if err := m.Repo.Update(&virtualServerRepo); err != nil {
			log.Infof("update vm %s err, %v", vmPod.GetName(), err)
			return err
		}
	}
	return nil
}

func InitProducts(repo pkg.VirtualServerRepo) (map[string]int, error) {
	vms, err := repo.List("")
	if err != nil {
		log.Infof("list vm from repo err, %v", err)
		return nil, err
	}
	// 查询各产品的使用数
	usedProducts := map[string]int{}
	for _, vm := range vms {
		// 已删除的虚拟机产品会被释放
		usedProducts[vm.Product] = usedProducts[vm.Product] + 1
	}
	return usedProducts, nil
}

type VirtulServerWorkFlow struct {
	Metadata v1.VirtualServer
	Works    []VirtulServerFlow
}

type VirtulServerFlow struct {
	WorkName string
	Work     WorkFlow
	RollBack WorkFlow
}

type WorkFlow func(v1.VirtualServer) error

const (
	WORK_CREATE_ESTORAGE       = "create storage"
	WORK_CHECK_STORAGESTATUS   = "check storage bound status"
	WORK_CREATE_KUBEVIRTVM     = "create kubevirt vm"
	WORK_CREATE_SSHPORTSERVICE = "create virtualserver ssh port service"
	WORK_DELETE_SSHPORTSERVICE = "delete ssh port service"
	WORK_DELETE_KUBEVIRTVM     = "delete kubevirt vm"
	WORK_DELETE_STORAGE        = "delete storage"
	WORK_UPDATE_STATUS         = "update status"
)

func (m *VirtualServerManager) AddVirtualServer(virtualserver v1.VirtualServer) error {
	// 0.将用户数据初始化到数据库，准备创建虚拟机
	if err := m.Repo.Store(v1.TransformVirtualServerToRepo(virtualserver)); err != nil {
		return err
	}
	virtulServerCreateWorkFlow := VirtulServerWorkFlow{
		Metadata: virtualserver,
		Works: []VirtulServerFlow{
			{
				WorkName: WORK_CREATE_ESTORAGE,
				Work:     m.CreateStorage, // 1. 创建存储,
				RollBack: m.RecoverStorage,
			},
			{
				WorkName: WORK_CHECK_STORAGESTATUS,
				Work:     workforTimeout(WORK_CHECK_STORAGESTATUS, 3, 60, m.CheckStorageBoundStatus), // 2. 检查存储卷是否绑定成功, 一直到超时失败
			},
			{
				WorkName: WORK_CREATE_KUBEVIRTVM,
				Work:     m.CreateKubevirtVm, // 3. kubevirt创建虚拟服务器
				RollBack: m.DeleteKubevirtVm,
			},
			{
				WorkName: WORK_CREATE_SSHPORTSERVICE,
				Work:     m.CreateSshPortService, // 4. 创建虚拟机服务端口映射
				RollBack: m.DeleteSshPortService,
			},
		},
	}
	go m.StartVirtualserverWorkFlow(virtulServerCreateWorkFlow)
	return nil
}

func (m *VirtualServerManager) RestartVirtualServer(virtualServer v1.VirtualServer) error {
	virtualServer.State.Status = v1.STATE_RESTARTING
	virtualServer.StartedTime = nil
	virtualServerRepo := v1.TransformVirtualServerToRepo(virtualServer)
	if err := m.Repo.Update(&virtualServerRepo); err != nil {
		return fmt.Errorf("update vm %s status err, %v", virtualServer.Name, err)
	}
	log.Infof("restart vm %s begin", virtualServer.Name)
	if err := m.RestartKubevirtVm(virtualServer); err != nil {
		log.Infof("stop vm %s err, %v", virtualServer.Name, err)
	}
	return nil
}

func (m *VirtualServerManager) StopVirtualServer(virtualServer v1.VirtualServer) error {
	virtualServer.State.Status = v1.STATE_STOPPING
	virtualServer.StartedTime = nil
	virtualServerRepo := v1.TransformVirtualServerToRepo(virtualServer)
	if err := m.Repo.Update(&virtualServerRepo); err != nil {
		return fmt.Errorf("update vm %s status err, %v", virtualServer.Name, err)
	}
	log.Infof("stop vm %s begin", virtualServer.Name)
	if err := m.StopKubevirtVm(virtualServer); err != nil {
		log.Infof("stop vm %s err, %v", virtualServer.Name, err)
		return err
	}
	return nil
}

func (m *VirtualServerManager) StartVirtualServer(virtualServer v1.VirtualServer) error {
	virtualServer.State.Status = v1.STATE_STARTING
	virtualServerRepo := v1.TransformVirtualServerToRepo(virtualServer)
	if err := m.Repo.Update(&virtualServerRepo); err != nil {
		return fmt.Errorf("update vm %s status err, %v", virtualServer.Name, err)
	}
	log.Infof("start vm %s begin", virtualServer.Name)
	if err := m.StartKubevirtVm(virtualServer); err != nil {
		log.Infof("start vm %s err, %v", virtualServer.Name, err)
		return err
	}
	return nil
}

// 修改数据库模型为已停止（数据库状态字段跟着watch delete事件走会被修改为已停止）
// 考虑发起删除后，虚拟机在已释放前应该还是在未释放的列表中。则修改模型为已释放放在watch delete事件处理中。发生删除事件并且状态为删除中，则修改状态为已停止，且已释放。
func (m *VirtualServerManager) DeleteVirtualServer(virtualServer v1.VirtualServer) error {
	// 0. 数据库状态记录为删除中（考虑异常流程处理一启动就把数据加载到内存，后续新增，删除都走内存查询。当前假设没有重启）
	virtualServer.State.Status = v1.STATE_DELETING
	virtualServer.DeleteTime = utils.TimeNow()
	virtualServerRepo := v1.TransformVirtualServerToRepo(virtualServer)
	if err := m.Repo.Update(&virtualServerRepo); err != nil {
		return fmt.Errorf("update vm %s status err, %v", virtualServer.Name, err)
	}

	virtulServerDeleteWorkFlow := VirtulServerWorkFlow{
		Metadata: virtualServer,
		Works: []VirtulServerFlow{
			{
				WorkName: WORK_DELETE_SSHPORTSERVICE,
				Work:     m.DeleteSshPortService, // 1. 删除虚拟机服务端口映射 (删除顺序和创建顺序相反，先删除端口映射即虚拟机会立即停止服务)
			},
			{
				WorkName: WORK_DELETE_KUBEVIRTVM,
				Work:     m.DeleteKubevirtVm, // 2. 删除kubevirt虚拟服务器
			},
			{
				WorkName: WORK_DELETE_STORAGE,
				Work:     m.DeleteStorage, // 3. 删除存储
			},
			{
				WorkName: WORK_UPDATE_STATUS,
				Work:     m.DeleteSuccess, // 4. 更新状态为已删除
			},
		},
	}

	go m.StartVirtualserverWorkFlow(virtulServerDeleteWorkFlow)
	return nil
}

func (m *VirtualServerManager) DeleteSuccess(virtualServer v1.VirtualServer) error {
	virtualServer.State.Status = v1.STATE_STOPPED
	virtualServer.Deleted = v1.DELETED_YES_INT
	virtualServerRepo := v1.TransformVirtualServerToRepo(virtualServer)
	if err := m.Repo.Update(&virtualServerRepo); err != nil {
		return fmt.Errorf("update vm %s status err, %v", virtualServer.Name, err)
	}
	return nil
}

func (m *VirtualServerManager) StartVirtualserverWorkFlow(workFlow VirtulServerWorkFlow) {
	rollBackWorks := []WorkFlow{}
	for _, work := range workFlow.Works {
		if work.RollBack != nil {
			rollBackWorks = append(rollBackWorks, work.RollBack)
		}
		virtualServer := workFlow.Metadata
		if err := work.Work(virtualServer); err != nil {
			log.Infof("virtualServer %s process %s do err,%v", workFlow.Metadata.Name, work.WorkName, err)
			// 更新失败状态
			failedStatus := &v1.FailedStatus{
				Reason:  fmt.Sprintf("%s err", work.WorkName),
				Message: err.Error(),
			}
			virtualServer.State.Status = v1.STATE_ABNORMITY
			virtualServer.State.FailedStatus = failedStatus
			if err := m.UpdateStatus(virtualServer); err != nil {
				fmt.Printf("update vm %s status err, %v", virtualServer.Name, err)
			}
			for _, rollBackWork := range rollBackWorks {
				if err := rollBackWork(workFlow.Metadata); err != nil {
					fmt.Printf("rollback %s process %s do err,%v", workFlow.Metadata.Name, work.WorkName, err)
				}
			}
			// 回滚完成即任务退出
			return
		}
	}
}

func (m *VirtualServerManager) CheckStorageBoundStatus(virtualserver v1.VirtualServer) error {
	pvcList, err := m.virtualServerDisks(virtualserver.State.InstanceId)
	if err != nil {
		return fmt.Errorf("kubecli list pvc %s err, %v", virtualserver.Name, err)
	}
	vmPvcs := map[string]*corev1.PersistentVolumeClaim{}
	for _, pvc := range pvcList {
		vmPvcs[pvc.Name] = pvc
	}
	// 检查所有新创建的卷都已经bound成功
	systemDiskPvc, ok := vmPvcs[virtualserver.Storage.SystemStorage.PvcName]
	if !ok {
		return fmt.Errorf("vm %s system storage pvc %s not found", virtualserver.Name, virtualserver.Storage.SystemStorage.PvcName)
	}
	if systemDiskPvc.Status.Phase != corev1.ClaimBound {
		return fmt.Errorf("pvc %s status %s not bound", systemDiskPvc.Name, systemDiskPvc.Status.Phase)
	}
	ddisk := []v1.StorageEntity{}
	for _, disk := range virtualserver.Storage.DataStorage {
		if disk.IsNew {
			diskPvc, ok := vmPvcs[disk.PvcName]
			if !ok {
				return fmt.Errorf("vm %s disk pvc %s not found", virtualserver.Name, disk.PvcName)
			}
			if systemDiskPvc.Status.Phase != corev1.ClaimBound {
				return fmt.Errorf("pvc %s status %s not bound", diskPvc.Name, diskPvc.Status.Phase)
			}
		}
		ddisk = append(ddisk, disk)
	}
	virtualserver.Storage.DataStorage = ddisk
	return nil
}

func workforTimeout(workName string, period int64, timeOut int64, work WorkFlow) WorkFlow {
	return func(vm v1.VirtualServer) error {
		for timeOut > 0 {
			if err := work(vm); err != nil {
				log.Infof("%s err, %v, continue work", workName, err)
				time.Sleep(time.Duration(period) * time.Second)
				timeOut--
				continue
			}
			return nil
		}
		return fmt.Errorf("%s until timeout", workName)
	}
}

func (m *VirtualServerManager) UpdateStatus(virtualServer v1.VirtualServer) error {
	virtualserverRepo := v1.TransformVirtualServerToRepo(virtualServer)
	if err := m.Repo.Update(&virtualserverRepo); err != nil {
		return err
	}
	return nil
}

func (m *VirtualServerManager) CreateStorage(virtualserver v1.VirtualServer) error {
	// 系统盘创建直接clone系统中的镜像盘
	systemDisk := virtualserver.Storage.SystemStorage
	if systemDisk.IsNew { // IsNew为true时创建卷，为false时更新标签
		if err := createVolume(context.Background(), &virtualserver, &systemDisk, m.Kubecli); err != nil {
			log.Infof("create vm %s system volume err, %v", virtualserver.Name, err)
			return err
		}
	} else {
		// 更新存储卷标签（被实例绑定）
		pvc, err := m.pvcLister.PersistentVolumeClaims(v1.VIRTUALSERVER_NAMESPACE).Get(systemDisk.PvcName)
		if errors.IsNotFound(err) {
			return fmt.Errorf("vm %s system volume pvc %s not found in cacha", virtualserver.Name, systemDisk.PvcName)
		}
		if err != nil {
			return err
		}
		pvc.Labels = utils.MergeMaps(pvc.Labels, systemDisk.Labels)
		if err := updateVolume(context.Background(), pvc, m.Kubecli); err != nil {
			log.Infof("update vm %s data volume %s err, %v", virtualserver.Name, systemDisk.PvcName, err)
			return err
		}
	}
	// 创建数据盘存储卷
	for _, disk := range virtualserver.Storage.DataStorage {
		if !disk.IsNew {
			// 数据盘IsNew=false，即不需要新建, 更新标签
			pvc, err := m.pvcLister.PersistentVolumeClaims(v1.VIRTUALSERVER_NAMESPACE).Get(disk.PvcName)
			if errors.IsNotFound(err) {
				return fmt.Errorf("vm %s disk volume pvc %s not found in cacha", virtualserver.Name, disk.PvcName)
			}
			if err != nil {
				return err
			}
			pvc.Labels = utils.MergeMaps(pvc.Labels, disk.Labels)
			if err := updateVolume(context.Background(), pvc, m.Kubecli); err != nil {
				log.Infof("update vm %s data volume %s err, %v", virtualserver.Name, disk.PvcName, err)
				return err
			}
			continue
		}
		if err := createVolume(context.Background(), &virtualserver, &disk, m.Kubecli); err != nil {
			log.Infof("create vm %s data volume %s err, %v", virtualserver.Name, disk.PvcName, err)
			return err
		}
	}
	return nil
}

func (m *VirtualServerManager) RecoverStorage(virtualserver v1.VirtualServer) error {
	pvcs, err := m.virtualServerDisks(virtualserver.State.InstanceId)
	if err != nil {
		return err
	}
	for _, pvc := range pvcs {
		if pvc.Labels[v1.LabelGCPCreateByInstance] == virtualserver.State.InstanceId {
			// 随实例创建的卷，回滚时删除
			log.Infof("recover create vm %s delete volume %s", virtualserver.Name, pvc.Name)
			if err := deleteVolume(context.Background(), &virtualserver, pvc.Name, m.Kubecli); err != nil {
				log.Infof("create vm %s data volume %s err, %v", virtualserver.Name, pvc.Name, err)
				return err
			}
		} else {
			// 不随实例创建的卷，即卸载实例标签
			log.Infof("recover create vm %s uninstall instance label form volume %s", virtualserver.Name, pvc.Name)
			pvc.Labels = utils.ReduceMaps(pvc.Labels, []string{v1.LabelGCPBindInstanceIdKey, v1.LabelGCPBindInstanceNameKey})
			pvc.Annotations = utils.ReduceMaps(pvc.Annotations, []string{v1.LabelGCPMountPointKey})
			if err := updateVolume(context.Background(), pvc, m.Kubecli); err != nil {
				log.Infof("update vm %s data volume %s err, %v", virtualserver.Name, pvc.Name, err)
				return err
			}
		}
	}
	return nil
}

func (m *VirtualServerManager) DeleteStorage(virtualserver v1.VirtualServer) error {
	pvcs, err := m.virtualServerDisks(virtualserver.State.InstanceId)
	if err != nil {
		return err
	}
	for _, pvc := range pvcs {
		if pvc.Labels[v1.LabelGCPReleaseByInstance] != "" {
			// 随实例释放
			if err := deleteVolume(context.Background(), &virtualserver, pvc.Name, m.Kubecli); err != nil {
				log.Infof("create vm %s data volume %s err, %v", virtualserver.Name, pvc.Name, err)
				return err
			}
		} else {
			// 不随实例释放，即卸载实例标签
			pvc.Labels = utils.ReduceMaps(pvc.Labels, []string{v1.LabelGCPBindInstanceIdKey, v1.LabelGCPBindInstanceNameKey})
			pvc.Annotations = utils.ReduceMaps(pvc.Annotations, []string{v1.LabelGCPMountPointKey})
			if err := updateVolume(context.Background(), pvc, m.Kubecli); err != nil {
				log.Infof("update vm %s data volume %s err, %v", virtualserver.Name, pvc.Name, err)
				return err
			}
		}
	}
	return nil
}

func updateVolume(ctx context.Context, pvc *corev1.PersistentVolumeClaim, Kubecli clientset.Interface) error {
	_, err := Kubecli.
		CoreV1().
		PersistentVolumeClaims(pvc.Namespace).
		Update(ctx, pvc, metav1.UpdateOptions{})

	return err
}

func createVolume(ctx context.Context, virtualServer *v1.VirtualServer, disk *v1.StorageEntity, Kubecli clientset.Interface) error {
	_, err := Kubecli.
		CoreV1().
		PersistentVolumeClaims(virtualServer.Namespace()).
		Create(ctx, disk.ConstructPvc(virtualServer), metav1.CreateOptions{})
	return err
}

func deleteVolume(ctx context.Context, virtualServer *v1.VirtualServer, pvcName string, Kubecli clientset.Interface) error {
	err := Kubecli.
		CoreV1().
		PersistentVolumeClaims(virtualServer.Namespace()).
		Delete(ctx, pvcName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func (m *VirtualServerManager) CreateKubevirtVm(virtualserver v1.VirtualServer) error {
	_, err := m.Kubevirt.VirtualMachine(virtualserver.Namespace()).Create(context.Background(), virtualserver.ConstructKubevirtVm())
	return err
}

func (m *VirtualServerManager) DeleteKubevirtVm(virtualserver v1.VirtualServer) error {
	err := m.Kubevirt.VirtualMachine(virtualserver.Namespace()).Delete(context.Background(), virtualserver.KubevirtVmName(), &metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func (m *VirtualServerManager) StopKubevirtVm(virtualserver v1.VirtualServer) error {
	err := m.Kubevirt.VirtualMachine(virtualserver.Namespace()).Stop(context.Background(), virtualserver.KubevirtVmName(), &kubevirtV1.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *VirtualServerManager) StartKubevirtVm(virtualserver v1.VirtualServer) error {
	err := m.Kubevirt.VirtualMachine(virtualserver.Namespace()).Start(context.Background(), virtualserver.KubevirtVmName(), &kubevirtV1.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *VirtualServerManager) RestartKubevirtVm(virtualserver v1.VirtualServer) error {
	err := m.Kubevirt.VirtualMachine(virtualserver.Namespace()).Restart(context.Background(), virtualserver.KubevirtVmName(), &kubevirtV1.RestartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *VirtualServerManager) CreateSshPortService(virtualserver v1.VirtualServer) error {
	_, err := m.Kubecli.
		CoreV1().
		Services(virtualserver.Namespace()).
		Create(context.Background(), virtualserver.ConstructSshPortService(), metav1.CreateOptions{})

	return err
}

func (m *VirtualServerManager) DeleteSshPortService(virtualserver v1.VirtualServer) error {
	err := m.Kubecli.
		CoreV1().
		Services(virtualserver.Namespace()).
		Delete(context.Background(), virtualserver.SshServiceName(), metav1.DeleteOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func (m *VirtualServerManager) virtualServerDisks(instanceId string) ([]*corev1.PersistentVolumeClaim, error) {
	setSelector := labels.SelectorFromSet(labels.Set(map[string]string{
		v1.LabelGCPBindInstanceIdKey: instanceId}))

	pvcList, err := m.pvcLister.PersistentVolumeClaims(v1.VIRTUALSERVER_NAMESPACE).List(setSelector)
	if err != nil {
		return nil, err
	}
	instancePvcList := []*corev1.PersistentVolumeClaim{}
	for _, pvc := range pvcList {
		if pvc.Labels[v1.LabelGCPBindInstanceIdKey] == instanceId {
			instancePvcList = append(instancePvcList, pvc)
		}
	}
	return instancePvcList, nil
}

func (m *VirtualServerManager) AddEnventToWorkqueue(obj interface{}) {
	if obj == nil {
		return
	}
	mObj := obj.(metav1.Object)
	log.Infof("New Pod Added to Store: %s", mObj.GetName())
	labels := mObj.GetLabels()
	instanceId := labels[v1.LabelGCPInstanceIdKey]
	if instanceId == "" {
		log.Infof("vm %s instanceId not found", mObj.GetName())
		return
	}
	m.workqueue.Add(instanceId)
}

func (m *VirtualServerManager) UpdateEnventToWorkqueue(oldObj, newObj interface{}) {
	if newObj == nil || oldObj == nil {
		return
	}
	oObj := oldObj.(metav1.Object)
	nObj := newObj.(metav1.Object)
	if oObj.GetResourceVersion() == nObj.GetResourceVersion() {
		return
	}
	log.Infof("New Pod Updated to Store: %s", nObj.GetName())
	labels := nObj.GetLabels()
	instanceId := labels[v1.LabelGCPInstanceIdKey]
	if instanceId == "" {
		log.Infof("vm %s instanceId not found", nObj.GetName())
		return
	}
	m.workqueue.Add(instanceId)
}

func (m *VirtualServerManager) getKubevirtVMPod(instanceId string) (*corev1.Pod, error) {
	setSelector := labels.SelectorFromSet(labels.Set(map[string]string{
		v1.LabelGCPInstanceIdKey: instanceId}))

	pods, err := m.podLister.Pods(v1.VIRTUALSERVER_NAMESPACE).List(setSelector)
	if err != nil {
		return nil, err
	}
	if len(pods) > 0 {
		return pods[0], nil
	}
	return nil, nil
}
