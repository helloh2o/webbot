package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"robot/bean"
	"robot/consts"
)

type conf struct {
	// target host
	TargetURL string `yaml:"TargetURL"` // 目标主机
	// Concurrent
	Concurrent     int `yaml:"Concurrent"` // 并发数量
	ConcurrentTick int `yaml:"ConcurrentTick"`
	// Timeout 超时(s) // 运行时间
	Timeout int `yaml:"Timeout"`
	// 动作有先后
	Actions []*bean.Action `yaml:"Actions"`
	// Total Request
	TotalRequest int
}

var c *conf

func Load(filename string) {
	c = &conf{}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		panic(err)
	} else if err = yaml.Unmarshal(yamlFile, c); err != nil {
		panic(err)
	}
	for _, act := range c.Actions {
		if act.Name != consts.Captcha {
			c.TotalRequest += c.Concurrent * act.Times
		}
		act.Url = c.TargetURL + act.Url
	}
}

func Get() *conf {
	return c
}
