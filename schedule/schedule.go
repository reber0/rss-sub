/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-02-04 23:33:28
 * @LastEditTime: 2022-02-09 00:27:20
 */
package schedule

import (
	"fmt"

	"github.com/robfig/cron"
)

func Start() {
	crontab := cron.New()

	specArticle := fmt.Sprintf("0 0 */%d * * *", 6) // 整点执行，6/12/18/24 点执行
	crontab.AddFunc(specArticle, checkArticle)

	specVideo := fmt.Sprintf("0 */%d * * * *", 30) // 整点执行，每 30 分钟执行一次
	crontab.AddFunc(specVideo, checkVideo)

	crontab.Start()
	select {}
}
