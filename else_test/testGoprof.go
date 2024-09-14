package main

import (
	"log"
	"net/http"
	//	导入pprof包以启用性能分析功能（_表示只导入不使用）
	_ "net/http/pprof"
)

//	func main() {
//		example()
//	}
func example() {
	//	开启一个HTTP服务器
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//	模拟业务逻辑
	for {

	}
}
