package config

type Mongo struct {
	URI         string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Username    string `mapstructure:"username" json:"username" yaml:"username"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	MinPoolSize int    `mapstructure:"min-pool-size" json:"minpoolsize" yaml:"min-pool-size"`
	MaxPoolSize int    `mapstructure:"max-pool-size" json:"maxpoolsize" yaml:"max-pool-size"`
}
