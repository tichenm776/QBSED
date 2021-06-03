package configs

import (
	"io/ioutil"
	"github.com/alecthomas/log4go"
	"gopkg.in/yaml.v2"
	"errors"
)

type Confs struct {
	ServerPort         int    `yaml:"server_port"` // 服务端口
	ServerPort_theme   int    `yaml:"server_port_theme"`//主题服务端口
	KoalaHost          string `yaml:"koala_host"`
	KoalaPort          int    `yaml:"koala_port"`
	KoalaCorePort      int    `yaml:"koala_core_port"`
	DBUsername         string `yaml:"db_username"`
	DBPassword         string `yaml:"db_password"`
	DBIP               string `yaml:"db_ip"`
	DBPort             int    `yaml:"db_port"`
	DBName             string `yaml:"db_name"`
	KoalaApisHost      string `yaml:"koala_apis_host"`
	KoalaUsername      string `yaml:"koala_username"`
	KoalaPassword      string `yaml:"koala_password"`
	ArmPath           string `yaml:"arm_path"`
	Amd64Path         string `yaml:"amd64_path"`
	Amd64Fullpath     string `yaml:"amd64_fullpath"`
}

// 配置参数全局变量

var Gconf Confs

func Init(file string) {
	_, err := Gconf.LoadConfig(file)
	if err != nil {
		log4go.Crash("读取配置文件失败！")
	}

}

// LoadConig 读取配置文件
func (c *Confs) LoadConfig(file string) (*Confs, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log4go.Info("打开 conf.yaml 失败: %v \n", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log4go.Error("yaml unmarshal error: %v \n", err)
		return nil, err
	}

	return c, nil
}


func (c *Confs)EditYaml(username string, password string) (*Confs, error){
	c.KoalaUsername = username
	c.KoalaPassword = password
	yaml_data, err := yaml.Marshal(&c)
	if err != nil {
		log4go.Error(err)
		return c, errors.New("生成yaml失败!")
	}
	err1 := ioutil.WriteFile("conf.yaml", yaml_data, 0644)
	if err1 != nil {
		log4go.Error(err1)
		return c, errors.New("写入yaml失败!")
	}
	return c, nil
}
