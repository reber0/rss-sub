/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 22:02:18
 * @LastEditTime: 2022-01-07 22:44:32
 */
package myreq

import (
	"net/http"
)

func New() *Client {
	client := &http.Client{}
	if client.Transport == nil {
		client.Transport = createTransport(nil)
	}

	return &Client{
		httpClient: client,
		Header:     http.Header{},
	}
}
