package log

import (
	"fmt"
	"go.uber.org/zap"
)

// 声明日志类，获取日志器

var Logger *zap.Logger

//	func init() {
//		var err error
//		Logger, err = NewLogger()
//		if err != nil {
//			panic(err)
//		}
//		Logger.Info("Logger初始化成功...")
//	}
func NewLogger() (*zap.Logger, error) {
	//	获取生产时日志配置器
	pro := zap.NewProductionConfig()
	//	自定义拼接日志输出地址
	pro.OutputPaths = append(pro.OutputPaths, "./productHandler.log", "stdout")
	//	构造日志器
	z, err := pro.Build()
	if err != nil {
		fmt.Println("zap初始化出现错误")
		panic(err)
	}
	return z, nil
}
