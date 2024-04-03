package main

import (
	"cpsql-webdav/config"
	_ "cpsql-webdav/config"
	"cpsql-webdav/db"
	_ "cpsql-webdav/log"
	_ "cpsql-webdav/webdav"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	if _, err := db.FindMySQLDumpPath(); err != nil {
		log.Fatalf("执行失败 %s", err)
	}
	ticker := time.NewTicker(time.Duration(config.Config.Server.BackupInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			_, err := db.CreateInstance()
			if err != nil {
				zap.S().Errorf("备份数据库失败 创建数据库连接失败 %s", err)
				break
			}
			err = db.BackupDatabase()
			if err != nil {
				log.Printf("备份数据库失败: %v", err)
			}
		}
	}
}
