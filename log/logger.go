package log

import "go.uber.org/zap"

// 声明日志类，获取日志器

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = NewLogger()
	if err != nil {
		panic(err)
	}
}
func NewLogger() (*zap.Logger, error) {
	//	获取生产时日志配置器
	pro := zap.NewProductionConfig()
	//	自定义拼接日志输出地址
	pro.OutputPaths = append(pro.OutputPaths, "./accountHandler.log")
	//	构造日志器
	return pro.Build()
}
