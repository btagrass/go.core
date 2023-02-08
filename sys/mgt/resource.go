package mgt

import (
	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @summary 获取资源
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} mdl.Resource
// @router /mgt/sys/resources/{id} [get]
func GetResource(c *gin.Context) {
	resource, err := svc.ResourceSvc.Get(c.Param("id"))
	r.J(c, resource, err)
}

// @summary 获取菜单集合
// @tags 系统
// @success 200 {object} []mdl.Resource
// @router /mgt/sys/resources/menu [get]
func ListMenus(c *gin.Context) {
	userId := cast.ToString(c.GetFloat64("userId"))
	menus, err := svc.ResourceSvc.ListMenus(userId)
	r.J(c, menus, err)
}

// @summary 获取资源集合
// @tags 系统
// @success 200 {object} []mdl.Resource
// @router /mgt/sys/resources [get]
func ListResources(c *gin.Context) {
	conds := r.Q(c)
	if _, ok := conds["user"]; ok {
		conds["user"] = c.GetInt64("id")
	}
	resources, err := svc.ResourceSvc.ListResources(conds)
	r.J(c, resources, err)
}

// @summary 分页资源集合
// @tags 系统
// @success 200 {object} []mdl.Resource
// @router /mgt/sys/resources [get]
func PageResources(c *gin.Context) {
	resources, count, err := svc.ResourceSvc.PageResources(r.Q(c))
	r.J(c, resources, count, err)
}

// @summary 移除资源集合
// @tags 系统
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/sys/resources/{ids} [delete]
func RemoveResources(c *gin.Context) {
	err := svc.ResourceSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, err)
}

// @summary 保存资源
// @tags 系统
// @param resource body mdl.Resource true "资源"
// @success 200 {object} bool
// @router /mgt/sys/resources [post]
func SaveResource(c *gin.Context) {
	var resource *mdl.Resource
	err := c.ShouldBind(&resource)
	if err == nil {
		err = svc.ResourceSvc.Save(resource)
	}
	r.J(c, resource.Id, err)
}
