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
	routerGroup.GET("/logout", adminLogin.AdminLoginOut)
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
	// 管理员登录业务说明
	// 1.Username查询取得管理员信息
	// 2.adminInfo.salt + params.Password sha256 ---> saltPassword
	// 2.saltPassword == adminInfo.password
	admin := &dao.Admin{}
	// 数据库连接池的获取 default 取自配置文件 conf/dev/mysql_map.toml
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 1002, err)
	}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 1003, err)
	}
	// 保存登录信息到session
	saveAdminInfoToSession(admin, c)
	out := &dto.AdminLoginOutput{
		Token: params.UserName,
	}
	middleware.ResponseSuccess(c, out)
}

// saveAdminInfoToSession 保存登录信息到session
func saveAdminInfoToSession(admin *dao.Admin, c *gin.Context) {
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
}

// AdminLoginOut action
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept json
// @Produce json
// Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminLogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(public.AdminSessionInfoKey)
	session.Save()
	middleware.ResponseSuccess(c, "success")
}
