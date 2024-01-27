package internal

type ProductSrvConfig struct {
	SrvName string   `mapstruct:"srvName" json:"srvName"`
	Host    string   `mapstrcuture:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tage" json:"tags"`
}

type ProductWebConfig struct {
	SrvName string   `mapstruct:"srvName" json:"srvName"`
	Host    string   `mapstrcuture:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tage" json:"tags"`
}

type AppConfig struct {
	RedisConfig      RedisConfig      `mapstructure:"redis" json:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul" json:"consul"`
	ProductSrvConfig ProductSrvConfig `mapstructure:"product_srv" json:"product_srv"`
	ProductWebConfig ProductWebConfig `mapstructure:"product_web" json:"product_web"`
	DBConfig         DBConfig         `mapstructure:"db" json:"db"`
	JWTConfig        JWTConfig        `mapstructure:"jwt" json:"jwt"`
}
