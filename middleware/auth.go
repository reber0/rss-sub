/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-07 10:01:01
 * @LastEditTime: 2022-02-10 16:03:24
 */
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义 JWTAuth 中间件，进行登录校验
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判定是否携带 token
		token := c.Request.Header.Get("access_token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 1001, // code 为 1001 时返回登录页
				"msg":  "请求未携带 token，无权限访问",
			})
			c.Abort()
			return
		}

		claims, err := ParserToken(token)
		if err != nil {
			c.JSON(401, gin.H{
				"code": 1001, // code 为 1001 时返回登录页
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		// 将解析后的有效载荷写入 gin.Context 引用对象中
		c.Set("uid", claims.UID)
		c.Set("uname", claims.UName)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Set("avatar", claims.Avatar)
	}
}

// 定义 RootAuth 中间件，校验用户是否为 root 权限
func RootAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
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
