package main

import "fmt"

// foo 指针逃逸
func foo() *int {
	x := 42
	return &x
}

// 栈空间不足逃逸
func largeStack() {
	var largeArray [1000000]int
	fmt.Println(largeArray)
}

// 动态类型逃逸
func dynamicEscape() interface{} {
	x := 23
	return x // x 逃逸到堆上，因为返回了接口类型
}

// 闭包引用对象逃逸
func closeureEscape() func() {
	x := 100
	return func() {
		fmt.Println(x)
	}
}
