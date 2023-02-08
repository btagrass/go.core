package dept

import (
	"github.com/spf13/cast"
	"go.core/svc"
	"go.core/sys/mdl"
	"gorm.io/gorm"
)

// 部门服务
type DeptSvc struct {
	*svc.Svc[mdl.Dept]
}

// 构造函数
func NewDeptSvc() *DeptSvc {
	return &DeptSvc{
		Svc: svc.NewSvc[mdl.Dept]("sys:depts"),
	}
}

// 分页部门集合
func (d *DeptSvc) PageDepts(conds map[string]any) ([]mdl.Dept, int64, error) {
	var depts []mdl.Dept
	var count int64
	size := cast.ToInt(conds["size"])
	current := cast.ToInt(conds["current"])
	delete(conds, "size")
	delete(conds, "current")
	err := d.Db.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Where("parent_id = 0 or parent_id is null").
		Limit(size).
		Offset(size*(current-1)).
		Find(&depts, conds).
		Order("sequence").
		Count(&count).Error
	if err != nil {
		return depts, count, err
	}

	return depts, count, nil
}
