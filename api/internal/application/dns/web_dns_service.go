package dns

import (
	dnsDTO "local_dns_proxy/internal/schema/dto/dns"
	dnsVO "local_dns_proxy/internal/schema/vo/dns"
	"local_dns_proxy/pkg/utils"
)

type WebDnsService struct {
	Repo Repository
}

func (r *WebDnsService) FindDnsList(orderBy string) ([]*dnsVO.WebDNSTableVO, error) {
	return r.Repo.GetDnsList(orderBy)
}

func (r *WebDnsService) DeleteDnsData(id int64) error {
	return r.Repo.DeleteDnsRow(id)
}

func (r *WebDnsService) SaveDnsData(data dnsDTO.WebSaveData) dnsVO.WebSaveDataVO {
	return r.Repo.SaveDnsData(data)
}

func (r *WebDnsService) NewDNSProxy() (*utils.DNSProxy, error) {
	return r.Repo.NewDNSProxy()
}
