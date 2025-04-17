package hooks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/loft-sh/vcluster-sdk/hook"
	sdksync "github.com/loft-sh/vcluster-sdk/syncer/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewPodHook(ctx *sdksync.RegisterContext) hook.ClientHook {
	return &podHook{
		*ctx,
	}
}

type podHook struct {
	rCtx sdksync.RegisterContext
}

func (p *podHook) Name() string {
	return "pod-hook"
}

func (p *podHook) Resource() client.Object {
	return &corev1.Pod{}
}

var _ hook.MutateCreatePhysical = &podHook{}

var _ hook.MutateUpdatePhysical = &podHook{}

// MutateCreatePhysical 用于创建 Pod 时的变更
func (p *podHook) MutateCreatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		klog.Errorf("Expected a Pod object but got: %T", obj)
		return nil, fmt.Errorf("object %v is not a pod", obj)
	}

	klog.Infof("Starting to mutate pod on create: %s in namespace: %s", pod.Name, pod.Namespace)

	pod = mutatePodForCreate(pod)

	return pod, nil
}

// MutateUpdatePhysical 用于更新 Pod 时的变更
func (p *podHook) MutateUpdatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		klog.Errorf("Expected a Pod object but got: %T", obj)
		return nil, fmt.Errorf("object %v is not a pod", obj)
	}

	klog.Infof("Starting to mutate pod on update: %s in namespace: %s", pod.Name, pod.Namespace)

	pod = mutatePodForUpdate(pod)

	return pod, nil
}

// mutatePod 执行共同的变更逻辑，包括注入环境变量和处理标签
func mutatePodForCreate(pod *corev1.Pod) *corev1.Pod {
	handleLabels(pod)
	injectEnv(pod)

	return pod
}

func mutatePodForUpdate(pod *corev1.Pod) *corev1.Pod {
	handleLabels(pod)

	return pod
}

// handleLabels 处理 Pod 的 Labels，包括自定义标签、注解转换等
func handleLabels(pod *corev1.Pod) {
	// 检查并初始化 Labels
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}

	pod.Labels["gcp.com/resource-type"] = "vcluster"
	pod.Labels["created-by-plugin"] = "pod-hook"

	// 如果不是手动初始化的集群，那么需要添加 instance-id 用于计量计费
	if os.Getenv("INITIALIZE") != "true" {
		instanceID := os.Getenv("instance-id")
		pod.Labels["dc.com/tenant.instance-id"] = instanceID
	}

	// 添加 dc.com/osm-vcluster-id，dc.com/osm-pod-name，dc.com/osm-pod-namespace 标签
	// 主要作用是能在宿主机通过标签匹配来获取 Pod 的信息
	managedBy := pod.Labels["vcluster.loft.sh/managed-by"] // vcluster id
	namespace := pod.Labels["vcluster.loft.sh/namespace"]  // vcluster 内部的 namespace
	var getPodNameFn = func(pod *corev1.Pod) string {
		var realPodName string

		// 从名称中提取实际 Pod 名称
		suffix := fmt.Sprintf("-x-%s-x-%s", namespace, managedBy)
		if strings.HasSuffix(pod.Name, suffix) {
			realPodName = strings.TrimSuffix(pod.Name, suffix)
		} else {
			klog.Warningf("Pod name does not match expected suffix pattern for pod %s", pod.Name)
		}
		return realPodName
	}

	if managedBy != "" && namespace != "" && getPodNameFn(pod) != "" {
		pod.Labels["dc.com/osm-vcluster-id"] = managedBy
		pod.Labels["dc.com/osm-pod-namespace"] = namespace
		pod.Labels["dc.com/osm-pod-name"] = getPodNameFn(pod)
	}

	// 将指定的 Annotations 转换为 Labels
	if annotations, exists := pod.Annotations["vcluster.loft.sh/labels"]; exists {
		lines := strings.Split(annotations, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				klog.Warningf("Malformed annotation line for pod %s: %s", pod.Name, line)
				continue
			}
			key := strings.Trim(parts[0], "\" ")
			value := strings.Trim(parts[1], "\" ")

			if strings.HasPrefix(line, "dc.com/") {
				pod.Labels[key] = value
			} else if strings.HasPrefix(line, "k8s-app") {
				if key == "k8s-app" && value == "kube-dns" {
					pod.Labels["release"] = strings.Split(pod.Namespace, "-")[1]
					pod.Labels["dc.com/tenant.source"] = "system"
				}
			}
		}
	}
}

// injectEnv 将环境变量注入到所有 containers 和 initContainers
// 包括：NODE_IP，VCLUSTER_ID，POD_NAME，POD_NAMESPACE
func injectEnv(pod *corev1.Pod) {
	// 公共环境变量
	commonEnvVars := []corev1.EnvVar{
		{
			Name: "NODE_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"},
			},
		},
		{
			Name:  "VCLUSTER_ID",
			Value: pod.Labels["dc.com/osm-vcluster-id"],
		},
		{
			Name:  "POD_NAME",
			Value: pod.Labels["dc.com/osm-pod-name"],
		},
		{
			Name:  "POD_NAMESPACE",
			Value: pod.Labels["dc.com/osm-pod-namespace"],
		},
	}

	// 将公共环境变量注入 Containers 和 InitContainers
	injectEnvToContainers := func(containers []corev1.Container) {
		for i := range containers {
			containers[i].Env = append(containers[i].Env, commonEnvVars...)
		}
	}

	injectEnvToContainers(pod.Spec.Containers)
	injectEnvToContainers(pod.Spec.InitContainers)
}
