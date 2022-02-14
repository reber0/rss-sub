/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 10:01:01
 * @LastEditTime: 2022-02-14 09:53:35
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
		uid := c.GetString("uid")

		var uname string
		result := global.Db.Model(&mydb.User{}).Select("uname").Where("uid = ?", uid).First(&uname)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
		}

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
