package db

import (
	"archive/zip"
	"cpsql-webdav/webdav"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"cpsql-webdav/config"
	_ "github.com/go-sql-driver/mysql"
)

const backup = "backup"

func CreateInstance() (db *sql.DB, err error) {
	mysqlConfig := config.Config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DatabaseName)
	db, _ = sql.Open("mysql", dsn)
	if err := db.Ping(); err != nil {
		zap.S().Errorf("无法连接到数据库: %v", err)
		return nil, err
	}
	return db, nil
}

func BackupDatabase() error {
	workDir, _ := os.Getwd()
	dir := filepath.Join(workDir, backup)
	if _, err := os.Stat(dir); err != nil {
		err := os.MkdirAll(dir, fs.ModePerm)
		if err != nil {
			zap.S().Errorf("无法创建目录： %v", err)
		}
	}
	mysqlConfig := config.Config.Database
	fileName := "backup-" + time.Now().Format("2006-01-02_15-04-05")
	backupFileName := filepath.Join(dir, fileName+".sql")
	backupZipName := filepath.Join(dir, fileName+".zip")
	cmd := exec.Command("mysqldump", "-u", mysqlConfig.Username, fmt.Sprintf("-p%s", mysqlConfig.Password), mysqlConfig.DatabaseName, "--result-file="+backupFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	if err = CreatZip(backupFileName, backupZipName); err != nil {
		return err
	}
	//删除sql保留zip
	_ = os.Remove(backupFileName)
	serverConfig := config.Config.Server
	remoteFileDirectory := filepath.ToSlash(filepath.Join(serverConfig.BackupFolder, serverConfig.ServerName, serverConfig.ProjectName))
	remoteFilePath := filepath.ToSlash(filepath.Join(remoteFileDirectory, fileName+".zip"))
	//将zip同步到webdav
	return webdav.UploadFileToWebDAV(webdav.Client, backupZipName, remoteFilePath)
}

func FindMySQLDumpPath() (path string, err error) {
	fmt.Printf("当前操作系统为 %s 寻找 mysqldump \n", runtime.GOOS)
	pathEnv := os.Getenv("PATH")
	mysqldump := "mysqldump"
	if runtime.GOOS == "windows" {
		mysqldump = "mysqldump.exe"
	}
	paths := filepath.SplitList(pathEnv)
	for _, p := range paths {
		mysqldumpPath := filepath.Join(p, mysqldump)
		if _, err := os.Stat(mysqldumpPath); err == nil {
			fmt.Printf("mysqldump 位置位于 %s \n", mysqldumpPath)
			return mysqldumpPath, nil
		}
	}
	return "", fmt.Errorf("未找到 mysqldump 工具，不能完成备份操作")
}

func CreatZip(path string, toPath string) error {
	// 打开要压缩的文件
	file, err := os.Open(path)
	if err != nil {
		zap.S().Errorf("CreatZip 打开文件时出错： %v", err)
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 创建 ZIP 文件
	zipFile, err := os.Create(toPath)
	if err != nil {
		zap.S().Errorf("CreatZip 创建 ZIP 文件时出错： %v", err)
		return err
	}
	defer func(zipFile *os.File) {
		_ = zipFile.Close()
	}(zipFile)

	// 创建 ZIP Writer
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	// 将文件添加到 ZIP 中
	fileInfo, _ := file.Stat()
	header := &zip.FileHeader{
		Name:   fileInfo.Name(),
		Method: zip.Deflate,
	}
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		zap.S().Errorf("CreatZip 创建 ZIP 头时出错： %v", err)
		return err
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		zap.S().Errorf("写入 ZIP 文件时出错： %v", err)
		return err
	}
	zap.S().Info("文件压缩成功", err)
	return nil
}
