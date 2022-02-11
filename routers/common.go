/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-06 17:14:01
 * @LastEditTime: 2022-02-11 10:42:56
 */
package routers

import (
	"RssSub/global"
	"RssSub/mydb"
)

// 根据 UserID 获取用户名、角色
func GetUserMsg(uid string) (uname, role string) {
	var data mydb.User
	result := global.Db.Model(&mydb.User{}).Select("uname, role").Where("uid = ?", uid).First(&data)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
	}
	return data.Uname, data.Role
}

// 根据 UserID 获取用户头像名
func GetAvatar(uid string) (avatar string) {
	result := global.Db.Model(&mydb.User{}).Select("avatar").Where("uid = ?", uid).First(&avatar)
	if result.Error != nil {
		global.Log.Error(result.Error.Error())
	}
	return avatar
}

// 记录后台计划任务、程序错误信息等
func loggerMsg(username, action, data string) {
	msg := mydb.Message{Uname: username, Action: action, Data: data, Status: "unread"}
	err := global.Db.Create(&msg).Error
	if err != nil {
		global.Log.Error(err.Error())
	}
}
