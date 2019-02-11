package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//==============================================================================
// Constant
//==============================================================================
const (
	DefaultZkReplicas = 3

	DefaultZkImage           = "jam2in/arcus:latest"
	DefaultZkImagePullPolicy = "Always"

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
	Zookeeper ZookeeperSpec `json:"zookeeper,omitempty"`
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
// ZookeeperSpec
//==============================================================================
type ZookeeperSpec struct {
	Replicas int32 `json:"replicas"`

	Image ContainerImage `json:"image,omitempty"`

	Ports ZookeeperPort `json:"ports,omitempty"`

	Configuration ZookeeperConfiguration `json:"configuration,omitempty"`
}

func (c *ZookeeperSpec) withDefaults(arcus *Arcus) (changed bool) {
	zkSpec := &arcus.Spec.Zookeeper
	zkImage := &zkSpec.Image
	zkPorts := &zkSpec.Ports
	zkConfig := &zkSpec.Configuration

	if zkSpec.Replicas == 0 {
		changed = true
		zkSpec.Replicas = DefaultZkReplicas
	}

	if zkImage.Name == "" {
		changed = true
		zkImage.Name = DefaultZkImage
	}
	if zkImage.PullPolicy == "" {
		changed = true
		zkImage.PullPolicy = DefaultZkImagePullPolicy
	}

	if zkPorts.Client == 0 {
		changed = true
		zkPorts.Client = DefaultZkClientPort
	}
	if zkPorts.Server == 0 {
		changed = true
		zkPorts.Server = DefaultZkServerPort
	}
	if zkPorts.LeaderElection == 0 {
		changed = true
		zkPorts.LeaderElection = DefaultZkLeaderElectionPort
	}

	if zkConfig.MaxClientCnxns == 0 {
		changed = true
		zkConfig.MaxClientCnxns = DefaultZkMaxClientCnxns
	}
	if zkConfig.TickTime == 0 {
		changed = true
		zkConfig.TickTime = DefaultZkTickTime
	}
	if zkConfig.InitLimit == 0 {
		changed = true
		zkConfig.InitLimit = DefaultZkInitLimit
	}
	if zkConfig.SyncLimit == 0 {
		changed = true
		zkConfig.SyncLimit = DefaultZkSyncLimit
	}
	if zkConfig.MinSessionTimeout == 0 {
		changed = true
		zkConfig.MinSessionTimeout = DefaultZkMinSessionTimeout
	}
	if zkConfig.MaxSessionTimeout == 0 {
		changed = true
		zkConfig.MaxSessionTimeout = DefaultZkMaxSessionTimeout
	}

	return changed
}

//==============================================================================
// ZookeeperConfiguration
//==============================================================================
type ZookeeperConfiguration struct {
	MaxClientCnxns    int32 `json:"maxClientCnxns"`
	TickTime          int32 `json:"tickTime"`
	InitLimit         int32 `json:"initLimit"`
	SyncLimit         int32 `json:"syncLimit"`
	MinSessionTimeout int32 `json:"minSessionTimeout"`
	MaxSessionTimeout int32 `json:"maxSessionTimeout"`
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
// ContainerImage
//==============================================================================
type ContainerImage struct {
	Name       string            `json:"name"`
	PullPolicy corev1.PullPolicy `json:"pullPolicy"`
}

//==============================================================================
// Private Functions
//==============================================================================
func init() {
	SchemeBuilder.Register(&Arcus{}, &ArcusList{})
}
