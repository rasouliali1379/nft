package config

type Etcd struct {
	Endpoints      []string `json:"endpoints" yaml:"etcd.endpoints" required:"true"`
	WatchList      []string `json:"watch_list" yaml:"etcd.watchList"`
	ProviderPrefix string   `json:"provider_prefix" yaml:"providerprefix"`
	Username       string   `json:"username" yaml:"etcd.username" required:"true"`
	Password       string   `json:"password" yaml:"etcd.password" required:"true"`
}
