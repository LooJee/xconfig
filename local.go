package xconfig

import (
	"strings"

	"github.com/spf13/viper"
)

func loadLocalConfig(file string) error {
	if len(file) == 0 {
		return nil
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func hasLocalConfig(appId, key string) bool {
	return viper.IsSet(strings.Join([]string{appId, key}, "."))
}

func getLocalConfig(appId, key string) interface{} {
	return viper.Get(strings.Join([]string{appId, key}, "."))
}
