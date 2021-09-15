#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-06 09:40:37
LastEditTime: 2021-09-15 17:49:39
'''

import re
import time
import requests
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import to_dict
from sqlmodule import User
from sqlmodule import Article
from sqlmodule import Data
from sqlmodule import Message

from lib.data import global_data
from lib.request import req, ReqException


def article_job():
    global_data.logger.info("==> article check...")

    logger_msg(msg_type="system", action="start article check", username="schedule")
    _article = ArticleClass()
    _article.check()

def logger_msg(username="-", action="-", data="-", msg_type=""):
    """
    记录后天计划任务、程序错误信息等
    """
    with session_maker(global_data.sqlite_uri) as db_session:
        msg = Message(username=username, action=action, data=data, status="unread", msg_type=msg_type)
        db_session.add(msg)

class ArticleClass(object):
    def __init__(self):
        super(ArticleClass, self).__init__()

    def get_all_site_msg(self):
        """
        获取所有站点的信息
        """
        site_msg_list = list()
        try:
            with session_maker(global_data.sqlite_uri) as db_session:
                results = db_session.query(Article).all()
                site_msg_list = to_dict(results)
        except Exception as e:
            global_data.logger.error(str(e))
        else:
            return site_msg_list

    def get_site_article_url(self, site_id):
        """
        获取一个博客的现有文章链接
        """
        site_article_list = list()
        article_url_list = list()
        try:
            with session_maker(global_data.sqlite_uri) as db_session:
                results = db_session.query(Data).filter_by(ref_id=site_id).filter_by(category='article').all()
                site_article_list = to_dict(results)
        except Exception as e:
            global_data.logger.error(str(e))
        else:
            if site_article_list:
                article_url_list = [article_msg["url"] for article_msg in site_article_list]
        return article_url_list

    def get_site_new_article_msg(self, username, name, link, regex, article_url_list):
        """
        获取一个站点的新文章信息
        """
        error_msg = ""
        new_article_msg_list = list()
        try:
            resp = req.get(url=link)
            html = resp.content
        except ReqException as error_msg:
            global_data.logger.error(error_msg)
            logger_msg(msg_type="system", username="schedule", action="article check: {}".format(name), data=str(error_msg))
            logger_msg(msg_type="user", username=username, action="{} 更新".format(name), data=str(error_msg))
        else:
            href_text_list = re.findall(regex, str(html, encoding='utf-8'), re.S|re.M)

            for href_text in href_text_list:
                a_href = href_text[0].strip()
                if not a_href.startswith("http"):
                    a_href = link.rstrip("/")+"/"+a_href.lstrip("/").lstrip(link)
                a_text = href_text[1].strip()

                if a_href not in article_url_list:
                    # resp = requests.get(url=a_href, verify=False)
                    # html = str(resp.content, encoding='utf-8')
                    # html = "<iframe src={}></iframe>".format(a_href)
                    html = "<h2>双指向左滑动</h2>"
                    new_article_msg_list.insert(0, (a_text, a_href, html))

        return new_article_msg_list

    def check(self):
        all_site_msg = self.get_all_site_msg()
        for site_msg in all_site_msg:
            site_id = site_msg.get("id")
            user_id = site_msg.get("user_id")
            name = site_msg.get("name")
            link = site_msg.get("link")
            regex = site_msg.get("regex")

            with session_maker(global_data.sqlite_uri) as db_session:
                result = db_session.query(User.username).filter_by(user_id=user_id).first()
                username = result.username

            # global_data.logger.info("Article check: {}".format(link))
            article_url_list = self.get_site_article_url(site_id)
            new_article_msg_list = self.get_site_new_article_msg(username, name, link, regex, article_url_list)

            article_data = list()
            for new_article_msg in new_article_msg_list[:30]: # 获取最新的 30 篇
                title, url, html = new_article_msg
                date = func.now()

                article = Data(ref_id=site_id, category="article", title=title, url=url, date=date, description=html, status="unread")
                article_data.append(article)

            if article_data:
                message_data_list = [x[0] for x in new_article_msg_list] # 获取 title 的 list 用于保存
                logger_msg(msg_type="system", username="schedule", action="article check: {}".format(name), data=str(message_data_list))
                logger_msg(msg_type="user", username=username, action="{} 已更新".format(name), data=str(message_data_list))
                try:
                    with session_maker(global_data.sqlite_uri) as db_session:
                        db_session.add_all(article_data)
                except Exception as e:
                    global_data.logger.error(str(e))
                else:
                    global_data.logger.info("Article check: {} 获取 {} 个新资源".format(name, len(article_data)))

