package settings

import (
	"github.com/OpenPeeDeeP/xdg"
	"github.com/spf13/viper"
)

// Init initalize settings and default values
func Init() {
	viper.AutomaticEnv()
	viper.BindEnv("ELASTIC_AGENT_CHANGELOG")

	setDefaults()
	setConstants()
}

func setDefaults() {
	viper.SetDefault("cache_dir", xdg.CacheHome())
	viper.SetDefault("config_dir", xdg.ConfigHome())
	viper.SetDefault("data_dir", xdg.DataHome())
}

func setConstants() {
	// viper.Set()
}
