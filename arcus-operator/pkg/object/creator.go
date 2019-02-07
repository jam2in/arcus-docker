package object

import (
	"strconv"

	jam2inv1 "github.com/jam2in/arcus-operator/pkg/apis/jam2in/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//==============================================================================
// Constant
//==============================================================================
const (
	ObjectNameConfigMap = "arcus-configmap"
)

//==============================================================================
// Public Function
//==============================================================================
func CreateConfigMap(arcus *jam2inv1.Arcus) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ObjectNameConfigMap,
			Namespace: arcus.Namespace,
		},
		Data: map[string]string{
			"maxClientCnxns":    strconv.Itoa(int(arcus.Spec.ZkMaxClientCnxns)),
			"tickTime":          strconv.Itoa(int(arcus.Spec.ZkTickTime)),
			"initLimit":         strconv.Itoa(int(arcus.Spec.ZkInitLimit)),
			"syncLimit":         strconv.Itoa(int(arcus.Spec.ZkSyncLimit)),
			"minSessionTimeout": strconv.Itoa(int(arcus.Spec.ZkMinSessionTimeout)),
			"maxSessionTimeout": strconv.Itoa(int(arcus.Spec.ZkMaxSessionTimeout)),
		},
	}
}
