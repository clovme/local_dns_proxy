package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"io"
	"local_dns_proxy/internal/bootstrap/boot"
	"local_dns_proxy/internal/bootstrap/database"
	"local_dns_proxy/internal/bootstrap/routers"
	"local_dns_proxy/internal/infrastructure/query"
	"local_dns_proxy/pkg/cfg"
	"local_dns_proxy/pkg/logger/log"
	"local_dns_proxy/pkg/utils"
	"local_dns_proxy/pkg/utils/file"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Start the work in a goroutine.
	p.exit = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	// 初始化表单验证器
	boot.InitializationFormValidate()

	// 连接数据库
	db := boot.InitializationDB()

	if !file.IsFileExist(cfg.Path) {
		// 初始化数据库迁移
		if err := database.AutoMigrate(db, query.Q); err != nil {
			log.Fatal().Err(err).Msg("数据库迁移失败")
		}
		cfg.SaveToIni()
		log.Info().Msg("数据初始化完成")
	}

	// 初始化路由
	engine := routers.Initialization(db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go engine.Run("localhost", cfg.C.Server.Port)
	<-c
	if err := utils.GetDNSProxy().RestoreDNS(cfg.C.Server.Iface); err != nil {
		log.Error().Err(err).Msg("系统代理DNS设置失败")
		return
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop 应尽快返回，通知 goroutine 退出
	close(p.exit)
	return nil
}

func init() {
	time.Local = time.UTC
	// 禁用 Gin 框架的日志输出
	gin.DefaultWriter = io.Discard
}

func printHelp() {
	str := `Local DNS Proxy - 本地 DNS 代理服务
--------------------------------
用法:
    LocalDNSProxy [命令]

命令:
    install      安装服务 (注册为 Windows 服务)
    uninstall    卸载服务
    start        启动已安装的服务
    stop         停止已安装的服务
    server       以前台控制台模式运行 (调试用)

说明:
    - 不带任何命令运行时, 默认以服务方式启动 (由系统服务管理器调用)
    - 安装/卸载/启动/停止 服务需要管理员权限
    - server 模式会在控制台运行, 输出日志, 适合开发调试

示例:
    LocalDNSProxy install
    LocalDNSProxy start
    LocalDNSProxy stop
    LocalDNSProxy uninstall
    LocalDNSProxy server
`
	fmt.Println(strings.ReplaceAll(str, "LocalDNSProxy", filepath.Base(os.Args[0])))
}

func main() {
	// 初始化系统日志
	boot.InitializationLogger(cfg.C.Logger)

	svcConfig := &service.Config{
		Name:        "local_dns_proxy",
		DisplayName: "本地DNS代理",
		Description: "本地DNS代理服务，用于代理本地DNS请求",
		Arguments:   []string{"server"},
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal().Err(err).Msgf("service.New error %s", svcConfig.Name)
		os.Exit(1)
	}

	if len(os.Args) <= 1 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "install":
		if err := s.Install(); err != nil {
			log.Info().Err(err).Msgf("service install error %s", svcConfig.Name)
		} else {
			log.Info().Msgf("service installed %s", svcConfig.Name)
		}
		return
	case "uninstall":
		if err := s.Uninstall(); err != nil {
			log.Info().Err(err).Msgf("service uninstall error %s", svcConfig.Name)
		} else {
			log.Info().Msgf("service uninstalled %s", svcConfig.Name)
		}
		return
	case "start":
		if err := s.Start(); err != nil {
			log.Info().Err(err).Msgf("service start error %s", svcConfig.Name)
		} else {
			log.Info().Msgf("service started %s", svcConfig.Name)
			log.Info().Msgf("访问地址：http://localhost:%d", cfg.C.Server.Port)
		}
		return
	case "stop":
		if err := s.Stop(); err != nil {
			log.Info().Err(err).Msgf("service stop error %s", svcConfig.Name)
		} else {
			log.Info().Msgf("service stopped %s", svcConfig.Name)
		}
		return
	case "server":
		// 以控制台模式运行（方便调试）
		if err := s.Run(); err != nil {
			log.Info().Err(err).Msgf("service run error %s", svcConfig.Name)
		}
		return
	}
}
