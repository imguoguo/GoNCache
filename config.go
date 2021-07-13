package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 网站缓存模型
type CacheListModel struct {
	Source string `yaml:"source"`
	IP string `yaml:"ip"`
	SeparateMobile bool `yaml:"separate_mobile"` // 是否将手机UA分开缓存
	ExpiredTime int `yaml:"expired_time"`
	NickName string `yaml:"nick_name"`
	Host []string `yaml:"host"`
	Scheme string `yaml:"scheme"`
	CacheBlackList []string `yaml:"cache_black_list"`
	TypeBlackList []string `yaml:"type_black_list"`
	VisitBlackList []string `yaml:"visit_black_list"`
}

// 配置文件格式
type Config struct {
	ListenAddress string `yaml:"listen_address"`
	ManagementAddress string `yaml:"management_address"`
	CacheList []CacheListModel `yaml:"cache_list"`
}

// 获取配置文件
// func (使用的结构体) 函数名(传入的参数) (返回的参数) {代码段}
func (conf *Config) getConf(filePath string) (error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return err
	}
	return nil
}