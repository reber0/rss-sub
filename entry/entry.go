/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-18 09:23:30
 * @LastEditTime: 2022-09-09 10:45:00
 */
package entry

import (
	"crypto/tls"
	"os"
	"rsssub/global"
	"rsssub/mydb"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/reber0/go-common/mylog"
	"github.com/reber0/go-common/utils"
)

func AppInit() {
	global.RootPath, _ = os.Getwd()
	global.Log = mylog.NewLogger()

	if !utils.IsFileExist(global.RootPath + "/data/data.db") {
		mydb.DbInit()
	}

	global.Db = mydb.GetDbConn()

	global.Client = resty.New()
	global.Client.SetTimeout(time.Duration(20) * time.Second)
	global.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	global.Client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:78.0) Gecko/20100101 Firefox/78.0")
}
