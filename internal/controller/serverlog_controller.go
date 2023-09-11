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
	logv1 "github.com/yshaojie/log-collector/api/v1"
	"github.com/yshaojie/log-collector/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

//额外添加权限
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

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

	var pod v1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		//不存在，则不处理
		if errors.IsNotFound(err) {
			result, err := r.processDelete(ctx, req)
			if err != nil {
				return result, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	//正在删除状态，不处理
	if !pod.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	serverLog := &logv1.ServerLog{}
	if err := r.Get(ctx, req.NamespacedName, serverLog); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("create server log, %v", req.Name)
			//不存在说明需要创建
			return r.processCreate(ctx, req, pod)
		}
		return ctrl.Result{}, errors.NewInternalError(err)
	}
	return r.processUpdate(ctx, serverLog, pod)
}

func getLogDir(pod v1.Pod) string {
	logDir := pod.GetObjectMeta().GetAnnotations()["server.xy.io/logDir"]
	if logDir == "" {
		logDir = "/data/log"
	}
	return logDir
}

func (r *ServerLogReconciler) processCreate(ctx context.Context, req ctrl.Request, pod v1.Pod) (ctrl.Result, error) {

	newServerLog := &logv1.ServerLog{}
	newServerLog.Spec.Dir = getLogDir(pod)
	newServerLog.Namespace = pod.GetNamespace()
	newServerLog.Name = pod.GetName()
	newServerLog.Status.Phase = "Init"
	//newServerLog.GetObjectMeta().SetFinalizers()
	if err := controllerutil.SetControllerReference(&pod, newServerLog, r.Scheme); err != nil {
		return ctrl.Result{}, errors.NewInternalError(err)
	}
	if err := r.Create(ctx, newServerLog); err != nil {
		err := r.Create(ctx, newServerLog)
		if err != nil {
			println(err.Error())
			return ctrl.Result{}, err
		}
		if errors.IsAlreadyExists(err) {
			serverLog := &logv1.ServerLog{}
			if err := r.Get(ctx, req.NamespacedName, serverLog); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, errors.NewInternalError(err)
	}

	r.EventRecorder.Event(newServerLog, "Normal", "Created", "create server log")
	return ctrl.Result{}, nil
}

func (r *ServerLogReconciler) processDelete(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServerLogReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//Pod为ServerLog的ownerReference，所以需要监听Pod和ServerLog
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Owns(&logv1.ServerLog{}).
		Complete(r)
}

func (r *ServerLogReconciler) processUpdate(ctx context.Context, serverLog *logv1.ServerLog, pod v1.Pod) (ctrl.Result, error) {
	needUpdated := false
	//添加Finalizer，用于清理资源
	if serverLog.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containString(serverLog.ObjectMeta.Finalizers, utils.FinalizerNameAgentHolder) {
			serverLog.ObjectMeta.Finalizers = append(serverLog.ObjectMeta.Finalizers, utils.FinalizerNameAgentHolder)
			needUpdated = true
		}
	}

	if serverLogChange(serverLog, pod) {
		serverLog.Spec.Dir = getLogDir(pod)
		needUpdated = true
	}
	if needUpdated {
		klog.Info("update serverlog ..", " name=", serverLog.Name, " version=", serverLog.ObjectMeta.ResourceVersion)
		err := r.Update(ctx, serverLog)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func containString(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func serverLogChange(log *logv1.ServerLog, pod v1.Pod) bool {
	logDir := getLogDir(pod)
	if logDir != log.Spec.Dir {
		return true
	}

	return false
}
