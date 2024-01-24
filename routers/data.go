/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:53:25
 * @LastEditTime: 2024-01-24 13:44:11
 */
package routers

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	"github.com/reber0/goutils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/middleware"
	"github.com/reber0/rss-sub/mydb"
)

// Data 相关路由
func DataRouter(r *gin.Engine) {
	dataGroup := r.Group("/api/data")

	dataGroup.GET("/rss/:uid/:category/:ref_id", getRss) // getRss 路由在中间件之前，不经过中间件处理

	dataGroup.Use(middleware.JWTAuth(), middleware.Action()) // 在路由组中应用中间件
	dataGroup.POST("/article/list", dataArticleList)
	dataGroup.POST("/article/update", dataArticleUpdate)
	dataGroup.POST("/article/delete", dataArticleDelete)
	dataGroup.POST("/video/list", dataVideoList)
	dataGroup.POST("/video/update", dataVideoUpdate)
	dataGroup.POST("/video/delete", dataVideoDelete)
}

func dataArticleList(c *gin.Context) {
	type PostData struct {
		PageIndex int    `json:"page"`
		PageSize  int    `json:"limit"`
		KeyWord   string `json:"keyword"`
		Title     string `json:"title"`
		Status    string `json:"status"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Name      string `json:"name" gorm:"column:name; type:text; comment:博客名字"`
		Title     string `json:"title" gorm:"column:title; type:text; comment:标题，文章名字、番剧每集名字"`
		URL       string `json:"url" gorm:"column:url; type:text; comment:网址，文章链接、番剧每集 URL"`
		Status    string `json:"status" gorm:"column:status; type:varchar(10); comment:状态，是否已读、已看"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; comment:添加时间"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	} else {
		userId := c.GetString("uid")
		_, role := GetUserMsg(userId)

		keyword := postJson.KeyWord
		title := postJson.Title
		status := postJson.Status

		var count int64
		var datas []RespData
		tx := global.Db.Model(&mydb.Data{}).Joins("JOIN article ON data.ref_id = article.id")
		tx = tx.Select("data.id", "article.name", "data.title", "data.url", "data.status", "data.created_at")
		tx = tx.Where("article.uid=? or ?='root'", userId, role)
		tx = tx.Where("data.category='article'")
		if keyword != "" {
			tx = tx.Where("article.name like ?", "%"+keyword+"%")
		}
		if title != "" {
			tx = tx.Where("data.title like ?", "%"+title+"%")
		}
		if status != "" {
			tx = tx.Where("data.status in ?", []string{status})
		}
		tx = tx.Count(&count).Order("data.id desc")
		tx = tx.Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

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

func dataArticleUpdate(c *gin.Context) {
	type PostData struct {
		UpdateIDList []int  `json:"id_list"`
		Status       string `json:"status"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		userId := c.GetString("uid")
		_, role := GetUserMsg(userId)

		updateIDList := postJson.UpdateIDList
		status := postJson.Status

		var count int64
		tx := global.Db.Model(&mydb.Data{}).Joins("JOIN article ON data.ref_id = article.id").Where(
			"(article.uid=? or ?='root') and data.id in ?", userId, role, updateIDList).Count(&count)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		if count == int64(len(updateIDList)) {
			tx := global.Db.Model(&mydb.Data{}).Where("id in ?", updateIDList).Update("status", status)
			if tx.Error != nil {
				global.Log.Error(tx.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "更新失败",
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
		})
	}
}

func dataArticleDelete(c *gin.Context) {
	type PostData struct {
		ID int `json:"id"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "删除失败",
		})
	} else {
		userId := c.GetString("uid")
		_, role := GetUserMsg(userId)

		id := postJson.ID

		var count int64
		tx := global.Db.Model(&mydb.Data{}).Joins("JOIN article ON data.ref_id = article.id").Where(
			"(article.uid=? or ?='root') and data.category='article' and data.id=?", userId, role, id).Count(&count)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "删除失败",
			})
			return
		}
		if count == 1 {
			tx := global.Db.Where("category='article' and id=?", id).Delete(&mydb.Data{})
			if tx.Error != nil {
				global.Log.Error(tx.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "删除失败",
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "删除成功",
		})
	}
}

func dataVideoList(c *gin.Context) {
	type PostData struct {
		PageIndex int    `json:"page"`
		PageSize  int    `json:"limit"`
		KeyWord   string `json:"keyword"`
		Status    string `json:"status"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Name      string `json:"name" gorm:"column:name; type:text; comment:番剧名字"`
		Title     string `json:"title" gorm:"column:title; type:text; comment:标题，文章名字、番剧每集名字"`
		URL       string `json:"url" gorm:"column:url; type:text; comment:网址，文章链接、番剧每集 URL"`
		Status    string `json:"status" gorm:"column:status; type:varchar(10); comment:状态，是否已读、已看"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; comment:添加时间"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "查询失败",
		})
	} else {
		userId := c.GetString("uid")
		_, role := GetUserMsg(userId)

		status := postJson.Status
		keyword := postJson.KeyWord

		var count int64
		var datas []RespData
		tx := global.Db.Model(&mydb.Data{}).Joins("JOIN video ON data.ref_id = video.id")
		tx = tx.Select("data.id", "video.name", "data.title", "data.url", "data.status", "data.created_at")
		tx = tx.Where("video.uid=? or ?='root'", userId, role)
		tx = tx.Where("data.category='video'")
		if keyword != "" {
			tx = tx.Where("video.name like ?", "%"+keyword+"%")
		}
		if status != "" {
			tx = tx.Where("data.status in ?", []string{status})
		}
		tx = tx.Count(&count).Order("data.id asc")
		tx = tx.Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

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

func dataVideoUpdate(c *gin.Context) {
	type PostData struct {
		UpdateIDList []int  `json:"id_list"`
		Status       string `json:"status"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "更新失败",
		})
	} else {
		userId := c.GetString("uid")
		_, role := GetUserMsg(userId)

		updateIDList := postJson.UpdateIDList
		status := postJson.Status

		var count int64
		tx := global.Db.Model(&mydb.Data{}).Joins("JOIN video ON data.ref_id = video.id").Where(
			"(video.uid=? or ?='root') and data.id in ?", userId, role, updateIDList).Count(&count)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			return
		}
		if count == int64(len(updateIDList)) {
			tx := global.Db.Model(&mydb.Data{}).Where("id in ?", updateIDList).Update("status", status)
			if tx.Error != nil {
				global.Log.Error(tx.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "更新失败",
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
		})
	}
}

func dataVideoDelete(c *gin.Context) {
	type PostData struct {
		ID int `json:"id"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "删除失败",
		})
	} else {
		userId := c.GetString("uid")
		_, role := GetUserMsg(userId)

		id := postJson.ID

		var count int64
		tx := global.Db.Model(&mydb.Data{}).Joins("JOIN video ON data.ref_id = video.id").Where(
			"(video.uid=? or ?='root') and data.category='video' and data.id=?", userId, role, id).Count(&count)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "删除失败",
			})
			return
		}
		if count == 1 {
			tx := global.Db.Where("category='video' and id=?", id).Delete(&mydb.Data{})
			if tx.Error != nil {
				global.Log.Error(tx.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "删除失败",
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "删除成功",
		})
	}
}

func getRss(c *gin.Context) {
	uid := c.Param("uid")
	category := c.Param("category")
	ref_id := c.Param("ref_id")

	type SiteMsg struct {
		Name string `json:"name" gorm:"column:name; type:varchar(100); comment:站点名字"`
		Link string `json:"link" gorm:"column:link; type:varchar(100); comment:站点链接"`
	}

	type Item struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Link        string   `xml:"link"`
		PubDate     string   `xml:"pubDate"`
		Description string   `xml:"description"`
	}

	type Channel struct {
		XMLName xml.Name `xml:"channel"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
		Items   []Item   `xml:"item"`
	}

	type Rss struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		Channel Channel  `xml:"channel"`
	}

	var site_msg SiteMsg
	var datas []mydb.Data
	if category == "article" {
		tx := global.Db.Model(&mydb.Article{}).Where(
			"uid=? and id=?", uid, ref_id).First(&site_msg)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}
		tx = global.Db.Model(&mydb.Data{}).Joins("JOIN article ON data.ref_id = article.id").Where(
			"article.uid=? and data.category='article' and data.ref_id=?", uid, ref_id).Order("data.id desc").Limit(30).Find(&datas)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}
	} else if category == "video" {
		tx := global.Db.Model(&mydb.Video{}).Where(
			"uid=? and id=?", uid, ref_id).First(&site_msg)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}
		tx = global.Db.Model(&mydb.Data{}).Joins("JOIN video ON data.ref_id = video.id").Where(
			"video.uid=? and data.category='video' and data.ref_id=?", uid, ref_id).Order("data.id desc").Limit(30).Find(&datas)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}
	}

	var items []Item
	for _, data := range datas {
		title := data.Title
		url := data.URL
		date, _ := goutils.Unix2Str(data.CreatedAt)
		description := data.Description

		var item Item
		item.Title = title
		item.Link = url
		item.PubDate = date
		item.Description = description

		items = append(items, item)
	}

	rss := Rss{Version: "2.0"}
	rss.Channel = Channel{
		Title: site_msg.Name,
		Link:  site_msg.Link,
		Items: items,
	}
	c.XML(200, rss)
}
