package internal

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

//	go grpc项目集成Prometheus

// 声明变量，存储grpc服务器指标（Metrics）
var grpcMetrics *grpc_prometheus.ServerMetrics

// InitMetrics 函数，初始化 Prometheus指标
func InitMetrics(server *grpc.Server) {
	//	启用处理时间直方图(HandlingTimeHistogram)
	grpcMetrics.EnableHandlingTimeHistogram()
	//	将指标注册到grpc服务器
	grpcMetrics.InitializeMetrics(server)
}

// StartMetricsServer 启动 Prometheus 服务
func StartMetricsServer(port int) {
	go func() {
		//	将 /Metrics 路径与 prometheus handler 关联（常用写法）
		http.Handle("/metrics", promhttp.Handler())
		//	打印 Prometheus 指标的可用地址信息
		fmt.Printf("Prometheus 指标的可用地址信息：%d\n", port)
		//	启动Http服务，监听指定的端口
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			zap.S().Error(err.Error())
		}
	}()
}

// GetGrpcMetrics 返回gRpc的 prometheus 指标
func GetGrpcMetrics() *grpc_prometheus.ServerMetrics {
	//	创建指标
	grpcMetrics = grpc_prometheus.NewServerMetrics()
	return grpcMetrics // 返回 gRPC 服务器的指标实例
}
