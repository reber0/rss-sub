#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-30 18:17:41
@LastEditTime : 2021-01-29 19:32:43
'''

import uuid
import random
import hashlib
from datetime import timezone
from datetime import timedelta

def get_uuid():
    _uid = uuid.uuid4()
    _uid = "".join(str(_uid).split("-"))
    return _uid

def get_md5(string):
    m = hashlib.md5()
    m.update(string.encode())
    data = m.hexdigest()
    return data

def utc2bj(utc_time):
    """
    UTC 时间转为北京时间
    """
    return utc_time.astimezone(timezone(timedelta(hours=-8))).strftime("%Y-%m-%d %H:%M:%S")
