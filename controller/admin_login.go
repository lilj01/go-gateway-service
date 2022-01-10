package controller

import (
	"encoding/json"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dao"
	"github.com/lilj_01/gin_gateway/dto"
	"github.com/lilj_01/gin_gateway/middleware"
	"github.com/lilj_01/gin_gateway/public"
	"time"
)

type AdminLoginController struct {
}

//AdminLoginRegister 注册login-controller
func AdminLoginRegister(routerGroup *gin.RouterGroup) {
	adminLogin := AdminLoginController{}
	routerGroup.POST("/login", adminLogin.AdminLogin)
}

// AdminLogin action
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	// 参数错误处理
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
	}
	// Username查询取得管理员信息
	// adminInfo.salt + params.Password sha256 ---> saltPassword
	// saltPassword == adminInfo.password
	admin := &dao.Admin{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 1002, err)
	}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 1003, err)
	}

	// 保存登录信息到session
	sessionInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		Username:  admin.Username,
		LoginTime: time.Now(),
	}
	sessBytes, err := json.Marshal(sessionInfo)
	if err != nil {
		middleware.ResponseError(c, 1004, err)
	}
	session := sessions.Default(c)
	session.Set(public.AdminSessionInfoKey, string(sessBytes))
	session.Save()
	out := &dto.AdminLoginOutput{
		Token: params.UserName,
	}
	middleware.ResponseSuccess(c, out)
}
