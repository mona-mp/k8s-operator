package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "k8s-operator/api/v1alpha1"
)

func labels(v *appsv1alpha1.Myapp) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app": v.Spec.Name,
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *MyappReconciler) ensureDeployment(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *MyappReconciler) backendDeployment(v *appsv1alpha1.Myapp) *appsv1.Deployment {
	a := corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: v.Spec.Name + "-pvc",
	}

	labels := labels(v)
	size := int32(1)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name,
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           v.Spec.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            v.Spec.Name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: v.Spec.Portnumber,
							Name:          v.Spec.Portname,
						}},
						Env: v.Spec.Envs,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      v.Spec.Name + "-storage",
							MountPath: v.Spec.MountPath,
						}},
					}},
					ImagePullSecrets: []corev1.LocalObjectReference{{Name: v.Spec.Name + "-imgsecret"}},
					Volumes: []corev1.Volume{{
						Name: v.Spec.Name + "-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: a.ClaimName,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

// 	Name: v.Spec.Name + "pvc",
// 	VolumeSource: corev1.VolumeSource{
// 		HostPath: &corev1.HostPathVolumeSource{
// 			Path: "/mnt/data",

// 	},
// }
