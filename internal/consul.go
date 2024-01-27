package internal

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	//	根据实际情况转换数据类型（int/int32）
	Port int `mapstructure:"port"`
}

// Reg consul服务注册
func Reg(host, name, id string, port int, tags []string) error {
	//	通过api获取默认consul Config对象
	defaultConfig := api.DefaultConfig()
	h := AppConf.ConsulConfig.Host
	p := AppConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	//	通过config对象，获取consul客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	//	通过api.AgentServiceRegistration 配置服务的注册信息
	agentServiceRegistration := new(api.AgentServiceRegistration)

	serverAddr := fmt.Sprintf("http://%s:%d/health", host, port)
	//  定义服务健康检查配置
	check := api.AgentServiceCheck{
		HTTP:                           serverAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	//	配置服务相关信息
	agentServiceRegistration.Address = defaultConfig.Address
	agentServiceRegistration.Port = port
	agentServiceRegistration.Name = name
	agentServiceRegistration.ID = id
	agentServiceRegistration.Tags = tags
	//	配置服务的健康检查，需要传入一个对象的指针
	agentServiceRegistration.Check = &check
	/*	注册服务
		client.Agent()  获取consul agent(代理人)对象，可以使用agent的相关功能
	*/
	return client.Agent().ServiceRegister(agentServiceRegistration)

}

// GetServiceList 获取服务列表
func GetServiceList() error {
	defaultConfig := api.DefaultConfig()
	h := AppConf.ConsulConfig.Host
	p := AppConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	//	获取consul的客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	//	获取本地注册的服务
	serviceList, err := client.Agent().Services()
	if err != nil {
		return err
	}
	for k, v := range serviceList {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("=========================")
	}
	return nil
}

// FilterService 配置服务过滤器
func FilterService() error {
	//	通过api获取默认consul Config对象
	defaultConfig := api.DefaultConfig()
	h := AppConf.ConsulConfig.Host
	p := AppConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	//	通过config对象，获取consul客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	//	ServicesWithFilter() 获取特定服务名称的服务列表
	serviceList, err := client.Agent().ServicesWithFilter("Service==account_web")
	if err != nil {
		return err
	}
	for k, v := range serviceList {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("=========================")
	}
	return nil
}
