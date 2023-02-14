package dept

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/spf13/cast"
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
func (s *DeptSvc) PageDepts(conds map[string]any) ([]mdl.Dept, int64, error) {
	var depts []mdl.Dept
	var count int64
	current := cast.ToInt(conds["current"])
	size := cast.ToInt(conds["size"])
	delete(conds, "current")
	delete(conds, "size")
	err := s.Db.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Limit(size).
		Offset(size*(current-1)).
		Where("parent_id = 0 or parent_id is null").
		Order("sequence").
		Find(&depts, conds).
		Limit(-1).
		Offset(-1).
		Count(&count).Error
	if err != nil {
		return depts, count, err
	}

	return depts, count, nil
}
