/*
Copyright 2023 Mahboob.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyWebAppSpec defines the desired state of MyWebApp
type MyWebAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Image will be read in reconcile. "image" will be passed in the spec
	Image       string           `json:"image"`
	Port        int              `json:"port"`
	StartHour   int              `json:"start"`
	EndHour     int              `json:"end"`
	Replicas    int              `json:"replica"`
	Deployments []NameSpacedName `json:"deployments"`
}

type NameSpacedName struct {
	Name      string `json:"name"`
	NameSpace string `json:"namespace"`
}

// MyWebAppStatus defines the observed state of MyWebApp
type MyWebAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyWebApp is the Schema for the mywebapps API
type MyWebApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyWebAppSpec   `json:"spec,omitempty"`
	Status MyWebAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyWebAppList contains a list of MyWebApp
type MyWebAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyWebApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyWebApp{}, &MyWebAppList{})
}
