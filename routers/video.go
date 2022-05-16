/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:54:15
 * @LastEditTime: 2022-04-30 12:44:21
 */
package routers

import (
	"RssSub/global"
	"RssSub/middleware"
	"RssSub/mydb"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/reber0/go-common/utils"
)

// Video Site 相关路由
func VideoRouter(r *gin.Engine) {
	videoGroup := r.Group("/api/video").Use(middleware.JWTAuth(), middleware.Action())
	{
		videoGroup.POST("/add", videoSiteAdd)
		videoGroup.POST("/list", videoSiteList)
		videoGroup.POST("/update", videoSiteUpdate)
		videoGroup.POST("/delete", videoSiteDelete)
		videoGroup.POST("/search", videoSiteSearch)
	}
}

func videoSiteAdd(c *gin.Context) {
	type PostData struct {
		Link string `form:"link" json:"link"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "添加失败",
		})
	} else {
		userId := c.GetString("uid")

		link := strings.TrimSpace(postJson.Link)
		name := getName(link)

		var domain string
		global.Db.Model(&mydb.Config{}).Select("value").Where("key='domain'").First(&domain)

		site := mydb.Video{UID: userId, Name: name, Link: link, Status: "连载中"}
		result := global.Db.Create(&site)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "添加失败",
			})
			return
		}
		refId := site.ID

		rssPath := fmt.Sprintf("/api/data/rss/%s/video/%d", userId, refId)
		result = global.Db.Model(&mydb.Video{}).Where("id=?", refId).Update("rss", rssPath)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "添加失败",
			})
			return
		}

		rss := fmt.Sprintf("%s/api/data/rss/%s/video/%d", strings.TrimRight(domain, "/"), userId, refId)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("得到 rss 链接: <br>%s", rss),
		})
	}
}

func videoSiteList(c *gin.Context) {
	type PostData struct {
		PageIndex int `form:"page" json:"page"`
		PageSize  int `form:"limit" json:"limit"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		Uname     string `json:"uname,omitempty" gorm:"column:uname; type:varchar(50); comment:用户名"`
		Name      string `json:"name" gorm:"column:name; type:varchar(100); not null; comment:系列名字，比如番剧名"`
		Link      string `json:"link" gorm:"column:link; type:varchar(100); not null; comment:主页目录，比如番剧主页、UP 主主页的 URL"`
		Status    string `json:"status" gorm:"column:status; type:varchar(100); comment:连载状态"`
		Rss       string `json:"rss" gorm:"column:rss; type:varchar(100); comment:RSS 地址"`
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
		_, role := GetUserMsg(userId)

		var count int64
		var datas []RespData
		result := global.Db.Model(&mydb.Video{}).Joins("JOIN user ON user.uid = video.uid").Select("video.id,user.uname,video.name,video.link,video.status,video.rss,video.created_at").Where(
			"video.uid=? or ?='root'", userId, role).Order("video.id desc").Count(&count).Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

		var domain string
		global.Db.Model(&mydb.Config{}).Select("value").Where("key='domain'").First(&domain)

		for index, data := range datas {
			if role == "user" {
				datas[index].Uname = ""
			}
			datas[index].Rss = strings.TrimRight(domain, "/") + data.Rss
			datas[index].CreatedAt = utils.Unix2String(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":  0,
			"data":  datas,
			"count": count,
		})
	}
}

func videoSiteUpdate(c *gin.Context) {
	type PostData struct {
		ID     int    `form:"id" json:"id"`
		Name   string `form:"name" json:"name"`
		Link   string `form:"link" json:"link"`
		Status string `form:"status" json:"status"`
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
		_, role := GetUserMsg(userId)

		updateData := mydb.Video{Name: postJson.Name, Link: postJson.Link, Status: postJson.Status}
		result := global.Db.Model(&mydb.Video{}).Where(
			"(uid=? or ?='root') and id=?", userId, role, postJson.ID).Updates(updateData)
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

func videoSiteDelete(c *gin.Context) {
	type PostData struct {
		ID int `form:"id" json:"id"`
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
		_, role := GetUserMsg(userId)

		id := postJson.ID

		result := global.Db.Where("(uid=? or ?='root') and id=?", userId, role, id).Delete(&mydb.Video{})
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "删除失败",
			})
			return
		}
		result = global.Db.Where("category='video' and ref_id=?", id).Delete(&mydb.Data{})
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

func videoSiteSearch(c *gin.Context) {
	type PostData struct {
		PageIndex int    `form:"page" json:"page"`
		PageSize  int    `form:"limit" json:"limit"`
		KeyWord   string `form:"keyword" json:"keyword"`
		Status    string `form:"status" json:"status"`
	}

	type RespData struct {
		ID        int    `json:"id"`
		UID       string `json:"uid" gorm:"column:uid; size:32; not null; unique; comment:用户唯一 id(uuid)"`
		Name      string `json:"name" gorm:"column:name; type:varchar(100); not null; comment:博客名字"`
		Link      string `json:"link" gorm:"column:link; type:varchar(100); not null; comment:文章网站的网址"`
		Regex     string `json:"regex" gorm:"column:regex; type:text; not null; comment:正则"`
		Rss       string `json:"rss" gorm:"column:rss; type:varchar(100); comment:RSS 地址"`
		CreatedAt string `json:"created_at" gorm:"column:created_at; type:varchar(100); comment:添加时间"`
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
		_, role := GetUserMsg(userId)

		keyword := fmt.Sprintf("%%%s%%", postJson.KeyWord)

		var count int64
		var datas []RespData
		result := global.Db.Model(&mydb.Video{}).Select("id", "name", "link", "status", "rss", "created_at").Where(
			"(uid=? or ?='root') and name like ?", userId, role, keyword).Order("id desc").Count(&count).Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			return
		}

		for index, data := range datas {
			datas[index].CreatedAt = utils.Unix2String(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":  0,
			"data":  datas,
			"count": count,
		})
	}
}

// 获取番剧的名字
func getName(targetURL string) string {
	var name string

	user_info_api := "https://api.bilibili.com/x/space/acc/info?mid={}"        // 用户信息接口
	bangumi_api := "https://api.bilibili.com/pgc/view/web/season?season_id={}" // 番剧接口

	if strings.HasPrefix(targetURL, "https://space.bilibili.com/") {
		uSlice := strings.Split(strings.Trim(targetURL, "/"), "/")
		mid := uSlice[len(uSlice)-1]

		url := strings.Replace(user_info_api, "{}", mid, 1)
		resp, _ := global.Client.Get(url)
		jsonData := resp.Json()

		code := jsonData.Get("code").MustInt()
		if code == 0 {
			name = jsonData.Get("data").Get("name").MustString()
		}
	} else if strings.HasPrefix(targetURL, "https://www.bilibili.com/bangumi") {
		resp, _ := global.Client.Get(targetURL)

		reg := regexp.MustCompile(`season_id":(\d+),`)
		m := reg.FindStringSubmatch(resp.Html())
		if len(m) > 0 {
			seasonId := m[1]
			url := strings.Replace(bangumi_api, "{}", seasonId, 1)
			resp, _ = global.Client.Get(url)
			jsonData := resp.Json()

			code := jsonData.Get("code").MustInt()
			if code == 0 {
				name = jsonData.Get("result").Get("season_title").MustString()
			}
		}
	} else if strings.HasPrefix(targetURL, "https://www.acfun.cn/u") {
		resp, _ := global.Client.Get(targetURL)

		reg := regexp.MustCompile(`<span class="name" data-username=(.*?)>`)
		m := reg.FindStringSubmatch(resp.Html())
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "https://www.acfun.cn/bangumi") {
		resp, _ := global.Client.Get(targetURL)

		reg := regexp.MustCompile(`bangumiTitle":"(.*?)",`)
		m := reg.FindStringSubmatch(resp.Html())
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "https://www.agemys.com/") {
		resp, _ := global.Client.Get(targetURL)

		reg := regexp.MustCompile(`detail_imform_name">(.*?)<`)
		m := reg.FindStringSubmatch(resp.Html())
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "http://www.yinghuacd.com/") {
		resp, _ := global.Client.Get(targetURL)

		reg := regexp.MustCompile(`<h1>(.*?)</h1>`)
		m := reg.FindStringSubmatch(resp.Html())
		if len(m) > 0 {
			name = m[1]
		}
	}

	return name
}
