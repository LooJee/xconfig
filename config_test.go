package xconfig

import (
	"testing"
	"time"
)

type Db struct {
	Postgres struct {
		Host  string
		Port  uint
		Debug bool
	}
}

type App struct {
	Grpc struct {
		Name string
		Port string
	}
}

type Config struct {
	Db  Db  `xconfig:"appId:db-config"`
	App App `xconfig:"appId:commodity-srv"`
}

func TestParseConfig(t *testing.T) {
	cfg := Config{}
	if err := NewConfig(&cfg, ApolloIp("apollo.api.thingyouwe.com"), LocalConfig("local.yaml")); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", cfg)
}

func TestEventChange(t *testing.T) {
	type Config struct {
		Db struct {
			Postgres struct {
				Host  string
				Port  int
				Debug bool
			}
		} `xconfig:"appId:db-config"`

		App struct {
			Grpc struct {
				Name  string
				Port  string
				Hello int `xconfig:"default:10"`
			}
		} `xconfig:"appId:commodity-srv"`
	}

	cfg := &Config{}

	if err := NewConfig(cfg, ApolloIp("apollo.api.thingyouwe.com"), LocalConfig("local.yaml")); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Minute)

	t.Logf("%+v", cfg)
}
