package initdata

import (
	"local_dns_proxy/internal/models"
	"local_dns_proxy/pkg/cfg"
	"local_dns_proxy/pkg/enums/status"
	"local_dns_proxy/pkg/logger/log"
	"strconv"
)

// Config 初始化配置
func (r *InitData) Config() {
	modelList := []models.DNSTable{
		{Protocol: "http", Domain: "localhost", IP: "127.0.0.1", Port: strconv.Itoa(cfg.C.Server.Port), Status: status.Enable},
		{Protocol: "http", Domain: "www.localhost.com", IP: "127.0.0.1", Port: strconv.Itoa(cfg.C.Server.Port), Status: status.Enable},
	}

	newModelList := insertIfNotExist[models.DNSTable](modelList, func(model models.DNSTable) (*models.DNSTable, error) {
		return r.Q.DNSTable.Where(r.Q.DNSTable.Domain.Eq(model.Domain)).Take()
	})

	if len(newModelList) <= 0 {
		return
	}

	if err := r.Q.DNSTable.CreateInBatches(newModelList, 100); err != nil {
		log.Error().Err(err).Msg("[DNS表]初始化失败！")
	} else {
		log.Info().Msgf("[DNS表]初始化成功，共%d条数据！", len(newModelList))
	}
}
