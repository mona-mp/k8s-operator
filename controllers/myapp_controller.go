package controllers

import (
	"context"
	appsv1alpha1 "k8s-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MyappReconciler reconciles a Myapp object
type MyappReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.my.domain,resources=myapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.my.domain,resources=myapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.my.domain,resources=myapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Myapp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *MyappReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("Myapp", req.NamespacedName)

	// Fetch the Myapp instance
	instance := &appsv1alpha1.Myapp{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	// Check if this PVC already exists
	if len(instance.Spec.Pvcstorage) > 0 {
		result, err = r.ensurePersistentVolumeClaim(req, instance, r.backendPersistentVolumeClaim(instance))
		if result != nil {
			log.Error(err, "PVC Not ready")
			return *result, err
		}
	}

	// Check if this Secret already exists
	if len(instance.Spec.Secretkey) > 0 {
		result, err = r.ensureSecret(req, instance, r.backendSecret(instance))
		if result != nil {
			log.Error(err, "Secret Not ready")
			return *result, err
		}
	}

	// Check if this ImgSecret already exists
	if len(instance.Spec.Dockerusername) > 0 {
		result, err = r.ensureImgSecret(req, instance, r.backendImgSecret(instance))
		if result != nil {
			log.Error(err, "ImgSecret Not ready")
			return *result, err
		}
	}

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)

	result, err = r.ensureDeployment(req, instance, r.backendDeployment(instance))
	if result != nil {
		log.Error(err, "Deployment Not ready")
		return *result, err
	}

	// Check if this Service already exists
	result, err = r.ensureService(req, instance, r.backendService(instance))
	if result != nil {
		log.Error(err, "Service Not ready")
		return *result, err
	}

	// Create Ingress
	if len(instance.Spec.Ingresshost) > 0 {
		result, err = r.ensureIngress(req, instance, r.backendIngress(instance))
		if result != nil {
			log.Error(err, "Ingress Not ready")
			return *result, err
		}
	}

	// Create SVCmonitor
	if instance.Spec.Servicemonitorenable {
		result, err = r.ensureSvcMonitor(req, instance, r.backendSvcMonitor(instance))
		if result != nil {
			log.Error(err, "Servicemonitor Not ready")
			return *result, err
		}
	}

	// Deployment and Service already exists - don't requeue
	log.Info("Skip reconcile: Deployment and service already exists",
		"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyappReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Myapp{}).
		Complete(r)
}
