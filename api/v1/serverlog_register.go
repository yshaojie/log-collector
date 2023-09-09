package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	scheme1 = runtime.NewScheme()
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
	println("init -----------------------------")
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme1))

	utilruntime.Must(AddToScheme1(scheme1))
}
