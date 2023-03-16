/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 21:12:25
 * @LastEditTime: 2023-03-16 21:12:43
 */
package schedule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/reber0/go-common/parse"
	"github.com/reber0/go-common/utils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/mydb"
	"github.com/tidwall/gjson"
)

func checkArticle() {
	global.Log.Info("checkArticle...")

	articleAllSiteMsg := getArticleAllSiteMsg()
	for _, siteMsg := range articleAllSiteMsg {
		site_id := siteMsg.ID
		user_id := siteMsg.UID
		name := siteMsg.Name
		link := siteMsg.Link
		_type := siteMsg.Type
		match := siteMsg.Match

		var username string
		result := global.Db.Model(&mydb.User{}).Select("uname").Where("uid=?", user_id).First(&username)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
		}

		global.Log.Info(fmt.Sprintf("schedule article check ==> %s", name))

		articleURLSlice := getArticleURL(site_id)

		var newArticleMsgList [][]string
		var err error
		if _type == "regex" {
			newArticleMsgList, err = getNewSiteArticleMsg(link, match, articleURLSlice)
		} else if _type == "wechat" {
			newArticleMsgList, err = getNewWeChatArticleMsg(link, match, articleURLSlice)
		}
		if err != nil {
			loggerMsg(username, name+" 更新", err.Error())
		} else {
			if len(newArticleMsgList) > 0 {
				// 保存新的文章信息
				var titleSlice []string
				for _, newArticleMsg := range newArticleMsgList {
					title := newArticleMsg[0]
					url := newArticleMsg[1]
					data := mydb.Data{RefId: site_id, Category: "article", Title: title, URL: url, Status: "unread"}
					global.Db.Create(&data)

					titleSlice = append(titleSlice, title)
				}
				loggerMsg(username, name+" 更新", utils.SliceToString(titleSlice))
			}
		}
	}
}

// 获取所有站点的信息
func getArticleAllSiteMsg() []mydb.Article {
	var siteMsgSlice []mydb.Article
	result := global.Db.Model(&mydb.Article{}).Find(&siteMsgSlice)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
	}
	return siteMsgSlice
}

// 获取一个博客已获取的文章链接
func getArticleURL(articleID int) []string {
	var articleURLSlice []string

	var dataSlice []mydb.Data
	result := global.Db.Model(&mydb.Data{}).Where("category='article' and ref_id = ?", articleID).Find(&dataSlice)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
	} else {
		for _, data := range dataSlice {
			articleURLSlice = append(articleURLSlice, data.URL)
		}
	}

	return articleURLSlice
}

// 获取一个博客更新的文章信息
func getNewSiteArticleMsg(link, match string, articleURLSlice []string) ([][]string, error) {
	var newArticleMsgList [][]string

	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newArticleMsgList, err
	} else {
		html := utils.EncodeToUTF8(resp)
		baseURL := parse.NewParseURL(link).BaseURL()

		reg := regexp.MustCompile(`(?sm)` + match)
		href_text_list := reg.FindAllStringSubmatch(html, -1)
		for _, href_text := range href_text_list {
			a_href := strings.TrimSpace(href_text[1])
			if !strings.HasPrefix(a_href, "http") {
				a_href = baseURL + strings.TrimLeft(a_href, "/")
			}
			a_text := strings.TrimSpace(href_text[2])

			// 暂存新文章
			if !utils.InSlice(a_href, articleURLSlice) {
				newArticleMsgList = append(newArticleMsgList, []string{a_text, a_href})
			}
		}

		newArticleMsgList = utils.SliceReverse(newArticleMsgList)
	}

	return newArticleMsgList, nil
}

// 获取一个公众号更新的文章信息
func getNewWeChatArticleMsg(link, match string, articleURLSlice []string) ([][]string, error) {
	var newArticleMsgList [][]string

	link = fmt.Sprintf("https://mp.weixin.qq.com/mp/appmsgalbum?action=getalbum&album_id=%s&count=10&is_reverse=0&f=json", match)

	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newArticleMsgList, err
	} else {
		html := utils.EncodeToUTF8(resp)

		article_list := gjson.Get(html, "getalbum_resp.article_list").Array()
		for _, article := range article_list {
			a_href := article.Get("url").String()
			a_href = strings.Split(a_href, "&chksm")[0]
			a_text := article.Get("title").String()

			// 暂存新文章
			if !utils.InSlice(a_href, articleURLSlice) {
				newArticleMsgList = append(newArticleMsgList, []string{a_text, a_href})
			}
		}

		newArticleMsgList = utils.SliceReverse(newArticleMsgList)
	}

	return newArticleMsgList, nil
}
