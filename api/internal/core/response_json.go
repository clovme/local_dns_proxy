package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"local_dns_proxy/pkg/constants"
	"local_dns_proxy/pkg/enums/code"
	"net/http"
)

type page struct {
	PageSize    int64 `json:"pageSize"`
	CurrentPage int64 `json:"currentPage"`
	Total       int64 `json:"total"`
}

// response 响应结构体
type response struct {
	Code    code.ResponseCode `json:"code"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
	Page    interface{}       `json:"page,omitempty"`
}

// Limit 设置分页信息
// 参数：
//   - pageSize: 每页数量
//   - currentPage: 当前页码
//   - total: 总记录数
//
// 返回值：
//   - *Context: 自定义gin.Context对象
func (r *Context) Limit(pageSize, currentPage, total int64) *Context {
	r.Context.Set(constants.LimitPage, &page{
		PageSize:    pageSize,
		CurrentPage: currentPage,
		Total:       total,
	})
	return r
}

// JsonSafe 安全响应
// 参数：
//   - c: gin.Context
//   - code: 响应码
//   - message: 响应消息
//   - data: 响应数据
//
// 返回值：
//   - 无
func (r *Context) JsonSafe(httpCode code.ResponseCode, message string, data interface{}) {
	r.Context.JSON(http.StatusOK, response{
		Code:    httpCode,
		Message: fmt.Sprintf("[%d] %s", httpCode.Int(), message),
		Data:    data,
		Page:    r.Get(constants.LimitPage),
	})
}

// JsonSafeDesc 安全响应
// 参数：
//   - c: gin.Context
//   - code: 响应码
//   - data: 响应数据
//
// 返回值：
//   - 无
func (r *Context) JsonSafeDesc(httpCode code.ResponseCode, data interface{}) {
	r.JsonSafe(httpCode, httpCode.Desc(), data)
}

// JsonSafeSuccess 安全响应成功
// 参数：
//   - c: gin.Context
//   - data: 响应数据
//
// 返回值：
//   - 无
func (r *Context) JsonSafeSuccess(data interface{}) {
	r.JsonSafe(code.Success, code.Success.Desc(), data)
}

func (r *Context) JsonSafeDnsStatus(first, running string) {
	r.JsonSafeSuccess(gin.H{
		"first":   first,
		"running": running,
	})
}

// JsonUnSafe 不安全响应
// 参数：
//   - c: gin.Context
//   - code: 响应码
//   - message: 响应消息
//   - data: 响应数据
//
// 返回值：
//   - 无
func (r *Context) JsonUnSafe(httpCode code.ResponseCode, message string, data interface{}) {
	r.Context.JSON(http.StatusOK, response{
		Code:    httpCode,
		Message: fmt.Sprintf("[%d] %s", httpCode.Int(), message),
		Data:    data,
		Page:    r.Get(constants.LimitPage),
	})
}

// JsonUnSafeDesc 不安全响应
// 参数：
//   - c: gin.Context
//   - code: 响应码
//   - data: 响应数据
//
// 返回值：
//   - 无
func (r *Context) JsonUnSafeDesc(httpCode code.ResponseCode, data interface{}) {
	r.JsonUnSafe(httpCode, httpCode.Desc(), data)
}

// JsonUnSafeSuccess 不安全响应成功
// 参数：
//   - c: gin.Context
//   - data: 响应数据
//
// 返回值：
//   - 无
func (r *Context) JsonUnSafeSuccess(data interface{}) {
	r.JsonUnSafe(code.Success, code.Success.Desc(), data)
}

func (r *Context) JsonUnSafeDnsStatus(first, running string) {
	r.JsonUnSafeSuccess(gin.H{
		"first":   first,
		"running": running,
	})
}
