package core

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Req struct {
	sync.RWMutex
	http.Client
	proxy  string
	header map[string]string
}

func (req *Req) AddHeader(key, val string) {
	req.Lock()
	defer req.Unlock()
	req.header[key] = val
}

func (req *Req) GetHeader() map[string]string {
	req.RLock()
	defer req.RUnlock()
	ret := make(map[string]string)
	for k, v := range req.header {
		ret[k] = v
	}
	return ret
}

var (
	agent       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.3578.98 Safari/537.36"
	agentMobile = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Mobile Safari/537.36"
)

// any args call mobile
func NewReq(proxyUrl string, args ...interface{}) *Req {
	req := new(Req)
	req.header = make(map[string]string)
	req.Client = http.Client{}
	req.Jar = new(Jar)
	req.Timeout = time.Second * 10
	tr := &http.Transport{}
	if len(args) > 0 {
		agent = agentMobile
	}
	if proxyUrl != "" {
		proxyFunc := func(r *http.Request) (*url.URL, error) {
			r.Header.Set("User-Agent", agent)
			return url.Parse(proxyUrl)
		}
		tr.Proxy = proxyFunc
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req.Transport = tr
	return req
}

func (req *Req) Send(targetUrl string, method string, body []byte, callback func(*url.URL, io.Reader)) bool {
	switch method {
	case "GET":
	case "POST":
	default:
		log.Printf("Unsupport method %s", method)
		return false
	}
	urlInfo, err := url.Parse(targetUrl)
	if err != nil {
		log.Printf("can't parse url error %v", err)
		return false
	}
	var bd io.Reader
	if body != nil {
		bd = bytes.NewBuffer(body)
	}
	reqInfo, err := http.NewRequest(method, targetUrl, bd)
	if err != nil {
		log.Printf("new http request for url %s error", targetUrl)
		return false
	}
	if req.header != nil {
		for k, v := range req.GetHeader() {
			reqInfo.Header.Add(k, v)
		}
	}
	reqInfo.Header.Add("Host", urlInfo.Host)
	reqInfo.Header.Add("User-Agent", agent)
	resp, err := req.Do(reqInfo)
	if err != nil {
		log.Printf("do request error %v", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		if callback != nil {
			callback(urlInfo, resp.Body)
		}
		return true
	}
	return false
}

type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}
