/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 20:53:53
 * @LastEditTime: 2024-01-24 13:46:35
 */
package routers

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reber0/goutils"
	"github.com/reber0/rss-sub/global"
	"github.com/reber0/rss-sub/middleware"
	"github.com/reber0/rss-sub/mydb"
)

// 设置相关路由(网站/邮箱/个人资料/密码)
func UploadRouter(r *gin.Engine) {
	setGroup := r.Group("/api/upload")

	setGroup.Use(middleware.JWTAuth(), middleware.Action())
	setGroup.POST("/avatar", upAvatar)
}

func upAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "上传失败",
		})
		return
	} else {
		userId := c.GetString("uid")

		fileExt := path.Ext(file.Filename)
		fileSize := file.Size / 1024

		// 判断上传文件后缀是否非法
		if !goutils.IsInCol(strings.Trim(fileExt, "."), []string{"png", "jpg", "jpeg"}) {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "文件非法",
			})
			return
		}

		// 判断上传图片是否过大
		var UploadMaxSize int64
		tx := global.Db.Model(&mydb.Config{}).Select("value").Where("key='upload_max_size'").First(&UploadMaxSize)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "上传失败",
			})
			return
		} else {
			if fileSize > UploadMaxSize {
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  "文件过大",
				})
				return
			}
		}

		// 保存图片
		avatar := goutils.Md5([]byte(uuid.New().String())) + fileExt
		filename := global.RootPath + "/avatar/" + avatar
		if err := c.SaveUploadedFile(file, filename); err != nil {
			global.Log.Error(err.Error())
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "更新失败",
			})
			return
		}

		// 更新 user 表中的 avatar
		tx = global.Db.Model(&mydb.User{}).Where("uid=?", userId).Update("avatar", avatar)
		if tx.Error != nil {
			global.Log.Error(tx.Error.Error())
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "更新失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"filename": avatar,
			},
			"msg": "更新头像成功",
		})
	}
}
