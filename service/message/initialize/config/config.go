package config

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	configViper config
	Config      config1
)

type config struct {
	Nacos
}

type config1 struct {
	DouYinService     `json:"douyinService"`
	Log               `json:"log"`
	Mysql             `json:"mysql"`
	Redis             `json:"redis"`
	SnowFlake         `json:"snowflake"`
	ConsulServer      `json:"consulServer"`
	ConsulRegister    `json:"consulRegister"`
	ConsulCheckHealth `json:"consulCheckHealth"`
	Rocketmq          `json:"rocketmq"`
}

type DouYinService struct {
	Name     string `mapstructure:"name" json:"name"`
	Mode     string `mapstructure:"mode" json:"mode"`
	Port     int    `mapstructure:"port" json:"port"`
	Protocol string `mapstructure:"protocol" json:"protocol"`
}

type Log struct {
	Level      string `mapstructure:"level" json:"level"`
	Filename   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
}

type Mysql struct {
	Host         string `mapstructure:"host" json:"host"`
	Port         int    `mapstructure:"port" json:"port"`
	User         string `mapstructure:"user" json:"user"`
	Password     string `mapstructure:"password" json:"password"`
	Dbname       string `mapstructure:"dbname" json:"dbname"`
	MaxIdleConns int    `mapstructure:"MaxIdleConns" json:"MaxIdleConns"`
	MaxOpenConns int    `mapstructure:"MaxOpenConns" json:"MaxOpenConns"`
}

type Redis struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	Db       int    `mapstructure:"db" json:"db"`
	PoolSize int    `mapstructure:"PoolSize" json:"PoolSize"`
}

type SnowFlake struct {
	StartTime string `mapstructure:"starttime" json:"starttime"`
	MachineID int64  `mapstructure:"machineID" json:"machineID"`
}

type ConsulServer struct {
	Ip   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulRegister struct {
	ID   string `mapstructure:"id" json:"id"`
	IP   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
	Tags string `mapstructure:"tags" json:"tags"`
}

type ConsulCheckHealth struct {
	IP                             string `mapstructure:"ip" json:"ip"`
	Port                           int    `mapstructure:"port" json:"port"`
	Timeout                        string `json:"timeout" json:"timeout"`
	Interval                       string `json:"interval" json:"interval"`
	DeregisterCriticalServiceAfter string `json:"deregisterCriticalServiceAfter" json:"deregisterCriticalServiceAfter"`
}

type Rocketmq struct {
	Address string `json:"address"`
}

type Nacos struct {
	IP                  string `mapstructure:"ip" json:"ip"`
	Port                int    `mapstructure:"port" json:"port"`
	Namespace           string `mapstructure:"namespace" json:"namespace"`
	TimeoutMs           int    `mapstructure:"timeoutMs" json:"timeoutMs"`
	NotLoadCacheAtStart bool   `mapstructure:"notLoadCacheAtStart" json:"notLoadCacheAtStart"`
	LogDir              string `mapstructure:"logDir" json:"logDir"`
	CacheDir            string `mapstructure:"cacheDir" json:"cacheDir"`
	LogLevel            string `mapstructure:"logLevel" json:"logLevel"`
	DataID              string `mapstructure:"dataID" json:"dataID"`
	Group               string `mapstructure:"group" json:"group"`
}

//使用viper管理配置文件

func Init(filePath string) error {
	//指定配置文件路径
	viper.SetConfigFile(filePath)
	//打开配置文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	//读取配置文件
	if err1 := viper.ReadConfig(file); err1 != nil {
		return err1
	}
	if err2 := viper.Unmarshal(&configViper); err2 != nil {
		return err2
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: configViper.Nacos.IP,
			Port:   uint64(configViper.Nacos.Port),
		},
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         configViper.Nacos.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           uint64(configViper.Nacos.TimeoutMs),
		NotLoadCacheAtStart: configViper.Nacos.NotLoadCacheAtStart,
		LogDir:              configViper.Nacos.LogDir,
		CacheDir:            configViper.Nacos.CacheDir,
		LogLevel:            configViper.Nacos.LogLevel,
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: configViper.Nacos.DataID,
		Group:  configViper.Nacos.Group})
	if err != nil {
		fmt.Println("GetConfig err: ", err)
		return err
	}
	if err := json.Unmarshal([]byte(content), &Config); err != nil {
		log.Fatal("解释失败： ", err)
		return err
	}
	//fmt.Printf("%v\n", Config)
	//监听配置
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: configViper.Nacos.DataID,
		Group:  configViper.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置中心的配置文件修改了")
			fmt.Println("group:" + group + ", dataId:" + dataId)
		},
	})
	if err != nil {
		return err
	}
	//监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("本地配置文件修改了")
		if err2 := viper.Unmarshal(&Config); err2 != nil {
			fmt.Println("config unmarshal failed!")
		}
	})
	return nil
}
