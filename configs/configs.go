package configs

import "github.com/spf13/viper"

var cfg *conf

type conf struct {
	DBUser              string `mapstructure:"DB_USER"`
	DBPassword          string `mapstructure:"DB_PASSWORD"`
	NuvemshopAPIToken   string `mapstructure:"NUVEMSHOP_API_TOKEN"`
	NuvemshopStoreID    string `mapstructure:"NUVEMSHOP_STORE_ID"`
	NuvemshopAPIBaseURL string `mapstructure:"NUVEMSHOP_API_BASE_URL"`
	NuvemshopUserAgent  string `mapstructure:"NUVEMSHOP_USER_AGENT"`
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return cfg, nil
}
