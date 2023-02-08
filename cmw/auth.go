package cmw

import (
	"fmt"

	"github.com/btagrass/go.core/r"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

// 认证
func Auth(perm *casbin.SyncedEnforcer, signedKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			r.J(c, fmt.Errorf("请携带有效的访问令牌"))
			return
		}
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			return signedKey, nil
		})
		if err != nil || !token.Valid {
			r.J(c, fmt.Errorf("请携带有效的访问令牌"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			r.J(c, fmt.Errorf("请携带有效的访问令牌"))
			return
		}
		userId, ok := claims["userId"]
		if !ok {
			r.J(c, fmt.Errorf("请携带有效的访问令牌"))
			return
		}
		ok, err = perm.Enforce(cast.ToString(userId), c.Request.URL.Path, c.Request.Method)
		if err != nil || !ok {
			r.J(c, fmt.Errorf("没有访问 %s 的权限", c.Request.URL.Path))
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
