/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 15:14:59
 * @LastEditTime: 2022-02-14 15:03:23
 */
package routers

import (
	"RssSub/global"
	"RssSub/middleware"
	"RssSub/mydb"
	"RssSub/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User 相关路由
func UserRouter(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/logout", logout)
		userGroup.POST("/login", login)
		userGroup.POST("/avatar", middleware.JWTAuth(), userAvatar)
		userGroup.POST("/list", middleware.JWTAuth(), middleware.Action(), middleware.RootAuth(), userList)
		userGroup.POST("/add", middleware.JWTAuth(), middleware.Action(), middleware.RootAuth(), userAdd)
		userGroup.POST("/update", middleware.JWTAuth(), middleware.Action(), middleware.RootAuth(), userUpdate)
		userGroup.POST("/delete", middleware.JWTAuth(), middleware.Action(), middleware.RootAuth(), userDelete)
	}
}

func login(c *gin.Context) {
	type User struct {
		// 客户端传入 {"username": "xxx", "password":"123456"}
		UserName string `form:"username" json:"username"`
		PassWord string `form:"password" json:"password"`
	}

	json := User{}
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(500, gin.H{"err": err})
	}
	md5_pwd := utils.Md5([]byte(json.PassWord))

	data := mydb.User{}
	result := global.Db.Where("uname = ?", json.UserName).First(&data)
	if result.Error != nil {
		loggerMsg("system", "/api/user/login", json.UserName+"登录失败")
		c.JSON(401, gin.H{ // 未查到用户名，返回用户名或密码错误
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	} else if data.PassWord == md5_pwd {
		token, err := middleware.CreateToken(data.UID)
		if err != nil {
			c.JSON(500, gin.H{ // 服务端生成 token 错误
				"code": 500,
				"msg":  "error",
			})
			return
		}

		loggerMsg("system", "/api/user/login", json.UserName+"登录成功")
		c.Set("user", data)
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"access_token": token,
			},
		})
	} else {
		loggerMsg("system", "/api/user/login", json.UserName+"登录失败")
		c.JSON(401, gin.H{ // 密码错误，返回用户名或密码错误
			"code": 401,
			"msg":  "用户名或密码错误",
		})
	}
}

func logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "退出登录",
	})
}

// 获取头像名字
func userAvatar(c *gin.Context) {
	userId := c.GetString("uid")
	avatar := GetAvatar(userId)

	c.JSON(200, gin.H{
		"code":   0,
		"avatar": avatar,
	})
}

func userList(c *gin.Context) {
	type PostData struct {
		PageIndex    int    `form:"page" json:"page"`
		PageSize     int    `form:"limit" json:"limit"`
		UserName     string `form:"username" json:"username"`
		Email        string `form:"email" json:"email"`
		Role         string `form:"role" json:"role"`
		ExportIdList []int  `form:"export_id_list" json:"export_id_list"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Uname     string `json:"uname" gorm:"column:uname; type:varchar(50); not null; comment:用户名"`
		Role      string `json:"role" gorm:"column:role; type:varchar(10); not null; comment:用户身份，root/user/e.g."`
		Email     string `json:"email" gorm:"column:email; type:varchar(100); not null; comment:用户邮箱"`
		Avatar    string `json:"avatar" gorm:"column:avatar; type:varchar(40); not null; comment:头像图片名"`
		CreatedAt string `json:"add_time" gorm:"column:create_at; comment:添加时间"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		username := fmt.Sprintf("%%%s%%", postJson.UserName)
		email := fmt.Sprintf("%%%s%%", postJson.Email)
		role := fmt.Sprintf("%%%s%%", postJson.Role)

		var count int64
		var datas []RespData
		if postJson.ExportIdList != nil {
			result := global.Db.Model(&mydb.User{}).Where(
				"id in ?", postJson.ExportIdList).Order("id desc").Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
			if result.Error != nil {
				global.Log.Error(result.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "查询失败",
				})
				return
			}
			count = result.RowsAffected
		} else {
			result := global.Db.Model(&mydb.User{}).Where(
				"uname like ? and email like ? and role like ?", username, email, role).Order("id desc").Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
			if result.Error != nil {
				global.Log.Error(result.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "查询失败",
				})
				return
			}
			count = result.RowsAffected
		}

		for index, data := range datas {
			datas[index].CreatedAt = utils.UnixToTime(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":  0,
			"data":  datas,
			"count": count,
		})
	}
}

func userAdd(c *gin.Context) {
	type PostData struct {
		UserName string `form:"username" json:"username"`
		PassWord string `form:"password" json:"password"`
		Role     string `form:"role" json:"role"`
		Email    string `form:"email" json:"email"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		user_id := strings.ReplaceAll(uuid.New().String(), "-", "")
		username := postJson.UserName
		password := utils.Md5([]byte(postJson.PassWord))
		role := postJson.Role
		email := postJson.Email
		avatar := strconv.Itoa(utils.RandInt(1, 9)) + ".png"

		u := mydb.User{UID: user_id, Uname: username, PassWord: password, Role: role, Email: email, Avatar: avatar}
		result := global.Db.Create(&u)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "添加失败",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "添加成功",
			})
		}
	}
}

func userUpdate(c *gin.Context) {
	type PostData struct {
		ID       int    `form:"id" json:"id"`
		Uname    string `form:"uname" json:"uname"`
		Password string `form:"password" json:"password"`
		Email    string `form:"email" json:"email"`
		Role     string `form:"role" json:"role"`
	}

	var postJson PostData
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		postJson.Password = utils.Md5([]byte(postJson.Password))

		updateData := mydb.User{Uname: postJson.Uname, PassWord: postJson.Password, Email: postJson.Email, Role: postJson.Role}
		result := global.Db.Model(&mydb.User{}).Where("id=?", postJson.ID).Updates(updateData)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "更新成功",
			})
		}
	}
}

func userDelete(c *gin.Context) {
	type PostData struct {
		DeleteIDList []int `form:"id_list" json:"id_list"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		uid := c.GetString("uid")

		deleteIDList := postJson.DeleteIDList

		if utils.InSlice(uid, postJson.DeleteIDList) {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "不能删除自己",
			})
			return
		}

		for _, deleteID := range deleteIDList {
			var deleteUserID string
			result := global.Db.Model(mydb.User{}).Select("uid").Where("id=?", deleteID).First(&deleteUserID)
			if result.Error != nil {
				global.Log.Error(result.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "删除失败",
				})
				return
			} else {
				// 删除用户
				result := global.Db.Where("uid=?", deleteUserID).Delete(&mydb.User{})
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				// 删除用户 Article 相关数据
				var articleIDList []int // 获取 articleIDList 为了在 data 表中删除数据
				result = global.Db.Model(&mydb.Article{}).Select("id").Where("uid=?", deleteUserID).Find(&articleIDList)
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
				result = global.Db.Where("uid=?", deleteUserID).Delete(&mydb.Article{})
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				// 删除用户 Video 相关数据
				var videoIDList []int // 获取 videoIDList 为了在 data 表中删除数据
				result = global.Db.Model(&mydb.Video{}).Select("id").Where("uid=?", deleteUserID).Find(&videoIDList)
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
				result = global.Db.Where("uid=?", deleteUserID).Delete(&mydb.Video{})
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				// 删除用户 Data 相关数据
				result = global.Db.Where("category=? and ref_id in ?", "article", articleIDList).Delete(&mydb.Data{})
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
				result = global.Db.Where("category=? and ref_id in ?", "video", videoIDList).Delete(&mydb.Data{})
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				c.JSON(200, gin.H{
					"code": 0,
					"msg":  "删除成功",
				})
			}
		}
	}
}
