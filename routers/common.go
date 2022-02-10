/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-06 17:14:01
 * @LastEditTime: 2022-02-08 23:02:20
 */
package routers

import (
	"RssSub/global"
	"RssSub/mydb"
)

// 记录后台计划任务、程序错误信息等
func loggerMsg(username, action, data string) {
	msg := mydb.Message{Uname: username, Action: action, Data: data, Status: "unread"}
	err := global.Db.Create(&msg).Error
	if err != nil {
		global.Log.Error(err.Error())
	}
}
