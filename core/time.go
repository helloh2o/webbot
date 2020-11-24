package core

import (
	"fmt"
	"time"
)

func GetTime2Time(start int64) string {
	end := time.Now().UnixNano()
	sub := time.Duration(end - start)
	ss := sub / time.Second
	nano := sub % time.Second
	cost := fmt.Sprintf("%d:%d s", ss, nano)
	return cost
}
