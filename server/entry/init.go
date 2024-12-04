// Package
// Author: reber
// Mail: reber0ask@qq.com
// Date: 2023-03-17 13:17:30
// LastEditTime: 2023-06-20 21:12:32

/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-18 09:23:30
 * @LastEditTime: 2024-12-04 12:33:26
 */
package entry

import (
	"crypto/tls"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/reber0/goutils"
	"github.com/reber0/rss-sub/server/global"
	"github.com/reber0/rss-sub/server/mydb"
)

// AppInit 初始化
func AppInit() {
	global.RootPath, _ = os.Getwd()
	global.Log = goutils.NewLog().IsShowCaller(true).IsToFile(true).L()

	if !goutils.PathExists(global.RootPath + "/data/data.db") {
		os.Mkdir(global.RootPath+"/data", 0755)
		os.Mkdir(global.RootPath+"/data/avatar", 0755)
		// 创建数据库
		file, err := os.Create(global.RootPath + "/data/data.db")
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
		file.Close()
		mydb.DbInit()
	}

	global.Db = mydb.GetDbConn()

	global.Client = resty.New()
	global.Client.SetTimeout(time.Duration(20) * time.Second)
	global.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	global.Client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:78.0) Gecko/20100101 Firefox/78.0")
	global.Client.SetHeader("Cookie", "buvid3=a") // b 站 api 获取需要带 cookie
	// global.Client.SetProxy("http://127.0.0.1:7890")
}
