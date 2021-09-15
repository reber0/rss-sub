#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-27 19:01:23
LastEditTime: 2021-09-15 19:45:28
'''
import demjson
from apscheduler.schedulers.background import BlockingScheduler

from sqlmodule import session_maker
from sqlmodule import Config

from lib.data import global_data
from lib.utils import initLogger

from .article import article_job
from .video import video_job


def startSchedule(init_config=None, global_queue=None):
    global_data.update(init_config)
    global_data.logger = initLogger(global_data.log_file, global_data.log_err_file)

    global_data.logger.info("Start schedule...")

    with session_maker(global_data.sqlite_uri) as db_session:
        check_article_rate = db_session.query(Config.value).filter_by(key="check_article_rate").one()
        check_video_rate = db_session.query(Config.value).filter_by(key="check_video_rate").one()

    check_article_rate = demjson.decode(check_article_rate[0])
    check_video_rate = demjson.decode(check_video_rate[0])

    article_h = check_article_rate.get("hours")
    article_m = check_article_rate.get("minutes")
    article_s = check_article_rate.get("seconds")
    video_h = check_video_rate.get("hours")
    video_m = check_video_rate.get("minutes")
    video_s = check_article_rate.get("seconds")

    scheduler = BlockingScheduler()
    scheduler.add_job(
        func=article_job, max_instances=1, trigger="interval",
        hours=article_h, minutes=article_m, seconds=article_s)
    scheduler.add_job(
        func=video_job, max_instances=1, trigger="interval",
        hours=video_h, minutes=video_m, seconds=video_s)
    scheduler.start()
