package object

import (
	"strconv"

	jam2inv1 "github.com/jam2in/arcus-operator/pkg/apis/jam2in/v1"
	"github.com/jam2in/arcus-operator/pkg/global"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
			Name:      global.ObjectNameConfigMap,
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

func CreateZkHeadlessService(arcus *jam2inv1.Arcus) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      global.ObjectNameZkHeadlessService,
			Namespace: arcus.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(arcus, schema.GroupVersionKind{
					Group:   jam2inv1.SchemeGroupVersion.Group,
					Version: jam2inv1.SchemeGroupVersion.Version,
					Kind:    global.KindNameArcus,
				}),
			},
			Labels: global.GetZkLabel(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Ports:     arcus.GetZkServicePorts(),
			Selector:  global.GetZkLabel(),
		},
	}
}
