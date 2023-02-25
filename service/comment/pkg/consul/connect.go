package consul

import "github.com/hashicorp/consul/api"

var Consul *api.Client

//连接consul服务

func Init(addr string) error {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	Consul = c
	return nil
}
