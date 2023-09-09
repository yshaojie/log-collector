package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	// SchemeBuilder points to a list of functions added to Scheme.
	SchemeBuilder1     = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder1
	// AddToScheme is a common registration function for mapping packaged scoped group & version keys to a scheme.
	AddToScheme1 = localSchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&ServerLog{},
		&ServerLogList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
