package network

import (
	"crypto/tls"
	"fmt"
	"local_dns_proxy/pkg/logger/log"
	"net"
	"net/http"
	"regexp"
	"time"
)

// GetLanIP 获取LAN IP地址
//
// 返回值：
//   - string LAN IP地址
func GetLanIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error().Err(err).Msg("获取本机IP失败")
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				// 跳过 169.254.x.x
				if !ip4.IsLinkLocalUnicast() {
					return ip4.String()
				}
			}
		}
	}
	log.Error().Msg("没有找到有效的IP地址")
	return "127.0.0.1"
}

// GetLocalIP 获取本地IP地址
//
// 返回值：
//   - string 本地IP地址
func GetLocalIP(ip string) string {
	ipv4Regex := regexp.MustCompile(`^((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)$`)
	if !ipv4Regex.MatchString(ip) && ip != "localhost" {
		ip = "127.0.0.1"
	}
	return ip
}

// IsMainWebStart 用 HTTP 请求检测主程序是否启动
// 参数：
//   - requestAddr 请求地址
//
// 返回值：
//   - bool 是否启动, true: 启动, false: 未启动
func IsMainWebStart(requestAddr string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
		},
	}

	resp, err := client.Get(requestAddr)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

// IsPortAvailable 检查端口是否可用
// 参数：
//   - port 端口号
//
// 返回值：
//   - bool 是否可用, true: 可用, false: 不可用
func IsPortAvailable(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), time.Millisecond*500)
	if err != nil {
		// 无法连接，说明没人监听，端口可用
		return true
	}
	conn.Close()
	return false
}

// GetPort 获取一个可用端口
// 参数：
//   - port 起始端口号
//
// 返回值：
//   - int 可用端口号
func GetPort(port int) int {
	for {
		if !IsPortAvailable(port) {
			port++
			continue
		}
		return port
	}
}
