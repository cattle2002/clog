package log

import (
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getZapCoreWithWriter(logDir string, logFile string, maxSize int, maxAge int, maxBackups int, compress bool, logLevel string) zapcore.Core {
	var level zapcore.LevelEnabler
	if logLevel == "debug" {
		level = zap.DebugLevel
	} else if logLevel == "info" {
		level = zap.InfoLevel
	} else if logLevel == "error" {
		level = zap.ErrorLevel
	} else if logLevel == "panic" {
		level = zap.PanicLevel
	} else if logLevel == "fatal" {
		level = zap.FatalLevel
	}
	pos := filepath.Join(logDir, logFile)
	writer := lumberjack.Logger{
		Filename:   pos,
		MaxSize:    maxSize, // 当日志文件大小超过此值时，将被分割，单位为MB，此处设置的是1MB。
		MaxAge:     maxAge,  // 历史日志的保留天数
		MaxBackups: maxBackups,
		LocalTime:  true,
		Compress:   compress, // 在实际生产环境中，往往需要压缩
	}

	cfg := zap.NewProductionEncoderConfig()

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewTee(zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(&writer), level))

	// return zapcore.NewTee(zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(&writer), zap.DebugLevel))

	// return zapcore.NewTee(zapcore.NewCore(encoder, zapcore.AddSync(&writer), zap.DebugLevel))

}
func NewLog(logDir string, logFile string, maxSize int, maxAge int, maxBackups int, compress bool, logLevel string) (*zap.Logger, error) {
	return zap.NewProduction(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return getZapCoreWithWriter(logDir, logFile, maxSize, maxAge, maxBackups, compress, logLevel)
	}), zap.AddStacktrace(zap.ErrorLevel))
}
