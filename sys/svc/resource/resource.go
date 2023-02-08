package resource

import (
	"github.com/spf13/cast"
	"go.core/svc"
	"go.core/sys/mdl"
	"go.core/sys/svc/user"
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
func (r *ResourceSvc) ListMenus(userId string) ([]mdl.Resource, error) {
	var resources []mdl.Resource
	if userId == "300000000000001" {
		err := r.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("sequence")
			}).
			Find(&resources, "parent_id is null and type = 1").
			Order("sequence").Error
		if err != nil {
			return resources, err
		}
	} else {
		var ids []string
		roles, _ := r.userSvc.Perm.GetRolesForUser(cast.ToString(userId))
		for _, role := range roles {
			permissions := r.userSvc.Perm.GetPermissionsForUser(role)
			for _, p := range permissions {
				ids = append(ids, p[3])
			}
		}
		err := r.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Where("id in ?", ids).Order("sequence")
			}).
			Find(&resources, "parent_id is null and type = 1 and id in ?", ids).
			Order("sequence").Error
		if err != nil {
			return resources, err
		}
	}

	return resources, nil
}

// 获取资源集合
func (r *ResourceSvc) ListResources(conds map[string]any) ([]mdl.Resource, error) {
	var resources []mdl.Resource
	user, ok := conds["user"]
	if ok {
		delete(conds, "user")
		var resourceIds []int64
		roles, _ := r.userSvc.Perm.GetRolesForUser(cast.ToString(user))
		for _, role := range roles {
			permissions := r.userSvc.Perm.GetPermissionsForUser(role)
			for _, p := range permissions {
				resourceIds = append(resourceIds, cast.ToInt64(p[3]))
			}
		}
		err := r.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Where("id in ?", resourceIds).Order("sequence")
			}).
			Where("parent_id is null and id in ?", resourceIds).
			Find(&resources, conds).
			Order("sequence").Error
		if err != nil {
			return resources, err
		}
	} else {
		err := r.Db.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("sequence")
			}).
			Where("parent_id is null").
			Find(&resources, conds).
			Order("sequence").Error
		if err != nil {
			return resources, err
		}
	}

	return resources, nil
}

// 分页资源集合
func (r *ResourceSvc) PageResources(conds map[string]any) ([]*mdl.Resource, int64, error) {
	var resources []*mdl.Resource
	var count int64
	size := cast.ToInt(conds["size"])
	current := cast.ToInt(conds["current"])
	delete(conds, "size")
	delete(conds, "current")
	err := r.Db.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Where("parent_id = 0 or parent_id is null").
		Limit(size).
		Offset(size*(current-1)).
		Find(&resources, conds).
		Order("sequence").
		Count(&count).Error
	if err != nil {
		return resources, count, err
	}

	return resources, count, nil
}
