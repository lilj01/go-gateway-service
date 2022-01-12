package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dto"
	"github.com/lilj_01/gin_gateway/middleware"
	"github.com/lilj_01/gin_gateway/public"
)

type AdminController struct {
}

func AdminRegister(routerGroup *gin.RouterGroup) {
	adminLogin := AdminController{}
	routerGroup.GET("/admin_info", adminLogin.AdminInfo)
}

// AdminInfo action
// @Summary 登录信息
// @Description 登录信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept json
// @Produce json
// Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminLogin *AdminController) AdminInfo(c *gin.Context) {
	//1、读取sessionKey对应的json，转换为结构体
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
	}
	//2、返回结构体
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.Username,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://bclz_xc.gitee.io/lilj_01-static/lilj/tx.jpg",
		Introduction: "super man",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}
