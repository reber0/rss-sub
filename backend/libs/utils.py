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


def random_useragent(condition=False):
    USER_AGENTS = [
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
        "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
        "Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
        "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
        "Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
        "Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
        "Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
        "Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
        "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
        "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
        "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
    ]
    if condition:
        return random.choice(USER_AGENTS)
    else:
        return USER_AGENTS[0]

# 随机X-Forwarded-For，动态IP
def random_x_forwarded_for(condition=False):
    if condition:
        tmp_num1 = random.randint(1, 254)
        tmp_num2 = random.randint(1, 254)
        tmp_num3 = random.randint(1, 254)
        tmp_num4 = random.randint(1, 254)
        return "{a}.{b}.{c}.{d}".format(
            a=tmp_num1, b=tmp_num2, c=tmp_num3, d=tmp_num4)
    else:
        return '8.8.8.8'

# HTTP 头设置
def get_headers(random_ua=False, random_xff=False):
    return {
        "Accept": "application/json, text/plain, */*",
        'User-Agent': random_useragent(random_ua),
        'X_FORWARDED_FOR': random_x_forwarded_for(random_xff)
    }
