package webdav

import (
	config2 "cpsql-webdav/config"
	"fmt"
	"github.com/studio-b12/gowebdav"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var Client *gowebdav.Client

func init() {
	config := config2.Config.WebDAV
	// 创建WebDAV客户端
	Client = gowebdav.NewClient(config.WebdavAddress, config.Username, config.Password)
	// 尝试获取远程目录的信息
	_, err := Client.Stat("/")
	if err != nil {
		log.Fatalf("无法连接到 WebDAV 服务器: %v", err)
	}
	zap.S().Info("接到 WebDAV 服务器成功")
}

func UploadFileToWebDAV(client *gowebdav.Client, localFilePath string, remoteFilePath string) error {
	remoteDir := filepath.ToSlash(filepath.Dir(remoteFilePath))
	if _, err := client.Stat(remoteDir); err != nil {
		zap.S().Errorf("webdav 上传目录不存在 %s : %v", remoteDir, err)
		err := client.MkdirAll(remoteDir, os.ModePerm)
		if err != nil {
			zap.S().Errorf("webdav 创建上传目录失败 %s  %v", remoteDir, err)
		}
	}
	//超过最大值先删除
	if err := CheckUploadMaxBackup(client, remoteDir, config2.Config.Server.MaxBackupData); err != nil {
		return err
	}

	// 读取本地文件内容
	data, err := os.ReadFile(localFilePath)
	if err != nil {
		zap.S().Errorf("无法读取本地文件: %v", err)
		return err
	}
	// 将本地文件内容写入到 WebDAV 服务器上的目标文件中
	err = client.Write(remoteFilePath, data, os.ModePerm)
	if err != nil {
		zap.S().Errorf("无法上传文件到 WebDAV 服务器: %v", err)
		return err
	}
	zap.S().Info("文件 %s 已同步到webdav 位置位于 %s", localFilePath, remoteFilePath)
	return nil
}

func CheckUploadMaxBackup(client *gowebdav.Client, path string, max int) error {
	fmt.Println(max)
	files, err := client.ReadDir(path)
	if err != nil {
		zap.S().Errorf("无法读取目录 %v", err)
		return err
	}
	fmt.Println(len(files))
	if len(files) >= max {
		//排序
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})
		filePath := filepath.ToSlash(filepath.Join(path, files[0].Name()))
		fmt.Println(filePath)
		err := client.Remove(filePath)
		if err != nil {
			zap.S().Errorf("无法删除文件 %v", err)
			return err
		}
		zap.S().Errorf("保留最大存储份数 : %d ; 删除文件成功 %s", max, filePath)
	}
	return nil
}
