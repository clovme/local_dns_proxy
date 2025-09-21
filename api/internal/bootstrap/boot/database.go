package boot

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"local_dns_proxy/internal/infrastructure/query"
	"local_dns_proxy/pkg/constants"
	"local_dns_proxy/pkg/logger"
	"local_dns_proxy/pkg/logger/log"
	"local_dns_proxy/pkg/utils/file"
	"os"
	"time"
)

// InitializationDB 打开数据库连接, 并设置连接池, 数据库链接统一入口
// 返回值：
//   - *gorm.DB 数据库连接对象
func InitializationDB() *gorm.DB {
	dbPath, err := file.GetFileAbsPath(fmt.Sprintf("%s_data.db", constants.ProjectName))
	if err != nil {
		return nil
	}
	var dsn = sqlite.Open(dbPath)

	db, err := gorm.Open(dsn, logger.GetGormLogger())
	if err != nil {
		log.Error().Err(err).Msg("打开SQLite失败")
		os.Exit(-1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("获取底层 sql.DB 失败")
		os.Exit(-1)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	query.SetDefault(db)
	return db
}
