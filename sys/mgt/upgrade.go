package mgt

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.core/sys/svc"
)

// @summary 升级
// @tags 系统
// @param ver path string true "版本"
// @success 200 {object} []byte
// @router /mgt/upgrades/{ver} [get]
func Upgrade(c *gin.Context) {
	ver := c.Param("ver")
	filePath, fileVer, err := svc.UpgradeSvc.Upgrade(ver)
	if err != nil || filePath == "" {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Header("ver", fileVer)
		c.FileAttachment(filePath, filepath.Base(filePath))
	}
}
