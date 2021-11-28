#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-06 09:44:54
LastEditTime: 2021-11-28 12:11:10
'''

import re
import time
import demjson
from lxml import etree
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import to_dict
from sqlmodule import User
from sqlmodule import Video
from sqlmodule import Data
from sqlmodule import Message

from lib.data import global_data
from lib.request import req, ReqException


def video_job():
    global_data.logger.info("==> video check...")

    logger_msg(msg_type="system", username="schedule", action="start video check")
    _video = VideoClass()
    _video.check()

def logger_msg(username="-", action="-", data="-", msg_type=""):
    """
    记录后天计划任务、程序错误信息等
    """
    with session_maker(global_data.sqlite_uri) as db_session:
        msg = Message(username=username, action=action, data=data, status="unread", msg_type=msg_type)
        db_session.add(msg)

class VideoClass(object):
    def __init__(self):
        super(VideoClass, self).__init__()

    def get_all_video_msg(self):
        """
        获取未完结的站点的信息
        """
        video_msg_list = list()
        try:
            with session_maker(global_data.sqlite_uri) as db_session:
                results = db_session.query(Video).filter(Video.status != "已完结").all()
                video_msg_list = to_dict(results)
        except Exception as e:
            global_data.logger.error(str(e))
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
            with session_maker(global_data.sqlite_uri) as db_session:
                results = db_session.query(Data).filter_by(ref_id=site_id).filter_by(category='video').all()
                video_msg_list = to_dict(results)
        except Exception as e:
            global_data.logger.error(str(e))
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
        video_type = ""

        if link.startswith("https://space.bilibili.com/"):
            href_text_list, status, video_type = self.bilibili_up(username, name, link)
        elif link.startswith("https://www.bilibili.com/bangumi"):
            href_text_list, status, video_type = self.bilibili_bangumi(username, name, link)
        elif link.startswith("https://www.acfun.cn/u"):
            href_text_list, status, video_type = self.acfun_up(username, name, link)
        elif link.startswith("https://www.acfun.cn/bangumi"):
            href_text_list, status, video_type = self.acfun_bangumi(username, name, link)
        elif link.startswith("https://www.agefans.vip/"):
            href_text_list, status, video_type = self.age(username, name, link)
        elif link.startswith("http://www.yhdm2.com/"):
            href_text_list, status, video_type = self.yhdm2(username, name, link)

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

            with session_maker(global_data.sqlite_uri) as db_session:
                result = db_session.query(User.username).filter_by(user_id=user_id).first()
                username = result.username

            # global_data.logger.info("Video check: {}".format(link))
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
                
                affect_num = 0
                try:
                    with session_maker(global_data.sqlite_uri) as db_session:
                        db_session.add_all(article_data)
                        if curr_status != status:
                            affect_num = db_session.query(Video).filter_by(id=site_id).update({"status": curr_status})
                            if affect_num:
                                global_data.logger.info("check {} success".format(name))
                except Exception as e:
                    global_data.logger.error(str(e))
                else:
                    if affect_num:
                        global_data.logger.info("Video check: {} 获取 {} 个新资源".format(name, len(article_data)))

    def bilibili_up(self, username, name, url):
        href_text_list = list()
        status = ""
        vedio_type = "bilibili_up"

        # up主视频接口
        up_video_api = "https://api.bilibili.com/x/space/arc/search?mid={}&ps=30&tid=0&pn=1&keyword=&order=pubdate&jsonp=jsonp"

        # 获取视频列表
        up_id = url.strip("/").split("/")[-1]
        try:
            resp = req.get(up_video_api.format(up_id))
            result = resp.json()
        except ReqException as e:
            global_data.logger.error(e)
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

        return href_text_list, status, vedio_type

    def bilibili_bangumi(self, username, name, url):
        href_text_list = list()
        status = ""
        vedio_type = "bilibili_bangumi"

        #番剧接口
        bangumi_api = "https://api.bilibili.com/pgc/view/web/season?season_id={}"

        # 获取视频列表
        try:
            resp = req.get(url)
            m = re.search(r'season_id":(\d+),', resp.text, re.M|re.S)
            season_id = m.group(1)
        except ReqException as e:
            global_data.logger.error(e)
        except Exception as e:
            # 番剧下架的话会获取不到 season_id
            # 出现 AttributeError: 'NoneType' object has no attribute 'group'
            if "NoneType' object has no attribute 'group" in str(e):
                msg = "未能成功获取 season_id，可能下架了..."
                global_data.logger.error(name+msg)
                logger_msg(msg_type="system", username="schedule", action="video check: {}".format(name), data=msg)
                logger_msg(msg_type="user", username=username, action="video check: {}".format(name), data=msg)
        else:
            resp = req.get(bangumi_api.format(season_id))
            result = resp.json()
            if result.get("code") == 0:
                result = result.get("result")
                _title = result.get("title")

                # 获取正片
                episodes = result.get("episodes")
                for episode in episodes:
                    badge = episode.get("badge") # 值为 预告 或为空
                    if "预告" not in badge:
                        title = episode.get("share_copy").replace("《" + _title + "》", "").strip()
                        url = episode.get("share_url")
                        href_text_list.append((url, title))

                status = result.get("publish").get("is_finish")
                status = "已完结" if status else "连载中"

        return href_text_list, status, vedio_type

    def acfun_up(self, username, name, url):
        href_text_list = list()
        status = ""
        vedio_type = "acfun_up"

        try:
            html = req.get(url).text
        except ReqException as e:
            global_data.logger.error(e)
        except Exception as e:
            global_data.logger.error(str(e))
        else:
            a_href_list = re.findall(r'<a href="(.*?)" target="_blank" class="ac-space-video', html)
            movurl_list = ["https://www.acfun.cn"+a_href for a_href in a_href_list]
            title_list = re.findall(r'<figcaption>.*?<p class="title line" title=".*?">(.*?)</p>', html, re.S|re.M)

            for index,title in enumerate(title_list):
                url = movurl_list[index]
                href_text_list.append((url, title))

            href_text_list = href_text_list[::-1]
            status = "连载中"

        return href_text_list, status, vedio_type

    def acfun_bangumi(self, username, name, url):
        href_text_list = list()
        status = ""
        vedio_type = "acfun_bangumi"

        try:
            html = req.get(url).text

            status = re.search(r'extendsStatus":"(.*?)",', html, re.S|re.M).group(1)
            bangumiList = re.search(r'window\.bangumiList = (.*?);', html).group(1)
        except ReqException as e:
            global_data.logger.error(e)
        except Exception as e:
            global_data.logger.error(str(e))
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

        return href_text_list, status, vedio_type

    def age(self, username, name, url):
        href_text_list = list()
        status = "连载中"
        vedio_type = "age"

        try:
            resp = req.get(url)
            resp.encoding = resp.apparent_encoding
            html = resp.text
        except ReqExceptin as error_msg:
            global_data.logger.error(error_msg)
            logger_msg(msg_type="system", username="schedule", action="video check: {}".format(name), data=str(error_msg))
            logger_msg(msg_type="user", username=username, action="{} 更新".format(name), data=str(error_msg))
        else:
            selector = etree.HTML(html)
            li_tag_list = selector.xpath('//*/div[@style="display:block"]/ul/li')
            for li_tag in li_tag_list:
                url = "https://www.agefans.vip"+li_tag.xpath('a/@href')[0]
                title = li_tag.xpath('a/text()')[0]
                if "pv" in title.lower() or "无字" in title:
                    continue
                href_text_list.append((url, title))

            detail_imform_kv_list = selector.xpath('//*/li[@class="detail_imform_kv"]')
            status_kv = detail_imform_kv_list[7]
            status = status_kv.xpath('span[2]/text()')[0]
            if status == "完结":
                status = "已完结"
            elif status == "连载":
                status = "连载中"

        return href_text_list, status, vedio_type

    def yhdm2(self, username, name, url):
        href_text_list = list()
        status = "连载中"
        vedio_type = "yhdm2"

        try:
            resp = req.get(url)
            resp.encoding = resp.apparent_encoding
            html = resp.text
        except ReqExceptin as error_msg:
            global_data.logger.error(error_msg)
            logger_msg(msg_type="system", username="schedule", action="video check: {}".format(name), data=str(error_msg))
            logger_msg(msg_type="user", username=username, action="{} 更新".format(name), data=str(error_msg))
        else:
            selector = etree.HTML(html)
            li_tag_list = selector.xpath('//*[@id="stab_1_71"]/ul/li')
            for li_tag in li_tag_list:
                url = "http://www.yhdm2.com"+li_tag.xpath('a/@href')[0]
                title = li_tag.xpath('a/text()')[0]
                if "pv" in title.lower() or "备用" in title or "英字" in title or "生肉" in title:
                    continue
                href_text_list.append((url, title))

            if status != "已完结":
                href_text_list = href_text_list[::-1]

        return href_text_list, status, vedio_type
