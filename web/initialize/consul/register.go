package consul

import (
	"douyin/web/initialize/config"
	"fmt"
	"github.com/hashicorp/consul/api"
)

func RegisterService() error {
	var check = api.AgentServiceCheck{
		HTTP: fmt.Sprintf("http://%s:%d/checkHealth", config.Config.ConsulCheckHealth.IP, config.Config.ConsulCheckHealth.Port),
		//超时时间
		Timeout: config.Config.ConsulCheckHealth.Timeout,
		//运行检查的频率
		Interval: config.Config.ConsulCheckHealth.Interval,
		//指定时间后自动注销不健康的服务节点
		DeregisterCriticalServiceAfter: config.Config.ConsulCheckHealth.DeregisterCriticalServiceAfter,
	}
	var srv = api.AgentServiceRegistration{
		ID:      config.Config.ConsulRegister.ID,
		Name:    config.Config.ConsulRegister.Name,
		Tags:    []string{config.Config.ConsulRegister.Tags},
		Address: config.Config.ConsulRegister.IP,
		Port:    config.Config.ConsulRegister.Port,
		Check:   &check,
	}
	if err := Consul.Agent().ServiceRegister(&srv); err != nil {
		return err
	}
	return nil
}
