package core

import (
	"io"
	"robot/bean"
	"robot/consts"
)

func (r *robot) getCaptcha(act bean.Action, next bean.Action) {
	act.Next = &next
	act.Next.Pre = &act
	// 动作回调
	act.Callback = func(reader io.Reader) {
		r.getJsonResult(reader, func(data interface{}) {
			dataMap := data.(map[string]interface{})
			r.captchaId, _ = dataMap["captchaId"].(string)
			r.captchaCode, _ = dataMap["captchaCode"].(string)
			// 是否获取到了验证码
			if r.captchaId != "" && r.captchaCode != "" {
				switch act.Next.Name {
				// 注册账号
				case consts.Signup:
					r.register(act)
					// 登录
				case consts.SignIn:
					// 用户名密码不齐全的情况，去注册一个
					if r.user.Username != "" && r.user.Password != "" {
						r.Login(act)
					} else {
						r.register(act)
					}
				}
			}
		})
	}
	// 执行获取验证码
	r.do(&act)
}

// 注册
func (r *robot) register(pre bean.Action) {
	reg := r.actionsMap[consts.Signup]
	reg.Pre = &pre
	// 注册后获取验证码、
	temp := r.actionsMap[consts.Captcha]
	reg.Next = &temp
	reg.Callback = func(reader io.Reader) {
		// 获取token
		r.getJsonResult(reader, r.registerSucceed)
	}
	// 用户名，昵称
	//{"username":"fa'd'fa'fdafadsf","nickname":"dfadfdasf","password":"123456","rpassword":"123456","ref":"/","captchaCode":"6048","captchaId":"pRs2e3EvckyunjMEupnR","email":"123456@qq.com"}
	rdstr := GetRandomString(8)
	email := rdstr + "@gmail.com"
	params := make(map[string]interface{})
	params["username"] = rdstr
	params["nickname"] = rdstr
	params["password"] = rdstr
	params["rpassword"] = rdstr
	params["email"] = email
	params["ref"] = "/"
	params["captchaCode"] = r.captchaCode
	params["captchaId"] = r.captchaId
	r.user.Password = rdstr
	r.createPostReq(reg, params)
}

// 登录 {"username":"aiyongay","password":"aiyongay","captchaCode":"3106","captchaId":"eNjafH1Dj3g7chnARKgt","ref":"/"}
func (r *robot) Login(pre bean.Action) {
	signIn := r.actionsMap[consts.SignIn]
	signIn.Pre = &pre
	params := make(map[string]interface{})
	params["username"] = r.user.Username
	params["password"] = r.user.Password
	params["captchaCode"] = r.captchaCode
	params["captchaId"] = r.captchaId
	params["ref"] = "/"
	signIn.Callback = func(reader io.Reader) {
		r.getJsonResult(reader, r.loginSucceed)
	}
	r.createPostReq(signIn, params)
}

func (r *robot) do(action *bean.Action) {
	go func() {
		r.ReqChan <- action
	}()
}
