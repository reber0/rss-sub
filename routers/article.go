/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:52:53
 * @LastEditTime: 2022-04-30 09:50:38
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
	"github.com/reber0/go-common/parse"
	"github.com/reber0/go-common/utils"
)

// Article Site 相关路由
func ArticleRouter(r *gin.Engine) {
	articleGroup := r.Group("/api/article").Use(middleware.JWTAuth(), middleware.Action())
	{
		articleGroup.POST("/check", articleCheckRegex)
		articleGroup.POST("/list", articleSiteList)
		articleGroup.POST("/add", articleSiteAdd)
		articleGroup.POST("/update", articleSiteUpdate)
		articleGroup.POST("/delete", articleSiteDelete)
		articleGroup.POST("/search", articleSiteSearch)
	}
}

func articleCheckRegex(c *gin.Context) {
	type PostData struct {
		Link  string `form:"link" json:"link"`
		Regex string `form:"regex" json:"regex"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		base_url := parse.NewParseURL(postJson.Link).BaseURL()
		resp, err := global.Client.Get(postJson.Link)
		if resp != nil {
			html := resp.Html()

			article_tag := make([]map[string]string, 0, 100)
			reg := regexp.MustCompile(`(?sm)` + postJson.Regex)
			result := reg.FindAllStringSubmatch(html, -1)
			for _, href_text := range result {
				a_href := strings.TrimSpace(href_text[1])
				if !strings.HasPrefix(a_href, "http") {
					a_href = base_url + strings.TrimLeft(a_href, "/")
				}
				a_text := strings.TrimSpace(href_text[2])
				article_tag = append(article_tag, map[string]string{"title": a_text, "url": a_href})
			}

			c.JSON(200, gin.H{
				"code": 0,
				"data": article_tag[0:5],
			})
		}
		if err != nil {
			global.Log.Error(err.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "error",
			})
		}
	}
}

func articleSiteAdd(c *gin.Context) {
	type PostData struct {
		Name  string `form:"name" json:"name"`
		Link  string `form:"link" json:"link"`
		Regex string `form:"regex" json:"regex"`
	}

	postJson := PostData{}
	if err := c.BindJSON(&postJson); err != nil {
		global.Log.Error(err.Error())
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "检查失败",
		})
	} else {
		userId := c.GetString("uid")

		var domain string
		global.Db.Model(&mydb.Config{}).Select("value").Where("key='domain'").First(&domain)

		site := mydb.Article{UID: userId, Name: postJson.Name, Link: postJson.Link, Regex: postJson.Regex}
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

		rssPath := fmt.Sprintf("/api/data/rss/%s/article/%d", userId, refId)
		result = global.Db.Model(&mydb.Article{}).Where("id=?", refId).Update("rss", rssPath)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "添加失败",
			})
			return
		}

		rss := fmt.Sprintf("%s/api/data/rss/%s/article/%d", strings.TrimRight(domain, "/"), userId, refId)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("得到 rss 链接: <br>%s", rss),
		})
	}
}

func articleSiteList(c *gin.Context) {
	type PostData struct {
		PageIndex int `form:"page" json:"page"`
		PageSize  int `form:"limit" json:"limit"`
	}

	type RespData struct {
		ID        int    `json:"id"`
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

		var count int64
		var datas []RespData
		result := global.Db.Model(&mydb.Article{}).Select("id", "name", "link", "regex", "rss", "created_at").Where(
			"uid=? or ?='root'", userId, role).Order("id desc").Count(&count).Limit(postJson.PageSize).Offset((postJson.PageIndex - 1) * postJson.PageSize).Find(&datas)
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
			datas[index].Rss = strings.TrimRight(domain, "/") + data.Rss
			datas[index].CreatedAt = utils.UnixToStr(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":  0,
			"data":  datas,
			"count": count,
		})
	}
}

func articleSiteUpdate(c *gin.Context) {
	type PostData struct {
		ID    int    `form:"id" json:"id"`
		Name  string `form:"name" json:"name"`
		Link  string `form:"link" json:"link"`
		Regex string `form:"regex" json:"regex"`
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

		updateData := mydb.Article{Name: postJson.Name, Link: postJson.Link, Regex: postJson.Regex}
		result := global.Db.Model(&mydb.Article{}).Where(
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

func articleSiteDelete(c *gin.Context) {
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

		result := global.Db.Where("(uid=? or ?='root') and id=?", userId, role, id).Delete(&mydb.Article{})
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "删除失败",
			})
			return
		}
		result = global.Db.Where("category='article' and ref_id=?", id).Delete(&mydb.Data{})
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

func articleSiteSearch(c *gin.Context) {
	type PostData struct {
		PageIndex int    `form:"page" json:"page"`
		PageSize  int    `form:"limit" json:"limit"`
		KeyWord   string `form:"keyword" json:"keyword"`
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
		result := global.Db.Model(&mydb.Article{}).Select("id", "name", "link", "regex", "rss", "created_at").Where(
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
			datas[index].CreatedAt = utils.UnixToStr(data.CreatedAt)
		}

		c.JSON(200, gin.H{
			"code":  0,
			"data":  datas,
			"count": count,
		})
	}
}
