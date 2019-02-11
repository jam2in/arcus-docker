package arcus

import (
	"context"

	"github.com/go-logr/logr"
	jam2inv1 "github.com/jam2in/arcus-operator/pkg/apis/jam2in/v1"
	"github.com/jam2in/arcus-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	err = r.reconcilePodDisruptionBudget(arcus)
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
	objectName := "ConfigMap"
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
			r.logFailedToGet(objectName, namespace, name, err)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newConfigMap, r.scheme)
		if err != nil {
			r.logFailedToSetControllerReference(objectName, namespace, name, err)
			return err
		}
		r.logFailedToCreate(objectName, namespace, name, err)
		err = r.client.Create(context.TODO(), newConfigMap)
		if err != nil {
			r.logFailedToCreate(objectName, namespace, name, err)
			return err
		}
		return nil
	}
	util.SynchronizeConfigMap(oldConfigMap, newConfigMap)
	err = r.client.Update(context.TODO(), oldConfigMap)
	if err != nil {
		r.logFailedToUpdate(objectName, namespace, name, err)
		return err
	}
	return nil
}

func (r *ReconcileArcus) reconcilePodDisruptionBudget(arcus *jam2inv1.Arcus) error {
	objectName := "PodDisruptionBudget"
	namespace := arcus.Namespace
	name := jam2inv1.GetObjectNamePodDisruptionBudget(arcus)

	newPDB := jam2inv1.CreatePodDisruptionBudget(arcus)
	oldPDB := &policyv1beta1.PodDisruptionBudget{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, oldPDB)

	if err != nil {
		if !errors.IsNotFound(err) {
			r.logFailedToGet(objectName, namespace, name, err)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newPDB, r.scheme)
		if err != nil {
			r.logFailedToSetControllerReference(objectName, namespace, name, err)
			return err
		}
		r.logCreating(objectName, namespace, name)
		err = r.client.Create(context.TODO(), newPDB)
		if err != nil {
			r.logFailedToCreate(objectName, namespace, name, err)
			return err
		}
		return nil
	}

	return nil
}

func (r *ReconcileArcus) reconcileZkHeadlessService(arcus *jam2inv1.Arcus) error {
	objectName := "Service"
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
			r.logFailedToGet(objectName, namespace, name, err)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newService, r.scheme)
		if err != nil {
			r.logFailedToSetControllerReference(objectName, namespace, name, err)
			return err
		}
		r.logCreating(objectName, namespace, name)
		err = r.client.Create(context.TODO(), newService)
		if err != nil {
			r.logFailedToCreate(objectName, namespace, name, err)
			return err
		}
		return nil
	}
	util.SynchronizeService(oldService, newService)
	err = r.client.Update(context.TODO(), oldService)
	if err != nil {
		r.logFailedToUpdate(objectName, namespace, name, err)
		return err
	}

	return nil
}

func (r *ReconcileArcus) reconcileZkStatefulSet(arcus *jam2inv1.Arcus) error {
	objectName := "StatefulSet"
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
			r.logFailedToGet(objectName, namespace, name, err)
			return err
		}
		err = controllerutil.SetControllerReference(arcus, newStatefulSet, r.scheme)
		if err != nil {
			r.logFailedToSetControllerReference(objectName, namespace, name, err)
			return err
		}
		r.logCreating(objectName, namespace, name)
		err = r.client.Create(context.TODO(), newStatefulSet)
		if err != nil {
			r.logFailedToCreate(objectName, namespace, name, err)
			return err
		}
		return nil
	}
	util.SynchronizeStatefulSet(oldStatefulSet, newStatefulSet)
	err = r.client.Update(context.TODO(), oldStatefulSet)
	if err != nil {
		r.logFailedToUpdate(objectName, namespace, name, err)
		return err
	}
	return nil
}

func (r *ReconcileArcus) logCreating(objectName string, namespace string, name string) {
	r.log.Info("Creating a new "+objectName+".",
		objectName+".Namespace", namespace,
		objectName+".Name", name)
}

func (r *ReconcileArcus) logFailedToGet(objectName string, namespace string, name string, err error) {
	r.log.Error(err, "Failed to get "+objectName+".",
		objectName+".Namespace", namespace,
		objectName+".Name", name)
}

func (r *ReconcileArcus) logFailedToSetControllerReference(objectName string, namespace string, name string, err error) {
	r.log.Error(err, "Failed to set controller reference for "+objectName+".",
		objectName+".Namespace", namespace,
		objectName+".Name", name)
}

func (r *ReconcileArcus) logFailedToCreate(objectName string, namespace string, name string, err error) {
	r.log.Error(err, "Failed to create "+objectName+".",
		objectName+".Namespace", namespace,
		objectName+".Name", name)
}

func (r *ReconcileArcus) logFailedToUpdate(objectName string, namespace string, name string, err error) {
	r.log.Error(err, "Failed to update "+objectName+".",
		objectName+".Namespace", namespace,
		objectName+".Name", name)
}
