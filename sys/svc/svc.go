package svc

import (
	"github.com/sirupsen/logrus"
	"go.core/svc"
	"go.core/sys/mdl"
	"go.core/sys/svc/dept"
	"go.core/sys/svc/dict"
	"go.core/sys/svc/resource"
	"go.core/sys/svc/role"
	"go.core/sys/svc/upgrade"
	"go.core/sys/svc/user"
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
		"INSERT INTO `sys_dept` VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, NULL, '公司', '', '', 1)",
		"INSERT INTO `sys_dept` VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '开发部', '', '', 1)",
		"INSERT INTO `sys_dept` VALUES (300000000000003, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, '财务部', '', '', 2)",
		"INSERT INTO `sys_dict` VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 'Resource', 1, '菜单', 1)",
		"INSERT INTO `sys_dict` VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, 'Resource', 2, '权限', 2)",
		"INSERT INTO `sys_resource` VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, NULL, 'sys', '系统设置', 1, 'setting', '/sys', NULL, 100)",
		"INSERT INTO `sys_resource` VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'sysResource', '资源管理', 1, 'menu', '/sys/resources', NULL, 1)",
		"INSERT INTO `sys_resource` VALUES (300000000000003, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'sysRole', '角色管理', 1, 'avatar', '/sys/roles', NULL, 2)",
		"INSERT INTO `sys_resource` VALUES (300000000000004, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'sysUser', '用户管理', 1, 'user', '/sys/users', NULL, 3)",
		"INSERT INTO `sys_resource` VALUES (300000000000005, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'sysDept', '部门管理', 1, 'office-building', '/sys/depts', NULL, 4)",
		"INSERT INTO `sys_resource` VALUES (300000000000006, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'sysDict', '字典管理', 1, 'files', '/sys/dicts', NULL, 5)",
		"INSERT INTO `sys_role` VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, '管理员')",
		"INSERT INTO `sys_user` VALUES (300000000000001, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000001, 'admin', NULL, '12345678900', '$2a$10$enX7NxYTZZo9yLJQN6jXF.B6FGg7d9Q5eTW5off94hJZSa5AO9av2', 0)",
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
