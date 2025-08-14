package config

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	Mode     int    `mapstructure:"mode" json:"mode" yaml:"mode"`             // 模式，0：单点，1：哨兵，2：集群
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
}
