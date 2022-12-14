/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyappSpec defines the desired state of Myapp
type MyappSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Myapp Name
	Name string `json:"name"`
	// Image address
	Image string `json:"image"`
	// Containerport and Serviceport
	Portnumber int32 `json:"portnumber"`
	// Environment variables
	Envs []corev1.EnvVar `json:"envs,omitempty"`
	// Container volume mountpath
	MountPath string `json:"volumemountpath,omitempty"`
	// Service type
	Servicetype string `json:"servicetype,omitempty"`
	// Service Nodeport
	Servicenodeport int32 `json:"servicenodeport,omitempty"`
	// Ingressclass
	Ingressclass string `json:"ingressclass,omitempty"`
	// Ingress hostname
	Ingresshost string `json:"ingresshost,omitempty"`
	// Secret key
	Secretkey string `json:"secretkey,omitempty"`
	// Secret value
	Secretvalue string `json:"secretvalue,omitempty"`
	// Imagepullsecret dockerconfig json
	Dockerconfigjson string `json:"dockerconfigjson,omitempty"`
	// StorageClass name
	Storageclass string `json:"storageclass,omitempty"`
	// PVC storage resource
	Pvcstorage string `json:"pvcstorage,omitempty"`
	// SVCmonitor enable
	Servicemonitorenable bool `json:"servicemonitorenable,omitempty"`
	// docker username
	Dockerusername string `json:"dockerusername,omitempty"`
	// docker password
	Dockerpassword string `json:"dockerpassword,omitempty"`
	// docker email
	Dockeremail string `json:"dockeremail,omitempty"`
}

// MyappStatus defines the observed state of Myapp
type MyappStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Myapp is the Schema for the myapps API
type Myapp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyappSpec   `json:"spec,omitempty"`
	Status MyappStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyappList contains a list of Myapp
type MyappList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Myapp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Myapp{}, &MyappList{})
}
