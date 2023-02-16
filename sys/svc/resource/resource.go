package resource

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc/user"
	"gorm.io/gorm"
)

// 资源服务
type ResourceSvc struct {
	*svc.Svc[mdl.Resource]
	userSvc *user.UserSvc
}

// 构造函数
func NewResourceSvc(userSvc *user.UserSvc) *ResourceSvc {
	return &ResourceSvc{
		Svc:     svc.NewSvc[mdl.Resource]("sys:resources"),
		userSvc: userSvc,
	}
}

// 获取菜单集合
func (s *ResourceSvc) ListMenus(userId string) ([]mdl.Resource, error) {
	var resources []mdl.Resource
	if userId == "300000000000001" {
		err := s.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("sequence")
			}).
			Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("sequence")
			}).
			Where("parent_id = 0 and type = 1").
			Order("sequence").
			Find(&resources).Error
		if err != nil {
			return resources, err
		}
	} else {
		var ids []string
		roles, _ := s.userSvc.Perm.GetRolesForUser(userId)
		for _, r := range roles {
			permissions := s.userSvc.Perm.GetPermissionsForUser(r)
			for _, p := range permissions {
				ids = append(ids, p[3])
			}
		}
		err := s.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Where("id in ?", ids).Order("sequence")
			}).
			Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
				return db.Where("id in ?", ids).Order("sequence")
			}).
			Where("parent_id = 0 and type = 1 and id in ?", ids).
			Order("sequence").
			Find(&resources).Error
		if err != nil {
			return resources, err
		}
	}

	return resources, nil
}

// 获取资源集合
func (s *ResourceSvc) ListResources(conds map[string]any) ([]mdl.Resource, int64, error) {
	var resources []mdl.Resource
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
		Find(&resources)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return resources, count, err
	}

	return resources, count, nil
}
