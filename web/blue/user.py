#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-25 16:05:49
LastEditTime: 2021-09-15 18:01:45
'''
import uuid
import random
from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import User
from sqlmodule import Article
from sqlmodule import Video
from sqlmodule import Data
from sqlmodule import Message
from sqlmodule import to_dict
from sqlmodule import utc2bj

from lib.data import global_data
from lib.utils import getMd5

from .common import is_login
from .common import is_admin
from .common import get_user_id
from .common import get_username
from .common import create_token
from .common import logger_msg
from .common import logger_user_action



user_blueprint = Blueprint('user', __name__)


@user_blueprint.route("/login", methods=['POST'])
def login():
    """
    用户登录
    """
    post_data = request.get_json()
    username = post_data.get("username", "")
    password = post_data.get("password", "")
    post_password = getMd5(password)

    r_data = dict()
    with session_maker(global_data.sqlite_uri) as db_session:
        results = db_session.query(User.user_id, User.password).filter_by(username=username).first()

    if results:
        curr_user_id = results[0]
        curr_password = results[1]
        if curr_password == post_password:
            token = create_token(curr_user_id)
            if token:
                r_data = {"code": 0, "data": {"access_token": token}}
                logger_msg(username=username, action="/api/user/login", data="登录成功", msg_type="user")
                return make_response(jsonify(r_data), 200)
            else:
                r_data = {"code": 1, "msg": "服务器处理错误"}
                logger_msg(username=username, action="/api/user/login", data="登录验证成功，但获取 token 失败。", msg_type="user")
                return make_response(jsonify(r_data), 500)
        else:
            r_data = {"code": 1, "msg": "用户名或密码错误"}
            logger_msg(username=username, action="/api/user/login", data="登录失败 {}".format(post_data), msg_type="user")
            return make_response(jsonify(r_data), 401)
    else:
        r_data = {"code": 1, "msg": "用户名或密码错误"}
        logger_msg(username=username, action="/api/user/login", data="登录失败 {}".format(post_data), msg_type="user")
        return make_response(jsonify(r_data), 401)

@user_blueprint.route("/logout", methods=['POST'])
@logger_user_action(msg_type="user")
def logout():
    r_data = {"code": 0, "msg": "退出登录"}
    return make_response(jsonify(r_data), 200)

@user_blueprint.route("/face", methods=['POST'])
@is_login
def get_user_face():
    """
    获取头像名字
    """
    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)

    if user_id:
        with session_maker(global_data.sqlite_uri) as db_session:
            result = db_session.query(User.face).filter_by(user_id=user_id).first()
            face = result.face

            r_data = {"code": 0, "msg": "", "data": {"face": face}}
            return r_data
    else:
        # 登录状态失效，返回 401，layui 自动跳转到 login
        r_data = {"code": 401, "msg": "Unauthorized"}
        return make_response(jsonify(r_data))

@user_blueprint.route("/list", methods=['POST'])
@is_login
@is_admin
def list_user():
    """
    显示所有用户
    """
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)
    username = data.get("username", "")
    email = data.get("email", "")
    role = data.get("role", "")
    export_id_list = data.get("export_id_list", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)

    r_data = dict()
    data = list()
    count = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            if export_id_list:
                count = db_session.query(func.count(User.id)).filter(User.id.in_(export_id_list)).scalar()
                results = db_session.query(User).filter(User.id.in_(export_id_list)).limit(pageSize).offset((pageIndex-1)*pageSize).all()
                data = to_dict(results)
            else:
                # User.user_id != user_id 不显示当前用户
                count = db_session.query(func.count(User.id)).filter(
                    User.user_id != user_id, User.username.like("%{}%".format(username)),
                    User.email.like("%{}%".format(email)), User.role.like("%{}%".format(role))).scalar()
                results = db_session.query(User).filter(
                    User.user_id != user_id, User.username.like("%{}%".format(username)),
                    User.email.like("%{}%".format(email)), User.role.like("%{}%".format(role))).order_by(User.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()
                data = to_dict(results)

        for index,value in enumerate(data):
            data[index]["add_time"] = utc2bj(value["add_time"])
            data[index].pop("user_id")
            data[index].pop("password")

    except Exception as e:
        r_data["code"] = 1
        global_data.logger.error(str(e))
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": data, "count": count}
        return make_response(jsonify(r_data), 200)

@user_blueprint.route("/add", methods=['POST'])
@is_login
@is_admin
@logger_user_action(msg_type="user")
def add_user():
    """
    添加用户
    """
    data = request.get_json()
    username = data.get("username", 0)
    password = data.get("password", 0)
    role = data.get("role", "user")
    email = data.get("email", "")
    face = str(random.randint(1,9))+".png"

    password = getMd5(password)
    _uid = uuid.uuid4()
    user_id = "".join(str(_uid).split("-"))

    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            user = User(user_id=user_id, username=username, password=password, role=role, email=email, face=face)
            db_session.add(user)
            db_session.flush()
            db_session.refresh(user)
            new_id = user.id
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "添加出错"}
        return make_response(jsonify(r_data), 500)
    else:
        if new_id:
            r_data = {"code": 0, "msg": "添加成功"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "添加失败"}
            return make_response(jsonify(r_data), 500)

@user_blueprint.route("/update", methods=['POST'])
@is_login
@is_admin
@logger_user_action(msg_type="user")
def update_user():
    data = request.get_json()
    _id = data.get("id", "")
    password = data.get("password", "")

    if password:
        data["password"] = getMd5(password)
    else:
        data.pop("password")
    data.pop("id")
    data.pop("add_time")

    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            affect_num = db_session.query(User).filter_by(id=_id).update(data)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "更新出错"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "更新成功"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "更新失败"}
            return make_response(jsonify(r_data), 403)

@user_blueprint.route("/delete", methods=['POST'])
@is_login
@is_admin
@logger_user_action(msg_type="user")
def delete_user():
    data = request.get_json()
    delete_id = data.get("delete_id", "")
    delete_id_list = data.get("delete_id_list", "")
    if delete_id:
        delete_id_list = [delete_id]

    access_token = request.headers.get('access_token')
    curr_user_id = get_user_id(access_token)
    username = get_username(access_token)

    with session_maker(global_data.sqlite_uri) as db_session:
        result = db_session.query(User.id).filter_by(user_id=curr_user_id).first()
        _id = result.id
        if _id in delete_id_list:
            r_data = {"code": 1, "msg": "不能删除自己(id 为 {})".format(_id)}
            return make_response(jsonify(r_data), 500)

    affect_num = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            # 删除用户及用户相关数据：Article、Video、Data、Message
            for delete_id in delete_id_list:
                result = db_session.query(User.user_id).filter_by(id=delete_id).first()
                delete_user_id = result.user_id

                # 删除用户
                delete_user_affect_num = db_session.query(User).filter_by(user_id=delete_user_id).delete()
                global_data.logger.info("delete user {}".format(delete_user_id))

                # 删除用户 Article 相关数据
                results = db_session.query(Article.id).filter_by(user_id=delete_user_id).all()
                article_id_list = [result.id for result in results]
                affect_num = db_session.query(Article).filter_by(user_id=delete_user_id).delete(synchronize_session=False)
                global_data.logger.info("delete user {} article: {}".format(delete_user_id, article_id_list))

                # 删除用户 Video 相关数据
                results = db_session.query(Video.id).filter_by(user_id=delete_user_id).all()
                video_id_list = [result.id for result in results]
                affect_num = db_session.query(Video).filter_by(user_id=delete_user_id).delete(synchronize_session=False)
                global_data.logger.info("delete user {} video: {}".format(delete_user_id, video_id_list))

                # 删除用户 Data 相关数据
                affect_num = db_session.query(Data).filter(Data.category=="article", Data.ref_id.in_(article_id_list)).delete(synchronize_session=False)
                affect_num = db_session.query(Data).filter(Data.category=="video", Data.ref_id.in_(video_id_list)).delete(synchronize_session=False)
                global_data.logger.info("delete user {} data".format(delete_user_id))
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "删除失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if delete_user_affect_num:
            r_data = {"code": 0, "msg": "删除成功"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "删除失败"}
            return make_response(jsonify(r_data), 500)
