package xconfig

import (
	"testing"

	"github.com/spf13/viper"
)

func TestLocalConfig(t *testing.T) {
	if err := loadLocalConfig("local.yaml"); err != nil {
		t.Fatal(err)
	}

	t.Log(viper.AllKeys())

	t.Log(viper.Get("db-config"))
}
