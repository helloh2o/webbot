package core

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"robot/bean"
	"robot/conf"
	"sync/atomic"
	"time"
)

// 发个简单的帖子
// {"node_id":3,"title":"这是机器人发的帖子","tags":["机器人","发帖"],"content":"<p>我是一个机器人，咿呀咿呀呦。</p><p><br></p>","ats":[],"type":0,"money_type":0,"amount":0,"reward_hours":0}
func (r *robot) PostPainTopic(post bean.Action) {
	// json
	params := make(map[string]interface{})
	params["node_id"] = 1 // 机器人模块
	params["title"] = fmt.Sprintf("这是机器人: %s, 发的帖子", r.user.Username)
	params["content"] = fmt.Sprintf("<p>我的名字叫: %s</p><p><br></p>", r.user.Username)
	params["tags"] = []string{"机器人", "发帖"}
	params["ats"] = []string{}
	params["type"] = 0
	params["money_type"] = 0
	params["amount"] = 0
	params["reward_hours"] = 0
	post.Callback = func(reader io.Reader) {
		r.getJsonResult(reader, func(data interface{}) {
			// Todo 发帖成功
			log.Println("发帖成功")
		})
	}
	r.createPostReq(post, params)
}

// 获取最新热帖
func (r *robot) GetNews(act bean.Action) {
	act.Callback = func(reader io.Reader) {
		r.getJsonResult(reader, func(data interface{}) {
			r.gotTopicNews(act, data)
		})
	}
	r.do(&act)
}

// 对某个帖子进行评论 {"Topic_id":7,"reply_id":0,"content":"<p>测试</p>","attachment_ids":[],"ats":[]}
// 随机触发一个
func (r *robot) CommentTopic(topics []bean.Topic) {
	// 时间间隔10秒
	if (time.Now().Unix() - atomic.LoadInt64(&r.lastComment)) <= 10 {
		//log.Printf("评论过快 ... ")
		return
	}
	// 标准洗牌算法
	//n := rand.Perm(len(topics) -1)
	rand.Seed(time.Now().UnixNano())
	topic := topics[rand.Intn(len(topics)-1)]
	params := make(map[string]interface{})
	params["Topic_id"] = topic.Id // 机器人模块
	params["reply_id"] = 0
	params["content"] = fmt.Sprintf("<p>这是，%s 发起的评论</p>", r.user.Username)
	params["attachment_ids"] = []string{}
	params["ats"] = []string{}
	action := bean.Action{
		Url:    conf.Get().TargetURL + "/api/comment/reply",
		Method: "POST",
		Callback: func(reader io.Reader) {
			r.getJsonResult(reader, func(data interface{}) {
				log.Printf("评论主题 %d, 成功！", topic.Id)
			})
		},
	}
	r.createPostReq(action, params)
	old := atomic.LoadInt64(&r.lastComment)
	if old > 0 {
		atomic.AddInt64(&r.lastComment, -old)
	}
	atomic.AddInt64(&r.lastComment, time.Now().Unix())
}
