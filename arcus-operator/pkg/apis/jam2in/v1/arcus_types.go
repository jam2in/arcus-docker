package v1

import (
	"github.com/jam2in/arcus-operator/pkg/global"
	corev1 "k8s.io/api/core/v1"
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
	ZkReplicas int32 `json:"zkReplicas"`

	ZkMaxClientCnxns    int32 `json:"zkMaxClientCnxns"`
	ZkTickTime          int32 `json:"zkTickTime"`
	ZkInitLimit         int32 `json:"zkInitLimit"`
	ZkSyncLimit         int32 `json:"zkSyncLimit"`
	ZkMinSessionTimeout int32 `json:"zkMinSessionTimeout"`
	ZkMaxSessionTimeout int32 `json:"zkMaxSessionTimeout"`

	Ports []corev1.ContainerPort `json:"ports,omitempty"`
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
	newPorts := arcus.Spec.Ports
	arcus.Spec.Ports = []corev1.ContainerPort{
		{
			Name:          global.PortNameZkClient,
			ContainerPort: DefaultZkClientPort,
		},
		{
			Name:          global.PortNameZkServer,
			ContainerPort: DefaultZkServerPort,
		},
		{
			Name:          global.PortNameZkLeaderElection,
			ContainerPort: DefaultZkLeaderElectionPort,
		},
	}
	for portIdx := range arcus.Spec.Ports {
		foundPortName := false
		for newPortIdx := range newPorts {
			if arcus.Spec.Ports[portIdx].Name == newPorts[newPortIdx].Name {
				foundPortName = true
				arcus.Spec.Ports[portIdx].ContainerPort = newPorts[newPortIdx].ContainerPort
				break
			}
		}
		if !foundPortName {
			changed = true
		}
	}

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

func (arcus *Arcus) GetZkServicePorts() []corev1.ServicePort {
	zkClientPort := int32(DefaultZkClientPort)
	zkServerPort := int32(DefaultZkServerPort)
	zkLeaderElectionPort := int32(DefaultZkLeaderElectionPort)

	for _, port := range arcus.Spec.Ports {
		switch port.Name {
		case global.PortNameZkClient:
			zkClientPort = port.ContainerPort
		case global.PortNameZkServer:
			zkServerPort = port.ContainerPort
		case global.PortNameZkLeaderElection:
			zkLeaderElectionPort = port.ContainerPort
		}
	}

	return []corev1.ServicePort{
		{
			Name: global.PortNameZkClient,
			Port: zkClientPort,
		},
		{
			Name: global.PortNameZkServer,
			Port: zkServerPort,
		},
		{
			Name: global.PortNameZkLeaderElection,
			Port: zkLeaderElectionPort,
		},
	}
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
