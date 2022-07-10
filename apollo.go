package xconfig

import (
	"fmt"
	"strings"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
)

type Client struct {
	client          agollo.Client
	Config          *config.AppConfig
	ChangeListeners []storage.ChangeListener
}

func NewClient(opts ...Option) (*Client, error) {
	var client = &Client{
		Config: &config.AppConfig{
			AppID:            "db-config",
			Cluster:          "default",
			NamespaceName:    storage.GetDefaultNamespace(),
			IP:               "http://apollo.api.thingyouwe.com",
			MustStart:        true,
			Secret:           "",
			IsBackupConfig:   true,
			BackupConfigPath: ".apollo",
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	var err error
	client.client, err = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return client.Config, nil
	})
	if err != nil {
		return nil, err
	}

	for _, listener := range client.ChangeListeners {
		client.client.AddChangeListener(listener)
	}

	return client, nil
}

func (op *Client) Init(aps *appSetter) (err error) {
	if err := aps.Reset(); err != nil {
		return err
	}

	op.client.GetDefaultConfigCache().Range(func(key, value interface{}) bool {
		if err = aps.SetValue(key.(string), value); err != nil {
			return false
		}

		return true
	})

	return nil
}

func checkKey(namespace string, client agollo.Client) {
	cache := client.GetConfigCache(namespace)
	count := 0
	cache.Range(func(key, value interface{}) bool {
		fmt.Println("key : ", key, ", value :", value)
		count++
		return true
	})
	if count < 1 {
		panic("config key can not be null")
	}
}

type Option func(option *Client)

func AppIp(appId string) Option {
	return func(option *Client) {
		option.Config.AppID = appId
	}
}

func Cluster(cluster string) Option {
	return func(option *Client) {
		if len(cluster) == 0 {
			cluster = "default"
		}
		option.Config.Cluster = cluster
	}
}

func Namespace(namespace string) Option {
	return func(option *Client) {
		option.Config.NamespaceName = namespace
	}
}

func IP(ip string) Option {
	return func(option *Client) {
		if len(ip) == 0 {
			ip = "apollo.api.thingyouwe.com"
		}

		if !strings.HasPrefix(ip, "http://") || !strings.HasPrefix(ip, "https://") {
			ip = "http://" + ip
		}

		option.Config.IP = ip
	}
}

func Secret(secret string) Option {
	return func(option *Client) {
		option.Config.Secret = secret
	}
}

func BackupConfigPath(path string) Option {
	return func(option *Client) {
		option.Config.IsBackupConfig = true
		option.Config.BackupConfigPath = path
	}
}

func ChangeListener(listener ...storage.ChangeListener) Option {
	return func(option *Client) {
		if len(listener) == 0 {
			return
		}
		option.ChangeListeners = append(option.ChangeListeners, listener...)
	}
}
