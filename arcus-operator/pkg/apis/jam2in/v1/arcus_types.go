package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//==============================================================================
// Constant
//==============================================================================
const (
	DefaultZkReplicas          = 3
	DefaultZkMaxClientCnxns    = 100
	DefaultZkTickTime          = 2000
	DefaultZkInitLimit         = 10
	DefaultZkSyncLimit         = 5
	DefaultZkMinSessionTimeout = 4000
	DefaultZkMaxSessionTimeout = 200000
)

//==============================================================================
// ArcusSpec
//==============================================================================
type ArcusSpec struct {
	ZkReplicas          int32 `json:"zkReplicas"`
	ZkMaxClientCnxns    int32 `json:"zkMaxClientCnxns"`
	ZkTickTime          int32 `json:"zkTickTime"`
	ZkInitLimit         int32 `json:"zkInitLimit"`
	ZkSyncLimit         int32 `json:"zkSyncLimit"`
	ZkMinSessionTimeout int32 `json:"zkMinSessionTimeout"`
	ZkMaxSessionTimeout int32 `json:"zkMaxSessionTimeout"`
}

//==============================================================================
// ArcusStatus
//==============================================================================
type ArcusStatus struct {
}

//==============================================================================
// Arcus
//==============================================================================
func (arcus *Arcus) WithDefaults() (changed bool) {
	if arcus.Spec.ZkReplicas == 0 {
		changed = true
		arcus.Spec.ZkReplicas = DefaultZkReplicas
	}
	if arcus.Spec.ZkMaxClientCnxns == 0 {
		changed = true
		arcus.Spec.ZkMaxClientCnxns = DefaultZkMaxClientCnxns
	}
	if arcus.Spec.ZkTickTime == 0 {
		changed = true
		arcus.Spec.ZkTickTime = DefaultZkTickTime
	}
	if arcus.Spec.ZkInitLimit == 0 {
		changed = true
		arcus.Spec.ZkInitLimit = DefaultZkInitLimit
	}
	if arcus.Spec.ZkSyncLimit == 0 {
		changed = true
		arcus.Spec.ZkSyncLimit = DefaultZkSyncLimit
	}
	if arcus.Spec.ZkMinSessionTimeout == 0 {
		changed = true
		arcus.Spec.ZkMinSessionTimeout = DefaultZkMinSessionTimeout
	}
	if arcus.Spec.ZkMaxSessionTimeout == 0 {
		changed = true
		arcus.Spec.ZkMaxSessionTimeout = DefaultZkMaxSessionTimeout
	}
	return changed
}

type Arcus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArcusSpec   `json:"spec,omitempty"`
	Status ArcusStatus `json:"status,omitempty"`
}

//==============================================================================
// ArcusList
//==============================================================================
type ArcusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Arcus `json:"items"`
}

//==============================================================================
// Private Function
//==============================================================================
func init() {
	SchemeBuilder.Register(&Arcus{}, &ArcusList{})
}
