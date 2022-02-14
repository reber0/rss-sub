/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2021-12-12 19:32:15
 * @LastEditTime: 2022-02-14 17:39:41
 */

package myreq

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	Header     http.Header
	FormData   url.Values
	Body       string // 发送原始 body
}

//设置 Transport
func (r *Client) SetTransport(transport http.RoundTripper) *Client {
	if transport != nil {
		r.httpClient.Transport = transport
	}
	return r
}

// 请求 https 网站是否跳过证书验证
func (r *Client) SkipVerify(IsSkipVerify bool) *Client {
	transport := r.httpClient.Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: IsSkipVerify}
	return r
}

// 设置代理
func (r *Client) SetProxy(proxyURL string) *Client {
	transport := r.httpClient.Transport.(*http.Transport)

	pURL, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("%v", err)
		return r
	}

	// 设置 transport 的 Proxy
	transport.Proxy = http.ProxyURL(pURL)

	return r
}

// 设置超时
func (r *Client) SetTimeout(timeout time.Duration) *Client {
	r.httpClient.Timeout = timeout
	return r
}

// 设置单个 header
func (r *Client) SetHeader(header, value string) *Client {
	r.Header.Set(header, value)
	return r
}

// 设置多个 header
func (r *Client) SetHeaders(headers map[string]string) *Client {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

// 设置 post 的数据
func (r *Client) SetFormData(data map[string]string) *Client {
	for k, v := range data {
		r.FormData.Set(k, v)
	}
	r.Body = r.FormData.Encode()
	return r
}

// 设置发送的原始 body
func (r *Client) SetBody(body string) *Client {
	r.Body = body
	return r
}

func (r *Client) Get(url string) (*Response, error) {
	return r.Execute("GET", url)
}

func (r *Client) Post(url string) (*Response, error) {
	return r.Execute("POST", url)
}

func (r *Client) PostJson(url string) (*Response, error) {
	return r.Execute("PostJson", url)
}

func (r *Client) Execute(method, url string) (*Response, error) {
	var request *http.Request
	var response *Response

	if method == "GET" {
		request, _ = http.NewRequest("GET", url, nil)
		request.Header = r.Header
	}
	if method == "POST" {
		request, _ = http.NewRequest("POST", url, strings.NewReader(r.Body))
		request.Header = r.Header
	}
	if method == "PostJson" {
		request, _ = http.NewRequest("POST", url, strings.NewReader(r.Body))
		r.SetHeader("Content-Type", "application/json")
		request.Header = r.Header
	}

	resp, err := r.httpClient.Do(request)
	if resp != nil {
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		response = &Response{
			RawRequest:  request,
			RawResponse: resp,
			body:        body,
		}
		return response, err
	}
	if err != nil {
		fmt.Println(err)
	}

	return response, err
}
