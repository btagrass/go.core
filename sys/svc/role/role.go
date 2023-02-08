package role

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc/user"
	"github.com/spf13/cast"
)

// 角色服务
type RoleSvc struct {
	*svc.Svc[mdl.Role]
	userSvc *user.UserSvc
}

// 构造函数
func NewRoleSvc(userSvc *user.UserSvc) *RoleSvc {
	return &RoleSvc{
		Svc:     svc.NewSvc[mdl.Role]("sys:roles"),
		userSvc: userSvc,
	}
}

// 获取角色资源集合
func (r *RoleSvc) ListRoleResources(id int64) ([]int64, error) {
	var resources []int64
	permissions := r.userSvc.Perm.GetPermissionsForUser(cast.ToString(id))
	for _, p := range permissions {
		resources = append(resources, cast.ToInt64(p[3]))
	}

	return resources, nil
}

// 保存角色资源集合
func (r *RoleSvc) SaveRoleResources(id int64, resources []mdl.Resource) error {
	_, err := r.userSvc.Perm.DeletePermissionsForUser(cast.ToString(id))
	if err != nil {
		return err
	}
	var rs [][]string
	for _, res := range resources {
		rs = append(rs, []string{
			res.Url,
			res.Act,
			cast.ToString(res.Id),
			cast.ToString(res.Type),
		})
	}
	_, err = r.userSvc.Perm.AddPermissionsForUser(cast.ToString(id), rs...)
	if err != nil {
		return err
	}

	return nil
}
