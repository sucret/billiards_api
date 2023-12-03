package config

type RechargeAmount struct {
	Amount        int32 `mapstructure:"amount" json:"amount" yaml:"amount"`
	BundledAmount int32 `mapstructure:"bundled_amount" json:"bundled_amount" yaml:"bundled_amount"`
}
