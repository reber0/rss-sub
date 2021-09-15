#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-31 15:17:06
LastEditTime: 2021-09-15 16:41:04
'''

import jwt
import time
import demjson
from functools import wraps

from flask import request
from flask import jsonify
from flask import make_response

from lib.data import global_data

from sqlmodule import session_maker
from sqlmodule import Config
from sqlmodule import User
from sqlmodule import Message


def logger_user_action(msg_type="user"):
    """
    记录用户动作：访问的 URI、POST 的数据
    """
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            data = request.get_json()
            action = request.path
            if data:
                data = "POST: {}".format(data)
            else:
                data = "-"

            access_token = request.headers.get('access_token')
            user_id = get_user_id(access_token)

            with session_maker(global_data.sqlite_uri) as db_session:
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
    with session_maker(global_data.sqlite_uri) as db_session:
        msg = Message(username=username, action=action, data=data, status="unread", msg_type=msg_type)
        db_session.add(msg)

def get_username(access_token=None, user_id=None):
    """
    通过 access_token 或 user_id 获取 username
    """

    if access_token:
        user_id = get_user_id(access_token)

    username = ""
    with session_maker(global_data.sqlite_uri) as db_session:
        result = db_session.query(User.username).filter_by(user_id=user_id).first()
        username = result.username
    return username

def get_user_role(access_token):
    """
    通过 access_token 获取用户身份
    """

    user_id = get_user_id(access_token)
    with session_maker(global_data.sqlite_uri) as db_session:
        role = db_session.query(User.role).filter_by(user_id=user_id).first()
        return role[0] if role else None

def get_user_id(access_token):
    """
    通过 access_token 获取用户 id
    """
    data = jwt_decode(access_token)
    return data.get("user_id")

def is_login(func):
    @wraps(func)
    def check_login(*args, **kwargs):
        access_token = request.headers.get('access_token')
        data = jwt_decode(access_token)
        if data:
            return func(*args, **kwargs)
        else:
            # 登录失效 return 401
            r_data = {"code": 401, "msg": "登录失效"}
            return make_response(jsonify(r_data), 401)
    return check_login

def is_admin(func):
    @wraps(func)
    def check_is_root(*args, **kwargs):
        access_token = request.headers.get('access_token')
        user_role = get_user_role(access_token)

        if user_role == "root":
            return func(*args, **kwargs)
        else:
            r_data = {"code": 403, "msg": "Permission denied"}
            return make_response(jsonify(r_data), 403)
    return check_is_root


def create_token(user_id):
    return jwt_encode(user_id)

def jwt_encode(user_id):
    with session_maker(global_data.sqlite_uri) as db_session:
        result = db_session.query(Config.value).filter_by(key="jwt_deadline").first()
        result = demjson.decode(result[0])
        hours = result.get("hours")
        minutes = result.get("minutes")

        payload = {
            "user_id": user_id,
            "iat": time.time(), # 该jwt的发布时间；unix 时间戳
            "exp": time.time()+60*minutes+60*60*hours # 该jwt销毁的时间；unix 时间戳，60分钟过期
        }
        access_token = jwt.encode(payload, key=global_data.jwt_key, algorithm='HS256')
        return access_token

def jwt_decode(access_token):
    data = dict()
    try:
        data = jwt.decode(access_token, key=global_data.jwt_key, algorithms=['HS256'])
    except jwt.exceptions.ExpiredSignatureError as e:
        global_data.logger.error("Signature has expired")
    except Exception as e:
        global_data.logger.error(str(e))
    else:
        return data
