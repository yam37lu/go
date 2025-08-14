package config

type Server struct {
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Gorm    Gorm    `mapstructure:"gorm" json:"gorm" yaml:"gorm"`
	License License `mapstructure:"license" json:"license" yaml:"license"`
}
