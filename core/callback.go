package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"robot/bean"
	"robot/simple"
	"strconv"
	"sync/atomic"
)

// 登录成功
func (r *robot) loginSucceed(data interface{}) {
	dataMap := data.(map[string]interface{})
	r.token, _ = dataMap["token"].(string)
	if r.token != "" {
		if user, ok := dataMap["user"].(map[string]interface{}); ok {
			uid, _ := user["id"].(float64)
			if r.UserId > 0 {
				atomic.AddInt64(&r.UserId, -r.UserId)
			}
			atomic.AddInt64(&r.UserId, int64(uid))
			r.Req.AddHeader("X-User-Token", r.token)
			r.Req.AddHeader("X-User-Id", strconv.Itoa(int(r.UserId)))
			log.Printf("用户 %s, 登录成功!", r.user.Username)

		}
	}
}

// 注册成功
func (r *robot) registerSucceed(data interface{}) {
	dataMap := data.(map[string]interface{})
	r.token, _ = dataMap["token"].(string)
	if r.token != "" {
		// 已登录，可以其他操作
		if user, ok := dataMap["user"].(map[string]interface{}); ok {
			uid, _ := user["id"].(float64)
			username, _ := user["username"].(string)
			if r.UserId > 0 {
				atomic.AddInt64(&r.UserId, -r.UserId)
			}
			atomic.AddInt64(&r.UserId, int64(uid))
			r.Req.AddHeader("X-User-Token", r.token)
			r.Req.AddHeader("X-User-Id", strconv.Itoa(int(r.UserId)))
			r.user.Username = username
			// write local file
			f, err := os.OpenFile("./accounts.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				log.Printf("open accounts file error %s", err)
			}
			bw := bufio.NewWriter(f)
			_, err = bw.WriteString(fmt.Sprintf("%s:%s\n", r.user.Username, r.user.Password))
			if err != nil {
				log.Printf("Write account error %s", err)
			} else if err = bw.Flush(); err == nil {
				log.Printf("Write new account (%s:%s) to txt ", r.user.Username, r.user.Password)
			}
		}
	}
}

// 获取到帖子
func (r *robot) gotTopicNews(action bean.Action, data interface{}) {
	dataMap := data.(map[string]interface{})
	if len(dataMap) == 2 {
		var page bean.Page
		var topics []bean.Topic
		pageValue := dataMap["page"]
		data, err := json.Marshal(pageValue)
		if err != nil {
			log.Printf("can't marshl page value")
		} else if err = json.Unmarshal(data, &page); err != nil {
			log.Printf("can't unmarshl page value to page")
			return
		}
		resultsValue := dataMap["results"]
		if resultsValue != nil {
			data, err = json.Marshal(resultsValue)
			if err != nil {
				log.Printf("can't marshl topics value")
			} else if err = json.Unmarshal(data, &topics); err != nil {
				log.Printf("can't unmarshl topics value to page")
				return
			}
		}
		log.Printf("第 %d 页面， 帖子 %d", page.Page, len(topics))
		// 如果已经没有数据了，回到第一页
		if len(topics) < 10 {
			page.Page = 1
		} else {
			// 随机评论
			r.CommentTopic(topics)
		}
		r.Lock()
		defer r.Unlock()
		r.Views[action.Name] = bean.View{
			Topics: topics,
			Page:   page,
		}

	}
}

func (r *robot) getJsonResult(reader io.Reader, successCall func(data interface{})) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("read json result error %v", err)
		return
	}
	var ret simple.JsonResult
	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Printf("Unmarshal json result error %v", err)
		return
	} else if ret.Success && successCall != nil {
		successCall(ret.Data)
	} else {
		log.Printf("Client %d result %v, error %v", r.Id, ret.Success, ret.Message)
	}
}
