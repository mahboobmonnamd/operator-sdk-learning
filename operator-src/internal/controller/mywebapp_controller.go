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

package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	webappv1 "github.com/mahboobmonnamd/operator-sdk-learning/api/v1"
)

// MyWebAppReconciler reconciles a MyWebApp object
type MyWebAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=web-app.4m.fyi,resources=mywebapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=web-app.4m.fyi,resources=mywebapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=web-app.4m.fyi,resources=mywebapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyWebApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MyWebAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info(fmt.Sprintf("%v \n", strings.Repeat("#", 10)))
	logger.Info("Reconcile Logic Started")
	logger.Info(fmt.Sprintf("%v \n", strings.Repeat("#", 10)))

	instance := &webappv1.MyWebApp{}
	// "r" refers to the reconciler object (ReconcileWebApp in the generated code).
	// "r.Get" is a method defined within the reconciler that retrieves the instance of the custom resource being reconciled.
	// "r.Client" refers to the Kubernetes client provided by the Operator SDK.
	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		logger.Error(err, "Error in client get")
		return ctrl.Result{}, err
	}

	// 1. Pull image and do the deployment
	foundDeployment, err := r.makeDeployment(instance)
	if err != nil {
		logger.Error(err, "Error during the deployment")
	}
	// 2. Create service and apply it
	if _, err := r.makeService(instance); err != nil {
		logger.Error(err, "Error during service creation. Performing rollback of the deployment")
		// If service creation have error delete the deployment
		r.Client.Delete(context.TODO(), foundDeployment)
	}

	// 3. Create replica set based on time
	if err := replicaSetup(instance, r); err != nil {
		logger.Error(err, "Scaling of replica set have error")
	}

	logger.Info(fmt.Sprintf("%v \n", strings.Repeat("#", 10)))
	logger.Info("Reconcile Logic Ended")
	logger.Info(fmt.Sprintf("%v \n", strings.Repeat("#", 10)))

	// Renconcile logic will be triggered every 30 seconds
	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyWebAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.MyWebApp{}).
		Complete(r)
}

// Will define the service definition of the deployment
func (r *MyWebAppReconciler) ServiceDefinition(instance *webappv1.MyWebApp) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": instance.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       int32(instance.Spec.Port), // Defines the host port
					TargetPort: intstr.FromInt(80),        //container port
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}

// Defining deployment metadata
func (r *MyWebAppReconciler) DeploymentDefinition(instance *webappv1.MyWebApp) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": instance.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": instance.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  instance.Name,
							Image: instance.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}
}

func replicaSetup(instance *webappv1.MyWebApp, r *MyWebAppReconciler) error {
	log.Log.Info("Start of replica setup")

	startTime := instance.Spec.StartHour
	endTime := instance.Spec.EndHour
	replicas := int32(instance.Spec.Replicas)

	currentHour := time.Now().UTC().Hour()
	log.Log.Info(fmt.Sprintf("Current Time %d", currentHour))

	if currentHour >= startTime && currentHour <= endTime {
		for _, deploy := range instance.Spec.Deployments {
			deployment := &appsv1.Deployment{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{
				Namespace: deploy.NameSpace,
				Name:      deploy.Name,
			}, deployment)
			if err != nil {
				return err
			}
			if deployment.Spec.Replicas != &replicas {
				log.Log.Info(fmt.Sprintf("Increasing replicas to %v", replicas))
				*deployment.Spec.Replicas = replicas
				if err := r.Client.Update(context.TODO(), deployment); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *MyWebAppReconciler) makeService(instance *webappv1.MyWebApp) (ctrl.Result, error) {
	// 1 Define Service
	service := r.ServiceDefinition(instance)
	// Set owner reference to enable garbage collection
	if err := controllerutil.SetControllerReference(instance, service, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Check if the service exists, create or update accordingly
	foundService := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, foundService)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create the service
			err = r.Client.Create(context.TODO(), service)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	}
	// Update the service
	foundService.Spec = service.Spec
	err = r.Client.Update(context.TODO(), foundService)
	return reconcile.Result{}, err
}

// Check if the deployment already exists, create or update accordingly
func (r *MyWebAppReconciler) makeDeployment(instance *webappv1.MyWebApp) (*appsv1.Deployment, error) {

	foundDeployment := &appsv1.Deployment{}

	// 1 Define the deployment meta data.
	deployment := r.DeploymentDefinition(instance)

	// Set owner reference to enable garbage collection
	if err := controllerutil.SetControllerReference(instance, deployment, r.Scheme); err != nil {
		return foundDeployment, err
	}

	// "r.Client.Get" is a method on the client that interacts directly with the Kubernetes API server to retrieve a resource.
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
	}, foundDeployment)

	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), deployment)
		if err != nil {
			return foundDeployment, err
		}
	} else if err == nil {
		// // commented generic updates, as i have specific updates to replicas alone
		// foundDeployment.Spec = deployment.Spec
		// err = r.Client.Update(context.TODO(), foundDeployment)
		// if err != nil {
		// 	return foundDeployment, err
		// }
		return foundDeployment, nil
	}
	return foundDeployment, err
}
