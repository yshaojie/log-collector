package v2

import (
	v1 "github.com/yshaojie/log-collector/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// 实现资源版本转换
var serverloglog = logf.Log.WithName("serverlog-resource")
var _ conversion.Convertible = &ServerLog{}

func (src *ServerLog) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.ServerLog)
	dst.Spec.Dir = src.Spec.Dir
	dst.Spec.NodeName = src.Spec.NodeName
	dst.ObjectMeta = src.ObjectMeta
	switch src.Status.Phase {
	case ServerLogInit:
		dst.Status.Phase = v1.ServerLogPending
	case ServerLogRunning:
		dst.Status.Phase = v1.ServerLogRunning
	case ServerLogCompleted:
		dst.Status.Phase = v1.ServerLogCompleted
	default:
		dst.Status.Phase = v1.ServerLogPending
	}
	return nil
}

func (dst *ServerLog) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.ServerLog)
	dst.Spec.Dir = src.Spec.Dir
	dst.Spec.NodeName = src.Spec.NodeName
	dst.ObjectMeta = src.ObjectMeta
	switch src.Status.Phase {
	case v1.ServerLogPending:
		dst.Status.Phase = ServerLogInit
	case v1.ServerLogRunning:
		dst.Status.Phase = ServerLogRunning
	case v1.ServerLogCompleted:
		dst.Status.Phase = ServerLogCompleted
	default:
		dst.Status.Phase = ServerLogInit
	}
	return nil
}
