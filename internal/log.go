package internal

import (
	"fmt"
	"go.uber.org/zap"
)

func init() {
	//	//zap.NewDevelopment()创建一个开发模式下的日志器，默认将其输出到控制台`stdout`，
	//	此处通过 NewDevelopmentConfig() 定制输出路径，除了控制台还有外部文件
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout", "./productHandler.log"}
	development, err := config.Build()
	if err != nil {
		panic(err)
	}
	//	ReplaceGlobals将全局Logger替换
	zap.ReplaceGlobals(development)
	fmt.Println("zap初始化完毕..")
}
