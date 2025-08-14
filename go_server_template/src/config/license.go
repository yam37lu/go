package config

type License struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	Status string `mapstructure:"status" json:"status" yaml:"status"`
	Module string `mapstructure:"module" json:"module" yaml:"module"`
}
