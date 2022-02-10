/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 20:53:43
 * @LastEditTime: 2022-01-07 22:12:35
 */
package myreq

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
	"golang.org/x/text/transform"
)

type Response struct {
	RawRequest  *http.Request
	RawResponse *http.Response

	body []byte
}

// 获取返回的状态码
func (r *Response) StatusCode() int {
	if r.RawResponse == nil {
		return 0
	}
	return r.RawResponse.StatusCode
}

// 获取返回包的 headers
func (r *Response) Header() http.Header {
	if r.RawResponse == nil {
		return http.Header{}
	}
	return r.RawResponse.Header
}

// 返回包的 body，类型为 []byte
func (r *Response) Body() []byte {
	if r.RawResponse == nil {
		return []byte{}
	}
	return r.body
}

// 返回包的 html， 类型为 string
func (r *Response) Html() string {
	if r.body == nil {
		return ""
	}

	return strings.TrimSpace(string(r.body))
}

// 返回包的 html， 类型为经过编码的 string
func (r *Response) String() string {
	if r.body == nil {
		return ""
	}

	e, name := determineEncodeing(r.body) // 获取编码
	if name != "utf-8" {
		bodyReader := bytes.NewReader(r.body)
		utf8Obj := transform.NewReader(bodyReader, e.NewDecoder()) // 转化为 utf8 格式
		body, _ := io.ReadAll(utf8Obj)
		return strings.TrimSpace(string(body))
	}

	return strings.TrimSpace(string(r.body))
}

// 返回包中的 json 数据
func (r *Response) Json() *simplejson.Json {
	if r.body == nil {
		return &simplejson.Json{}
	}

	res, err := simplejson.NewJson(r.body)
	if err != nil {
		fmt.Printf("(req.Response).Json() %v\n", err)
		return &simplejson.Json{}
	}

	return res
}
