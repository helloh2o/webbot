package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"robot/bean"
	"robot/consts"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type robot struct {
	Id      int
	UserId  int64
	Req     *Req
	ReqChan chan *bean.Action
	actions []*bean.Action
	wg      *sync.WaitGroup
	// 验证码
	captchaId   string
	captchaCode string
	// token
	token string
	// 动作通道
	actionChannels []chan *bean.Action
	// 动作映射
	actionsMap map[string]bean.Action
	// 用户数据
	user *User
	// 浏览页面
	sync.RWMutex
	Views map[string]bean.View
	// 上次回复时间
	lastComment int64
}

type User struct {
	Username     string `json:"username"`
	Password     string `json:"-"`
	Gold         string `json:"gold"`
	Score        string `json:"score"`
	Diamond      string `json:"diamond"`
	TopicCount   string `json:"topicCount"`
	CommentCount string `json:"commentCount"`
}

func NewRobot(id int, in chan *bean.Action, wg *sync.WaitGroup) *robot {
	r := &robot{Req: NewReq("")}
	r.ReqChan = in
	r.Id = id
	r.wg = wg
	r.Views = make(map[string]bean.View)
	// 现有账号
	if len(accounts) >= id {
		r.user = &accounts[id-1]
	} else {
		r.user = new(User)
	}
	return r
}

// 设置动作
func (r *robot) SetActions(list []*bean.Action) {
	r.actions = list
	r.actionChannels = make([]chan *bean.Action, 0)
	r.actionsMap = make(map[string]bean.Action)
	for _, a := range list {
		a.Current = a
		r.actionChannels = append(r.actionChannels, make(chan *bean.Action))
		r.actionsMap[a.Name] = *a
	}
}

func (r *robot) SetTimeoutReq(seconds int) {
	r.Req.Timeout = time.Second * time.Duration(seconds)
}

func (r *robot) Run() {
	// 执行
	go func() {
		for {
			action := <-r.ReqChan
			if action.Page != 0 {
				r.countPage(action)
			}
			AddSent()
			start := time.Now().UnixNano()
			if ok := r.Req.Send(action.Url, action.Method, action.Body, func(url *url.URL, reader io.Reader) {
				AddSucceed()
				log.Printf("clinet %d req %s ok, cost %s ", r.Id, action.Url, GetTime2Time(start))
				// 动作成功回调
				if action.Callback != nil {
					action.Callback(reader)
				}
			}); !ok {
				AddFailed()
				log.Printf("==>clinet %d req %s failed. cost time %s ", r.Id, action.Url, GetTime2Time(start))
			}
			if r.wg != nil && action.Name != consts.Captcha && action.Name != "" {
				r.wg.Done()
			}
			// 稍微正常点
			time.Sleep(time.Millisecond * 200)
		}
	}()
	// 动作
	go func() {
		for _, act := range r.actions {
			if act.Times == 0 {
				continue
			} else {
				go r.runTick(*act)
			}
		}
	}()
}

// copy action
func (r *robot) runTick(action bean.Action) {
	switch action.Name {
	case consts.Captcha, consts.SignIn, consts.Signup:
		if action.Tick == 0 {
			action.Tick = 1
		}
		goto tickStart
	default:
		// wait until robot login
		for {
			if atomic.LoadInt64(&r.UserId) != 0 {
				goto tickStart
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
tickStart:
	tk := time.NewTicker(time.Millisecond * time.Duration(action.Tick))
	for {
		<-tk.C
		if action.Times == 0 {
			break
		}
		r.rSelect(action)
		action.Times--
		//log.Printf("Client %d, Action ==>%s tick ==>%d", r.Id, action.Name, action.Times)
	}
}

// 选择动作，加载逻辑和数据
func (r *robot) rSelect(act bean.Action) {
	getCaptcha := r.actionsMap[consts.Captcha]
	switch act.Name {
	// 注册/登录 先获取验证码
	case consts.Signup, consts.SignIn:
		// cp action
		r.getCaptcha(getCaptcha, act)
	case consts.PostTopic:
		r.PostPainTopic(act)
	case consts.GetNews:
		r.GetNews(act)
	default:
		r.ReqChan <- &act
	}
}

// 创建post请求
func (r *robot) createPostReq(action bean.Action, params map[string]interface{}) {
	var err error
	action.Body, err = json.Marshal(&params)
	if err != nil {
		log.Printf("%s error %s", action.Name, err.Error())
	} else {
		r.do(&action)
	}
}

// 计算页面
func (r *robot) countPage(action *bean.Action) {
	r.RLock()
	defer r.RUnlock()
	if action.Page == 0 {
		return
	}
	if action.PageName == "" {
		action.PageName = "page"
	}
	// 是否有浏览记录
	view, ok := r.Views[action.Name]
	if ok {
		urlInfo, err := url.Parse(action.Url)
		if err != nil {
			log.Printf("parse url %s error %v", action.Url, err)
		} else {
			// 下一页
			next := strconv.Itoa(view.Page.Page + 1)
			if len(urlInfo.Query()) == 0 {
				action.Url = action.Url + fmt.Sprintf("?%s=%s", action.PageName, next)
			} else {
				action.Url = action.Url + fmt.Sprintf("&%s=%s", action.PageName, next)
			}
		}
	}
}
