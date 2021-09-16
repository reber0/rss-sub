#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-17 17:03:50
LastEditTime: 2021-09-16 11:21:51
'''

import re
import demjson
import requests
from urllib.parse import urlparse, urlunparse
from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import Config
from sqlmodule import Article
from sqlmodule import Data
from sqlmodule import to_dict
from sqlmodule import utc2bj

from lib.data import global_data
from lib.request import req

from .common import is_login
from .common import get_user_id
from .common import get_user_role
from .common import logger_user_action


article_blueprint = Blueprint('article', __name__)


@article_blueprint.route("/check", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def check_regex():
    data = request.get_json()
    name = data.get("name", "")
    link = data.get("link", "")
    regex = data.get("regex", "")

    o = urlparse(link)
    base_url = urlunparse((o.scheme, o.netloc, "/", "", "", ""))

    article_tag = list()
    try:
        resp = req.get(url=link)
        href_text_list = re.findall(regex, str(resp.content, encoding='utf-8'), re.S|re.M)
        if href_text_list:
            for href_text in href_text_list:
                a_href = href_text[0].strip()
                if not a_href.startswith("http"):
                    a_href = base_url+a_href.lstrip("/")
                a_text = href_text[1].strip()
                article_tag.append({"title": a_text, "url": a_href})
    except requests.exceptions.ConnectionError as e:
        r_data = {"code": 1, "msg": "requests.exceptions.ConnectionError"}
        return make_response(jsonify(r_data), 500)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "server error"}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "msg": "", "data": article_tag[:5], "count": len(article_tag)}
        return make_response(jsonify(r_data), 200)


@article_blueprint.route("/add", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def add_blog():
    post_data = request.get_json()
    name = post_data.get("name", "")
    link = post_data.get("link", "")
    regex = post_data.get("regex", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)

    r_data = dict()
    affect_num = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="website") .first()
            result = demjson.decode(result[0])
            rss_site = result.get("domain")

            site = Article(user_id=user_id, name=name, link=link, regex=regex)
            db_session.add(site)
            db_session.flush()
            db_session.refresh(site)
            ref_id = site.id

            rss = "{}/api/data/rss/{}/article/{}".format(rss_site.rstrip("/"), user_id, ref_id)
            affect_num = db_session.query(Article).filter_by(id=ref_id).update({"rss": rss})

    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "add site error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "得到 rss 链接: <br>{}".format(rss)}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "add site error"}
            return make_response(jsonify(r_data), 500)

@article_blueprint.route("/list", methods=['POST'])
@is_login
def list_blog():
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    r_data = dict()
    data = list()
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            count = db_session.query(func.count(Article.id)).filter(
                (Article.user_id==user_id) | (user_role=="root")).scalar()
            results = db_session.query(Article).filter(
                (Article.user_id==user_id) | (user_role=="root")).order_by(Article.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()
            _data = to_dict(results)
        if _data:
            for x in _data:
                x["add_time"] = utc2bj(x["add_time"])
                data.append(x)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "msg": "", "data": data, "count": count}
        return make_response(jsonify(r_data), 200)

@article_blueprint.route("/update", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def update_blog():
    data = request.get_json()
    _id = data.get("id", "")
    data.pop("id")
    data.pop("add_time")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    try:
        from sqlmodule import User
        from sqlalchemy import or_, and_
        with session_maker(global_data.sqlite_uri) as db_session:
            affect_num = db_session.query(Article).filter(
                (Article.user_id==user_id) | (user_role=="root")).filter_by(id=_id).update(data, synchronize_session=False)
    except Exception as e:
        # r_data["msg"] = str(e).split('\n')[0]
        global_data.logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "Permission denied"}
            return make_response(jsonify(r_data), 403)

@article_blueprint.route("/delete", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def delete_blog():
    data = request.get_json()
    _id = data.get("id")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    affect_num = 0
    r_data = dict()
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            affect_num = db_session.query(Article).filter(
                (Article.user_id==user_id) | (user_role=="root")).filter(Article.id==_id).delete(synchronize_session=False)
            if affect_num:
                count = db_session.query(func.count(Data.id)).filter(Data.category=="article").filter(Data.ref_id==_id).scalar()
                if count:
                    affect_num = db_session.query(Data).filter(Data.category=="article").filter(Data.ref_id==_id).delete()
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
            return make_response(jsonify(r_data), 403)

@article_blueprint.route("/search", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def search_article_site():
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)
    name = data.get("name", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    r_data = dict()
    data = list()
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            count = db_session.query(func.count(Article.id)).filter(
                (Article.user_id==user_id) | (user_role=="root"), Article.name.like("%{}%".format(name))).scalar()
            results = db_session.query(Article).filter(
                (Article.user_id==user_id) | (user_role=="root"), Article.name.like("%{}%".format(name))).order_by(Article.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()
            _data = to_dict(results)

        if _data:
            for x in _data:
                x["add_time"] = utc2bj(x["add_time"])
                data.append(x)
    except Exception as e:
        logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "msg": "", "data": data, "count": count}
        return make_response(jsonify(r_data), 200)
