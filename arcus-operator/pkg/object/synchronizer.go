package object

import (
	corev1 "k8s.io/api/core/v1"
)

//==============================================================================
// Public Function
//==============================================================================
func SynchronizeService(dst *corev1.Service, src *corev1.Service) {
	dst.Spec.Ports = src.Spec.Ports
	dst.Spec.Type = src.Spec.Type
}

func SynchronizeConfigMap(dst *corev1.ConfigMap, src *corev1.ConfigMap) {
	dst.Data = src.Data
	dst.BinaryData = src.BinaryData
}
