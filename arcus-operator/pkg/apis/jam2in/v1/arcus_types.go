package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//==============================================================================
// Constants
//==============================================================================
const (
	// Default ZooekeperSpec
	DefaultZkReplicas = 3

	// Default ContainerImage
	DefaultZkImage           = "zookeeper/3.4.13:latest"
	DefaultZkImagePullPolicy = "Always"

	// Default PodPolicy
	DefaultZkTerminationGracePeriodSeconds = corev1.DefaultTerminationGracePeriodSeconds

	// Default ZookeeperPort
	DefaultZkClientPort         = 2181
	DefaultZkServerPort         = 2888
	DefaultZkLeaderElectionPort = 3888

	// Default ZookeeperDirectory
	DefaultZkDirHome = "/zookeeper-3.4.13"
	DefaultZkDirConf = "/conf"
	DefaultZkDirData = "/data"

	// Default ZookeeperConfiguration
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

	Pod PodPolicy `json:"pod,omitempty"`

	Ports ZookeeperPort `json:"ports,omitempty"`

	Directory ZookeeperDirectory `json:directory,omitempty`

	Configuration ZookeeperConfiguration `json:"configuration,omitempty"`
}

func (spec *ZookeeperSpec) withDefaults(arcus *Arcus) (changed bool) {
	if arcus.Spec.Zookeeper.Replicas == 0 {
		changed = true
		arcus.Spec.Zookeeper.Replicas = DefaultZkReplicas
	}

	if arcus.Spec.Zookeeper.Image.withZkDefaults() {
		changed = true
	}

	if arcus.Spec.Zookeeper.Pod.withZkDefaults() {
		changed = true
	}

	if arcus.Spec.Zookeeper.Ports.withDefaults() {
		changed = true
	}

	if arcus.Spec.Zookeeper.Directory.withDefaults() {
		changed = true
	}

	if arcus.Spec.Zookeeper.Configuration.withDefaults() {
		changed = true
	}

	return changed
}

//==============================================================================
// ContainerImage
//==============================================================================
type ContainerImage struct {
	Name       string            `json:"name"`
	PullPolicy corev1.PullPolicy `json:"pullPolicy"`
}

func (image *ContainerImage) withZkDefaults() (changed bool) {
	if image.Name == "" {
		changed = true
		image.Name = DefaultZkImage
	}
	if image.PullPolicy == "" {
		changed = true
		image.PullPolicy = DefaultZkImagePullPolicy
	}

	return changed
}

//==============================================================================
// PodPolicy
//==============================================================================
type PodPolicy struct {
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	Resources corev1.ResourceRequirements `json:resources,omitempty`

	Toleration []corev1.Toleration `json:toleration,omitempty`

	Env []corev1.EnvVar `json:"env.omitempty"`

	TerminationGracePeriodSeconds int64 `json:"terminationGracePeriodSeconds"`
}

func (pod *PodPolicy) withZkDefaults() (changed bool) {
	if pod.TerminationGracePeriodSeconds == 0 {
		changed = true
		pod.TerminationGracePeriodSeconds = DefaultZkTerminationGracePeriodSeconds
	}

	if pod.Affinity == nil {
		changed = true
		pod.Affinity = &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					{
						Weight: 100,
						PodAffinityTerm: corev1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      LabelKeyApp,
										Operator: metav1.LabelSelectorOpIn,
										Values:   []string{LabelValueZk},
									},
								},
							},
							TopologyKey: "kubernetes.io/hostname",
						},
					},
				},
			},
		}
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

func (ports *ZookeeperPort) withDefaults() (changed bool) {
	if ports.Client == 0 {
		changed = true
		ports.Client = DefaultZkClientPort
	}
	if ports.Server == 0 {
		changed = true
		ports.Server = DefaultZkServerPort
	}
	if ports.LeaderElection == 0 {
		changed = true
		ports.LeaderElection = DefaultZkLeaderElectionPort
	}

	return changed
}

//==============================================================================
// ZookeeperDirectory
//==============================================================================
type ZookeeperDirectory struct {
	Home string `json:"home"`
	Conf string `json:"conf"`
	Data string `json:"data"`
}

func (directory *ZookeeperDirectory) withDefaults() (changed bool) {
	if directory.Home == "" {
		changed = true
		directory.Home = DefaultZkDirHome
	}
	if directory.Conf == "" {
		changed = true
		directory.Conf = DefaultZkDirConf
	}
	if directory.Data == "" {
		changed = true
		directory.Data = DefaultZkDirData
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

func (configuration *ZookeeperConfiguration) withDefaults() (changed bool) {
	if configuration.MaxClientCnxns == 0 {
		changed = true
		configuration.MaxClientCnxns = DefaultZkMaxClientCnxns
	}
	if configuration.TickTime == 0 {
		changed = true
		configuration.TickTime = DefaultZkTickTime
	}
	if configuration.InitLimit == 0 {
		changed = true
		configuration.InitLimit = DefaultZkInitLimit
	}
	if configuration.SyncLimit == 0 {
		changed = true
		configuration.SyncLimit = DefaultZkSyncLimit
	}
	if configuration.MinSessionTimeout == 0 {
		changed = true
		configuration.MinSessionTimeout = DefaultZkMinSessionTimeout
	}
	if configuration.MaxSessionTimeout == 0 {
		changed = true
		configuration.MaxSessionTimeout = DefaultZkMaxSessionTimeout
	}

	return changed
}

//==============================================================================
// Private Functions
//==============================================================================
func init() {
	SchemeBuilder.Register(&Arcus{}, &ArcusList{})
}
