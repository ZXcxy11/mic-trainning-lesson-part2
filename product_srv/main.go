package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/product_srv/biz"
	"mic-trainning-lesson-part2/proto/pb"
	"mic-trainning-lesson-part2/util"
	"net"
)

// todo grpc服务器
func init() {
	internal.InitDB()
}

/*
*
快速启动srv服务
1. 生成proto对应文件
2. 生成biz文件，实现接口
3. 拷贝先前例子的main函数
*/
func main() {

	/*
		flag包是用于处理命令行参数的
		flag.String()用于处理字符串类型的参数，返回值类型是该类型值的指针，参数1：命令名称，参数2：默认值，参数3：参数备注
		...
	*/

	//	使用随机端口改造grpc服务器地址
	port := util.GenRandomPort()
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ProductSrvConfig.Host, port)
	fmt.Println("addr:", addr)

	//	获取grpc服务器
	server := grpc.NewServer()
	//	注册服务器
	pb.RegisterProductServiceServer(server, &biz.ProductServer{})
	//	获取监听器
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("product_srv启动异常：" + err.Error())
		panic(err)
	}
	//	todo 将GRPC服务注册到consul
	// 为当前grpc服务添加健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	//	获取consul配置
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port)
	//	获取consul客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	//	设置检查路由
	checkAddr := fmt.Sprintf("%s:%d",
		internal.AppConf.ProductSrvConfig.Host,
		port)
	//  配置健康检查信息
	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	randUUID := uuid.New().String()
	//	注册代理服务
	reg := api.AgentServiceRegistration{
		Name:    internal.AppConf.ProductSrvConfig.SrvName,
		ID:      randUUID,
		Port:    port,
		Tags:    internal.AppConf.ProductSrvConfig.Tags,
		Address: internal.AppConf.ProductSrvConfig.Host,
		Check:   check,
	}
	fmt.Println(fmt.Sprintf("%s启动在%d", randUUID, port))
	//	consul客户端调用代理服务
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
