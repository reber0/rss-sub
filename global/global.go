/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-05 11:23:54
 * @LastEditTime: 2022-05-31 15:20:24
 */
package global

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	RootPath string
	Log      *zap.Logger

	Db     *gorm.DB
	Client *resty.Client
)
