package dict

import (
	"go.core/svc"
	"go.core/sys/mdl"
)

// 字典服务
type DictSvc struct {
	*svc.Svc[mdl.Dict]
}

// 构造函数
func NewDictSvc() *DictSvc {
	return &DictSvc{
		Svc: svc.NewSvc[mdl.Dict]("sys:dicts"),
	}
}
