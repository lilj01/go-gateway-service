package models

import (
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dto"
	"gorm.io/gorm"
	"time"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := t.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}
	httpRule := &HttpRule{ServiceID: search.ID}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tcpRule := &TcpRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	grpcRule := &GrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

// Find 查询
func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.WithContext(c).Where(search).Where("is_delete = 0").Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PageList 分页查询
func (t *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, params *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	var count int64 = 0
	var list []ServiceInfo
	// 偏移量计算
	offset := (params.PageNo - 1) * params.PageSize
	query := tx.WithContext(c)
	// 设置TableName是因为count时需要
	query = query.Table(t.TableName()).Where("is_delete = ?", 0)
	if params.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
	}
	err := query.Limit(params.PageSize).Offset(offset).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Table(t.TableName()).Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
