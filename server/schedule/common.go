/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-06 17:14:01
 * @LastEditTime: 2024-12-04 13:24:48
 */
package schedule

import (
	"github.com/reber0/rss-sub/server/global"
	"github.com/reber0/rss-sub/server/mydb"
)

// 记录后台计划任务、程序错误信息等
func loggerMsg(username, action, data string) {
	msg := mydb.Message{Uname: username, Action: action, Data: data, Status: "unread"}
	err := global.Db.Create(&msg).Error
	if err != nil {
		global.Log.Error(err.Error())
	}
}
