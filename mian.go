package main

import (
	"cpsql-webdav/config"
	_ "cpsql-webdav/config"
	_ "cpsql-webdav/log"
	"cpsql-webdav/save"
	"cpsql-webdav/webdav"
	"flag"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

const stringEmpty = ""

func main() {
	//解析命令行参数
	serverName := flag.String("server-name", stringEmpty, "server name")
	projectName := flag.String("project-name", stringEmpty, "project name")

	webdevAddress := flag.String("webdav-address", stringEmpty, "webdav address")
	webdavUsername := flag.String("webdav-username", stringEmpty, "webdav username")
	webdavPassword := flag.String("webdav-password", stringEmpty, "webdav password")

	sqlDatabase := flag.String("sql-database", stringEmpty, "SQL database name")
	sqlUser := flag.String("sql-user", stringEmpty, "SQL user name")
	sqlPassword := flag.String("sql-password", stringEmpty, "SQL password")

	flag.Parse()

	if *serverName != stringEmpty {
		config.Config.Server.ServerName = *serverName
	}
	if *projectName != stringEmpty {
		config.Config.Server.ProjectName = *projectName
	}
	if *webdevAddress != stringEmpty {
		config.Config.WebDAV.WebdavAddress = *webdevAddress
	}
	if *webdavUsername != stringEmpty {
		config.Config.WebDAV.Username = *webdavUsername
	}
	if *webdavPassword != stringEmpty {
		config.Config.WebDAV.Password = *webdavPassword
	}
	if *sqlDatabase != stringEmpty {
		config.Config.Database.DatabaseName = *sqlDatabase
	}
	if *sqlUser != stringEmpty {
		config.Config.Database.Username = *sqlUser
	}
	if *sqlPassword != stringEmpty {
		config.Config.Database.Password = *sqlPassword
	}
	marshal, _ := yaml.Marshal(config.Config)
	_ = os.WriteFile("config.yaml", marshal, os.ModePerm)

	config.PrintFields(config.Config)
	webdav.CreateClient()

	if _, err := save.FindMySQLDumpPath(); err != nil {
		log.Fatalf("执行失败 %s", err)
	}
	ticker := time.NewTicker(time.Duration(config.Config.Server.BackupInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			_, err := save.CreateInstance()
			if err != nil {
				zap.S().Errorf("备份数据库失败 创建数据库连接失败 %s", err)
				break
			}
			err = save.BackupDatabase()
			if err != nil {
				log.Printf("备份数据库失败: %v", err)
			}
		}
	}
}
