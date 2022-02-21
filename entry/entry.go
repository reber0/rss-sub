/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-18 09:23:30
 * @LastEditTime: 2022-02-21 16:04:59
 */
package entry

import (
	"RssSub/global"
	"RssSub/mydb"
	"net/http"
	"os"
	"time"

	"github.com/reber0/go-common/mylog"
	"github.com/reber0/go-common/myreq"
	"github.com/reber0/go-common/utils"
)

func AppInit() {
	global.RootPath, _ = os.Getwd()
	global.Log = mylog.NewLogger()

	if !utils.IsExist(global.RootPath + "/data/data.db") {
		mydb.CreateDb()
	}

	global.Db = mydb.GetDbConn()

	transport := http.Transport{
		DisableKeepAlives:   true, // 向同一服务器发大量请求，设置为 false 保持长连接，防止 socket too many open file
		MaxIdleConns:        500,  // 所有 host 的连接池最大连接数量，默认无穷大
		MaxIdleConnsPerHost: 500,  // 每个 host 的连接池最大空闲连接数，默认 2
		MaxConnsPerHost:     500,  // 每个 host 的最大连接数量，防止 connection reset by peer
	}
	global.Client = myreq.New().SetTransport(&transport).SkipVerify(true).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:78.0) Gecko/20100101 Firefox/78.0").
		SetTimeout(time.Duration(20) * time.Second)
}
