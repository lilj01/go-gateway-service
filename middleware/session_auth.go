package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/public"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		adminInfo, ok := session.Get(public.AdminSessionInfoKey).(string)
		fmt.Printf("session adminInfo = %s\n", adminInfo)
		if !ok || adminInfo == "" {
			ResponseError(c, InternalErrorCode, errors.New("user not login"))
			c.Abort()
			return
		}
		c.Next()
	}
}
