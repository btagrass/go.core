package mgt

import (
	"net/http"
	"path/filepath"

	"github.com/btagrass/go.core/sys/svc"
	"github.com/gin-gonic/gin"
)

// @summary 升级
// @tags 系统
// @param ver path string true "版本"
// @success 200 {object} []byte
// @router /mgt/upgrades/{ver} [get]
func Upgrade(c *gin.Context) {
	filePath, fileVer, err := svc.UpgradeSvc.Upgrade(c.Param("ver"))
	if err != nil || filePath == "" {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Header("ver", fileVer)
		c.FileAttachment(filePath, filepath.Base(filePath))
	}
}
