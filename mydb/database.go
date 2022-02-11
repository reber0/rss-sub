/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2022-01-04 21:19:24
 * @LastEditTime: 2022-02-11 17:55:34
 */
package mydb

/*
column	指定 db 列名，否则 userid 自动转为 user_id、username 自动转为 user_name
size	指定列大小
type	列数据类型，同时可设置长度（其中 string 对应的默认为 varchar）
index	根据参数创建索引
unique	指定列为唯一，值不能重复
primaryKey	指定列为主键，这里 ID 自动为主键，其实可以不用设置，int 后面为空即可
*/

type User struct {
	// gorm.Model // 加上此行会自动添加字段 id、created_at、updated_at、deleted_at
	ID        int
	UID       string `gorm:"column:uid; size:32; not null; unique; comment:用户唯一 id(uuid)"`
	Uname     string `gorm:"column:uname; type:varchar(50); not null; comment:用户名"`
	PassWord  string `gorm:"column:password; type:varchar(32); not null; comment:密码"`
	Role      string `gorm:"column:role; type:varchar(10); not null; comment:用户身份，root/user/e.g."`
	Email     string `gorm:"column:email; type:varchar(100); not null; comment:用户邮箱"`
	Avatar    string `gorm:"column:avatar; type:varchar(40); not null; comment:头像图片名"`
	CreatedAt int    `gorm:"column:create_at; comment:添加时间"`
}

//设置表名为 user，否则自动为 users（默认是结构体的名的复数形式）
func (User) TableName() string {
	return "user"
}

type Config struct {
	// gorm.Model        // 加上此行会自动添加字段 id、created_at、updated_at、deleted_at
	ID    int
	Key   string `gorm:"column:key; type:varchar(30); not null; comment:配置的名字"`
	Value string `gorm:"column:value; type:text; not null; comment:配置的值"`
}

//设置表名为 config，否则自动为 configs（默认是结构体的名的复数形式）
func (Config) TableName() string {
	return "config"
}

type Article struct {
	ID        int
	UID       string `gorm:"column:uid; size:32; not null; comment:用户唯一 id(uuid)"`
	Name      string `gorm:"column:name; type:varchar(100); not null; comment:博客名字"`
	Link      string `gorm:"column:link; type:varchar(100); not null; comment:文章网站的网址"`
	Regex     string `gorm:"column:regex; type:text; not null; comment:正则"`
	Rss       string `gorm:"column:rss; type:varchar(100); comment:RSS 地址"`
	CreatedAt int    `gorm:"autoCreateTime; column:create_at; comment:添加时间"`
}

func (Article) TableName() string {
	return "article"
}

type Video struct {
	ID        int
	UID       string `gorm:"column:uid; size:32; not null; comment:用户唯一 id(uuid)"`
	Name      string `gorm:"column:name; type:varchar(100); not null; comment:系列名字，比如番剧名"`
	Link      string `gorm:"column:link; type:varchar(100); not null; comment:主页目录，比如番剧主页、UP 主主页的 URL"`
	Status    string `gorm:"column:status; type:varchar(100); comment:连载状态"`
	Rss       string `gorm:"column:rss; type:varchar(100); comment:RSS 地址"`
	CreatedAt int    `gorm:"column:create_at; comment:添加时间"`
}

func (Video) TableName() string {
	return "video"
}

type Data struct {
	ID          int
	RefId       int    `gorm:"column:ref_id; comment:article、video 表中站点对应的 id"`
	Category    string `gorm:"column:category; type:varchar(10); comment:表名，article、video，用于区分该条记录对应那个表里的数据"`
	Title       string `gorm:"column:title; type:text; comment:标题，文章名字、番剧每集名字"`
	URL         string `gorm:"column:url; type:text; comment:网址，文章链接、番剧每集 URL"`
	Description string `gorm:"column:description; type:text; comment:简介"`
	Status      string `gorm:"column:status; type:varchar(10); comment:状态，是否已读、已看"`
	CreatedAt   int    `gorm:"column:create_at; comment:添加时间"`
}

func (Data) TableName() string {
	return "data"
}

type Message struct {
	ID        int
	Uname     string `gorm:"column:uname; type:varchar(50); comment:用户名"`
	Action    string `gorm:"column:action; type:text; comment:执行的操作触发的 URI、计划任务动作"`
	Data      string `gorm:"column:data; type:text; comment:POST 的数据、得到的数据"`
	Status    string `gorm:"column:status; type:varchar(10); comment:状态，是否已读、已看"`
	CreatedAt int    `gorm:"column:create_at; comment:添加时间"`
}

func (Message) TableName() string {
	return "message"
}
