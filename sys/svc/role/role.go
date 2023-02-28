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
func (s *RoleSvc) ListRoleResources(id string) ([]int64, error) {
	resources := make([]int64, 0)
	permissions := s.userSvc.Perm.GetPermissionsForUser(id)
	for _, p := range permissions {
		resources = append(resources, cast.ToInt64(p[3]))
	}

	return resources, nil
}

// 保存角色资源集合
func (s *RoleSvc) SaveRoleResources(id string, resources []mdl.Resource) error {
	_, err := s.userSvc.Perm.DeletePermissionsForUser(id)
	if err != nil {
		return err
	}
	var rs [][]string
	for _, r := range resources {
		rs = append(rs, []string{
			r.Uri,
			r.Act,
			cast.ToString(r.Id),
		})
	}
	_, err = s.userSvc.Perm.AddPermissionsForUser(id, rs...)
	if err != nil {
		return err
	}

	return nil
}
