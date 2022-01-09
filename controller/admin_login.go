package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dto"
	"github.com/lilj_01/gin_gateway/middleware"
)

type AdminLoginController struct {
}

//AdminLoginRegister 注册login-controller
func AdminLoginRegister(routerGroup *gin.RouterGroup) {
	adminLogin := AdminLoginController{}
	routerGroup.POST("/login", adminLogin.AdminLogin)
}

// AdminLogin action
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	// 参数错误处理
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
	}
	data := struct {
		Message string `json:"message"`
	}{
		Message: "success",
	}
	middleware.ResponseSuccess(c, data)
}
