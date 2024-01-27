package internal

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// todo 配置viper，用于读取配置文件(经过nacos改造)
var AppConf AppConfig
var NacosConf NacosConfig

// ViperConf用于存储从配置文件中读取的各个信息，并作为全局变量使用
// var ViperConf ViperConfig
var fileName = "./dev-config.yaml"

// Nacos的初始化
func initNacos() {
	//	获取viper对象
	v := viper.New()
	//	设置配置文件路径
	v.SetConfigFile(fileName)
	//	导入配置文件
	v.ReadInConfig()
	//	读取配置文件内容
	v.Unmarshal(&NacosConf)
	fmt.Println("nacos数据：", NacosConf)

}

// 从Nacos中获取数据
func initFromNacos() {
	//	声明nacos服务器配置信息
	serverConfigs := []constant.ServerConfig{
		{
			//	当前主机的IP地址和端口
			IpAddr: NacosConf.Host,
			Port:   NacosConf.Port,
		},
	}
	//	声明nacos客户端配置
	clientConfig := constant.ClientConfig{
		//	配置的命名空间ID
		//NamespaceId: "c6bff0ef-3e55-4fa8-931a-8afe66911f0e",
		NamespaceId: NacosConf.Namespace,
		//	超时时间
		TimeoutMs: NacosConf.TimeoutMS,
		//	是否从缓存获取
		NotLoadCacheAtStart: NacosConf.NotLoadCacheAtStart,
		//	日志保存目录
		LogDir: NacosConf.LogDir,
		//	缓存保存目录
		CacheDir: NacosConf.CacheDir,
		//	日志等级
		LogLevel: NacosConf.LogLevel,
	}
	//	获取nacos客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		//	注：键值不可随意更改，否则会获取失败，例："serverConfigs"不可少s
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	//	从配置中心nacos中获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		//	配置名称
		DataId: NacosConf.DataId,
		//	配置所在组
		Group: NacosConf.Group,
	})
	if err != nil {
		panic(err)
	}
	//	打印配置
	json.Unmarshal([]byte(content), &AppConf)
	fmt.Println("AppConf:", AppConf)
}

func init() {
	initNacos()
	initFromNacos()
	fmt.Println("初始化完成")
	InitRedis()
}
