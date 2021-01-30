# RssSub

[![platform](https://img.shields.io/static/v1?label=platform&message=macOS%20|%20Linux&color=172b43)](https://github.com/reber0/Rpscan/tree/master)
[![python](https://img.shields.io/static/v1?label=python&message=3.9&color=346fb0)](https://www.python.org/)

&emsp;&emsp;写这个代码的起因是 irreader 现在订阅超过 10 个站就需要收费，原来使用 irreader 主要是看中它可以正则获取网站信息，现在收费就不用了，不要问为什么，问就是穷。  
&emsp;&emsp;一直安装的有 Reeder，而 Reeder 只能订阅有 RSS 文件的站而不能通过正则获取信息，然后发现 Inoreader 等一些 RSS 订阅网站被墙导致 Reeder 获取不到更新，后续搜了下虽然发现了 RSSHub、TTRSS 等工具，但感觉不好用，索性自己写代码实现通过正则获取网站文章信息生成 xml 从而订阅这个功能，所以就有了 RssSub。

### 安装模块
pip3 install -r requirements.txt

### 修改配置
setting.py 配置说明：
* jwt_key  JWT 生成 token 的 key
* screct_key  flask 的 secret_key

### 使用
* ~~生成数据库~~
    * ~~cd /path/to/RssSub/backend/sqlmodule && python3 module.py~~
    * ~~连接数据库，执行 /path/to/RssSub/sql.sql 中的两条语句~~
* 运行程序
    * cd /path/to/RssSub
    * sudo supervisord -c supervisor.conf
    * sudo supervisorctl -c supervisor.conf start all
* 通过 `http://127.0.0.1:8083/` 访问即可（可在 supervisor.conf 里修改端口）
    * 默认账户：管理员（admin/admin）、普通用户（test/test）
    * 登录后改密码、改配置中的网站域名即可
* 添加 Blog
    * 通过正则添加博客，拿到生成的链接即可使用
