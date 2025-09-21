package models

import (
	"local_dns_proxy/pkg/enums/status"
	"time"
)

// DNSTable DNS代理配置项表
type DNSTable struct {
	ID        int64         `gorm:"primaryKey;autoIncrement;comment:ID，主键"`
	Protocol  string        `gorm:"type:varchar(10);not null;comment:协议"`
	Domain    string        `gorm:"type:varchar(50);not null;unique;comment:域名"`
	IP        string        `gorm:"type:varchar(20);not null;comment:IP地址"`
	Port      string        `gorm:"type:varchar(6);default:53;comment:端口"`
	Status    status.Status `gorm:"default:1;comment:状态：Enable启用，Disable禁用"`
	CreatedAt *time.Time    `gorm:"column:created_at;autoCreateTime:nano;comment:创建时间"`
	UpdatedAt *time.Time    `gorm:"column:updated_at;autoUpdateTime:nano;comment:更新时间"`
}

func (r *DNSTable) TableComment() string {
	return "DNS代理配置表"
}
