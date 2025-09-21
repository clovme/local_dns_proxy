package core

import (
	"local_dns_proxy/pkg/constants"
	"net/http"
)

// HTMLResponse 自定义HTML响应
type viewData struct {
	Data          any
	IsEnableEmail bool
	IsLogin       bool
	WebTitle      string
	ClientID      string
}

// HTML 加载HTML模板
//
// 参数:
//   - name: 模板名称
//   - title: 页面标题
//   - data: 页面数据
//
// 说明:
//   - 加载HTML模板，渲染页面数据。
func (r *Context) HTML(name string, data any) {
	r.Context.HTML(http.StatusOK, name, viewData{
		Data:     data,
		WebTitle: constants.WebTitle,
	})
}
