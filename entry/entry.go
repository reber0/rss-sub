/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-18 09:23:30
 * @LastEditTime: 2023-03-26 16:09:50
 */
package entry

import (
	"crypto/tls"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/reber0/go-common/mylog"
	"github.com/reber0/go-common/utils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/mydb"
)

func AppInit() {
	global.RootPath, _ = os.Getwd()
	global.Log = mylog.New().IsShowCaller(true).IsToFile(true).Logger()

	if !utils.IsFileExist(global.RootPath + "/data/data.db") {
		mydb.DbInit()
	}

	global.Db = mydb.GetDbConn()

	global.Client = resty.New()
	global.Client.SetTimeout(time.Duration(20) * time.Second)
	global.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	global.Client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:78.0) Gecko/20100101 Firefox/78.0")
	global.Client.SetHeader("Cookie", "buvid3=a") // b 站 api 获取需要带 cookie
}
