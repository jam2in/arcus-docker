package util

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

//==============================================================================
// Public Functions
//==============================================================================
func SynchronizeStatefulSet(dst *appsv1.StatefulSet, src *appsv1.StatefulSet) {
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Spec.Template = src.Spec.Template
	dst.Spec.UpdateStrategy = src.Spec.UpdateStrategy
}

func SynchronizeService(dst *corev1.Service, src *corev1.Service) {
	dst.Spec.Ports = src.Spec.Ports
	dst.Spec.Type = src.Spec.Type
}

func SynchronizeConfigMap(dst *corev1.ConfigMap, src *corev1.ConfigMap) {
	dst.Data = src.Data
	dst.BinaryData = src.BinaryData
}

func MakeFilePath(folder string, filename string) string {
	return strings.Join([]string{folder, filename}, "/")
}
