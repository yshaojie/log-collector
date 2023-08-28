/*
Copyright 2023.

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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	logv1 "4yxy.io/log-collector/api/v1"
)

// ServerLogReconciler reconciles a ServerLog object
type ServerLogReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=log.4yxy.io,resources=serverlogs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=log.4yxy.io,resources=serverlogs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=log.4yxy.io,resources=serverlogs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServerLog object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ServerLogReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	serverLog := &logv1.ServerLog{}
	var pod v1.Pod
	println(req.Name)
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		//不存在，则不处理
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}
	logDir := pod.GetObjectMeta().GetAnnotations()["server.xy.io/logDir"]
	println("logDir: ", logDir)
	if logDir == "" {
		logDir = "/data/log"
	}
	if err := r.Get(ctx, req.NamespacedName, serverLog); err != nil {
		//不存在说明需要创建
		if errors.IsNotFound(err) {
			newServerLog := &logv1.ServerLog{}
			newServerLog.Spec.Dir = logDir
			newServerLog.Namespace = req.Namespace
			newServerLog.Name = req.Name
			newServerLog.Status.Phase = "Init"
			//newServerLog.GetObjectMeta().SetFinalizers()
			if err := controllerutil.SetControllerReference(&pod, newServerLog, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.Create(ctx, newServerLog); err != nil {
				if errors.IsAlreadyExists(err) {
					if err := r.Get(ctx, req.NamespacedName, serverLog); err != nil {
						return ctrl.Result{}, err
					}
				}
			}
			r.EventRecorder.Event(newServerLog, "Normal", "Created", "create server log")
			return ctrl.Result{}, err
		}
	}
	//日志目录变更，更改serverLog
	if logDir != serverLog.Spec.Dir {
		serverLog.Spec.Dir = logDir
		err := r.Update(ctx, serverLog)

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServerLogReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Owns(&logv1.ServerLog{}).
		Complete(r)
}
