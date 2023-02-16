package mdl

import "github.com/btagrass/go.core/mdl"

// 字典
type Dict struct {
	mdl.Mdl
	Type     string `gorm:"size:50;not null;comment:类型" json:"type"` // 类型
	Code     int8   `gorm:"comment:代码" json:"code"`                  // 代码
	Name     string `gorm:"size:50;not null;comment:名称" json:"name"` // 名称
	Sequence int    `gorm:"comment:次序" json:"sequence"`              // 次序
}

func (m Dict) TableName() string {
	return "sys_dict"
}
