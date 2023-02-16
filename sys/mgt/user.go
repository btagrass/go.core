package mgt

import (
	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
)

// @summary 获取用户
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} mdl.User
// @router /mgt/sys/users/{id} [get]
func GetUser(c *gin.Context) {
	user, err := svc.UserSvc.Get(c.Param("id"))
	r.J(c, user, err)
}

// @summary 获取取用户集合
// @tags 系统
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.User
// @router /mgt/sys/users [get]
func ListUsers(c *gin.Context) {
	users, count, err := svc.UserSvc.ListUsers(r.Q(c))
	r.J(c, users, count, err)
}

// @summary 移除用户集合
// @tags 系统
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/sys/users/{ids} [delete]
func RemoveUsers(c *gin.Context) {
	err := svc.UserSvc.RemoveUsers(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

// @summary 保存用户
// @tags 系统
// @param user body mdl.User true "用户"
// @success 200 {object} int
// @router /mgt/sys/users [post]
func SaveUser(c *gin.Context) {
	var user mdl.User
	err := c.ShouldBind(&user)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.UserSvc.SaveUser(user)
	r.J(c, user.GetId(), err)
}

// @summary 获取用户角色集合
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} []int
// @router /mgt/sys/users/{id}/roles [get]
func ListUserRoles(c *gin.Context) {
	roles, err := svc.UserSvc.ListUserRoles(c.Param("id"))
	r.J(c, roles, err)
}

// @summary 保存用户角色集合
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} bool
// @router /mgt/sys/users/{id}/roles [post]
func SaveUserRoles(c *gin.Context) {
	var roles []int64
	err := c.ShouldBind(&roles)
	if err == nil {
		err = svc.UserSvc.SaveUserRoles(c.Param("id"), roles)
	}
	r.J(c, true, err)
}

// @summary 登录
// @tags 系统
// @param userName formData string true "用户名"
// @param password formData string true "密码"
// @success 200 {object} mdl.User
// @router /mgt/login [post]
func Login(c *gin.Context) {
	var user *mdl.User
	var p struct {
		UserName string `form:"userName" json:"userName" binding:"required"` // 用户名
		Password string `form:"password" json:"password" binding:"required"` // 密码
	}
	err := c.ShouldBind(&p)
	if err == nil {
		user, err = svc.UserSvc.Login(p.UserName, p.Password)
	}
	r.J(c, user, err)
}
