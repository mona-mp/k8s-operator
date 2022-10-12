package controllers

import (
	"context"
	"fmt"

	appsv1alpha1 "k8s-operator/api/v1alpha1"

	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureSvcMonitor ensures SvcMonitor is Running in a namespace.
func (r *MyappReconciler) ensureSvcMonitor(request reconcile.Request,
	instance *appsv1alpha1.Myapp,
	svcmonitor *monitoring.ServiceMonitor,
) (*reconcile.Result, error) {

	// See if SvcMonitor already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      svcmonitor.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the SvcMonitor
		err = r.Create(context.TODO(), svcmonitor)

		if err != nil {
			// SvcMonitor creation failed
			return &reconcile.Result{}, err
		} else {
			// SvcMonitor creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the SvcMonitor not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendSvcMonitor is a code for creating a SvcMonitor
func (r *MyappReconciler) backendSvcMonitor(v *appsv1alpha1.Myapp) *monitoring.ServiceMonitor {

	svcmonitor := &monitoring.ServiceMonitor{

		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.Name + "-svcmonitor",
			Namespace: v.Namespace},
		Spec: monitoring.ServiceMonitorSpec{
			Endpoints: []monitoring.Endpoint{{
				Port: v.Spec.Name,
			}},
			Selector: metav1.LabelSelector{
				MatchLabels: labels(v),
			},
		},
	}

	controllerutil.SetControllerReference(v, svcmonitor, r.Scheme)
	yamlData, _ := yaml.Marshal(&svcmonitor)
	fmt.Println(string(yamlData))
	return svcmonitor
}
