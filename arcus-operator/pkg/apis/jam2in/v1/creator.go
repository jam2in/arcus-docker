package v1

import (
	"fmt"
	"strconv"

	"github.com/jam2in/arcus-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//==============================================================================
// Public Functions
//==============================================================================
func CreateConfigMap(arcus *Arcus) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetObjectNameConfigMap(arcus),
			Namespace: arcus.Namespace,
		},
		Data: map[string]string{
			FileZkInitScript: createZkInitScript(arcus),
			FileZkOkScript:   createZkOkScript(arcus),
		},
	}
}

func CreateZkHeadlessService(arcus *Arcus) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetObjectNameZkHeadlessService(arcus),
			Namespace: arcus.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(arcus, schema.GroupVersionKind{
					Group:   SchemeGroupVersion.Group,
					Version: SchemeGroupVersion.Version,
					Kind:    KindNameArcus,
				}),
			},
			Labels: map[string]string{
				LabelKeyApp: LabelValueZk,
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Name: PortNameZkClient,
					Port: arcus.Spec.Zookeeper.Ports.Client,
				},
				{
					Name: PortNameZkServer,
					Port: arcus.Spec.Zookeeper.Ports.Server,
				},
				{
					Name: PortNameZkLeaderElection,
					Port: arcus.Spec.Zookeeper.Ports.LeaderElection,
				},
			},
			Selector: map[string]string{
				LabelKeyApp: LabelValueZk,
			},
		},
	}
}

func CreateZkStatefulSet(arcus *Arcus) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetObjectNameZkStatefulSet(arcus),
			Namespace: arcus.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: GetObjectNameZkHeadlessService(arcus),
			Replicas:    &arcus.Spec.Zookeeper.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					LabelKeyApp: LabelValueZk,
				},
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						LabelKeyApp: LabelValueZk,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector:                  arcus.Spec.Zookeeper.Pod.NodeSelector,
					Affinity:                      arcus.Spec.Zookeeper.Pod.Affinity,
					Tolerations:                   arcus.Spec.Zookeeper.Pod.Toleration,
					TerminationGracePeriodSeconds: &arcus.Spec.Zookeeper.Pod.TerminationGracePeriodSeconds,
					Volumes: []corev1.Volume{
						{
							Name: GetObjectNameZkVolume(arcus),
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: defaultMode(),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: GetObjectNameConfigMap(arcus),
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:            GetObjectNameZkStatefulSet(arcus),
							Image:           arcus.Spec.Zookeeper.Image.Name,
							ImagePullPolicy: arcus.Spec.Zookeeper.Image.PullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          PortNameZkClient,
									ContainerPort: arcus.Spec.Zookeeper.Ports.Client,
								},
								{
									Name:          PortNameZkServer,
									ContainerPort: arcus.Spec.Zookeeper.Ports.Server,
								},
								{
									Name:          PortNameZkLeaderElection,
									ContainerPort: arcus.Spec.Zookeeper.Ports.LeaderElection,
								},
							},
							Resources: arcus.Spec.Zookeeper.Pod.Resources,
							Env:       arcus.Spec.Zookeeper.Pod.Env,
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: ProbeInitialDelaySecondsZk,
								TimeoutSeconds:      ProbeTimeoutSecondsZk,
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{Command: []string{util.MakeFilePath(DirPathZkScripts, FileZkOkScript)}},
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: ProbeInitialDelaySecondsZk,
								TimeoutSeconds:      ProbeTimeoutSecondsZk,
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{Command: []string{util.MakeFilePath(DirPathZkScripts, FileZkOkScript)}},
								},
							},
							Command: []string{"/bin/bash"},
							Args:    []string{"-c", util.MakeFilePath(DirPathZkScripts, FileZkInitScript) + " && " + arcus.Spec.Zookeeper.Directory.Home + "/bin/zkServer.sh start-foreground"},
							VolumeMounts: []corev1.VolumeMount{
								{Name: GetObjectNameZkVolume(arcus), MountPath: DirPathZkScripts},
							},
						},
					},
				},
			},
		},
	}
}

func CreatePodDisruptionBudget(arcus *Arcus) *policyv1beta1.PodDisruptionBudget {
	pdbCount := intstr.FromInt(PDBMaxUnavailble)
	return &policyv1beta1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: "policy/v1beta",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetObjectNamePodDisruptionBudget(arcus),
			Namespace: arcus.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(arcus, schema.GroupVersionKind{
					Group:   SchemeGroupVersion.Group,
					Version: SchemeGroupVersion.Version,
					Kind:    KindNameArcus,
				}),
			},
		},
		Spec: policyv1beta1.PodDisruptionBudgetSpec{
			MaxUnavailable: &pdbCount,
			Selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      LabelKeyApp,
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{LabelValueZk, LabelValueMc},
					},
				},
			},
		},
	}
}

//==============================================================================
// Private Functions
//==============================================================================
func headlessDomain(arcus *Arcus) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", GetObjectNameZkHeadlessService(arcus), arcus.Namespace)
}

func defaultMode() *int32 {
	defaultMode := new(int32)
	*defaultMode = int32(0755)
	return defaultMode
}

func createZkInitScript(arcus *Arcus) string {
	return "#!/bin/bash\n\n" +
		"set -ex\n\n" +
		"ZOOKEEPER_CONF_DIR=" + arcus.Spec.Zookeeper.Directory.Conf + "\n" +
		"ZOOKEEPER_CONF_FILE=$ZOOKEEPER_CONF_DIR/zoo.cfg\n" +
		"ZOOKEEPER_DATA_DIR=" + arcus.Spec.Zookeeper.Directory.Data + "\n" +
		"ZOOKEEPER_MYID_FILE=$ZOOKEEPER_DATA_DIR/myid\n\n" +
		"HOST_SHORT_NAME=`hostname -s`\n" +
		"HOST_DOMAIN_NAME=" + headlessDomain(arcus) + "\n\n" +
		"if [[ $HOST_SHORT_NAME =~ (.*)-([0-9]+)$ ]]; then\n" +
		"    HOST_NAME=${BASH_REMATCH[1]}\n" +
		"    HOST_ORG=${BASH_REMATCH[2]}\n" +
		"fi\n\n" +
		"function create_config() {\n" +
		"    set +e\n" +
		"    mkdir -p $ZOOKEEPER_CONF_DIR\n" +
		"    set -e\n" +
		"    echo \"maxClientCnxns=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Configuration.MaxClientCnxns)) + "\" > $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"tickTime=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Configuration.TickTime)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"initLimit=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Configuration.InitLimit)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"syncLimit=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Configuration.SyncLimit)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"minSessionTimeout=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Configuration.MinSessionTimeout)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"maxSessionTimeout=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Configuration.MaxSessionTimeout)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"clientPort=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.Client)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"dataDir=$ZOOKEEPER_DATA_DIR\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    for (( i=0; i<" + strconv.Itoa(int(arcus.Spec.Zookeeper.Replicas)) + "; i++ ))\n" +
		"    do\n" +
		"         if [ $i -eq $HOST_ORG ]; then\n" +
		"							# Zookeeper UnknownHostException issue in Kubernetes\n" +
		"             # https://github.com/kubernetes/contrib/issues/2737\n" +
		"							# https://stackoverflow.com/questions/46605686/zookeeper-hostname-resolution-fails\n" +
		"             # https://hub.docker.com/_/zookeeper?tab=description\n" +
		"             echo \"server.$((i+1))=0.0.0.0:" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.Server)) + ":" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.LeaderElection)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"         else\n" +
		"             echo \"server.$((i+1))=$HOST_NAME-$i.$HOST_DOMAIN_NAME:" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.Server)) + ":" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.LeaderElection)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"         fi\n" +
		"    done\n" +
		"}\n\n" +
		"function create_myid() {\n" +
		"    set +e\n" +
		"    mkdir -p $ZOOKEEPER_DATA_DIR\n" +
		"    set -e\n" +
		"    echo $((HOST_ORG+1)) > $ZOOKEEPER_MYID_FILE\n" +
		"}\n\n" +
		"create_config && create_myid"
}

func createZkOkScript(arcus *Arcus) string {
	return "#!/bin/bash\n\n" +
		"ZK_CLIENT_PORT=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.Client)) + "\n" +
		"OK=$(echo ruok | nc 127.0.0.1 $ZK_CLIENT_PORT)\n" +
		"if [ \"$OK\" == \"imok\" ]; then\n" +
		"    echo \"Zookeeper service is available.\"\n" +
		"    exit 0\n" +
		"else\n" +
		"    echo \"Zookeeper service is not available for request.\"\n" +
		"    exit 1\n" +
		"fi"
}
