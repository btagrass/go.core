package resource

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc/user"
	"github.com/spf13/cast"
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
			Order("sequence").
			Find(&resources, "parent_id = 0 and type = 1").Error
		if err != nil {
			return resources, err
		}
	} else {
		var ids []string
		roles, _ := s.userSvc.Perm.GetRolesForUser(cast.ToString(userId))
		for _, role := range roles {
			permissions := s.userSvc.Perm.GetPermissionsForUser(role)
			for _, p := range permissions {
				ids = append(ids, p[3])
			}
		}
		err := s.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Where("id in ?", ids).Order("sequence")
			}).
			Find(&resources, "parent_id = 0 and type = 1 and id in ?", ids).
			Order("sequence").Error
		if err != nil {
			return resources, err
		}
	}

	return resources, nil
}

// 获取资源集合
func (s *ResourceSvc) ListResources(conds map[string]any) ([]mdl.Resource, error) {
	var resources []mdl.Resource
	user, ok := conds["user"]
	if ok {
		delete(conds, "user")
		var resourceIds []int64
		roles, _ := s.userSvc.Perm.GetRolesForUser(cast.ToString(user))
		for _, role := range roles {
			permissions := s.userSvc.Perm.GetPermissionsForUser(role)
			for _, p := range permissions {
				resourceIds = append(resourceIds, cast.ToInt64(p[3]))
			}
		}
		err := s.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Where("id in ?", resourceIds).Order("sequence")
			}).
			Order("sequence").
			Where("parent_id = 0 and id in ?", resourceIds).
			Find(&resources, conds).Error
		if err != nil {
			return resources, err
		}
	} else {
		err := s.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("sequence")
			}).
			Order("sequence").
			Where("parent_id = 0").
			Find(&resources, conds).Error
		if err != nil {
			return resources, err
		}
	}

	return resources, nil
}

// 分页资源集合
func (s *ResourceSvc) PageResources(conds map[string]any) ([]mdl.Resource, int64, error) {
	var resources []mdl.Resource
	var count int64
	current := cast.ToInt(conds["current"])
	size := cast.ToInt(conds["size"])
	delete(conds, "current")
	delete(conds, "size")
	err := s.Db.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Limit(size).
		Offset(size*(current-1)).
		Order("sequence").
		Where("parent_id = 0").
		Find(&resources, conds).
		Count(&count).Error
	if err != nil {
		return resources, count, err
	}

	return resources, count, nil
}
