/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 21:12:34
 * @LastEditTime: 2022-05-31 15:38:15
 */
package schedule

import (
	"fmt"
	"regexp"
	"rsssub/global"
	"rsssub/mydb"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/reber0/go-common/parse"
	"github.com/reber0/go-common/utils"
	"github.com/tidwall/gjson"
)

func checkVideo() {
	global.Log.Info("checkVideo...")

	videoAllSiteMsg := getVideoAllSiteMsg()
	for _, siteMsg := range videoAllSiteMsg {
		site_id := siteMsg.ID
		user_id := siteMsg.UID
		name := siteMsg.Name
		link := siteMsg.Link

		var username string
		result := global.Db.Model(&mydb.User{}).Select("uname").Where("uid=?", user_id).First(&username)
		if result.Error != nil {
			global.Log.Error(result.Error.Error())
		}

		global.Log.Info(fmt.Sprintf("schedule video check ==> %s\n", name))
		videoURLSlice := getVideoURL(site_id)
		newVideoMsgList, status, err := getNewVideoMsg(link, videoURLSlice)
		if err != nil {
			loggerMsg(username, name+" 更新", err.Error())
		} else {
			if len(newVideoMsgList) > 0 {
				if len(newVideoMsgList) > 30 {
					newVideoMsgList = newVideoMsgList[0:30]
				}

				// 保存新的番剧信息
				var titleSlice []string
				for _, newVideoMsg := range newVideoMsgList {
					title := newVideoMsg[0]
					url := newVideoMsg[1]
					data := mydb.Data{RefId: site_id, Category: "video", Title: title, URL: url, Status: "unread"}
					global.Db.Create(&data)

					titleSlice = append(titleSlice, title)
				}
				loggerMsg(username, name+" 更新", utils.SliceToString(titleSlice))

				// 更新番剧状态
				result := global.Db.Model(&mydb.Video{}).Where("id = ?", site_id).Update("status", status)
				if result.Error != nil {
					global.Log.Error(result.Error.Error())

				}
			}
		}
	}
}

// 获取未完结的站点的信息
func getVideoAllSiteMsg() []mydb.Video {
	var siteMsgSlice []mydb.Video
	result := global.Db.Model(&mydb.Video{}).Where("status != '已完结'").Find(&siteMsgSlice)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
	}
	return siteMsgSlice
}

// 获取一个番剧已获取的链接
func getVideoURL(videoID int) []string {
	var videoURLSlice []string

	var dataSlice []mydb.Data
	result := global.Db.Model(&mydb.Data{}).Where("category='video' and ref_id = ?", videoID).Find(&dataSlice)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
	} else {
		for _, data := range dataSlice {
			videoURLSlice = append(videoURLSlice, data.URL)
		}
	}

	return videoURLSlice
}

// 获取一个番剧更新的信息
func getNewVideoMsg(link string, videoURLSlice []string) ([][]string, string, error) {
	var status string
	var newVideoMsgList [][]string

	baseURL := parse.NewParseURL(link).BaseURL()

	var err error
	var href_text_list [][]string
	if strings.HasPrefix(link, "https://space.bilibili.com/") {
		href_text_list, status, err = bilibiliUp(link)
	} else if strings.HasPrefix(link, "https://www.bilibili.com/bangumi") {
		href_text_list, status, err = bilibiliBangumi(link)
	} else if strings.HasPrefix(link, "https://www.acfun.cn/u") {
		href_text_list, status, err = acfunUp(link)
	} else if strings.HasPrefix(link, "https://www.acfun.cn/bangumi") {
		href_text_list, status, err = acfunBangumi(link)
	} else if strings.HasPrefix(link, "https://www.agemys.com/") {
		href_text_list, status, err = age(link)
	} else if strings.HasPrefix(link, "http://www.yinghuacd.com/") {
		href_text_list, status, err = yinghuacd(link)
	}
	if err != nil {
		return newVideoMsgList, status, err
	} else {
		for _, href_text := range href_text_list {
			a_text := strings.TrimSpace(href_text[0])

			a_href := strings.TrimSpace(href_text[1])
			if !strings.HasPrefix(a_href, "http") {
				a_href = baseURL + strings.TrimLeft(a_href, "/")
			}

			// 暂存新剧集
			if !utils.InSlice(a_href, videoURLSlice) {
				newVideoMsgList = append(newVideoMsgList, []string{a_text, a_href})
			}
		}
	}

	return newVideoMsgList, status, nil
}

func bilibiliUp(link string) ([][]string, string, error) {
	var newVideoMsgList [][]string

	// up主视频接口
	up_video_api := "https://api.bilibili.com/x/space/arc/search?mid={mid}&ps=30&tid=0&pn=1&keyword=&order=pubdate&jsonp=jsonp"

	// 获取视频列表
	_tmp := strings.Split(strings.Trim(link, "/"), "/")
	up_id := _tmp[len(_tmp)-1]
	resp, err := global.Client.R().Get(strings.ReplaceAll(up_video_api, "{mid}", up_id))
	if err != nil {
		global.Log.Error(err.Error())
		return newVideoMsgList, "连载中", err
	} else {
		html := utils.EncodeToUTF8(resp)

		code := gjson.Get(html, "code").Int()
		if code == 0 {
			vlist := gjson.Get("html", "data.list.vlist").Array()
			for _, v := range vlist {
				title := v.Get("title").String()
				bvid := v.Get("bvid").String()
				url := strings.ReplaceAll("https://www.bilibili.com/video/{bvid}", "{bvid}", bvid)
				newVideoMsgList = append(newVideoMsgList, []string{title, url})
			}
			newVideoMsgList = utils.SliceReverse(newVideoMsgList)
		}
	}

	return newVideoMsgList, "连载中", nil
}

func bilibiliBangumi(link string) ([][]string, string, error) {
	var status string
	var newVideoMsgList [][]string

	// 番剧接口
	bangumi_api := "https://api.bilibili.com/pgc/view/web/season?season_id={season_id}"

	// 获取视频列表
	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newVideoMsgList, "连载中", err
	} else {
		html := utils.EncodeToUTF8(resp)

		reg := regexp.MustCompile(`season_id":(\d+),`)
		m := reg.FindStringSubmatch(html)
		season_id := m[1]

		resp, err := global.Client.R().Get(strings.ReplaceAll(bangumi_api, "{season_id}", season_id))
		if err != nil {
			global.Log.Error(err.Error())
			return newVideoMsgList, "连载中", err
		} else {
			html := utils.EncodeToUTF8(resp)

			code := gjson.Get(html, "code").Int()
			if code == 0 {
				// result := gjson.Get(html, "result")
				_title := gjson.Get(html, "result.title").String()

				// 获取正片
				episodes := gjson.Get(html, "result.episodes").Array()
				for _, episode := range episodes {
					badge := episode.Get("badge").String()
					if !strings.Contains(badge, "预告") {
						share_copy := episode.Get("share_copy").String()
						title := strings.TrimSpace(strings.ReplaceAll(share_copy, "《"+_title+"》", ""))
						url := episode.Get("share_url").String()
						newVideoMsgList = append(newVideoMsgList, []string{title, url})
					}
				}
				if gjson.Get(html, "result.publish.is_finish").Int() == 1 {
					status = "已完结"
				} else {
					status = "连载中"
				}
			}
		}

	}

	return newVideoMsgList, status, nil
}

func acfunUp(link string) ([][]string, string, error) {
	var newVideoMsgList [][]string

	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newVideoMsgList, "连载中", err
	} else {
		html := utils.EncodeToUTF8(resp)

		var movurl_list []string
		reg1 := regexp.MustCompile(`<a href="(.*?)" target="_blank" class="ac-space-video`)
		for _, a_href := range reg1.FindAllStringSubmatch(html, -1) {
			movurl_list = append(movurl_list, a_href[1])
		}

		reg2 := regexp.MustCompile(`(?sm)<figcaption>.*?<p class="title line" title="(.*?)">`)
		for index, title := range reg2.FindAllStringSubmatch(html, -1) {
			url := movurl_list[index]
			newVideoMsgList = append(newVideoMsgList, []string{title[1], url})
		}
		newVideoMsgList = utils.SliceReverse(newVideoMsgList)
	}

	return newVideoMsgList, "连载中", nil
}

func acfunBangumi(link string) ([][]string, string, error) {
	var status string
	var newVideoMsgList [][]string

	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newVideoMsgList, "连载中", err
	} else {
		html := utils.EncodeToUTF8(resp)

		reg1 := regexp.MustCompile(`(?sm)extendsStatus":"(.*?)",`)
		m1 := reg1.FindStringSubmatch(html)
		status = m1[1]

		reg2 := regexp.MustCompile(`window\.bangumiList = (.*?);`)
		m2 := reg2.FindStringSubmatch(html)
		bangumiList := m2[1]

		res, err := simplejson.NewJson([]byte(bangumiList))
		if err != nil {
			global.Log.Error(err.Error())
		}
		for _, item := range res.Get("items").MustArray() {
			item := item.(map[string]interface{})

			title := fmt.Sprintf("%s %s", item["title"], item["episodeName"])

			bangumiId := item["bangumiId"]
			priority := item["priority"]
			itemId := item["itemId"]
			if priority != "1000" {
				url := fmt.Sprintf("https://www.acfun.cn/bangumi/aa%s_%s_%s", bangumiId, priority, itemId)
				newVideoMsgList = append(newVideoMsgList, []string{title, url})
			}
		}
	}

	return newVideoMsgList, status, nil
}

func age(link string) ([][]string, string, error) {
	var status string
	var newVideoMsgList [][]string

	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newVideoMsgList, "连载中", err
	} else {
		html := utils.EncodeToUTF8(resp)
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			global.Log.Error(err.Error())
		}

		dom.Find(`div[style="display:block"]>ul>li>a`).Each(func(i int, node *goquery.Selection) {
			url, _ := node.Attr("href")
			title := node.Text()
			if !strings.Contains(strings.ToLower(title), "pv") && !strings.Contains(strings.ToLower(title), "无字") {
				newVideoMsgList = append(newVideoMsgList, []string{title, url})
			}
		})

		reg := regexp.MustCompile(`(?sm)播放状态.*?value">(.*?)</span>`)
		m := reg.FindStringSubmatch(html)
		status = m[1]
		if status == "完结" {
			status = "已完结"
		} else {
			status = "连载中"
		}
	}

	return newVideoMsgList, status, nil
}

func yinghuacd(link string) ([][]string, string, error) {
	var status string
	var newVideoMsgList [][]string

	resp, err := global.Client.R().Get(link)
	if err != nil {
		global.Log.Error(err.Error())
		return newVideoMsgList, "连载中", err
	} else {
		html := utils.EncodeToUTF8(resp)

		dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			global.Log.Error(err.Error())
		}

		movurl := dom.Find(`div[class="movurl"]`)
		movurl.Find(`ul>li>a`).Each(func(i int, node *goquery.Selection) {
			url, _ := node.Attr("href")
			title := node.Text()
			if !strings.Contains(strings.ToLower(title), "pv") {
				newVideoMsgList = append(newVideoMsgList, []string{title, url})
			}
		})

		title := dom.Find("title").Text()
		if strings.Contains(title, "全集") {
			status = "已完结"
		} else {
			status = "连载中"
			newVideoMsgList = utils.SliceReverse(newVideoMsgList)
		}
	}

	return newVideoMsgList, status, nil
}
