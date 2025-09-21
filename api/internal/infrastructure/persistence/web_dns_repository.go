package persistence

import (
	"context"
	"fmt"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"local_dns_proxy/internal/infrastructure/query"
	dnsDTO "local_dns_proxy/internal/schema/dto/dns"
	dnsVO "local_dns_proxy/internal/schema/vo/dns"
	"local_dns_proxy/pkg/enums/status"
	"local_dns_proxy/pkg/utils"
	"strings"
)

type WebDnsRepository struct {
	DB *gorm.DB
	Q  *query.Query
}

func (r *WebDnsRepository) GetDnsList(orderBy string) ([]*dnsVO.WebDNSTableVO, error) {
	var results []*dnsVO.WebDNSTableVO

	orders := make([]field.Expr, 0)

	for _, p := range strings.Split(orderBy, ",") { // e.g. "domain|asc,ip|asc,status|desc"
		kv := strings.Split(p, "|")
		if len(kv) != 2 {
			continue
		}
		// 从生成的 query 里找字段
		col, ok := r.Q.DNSTable.GetFieldByName(kv[0])
		if !ok {
			continue // 前端传了不存在的字段，跳过
		}

		// 拼接排序表达式
		if strings.EqualFold(kv[1], "desc") {
			orders = append(orders, col.Desc())
		} else {
			orders = append(orders, col.Asc()) // 默认 ASC
		}
	}

	queryDNSTable := r.Q.DNSTable.WithContext(context.Background())
	if len(orders) > 0 {
		queryDNSTable = queryDNSTable.Order(orders...)
	}

	err := queryDNSTable.Scan(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *WebDnsRepository) DeleteDnsRow(id int64) error {
	info, err := r.Q.DNSTable.Unscoped().Where(r.Q.DNSTable.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	if info.Error != nil {
		return info.Error
	}
	if info.RowsAffected <= 0 {
		return fmt.Errorf("删除失败，未影响到任何数据")
	}
	return nil
}

func (r *WebDnsRepository) SaveDnsData(data dnsDTO.WebSaveData) (vo dnsVO.WebSaveDataVO) {
	t := r.Q.DNSTable
	vo.InsertCount = func() int {
		if err := t.CreateInBatches(data.InsertRecords, 100); err != nil {
			return 0
		}
		return len(data.InsertRecords)
	}()

	vo.UpdateCount = func() int {
		recordCount := 0
		for _, item := range data.UpdateRecords {
			if _, err := t.Where(t.ID.Eq(item.ID)).Select(t.Protocol, t.Domain, t.IP, t.Port, t.Status).Updates(item); err != nil {
				continue
			}
			recordCount++
		}
		return recordCount
	}()

	vo.DeleteCount = func() int64 {
		ids := make([]int64, 0)
		for _, item := range data.PendingRecords {
			ids = append(ids, item.ID)
		}
		for _, item := range data.RemoveRecords {
			ids = append(ids, item.ID)
		}
		if len(ids) <= 0 {
			return 0
		}

		info, err := t.Unscoped().Where(t.ID.In(ids...)).Delete()
		if err != nil {
			return 0
		}

		return info.RowsAffected
	}()
	return vo
}

func (r *WebDnsRepository) NewDNSProxy() (*utils.DNSProxy, error) {
	t := r.Q.DNSTable
	var dnsList []*dnsVO.WebDNSTableVO

	if err := t.Where(t.Status.Eq(status.Enable.Int())).Scan(&dnsList); err != nil {
		return nil, fmt.Errorf("未找到可映射的 DNS 记录，%+v", err)
	}
	if len(dnsList) <= 0 {
		return nil, fmt.Errorf("未找到可映射的 DNS 记录，请先添加/开启 DNS 记录")
	}

	hosts := make(map[string]string)
	for _, dns := range dnsList {
		hosts[strings.ToLower(dns.Domain)] = dns.IP
	}

	return utils.NewDNSProxy(hosts), nil
}
