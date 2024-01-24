/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 15:14:59
 * @LastEditTime: 2024-01-24 13:47:47
 */
package routers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reber0/goutils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/middleware"
	"github.com/reber0/rss-sub/mydb"
)

// User 相关路由
func UserRouter(r *gin.Engine) {
	userGroup := r.Group("/api/user")

	userGroup.POST("/logout", logout)
	userGroup.POST("/login", login)
	userGroup.POST("/avatar", middleware.JWTAuth(), userAvatar)

	userGroup.Use(middleware.JWTAuth(), middleware.Action(), middleware.RootAuth())
	userGroup.POST("/list", userList)
	userGroup.POST("/add", userAdd)
	userGroup.POST("/update", userUpdate)
	userGroup.POST("/delete", userDelete)
}

func login(c *gin.Context) {
	type User struct {
		// 客户端传入 {"username": "xxx", "password":"123456"}
		UserName string `json:"username"`
		PassWord string `json:"password"`
	}

	json := User{}
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(500, gin.H{"err": err})
	}
	md5_pwd := goutils.Md5([]byte(json.PassWord))

	data := mydb.User{}
	tx := global.Db.Where("uname = ?", json.UserName).First(&data)
	if tx.Error != nil {
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
				"Token": token,
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
		PageIndex    int    `json:"page"`
		PageSize     int    `json:"limit"`
		UserName     string `json:"username"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		ExportIdList []int  `json:"export_id_list"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Uname     string `json:"uname" gorm:"column:uname; type:varchar(50); not null; comment:用户名"`
		Role      string `json:"role" gorm:"column:role; type:varchar(10); not null; comment:用户身份，root/user/e.g."`
		Email     string `json:"email" gorm:"column:email; type:varchar(100); not null; comment:用户邮箱"`
		Avatar    string `json:"avatar" gorm:"column:avatar; type:varchar(40); not null; comment:头像图片名"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; comment:添加时间"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		username := postJson.UserName
		email := postJson.Email
		role := postJson.Role

		var count int64
		var datas []RespData
		tx := global.Db.Model(&mydb.User{})
		if postJson.ExportIdList != nil {
			tx.Where("id in ?", postJson.ExportIdList)
		} else {
			if username != "" {
				tx = tx.Where("uname like ?", "%"+username+"%")
			}
			if email != "" {
				tx = tx.Where("email like ?", "%"+email+"%")
			}
			if role != "" {
				tx = tx.Where("role in ?", []string{role})
			}
		}
		tx = tx.Count(&count).Order("id desc")
		tx = tx.Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)

		for index, data := range datas {
			datas[index].CreatedAt, _ = goutils.Unix2Str(data.CreatedAt)
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
		UName    string `json:"uname"`
		PassWord string `json:"password"`
		Role     string `json:"role"`
		Email    string `json:"email"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		user_id := strings.ReplaceAll(uuid.New().String(), "-", "")
		uname := postJson.UName
		password := goutils.Md5([]byte(postJson.PassWord))
		role := postJson.Role
		email := postJson.Email
		avatar := strconv.Itoa(goutils.RandomInt(1, 9)) + ".png"

		u := mydb.User{UID: user_id, Uname: uname, PassWord: password, Role: role, Email: email, Avatar: avatar}
		tx := global.Db.Create(&u)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
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
		ID       int    `json:"id"`
		Uname    string `json:"uname"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	var postJson PostData
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		uname := postJson.Uname
		password := goutils.Md5([]byte(postJson.Password))

		updateData := make(map[string]interface{})
		if uname != "" {
			updateData["uname"] = uname
		}
		if password != "" {
			updateData["password"] = password
		}
		tx := global.Db.Model(&mydb.User{}).Where("id=?", postJson.ID).Updates(updateData)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
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
		DeleteIDList []int `json:"id_list"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		uid := c.GetString("uid")

		deleteIDList := postJson.DeleteIDList

		if goutils.IsInCol(uid, postJson.DeleteIDList) {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "不能删除自己",
			})
			return
		}

		for _, deleteID := range deleteIDList {
			var deleteUserID string
			tx := global.Db.Model(mydb.User{}).Select("uid").Where("id=?", deleteID).First(&deleteUserID)
			if tx.Error != nil {
				global.Log.Error(tx.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "删除失败",
				})
				return
			} else {
				// 删除用户
				tx := global.Db.Where("uid=?", deleteUserID).Delete(&mydb.User{})
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				// 删除用户 Article 相关数据
				var articleIDList []int // 获取 articleIDList 为了在 data 表中删除数据
				tx = global.Db.Model(&mydb.Article{}).Select("id").Where("uid=?", deleteUserID).Find(&articleIDList)
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
				tx = global.Db.Where("uid=?", deleteUserID).Delete(&mydb.Article{})
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				// 删除用户 Video 相关数据
				var videoIDList []int // 获取 videoIDList 为了在 data 表中删除数据
				tx = global.Db.Model(&mydb.Video{}).Select("id").Where("uid=?", deleteUserID).Find(&videoIDList)
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
				tx = global.Db.Where("uid=?", deleteUserID).Delete(&mydb.Video{})
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}

				// 删除用户 Data 相关数据
				tx = global.Db.Where("category=? and ref_id in ?", "article", articleIDList).Delete(&mydb.Data{})
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
				tx = global.Db.Where("category=? and ref_id in ?", "video", videoIDList).Delete(&mydb.Data{})
				if tx.Error != nil {
					global.Log.Error(tx.Error.Error())
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
