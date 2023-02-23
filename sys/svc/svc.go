package svc

import (
	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc/dept"
	"github.com/btagrass/go.core/sys/svc/dict"
	"github.com/btagrass/go.core/sys/svc/resource"
	"github.com/btagrass/go.core/sys/svc/role"
	"github.com/btagrass/go.core/sys/svc/upgrade"
	"github.com/btagrass/go.core/sys/svc/user"
	"github.com/sirupsen/logrus"
)

var (
	DeptSvc     *dept.DeptSvc         // 部门服务
	DictSvc     *dict.DictSvc         // 字典服务
	ResourceSvc *resource.ResourceSvc // 资源服务
	RoleSvc     *role.RoleSvc         // 角色服务
	UpgradeSvc  *upgrade.UpgradeSvc   // 升级服务
	UserSvc     *user.UserSvc         // 用户服务
)

// 初始化
func init() {
	// 迁移
	err := svc.Migrate(
		[]any{
			&mdl.Dept{},
			&mdl.Dict{},
			&mdl.Resource{},
			&mdl.Role{},
			&mdl.User{},
		},
		"INSERT INTO sys_dept VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 0, '公司', '', '', 1)",
		"INSERT INTO sys_dept VALUES (300000000000101, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '开发部', '', '', 1)",
		"INSERT INTO sys_dept VALUES (300000000000102, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '财务部', '', '', 2)",
		"INSERT INTO sys_dict VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 'Resource', 1, '菜单', 1)",
		"INSERT INTO sys_dict VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, 'Resource', 2, '权限', 2)",
		"INSERT INTO sys_dict VALUES (300000000000003, '2023-01-29 00:00:00.000', NULL, NULL, 'Act', 1, 'GET', 1)",
		"INSERT INTO sys_dict VALUES (300000000000004, '2023-01-29 00:00:00.000', NULL, NULL, 'Act', 2, 'DELETE', 2)",
		"INSERT INTO sys_dict VALUES (300000000000005, '2023-01-29 00:00:00.000', NULL, NULL, 'Act', 3, 'POST', 3)",
		"INSERT INTO sys_resource VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 0, '系统设置', 1, 'Setting', '/sys', NULL, 100)",
		"INSERT INTO sys_resource VALUES (300000000000101, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '资源管理', 1, 'Menu', '/sys/resources', NULL, 1)",
		"INSERT INTO sys_resource VALUES (300000000000102, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '角色管理', 1, 'Avatar', '/sys/roles', NULL, 2)",
		"INSERT INTO sys_resource VALUES (300000000000103, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '用户管理', 1, 'User', '/sys/users', NULL, 3)",
		"INSERT INTO sys_resource VALUES (300000000000104, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '部门管理', 1, 'OfficeBuilding', '/sys/depts', NULL, 4)",
		"INSERT INTO sys_resource VALUES (300000000000105, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '字典管理', 1, 'Files', '/sys/dicts', NULL, 5)",
		"INSERT INTO sys_role VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, '管理员')",
		"INSERT INTO sys_role VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, '测试')",
		"INSERT INTO sys_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ('g', '300000000000002', '300000000000002', '', '', '', '')",
		"INSERT INTO sys_user VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'admin', NULL, '15800000000', '$2a$10$enX7NxYTZZo9yLJQN6jXF.B6FGg7d9Q5eTW5off94hJZSa5AO9av2', 0)",
		"INSERT INTO sys_user VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'test', NULL, '15800000000', '$2a$10$jM3N9CtuC0hd7EH1ybsoAe5U.znR0dvoLU8KXGSTiCzsYiZyCsGGi', 0)",
	)
	if err != nil {
		logrus.Fatal(err)
	}
	// 服务
	DeptSvc = dept.NewDeptSvc()
	DictSvc = dict.NewDictSvc()
	UserSvc = user.NewUserSvc()
	ResourceSvc = resource.NewResourceSvc(UserSvc)
	RoleSvc = role.NewRoleSvc(UserSvc)
	UpgradeSvc = upgrade.NewUpgradeSvc()
}
