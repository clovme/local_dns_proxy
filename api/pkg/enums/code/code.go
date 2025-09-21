package code

import (
	"local_dns_proxy/pkg/enums"
	"sort"
)

type ResponseCode int

const Name = "http_status_code"

const (
	// 正常返回
	Success ResponseCode = iota + 10000
	Fail

	// 业务错误
	ServiceVerifyError ResponseCode = iota + 19998
	ServiceInsertError
	ServiceDeleteError
	ServiceUpdateError

	// 请求错误
	RequestBadRequest ResponseCode = iota + 39984
	RequestUnauthorized
	RequestForbidden
	RequestNotFound
	RequestUnknown

	// 服务器内部错误
	ServerInternalError ResponseCode = iota + 49989
)

var (
	initiate = map[ResponseCode]enums.Enums{
		// 正常返回
		Success: {Key: "Success", Name: "成功", Desc: "请求已成功处理！"},
		Fail:    {Key: "Fail", Name: "操作失败", Desc: "操作失败，请稍后重试！"},

		// 业务错误
		ServiceVerifyError: {Key: "ServiceVerifyError", Name: "验证失败", Desc: "数据验证失败，请检查输入数据！"},
		ServiceInsertError: {Key: "ServiceInsertError", Name: "创建失败", Desc: "数据创建失败，请重试！"},
		ServiceDeleteError: {Key: "ServiceDeleteError", Name: "删除失败", Desc: "数据删除失败，请重试！"},
		ServiceUpdateError: {Key: "ServiceUpdateError", Name: "更新失败", Desc: "数据更新失败，请重试！"},

		// 请求错误
		RequestBadRequest:   {Key: "RequestBadRequest", Name: "错误请求", Desc: "请求参数格式错误或缺失，服务器无法处理！"},
		RequestUnauthorized: {Key: "RequestUnauthorized", Name: "未认证", Desc: "当前请求未认证或者您没有权限访问！"},
		RequestForbidden:    {Key: "RequestForbidden", Name: "禁止访问", Desc: "拒绝访问此资源！"},
		RequestNotFound:     {Key: "RequestNotFound", Name: "资源不存在", Desc: "请求的资源不存在或已被删除！"},
		RequestUnknown:      {Key: "RequestUnknown", Name: "未知错误", Desc: "未知错误或异常，请检查请求参数或配置！"},

		// 服务器内部错误
		ServerInternalError: {Key: "ServerInternalError", Name: "服务器内部错误", Desc: "服务器开小差了，请稍后再试！"},
	}

	enumToValue = make(map[string]ResponseCode)
)

func init() {
	for code, meta := range initiate {
		enumToValue[meta.Key] = code
	}
}

// Key 获取enums.Key
func (c ResponseCode) Key() string {
	if meta, ok := initiate[c]; ok {
		return meta.Key
	}
	return "Unknown"
}

// Name 获取枚举名称
func (c ResponseCode) Name() string {
	if meta, ok := initiate[c]; ok {
		return meta.Name
	}
	return "未知错误"
}

// Desc 获取枚举描述
func (c ResponseCode) Desc() string {
	if meta, ok := initiate[c]; ok {
		return meta.Desc
	}
	return "未知错误或异常，请检查请求参数或联系管理员"
}

// Int 获取枚举值
func (c ResponseCode) Int() int {
	return int(c)
}

// Is 比较枚举值
func (c ResponseCode) Is(v ResponseCode) bool {
	return v == c
}

// Code 获取Code
func Code(key string) ResponseCode {
	if enum, ok := enumToValue[key]; ok {
		return enum
	}
	return RequestUnknown
}

// Values 获取所有枚举
func Values() []ResponseCode {
	values := make([]ResponseCode, 0, len(initiate))
	for k := range initiate {
		values = append(values, k)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}
