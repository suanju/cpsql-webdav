package log

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	// 配置日志文件名及路径
	logFileName := fmt.Sprintf("logs/%s/log_%s.log", time.Now().Format("2006-01"), time.Now().Format("2006-01-02"))

	dir := filepath.Dir(logFileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
			panic(fmt.Sprintf("无法创建日志目录: %v", err))
		}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/1/2 15:04:05"))
	}
	// 配置日志文件分割器
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(getLogWriter(logFileName)),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	Logger = zap.New(core)
	zap.ReplaceGlobals(Logger)

}

// 获取日志写入器
func getLogWriter(logFileName string) zapcore.WriteSyncer {
	file, _ := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return zapcore.AddSync(file)
}
