package core

import (
	"log"
	"sync/atomic"
)

var (
	Req_sent    int64
	Req_succeed int64
	Req_failed  int64
)

func AddSent() {
	atomic.AddInt64(&Req_sent, 1)
}

func AddSucceed() {
	atomic.AddInt64(&Req_succeed, 1)
}

func AddFailed() {
	atomic.AddInt64(&Req_failed, 1)
}

func Done() {
	log.Println("========================= OVER =========================")
	log.Println("Req_sent::", atomic.LoadInt64(&Req_sent))
	log.Println("Req_succeed::", atomic.LoadInt64(&Req_succeed))
	log.Println("Req_failed::", atomic.LoadInt64(&Req_failed))
}
