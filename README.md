# local_dns_proxy

> 本地 DNS 代理管理服务

```sh
Local DNS Proxy - 本地 DNS 代理服务
--------------------------------
用法:
    local_dns_proxy.exe [命令]

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
    local_dns_proxy.exe install
    local_dns_proxy.exe start
    local_dns_proxy.exe stop
    local_dns_proxy.exe uninstall
    local_dns_proxy.exe server
```

![](https://private-user-images.githubusercontent.com/158544427/492040465-13bdb832-4e62-4a79-affb-2d4de656eb94.png?jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NTg0NzI1NTksIm5iZiI6MTc1ODQ3MjI1OSwicGF0aCI6Ii8xNTg1NDQ0MjcvNDkyMDQwNDY1LTEzYmRiODMyLTRlNjItNGE3OS1hZmZiLTJkNGRlNjU2ZWI5NC5wbmc_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1BS0lBVkNPRFlMU0E1M1BRSzRaQSUyRjIwMjUwOTIxJTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI1MDkyMVQxNjMwNTlaJlgtQW16LUV4cGlyZXM9MzAwJlgtQW16LVNpZ25hdHVyZT03M2Y0MzZmOGI0NDQxMWU5OWZmMzU4NGQ0ZmRiYmRhODBkODgzYTViOWFiZGJjOGE5YjZmYzc3YmNkNjE3MGJkJlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCJ9.gTIE0sv_D64IEGdgjulFvLoyXCKlEvS5nR_3jb-Q6JI)