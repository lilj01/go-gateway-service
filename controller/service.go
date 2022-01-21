package controller

import (
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/lilj_01/gin_gateway/dto"
	"github.com/lilj_01/gin_gateway/middleware"
	"github.com/lilj_01/gin_gateway/models"
	"github.com/lilj_01/gin_gateway/public"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.POST("/service_add_http", service.ServiceAddHTTP)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (*ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3000, err)
		return
	}
	serviceModel := &models.ServiceInfo{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 3001, err)
		return
	}
	list, count, err := serviceModel.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 3002, err)
		return
	}
	var outList []dto.ServiceListItemOutput
	outList, err = convert(list, tx, c)
	if err != nil {
		middleware.ResponseError(c, 3003, err)
		return
	}
	out := dto.ServiceListOutput{
		Total: count,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)
}

// convert
//1、http后缀接入 clusterIP+clusterPort+path
//2、http域名接入 domain
//3、tcp、grpc接入 clusterIP+servicePort
func convert(list []models.ServiceInfo, tx *gorm.DB, c *gin.Context) ([]dto.ServiceListItemOutput, error) {
	var outList []dto.ServiceListItemOutput
	clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
	clusterPort := lib.GetStringConf("base.cluster.cluster_port")
	clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			return nil, err
		}
		serviceAddr := "unKnow"
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			TotalNode:   len(ipList),
			Qpd:         0,
			Qps:         0,
		}
		outList = append(outList, outItem)
	}
	return outList, nil
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (*ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3005, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 3006, err)
		return
	}
	//读取基本信息
	serviceInfo := &models.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 3007, err)
		return
	}
	//软删除
	serviceInfo.IsDelete = public.UnValid
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 3008, err)
		return
	}
	middleware.ResponseSuccess(c, "success")
}

// ServiceAddHTTP godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (*ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3009, err)
		return
	}
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 3013, errors.New("IP列表与权重列表数量不一致"))
		return
	}
	saveServiceHttp(c, params)
}

//saveServiceHttp 添加HTTP服务
func saveServiceHttp(c *gin.Context, params *dto.ServiceAddHTTPInput) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 3010, err)
		return
	}
	tx = tx.Begin()
	serviceInfo := &models.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if serviceInfo != nil && strings.Compare(serviceInfo.ServiceName, params.ServiceName) == 0 {
		tx.Rollback()
		middleware.ResponseError(c, 3011, errors.New("服务已存在"))
		return
	}
	httpUrl := &models.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	httpUrl, err = httpUrl.Find(c, tx, httpUrl)
	if httpUrl != nil && httpUrl.RuleType == params.RuleType && httpUrl.Rule == params.Rule {
		tx.Rollback()
		middleware.ResponseError(c, 3012, errors.New("服务接入前缀或域名已存在"))
		return
	}
	serviceModel := &models.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
		IsDelete:    public.Valid,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 3013, err)
		return
	}
	httpRule := &models.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	accessControl := &models.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	loadBalance := &models.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
}
