/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 10:01:01
 * @LastEditTime: 2022-02-14 10:00:40
 */
package middleware

import (
	"RssSub/global"
	"RssSub/mydb"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义 JWTAuth 中间件，进行登录校验
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判定是否携带 token
		token := c.GetHeader("access_token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "请求未携带 Token",
			})
			c.Abort()
			return
		}

		claims, err := ParserToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		// 将解析后的有效载荷写入 gin.Context 引用对象中
		c.Set("uid", claims.UID)
	}
}

// 定义 RootAuth 中间件，校验用户是否为 root 权限
func RootAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetString("uid")

		var role string
		result := global.Db.Model(&mydb.User{}).Select("role").Where("uid = ?", uid).First(&role)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
		}

		if role != "root" {
			c.JSON(403, gin.H{
				"code": 403,
				"msg":  "无权限访问",
			})
			c.Abort()
			return
		}
	}
}
