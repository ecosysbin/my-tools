package hooks

import (
	"os"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodHook_MutateCreatePhysical(t *testing.T) {
	// 设置环境变量
	os.Setenv("instance-id", "test-instance-id")
	defer os.Unsetenv("instance-id")

	tests := []struct {
		name           string
		pod            *corev1.Pod
		expectedLabels map[string]string
		expectedError  bool
	}{
		{
			name: "No annotations",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Labels:    map[string]string{},
				},
			},
			expectedLabels: map[string]string{
				"gcp.com/instance-id":   "test-instance-id",
				"gcp.com/resource-type": "vcluster",
				"created-by-plugin":     "pod-hook",
			},
			expectedError: false,
		},
		{
			name: "With annotations and dc.com labels",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Labels:    map[string]string{},
					Annotations: map[string]string{
						"vcluster.loft.sh/labels": `
							dc.com/foo="bar"
							dc.com/another="value"
							other.com/key="ignore"`,
					},
				},
			},
			expectedLabels: map[string]string{
				"gcp.com/instance-id":   "test-instance-id",
				"gcp.com/resource-type": "vcluster",
				"created-by-plugin":     "pod-hook",
				"dc.com/foo":            "bar",
				"dc.com/another":        "value",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//hook := NewPodHook()

			//pod, err := hook.(*podHook).MutateCreatePhysical(context.TODO(), tt.pod)
			//
			//if tt.expectedError {
			//	assert.Error(t, err)
			//} else {
			//	assert.NoError(t, err)
			//	assert.Equal(t, tt.expectedLabels, pod.GetLabels())
			//}
		})
	}
}
