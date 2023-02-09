package mdl

import "github.com/btagrass/go.core/mdl"

// 角色
type Role struct {
	mdl.Mdl
	Name string `gorm:"size:50;not null;comment:名称" json:"name"` // 名称
}

func (m Role) TableName() string {
	return "sys_role"
}
