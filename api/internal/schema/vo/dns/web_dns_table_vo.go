package dns

import (
	"local_dns_proxy/pkg/enums/status"
	"time"
)

type WebDNSTableVO struct {
	ID        int64         `json:"id"`
	Protocol  string        `json:"protocol"`
	Domain    string        `json:"domain"`
	IP        string        `json:"ip"`
	Status    status.Status `json:"status"`
	Port      string        `json:"port"`
	CreatedAt *time.Time    `json:"createdAt"`
	UpdatedAt *time.Time    `json:"updatedAt"`
}

type WebSaveDataVO struct {
	InsertCount int   `json:"insertCount"`
	UpdateCount int   `json:"updateCount"`
	DeleteCount int64 `json:"deleteCount"`
}
