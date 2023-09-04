package v1

import (
	apiv1 "github.com/yshaojie/log-collector/api/v1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ServerLogLister helps list ServerLogs.
// All objects returned here must be treated as read-only.
type ServerLogLister interface {
	// List lists all ServerLogs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*apiv1.ServerLog, err error)
	// ServerLogs returns an object that can list and get ServerLogs.
	ServerLogs(namespace string) ServerLogNamespaceLister
}

// ServerLogLister implements the ServerLogLister interface.
type serverLogLister struct {
	indexer cache.Indexer
}

// NewServerLogLister returns a new ServerLogLister.
func NewServerLogLister(indexer cache.Indexer) ServerLogLister {
	return &serverLogLister{indexer: indexer}
}

// List lists all ServerLogs in the indexer.
func (s *serverLogLister) List(selector labels.Selector) (ret []*apiv1.ServerLog, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*apiv1.ServerLog))
	})
	return ret, err
}

// ServerLogs returns an object that can list and get ServerLogs.
func (s *serverLogLister) ServerLogs(namespace string) ServerLogNamespaceLister {
	return serverLogNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServerLogNamespaceLister helps list and get ServerLogs.
// All objects returned here must be treated as read-only.
type ServerLogNamespaceLister interface {
	// List lists all ServerLogs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*apiv1.ServerLog, err error)
	// Get retrieves the ServerLog from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*apiv1.ServerLog, error)
}

// ServerLogNamespaceLister implements the ServerLogNamespaceLister
// interface.
type serverLogNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ServerLogs in the indexer for a given namespace.
func (s serverLogNamespaceLister) List(selector labels.Selector) (ret []*apiv1.ServerLog, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*apiv1.ServerLog))
	})
	return ret, err
}

// Get retrieves the ServerLog from the indexer for a given namespace and name.
func (s serverLogNamespaceLister) Get(name string) (*apiv1.ServerLog, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("ServerLog"), name)
	}
	return obj.(*apiv1.ServerLog), nil
}
