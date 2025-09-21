package utils

import (
	"strings"
	"unicode"
)

// Capitalize 首字符大写
// 参数：
//   - s 字符串
//
// 返回值：
//   - string 首字符大写后的字符串
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	// 转成 rune 切片，防止中文/多字节字符乱码
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CamelToSnake 智能将 PascalCase → lower_camel_case，保留开头缩写正确形式
// 例如：UserID → user_id，HTTPResponse → http_response
// 参数：
//   - s 字符串
//
// 返回值：
//   - string 转换后的字符串
func CamelToSnake(s string) string {
	var result []rune
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		if i > 0 {
			// 当前是大写，前面是小写，或者当前是大写，前面是大写，后面是小写
			if unicode.IsUpper(runes[i]) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
				result = append(result, '_')
			}
		}
		result = append(result, unicode.ToLower(runes[i]))
	}
	temp := strings.ToLower(string(result))
	if strings.HasPrefix(temp, "_") {
		temp = temp[1:]
	}
	return temp
}
