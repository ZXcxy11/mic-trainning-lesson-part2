package main

import (
	"fmt"
	"runtime"
	"time"
)

// 创建垃圾数据
func createGarbage() {
	for i := 0; i < 100000; i++ {
		_ = make([]byte, 1024*10)
	}
}
func stressTest() {
	for j := 0; j < 100; j++ {
		go createGarbage()
		time.Sleep(100 * time.Millisecond)
	}
}

// 检测各阶段内存的占用
func checkMem() {
	var memStats runtime.MemStats

	//	检测占用前的内存使用
	runtime.ReadMemStats(&memStats)
	//Alloc时分配的堆的字节数
	fmt.Println("初始内存使用：", memStats.Alloc/1024, "KB")
	fmt.Println("系统分配的内存(HeapSys)：", memStats.HeapSys/1024, "KB")
	fmt.Println("未使用的内存(HeapIdle)：", memStats.HeapIdle/1024, "KB")
	fmt.Println("GC的次数：", memStats.NumGC)

	//	生成垃圾
	stressTest()

	//	生成后内存占用
	runtime.ReadMemStats(&memStats)
	fmt.Println("生成垃圾后的内存占用：", memStats.Alloc/1024, "KB")
	fmt.Println("系统分配的内存(HeapSys)：", memStats.HeapSys/1024, "KB")
	fmt.Println("未使用的内存(HeapIdle)：", memStats.HeapIdle/1024, "KB")
	fmt.Println("GC的次数：", memStats.NumGC)

	//	手动调用GC后的内存占用
	runtime.GC()

	//	让主线程停止，等待GC回收器完成
	time.Sleep(5 * time.Second)

	runtime.ReadMemStats(&memStats)
	fmt.Println("调用GC后的内存占用：", memStats.Alloc/1024, "KB")
	fmt.Println("系统分配的内存(HeapSys)：", memStats.HeapSys/1024, "KB")
	fmt.Println("未使用的内存(HeapIdle)：", memStats.HeapIdle/1024, "KB")
	fmt.Println("GC的次数：", memStats.NumGC)

}

//func main() {
//	checkMem()
//}
