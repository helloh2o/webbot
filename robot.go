package main

import (
	"flag"
	"log"
	"robot/bean"
	"robot/conf"
	"robot/core"
	"sync"
	"time"
)

var (
	confPath = flag.String("conf", "./config.yaml", "the config file path")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	conf.Load(*confPath)
	c := conf.Get()
	core.ReadAccounts()
	log.Printf("Concurrent User %d, Total Request %d", c.Concurrent, c.TotalRequest)
	var wg *sync.WaitGroup
	if c.Timeout == 0 {
		wg = new(sync.WaitGroup)
		wg.Add(c.TotalRequest)
	}
	// 开启请求
	for i := 0; i < c.Concurrent; i++ {
		r := core.NewRobot(i+1, make(chan *bean.Action), wg)
		r.SetActions(c.Actions)
		r.Run()
		if c.ConcurrentTick != 0 {
			time.Sleep(time.Millisecond * time.Duration(c.ConcurrentTick))
		}
	}
	start := time.Now().UnixNano()
	// 超时
	if c.Timeout != 0 {
		timeout := time.Tick(time.Second * time.Duration(c.Timeout))
		<-timeout
	} else if wg != nil {
		wg.Wait()
	}
	log.Println("========================= TEST =========================")
	log.Printf("并发用户 %d, 总请求次数 %d, 在线时间 %s", c.Concurrent, c.TotalRequest, core.GetTime2Time(start))
	for _, act := range c.Actions {
		if act.Times == 0 {
			continue
		}
		log.Printf("Action %s %d times, url=> %s", act.Name, act.Times*c.Concurrent, act.Url)
	}
	core.Done()
}
