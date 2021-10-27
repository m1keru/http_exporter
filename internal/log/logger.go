package log

import (
	"os"

	"github.com/m1keru/http_exporter/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//SharedLogger - Logger Instance
var SharedLogger *zap.SugaredLogger

//InitLogger - InitLogger
func InitLogger(logConfig config.Log) {
	writeSyncer := getLogWriter(logConfig.Path)
	encoder := getEncoder()
	var logLevel zapcore.Level
	logLevel.Set(logConfig.Level)
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	instance := zap.New(core)
	SharedLogger = instance.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(logPath string) zapcore.WriteSyncer {
	file, _ := os.Create(logPath)
	return zapcore.AddSync(file)
}
