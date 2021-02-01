#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-06 09:44:54
@LastEditTime : 2021-01-31 17:31:54
'''

import re
import time
import demjson
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import to_dict
from sqlmodule import Video
from sqlmodule import Data

from libs.auth import get_username
from libs.request import req, ReqExceptin
from libs.common import logger_msg

from setting import rss_sqlite_uri
from setting import logger


class VideoClass(object):
    def __init__(self):
        super(VideoClass, self).__init__()

    def get_all_video_msg(self):
        """
        获取未完结的站点的信息
        """
        video_msg_list = list()
        try:
            with session_maker(rss_sqlite_uri) as db_session:
                results = db_session.query(Video).filter(Video.status != "已完结").all()
                video_msg_list = to_dict(results)
        except Exception as e:
            logger.error(str(e))
            return []
        else:
            # print(video_msg_list)
            return video_msg_list

    def get_video_exist_url(self, site_id):
        """
        获取一个动漫的现有链接
        """
        video_msg_list = list()
        video_url_list = list()
        try:
            with session_maker(rss_sqlite_uri) as db_session:
                results = db_session.query(Data).filter_by(ref_id=site_id).filter_by(category='video').all()
                video_msg_list = to_dict(results)
        except Exception as e:
            logger.error(str(e))
        else:
            if video_msg_list:
                video_url_list = [video_msg["url"] for video_msg in video_msg_list]
        return video_url_list

    def get_new_video_msg(self, username, name, link, video_url_list):
        """
        获取一个动漫的新剧集的信息
        """
        href_text_list = list()
        status = ""

        if "bilibili" in link:
            href_text_list, status, video_type = self.bilibili(link)
        elif "acfun" in link:
            href_text_list, status, video_type = self.acfun(link)
        elif "yhdm" in link:
            href_text_list, status, video_type = self.yhdm(username, name, link)

        new_video_msg_list = list()
        for href_text in href_text_list:
            url, title = href_text
            # html = "<iframe src={}></iframe>".format(url)
            html = "<h2>双指向左滑动</h2>"
            if url not in video_url_list:
                new_video_msg_list.append((title, url, html))

        return new_video_msg_list, status, video_type

    def check(self):
        video_msg_list = self.get_all_video_msg()
        for video_msg in video_msg_list:
            site_id = video_msg.get("id")
            user_id = video_msg.get("user_id")
            name = video_msg.get("name")
            link = video_msg.get("link")
            status = video_msg.get("status")
            username = get_username(user_id=user_id)

            # logger.info("Video check: {}".format(link))
            video_url_list = self.get_video_exist_url(site_id)
            new_video_msg_list, curr_status, video_type = self.get_new_video_msg(username, name, link, video_url_list)

            article_data = list()
            for new_video_msg in new_video_msg_list:
                title, url, html = new_video_msg
                date = func.now()

                article = Data(ref_id=site_id, category="video", title=title, url=url, date=date, description=html, status="unread")
                article_data.append(article)

            if article_data:
                message_data_list = [x[0] for x in new_video_msg_list] # 获取 title 的 list 用于保存
                logger_msg(msg_type="system", username="schedule", action="video check: {}".format(name), data=str(message_data_list))
                logger_msg(msg_type="user", username=username, action="{} 已更新".format(name), data=str(message_data_list))

                try:
                    with session_maker(rss_sqlite_uri) as db_session:
                        db_session.add_all(article_data)
                        if curr_status != status:
                            affect_num = db_session.query(Video).filter_by(id=site_id).update({"status": curr_status})
                            if affect_num:
                                logger.info("check {} success".format(name))
                except Exception as e:
                    logger.error(str(e))
                else:
                    logger.info("Video check: {} 获取 {} 个新资源".format(name, len(article_data)))

    def bilibili(self, url):
        href_text_list = list()
        status = ""
        vedio_type = ""

        # 用户信息接口
        user_info_api = "https://api.bilibili.com/x/space/acc/info?mid={}"
        # up主视频接口
        up_video_api = "https://api.bilibili.com/x/space/arc/search?mid={}&ps=50&tid=0&pn=1&keyword=&order=pubdate&jsonp=jsonp"
        # up主投稿文章接口
        up_article_api = "https://api.bilibili.com/x/space/article?mid={}&pn=1&ps=12&sort=publish_time&jsonp=jsonp"
        #番剧接口
        bangumi_api = "https://api.bilibili.com/pgc/view/web/season?season_id={}"

        # 获取视频列表
        if "space.bilibili.com" in url:
            up_id = url.strip("/").split("/")[-1]
            try:
                resp = req.get(up_video_api.format(up_id))
                result = resp.json()
            except ReqException as e:
                logger.error(e)
            else:
                if result.get("code") == 0:
                    vlist = result.get("data").get("list").get("vlist")
                    for v in vlist:
                        title = v.get("title")
                        bvid = v.get("bvid")
                        url = "https://www.bilibili.com/video/{}".format(bvid)
                        href_text_list.append((url, title))

                href_text_list = href_text_list[::-1]
                status = "连载中"
                vedio_type = "bilibili_up"
        elif "www.bilibili.com/bangumi" in url:
            try:
                resp = req.get(url)
                m = re.search(r'season_id":(\d+),', resp.text, re.M|re.S)
                season_id = m.group(1)
            except ReqException as e:
                logger.error(e)
            else:
                resp = req.get(bangumi_api.format(season_id))
                result = resp.json()
                if result.get("code") == 0:
                    result = result.get("result")
                    series_title = result.get("share_copy")

                    episodes = result.get("episodes")
                    for episode in episodes:
                        badge = episode.get("badge") # 值为 预告 或为空
                        if badge != "预告":
                            title = episode.get("share_copy").replace(series_title, "")
                            url = episode.get("share_url")
                            href_text_list.append((url, title))

                    status = result.get("publish").get("is_finish")
                    status = "已完结" if status else "连载中"
                    vedio_type = "bilibili_bangumi"
        return href_text_list, status, vedio_type

    def acfun(self, url):
        href_text_list = list()
        status = ""
        vedio_type = ""

        if "www.acfun.cn/bangumi" in url:
            try:
                html = req.get(url).text

                status = re.search(r'extendsStatus":"(.*?)",', html, re.S|re.M).group(1)
                bangumiList = re.search(r'window\.bangumiList = (.*?);', html).group(1)
            except ReqException as e:
                logger.error(e)
            else:
                items = demjson.decode(bangumiList).get("items")
                items = sorted(items, key=lambda x:x["priority"], reverse=True)
                for item in items:
                    title = item.get("title")
                    episodeName = item.get("episodeName")
                    title = "{} {}".format(title, episodeName)

                    bangumiId = item.get("bangumiId")
                    itemId = item.get("itemId")
                    priority = int(item.get("priority"))
                    if priority != 1000:
                        tmp_id = "{:1<6d}".format(priority)
                        url = "https://www.acfun.cn/bangumi/aa{}_{}_{}".format(bangumiId, tmp_id, itemId)
                        href_text_list.append((url, title))

                href_text_list = href_text_list[::-1]
                vedio_type = "acfun_bangumi"
        elif "www.acfun.cn/u" in url:
            try:
                html = req.get(url).text
            except Exception as e:
                logger.error(str(e))
            else:
                a_href_list = re.findall(r'<a href="(.*?)" target="_blank" class="ac-space-video', html)
                movurl_list = ["https://www.acfun.cn"+a_href for a_href in a_href_list]
                title_list = re.findall(r'<figcaption>.*?<p class="title line" title=".*?">(.*?)</p>', html, re.S|re.M)

                for index,title in enumerate(title_list):
                    url = movurl_list[index]
                    href_text_list.append((url, title))

                href_text_list = href_text_list[::-1]
                status = "连载中"
                vedio_type = "acfun_up"

        return href_text_list, status, vedio_type

    def yhdm(self, username, name, url):
        href_text_list = list()
        status = ""
        vedio_type = "yhdm"

        try:
            resp = req.get(url)
            resp.encoding = resp.apparent_encoding
            html = resp.text
        except ReqExceptin as error_msg:
            logger.error(error_msg)
            logger_msg(msg_type="system", username="schedule", action="video check: {}".format(name), data=str(error_msg))
            logger_msg(msg_type="user", username=username, action="{} 更新".format(name), data=str(error_msg))
        else:
            href_text_tag_list = re.findall(r'<li><a href="(/v/.*?)" target="_blank">(.*?)</a>', html, re.S|re.M)
            for href_text in href_text_tag_list:
                href = href_text[0]
                url = "http://www.yhdm.io"+href
                title = href_text[1]
                if "-pv" not in url:
                    href_text_list.append((url, title))

            m = re.search(r'href="#commen">.*?<p>(.*?)</p>', html, re.S|re.M)
            if m:
                status = m.group(1)
                if status:
                    status = "已完结" if "全" in status else "连载中"
                else:
                    status = "即将上映"

                if status != "已完结":
                    # print('dddddd')
                    href_text_list = href_text_list[::-1]
            # print(href_text_list)
        return href_text_list, status, vedio_type

def video_check():
    logger.info("==> video check...")

    logger_msg(msg_type="system", username="schedule", action="start video check")
    _video = VideoClass()
    _video.check()
