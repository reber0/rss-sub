/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 20:52:21
 * @LastEditTime: 2022-02-14 17:45:28
 */
package myreq

import (
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

// 获取编码格式
func determineEncodeing(data []byte) (encoding.Encoding, string) {
	// 获取数据,Peek返回输入流的下n个字节
	var bytes []byte
	if len(data) > 1024 {
		bytes = data[:1024]
	} else {
		bytes = data
	}

	// 调用DEtermineEncoding函数，确定编码通过检查最多前 1024 个字节的内容和声明的内容类型来确定 HTML 文档的编码。
	// e, name, certain := charset.DetermineEncoding(bytes, "")
	// fmt.Println("DetermineEncoding: ", e, name, certain)
	e, name, _ := charset.DetermineEncoding(bytes, "")
	return e, name
}
