/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 22:20:57
 * @LastEditTime: 2022-06-02 00:00:21
 */
package mydb

import (
	"rsssub/global"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 返回数据库连接句柄
func GetDbConn() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(global.RootPath+"/data/data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func CreateDb() {
	Db := GetDbConn()

	// 自动迁移，创建表
	Db.AutoMigrate(&User{})
	Db.AutoMigrate(&Config{})
	Db.AutoMigrate(&Article{})
	Db.AutoMigrate(&Video{})
	Db.AutoMigrate(&Data{})
	Db.AutoMigrate(&Message{})

	// User 表添加数据
	Db.Create(&User{
		UID:      "dacf5b7314d5458ba1461b0dc409ad39",
		Uname:    "admin",
		PassWord: "21232f297a57a5a743894a0e4a801fc3",
		Role:     "root",
		Email:    "admin@123.com",
		Avatar:   "61f7de565bc44630ab3bc2787e980252.png",
	})
	Db.Create(&User{
		UID:      "4bad3c54f1ce430382d0786fca4f4daf",
		Uname:    "test",
		PassWord: "098f6bcd4621d373cade4e832627b4f6",
		Role:     "user",
		Email:    "test@test.com",
		Avatar:   "2.png",
	})

	// Config 表添加数据
	Db.Create(&Config{Key: "sitename", Value: "RssSub"})
	Db.Create(&Config{Key: "domain", Value: "http://127.0.0.1:8083/"})
	Db.Create(&Config{Key: "upload_max_size", Value: "1024"})
	Db.Create(&Config{Key: "title", Value: "RssSub"})
	Db.Create(&Config{Key: "keyword", Value: "RssSub, 博客, 番剧"})
	Db.Create(&Config{Key: "descript", Value: "主要是用来生成 RSS 订阅，可以获取博客的新文章(定期正则爬取)，可以获取 B 站、A 站、AGE.tv 的番剧更新情况。"})
	Db.Create(&Config{Key: "copyright", Value: "Copyright © 2020-2022 Reber. All Rights Reserved."})

	Db.Create(&Config{Key: "smtp_server", Value: "smtp.163.com"})
	Db.Create(&Config{Key: "smtp_port", Value: "25"})
	Db.Create(&Config{Key: "send_email", Value: "test@163.com"})
	Db.Create(&Config{Key: "send_nickname", Value: "reber"})
	Db.Create(&Config{Key: "send_email_pwd", Value: "123456"})

	Db.Create(&Config{Key: "check_article_rate", Value: `{"hours": 6, "minutes": 0, "seconds": 0}`})
	Db.Create(&Config{Key: "check_video_rate", Value: `{"hours": 0, "minutes": 30, "seconds": 0}`})
}
