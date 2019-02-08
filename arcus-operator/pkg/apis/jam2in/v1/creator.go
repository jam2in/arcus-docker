package v1

import (
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
			"arcus-zk-init.sh": createZkInitScript(arcus),
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
			Labels: GetZkLabel(arcus),
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
			Selector: GetZkLabel(arcus),
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
			Labels:    GetZkLabel(arcus),
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: GetObjectNameZkHeadlessService(arcus),
			Replicas:    &arcus.Spec.Zookeeper.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: GetZkLabel(arcus),
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: GetZkLabel(arcus),
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: GetObjectNameZkVolume(arcus),
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: GetDefaultMode(),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: GetObjectNameConfigMap(arcus),
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:    GetObjectNameZkStatefulSet(arcus),
							Image:   "busybox",
							Command: []string{"sleep", "3600"},
							VolumeMounts: []corev1.VolumeMount{
								{Name: GetObjectNameZkVolume(arcus), MountPath: ZkVolumeMountPath},
							},
						},
					},
				},
			},
		},
	}
}

//==============================================================================
// Private Functions
//==============================================================================
func createZkInitScript(arcus *Arcus) string {
	return "#!/usr/bin/env bash\n\n" +
		"set -ex\n\n" +
		"ZOOKEEPER_CONF_FILE=$ZOOKEEPER_PATH/conf/zoo.cfg\n" +
		"ZOOKEEPER_MY_ID_FILE=$ZOOKEEPER_PATH/data/myid\n\n" +
		"HOST_NAME_SHORT=`hostname -s`\n" +
		"HOST_NAME_DOMAIN=`hostname -d`\n\n" +
		"if [[ $HOST_NAME_SHORT =~ (.*)-([0-9]+)$ ]]; then\n" +
		"    HOST_NAME=${BASH_REMATCH[1]}\n" +
		"    HOST_ORG=${BASH_REMATCH[2]}\n" +
		"fi\n\n" +
		"function create_config() {\n" +
		"    mkdir -p $ZOOKEEPER_PATH/conf\n" +
		"    rm -rf $ZOOKEEPER_PATH/conf\n" +
		"    echo \"maxClientCnxns=" + strconv.Itoa(int(arcus.Spec.Zookeeper.MaxClientCnxns)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"tickTime=" + strconv.Itoa(int(arcus.Spec.Zookeeper.TickTime)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"initLimit=" + strconv.Itoa(int(arcus.Spec.Zookeeper.InitLimit)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"syncLimit=" + strconv.Itoa(int(arcus.Spec.Zookeeper.SyncLimit)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"minSessionTimeout=" + strconv.Itoa(int(arcus.Spec.Zookeeper.MinSessionTimeout)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"maxSessionTimeout=" + strconv.Itoa(int(arcus.Spec.Zookeeper.MaxSessionTimeout)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"clientPort=" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.Client)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    echo \"dataDir=$ZOOKEEPER_PATH/data\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    for (( i=0; i<" + strconv.Itoa(int(arcus.Spec.Zookeeper.Replicas)) + "; i++ ))\n" +
		"    do\n" +
		"         echo \"server.$((i+1))=$HOST_NAME-$i.$HOST_NAME_DOMAIN:" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.Server)) + ":" + strconv.Itoa(int(arcus.Spec.Zookeeper.Ports.LeaderElection)) + "\" >> $ZOOKEEPER_CONF_FILE\n" +
		"    done\n\n" +
		"function create_myid() {\n" +
		"    mkdir -p $ZOOKEEPER_PATH/data\n" +
		"    echo $((HOST_ORG+1)) > $ZOOKEEPER_MY_ID_FILE\n" +
		"}\n\n" +
		"create_config && create_myid"
}
