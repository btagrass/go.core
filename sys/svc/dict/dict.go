package dict

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
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
