package mdl

import "go.core/mdl"

// 资源
type Resource struct {
	mdl.Model
	ParentId int64       `gorm:"comment:父编码" json:"parentId"`             // 父编码
	Code     string      `gorm:"size:20;not null;comment:代码" json:"code"` // 代码
	Name     string      `gorm:"size:50;not null;comment:名称" json:"name"` // 名称
	Type     int8        `gorm:"comment:类型" json:"type"`                  // 类型
	Icon     string      `gorm:"size:50;comment:图标" json:"icon"`          // 图标
	Url      string      `gorm:"size:100;comment:网址" json:"url"`          // 网址
	Act      string      `gorm:"size:20;comment:动作" json:"act"`           // 动作
	Sequence int         `gorm:"comment:次序" json:"sequence"`              // 次序
	Children []*Resource `gorm:"foreignKey:ParentId" json:"children"`     // 子资源集合
}

func (Resource) TableName() string {
	return "sys_resource"
}
