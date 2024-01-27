package internal

type JWTConfig struct {
	//	用于读取配置文件中的密钥
	SigningKey string `mapstructure:"key"`
}
