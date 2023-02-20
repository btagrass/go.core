package mgt

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
)

// @summary 保存文件
// @tags 系统
// @param dir path string false "目录"
// @success 200 {object} string
// @router /mgt/sys/files/{dir} [post]
func SaveFile(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		c.Abort()
	}
	fileDir := filepath.Join("data/files", c.Param("dir"))
	err = utl.MakeDir(fileDir)
	if err != nil {
		c.Abort()
	}
	fileName := fmt.Sprintf("%s/%s%s", fileDir, utl.TimeId(), filepath.Ext(header.Filename))
	err = c.SaveUploadedFile(header, fileName)
	if err != nil {
		c.Abort()
	}

	c.String(http.StatusOK, fileName)
}
