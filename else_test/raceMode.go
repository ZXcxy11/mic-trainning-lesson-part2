package main

import (
	"fmt"
	"sync"
)

// 全局变量
var a int32

// 变量更改行为
func change(wg *sync.WaitGroup) {
	defer wg.Done() //	延迟完成任务
	for i := 0; i <= 1000; i++ {
		a++
	}
}

// 创建协程
func raceMode() {
	var wg sync.WaitGroup
	//	设定计数器
	wg.Add(4)
	go change(&wg)
	go change(&wg)
	go change(&wg)
	go change(&wg)
	//	等待全部协程完成
	wg.Wait()
	fmt.Println("最终值为：", a)
}

//func main() {
//	raceMode()
//}
