/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:53:07
 * @LastEditTime: 2022-09-20 10:17:06
 */
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/reber0/go-common/utils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/middleware"
	"github.com/reber0/rss-sub/mydb"
	"gorm.io/gorm"
)

// Message 相关路由
func MessageRouter(r *gin.Engine) {
	msgGroup := r.Group("/api/message").Use(middleware.JWTAuth())
	{
		msgGroup.POST("/tabs", msgTabs)
		msgGroup.POST("/count", msgCount)
		msgGroup.POST("/user_list", msgUserList)
		msgGroup.POST("/api_list", msgApiList)
		msgGroup.POST("/update", msgUpdate)
		msgGroup.POST("/read_all", msgReadAll)
		msgGroup.POST("/delete", msgDelete)
		msgGroup.POST("/delete_all", msgDeleteAll)
	}
}

// 显示消息的 tab
func msgTabs(c *gin.Context) {
	userId := c.GetString("uid")
	_, role := GetUserMsg(userId)

	if role == "root" {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"tabs": []string{"api", "user"},
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"tabs": []string{"user"},
			},
		})
	}
}

// 获取未读消息条数
func msgCount(c *gin.Context) {
	userId := c.GetString("uid")
	uname, role := GetUserMsg(userId)

	var count int64
	if role == "root" {
		global.Db.Model(&mydb.Message{}).Where("status='unread'").Count(&count)
	} else if role == "user" {
		global.Db.Model(&mydb.Message{}).Where("status='unread' and uname=? and action like '%更新%'", uname).Count(&count)
	}

	c.JSON(200, gin.H{
		"code":         0,
		"unread_count": count,
	})
}

// 列用户未读消息，主要是更新的消息
func msgUserList(c *gin.Context) {
	type PostData struct {
		PageIndex int `form:"page" json:"page"`
		PageSize  int `form:"limit" json:"limit"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Uname     string `json:"username" gorm:"column:uname; type:varchar(50); comment:msg_type 为 user 时是用户名，为 system 时是 schedule/error"`
		Action    string `json:"action" gorm:"column:action; type:text; comment:执行的操作触发的 URI、计划任务动作"`
		Data      string `json:"data" gorm:"column:data; type:text; comment:POST 的数据、得到的数据"`
		Status    string `json:"status" gorm:"column:status; type:varchar(10); comment:状态，是否已读、已看"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; comment:添加时间"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	} else {
		userId := c.GetString("uid")
		uname, role := GetUserMsg(userId)

		var count int64
		var datas []RespData
		result := global.Db.Model(&mydb.Message{}).Where(
			"(uname=? or ?='root') and action like '%更新%'", uname, role).Count(&count).Order("id desc").Count(&count).Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

		var unread_count int64
		result = global.Db.Model(&mydb.Message{}).Where(
			"(uname=? or ?='root') and action like '%更新%' and status='unread'", uname, role).Count(&unread_count)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

		for index, data := range datas {
			datas[index].CreatedAt = utils.Unix2Str(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":         0,
			"data":         datas,
			"count":        count,
			"unread_count": unread_count,
		})
	}
}

// 列系统相关未读消息，主要是 api 访问记录消息
func msgApiList(c *gin.Context) {
	type PostData struct {
		PageIndex int `form:"page" json:"page"`
		PageSize  int `form:"limit" json:"limit"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Uname     string `json:"username" gorm:"column:uname; type:varchar(50); comment:msg_type 为 user 时是用户名，为 system 时是 schedule/error"`
		Action    string `json:"action" gorm:"column:action; type:text; comment:执行的操作触发的 URI、计划任务动作"`
		Data      string `json:"data" gorm:"column:data; type:text; comment:POST 的数据、得到的数据"`
		Status    string `json:"status" gorm:"column:status; type:varchar(10); comment:状态，是否已读、已看"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; comment:添加时间"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	} else {
		userId := c.GetString("uid")
		uname, role := GetUserMsg(userId)

		var count int64
		var datas []RespData
		result := global.Db.Model(&mydb.Message{}).Where(
			"(uname=? or ?='root') and action like '/api/%'", uname, role).Count(&count).Order("id desc").Count(&count).Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

		var unread_count int64
		result = global.Db.Model(&mydb.Message{}).Where(
			"(uname=? or ?='root') and action like '/api/%' and status='unread'", uname, role).Count(&unread_count)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

		for index, data := range datas {
			datas[index].CreatedAt = utils.Unix2Str(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":         0,
			"data":         datas,
			"count":        count,
			"unread_count": unread_count,
		})
	}
}

// 更新一条或几条消息状态
func msgUpdate(c *gin.Context) {
	type PostData struct {
		UpdateIDList []int  `form:"id_list" json:"id_list"`
		Status       string `form:"status" json:"status"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		userId := c.GetString("uid")
		uname, role := GetUserMsg(userId)

		updateIDList := postJson.UpdateIDList
		status := postJson.Status

		result := global.Db.Model(&mydb.Message{}).Where(
			"(uname=? or ?='root') and id in ?", uname, role, updateIDList).Update("status", status)
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

// 更新所有消息状态为已读
func msgReadAll(c *gin.Context) {
	type PostData struct {
		MsgType string `form:"msgtype" json:"msgtype"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		userId := c.GetString("uid")
		uname, role := GetUserMsg(userId)

		msgtype := postJson.MsgType

		var result *gorm.DB
		if msgtype == "api" {
			result = global.Db.Model(&mydb.Message{}).Where(
				"(uname=? or ?='root') and action like '/api/%'", uname, role).Update("status", "read")
		} else if msgtype == "user" {
			result = global.Db.Model(&mydb.Message{}).Where(
				"(uname=? or ?='root') and action like '%更新%'", uname, role).Update("status", "read")
		}
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

// 删一条或几条消息
func msgDelete(c *gin.Context) {
	type PostData struct {
		DeleteAll    bool  `form:"delete_all" json:"delete_all"`
		DeleteIDList []int `form:"id_list" json:"id_list"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "删除失败",
		})
	} else {
		userId := c.GetString("uid")
		uname, role := GetUserMsg(userId)

		deleteIDList := postJson.DeleteIDList

		result := global.Db.Where(
			"(uname=? or ?='root') and id in ?", uname, role, deleteIDList).Delete(&mydb.Message{})
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

// 删除所有消息
func msgDeleteAll(c *gin.Context) {
	type PostData struct {
		MsgType string `form:"msgtype" json:"msgtype"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		userId := c.GetString("uid")
		uname, role := GetUserMsg(userId)

		msgtype := postJson.MsgType

		var result *gorm.DB
		if msgtype == "api" {
			result = global.Db.Where(
				"(uname=? or ?='root') and action like '/api/%'", uname, role).Delete(&mydb.Message{})
		} else if msgtype == "user" {
			result = global.Db.Where(
				"(uname=? or ?='root') and action like '%更新%'", uname, role).Delete(&mydb.Message{})
		}
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
