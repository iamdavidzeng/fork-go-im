package im

import (
	"fork_go_im/pkg/config"
	"fork_go_im/pkg/model"
	"fork_go_im/pkg/pool"
	"fork_go_im/pkg/redis"
	"time"
)

func SetupPool() {
	db := model.ConnectDB()
	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_file_seconds")) * time.Second)

	redis.InitClient()

	pool.ConnectPool()
}
