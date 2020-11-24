package simple

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"log"
	"reflect"
	"strings"
)

func FormatJson(obj interface{}) (str string, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}

func ParseJson(str string, t interface{}) error {
	return json.Unmarshal([]byte(str), t)
}

/**
严格按Json字段反序列化
*/
func JsonUnmarshal(str string, data interface{}) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		return errors.New("参数结构必须是指针")
	}

	dec := json.NewDecoder(strings.NewReader(str))
	dec.DisallowUnknownFields()
	if err := dec.Decode(data); err != nil {
		msg := fmt.Sprintf("type:[%s] err:[%+v] json:%+v ", strings.ToLower(t.String()), err, str)
		return errors.New(msg)
	}

	return nil
}

func ReadBody(ctx iris.Context, outptr interface{}) error {
	data, err := ctx.GetBody()
	if err != nil {
		log.Println(err)
		return errors.New("读取body内容错误")
	}

	str := string(data)
	if err := JsonUnmarshal(str, outptr); err != nil {
		log.Println("参数不合法")
		log.Println(err)
		return err
	}

	return nil
}
