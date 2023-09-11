package v1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	// SchemeBuilder points to a list of functions added to Scheme.
	SchemeBuilder1     = runtime.NewSchemeBuilder(addKnownTypes, addConversionFuncs)
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

func addConversionFuncs(scheme *runtime.Scheme) error {
	// Add field label conversions for kinds having selectable nothing but ObjectMeta fields.
	if err := scheme.AddFieldLabelConversionFunc(GroupVersion.WithKind("ServerLog"),
		func(label, value string) (string, string, error) {
			switch label {
			case "spec.nodeName":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported for ServerLog: %s", label)
			}
		}); err != nil {
		return err
	}

	return nil
}

func init() {
	//将ServerLog ServerLogList注册到clientgoscheme.Scheme上
	utilruntime.Must(AddToScheme1(clientgoscheme.Scheme))
}
