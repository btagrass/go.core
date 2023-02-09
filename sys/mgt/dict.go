package mgt

import (
	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
)

// @summary 获取字典
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} mdl.Dict
// @router /mgt/sys/dicts/{id} [get]
func GetDict(c *gin.Context) {
	dict, err := svc.DictSvc.Get(c.Param("id"))
	r.J(c, dict, err)
}

// @summary 分页字典集合
// @tags 系统
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.Dict
// @router /mgt/sys/dicts [get]
func PageDicts(c *gin.Context) {
	dicts, count, err := svc.DictSvc.Page(r.Q(c))
	r.J(c, dicts, count, err)
}

// @summary 移除字典集合
// @tags 系统
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/sys/dicts/{ids} [delete]
func RemoveDicts(c *gin.Context) {
	err := svc.DictSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, err)
}

// @summary 保存字典
// @tags 系统
// @param dict body mdl.Dict true "字典"
// @success 200 {object} bool
// @router /mgt/sys/dicts [post]
func SaveDict(c *gin.Context) {
	var dict mdl.Dict
	err := c.ShouldBind(&dict)
	if err == nil {
		err = svc.DictSvc.Save(dict)
	}
	r.J(c, dict.Id, err)
}
