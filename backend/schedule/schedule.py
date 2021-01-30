#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-15 17:24:24
@LastEditTime : 2021-01-30 22:51:30
'''

import demjson
from apscheduler.schedulers.background import BlockingScheduler

from schedule.article import article_check
from schedule.video import video_check

from sqlmodule import session_maker
from sqlmodule import Config

from setting import rss_sqlite_uri


def start_schedule():
    with session_maker(rss_sqlite_uri) as db_session:
        check_article_rate = db_session.query(Config.value).filter_by(key="check_article_rate").one()
        check_video_rate = db_session.query(Config.value).filter_by(key="check_video_rate").one()

    check_article_rate = demjson.decode(check_article_rate[0])
    check_video_rate = demjson.decode(check_video_rate[0])

    article_h = check_article_rate.get("hours")
    article_m = check_article_rate.get("minutes")
    video_h = check_video_rate.get("hours")
    video_m = check_video_rate.get("minutes")
    # print(article_h, article_m, video_h, video_m)

    scheduler = BlockingScheduler()
    scheduler.add_job(func=article_check, trigger="interval", hours=article_h, minutes=article_m)
    scheduler.add_job(func=video_check, trigger="interval", hours=video_h, minutes=video_m)
    scheduler.start()
