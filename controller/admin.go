package controller

import (
	"encoding/json"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dto"
	"github.com/lilj_01/gin_gateway/middleware"
	"github.com/lilj_01/gin_gateway/models"
	"github.com/lilj_01/gin_gateway/public"
)

type AdminController struct {
}

func AdminRegister(routerGroup *gin.RouterGroup) {
	adminLogin := AdminController{}
	routerGroup.GET("/admin_info", adminLogin.AdminInfo)
	routerGroup.POST("/change_pwd", adminLogin.ChangePwd)
}

// AdminInfo godoc
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
		return
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

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept json
// @Produce json
// @Param body body dto.ChangePwdInput true "body"
// Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminLogin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	// 参数错误处理
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//1.session读取用户信息到结构体
	//2.sessionInfo.id 读取数据库信息 adminInfo
	//3.params.password + admin.salt sha256 saltPassword
	//4.saltPassword => save adminInfo.password 执行数据保存
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	adminModel := &models.Admin{}
	adminModel, err = adminModel.Find(c, tx, &models.Admin{
		Username: adminSessionInfo.Username,
		Id:       adminSessionInfo.ID,
		IsDelete: 1,
	})
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	saltPwd := public.GenSaltPassword(adminModel.Salt, params.Password)
	adminModel.Password = saltPwd
	if err = adminModel.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2005, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
}
