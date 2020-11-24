package bean

import "io"

type Action struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`    // 路径
	Tick     int               `json:"tick"`   // 频率（mill_second）
	Times    int               `json:"times"`  // 总次数
	Method   string            `json:"method"` // 请求方法
	Page     int               `json:"page"`
	PageName string            `json:"page_name"`
	Params   map[string]string `json:"-"`
	Body     []byte            `json:"-"`
	Pre      *Action           `json:"-"` // 上一个动作
	Current  *Action           `json:"-"` // 当前动作
	Next     *Action           `json:"-"` // 下一个动作
	Callback func(reader io.Reader)
}
