package mgt

import (
	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @summary 获取角色
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} mdl.Role
// @router /mgt/sys/roles/{id} [get]
func GetRole(c *gin.Context) {
	role, err := svc.RoleSvc.Get(c.Param("id"))
	r.J(c, role, err)
}

// @summary 获取角色资源集合
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} []int
// @router /mgt/sys/roles/{id}/resources [get]
func ListRoleResources(c *gin.Context) {
	id := cast.ToInt64(c.Param("id"))
	resources, err := svc.RoleSvc.ListRoleResources(id)
	r.J(c, resources, err)
}

// @summary 分页角色集合
// @tags 系统
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.Role
// @router /mgt/sys/roles [get]
func PageRoles(c *gin.Context) {
	roles, count, err := svc.RoleSvc.Page(r.Q(c))
	r.J(c, roles, count, err)
}

// @summary 移除角色集合
// @tags 系统
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/sys/roles/{ids} [delete]
func RemoveRoles(c *gin.Context) {
	err := svc.RoleSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, err)
}

// @summary 保存角色
// @tags 系统
// @param role body mdl.Role true "角色"
// @success 200 {object} bool
// @router /mgt/sys/roles [post]
func SaveRole(c *gin.Context) {
	var role *mdl.Role
	err := c.ShouldBind(&role)
	if err == nil {
		err = svc.RoleSvc.Save(role)
	}
	r.J(c, role.Id, err)
}

// @summary 保存角色资源集合
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} bool
// @router /mgt/sys/roles/{id}/resources [post]
func SaveRoleResources(c *gin.Context) {
	id := cast.ToInt64(c.Param("id"))
	var resources []mdl.Resource
	err := c.ShouldBind(&resources)
	if err == nil {
		err = svc.RoleSvc.SaveRoleResources(id, resources)
	}
	r.J(c, true, err)
}
