package controllers

import (
	"context"
	appsv1alpha1 "k8s-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureService ensures Service is Running in a namespace.
func (r *MyappReconciler) ensureService(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	service *corev1.Service,
) (*reconcile.Result, error) {

	// See if service already exists and create if it doesn't
	found := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      service.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		err = r.Create(context.TODO(), service)

		if err != nil {
			// Service creation failed
			return &reconcile.Result{}, err
		} else {
			// Service creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendService is a code for creating a Service
func (r *MyappReconciler) backendService(v *appsv1alpha1.Myapp) *corev1.Service {
	labels := labels(v)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name + "-svc",
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Name:       v.Spec.Name,
				Port:       v.Spec.Portnumber,
				TargetPort: intstr.FromInt(int(v.Spec.Portnumber)),
				NodePort:   v.Spec.Servicenodeport,
			}},
			Type: corev1.ServiceType(v.Spec.Servicetype),
		},
	}

	controllerutil.SetControllerReference(v, service, r.Scheme)

	return service
}
