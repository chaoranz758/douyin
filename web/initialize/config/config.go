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
	DouYinWeb         `json:"douyinWeb"`
	Log               `json:"log"`
	ConsulServer      `json:"consulServer"`
	ConsulRegister    `json:"consulRegister"`
	ConsulCheckHealth `json:"consulCheckHealth"`
	RequestGRPCServer `json:"requestGRPCServer"`
	Rocketmq          `json:"rocketmq"`
	Jwt               `json:"jwt"`
	Oss               `json:"oss"`
}

type DouYinWeb struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
	Port int    `json:"port"`
}

type Log struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

type ConsulServer struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type ConsulRegister struct {
	ID   string `json:"id"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Name string `json:"name"`
	Tags string `json:"tags"`
}

type ConsulCheckHealth struct {
	IP                             string `json:"ip"`
	Port                           int    `json:"port"`
	Timeout                        string `json:"timeout"`
	Interval                       string `json:"interval"`
	DeregisterCriticalServiceAfter string `json:"deregisterCriticalServiceAfter"`
}

type RequestGRPCServer struct {
	UserService struct {
		Name string `json:"name"`
	} `json:"userService"`
	VideoService struct {
		Name string `json:"name"`
	} `json:"videoService"`
	CommitService struct {
		Name string `json:"name"`
	} `json:"commitService"`
	FavoriteService struct {
		Name string `json:"name"`
	} `json:"favoriteService"`
	FollowService struct {
		Name string `json:"name"`
	} `json:"followService"`
	MessageService struct {
		Name string `json:"name"`
	} `json:"messageService"`
}

type Rocketmq struct {
	Address string `json:"address"`
}

type Jwt struct {
	Secret string `json:"secret"`
}

type Oss struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
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
		fmt.Println("?????????????????????")
		if err2 := viper.Unmarshal(&Config); err2 != nil {
			fmt.Println("config unmarshal failed!")
		}
	})
	return nil
}
