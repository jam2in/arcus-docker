package arcus

import (
	"context"

	"github.com/go-logr/logr"
	jam2inv1 "github.com/jam2in/arcus-operator/pkg/apis/jam2in/v1"
	"github.com/jam2in/arcus-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_arcus")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileArcus{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("arcus-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &jam2inv1.Arcus{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jam2inv1.Arcus{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jam2inv1.Arcus{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jam2inv1.Arcus{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jam2inv1.Arcus{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileArcus{}

type ReconcileArcus struct {
	client client.Client
	scheme *runtime.Scheme
	log    logr.Logger
}

func (r *ReconcileArcus) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	r.log = log.WithValues(
		"Request.Namespace", request.Namespace,
		"Request.Name", request.Name)
	r.log.Info("Reconciling Arcus")

	arcus := &jam2inv1.Arcus{}
	result, err := r.reconcileArcus(arcus, &request)
	if result != nil {
		return *result, err
	}

	err = r.reconcileConfigMap(arcus)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileZkHeadlessService(arcus)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.reconcileZkStatefulSet(arcus)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileArcus) reconcileArcus(arcus *jam2inv1.Arcus, request *reconcile.Request) (*reconcile.Result, error) {
	namespace := arcus.Namespace
	name := arcus.Name

	err := r.client.Get(context.TODO(), request.NamespacedName, arcus)
	if err != nil {
		if errors.IsNotFound(err) {
			r.log.Error(err, "Arcus resource not found. Ignoring since object must be deleted.")
			return &reconcile.Result{}, nil
		}
		r.log.Error(err, "Failed to get Arcus.")
		return &reconcile.Result{}, err
	}

	changed := arcus.WithDefaults()
	if changed {
		r.log.Info("Setting default settings for Arcus.")
		if err := r.client.Update(context.TODO(), arcus); err != nil {
			r.log.Error(err, "Failed to update Arcus.",
				"Arcus.Namespace", namespace,
				"Arcus.Name", name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil
	}
	return nil, nil
}

func (r *ReconcileArcus) reconcileConfigMap(arcus *jam2inv1.Arcus) error {
	namespace := arcus.Namespace
	name := jam2inv1.GetObjectNameConfigMap(arcus)

	newConfigMap := jam2inv1.CreateConfigMap(arcus)
	oldConfigMap := &corev1.ConfigMap{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, oldConfigMap)
	if err != nil {
		if !errors.IsNotFound(err) {
			r.log.Error(err, "Failed to get ConfigMap.",
				"ConfigMap.Namespace", namespace,
				"ConfigMap.Name", name)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newConfigMap, r.scheme)
		if err != nil {
			r.log.Error(err, "Failed to set controller reference for ConfigMap.",
				"ConfigMap.Namespace", namespace,
				"ConfigMap.Name", name)
			return err
		}
		r.log.Info("Creating a new ConfigMap",
			"ConfigMap.Namespace", namespace,
			"ConfigMap.Name", name)
		err = r.client.Create(context.TODO(), newConfigMap)
		if err != nil {
			r.log.Error(err, "Failed to create ConfigMap.",
				"ConfigMap.Namespace", namespace,
				"ConfigMap.Name", name)
			return err
		}
		return nil
	}
	r.log.Info("Updating the existing ConfigMap",
		"ConfigMap.Namespace", namespace,
		"ConfigMap.Name", name)
	util.SynchronizeConfigMap(oldConfigMap, newConfigMap)
	err = r.client.Update(context.TODO(), oldConfigMap)
	if err != nil {
		r.log.Error(err, "Failed to update ConfigMap.",
			"ConfigMap.Namespace", namespace,
			"ConfigMap.Name", name)
		return err
	}
	return nil
}

func (r *ReconcileArcus) reconcileZkHeadlessService(arcus *jam2inv1.Arcus) error {
	namespace := arcus.Namespace
	name := jam2inv1.GetObjectNameZkHeadlessService(arcus)

	newService := jam2inv1.CreateZkHeadlessService(arcus)
	oldService := &corev1.Service{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, oldService)
	if err != nil {
		if !errors.IsNotFound(err) {
			r.log.Error(err, "Failed to get Service.",
				"Service.Namespace", namespace,
				"Service.Name", name)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newService, r.scheme)
		if err != nil {
			r.log.Error(err, "Failed to set controller reference for Service.",
				"Service.Namespace", namespace,
				"Service.Name", name)
			return err
		}
		r.log.Info("Creating a new Service",
			"Service.Namespace", namespace,
			"Service.Name", name)
		err = r.client.Create(context.TODO(), newService)
		if err != nil {
			r.log.Error(err, "Failed to create Service.",
				"Service.Namespace", namespace,
				"Service.Name", name)
			return err
		}
		return nil
	}
	r.log.Info("Updating the existing Service",
		"Service.Namespace", namespace,
		"Service.Name", name)
	util.SynchronizeService(oldService, newService)
	err = r.client.Update(context.TODO(), oldService)
	if err != nil {
		r.log.Error(err, "Failed to update Service.",
			"Service.Namespace", namespace,
			"Service.Name", name)
		return err
	}

	return nil
}

func (r *ReconcileArcus) reconcileZkStatefulSet(arcus *jam2inv1.Arcus) error {
	namespace := arcus.Namespace
	name := jam2inv1.GetObjectNameZkStatefulSet(arcus)

	newStatefulSet := jam2inv1.CreateZkStatefulSet(arcus)
	oldStatefulSet := &appsv1.StatefulSet{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, oldStatefulSet)
	if err != nil {
		if !errors.IsNotFound(err) {
			r.log.Error(err, "Failed to get StatefulSet.",
				"StatefulSet.Namespace", namespace,
				"StatefulSet.Name", name)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newStatefulSet, r.scheme)
		if err != nil {
			r.log.Error(err, "Failed to set controller reference for StatefulSet.",
				"StatefulSet.Namespace", namespace,
				"StatefulSet.Name", name)
			return err
		}
		r.log.Info("Creating a new StatefulSet",
			"StatefulSet.Namespace", namespace,
			"StatefulSet.Name", name)
		err = r.client.Create(context.TODO(), newStatefulSet)
		if err != nil {
			r.log.Error(err, "Failed to create StatefulSet.",
				"StatefulSet.Namespace", namespace,
				"StatefulSet.Name", name)
			return err
		}
		return nil
	}
	r.log.Info("Updating the existing StatefulSet",
		"StatefulSet.Namespace", namespace,
		"StatefulSet.Name", name)
	util.SynchronizeStatefulSet(oldStatefulSet, newStatefulSet)
	err = r.client.Update(context.TODO(), oldStatefulSet)
	if err != nil {
		r.log.Error(err, "Failed to update StatefulSet.",
			"StatefulSet.Namespace", namespace,
			"StatefulSet.Name", name)
		return err
	}
	return nil
}

/*
func (r *ReconcileArcus) reconcileStatefulSet(cr *jam2inv1.Arcus) error {
	sts := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, sts)
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new Arcus StatefulSet",
			"StatefulSet.Namespace", sts.Namespace,
			"StatefulSet.Name", sts.Name)
		err = r.client.Create(context.TODO(), newStatefulSet(cr))
	}
}

func newStatefulSet(cr *jam2inv1.Arcus) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": cr.Name,
			},
			Spec: appsv1.StatefulSetSpec{
				ServiceName: "arcus-headless",
				Replicas:    &cr.Spec.ZkReplicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": cr.Name,
					},
				},
				// TODO: UpdateStrategy
				// TODO: PodManagementPolicy
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: cr.GetName(),
						Labels: map[string]string{
							"app": cr.Name,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							corev1.Container{
								Name:  "zookeeper",
								Image: "emccorp/zookeeper:3.5.4-beta-operator",
								Ports: []corev1.ContainerPort{
									{
										Name:          "client",
										ContainerPort: 2181,
									},
									{
										Name:          "quorum",
										ContainerPort: 2888,
									},
									{
										Name:          "leader-election",
										ContainerPort: 3888,
									},
								},
								ImagePullPolicy: "Always",
								Command:         []string{"/usr/local/bin/zookeeperStart.sh"},
							},
						},
					},
				},
			},
		},
	}
}
*/

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *jam2inv1.Arcus) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
