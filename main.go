/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 15:14:04
 * @LastEditTime: 2022-02-11 14:57:29
 */
package main

import (
	"RssSub/entry"
	"RssSub/global"
	"RssSub/routers"
	"RssSub/schedule"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "embed"
)

//go:embed web
var web embed.FS

func main() {
	entry.AppInit()

	// 启动计划任务
	go schedule.Start()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 添加路由
	routers.OtherRouter(r)
	routers.ArticleRouter(r)
	routers.VideoRouter(r)
	routers.DataRouter(r)
	routers.MessageRouter(r)
	routers.UserRouter(r)
	routers.SetRouter(r)
	routers.UploadRouter(r)

	// 映射头像文件
	r.Static("/avatar", global.RootPath+"/data/avatar")

	// 处理 web/dist 下的文件
	dist, _ := fs.Sub(web, "web/dist")
	r.StaticFS("/dist", http.FS(dist))
	layui, _ := fs.Sub(web, "web/start/layui")
	r.StaticFS("/layui", http.FS(layui))

	// 处理 web/start 下的文件
	start, _ := fs.Sub(web, "web/start")
	r.GET("/", gin.WrapH(http.FileServer(http.FS(start))))

	global.Log.Info("startup service: 0.0.0.0:8082")
	if err := r.Run("0.0.0.0:8082"); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
