package dns

import (
	"local_dns_proxy/internal/models"
)

type WebDeleteRow struct {
	ID int64 `json:"id"`
}

type WebSaveData struct {
	InsertRecords  []*models.DNSTable `json:"insertRecords"`
	RemoveRecords  []*models.DNSTable `json:"removeRecords"`
	UpdateRecords  []*models.DNSTable `json:"updateRecords"`
	PendingRecords []*models.DNSTable `json:"pendingRecords"`
}
