
# RssSub

[![platform](https://img.shields.io/static/v1?label=platform&message=macOS%20|%20Linux&color=172b43)](https://github.com/reber0/Rpscan/tree/master)
[![golang](https://img.shields.io/static/v1?label=golang&message=1.18&color=346fb0)](https://go.dev/)

  写这个代码的起因是 irreader 现在订阅超过 10 个站就需要收费，原来使用 irreader 主要是看中它可以正则获取网站信息，现在收费就不用了，不要问为什么，问就是穷。  
  一直安装的有 Reeder，而 Reeder 只能订阅有 RSS 文件的站而不能通过正则获取信息，然后发现 Inoreader 等一些 RSS 订阅网站被墙导致 Reeder 获取不到更新，后续搜了下虽然发现了 RSSHub、TTRSS 等工具，但感觉不好用，索性自己写代码实现通过正则获取网站文章信息生成 xml 从而订阅这个功能，所以就有了 RssSub。

### 运行
* go run main.go

### 使用
* 通过 `http://127.0.0.1:8082/` 访问即可
    * 默认账户：管理员（admin/admin）、普通用户（test/test）
    * 管理员登录后，通过 设置-系统设置-网站设置 改网站域名
    * 修改密码
* 添加 Blog
    * 通过正则添加博客，添加成功即可获取对应的 RSS 链接
