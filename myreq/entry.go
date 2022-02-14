/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 22:02:18
 * @LastEditTime: 2022-02-14 16:55:13
 */
package myreq

import (
	"net/http"
	"net/url"
)

func New() *Client {
	client := &http.Client{}
	if client.Transport == nil {
		client.Transport = createTransport(nil)
	}

	return &Client{
		httpClient: client,
		Header:     http.Header{},
		FormData:   url.Values{},
	}
}
