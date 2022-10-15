package controllers

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	appsv1alpha1 "k8s-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureSecret ensures ImageSecret is Running in a namespace.
func (r *MyappReconciler) ensureImgSecret(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	imgsecret *corev1.Secret,
) (*reconcile.Result, error) {

	// See if ImageSecret already exists and create if it doesn't
	found := &corev1.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      imgsecret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the ImageSecret
		err = r.Create(context.TODO(), imgsecret)

		if err != nil {
			// ImageSecret creation failed
			return &reconcile.Result{}, err
		} else {
			// ImageSecret creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the ImageSecret not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendImageSecret is a code for creating a ImageSecret
func (r *MyappReconciler) backendImgSecret(v *appsv1alpha1.Myapp) *corev1.Secret {

	userdata := v.Spec.Dockerusername + ":" + v.Spec.Dockerpassword
	auth := b64.StdEncoding.EncodeToString([]byte(userdata))
	auths := "{\"auths\":{\"https://index.docker.io/v1/\":{\"username\":\"" + v.Spec.Dockerusername + "\",\"password\":\"" + v.Spec.Dockerpassword + "\",\"email\":\"" + v.Spec.Dockeremail + "\",\"auth\":\"" + auth + "\"}}}"
	dockerconfigjson := b64.StdEncoding.EncodeToString([]byte(auths))
	Data := make(map[string][]byte)
	key := ".dockerconfigjson"
	Data[key] = []byte(dockerconfigjson)
	imgsecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name + "-imgsecret",
			Namespace: v.Namespace},
		Data: Data,
		Type: corev1.SecretTypeDockerConfigJson,
	}
	fmt.Println(dockerconfigjson)
	controllerutil.SetControllerReference(v, imgsecret, r.Scheme)
	return imgsecret
}
