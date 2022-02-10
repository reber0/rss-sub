/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 10:01:01
 * @LastEditTime: 2022-02-10 00:37:17
 */
package middleware

import (
	"RssSub/global"
	"RssSub/mydb"
	"bytes"
	"io"
	"regexp"

	"github.com/gin-gonic/gin"
)

// 定义 Action 中间件，记录用户操作
func Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		uname := c.GetString("uname")
		action := c.Request.URL.Path

		body, _ := io.ReadAll(c.Request.Body)
		postData := string(body)

		reg := regexp.MustCompile(`password":".*?"`)
		postData = reg.ReplaceAllString(postData, `password":"******"`)
		reg = regexp.MustCompile(`old_pwd":".*?"`)
		postData = reg.ReplaceAllString(postData, `old_pwd":"******"`)
		data := mydb.Message{Uname: uname, Action: action, Data: postData, Status: "unread"}
		global.Db.Create(&data)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}
