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

// @summary 获取部门集合
// @tags 系统
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.Dept
// @router /mgt/sys/depts [get]
func ListDepts(c *gin.Context) {
	depts, count, err := svc.DeptSvc.ListDepts(r.Q(c))
	r.J(c, depts, count, err)
}

// @summary 移除部门集合
// @tags 系统
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/sys/depts/{ids} [delete]
func RemoveDepts(c *gin.Context) {
	err := svc.DeptSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

// @summary 保存部门
// @tags 系统
// @param dept body mdl.Dept true "部门"
// @success 200 {object} int
// @router /mgt/sys/depts [post]
func SaveDept(c *gin.Context) {
	var dept mdl.Dept
	err := c.ShouldBind(&dept)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.DeptSvc.Save(dept)
	r.J(c, dept.GetId(), err)
}
