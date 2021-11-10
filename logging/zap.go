package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ztime "gogs.buffalo-robot.com/zouhy/micro/time"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var date string
var logger *zap.Logger

func NewSimpleLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	atom := zap.NewAtomicLevelAt(zap.WarnLevel)

	config := zap.Config{
		Level:            atom,                                                // 日志级别
		Development:      true,                                                // 开发模式，堆栈跟踪
		Encoding:         "json",                                              // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                                       // 编码器配置
		InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"},                                  // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
	return logger
}

func GetLogger(fileName, level string) *zap.Logger {
	once.Do(func() {
		logger = NewProdLoggger(fileName, level)
	})
	return logger
}

func NewProdLoggger(fileName, level string, ioWriters ...io.Writer) *zap.Logger {
	date = ztime.Now().Date()
	hook := lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     90,       // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	if level == "info" {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else if level == "debug" {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.DebugLevel)
	}
	var writers zapcore.WriteSyncer

	syncWriters := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)}
	for _, writer := range ioWriters {
		syncWriters = append(syncWriters, zapcore.AddSync(writer))
	}
	writers = zapcore.NewMultiWriteSyncer(syncWriters...)

	// }
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		writers,                               // 打印到控制台和文件
		atomicLevel,                           // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	logger := zap.New(core, caller, development, zap.AddStacktrace(zapcore.WarnLevel))

	timer := time.NewTicker(time.Minute)

	go func() {
		for {
			select {
			case <-timer.C:
				if date != ztime.Now().Date() {
					hook.Rotate()
					date = ztime.Now().Date()
				}
			}
		}
	}()

	return logger
}

type Emiter interface {
	Emit(event string, payload string)
}

type SocketLogger struct {
	emiter Emiter
	event  string
}

func NewSocketLogger(event string, emiter Emiter) *SocketLogger {
	return &SocketLogger{
		event:  event,
		emiter: emiter,
	}
}

func (s *SocketLogger) Write(p []byte) (n int, err error) {
	s.emiter.Emit(s.event, string(p))
	n = len(p)
	return
}
