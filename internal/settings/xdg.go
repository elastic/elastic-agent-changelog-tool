package settings

import (
	"github.com/spf13/viper"
)

func CacheDir() string {
	return viper.GetString("cache_dir")
}

func ConfigDir() string {
	return viper.GetString("config_dir")
}

func DataDir() string {
	return viper.GetString("data_dir")
}
