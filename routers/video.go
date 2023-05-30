/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:54:15
 * @LastEditTime: 2023-05-16 16:51:02
 */
package routers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/reber0/go-common/utils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/middleware"
	"github.com/reber0/rss-sub/mydb"
	"github.com/tidwall/gjson"
)

// Video Site 相关路由
func VideoRouter(r *gin.Engine) {
	videoGroup := r.Group("/api/video")

	videoGroup.Use(middleware.JWTAuth(), middleware.Action())
	videoGroup.POST("/add", videoSiteAdd)
	videoGroup.POST("/list", videoSiteList)
	videoGroup.POST("/update", videoSiteUpdate)
	videoGroup.POST("/delete", videoSiteDelete)
}

func videoSiteAdd(c *gin.Context) {
	type PostData struct {
		Link string `form:"link" json:"link"`
	}

	postJson := PostData{}
	if err := c.ShouldBindJSON(&postJson); err != nil {
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
		PageIndex    int    `form:"page" json:"page"`
		PageSize     int    `form:"limit" json:"limit"`
		KeyWord      string `from:"keyword" json:"keyword"`
		Status       string `form:"status" json:"status"`
		ExportIdList []int  `form:"export_id_list" json:"export_id_list"`
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
		status := postJson.Status
		exportIdList := postJson.ExportIdList

		var count int64
		var datas []RespData
		tx := global.Db.Model(&mydb.Video{}).Joins("JOIN user ON user.uid = video.uid")
		tx = tx.Select("video.id,user.uname,video.name,video.link,video.status,video.rss,video.created_at")
		tx = tx.Where("video.uid=? or ?='root'", userId, role)

		if exportIdList != nil {
			tx.Where("video.id in ?", exportIdList)
		}
		if keyword != "" {
			tx = tx.Where("video.name like ?", "%"+keyword+"%")
		}
		if status != "" {
			tx = tx.Where("video.status=?", status)
		}

		tx = tx.Order("video.id desc").Count(&count)
		tx = tx.Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
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
			datas[index].CreatedAt = utils.Unix2Str(data.CreatedAt)
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
	if err := c.ShouldBindJSON(&postJson); err != nil {
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
		DeleteIDList []int `form:"target_id_list" json:"target_id_list"`
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

		deleteIDList := postJson.DeleteIDList

		for _, deleteID := range deleteIDList {
			result := global.Db.Where("(uid=? or ?='root') and id=?", userId, role, deleteID).Delete(&mydb.Video{})
			if result.Error != nil {
				global.Log.Error(result.Error.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "删除失败",
				})
				return
			} else {
				result = global.Db.Where("category='video' and ref_id=?", deleteID).Delete(&mydb.Data{})
				if result.Error != nil {
					global.Log.Error(result.Error.Error())
					c.JSON(500, gin.H{
						"code": 500,
						"msg":  "删除失败",
					})
					return
				}
			}
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "删除成功",
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
		resp, _ := global.Client.R().Get(url)
		html := utils.EncodeToUTF8(resp)

		code := gjson.Get(html, "code").Int()
		if code == 0 {
			name = gjson.Get(html, "data.name").String()
		}
	} else if strings.HasPrefix(targetURL, "https://www.bilibili.com/bangumi") {
		resp, _ := global.Client.R().Get(targetURL)
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`season_id":(\d+),`)
		m := reg.FindStringSubmatch(html)
		if len(m) > 0 {
			seasonId := m[1]
			url := strings.Replace(bangumi_api, "{}", seasonId, 1)
			resp, _ = global.Client.R().Get(url)
			html := utils.EncodeToUTF8(resp)

			code := gjson.Get(html, "code").Int()
			if code == 0 {
				name = gjson.Get(html, "result.season_title").String()
			}
		}
	} else if strings.HasPrefix(targetURL, "https://www.acfun.cn/u") {
		resp, _ := global.Client.R().Get(targetURL)
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`<span class="name" data-username=(.*?)>`)
		m := reg.FindStringSubmatch(html)
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "https://www.acfun.cn/bangumi") {
		resp, _ := global.Client.R().Get(targetURL)
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`bangumiTitle":"(.*?)",`)
		m := reg.FindStringSubmatch(html)
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "https://www.ysjdm.net/") {
		resp, _ := global.Client.R().Get(targetURL)
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`h2 class="title">\s+(.*?)\s+</h2>`)
		m := reg.FindStringSubmatch(html)
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "http://www.yinghuacd.com/") {
		resp, _ := global.Client.R().Get(targetURL)
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`<h1>(.*?)</h1>`)
		m := reg.FindStringSubmatch(html)
		if len(m) > 0 {
			name = m[1]
		}
	} else if strings.HasPrefix(targetURL, "https://www.agemys.vip/") {
		resp, _ := global.Client.R().Get(targetURL)
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`detail_imform_name">(.*?)</h4`)
		m := reg.FindStringSubmatch(html)
		if len(m) > 0 {
			name = m[1]
		}
	}

	return strings.TrimSpace(name)
}
