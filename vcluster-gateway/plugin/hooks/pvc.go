package hooks

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	sdksync "github.com/loft-sh/vcluster-sdk/syncer/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewPvcHook(ctx *sdksync.RegisterContext) hook.ClientHook {
	return &pvcHook{
		rCtx: *ctx,
	}
}

type pvcHook struct {
	rCtx sdksync.RegisterContext
}

func (h *pvcHook) Name() string {
	return "pvc-hook"
}

func (h *pvcHook) Resource() client.Object {
	return &corev1.PersistentVolumeClaim{}
}

var _ hook.MutateCreatePhysical = &pvcHook{}

func (h *pvcHook) MutateCreatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	klog.Infof("🍽 Starting to mutate PVC: %s", obj.GetName())

	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		klog.Errorf("Expected a PersistentVolumeClaim object but got: %T", obj)
		return nil, fmt.Errorf("object %v is not a PersistentVolumeClaim", obj)
	}

	// 检查 PVC 是否有指定的 label
	var setOwnerReference bool
	for _, v := range pvc.Labels {
		if v == "alayaNeW" {
			setOwnerReference = true
			break
		}
	}

	if !setOwnerReference {
		klog.Infof("PVC %s/%s does not have required label, skipping", pvc.Namespace, pvc.Name)
		return pvc, nil
	}

	vclusterID, exists := pvc.Labels["vcluster.loft.sh/managed-by"]
	if !exists || vclusterID == "" {
		klog.Infof("PVC %s/%s does not have required label or label value is empty, skipping", pvc.Namespace, pvc.Name)
		return pvc, nil
	}

	pvName := pvc.Spec.VolumeName
	if pvName == "" {
		klog.Infof("PVC %s/%s does not have a volumeName yet, skipping PV owner reference setting", pvc.Namespace, pvc.Name)
		return pvc, nil
	}

	// 获取物理集群中的 client
	clientSet := h.rCtx.PhysicalManager.GetClient()

	// 首先获取物理集群中的命名空间
	physicalNs := &corev1.Namespace{}
	nsName := fmt.Sprintf("vcluster-%s", vclusterID)
	err := clientSet.Get(ctx, client.ObjectKey{Name: nsName}, physicalNs)
	if err != nil {
		klog.Errorf("Failed to get namespace %s: %v", nsName, err)
		return nil, err
	}

	// 获取物理集群中的 PV
	pv := &corev1.PersistentVolume{}
	err = clientSet.Get(ctx, client.ObjectKey{Name: pvName}, pv)
	if err != nil {
		klog.Errorf("Failed to retrieve PV %s: %v", pvName, err)
		return nil, err
	}

	// 创建 patch 操作
	patch := client.MergeFrom(pv.DeepCopy())

	// 设置 OwnerReference
	trueVar := true
	pv.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         "v1",
			Kind:               "Namespace",
			Name:               physicalNs.Name,
			UID:                physicalNs.UID,
			BlockOwnerDeletion: &trueVar,
			Controller:         &trueVar,
		},
	}

	// 提交变更
	if err := clientSet.Patch(ctx, pv, patch); err != nil {
		klog.Errorf("Failed to patch PV %s with owner reference: %v", pvName, err)
		return nil, err
	}

	klog.Infof("Successfully set owner reference for PV %s to namespace %s", pvName, nsName)

	return pvc, nil
}
