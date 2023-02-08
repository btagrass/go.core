package mgt

import (
	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
)

// @summary 获取部门
// @tags 系统
// @param id path int true "编码"
// @success 200 {object} mdl.Dept
// @router /mgt/sys/depts/{id} [get]
func GetDept(c *gin.Context) {
	dept, err := svc.DeptSvc.Get(c.Param("id"))
	r.J(c, dept, err)
}

// @summary 分页部门集合
// @tags 系统
// @success 200 {object} []mdl.Dept
// @router /mgt/sys/depts [get]
func PageDepts(c *gin.Context) {
	depts, count, err := svc.DeptSvc.PageDepts(r.Q(c))
	r.J(c, depts, count, err)
}

// @summary 移除部门集合
// @tags 系统
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/sys/depts/{ids} [delete]
func RemoveDepts(c *gin.Context) {
	err := svc.DeptSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, err)
}

// @summary 保存部门
// @tags 系统
// @param dept body mdl.Dept true "部门"
// @success 200 {object} bool
// @router /mgt/sys/depts [post]
func SaveDept(c *gin.Context) {
	var dept *mdl.Dept
	err := c.ShouldBind(&dept)
	if err == nil {
		err = svc.DeptSvc.Save(dept)
	}
	r.J(c, dept.Id, err)
}
