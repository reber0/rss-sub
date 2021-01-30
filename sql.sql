-- -------------------------------------------------------------
-- TablePlus 2.10(268)
--
-- https://tableplus.com/
--
-- Database: rss.sqlite3
-- Generation Time: 2021-01-26 14:33:20.7640
-- -------------------------------------------------------------


INSERT INTO "user" ("id", "user_id", "username", "password", "role", "email", "face", "add_time") VALUES ('1', 'dacf5b7314d5458ba1461b0dc409ad39', 'admin', '21232f297a57a5a743894a0e4a801fc3', 'root', 'admin@123.com', '61f7de565bc44630ab3bc2787e980252.png', '2021-01-25 13:03:40'),
('2', '4bad3c54f1ce430382d0786fca4f4daf', 'test', '098f6bcd4621d373cade4e832627b4f6', 'user', 'test1@test1.com', '2.png', '2021-01-25 18:41:54');


INSERT INTO "config" ("id", "key", "value") VALUES ('1', 'website', '{''sitename'': ''RssSub'', ''domain'': ''http://123.com:8083/'', ''upload_max_size'': ''1024'', ''title'': ''RssSub'', ''keywords'': ''RssSub, 博客, 番剧'', ''descript'': ''主要是用来生成 RSS 订阅，可以获取博客的新文章(定期正则爬取)，可以获取 B 站、A 站、樱花动漫的番剧更新情况。'', ''copyright'': ''Copyright © 2020-2020 Reber. All Rights Reserved.''}'),
('2', 'email', '{''smtp_server'': ''smtp.163.com'', ''smtp_port'': ''25'', ''send_email'': ''test@163.com'', ''send_nickname'': ''reber'', ''send_email_pwd'': ''123456''}'),
('3', 'root_menu', '[{
    "title": "主页",
    "icon": "layui-icon-home",
    "jump": "/"
}, {
    "title": "Article",
    "icon": "layui-icon-read",
    "list": [{
        "name": "article_add_site",
        "title": "添加 Site",
        "jump": "article/add_site"
    }, {
        "name": "article_list_site",
        "title": "管理 Site",
        "jump": "article/list_site"
    }, {
        "name": "article_list_article",
        "title": "管理 Article",
        "jump": "article/list_article"
    }]
}, {
    "title": "Video",
    "icon": "layui-icon-video",
    "list": [{
        "name": "video_list_site",
        "title": "管理 Site",
        "jump": "video/list_site"
    }, {
        "name": "video_list_video",
        "title": "管理 Video",
        "jump": "video/list_video"
    }]
}, {
    "name": "user",
    "title": "用户",
    "icon": "layui-icon-user",
    "list": [{
        "name": "user",
        "title": "用户管理",
        "jump": "user/list_user"
    }]
}, {
    "name": "set",
    "title": "设置",
    "icon": "layui-icon-set",
    "list": [{
        "name": "system",
        "title": "系统设置",
        "spread": true,
        "list": [{
            "name": "website",
            "title": "网站设置"
        }, {
            "name": "email",
            "title": "邮件服务"
        }]
    }, {
        "name": "user",
        "title": "我的设置",
        "spread": true,
        "list": [{
            "name": "info",
            "title": "基本资料"
        }, {
            "name": "password",
            "title": "修改密码"
        }]
    }]
}]'),
('4', 'user_menu', '[{
    "title": "主页",
    "icon": "layui-icon-home",
    "jump": "/"
}, {
    "title": "Article",
    "icon": "layui-icon-read",
    "list": [{
        "name": "article_add_site",
        "title": "添加 Site",
        "jump": "article/add_site"
    }, {
        "name": "article_list_site",
        "title": "管理 Site",
        "jump": "article/list_site"
    }, {
        "name": "article_list_article",
        "title": "管理 Article",
        "jump": "article/list_article"
    }]
}, {
    "title": "Video",
    "icon": "layui-icon-video",
    "list": [{
        "name": "video_list_site",
        "title": "管理 Site",
        "jump": "video/list_site"
    }, {
        "name": "video_list_video",
        "title": "管理 Video",
        "jump": "video/list_video"
    }]
}, {
    "name": "set",
    "title": "设置",
    "icon": "layui-icon-set",
    "list": [{
        "name": "user",
        "title": "我的设置",
        "spread": true,
        "list": [{
            "name": "info",
            "title": "基本资料"
        }, {
            "name": "password",
            "title": "修改密码"
        }]
    }]
}]'),
('5', 'jwt_deadline', '{
  "hours": 12,
  "minutes": 0
}'),
('6', 'base_xml', '<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
    <channel>
        <title>{name}</title>
        <link>{link}</link>
        {items}
    </channel>
</rss>'),
('7', 'base_item', '<item>
    <title><![CDATA[{title}]]></title>
    <link><![CDATA[{url}]]></link>
    <pubDate>{date}</pubDate>
    <description><![CDATA[{description}]]></description>
</item>'),
('8', 'check_article_rate', '{
  "hours": 6,
  "minutes": 0
}'),
('9', 'check_video_rate', '{
  "hours": 3,
  "minutes": 0
}');

