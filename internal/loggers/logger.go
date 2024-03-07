package loggers

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

// RunLoggerConfig configures logrus by setting the logging level,
// the format of logged information, and configuring log writing to a file.
func RunLoggerConfig(EnvLogs string) {

	logLevel, err := logrus.ParseLevel(EnvLogs)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	//Настраиваем формат логируемой информации
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			_, filename := path.Split(f.File)
			filename = fmt.Sprintf("%s.%d.%s", filename, f.Line, f.Function)
			return "", filename
		},
	})
	// Настраиваем запись логов в файл
	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "firstService.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
	})
	logrus.SetOutput(mw)
}
