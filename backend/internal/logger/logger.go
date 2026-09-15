package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() {

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(
		"logs/app.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	fileSyncer := zapcore.AddSync(file)

	consoleSyncer := zapcore.AddSync(os.Stdout)

	writeSyncer := zapcore.NewMultiWriteSyncer(
		fileSyncer,
		consoleSyncer,
	)

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zap.InfoLevel,
	)

	Log = zap.New(core)

}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
