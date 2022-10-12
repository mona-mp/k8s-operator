package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "k8s-operator/api/v1alpha1"
)

const ()

// ensureIngress ensures Ingress is Running in a namespace.
func (r *MyappReconciler) ensureIngress(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	ingress *networkingv1.Ingress,
) (*reconcile.Result, error) {

	// See if Ingress already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      ingress.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the ingress
		err = r.Create(context.TODO(), ingress)

		if err != nil {
			// Ingress creation failed
			return &reconcile.Result{}, err
		} else {
			// Ingress creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the ingress not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *MyappReconciler) backendIngress(v *appsv1alpha1.Myapp) *networkingv1.Ingress {
	pathtype := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name + "-ingress",
			Namespace: v.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &v.Spec.Ingressclass,
			Rules: []networkingv1.IngressRule{
				{
					Host: v.Spec.Ingresshost,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathtype,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: v.Spec.Name + "-svc",
											Port: networkingv1.ServiceBackendPort{
												Number: v.Spec.Portnumber,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, ingress, r.Scheme)
	return ingress
}
