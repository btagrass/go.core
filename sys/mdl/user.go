package mdl

import "github.com/btagrass/go.core/mdl"

// 用户
type User struct {
	mdl.Mdl
	DeptId   int64  `gorm:"not null;comment:部门编码" json:"deptId"`                      // 部门编码
	UserName string `gorm:"uniqueIndex;size:50;not null;comment:用户名" json:"userName"` // 用户名
	FullName string `gorm:"size:50;comment:姓名" json:"fullName"`                       // 姓名
	Mobile   string `gorm:"size:50;not null;comment:手机" json:"mobile"`                // 手机
	Password string `gorm:"size:60;not null;comment:密码" json:"password"`              // 密码
	Frozen   bool   `gorm:"comment:是否冻结" json:"frozen"`                               // 是否冻结
	Token    string `gorm:"-" json:"token"`                                           // 令牌
	Dept     *Dept  `json:"dept"`                                                     // 部门
}

func (m User) TableName() string {
	return "sys_user"
}
