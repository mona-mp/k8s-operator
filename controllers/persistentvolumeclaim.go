package controllers

import (
	"context"
	appsv1alpha1 "k8s-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensurepersistentvolumeclaim ensures persistentvolumeclaim is Running in a namespace.
func (r *MyappReconciler) ensurePersistentVolumeClaim(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	pvc *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {

	// See if persistentvolumeclaim already exists and create if it doesn't
	found := &corev1.PersistentVolumeClaim{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      pvc.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the persistentvolumeclaim
		err = r.Create(context.TODO(), pvc)

		if err != nil {
			// persistentvolumeclaim creation failed
			return &reconcile.Result{}, err
		} else {
			// persistentvolumeclaim creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the persistentvolumeclaim not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendpersistentvolumeclaim is a code for creating a persistentvolumeclaim
func (r *MyappReconciler) backendPersistentVolumeClaim(v *appsv1alpha1.Myapp) *corev1.PersistentVolumeClaim {

	storage := make(corev1.ResourceList)
	q, _ := resource.ParseQuantity(v.Spec.Pvcstorage)
	storage[corev1.ResourceStorage] = q

	storageclass := v.Spec.Storageclass
	pvcmode := corev1.PersistentVolumeFilesystem
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name + "-pvc",
			Namespace: v.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName: v.Spec.Name + "-pv",
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: storage,
			},
			StorageClassName: &storageclass,
			VolumeMode:       &pvcmode,
		},
	}

	controllerutil.SetControllerReference(v, pvc, r.Scheme)
	return pvc
}
