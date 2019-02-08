package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//==============================================================================
// Constant
//==============================================================================
const (
	DefaultZkReplicas = 3

	DefaultZkMaxClientCnxns    = 100
	DefaultZkTickTime          = 2000
	DefaultZkInitLimit         = 10
	DefaultZkSyncLimit         = 5
	DefaultZkMinSessionTimeout = 4000
	DefaultZkMaxSessionTimeout = 200000

	DefaultZkClientPort         = 2181
	DefaultZkServerPort         = 2888
	DefaultZkLeaderElectionPort = 3888
)

//==============================================================================
// ArcusSpec
//==============================================================================
type ArcusSpec struct {
	Zookeeper ZookeeperConfig `json:"zookeeper,omitempty"`
}

//==============================================================================
// ArcusStatus
//==============================================================================
type ArcusStatus struct {
}

//==============================================================================
// Arcus
//==============================================================================
type Arcus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArcusSpec   `json:"spec,omitempty"`
	Status ArcusStatus `json:"status,omitempty"`
}

func (arcus *Arcus) WithDefaults() (changed bool) {
	return arcus.Spec.Zookeeper.withDefaults(arcus)
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
// ZookeeperConfig
//==============================================================================
type ZookeeperConfig struct {
	Replicas          int32 `json:"replicas"`
	MaxClientCnxns    int32 `json:"maxClientCnxns"`
	TickTime          int32 `json:"tickTime"`
	InitLimit         int32 `json:"initLimit"`
	SyncLimit         int32 `json:"syncLimit"`
	MinSessionTimeout int32 `json:"minSessionTimeout"`
	MaxSessionTimeout int32 `json:"maxSessionTimeout"`

	Ports ZookeeperPort `json:"ports,omitempty"`
}

func (c *ZookeeperConfig) withDefaults(arcus *Arcus) (changed bool) {
	if arcus.Spec.Zookeeper.Replicas == 0 {
		changed = true
		arcus.Spec.Zookeeper.Replicas = DefaultZkReplicas
	}
	if arcus.Spec.Zookeeper.MaxClientCnxns == 0 {
		changed = true
		arcus.Spec.Zookeeper.MaxClientCnxns = DefaultZkMaxClientCnxns
	}
	if arcus.Spec.Zookeeper.TickTime == 0 {
		changed = true
		arcus.Spec.Zookeeper.TickTime = DefaultZkTickTime
	}
	if arcus.Spec.Zookeeper.InitLimit == 0 {
		changed = true
		arcus.Spec.Zookeeper.InitLimit = DefaultZkInitLimit
	}
	if arcus.Spec.Zookeeper.SyncLimit == 0 {
		changed = true
		arcus.Spec.Zookeeper.SyncLimit = DefaultZkSyncLimit
	}
	if arcus.Spec.Zookeeper.MinSessionTimeout == 0 {
		changed = true
		arcus.Spec.Zookeeper.MinSessionTimeout = DefaultZkMinSessionTimeout
	}
	if arcus.Spec.Zookeeper.MaxSessionTimeout == 0 {
		changed = true
		arcus.Spec.Zookeeper.MaxSessionTimeout = DefaultZkMaxSessionTimeout
	}

	if arcus.Spec.Zookeeper.Ports.Client == 0 {
		changed = true
		arcus.Spec.Zookeeper.Ports.Client = DefaultZkClientPort
	}
	if arcus.Spec.Zookeeper.Ports.Server == 0 {
		changed = true
		arcus.Spec.Zookeeper.Ports.Server = DefaultZkServerPort
	}
	if arcus.Spec.Zookeeper.Ports.LeaderElection == 0 {
		changed = true
		arcus.Spec.Zookeeper.Ports.LeaderElection = DefaultZkLeaderElectionPort
	}

	return changed
}

//==============================================================================
// ZookeeperPort
//==============================================================================
type ZookeeperPort struct {
	Client         int32 `json:"client"`
	Server         int32 `json:"server"`
	LeaderElection int32 `json:"leaderElection"`
}

//==============================================================================
// Private Functions
//==============================================================================
func init() {
	SchemeBuilder.Register(&Arcus{}, &ArcusList{})
}
