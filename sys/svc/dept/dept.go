package dept

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
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

// 获取部门集合
func (s *DeptSvc) ListDepts(conds map[string]any) ([]mdl.Dept, int64, error) {
	var depts []mdl.Dept
	var count int64
	db := s.
		Make(conds).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Where("parent_id = 0").
		Order("sequence").
		Find(&depts)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return depts, count, err
	}

	return depts, count, nil
}
