package v2

import (
	v1 "github.com/yshaojie/log-collector/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// 实现资源版本转换
var _ conversion.Convertible

func (src *ServerLog) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.ServerLog)
	dst.Spec.Dir = src.Spec.Dir
	dst.Spec.NodeName = src.Spec.NodeName
	dst.Spec.FileFilter = ""
	dst.Spec.Pattern = ""
	return nil
}

func (dst *ServerLog) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.ServerLog)
	dst.Spec.Dir = src.Spec.Dir
	dst.Spec.NodeName = src.Spec.NodeName
	return nil
}
