package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/public"
	"time"
)

// AdminLoginInput 管理员登陆输入结构体
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,is_valid_username"` //管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   //密码
}

//BindValidParam 结构体校验方法
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

// AdminLoginOutput 管理员登陆反参
type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token"` //token
}

// AdminSessionInfo 登录保存session结构体
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}
