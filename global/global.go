/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-05 11:23:54
 * @LastEditTime: 2022-01-20 17:30:31
 */
package global

import (
	"RssSub/myreq"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	RootPath string
	Log      *zap.Logger

	Db     *gorm.DB
	Client *myreq.Client
)
