/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:53:42
 * @LastEditTime: 2022-04-30 09:51:26
 */
package routers

import (
	"RssSub/global"
	"RssSub/middleware"
	"RssSub/mydb"

	"github.com/gin-gonic/gin"
	"github.com/reber0/go-common/utils"
)

// 设置相关路由(网站/邮箱/个人资料/密码)
func SetRouter(r *gin.Engine) {
	setGroup := r.Group("/api/set").Use(middleware.JWTAuth())
	{
		setGroup.POST("/user/info", getInfo)
		setGroup.POST("/user/info_update", middleware.Action(), updateInfo)
		setGroup.POST("/user/pwd_update", middleware.Action(), updatePwd)

		setGroup.POST("/system/website", middleware.RootAuth(), getWebSite)
		setGroup.POST("/system/website_update", middleware.RootAuth(), middleware.Action(), updateWebSite)
		setGroup.POST("/system/email", middleware.RootAuth(), getEmail)
		setGroup.POST("/system/email_update", middleware.RootAuth(), middleware.Action(), updateEmail)
	}
}

func getWebSite(c *gin.Context) {
	var configs []mydb.Config
	result := global.Db.Model(&mydb.Config{}).Find(&configs)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	datas := make(map[string]string, 10)
	for _, config := range configs {
		k, v := config.Key, config.Value
		if k == "sitename" || k == "domain" || k == "upload_max_size" || k == "title" || k == "keyword" || k == "descript" || k == "copyright" {
			datas[k] = v
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": datas,
	})
}

func updateWebSite(c *gin.Context) {
	type PostData struct {
		SiteName      string `form:"sitename" json:"sitename"`
		Domain        string `form:"domain" json:"domain"`
		UploadMaxSize string `form:"upload_max_size" json:"upload_max_size"`
		Title         string `form:"title" json:"title"`
		KeyWord       string `form:"keyword" json:"keyword"`
		Descript      string `form:"descript" json:"descript"`
		Copyright     string `form:"copyright" json:"copyright"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		result := global.Db.Model(&mydb.Config{}).Where("key='sitename'").Update("value", postJson.SiteName)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='domain'").Update("value", postJson.Domain)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='upload_max_size'").Update("value", postJson.UploadMaxSize)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='title'").Update("value", postJson.Title)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='keyword'").Update("value", postJson.KeyWord)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='descript'").Update("value", postJson.Descript)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='copyright'").Update("value", postJson.Copyright)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
		})
	}
}

func getEmail(c *gin.Context) {
	var configs []mydb.Config
	result := global.Db.Model(&mydb.Config{}).Find(&configs)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	datas := make(map[string]string, 10)
	for _, config := range configs {
		k, v := config.Key, config.Value
		if k == "send_email_pwd" {
			datas[k] = "******"
		}
		if k == "send_email" || k == "send_nickname" || k == "smtp_port" || k == "smtp_server" {
			datas[k] = v
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": datas,
	})

}

func updateEmail(c *gin.Context) {
	type PostData struct {
		SendEmail    string `form:"send_email" json:"send_email"`
		SendEmailPwd string `form:"send_email_pwd" json:"send_email_pwd"`
		SendNickName string `form:"send_nickname" json:"send_nickname"`
		SmtpPort     string `form:"smtp_port" json:"smtp_port"`
		SmtpServer   string `form:"smtp_server" json:"smtp_server"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		result := global.Db.Model(&mydb.Config{}).Where("key='send_email'").Update("send_email", postJson.SendEmail)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='send_email_pwd'").Update("send_email_pwd", postJson.SendEmailPwd)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='send_nickname'").Update("send_nickname", postJson.SendNickName)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='smtp_port'").Update("smtp_port", postJson.SmtpPort)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		result = global.Db.Model(&mydb.Config{}).Where("key='smtp_server'").Update("smtp_server", postJson.SmtpServer)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
		})
	}
}

func getInfo(c *gin.Context) {
	type RespData struct {
		Uname     string `json:"uname" gorm:"column:uname; type:varchar(50); not null; comment:用户名"`
		Role      string `json:"role" gorm:"column:role; type:varchar(10); not null; comment:用户身份，root/user/e.g."`
		Email     string `json:"email" gorm:"column:email; type:varchar(100); not null; comment:用户邮箱"`
		Avatar    string `json:"avatar" gorm:"column:avatar; type:varchar(40); not null; comment:头像图片名"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; comment:添加时间"`
	}

	userId := c.GetString("uid")

	var data RespData
	result := global.Db.Model(&mydb.User{}).Where("uid=?", userId).First(&data)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	data.CreatedAt = utils.Unix2String(data.CreatedAt)

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

func updateInfo(c *gin.Context) {
	type PostData struct {
		Uname string `form:"uname" json:"uname"`
		Email string `form:"email" json:"email"`
	}

	postJson := mydb.User{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		uid := c.GetString("uid")

		updateData := mydb.User{Uname: postJson.Uname, Email: postJson.Email}
		result := global.Db.Model(&mydb.User{}).Where("uid=?", uid).Updates(updateData)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
		})
	}
}

func updatePwd(c *gin.Context) {
	type PostData struct {
		OldPwd string `form:"old_pwd" json:"old_pwd"`
		NewPwd string `form:"password" json:"password"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		old_pwd := utils.Md5(postJson.OldPwd)
		new_pwd := utils.Md5(postJson.NewPwd)

		userId := c.GetString("uid")

		var curr_pwd string
		result := global.Db.Model(&mydb.User{}).Select("password").Where("uid=?", userId).First(&curr_pwd)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "修改失败",
			})
			return
		}

		if curr_pwd == old_pwd {
			result = global.Db.Model(&mydb.User{}).Where("uid=?", userId).Update("password", new_pwd)
			if result.Error != nil {
				global.Log.Error(result.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "更新失败",
				})
				return
			}
		} else {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "旧密码错误",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
		})
	}
}
