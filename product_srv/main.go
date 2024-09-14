package main

import (
	"fmt"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-trainning-lesson-part2/internal"
	"mic-trainning-lesson-part2/product_srv/biz"
	"mic-trainning-lesson-part2/proto/pb"
	"net"
)

func init() {
	internal.InitDB() // 初始化数据库
}

func main() {
	// 生成随机端口(使用环境变量)
	//port := util.GenRandomPort()
	//addr := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	port := internal.AppConf.ProductSrvConfig.Port
	addr := fmt.Sprintf("%s:%d", "192.168.150.11", port)
	fmt.Println("addr:", addr)
	//2. 使用 gRPC 中间件构建服务器，添加 Prometheus 拦截器

	server := grpc.NewServer(
		// 添加一元拦截器
		grpc.UnaryInterceptor( // 添加ChainUnaryServer 拦截器
			grpc_middleware.ChainUnaryServer(
				internal.GetGrpcMetrics().UnaryServerInterceptor(), // 连接 Prometheus 的 Unary 拦截器
			),
		),
		grpc.StreamInterceptor( // 添加 Stream 拦截器
			grpc_middleware.ChainStreamServer(
				internal.GetGrpcMetrics().StreamServerInterceptor(), // 连接 Prometheus 的 Stream 拦截器
			),
		),
	)
	//server := grpc.NewServer()
	// 3. 注册 gRPC 服务
	pb.RegisterProductServiceServer(server, &biz.ProductServer{}) // 注册产品服务

	// 4. 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer()) // 注册健康检查服务

	// 5. 创建监听器
	listen, err := net.Listen("tcp", addr) // 在指定地址上创建 TCP 监听器
	if err != nil {
		zap.S().Error("product_srv启动异常：" + err.Error())
		panic(err)
	}

	// 6. 注册 Prometheus 的 gRPC 服务器指标
	internal.InitMetrics(server) // 初始化 Prometheus 指标

	// 7. 启动 Prometheus 的 HTTP 服务，用于暴露指标数据
	internal.StartMetricsServer(9091)

	// 8. 将 gRPC 服务注册到 Consul
	defaultConfig := api.DefaultConfig() // 使用默认配置
	defaultConfig.Address = fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port) // 设置 Consul 地址

	client, err := api.NewClient(defaultConfig) // 创建 Consul 客户端
	if err != nil {
		panic(err)
	}

	checkAddr := fmt.Sprintf("%s:%d",
		internal.AppConf.ProductSrvConfig.Host,
		port) // 健康检查的地址

	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr, // 设置健康检查的 gRPC 地址
		Timeout:                        "3s",      // 健康检查超时设置
		Interval:                       "1s",      // 健康检查间隔设置
		DeregisterCriticalServiceAfter: "5s",      // 如果服务未响应，5秒后注销该服务
	}

	randUUID := uuid.New().String() // 生成随机的服务 ID

	reg := api.AgentServiceRegistration{
		Name:    internal.AppConf.ProductSrvConfig.SrvName, // 服务名称
		ID:      randUUID,                                  // 服务 ID
		Port:    port,                                      // 服务端口
		Tags:    internal.AppConf.ProductSrvConfig.Tags,    // 服务标签
		Address: internal.AppConf.ProductSrvConfig.Host,    // 服务地址
		Check:   check,                                     // 健康检查配置
	}

	fmt.Println(fmt.Sprintf("%s启动在%d", randUUID, port)) // 打印服务启动信息

	err = client.Agent().ServiceRegister(&reg) // 注册服务到 Consul
	if err != nil {
		panic(err)
	}

	// 9. 启动 gRPC 服务
	err = server.Serve(listen) // 启动 gRPC 服务器
	if err != nil {
		panic(err)
	}
}
