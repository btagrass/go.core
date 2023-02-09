package mdl

import "github.com/btagrass/go.core/mdl"

// 部门
type Dept struct {
	mdl.Mdl
	ParentId int64   `gorm:"comment:父编码" json:"parentId"`             // 父编码
	Name     string  `gorm:"size:50;not null;comment:名称" json:"name"` // 名称
	Phone    string  `gorm:"size:50;comment:电话" json:"phone"`         // 电话
	Addr     string  `gorm:"size:100;comment:地址" json:"addr"`         // 地址
	Sequence int     `gorm:"comment:次序" json:"sequence"`              // 次序
	Children []*Dept `gorm:"foreignKey:ParentId" json:"children"`     // 子部门集合
}

func (m Dept) TableName() string {
	return "sys_dept"
}
