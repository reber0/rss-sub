#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-27 01:01:53
@LastEditTime : 2021-01-31 15:26:38
'''

import re
from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import User
from sqlmodule import Message
from sqlmodule import to_dict

from libs.auth import is_login
from libs.auth import is_admin
from libs.auth import get_username
from libs.auth import get_user_id
from libs.auth import get_user_role
from libs.utils import utc2bj

from setting import rss_sqlite_uri
from setting import logger

message_blueprint = Blueprint('message', __name__)


@message_blueprint.route("/tabs", methods=['POST'])
@is_login
def get_msg_tabs():
    """
    获取当前用户可以看到的消息 tabs
    root(可以看到 schedule、user、new_data等，即可以看到所有新消息)
    user(可以看到 new_data，即新番剧、新文章相关的消息)
    """

    access_token = request.headers.get('access_token')
    user_role = get_user_role(access_token)

    if user_role == "root":
        r_data = {"code": 0, "data": {"tabs": ['system', 'user']}}
        return make_response(jsonify(r_data))
    else:
        r_data = {"code": 0, "data": {"tabs": ['user']}}
        return make_response(jsonify(r_data))

@message_blueprint.route("/newmsg", methods=['POST'])
@is_login
def is_has_new_msg():
    """
    获取未读消息条数
    """

    access_token = request.headers.get('access_token')
    user_role = get_user_role(access_token)
    username = get_username(access_token)

    data = list()
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            if user_role == "root":
                unread_count = db_session.query(func.count(Message.id)).filter(
                    Message.status=="unread").scalar()
            else:
                unread_count = db_session.query(func.count(Message.id)).filter(
                    Message.username==username, Message.msg_type=="user",
                    (~Message.action.contains("/api/")) | (user_role=="root"), Message.status=="unread").scalar()
    except Exception as e:
        logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": {"newmsg": unread_count}}
        return make_response(jsonify(r_data))

@message_blueprint.route("/status/update", methods=['POST'])
@is_login
def update_msg_status():
    """
    更新消息状态：已读/未读
    """
    data = request.get_json()
    update_id = data.get("update_id", "") # 单个更新的 id
    status = data.get("status", "read")

    msg_type = data.get("msg_type", "") # 全部更新的 msg_type
    update_all = data.get("update_all", "")
    update_id_list = data.get("update_id_list", "") # 更新选中时的 id 列表

    if update_id:
        update_id_list = [update_id]

    access_token = request.headers.get('access_token')
    user_role = get_user_role(access_token)
    username = get_username(access_token)

    affect_num = 0
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            if update_id_list:
                affect_num = db_session.query(Message).filter(
                    (Message.username == username) | (user_role=="root"), Message.msg_type==msg_type, Message.id.in_(update_id_list)).update({"status": status}, synchronize_session=False)
            if update_all:
                affect_num = db_session.query(Message).filter(
                    (Message.username == username) | (user_role=="root"), Message.msg_type==msg_type).update({"status": status}, synchronize_session=False)

            # 获取剩余未读条数
            user_unread_count = db_session.query(func.count(Message.id)).filter(
                (Message.username == username) | (user_role=="root"), 
                (~Message.action.contains("/api/")) | (user_role=="root"), Message.status=="unread").scalar()
            schedule_unread_count = db_session.query(func.count(Message.id)).filter(
                user_role=="root", Message.status=="unread").scalar()

            unread_count = user_unread_count + schedule_unread_count
    except Exception as e:
        logger.error(str(e))
        r_data = {"code": 1, "msg": "操作失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": status, "unread_count": unread_count}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "操作失败"}
            return make_response(jsonify(r_data), 500)

@message_blueprint.route("/user/list", methods=['POST'])
@is_login
def get_user_action_msg():
    """
    获取用户未读消息
    """
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)

    access_token = request.headers.get('access_token')
    user_role = get_user_role(access_token)
    username = get_username(access_token)

    data = list()
    count = 0
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            unread_count = db_session.query(func.count(Message.id)).filter(
                (Message.username == username) | (user_role=="root"), Message.msg_type == "user", 
                (~Message.action.contains("/api/")) | (user_role=="root"), Message.status=="unread").scalar()
            count = db_session.query(func.count(Message.id)).filter(
                (Message.username == username) | (user_role=="root"), Message.msg_type == "user",
                (~Message.action.contains("/api/")) | (user_role=="root")).scalar()
            results = db_session.query(Message).filter(
                (Message.username == username) | (user_role=="root"), Message.msg_type == "user",
                (~Message.action.contains("/api/")) | (user_role=="root")).order_by(Message.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()
            data = to_dict(results)
        for index,value in enumerate(data):
            curr_data = data[index]["data"]
            curr_data = re.sub(r'\'password\': \'.*?\'', "'password': '******'", curr_data)
            curr_data = re.sub(r'\'email\': \'.*?\'', "'email': '******'", curr_data)
            data[index]["data"] = curr_data
            data[index]["add_time"] = utc2bj(data[index]["add_time"])
    except Exception as e:
        logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": data, "count": count, "unread_count": unread_count}
        return make_response(jsonify(r_data), 200)

@message_blueprint.route("/system/list", methods=['POST'])
@is_login
@is_admin
def get_system_msg():
    """
    获取系统未读消息
    """
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    data = list()
    count = 0
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            unread_count = db_session.query(func.count(Message.id)).filter(
                Message.msg_type == "system", Message.status=="unread").scalar()
            count = db_session.query(func.count(Message.id)).filter(
                Message.msg_type == "system").scalar()
            results = db_session.query(Message).filter_by(
                msg_type="system").order_by(Message.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()
            data = to_dict(results)
        for index,value in enumerate(data):
            data[index]["add_time"] = utc2bj(data[index]["add_time"])
    except Exception as e:
        logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": data, "count": count, "unread_count": unread_count}
        return make_response(jsonify(r_data), 200)

@message_blueprint.route("/delete", methods=['POST'])
@is_login
def delete_msg():
    data = request.get_json()
    delete_id = data.get("delete_id", "") # 单个删除的 id

    msg_type = data.get("msg_type", "") # 全部删除的 msg_type
    delete_all = data.get("delete_all", "")
    delete_id_list = data.get("delete_id_list", "") # 删除选中时的 id 列表

    if delete_id:
        delete_id_list = [delete_id]

    access_token = request.headers.get('access_token')
    user_role = get_user_role(access_token)
    username = get_username(access_token)

    affect_num = 0
    try:
        with session_maker(rss_sqlite_uri) as db_session:
            if delete_id_list:
                affect_num = db_session.query(Message).filter(
                    (Message.username == username) | (user_role=="root"), Message.id.in_(delete_id_list)).delete(synchronize_session=False)
            if delete_all:
                affect_num = db_session.query(Message).filter(
                    (Message.username == username) | (user_role=="root"), Message.msg_type==msg_type).delete(synchronize_session=False)

            # 获取剩余未读条数
            user_unread_count = db_session.query(func.count(Message.id)).filter(
                (Message.username == username) | (user_role=="root"), Message.status=="unread").scalar()
            schedule_unread_count = db_session.query(func.count(Message.id)).filter(
                user_role=="root", Message.status=="unread").scalar()

            if user_role == "root":
                unread_count = user_unread_count + schedule_unread_count
            else:
                unread_count = user_unread_count
    except Exception as e:
        logger.error(str(e))
        r_data = {"code": 1, "msg": "delete error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "delete success", "unread_count": unread_count}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "delete error"}
            return make_response(jsonify(r_data), 500)

