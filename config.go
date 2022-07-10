package xconfig

func NewConfig(obj interface{}, opts ...ConfigOption) error {
	cfg := &innerCfg{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}

	//解析配置对象
	apss, err := newParser().parse(obj)
	if err != nil {
		return err
	}

	//启动apollo
	for key, aps := range apss {
		if cli, err := NewClient(
			AppIp(key),
			IP(cfg.ApolloIp),
			Cluster(cfg.ApolloCluster),
			ChangeListener(NewWatcher(aps, key))); err != nil {
			return err
		} else if err := cli.Init(aps); err != nil {
			return err
		}
	}

	return nil
}

type innerCfg struct {
	ApolloIp      string
	ApolloCluster string
	LocalConfig   string
}

type ConfigOption func(cfg *innerCfg) error

func ApolloIp(ip string) ConfigOption {
	return func(cfg *innerCfg) error {
		cfg.ApolloIp = ip

		return nil
	}
}

func ApolloCluster(cluster string) ConfigOption {
	return func(cfg *innerCfg) error {
		if len(cluster) == 0 {
			cluster = "default"
		}
		cfg.ApolloCluster = cluster

		return nil
	}
}

func LocalConfig(lc string) ConfigOption {
	return func(cfg *innerCfg) error {
		return loadLocalConfig(lc)
	}
}
