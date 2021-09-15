#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-17 17:03:50
LastEditTime: 2021-09-15 16:43:41
'''

from datetime import timedelta

from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response
from flask import render_template
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import Config
from sqlmodule import Article
from sqlmodule import Video
from sqlmodule import Data
from sqlmodule import to_dict
from sqlmodule import utc2bj

from lib.data import global_data

from .common import is_login
from .common import get_user_id
from .common import get_user_role
from .common import logger_user_action


data_blueprint = Blueprint('data', __name__)


@data_blueprint.route("/article/list", methods=['POST'])
@is_login
def list_article_data():
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)
    name = data.get("name", "")
    status = data.get("status", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    r_data = dict()
    data = list()
    try:
        count = 0
        results = list()
        with session_maker(global_data.sqlite_uri) as db_session:
            count = db_session.query(func.count(Data.id)).join(Article, Data.ref_id==Article.id).filter(
                (Article.user_id==user_id) | (user_role=="root"), Data.category=='article',
                Data.status.like("{}%".format(status)), Article.name.like("%{}%".format(name))).scalar()
            results = db_session.query(Data.id, Article.name, Data.title, Data.url, Data.date, Data.status, Data.add_time).join(Article, Data.ref_id==Article.id).filter(
                (Article.user_id==user_id) | (user_role=="root"), Data.category=='article',
                Data.status.like("{}%".format(status)), Article.name.like("%{}%".format(name))).order_by(Data.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()

        for x in results:
            _tmp = dict()
            _tmp["id"] = x[0]
            _tmp["name"] = x[1]
            _tmp["title"] = x[2]
            _tmp["url"] = x[3]
            _tmp["date"] = utc2bj(x[4])
            _tmp["status"] = x[5]
            _tmp["add_time"] = utc2bj(x[6])
            data.append(_tmp)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": data, "count": count}
        return make_response(jsonify(r_data), 200)

@data_blueprint.route("/video/list", methods=['POST'])
@is_login
def list_video_data():
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)
    name = data.get("name", "")
    status = data.get("status", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    r_data = dict()
    data = list()
    try:
        count = 0
        results = list()
        with session_maker(global_data.sqlite_uri) as db_session:
            count = db_session.query(func.count(Data.id)).join(Video, Data.ref_id==Video.id).filter(
                (Video.user_id==user_id) | (user_role=="root"), Data.category=='video',
                Data.status.like("{}%".format(status)), Video.name.like("%{}%".format(name))).scalar()
            results = db_session.query(Data.id, Video.name, Data.title, Data.url, Data.date, Data.status, Data.add_time).join(Video, Data.ref_id==Video.id).filter(
                (Video.user_id==user_id) | (user_role=="root"), Data.category=='video',
                Data.status.like("{}%".format(status)), Video.name.like("%{}%".format(name))).limit(pageSize).offset((pageIndex-1)*pageSize).all()

        for x in results:
            _tmp = dict()
            _tmp["id"] = x[0]
            _tmp["name"] = x[1]
            _tmp["title"] = x[2]
            _tmp["url"] = x[3]
            _tmp["date"] = utc2bj(x[4])
            _tmp["status"] = x[5]
            _tmp["add_time"] = utc2bj(x[6])
            data.append(_tmp)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": data, "count": count}
        return make_response(jsonify(r_data), 200)

@data_blueprint.route("/article/delete", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def delete_article_data():
    data = request.get_json()
    data_id = data.get("id", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    affect_num = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            results = db_session.query(Data).join(Article, Data.ref_id==Article.id).filter(
                (Article.user_id==user_id) | (user_role=="root"),
                Data.category=='article', Data.id==data_id).first()
            if len(to_dict(results)) == 1:
                affect_num = db_session.query(Data).filter_by(category="article", id=data_id).delete()
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "delete error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "delete success"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "delete error"}
            return make_response(jsonify(r_data), 500)

@data_blueprint.route("/video/delete", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def delete_video_data():
    data = request.get_json()
    data_id = data.get("id", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            results = db_session.query(Data).join(Video, Data.ref_id==Video.id).filter(
                (Video.user_id==user_id) | (user_role=="root"),
                Data.category=='video', Data.id==data_id).first()
            if len(to_dict(results)) == 1:
                affect_num = db_session.query(Data).filter_by(category="video", id=data_id).delete()
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "delete error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "delete success"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "delete error"}
            return make_response(jsonify(r_data), 500)

# 改过变 article 的状态：已读/未读
@data_blueprint.route("/article/status/update", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def update_article_data():
    data = request.get_json()
    data_id = data.get("id", "")
    status = data.get("status", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    affect_num = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            results = db_session.query(Data).join(Article, Data.ref_id==Article.id).filter(
                (Article.user_id==user_id) | (user_role=="root"), Data.id==data_id).first()
            if len(to_dict(results)) == 1:
                affect_num = db_session.query(Data).filter_by(
                    category="article", id=data_id).update({"status": status})
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "操作失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": status}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "操作失败"}
            return make_response(jsonify(r_data), 500)

# 改变 video 的状态：已读/未读
@data_blueprint.route("/video/status/update", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def update_video_data():
    post_data = request.get_json()
    update_id = post_data.get("update_id", "")
    update_id_list = post_data.get("update_id_list", "") # 更新选中时的 id 列表
    status = post_data.get("status", "read")

    if update_id:
        update_id_list = [update_id]

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    affect_num = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            results = db_session.query(Data).join(Video, Data.ref_id==Video.id).filter(
                (Video.user_id==user_id) | (user_role=="root"), Data.id.in_(update_id_list)).first()
            if len(to_dict(results)):
                affect_num = db_session.query(Data).filter_by(
                    category="video").filter(
                        Data.id.in_(update_id_list)).update({"status": status}, synchronize_session=False)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "操作失败"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": status}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "操作失败"}
            return make_response(jsonify(r_data), 500)

@data_blueprint.route("/rss/<string:user_id>/<string:category>/<string:ref_id>", methods=['GET'])
def rss_link(user_id, category, ref_id):
    r_data = dict()
    site_msg = list()
    datas = list()
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            base_xml = db_session.query(Config.value).filter_by(key="base_xml").one()
            base_item = db_session.query(Config.value).filter_by(key="base_item").one()
            base_xml = base_xml[0]
            base_item = base_item[0]

            if category == "article":
                site_msg = db_session.query(Article.name, Article.link).filter(Article.user_id==user_id, Article.id==ref_id).first()
                datas = db_session.query(Data).join(Article, Data.ref_id==Article.id).filter(
                    Article.user_id==user_id, Data.category=='article', Data.ref_id==ref_id).order_by(Data.id.desc()).limit(30).all()
                datas = to_dict(datas)
            elif category == "video":
                site_msg = db_session.query(Video.name, Video.link).filter(Video.user_id==user_id, Video.id==ref_id).first()
                datas = db_session.query(Data).join(Video, Data.ref_id==Video.id).filter(
                    Video.user_id==user_id, Data.category=='video', Data.ref_id==ref_id).order_by(Data.id.desc()).limit(30).all()
                datas = to_dict(datas)
    except Exception as e:
        r_data["code"] = 1
        # r_data["msg"] = str(e).split('\n')[0]
        global_data.logger.error(str(e))
        return make_response(jsonify(r_data), 400)
    else:
        items = ""
        for index,data in enumerate(datas):
            title = data.get("title")
            url = data.get("url")
            date = data.get("date")-timedelta(seconds=index)
            description = data.get("description")
            items += base_item.format(title=title, url=url, date=date, description=description)

        rss_xml = base_xml.format(name=site_msg[0], link=site_msg[1], items=items)

        response = make_response(rss_xml, 200)
        response.headers["Content-Type"] = 'text/xml; charset=utf-8'
        return response
