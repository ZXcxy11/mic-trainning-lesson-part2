package register

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"mic-trainning-lesson-part2/internal"
)

// IRegister 接口定义了服务注册和注销的行为
type IRegister interface {
	// Register 注册服务到服务发现系统
	// name: 服务名称
	// id: 服务ID
	// port: 服务监听的端口
	// tags: 服务的标签列表
	// 返回错误信息，如果注册失败
	Register(name, id string, port int, tags []string) error

	// DeRegister 注销服务
	// serviceId: 要注销的服务ID
	// 返回错误信息，如果注销失败
	DeRegister(serviceId string) error
}

// ConsulRegistry 是一个使用 Consul 作为服务注册中心的实现
type ConsulRegistry struct {
	Host string // Consul 服务器的主机地址
	Port int    // Consul 服务器的端口号
}

// NewConsulRegistry 是 ConsulRegistry 的构造函数
// 它接收主机和端口参数，并返回一个初始化后的 ConsulRegistry 实例
func NewConsulRegistry(host string, port int) ConsulRegistry {
	return ConsulRegistry{
		Host: host,
		Port: port,
	}
}

// Register 方法实现了 IRegister 接口的 Register 方法
// 将服务注册到 Consul 服务注册中心
func (cr ConsulRegistry) Register(name, id string, port int, tags []string) error {
	// 获取 Consul 的默认配置
	defaultConfig := api.DefaultConfig()

	// 从应用配置中读取 Consul 的 Host 和 Port
	h := internal.AppConf.ConsulConfig.Host
	p := internal.AppConf.ConsulConfig.Port

	// 设置 Consul 的地址为从配置中读取的 Host 和 Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)

	// 创建一个新的 Consul 客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		// 如果创建客户端失败，记录错误日志并返回错误
		zap.S().Error(err)
		return err
	}

	// 创建一个新的服务注册实例
	agentServiceReg := new(api.AgentServiceRegistration)

	// 设置服务注册信息
	agentServiceReg.Port = port                     // 服务监听的端口
	agentServiceReg.Address = defaultConfig.Address // 服务的地址
	agentServiceReg.ID = id                         // 服务的ID
	agentServiceReg.Name = name                     // 服务的名称
	agentServiceReg.Tags = tags                     // 服务的标签

	// 设置服务的健康检查配置
	serverAddr := fmt.Sprintf("http://%s:%d/health", internal.AppConf.ProductWebConfig.Host,
		internal.AppConf.ProductWebConfig.Port)
	check := api.AgentServiceCheck{
		HTTP:                           serverAddr, // 健康检查的 HTTP 地址
		Timeout:                        "3s",       // 健康检查的超时时间
		Interval:                       "1s",       // 健康检查的间隔时间
		DeregisterCriticalServiceAfter: "5s",       // 在服务不健康时，自动注销服务的时间
	}
	agentServiceReg.Check = &check

	// 向 Consul 注册服务
	err = client.Agent().ServiceRegister(agentServiceReg)
	if err != nil {
		// 如果注册失败，记录错误日志并返回错误
		zap.S().Error(err)
		return err
	}

	// 如果服务注册成功，返回 nil 表示无错误
	return nil
}

// DeRegister 方法实现了 IRegister 接口的 DeRegister 方法
// 从 Consul 服务注册中心注销服务
func (cr ConsulRegistry) DeRegister(serviceId string) error {
	// 获取 Consul 的默认配置
	defaultConfig := api.DefaultConfig()

	// 从应用配置中读取 Consul 的 Host 和 Port
	h := internal.AppConf.ConsulConfig.Host
	p := internal.AppConf.ConsulConfig.Port

	// 设置 Consul 的地址为从配置中读取的 Host 和 Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)

	// 创建一个新的 Consul 客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		// 如果创建客户端失败，记录错误日志并返回错误
		zap.S().Error(err)
		return err
	}

	// 调用 Consul 客户端的 ServiceDeregister 方法注销服务
	return client.Agent().ServiceDeregister(serviceId)
}
