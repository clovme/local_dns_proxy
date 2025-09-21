package dns

import (
	dnsDTO "local_dns_proxy/internal/schema/dto/dns"
	dnsVO "local_dns_proxy/internal/schema/vo/dns"
	"local_dns_proxy/pkg/utils"
)

type Repository interface {
	GetDnsList(orderBy string) ([]*dnsVO.WebDNSTableVO, error)
	DeleteDnsRow(id int64) error
	SaveDnsData(data dnsDTO.WebSaveData) dnsVO.WebSaveDataVO
	NewDNSProxy() (*utils.DNSProxy, error)
}
