package models

import (
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dto"
	"github.com/lilj_01/gin_gateway/public"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	Id       int       `json:"id" gorm:"column:id"`
	Username string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt     string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password string    `json:"password" gorm:"column:password" description:"密码"`
	CreateAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	UpdateAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

// Find 查询
func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	query := tx.WithContext(c).Where(search)
	err := query.Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginCheck 登录检查
func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := t.Find(c, tx, &Admin{Username: param.UserName, IsDelete: 1})
	if err != nil {
		return nil, errors.New("管理员信息不存在")
	}
	saltPassword := public.GenSaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("password error")
	}
	return adminInfo, nil
}

// Save 保存
func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}
