#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-28 22:51:00
@LastEditTime : 2021-01-31 05:18:35
'''

from functools import wraps
from flask import request
from flask import jsonify
from flask import make_response

from sqlmodule import session_maker
from sqlmodule import User
from sqlmodule import Message

from setting import rss_sqlite_uri
from setting import logger

from libs.auth import get_user_id


def logger_user_action(msg_type="user"):
    """
    记录用户动作：访问的 URI、POST 的数据
    """
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            data = request.get_json()
            action = request.path
            if ("/api/user/update" in action) or ("/api/user/add" in action):
                data["password"] = "*"*8
                data = "POST: {}".format(data)
            elif data:
                data = "POST: {}".format(data)
            else:
                data = "-"

            access_token = request.headers.get('access_token')
            user_id = get_user_id(access_token)

            with session_maker(rss_sqlite_uri) as db_session:
                result = db_session.query(User.username).filter_by(user_id=user_id).first()
                username = result.username

                msg = Message(username=username, action=action, data=data, status="unread", msg_type=msg_type)
                db_session.add(msg)
                db_session.flush()
                db_session.refresh(msg)
                msg_id = msg.id

            if msg_id:
                return func(*args, **kwargs)
            else:
                # r_data = {"code": 1, "msg": "记录动作出错"}
                r_data = {"code": 1, "msg": "error"}
                return make_response(jsonify(r_data), 500)
        return wrapper
    return decorator


def logger_msg(username="-", action="-", data="-", msg_type=""):
    """
    记录后天计划任务、程序错误信息等
    """
    with session_maker(rss_sqlite_uri) as db_session:
        msg = Message(username=username, action=action, data=data, status="unread", msg_type=msg_type)
        db_session.add(msg)
