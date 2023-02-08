package r

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// 分页参数
type PageP struct {
	Current int `form:"current" json:"current" binding:"required"` // 当前页
	Size    int `form:"size" json:"size" binding:"required"`       // 页大小
}

// 结果
type R struct {
	Code any    `json:"code"` // 代码
	Data any    `json:"data"` // 数据
	Msg  string `json:"msg"`  // 消息
}

// 查询
func Q(c *gin.Context) map[string]any {
	params := make(map[string]any)
	for k := range c.Request.URL.Query() {
		params[k] = c.Query(k)
	}

	return params
}

// Json
func J(c *gin.Context, data ...any) {
	r := R{}
	err, ok := data[len(data)-1].(error)
	if ok {
		r.Code = http.StatusInternalServerError
		r.Msg = err.Error()
	} else {
		r.Code = http.StatusOK
	}
	ds := data[:len(data)-1]
	if len(ds) == 1 {
		r.Data = ds[0]
	} else if len(ds) == 2 {
		count := cast.ToInt64(ds[1])
		if count == 0 {
			r.Data = ds[0]
		} else {
			r.Data = map[string]any{
				"records": ds[0],
				"total":   ds[1],
			}
		}
	}
	c.JSON(http.StatusOK, r)
	c.Abort()
}
