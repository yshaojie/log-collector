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

package v1

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var serverloglog = logf.Log.WithName("serverlog-resource")

func (r *ServerLog) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-log-4yxy-io-v1-serverlog,mutating=true,failurePolicy=fail,sideEffects=None,groups=log.4yxy.io,resources=serverlogs,verbs=create;update,versions=v1,name=mserverlog.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ServerLog{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ServerLog) Default() {
	serverloglog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-log-4yxy-io-v1-serverlog,mutating=false,failurePolicy=fail,sideEffects=None,groups=log.4yxy.io,resources=serverlogs,verbs=create;update,versions=v1,name=vserverlog.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ServerLog{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ServerLog) ValidateCreate() (admission.Warnings, error) {
	serverloglog.Info("validate create", "name", r.Name)
	return r.validateCommonField()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ServerLog) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	serverloglog.Info("validate update", "name", r.Name)
	return r.validateCommonField()
}

func (r *ServerLog) validateCommonField() (admission.Warnings, error) {
	if len(r.Spec.Dir) < 2 {
		return admission.Warnings{"waring1...", "waring2..."}, errors.New("spec.dir length< 2")
	}
	return admission.Warnings{"new server log"}, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ServerLog) ValidateDelete() (admission.Warnings, error) {
	serverloglog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
