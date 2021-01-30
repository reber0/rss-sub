#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-15 17:24:24
@LastEditTime : 2021-01-11 16:15:13
'''

import sys
import pathlib
from loguru import logger

jwt_key = "7sMmCB6JYg59"
screct_key = "7sMmCB6JYg59"

abs_path = pathlib.Path(__file__).resolve().parent

rss_sqlite_uri = "sqlite:///{}/rss.sqlite3".format(abs_path)

# 设置日志路径
log_file_path = abs_path.joinpath("log/result.log")
err_log_file_path = abs_path.joinpath("log/error.log")

# 初始化日志
logger.remove()
logger_format1 = "[<green>{time:HH:mm:ss}</green>] <level>{message}</level>"
logger_format2 = "<green>{time:YYYY-MM-DD HH:mm:ss,SSS}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>"
logger.add(sys.stdout, format=logger_format2, level="INFO")
logger.add(log_file_path, format=logger_format2, level="INFO", rotation="10 MB", enqueue=True, encoding="utf-8", errors="ignore")
logger.add(err_log_file_path, rotation="10 MB", level="ERROR", enqueue=True, encoding="utf-8", errors="ignore")
