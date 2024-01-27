package internal

type NacosConfig struct {
	Host                string `mapstructure:"host"`
	Port                uint64 `mapstructure:"port"`
	Namespace           string `mapstructure:"namespace"`
	DataId              string `mapstructure:"dataId"`
	Group               string `mapstructure:"group"`
	TimeoutMS           uint64 `mapstructure:"timeoutMs"`
	NotLoadCacheAtStart bool   `mapstructure:"notLoadCacheAtStart"`
	LogDir              string `mapstructure:"logDir"`
	CacheDir            string `mapstructure:"cacheDir"`
	LogLevel            string `mapstructure:"logLevel"`
}
