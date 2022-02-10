/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2021-12-12 19:32:15
 * @LastEditTime: 2022-01-07 22:55:55
 */

package myreq

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	httpClient *http.Client
	Header     http.Header
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

func (r *Client) Get(url string) (*Response, error) {
	return r.Execute("GET", url)
}

func (r *Client) POST(url string) (*Response, error) {
	return r.Execute("POST", url)
}

func (r *Client) Execute(method, url string) (*Response, error) {
	var response *Response

	request, _ := http.NewRequest(method, url, nil)
	request.Header = r.Header

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
