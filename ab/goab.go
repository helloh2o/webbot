package main

import (
	"flag"
	"log"
	"net/url"
	"robot/bean"
	"robot/core"
	"sync"
	"time"
)

var (
	u  = flag.String("u", "http://192.168.1.201:8082/api/topic/node/news?nodeId=1", "the url to request")
	c  = flag.Int("c", 10, "concurrent size")
	n  = flag.Int("n", 100, "the number of request times")
	tp = flag.Int("tp", 5, "timeout (s) for per request")
	t  = flag.Int("t", 0, "timeout (s) for all the request")
)

func main() {
	flag.Parse()
	_, err := url.Parse(*u)
	if err != nil {
		panic(err)
	}
	log.Println("Total Request::", *n)
	log.Println("Concurrent Request::", *c)
	reqChan := make(chan *bean.Action)
	var wg *sync.WaitGroup
	if *t == 0 {
		wg = new(sync.WaitGroup)
		wg.Add(*n)
	}
	// 开启请求
	for i := 0; i < *c; i++ {
		r := core.NewRobot(i+1, reqChan, wg)
		r.SetTimeoutReq(*tp)
		r.Run()
	}
	start := time.Now().UnixNano()
	// 发射请求
	go func() {
		for i := 0; i < *n; i++ {
			reqChan <- &bean.Action{
				Method: "GET",
				Url:    *u,
			}
		}
	}()
	// 超时
	if *t == 0 {
		wg.Wait() // 等待结束
	} else {
		timeout := time.Tick(time.Second * time.Duration(*t))
		<-timeout
	}
	log.Println("=======================================================")
	log.Printf("Concurrent %d, Requests %d, Total cost %s", *c, *n, core.GetTime2Time(start))
	core.Done()
}
