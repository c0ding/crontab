package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 配置 1
type Config struct {
	ApiPort         int      `json:"api_port"`
	ApiReadTimeout  int      `json:"api_read_timeout"`
	ApiWriteTimeout int      `json:"api_write_timeout"`
	EtcdEndpoints   []string `json:"etcd_endpoints"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout"`
}

var (
	G_config *Config
)

// 加载配置 2
func InitConfig(filename string) (err error) {
	var (
		fileContent []byte
		config      Config
	)

	// 把配置文件读进来 1 ，得到一个二进制数组
	if fileContent, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 做json反序列化 2 ，把二进制数组 配置到Config 结构体中
	if err = json.Unmarshal(fileContent, &config); err != nil {
		return
	}

	// 赋值单例 3
	G_config = &config

	fmt.Println(config)
	return
}
