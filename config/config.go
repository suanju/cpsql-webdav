package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
)

var Config T

type ServerConfig struct {
	BackupFolder   string `yaml:"backup_folder" name:"webdav 备份文件夹名"`
	ServerName     string `yaml:"server_name" name:"服务器名称"`
	ProjectName    string `yaml:"project_name" name:"项目名称"`
	BackupInterval int    `yaml:"backup_interval" name:"备份间隔(秒)"`
	MaxBackupData  int    `yaml:"max_backup_data" name:"最大保存备份数"`
	SaveLocal      bool   `yaml:"save_local" name:"是否保留本地备份(ture|false)"`
}

type WebDAVConfig struct {
	WebdavAddress string `yaml:"webdav_address" name:"webdav 服务地址"`
	Username      string `yaml:"username" name:"webdav 服务账户"`
	Password      string `yaml:"password" name:"webdav 服务密码"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host" name:"数据库地址"`
	Port         int    `yaml:"port" name:"数据库端口"`
	DatabaseName string `yaml:"database_name" name:"数据库名称"`
	Username     string `yaml:"username" name:"数据库用户名"`
	Password     string `yaml:"password" name:"数据库密码"`
}

type T struct {
	Server   ServerConfig   `yaml:"server"`
	WebDAV   WebDAVConfig   `yaml:"webdav"`
	Database DatabaseConfig `yaml:"database"`
}

func init() {
	// 读取配置文件
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("无法读取配置文件 或配置文件不存在: %v", err)
	}

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("无法解析配置文件 或者格式错误: %v", err)
	}

}

func PrintFields(data interface{}) {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		if fieldVal.Kind() == reflect.Struct {
			PrintFields(fieldVal.Interface())
		} else {
			if reflect.ValueOf(fieldVal.Interface()).Kind() == reflect.Bool {
				continue
			}
			if reflect.ValueOf(fieldVal.Interface()).IsZero() {
				log.Fatalf("配置文件中字段 : %s 未填写或为空值 \n", fieldTyp.Name)
			}
		}
	}
}
