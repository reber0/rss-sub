#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-30 18:17:41
LastEditTime: 2021-09-15 16:20:30
'''
import sys
import uuid
import random
import string
import hashlib
from loguru import logger

def get_uuid():
    _uid = uuid.uuid4()
    _uid = "".join(str(_uid).split("-"))
    return _uid

# 初始化日志
def initLogger(log_file, log_err_file):
    logger.remove()
    logger_format1 = "[<green>{time:YYYY-MM-DD HH:mm:ss}</green>] [<level>{level}</level>] <level>{message}</level>"
    logger_format2 = "[<green>{time:YYYY-MM-DD HH:mm:ss}</green>] [<level>{level}</level>] <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>"
    logger.add(sys.stdout, format=logger_format2, level="INFO")
    logger.add(log_file, format=logger_format2, level="INFO", rotation="10 MB", enqueue=True, encoding="utf-8", errors="ignore")
    # logger.add(paths.log_file_path, format=logger_format2, level="INFO", rotation="00:00", enqueue=True, encoding="utf-8", errors="ignore")
    logger.add(log_err_file, rotation="10 MB", level="ERROR", enqueue=True, encoding="utf-8", errors="ignore")
    return logger

def fileGetContents(file_name):
    data = ""
    f_obj = None
    try:
        f_obj = open(file_name, 'r', encoding='UTF-8')
        data = f_obj.read()
    except Exception as e:
        return False
    else:
        return data
    finally:
        if f_obj:
            f_obj.close()

def getMd5(content):
    """
    返回传入值的 md5
    """
    if isinstance(content, str):
        content = content.encode(encoding='utf-8')
    m2 = hashlib.md5()
    m2.update(content)
    md5_text = m2.hexdigest()
    return md5_text

def randomStr(_length=8):
    """
    默认返回长度为 8 的字符串
    """
    return ''.join(random.sample(string.ascii_letters, _length))

def randomInt(_length=8):
    """
    默认返回长度为 8 的数字串
    """
    return ''.join(random.sample(string.digits, _length))

def getHeaders(random_ua=False, random_xff=False):
    """
    随机生成 User-Agent
    """
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
    if random_ua:
        user_agent = random.choice(USER_AGENTS)
    else:
        user_agent = USER_AGENTS[0]

    if random_xff:
        tmp_num1 = random.randint(1, 254)
        tmp_num2 = random.randint(1, 254)
        tmp_num3 = random.randint(1, 254)
        tmp_num4 = random.randint(1, 254)
        x_forwarded_for = "{a}.{b}.{c}.{d}".format(a=tmp_num1, b=tmp_num2, c=tmp_num3, d=tmp_num4)
    else:
        x_forwarded_for = '8.8.8.8'

    return {
        "Content-Type": "application/x-www-form-urlencoded",
        "Accept": "application/json, text/html, text/plain, */*",
        'User-Agent': user_agent,
        'X_FORWARDED_FOR': x_forwarded_for
    }
