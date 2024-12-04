/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-10 13:35:18
 * @LastEditTime: 2024-12-04 13:23:55
 */
package routers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/reber0/rss-sub/server/global"
	"github.com/reber0/rss-sub/server/middleware"
	"github.com/reber0/rss-sub/server/mydb"
)

var menu_root = `[{
    "title": "主页",
    "icon": "layui-icon-home",
    "jump": "/"
}, {
    "title": "文章",
    "icon": "layui-icon-read",
    "list": [{
        "name": "article_site_add",
        "title": "添加 Site",
        "jump": "article/site/add"
    }, {
        "name": "article_site_manage",
        "title": "管理 Site",
        "jump": "article/site/manage"
    }, {
        "name": "article_data_list",
        "title": "管理 Article",
        "jump": "article/data/list"
    }]
}, {
    "title": "视频",
    "icon": "layui-icon-video",
    "list": [{
        "name": "video_site_manage",
        "title": "管理 Site",
        "jump": "video/site/manage"
    }, {
        "name": "video_data_list",
        "title": "管理 Video",
        "jump": "video/data/list"
    }]
}, {
    "name": "user",
    "title": "用户",
    "icon": "layui-icon-user",
    "list": [{
        "name": "user",
        "title": "用户管理",
        "jump": "user/list"
    }]
}, {
    "name": "set",
    "title": "设置",
    "icon": "layui-icon-set",
    "list": [{
        "name": "system",
        "title": "系统设置",
        "spread": true,
        "list": [{
            "name": "website",
            "title": "网站设置"
        }, {
            "name": "email",
            "title": "邮件服务"
        }]
    }, {
        "name": "user",
        "title": "我的设置",
        "spread": true,
        "list": [{
            "name": "info",
            "title": "基本资料"
        }, {
            "name": "password",
            "title": "修改密码"
        }]
    }]
}]`

var menu_user = `[{
    "title": "主页",
    "icon": "layui-icon-home",
    "jump": "/"
}, {
    "title": "文章",
    "icon": "layui-icon-read",
    "list": [{
        "name": "article_site_add",
        "title": "添加 Site",
        "jump": "article/site/add"
    }, {
        "name": "article_site_list",
        "title": "管理 Site",
        "jump": "article/site/list"
    }, {
        "name": "article_data_list",
        "title": "管理 Article",
        "jump": "article/data/list"
    }]
}, {
    "title": "视频",
    "icon": "layui-icon-video",
    "list": [{
        "name": "video_site_list",
        "title": "管理 Site",
        "jump": "video/site/list"
    }, {
        "name": "video_data_list",
        "title": "管理 Video",
        "jump": "video/data/list"
    }]
}, {
    "name": "set",
    "title": "设置",
    "icon": "layui-icon-set",
    "list": [{
        "name": "user",
        "title": "我的设置",
        "spread": true,
        "list": [{
            "name": "info",
            "title": "基本资料"
        }, {
            "name": "password",
            "title": "修改密码"
        }]
    }]
}]`

// 其他路由(copyright/menu)
func OtherRouter(r *gin.Engine) {
	setGroup := r.Group("/api/other")

	setGroup.POST("/copyright", copyright)
	setGroup.POST("/menu", middleware.JWTAuth(), menu)
}

// copyright
func copyright(c *gin.Context) {
	var copyright string
	err := global.Db.Model(&mydb.Config{}).Select("value").Where("key = ?", "copyright").First(&copyright).Error
	if err != nil {
		global.Log.Error(err.Error())
	}

	year := strconv.Itoa(time.Now().Year())
	copyright = strings.Replace(copyright, "{year}", year, 1)

	c.JSON(200, gin.H{
		"code":      0,
		"copyright": copyright,
	})
}

// 获取左侧目录
func menu(c *gin.Context) {
	userId := c.GetString("uid")
	_, role := GetUserMsg(userId)

	if role == "root" {
		res, err := simplejson.NewJson([]byte(menu_root))
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		c.JSON(200, gin.H{
			"code": 0,
			"data": res,
		})
	} else if role == "user" {
		res, err := simplejson.NewJson([]byte(menu_user))
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		c.JSON(200, gin.H{
			"code": 0,
			"data": res,
		})
	}
}
