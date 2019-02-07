package object

import (
	corev1 "k8s.io/api/core/v1"
)

func SynchronizeConfigMap(dst *corev1.ConfigMap, src *corev1.ConfigMap) {
	dst.Data = src.Data
	dst.BinaryData = src.BinaryData
}
