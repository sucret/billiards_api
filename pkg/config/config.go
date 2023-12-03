package config

type Config struct {
	App            App              `mapstructure:"app" json:"app" yaml:"app"`
	Log            Log              `mapstructure:"log" json:"log" yaml:"log"`
	Database       Database         `mapstructure:"database" json:"database" yaml:"database"`
	Jwt            Jwt              `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Qiniu          Qiniu            `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
	Redis          Redis            `mapstructure:"redis" json:"redis" yaml:"redis"`
	WeApp          WeApp            `mapstructure:"weapp" json:"weapp" yaml:"weapp"`
	RechargeAmount []RechargeAmount `mapstructure:"recharge_amount" json:"recharge_amount" yaml:"recharge_amount"`
}
