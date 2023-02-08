package app

import (
	"io"
	"os"
	"path/filepath"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Dir     string    // 目录
	LogFile io.Writer // 日志文件
)

// 初始化
func init() {
	// 配置
	var err error
	Dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Fatal(err)
	}
	confPath := filepath.Join(Dir, "conf/app.*")
	confFiles, err := filepath.Glob(confPath)
	if err != nil {
		logrus.Fatal(err)
	}
	for _, f := range confFiles {
		v := viper.New()
		v.SetConfigFile(f)
		err = v.ReadInConfig()
		if err != nil {
			logrus.Fatal(err)
		}
		viper.MergeConfigMap(v.AllSettings())
	}
	// 日志
	logPath := filepath.Join(Dir, "logs/log.%Y%m%d.log")
	logCount := viper.GetUint("log.count")
	if logCount == 0 {
		logCount = 7
	}
	LogFile, err = rotatelogs.New(
		logPath,
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(logCount),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, LogFile))
	logLevel, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err == nil {
		logrus.SetLevel(logLevel)
	}
	logrus.SetReportCaller(true)
}
