package web

import (
	"github.com/gin-gonic/gin"
	dnsService "local_dns_proxy/internal/application/dns"
	"local_dns_proxy/internal/core"
	"local_dns_proxy/internal/schema/dto/dns"
	"local_dns_proxy/pkg/cfg"
	"local_dns_proxy/pkg/constants"
	"local_dns_proxy/pkg/copyright"
	"local_dns_proxy/pkg/enums/code"
	"local_dns_proxy/pkg/utils"
	"strings"
)

type DnsHandler struct {
	Service *dnsService.WebDnsService
}

// GetViewsIndexHandler
// @Type			web
// @Group 			dnsView
// @Router			/ [GET]
// @Name			indexView
// @Summary			首页视图
func (r *DnsHandler) GetViewsIndexHandler(c *core.Context) {
	c.HTML("index.html", nil)
}

// CopyrightHandler
// @Type			api
// @Group 			dnsApi
// @Router			/copyright [GET]
// @Name			copyright
// @Summary			版权
func (r *DnsHandler) CopyrightHandler(c *core.Context) {
	c.JsonUnSafeSuccess(copyright.NewCopyright())
}

// PageHandler
// @Type			api
// @Group 			dnsApi
// @Router			/list [GET]
// @Name			dnsList
// @Summary			获取DNS列表
func (r *DnsHandler) PageHandler(c *core.Context) {
	query, ok := c.GetQuery("orderBy")
	if !ok {
		query = "createdAt|asc"
	}
	dataList, err := r.Service.FindDnsList(query)
	if err != nil {
		c.JsonSafeDesc(code.Fail, err)
		return
	}
	c.JsonSafeSuccess(dataList)
}

// SaveHandler
// @Type			api
// @Group 			dnsApi
// @Router			/save [POST]
// @Name			dnsSave
// @Summary			保存DNS数据
func (r *DnsHandler) SaveHandler(c *core.Context) {
	var data dns.WebSaveData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JsonSafeDesc(code.ServiceInsertError, err)
		return
	}
	c.JsonUnSafeSuccess(r.Service.SaveDnsData(data))
}

// DeleteHandler
// @Type			api
// @Group 			dnsApi
// @Router			/delete [DELETE]
// @Name			dnsDelete
// @Summary			删除DNS数据
func (r *DnsHandler) DeleteHandler(c *core.Context) {
	var row dns.WebDeleteRow
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JsonSafeDesc(code.ServiceDeleteError, err)
		return
	}
	if err := r.Service.DeleteDnsData(row.ID); err != nil {
		c.JsonSafeDesc(code.ServiceDeleteError, err)
		return
	}
	c.JsonSafe(code.Success, "数据删除成功", nil)
}

// ServiceRunningHandler
// @Type			api
// @Group 			dnsApi
// @Router			/service/running/{first}/{iface} [POST]
// @Name			dnsServiceRunning
// @Summary			启动DNS服务
func (r *DnsHandler) ServiceRunningHandler(c *core.Context) {
	if strings.EqualFold(c.Param("first"), "first") {
		c.JsonSafeDnsStatus("first", cfg.C.Server.DNSRunning)
		return
	}

	cfg.C.Server.Iface = c.Param("iface")
	cfg.C.Server.DNSRunning = constants.DNSRunning

	dnsProxy, err := r.Service.NewDNSProxy()
	if err != nil {
		c.JsonSafe(code.RequestNotFound, err.Error(), nil)
		return
	}
	if err := dnsProxy.SetLocalDNS(cfg.C.Server.Iface); err != nil {
		c.JsonSafe(code.RequestNotFound, "设置系统DNS(127.0.0.1)失败", err)
		return
	}

	go dnsProxy.StartDnsServer()

	if !dnsProxy.IsRunning {
		c.JsonSafeDnsStatus("third", constants.DNSStop)
		return
	}

	<-dnsProxy.CtxStatus.Done()

	c.JsonSafeDnsStatus("third", constants.DNSRunning)
	cfg.SaveToIni()
}

// ServiceStopHandler
// @Type			api
// @Group 			dnsApi
// @Router			/service/stop/{first}/{iface} [POST]
// @Name			dnsServiceStop
// @Summary			禁用DNS服务
func (r *DnsHandler) ServiceStopHandler(c *core.Context) {
	if strings.EqualFold(c.Param("first"), "first") {
		c.JsonSafeDnsStatus("first", cfg.C.Server.DNSRunning)
		return
	}

	cfg.C.Server.DNSRunning = constants.DNSStop
	cfg.C.Server.Iface = c.Param("iface")

	dnsProxy, err := r.Service.NewDNSProxy()
	if err != nil {
		c.JsonSafe(code.RequestNotFound, "获取网卡列表失败", err)
		return
	}
	if err := dnsProxy.RestoreDNS(cfg.C.Server.Iface); err != nil {
		c.JsonSafe(code.RequestNotFound, "系统代理DNS设置失败", err)
		return
	}

	c.JsonSafeDnsStatus("third", constants.DNSStop)
	cfg.SaveToIni()
}

// GetNetIfaceHandler
// @Type			api
// @Group 			dnsApi
// @Router			/network/interfaces [GET]
// @Name			dnsNetIface
// @Summary			获取网络接口列表
func (r *DnsHandler) GetNetIfaceHandler(c *core.Context) {
	ifaces, err := utils.GetNetworkInterfaces()
	if err != nil {
		c.JsonSafe(code.RequestNotFound, "获取网卡列表失败", err)
		return
	}
	c.JsonSafeSuccess(gin.H{
		"iface":   cfg.C.Server.Iface,
		"running": cfg.C.Server.DNSRunning,
		"ifaces":  ifaces,
	})
}
