#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-11 16:26:10
@LastEditTime : 2021-01-30 23:31:11
'''

import demjson
from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response

from setting import rss_sqlite_uri
from setting import logger

from sqlmodule import session_maker
from sqlmodule import Article
from sqlmodule import Video
from sqlmodule import Config
from sqlmodule import User
from sqlmodule import to_dict

from libs.auth import is_login
from libs.auth import is_admin
from libs.auth import get_user_id
from libs.auth import get_user_role

from libs.utils import get_md5


set_blueprint = Blueprint('set', __name__)


@set_blueprint.route("/copyright", methods=['POST'])
def get_copyright():
    """
    获取 Copyright
    """
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="website").first()
            website = demjson.decode(result.value)
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1,  "msg": "server error"}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": {"copyright": website.get("copyright")}}
        return make_response(jsonify(r_data))

@set_blueprint.route("/menu", methods=['POST'])
@is_login
def get_menu():
    """
    获取左侧目录
    """
    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    user_role = ""
    with session_maker(rss_sqlite_uri) as db_session:
        results = db_session.query(User.role).filter_by(user_id=user_id).first()
        user_role = results[0]

    menu = dict()
    if user_role == "root":
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="root_menu").first()
            menu = demjson.decode(result.value)
    elif user_role == "user":
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="user_menu").first()
            menu = demjson.decode(result.value)
    r_data = {"code": 0, "msg": "", "data": menu}
    return make_response(jsonify(r_data))

@set_blueprint.route("/user/info", methods=['POST'])
@is_login
def get_user_info():
    """
    获取用户信息
    """
    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)

    result = dict()
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            results = db_session.query(User).filter_by(user_id=user_id).first()
            results = to_dict(results)
            results = results[0]

        results.pop("password")
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1,  "msg": "获取用户信息失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if len(results):
            r_data = {"code": 0, "data": results}
            return make_response(jsonify(r_data))

@set_blueprint.route("/user/info/update", methods=['POST'])
@is_login
def update_user_info():
    """
    更新用户信息
    """
    data = request.get_json()
    data.pop("file")
    data.pop("add_time")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)

    try:
        with session_maker(rss_sqlite_uri) as db_session:
            affect_num = db_session.query(User).filter_by(user_id=user_id).update(data)
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1, "msg": "更新失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "更新成功"}
            return make_response(jsonify(r_data))
        else:
            r_data = {"code": 1, "msg": "更新失败"}
            return make_response(jsonify(r_data), 500)

@set_blueprint.route("/user/password/update", methods=['POST'])
@is_login
def update_user_password():
    """
    更新用户密码
    """
    data = request.get_json()
    old_password = data.get("oldPassword", "")
    old_password = get_md5(old_password)
    new_password = data.get("password", "")
    new_password = get_md5(new_password)

    if old_password and new_password:
        access_token = request.headers.get('access_token')
        user_id = get_user_id(access_token)

        affect_num = 0
        try:
            with session_maker(rss_sqlite_uri) as db_session:
                result = db_session.query(User.password).filter_by(user_id=user_id).first()
                if old_password == result.password:
                    affect_num = db_session.query(User).filter_by(user_id=user_id).update({"password": new_password})
                else:
                    affect_num = 2
        except Exception as e:
            logger.error(e)
            r_data = {"code": 1, "msg": "修改失败"}
            return make_response(jsonify(r_data), 500)
        else:
            if affect_num == 1:
                r_data = {"code": 0, "msg": "修改成功"}
                return make_response(jsonify(r_data))
            elif affect_num == 2:
                r_data = {"code": 1, "msg": "旧密码错误"}
                return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 1, "msg": "密码不能为空"}
        return make_response(jsonify(r_data), 400)

@set_blueprint.route("/system/website", methods=['POST'])
@is_login
@is_admin
def get_website_info():
    """
    获取网站信息：domain、title、keywords、copyright 等
    """
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="website").first()
            website = demjson.decode(result.value)
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1,  "msg": "server error"}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": website}
        return make_response(jsonify(r_data))

@set_blueprint.route("/system/website/update", methods=['POST'])
@is_login
@is_admin
def update_website_info():
    """
    更新网站信息
    """
    data = request.get_json()
    domain = data.get("domain", "")

    affect_num = 0
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="website").first()
            website = demjson.decode(result.value)
            old_domain = website.get("domain")

            # 更新
            affect_num = db_session.query(Config).filter_by(key="website").update({"value": str(data)})

            # 更新已有 rss 链接中的域名
            if affect_num and domain:
                article_list = db_session.query(Article).all()
                for article in article_list:
                    article.rss = article.rss.replace(old_domain, domain)

                video_list = db_session.query(Video).all()
                for video in video_list:
                    video.rss = video.rss.replace(old_domain, domain)
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1,  "msg": "保存失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "保存成功"}
            return make_response(jsonify(r_data))
        else:
            r_data = {"code": 1,  "msg": "保存失败"}
            return make_response(jsonify(r_data), 500)

@set_blueprint.route("/system/email", methods=['POST'])
@is_login
@is_admin
def get_email_info():
    """
    获取邮箱信息：smtp、port、nickname、email、email_pwd
    """
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="email").first()
            email = demjson.decode(result.value)
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1,  "msg": "server error"}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": email}
        return make_response(jsonify(r_data))

@set_blueprint.route("/system/email/update", methods=['POST'])
@is_login
@is_admin
def update_email_info():
    """
    更新邮箱信息
    """
    data = request.get_json()
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            affect_num = db_session.query(Config).filter_by(key="email").update({"value": str(data)})
    except Exception as e:
        logger.error(e)
        r_data = {"code": 1,  "msg": "保存失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "保存成功"}
            return make_response(jsonify(r_data))
        else:
            r_data = {"code": 1,  "msg": "保存失败"}
            return make_response(jsonify(r_data), 500)

