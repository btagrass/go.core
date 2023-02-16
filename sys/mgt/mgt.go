package mgt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/btagrass/go.core/cmw"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// 管理
func Mgt() *gin.Engine {
	e := gin.Default()
	// 跨域
	e.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))
	// 调试
	d := e.Group("/debug")
	{
		d.GET("/cache", func(c *gin.Context) {
			keyword := c.Query("k")
			i := 0
			for k, v := range svc.UserSvc.Cache.Items() {
				if strings.Contains(k, keyword) {
					data, _ := json.Marshal(v)
					c.String(http.StatusOK, fmt.Sprintf("%d. %s    %+v\r\n", i+1, k, string(data)))
					i++
				}
			}
		})
		d.GET("/redis", func(c *gin.Context) {
			keyword := fmt.Sprintf("*%s*", c.Query("k"))
			keys := svc.UserSvc.Redis.Keys(context.Background(), keyword).Val()
			for i, k := range keys {
				var v any
				t := svc.UserSvc.Redis.Type(context.Background(), k).Val()
				if t == "hash" {
					v = svc.UserSvc.Redis.HGetAll(context.Background(), k).Val()
				}
				c.String(http.StatusOK, fmt.Sprintf("%d. %s    %+v\r\n", i+1, k, v))
			}
		})
	}
	pprof.Register(e)
	// 文档 (http://ip:port+1/swagger/index.html)
	e.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
		c.InstanceName = "mgt"
		c.Title = viper.GetString("app.name")
	}))
	// 管理
	m := e.Group("/mgt")
	{
		// 登录
		m.POST("/login", Login)
		// 升级
		m.GET("/upgrades/:ver", Upgrade)
	}
	// 系统
	s := m.Group("/sys").Use(cmw.Auth(svc.UserSvc.Perm, svc.UserSvc.SignedKey))
	{
		// 字典
		s.GET("/dicts/:id", GetDict)
		s.GET("/dicts", ListDicts)
		s.DELETE("/dicts/:ids", RemoveDicts)
		s.POST("/dicts", SaveDict)
		// 部门
		s.GET("/depts/:id", GetDept)
		s.GET("/depts", ListDepts)
		s.DELETE("/depts/:ids", RemoveDepts)
		s.POST("/depts", SaveDept)
		// 资源
		s.GET("/resources/:id", func(c *gin.Context) {
			if c.Param("id") == "menu" {
				ListMenus(c)
			} else {
				GetResource(c)
			}
		})
		s.GET("/resources", ListResources)
		s.DELETE("/resources/:ids", RemoveResources)
		s.POST("/resources", SaveResource)
		// 角色
		s.GET("/roles/:id", GetRole)
		s.GET("/roles", ListRoles)
		s.DELETE("/roles/:ids", RemoveRoles)
		s.POST("/roles", SaveRole)
		s.GET("/roles/:id/resources", ListRoleResources)
		s.POST("/roles/:id/resources", SaveRoleResources)
		// 用户
		s.GET("/users/:id", GetUser)
		s.GET("/users", ListUsers)
		s.DELETE("/users/:ids", RemoveUsers)
		s.POST("/users", SaveUser)
		s.GET("/users/:id/roles", ListUserRoles)
		s.POST("/users/:id/roles", SaveUserRoles)
	}

	return e
}
