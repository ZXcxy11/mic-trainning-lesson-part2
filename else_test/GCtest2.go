package main

import (
	"fmt"
	"runtime"
)

// 模拟对象的创建和联系，观察垃圾回收算法
// 模拟对象结构体
type Obj struct {
	id int
	//	指向另一个对象的引用
	relate *Obj
}

func testStart() {
	//	对象切片
	objs := make([]*Obj, 0)

	//	模拟多个对象的创建和引用
	for i := 0; i < 1000; i++ {
		obj := &Obj{id: i}
		objs = append(objs, obj)

		//	创建循环引用
		if i > 0 {
			//	每一个对象引用上一个对象
			obj.relate = objs[i-1]
		}
	}
	//	打印当前内存情况
	printMemStats("分配对象前")

	//	释放一部分对象的引用，模拟垃圾
	for j := 0; j < 250; j++ {
		//	断开前250个对象的引用
		objs[j].relate = nil
	}

	//	手动触发垃圾回收
	runtime.GC()

	//	再次查看内存状态
	printMemStats("第一次删除引用后：")

}
func printMemStats(stage string) {
	var memStats runtime.MemStats
	//	读取内存状态
	runtime.ReadMemStats(&memStats)
	fmt.Printf("%s:\n", stage)
	//	当前内存占用
	fmt.Printf("当前内存占用 Alloc = %v KB\n", memStats.Alloc/1024)
	//	自启动以来的总内存分配数
	fmt.Printf("总内存分配数 TotalAlloc = %v KB\n", memStats.TotalAlloc/1024)
	//	从操作系统获得的总内存量
	fmt.Printf("获得的总内存量 Sys = %v KB\n", memStats.Sys/1024)
	//	GC执行次数
	fmt.Printf("GC执行次数 NumGC = %v\n\n", memStats.NumGC)
}

//func main() {
//	testStart()
//}
