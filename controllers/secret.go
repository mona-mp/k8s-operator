package controllers

import (
	"context"

	appsv1alpha1 "k8s-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureSecret ensures Secret is Running in a namespace.
func (r *MyappReconciler) ensureSecret(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	secret *corev1.Secret,
) (*reconcile.Result, error) {

	// See if Secret already exists and create if it doesn't
	found := &corev1.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the secret
		err = r.Create(context.TODO(), secret)

		if err != nil {
			// secret creation failed
			return &reconcile.Result{}, err
		} else {
			// secret creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendSecret is a code for creating a Secret
func (r *MyappReconciler) backendSecret(v *appsv1alpha1.Myapp) *corev1.Secret {
	Data := make(map[string][]byte)
	Data[v.Spec.Secretkey] = []byte(v.Spec.Secretvalue)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name + "-secret",
			Namespace: v.Namespace},
		Data: Data,
		Type: corev1.SecretTypeOpaque,
	}

	controllerutil.SetControllerReference(v, secret, r.Scheme)
	return secret
}
