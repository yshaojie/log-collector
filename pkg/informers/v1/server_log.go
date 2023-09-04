package v1

import (
	"context"
	apiv1 "github.com/yshaojie/log-collector/api/v1"
	v12 "github.com/yshaojie/log-collector/pkg/listers/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

// ServerLogInformer provides access to a shared informer and lister for
// Deployments.
type ServerLogInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.DeploymentLister
}

type serverLogInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewServerLogInformer constructs a new informer for Deployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewServerLogInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredServerLogInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredServerLogInformer constructs a new informer for Deployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredServerLogInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				var timeout time.Duration
				if options.TimeoutSeconds != nil {
					timeout = time.Duration(*options.TimeoutSeconds) * time.Second
				}

				result := &apiv1.ServerLogList{}
				err := client.DiscoveryV1().RESTClient().Get().
					Namespace(namespace).
					Resource("serverlogs").
					VersionedParams(&options, scheme.ParameterCodec).
					Timeout(timeout).
					Do(context.TODO()).
					Into(result)
				return result, err
			},
			WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&opts)
				}
				var timeout time.Duration
				if opts.TimeoutSeconds != nil {
					timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
				}
				opts.Watch = true
				return client.DiscoveryV1().RESTClient().Get().
					Namespace(namespace).
					Resource("serverlogs").
					VersionedParams(&opts, scheme.ParameterCodec).
					Timeout(timeout).
					Watch(context.TODO())
			},
		},
		&apiv1.ServerLog{},
		resyncPeriod,
		indexers,
	)
}

func (f *serverLogInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredServerLogInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *serverLogInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiv1.ServerLog{}, f.defaultInformer)
}

func (f *serverLogInformer) Lister() v12.ServerLogLister {
	return v12.NewServerLogLister(f.Informer().GetIndexer())
}
