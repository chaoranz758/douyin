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
	RequestGRPCServer `json:"requestGRPCServer"`
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

type RequestGRPCServer struct {
	UserService struct {
		Name string `json:"name"`
	} `json:"userService"`
	VideoService struct {
		Name string `json:"name"`
	} `json:"videoService"`
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

//??????viper??????????????????

func Init(filePath string) error {
	//????????????????????????
	viper.SetConfigFile(filePath)
	//??????????????????
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	//??????????????????
	if err1 := viper.ReadConfig(file); err1 != nil {
		return err1
	}
	if err2 := viper.Unmarshal(&configViper); err2 != nil {
		return err2
	}
	// ????????????ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: configViper.Nacos.IP,
			Port:   uint64(configViper.Nacos.Port),
		},
	}
	// ??????clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         configViper.Nacos.Namespace, // ?????????????????????namespace???????????????????????????client,??????????????????NamespaceId??????namespace???public??????????????????????????????
		TimeoutMs:           uint64(configViper.Nacos.TimeoutMs),
		NotLoadCacheAtStart: configViper.Nacos.NotLoadCacheAtStart,
		LogDir:              configViper.Nacos.LogDir,
		CacheDir:            configViper.Nacos.CacheDir,
		LogLevel:            configViper.Nacos.LogLevel,
	}
	// ????????????????????????????????????????????? (??????)
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
	//??????????????????
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: configViper.Nacos.DataID,
		Group:  configViper.Nacos.Group})
	if err != nil {
		fmt.Println("GetConfig err: ", err)
		return err
	}
	if err := json.Unmarshal([]byte(content), &Config); err != nil {
		log.Fatal("??????????????? ", err)
		return err
	}
	//fmt.Printf("%v\n", Config)
	//????????????
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: configViper.Nacos.DataID,
		Group:  configViper.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("????????????????????????????????????")
			fmt.Println("group:" + group + ", dataId:" + dataId)
		},
	})
	if err != nil {
		return err
	}
	//????????????????????????
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// ??????????????????????????????????????????????????????
		fmt.Println("???????????????????????????")
		if err2 := viper.Unmarshal(&Config); err2 != nil {
			fmt.Println("config unmarshal failed!")
		}
	})
	return nil
}
