package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
)

var Config T

type ServerConfig struct {
	BackupFolder   string `yaml:"backup_folder"`
	ServerName     string `yaml:"server_name"`
	ProjectName    string `yaml:"project_name"`
	BackupInterval int    `yaml:"backup_interval"`
	MaxBackupData  int    `yaml:"max_backup_data"`
}

type WebDAVConfig struct {
	WebdavAddress string `yaml:"webdav_address"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DatabaseName string `yaml:"database_name"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
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
	printFields(Config)
}

func printFields(data interface{}) {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		if fieldVal.Kind() == reflect.Struct {
			printFields(fieldVal.Interface())
		} else {
			if reflect.ValueOf(fieldVal.Interface()).IsZero() {
				fmt.Printf("配置文件中字段 : %s 未填写或为空值", fieldTyp.Name)
			}
		}
	}
}
