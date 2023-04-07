/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-06 14:59:35
 * @LastEditTime: 2023-04-07 11:10:53
 */
package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/reber0/go-common/utils"
)

var jwtSecret = []byte(utils.RandomString(8))

// var jwtSecret = []byte("12345678")

// 自定义有效载荷(这里采用自定义的 UserID 作为有效载荷的一部分)
type CustomClaims struct {
	UID                string `json:"uid"`
	jwt.StandardClaims        // StandardClaims 结构体实现了 Claims 接口
}

// token 生成
func CreateToken(uid string) (string, error) {
	// expireTime := time.Now().Add(7 * 24 * time.Hour).Unix() // 设置 token 有效时间为 7 天
	expireTime := time.Now().AddDate(0, 1, 0).Unix() // 设置 token 有效时间为 1 个月后

	claims := CustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,                  // token 过期时间
			Issuer:    "github.com/reber0/rss-sub", // token 发行人
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// token 解码
func ParserToken(tokenString string) (*CustomClaims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回 *Token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("Token 不可用")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("Token 过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("Token 无效")
			} else {
				return nil, fmt.Errorf("Token 不可用")
			}
		}
	}

	// 将 token 中的 claims 信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Token 无效")
}
