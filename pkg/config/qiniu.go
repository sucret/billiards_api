package config

type Qiniu struct {
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Domain    string `mapstructure:"domain" json:"domain" yaml:"domain"`
}
