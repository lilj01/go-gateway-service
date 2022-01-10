package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/public"
)

// AdminLoginInput 管理员登陆输入结构体
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,is_valid_username"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

//BindValidParam 结构体校验方法
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
