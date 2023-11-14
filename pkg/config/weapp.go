package config

type WeApp struct {
	AppId     string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`
	AppSecret string `mapstructure:"app_secret" json:"app_secret" yaml:"app_secret"`
}
