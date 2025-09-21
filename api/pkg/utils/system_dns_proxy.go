package utils

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"local_dns_proxy/pkg/logger/log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type DNSProxy struct {
	originalDNS []string
	hosts       map[string]string
	IsRunning   bool

	CtxStatus  context.Context
	CancelStop context.CancelFunc
}

// NetIface 网卡信息结构体
type NetIface struct {
	Name         string    `json:"name"`
	HardwareAddr string    `json:"hardware_addr"`
	Flags        net.Flags `json:"flags"`
	Addrs        []string  `json:"addrs"`
}

var (
	ctxStatus  context.Context
	cancelStop context.CancelFunc
	dnsProxy   *DNSProxy
)

// handleDNS 处理 DNS 请求
func (r *DNSProxy) handleDNS(w dns.ResponseWriter, msg *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(msg)
	upstream := "114.114.114.114:53"

	for _, q := range msg.Question {
		hostKey := strings.TrimSuffix(strings.ToLower(q.Name), ".")
		switch q.Qtype {
		case dns.TypeA:
			if ip, ok := r.hosts[hostKey]; ok {
				rr, _ := dns.NewRR(q.Name + " A " + ip)
				m.Answer = append(m.Answer, rr)
				log.Debug().Msgf("[命中本地DNS映射] %s -> %s", hostKey, ip)
			} else {
				// 转发到上游
				resp, err := dns.Exchange(msg, upstream)
				if err == nil {
					m.Answer = append(m.Answer, resp.Answer...)
				} else {
					m.Rcode = dns.RcodeNameError
					log.Debug().Err(err).Msgf("[上游DNS查询失败] %s", hostKey)
				}
			}
		default:
			// 其它类型请求直接转发
			resp, err := dns.Exchange(msg, upstream)
			if err == nil {
				m.Answer = append(m.Answer, resp.Answer...)
			} else {
				m.Rcode = dns.RcodeNameError
				log.Debug().Err(err).Msgf("[上游DNS查询失败-非A记录] %s", hostKey)
			}
		}
	}
	_ = w.WriteMsg(m)
}

// StartDnsServer 启动 DNS 服务器
func (r *DNSProxy) StartDnsServer() {
	// 启动 DNS 服务
	udpServer := &dns.Server{Addr: ":53", Net: "udp"}
	tcpServer := &dns.Server{Addr: ":53", Net: "tcp"}

	dns.HandleFunc(".", r.handleDNS)

	go func() {
		log.Info().Msg("[DNS代理] 启动本地 DNS 服务 (UDP) on *:53")
		if err := udpServer.ListenAndServe(); err != nil {
			r.IsRunning = false
			log.Error().Err(err).Msg("本地 UDP DNS 启动失败！")
		}
	}()
	go func() {
		log.Info().Msg("[DNS代理] 启动本地 DNS 服务 (TCP) on *:53")
		if err := tcpServer.ListenAndServe(); err != nil {
			r.IsRunning = false
			log.Error().Err(err).Msg("本地 TCP DNS 启动失败！")
		}
	}()
	time.Sleep(500 * time.Millisecond)
	for domain, ip := range r.hosts {
		log.Info().Msgf("[DNS代理] 已注册临时域名映射: %s -> %s", domain, ip)
	}

	r.CancelStop()
	<-ctxStatus.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	flag := true
	if err := udpServer.ShutdownContext(shutdownCtx); err != nil {
		flag = false
		log.Error().Err(err).Msg("本地 UDP DNS 关闭失败！")
	}
	if err := tcpServer.ShutdownContext(shutdownCtx); err != nil {
		flag = false
		log.Error().Err(err).Msg("本地 TCP DNS 关闭失败！")
	}
	if flag {
		dnsProxy = nil
		ctxStatus = nil
		cancelStop = nil

		log.Info().Msg("[DNS代理] 本地 DNS 服务已停止")
	}
}

// SetLocalDNS 设置系统 DNS 为 127.0.0.1
func (r *DNSProxy) SetLocalDNS(ifaceName string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("netsh", "interface", "ipv4", "set", "dns", `name="`+ifaceName+`"`, "static", "127.0.0.1")
		log.Debug().Str("Name", ifaceName).Str("OS", runtime.GOOS).Str("cmd", strings.Join(cmd.Args, " ")).Msg("系统代理DNS设置")
		if err := cmd.Run(); err != nil {
			return err
		}
	case "darwin":
		cmd := exec.Command("networksetup", "-setdnsservers", ifaceName, "127.0.0.1")
		log.Debug().Str("Name", ifaceName).Str("OS", runtime.GOOS).Str("cmd", strings.Join(cmd.Args, " ")).Msg("系统代理DNS设置")
		if err := cmd.Run(); err != nil {
			return err
		}
	case "linux":
		backup, _ := os.ReadFile("/etc/resolv.conf")
		r.originalDNS = strings.Split(string(backup), "\n")
		_ = os.WriteFile("/etc/resolv.conf", []byte("nameserver 127.0.0.1\n"), 0644)
	default:
		return fmt.Errorf("不支持的操作系统")
	}
	log.Info().Msg("[DNS代理] 系统 DNS 设置为 127.0.0.1")
	return nil
}

// RestoreDNS 恢复系统 DNS
//
// 说明:
//   - 恢复系统 DNS 为之前的配置。
//     c := make(chan os.Signal, 1)
//     signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//     <-c
//     RestoreDNS()
func (r *DNSProxy) RestoreDNS(ifaceName string) error {
	if dnsProxy == nil {
		return nil
	}
	cancelStop()
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("netsh", "interface", "ipv4", "set", "dns", ifaceName, "dhcp")
		log.Debug().Str("Name", ifaceName).Str("OS", runtime.GOOS).Str("cmd", strings.Join(cmd.Args, " ")).Msg("[DNS代理] 系统代理DNS恢复")
		if err := cmd.Run(); err != nil {
			log.Error().Err(err).Str("网卡名称", ifaceName).Msg("系统代理DNS设置失败")
			return err
		}
	case "darwin":
		cmd := exec.Command("networksetup", "-setdnsservers", ifaceName, "Empty")
		log.Debug().Str("网卡名称", ifaceName).Str("OS", runtime.GOOS).Str("cmd", strings.Join(cmd.Args, " ")).Msg("[DNS代理] 系统代理DNS恢复")
		if err := cmd.Run(); err != nil {
			log.Error().Err(err).Str("网卡名称", ifaceName).Msg("系统代理DNS设置失败")
			return err
		}
	case "linux":
		if len(r.originalDNS) > 0 {
			if err := os.WriteFile("/etc/resolv.conf", []byte(strings.Join(r.originalDNS, "\n")), 0644); err != nil {
				log.Error().Err(err).Msg("系统代理DNS恢复失败")
				return err
			}
		}
	default:
		log.Error().Msgf("不支持的操作系统：%s", runtime.GOOS)
		return fmt.Errorf("不支持的操作系统：%s", runtime.GOOS)
	}
	log.Info().Msg("[DNS代理] 恢复系统 DNS 代理")
	return nil
}

// GetNetworkInterfaces 获取网络接口列表
func GetNetworkInterfaces() ([]NetIface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var results []NetIface
	for _, iface := range ifaces {
		// 过滤掉未启用或回环网卡
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		var addrStrs []string
		for _, a := range addrs {
			addrStrs = append(addrStrs, a.String())
		}

		results = append(results, NetIface{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        iface.Flags,
			Addrs:        addrStrs,
		})
	}

	return results, nil
}

func NewDNSProxy(hosts map[string]string) *DNSProxy {
	if dnsProxy == nil {
		ctxStatus, cancelStop = context.WithCancel(context.Background())
		dnsProxy = &DNSProxy{IsRunning: true}
		dnsProxy.CtxStatus, dnsProxy.CancelStop = context.WithCancel(context.Background())
	}
	dnsProxy.hosts = hosts
	return dnsProxy
}

func GetDNSProxy() *DNSProxy {
	return dnsProxy
}
