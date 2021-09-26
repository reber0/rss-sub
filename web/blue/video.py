#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-05 16:44:38
LastEditTime: 2021-09-26 18:00:37
'''

import re
import demjson
from lxml import etree
from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response
from sqlalchemy import func

from sqlmodule import session_maker
from sqlmodule import Config
from sqlmodule import Video
from sqlmodule import Data
from sqlmodule import to_dict
from sqlmodule import utc2bj

from lib.data import global_data
from lib.request import req

from .common import is_login
from .common import get_user_id
from .common import get_user_role
from .common import logger_user_action


video_blueprint = Blueprint('video', __name__)


@video_blueprint.route("/list", methods=['POST'])
@is_login
def list_site():
    data = request.get_json()
    pageIndex = data.get("page", 0)
    pageSize = data.get("limit", 0)
    keyword = data.get("keyword", "")
    target_id_list = data.get("target_id_list", "")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    data = list()
    count = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            if target_id_list:
                count = db_session.query(func.count(Video.id)).filter(
                    (Video.user_id==user_id) | (user_role=="root")).filter(Video.id.in_(target_id_list)).scalar()
                results = db_session.query(Video).filter_by(user_id=user_id).filter(Video.id.in_(target_id_list)).limit(pageSize).offset((pageIndex-1)*pageSize).all()
                data = to_dict(results)
            else:
                count = db_session.query(func.count(Video.id)).filter(
                    (Video.user_id==user_id) | (user_role=="root")).filter(Video.name.like("%{}%".format(keyword))).scalar()
                results = db_session.query(Video).filter((Video.user_id==user_id) | (user_role=="root")).filter(
                    Video.name.like("%{}%".format(keyword))).order_by(Video.id.desc()).limit(pageSize).offset((pageIndex-1)*pageSize).all()
                data = to_dict(results)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1}
        return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 0, "data": data, "count": count}
        return make_response(jsonify(r_data), 200)

@video_blueprint.route("/add", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def add_site():
    data = request.get_json()
    link = data.get("link", "").strip()
    global_data.logger.info("{}".format(link))

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    name = get_name(link)

    r_data = dict()
    affect_num = 0
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by(key="website") .first()
            result = demjson.decode(result[0])
            rss_site = result.get("domain")

            site = Video(user_id=user_id, name=name, link=link, status="连载中")
            db_session.add(site)
            db_session.flush()
            db_session.refresh(site)
            ref_id = site.id

            rss = "{}/api/data/rss/{}/video/{}".format(rss_site.rstrip("/"), user_id, ref_id)
            affect_num = db_session.query(Video).filter_by(id=ref_id).update({"rss": rss})

    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "rss add error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": 'add success', "data": {"rss": rss}}
            return make_response(jsonify(r_data))
        else:
            r_data = {"code": 1, "data": "rss add error"}
            return make_response(jsonify(r_data), 500)

@video_blueprint.route("/update", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def update_site():
    data = request.get_json()
    _id = data.get("id", "")
    data.pop("id")
    data.pop("rss")
    data.pop("add_time")

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            affect_num = db_session.query(Video).filter(
                (Video.user_id==user_id) | (user_role=="root"), Video.id==_id).update(data, synchronize_session=False)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "update error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num:
            r_data = {"code": 0, "msg": "update success"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "Permission denied"}
            return make_response(jsonify(r_data), 403)

@video_blueprint.route("/delete", methods=['POST'])
@is_login
@logger_user_action(msg_type="user")
def delete_site():
    data = request.get_json()
    target_id = data.get("target_id", "")
    target_id_list = data.get("target_id_list", "")
    if target_id:
        target_id_list = [target_id]

    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)
    user_role = get_user_role(access_token)

    affect_num = len(target_id_list)
    try:
        with session_maker(global_data.sqlite_uri) as db_session:
            affect_num = db_session.query(Video).filter(
                (Video.user_id==user_id) | (user_role=="root")).filter(Video.id.in_(target_id_list)).delete(synchronize_session=False)
            if affect_num == len(target_id_list):
                count = db_session.query(func.count(Data.id)).filter_by(category="video").filter(Data.ref_id.in_(target_id_list)).scalar()
                affect_num = db_session.query(Data).filter_by(category="video").filter(Data.ref_id.in_(target_id_list)).delete(synchronize_session=False)
    except Exception as e:
        global_data.logger.error(str(e))
        r_data = {"code": 1, "msg": "delete error"}
        return make_response(jsonify(r_data), 500)
    else:
        if affect_num == count:
            r_data = {"code": 0, "msg": "delete success"}
            return make_response(jsonify(r_data), 200)
        else:
            r_data = {"code": 1, "msg": "delete error"}
            return make_response(jsonify(r_data), 500)

def get_name(url):
    """
    获取番剧的名字
    """
    # 用户信息接口
    user_info_api = "https://api.bilibili.com/x/space/acc/info?mid={}"
    #番剧接口
    bangumi_api = "https://api.bilibili.com/pgc/view/web/season?season_id={}"

    if "space.bilibili.com" in url:
        mid = url.strip("/").split("/")[-1]
        url = user_info_api.format(mid)
        resp = req.get(url=url)
        result = resp.json()
        if result.get("code") == 0:
            name = result.get("data").get("name")
            return name
    elif "www.bilibili.com/bangumi" in url:
        resp = req.get(url=url)
        m = re.search(r'season_id":(\d+),', resp.text, re.S|re.M)
        if m:
            season_id = m.group(1)
            resp = req.get(url=bangumi_api.format(season_id))
            result = resp.json()
            if result.get("code") == 0:
                result = result.get("result")
                name = result.get("season_title")
        return name
    elif "www.acfun.cn/bangumi" in url:
        html = req.get(url=url).text
        m = re.search(r'bangumiTitle":"(.*?)",', html, re.S|re.M)
        if m:
            name = m1.group(1)
            return name
    elif "www.acfun.cn/u" in url:
        html = req.get(url=url).text
        m = re.search(r'<span class="name" data-username=(.*?)>', html, re.S|re.M)
        if m:
            name = m.group(1)
            return name
    elif "www.yhdm.so" in url:
        resp = req.get(url=url)
        resp.encoding = resp.apparent_encoding
        html = resp.text

        name = re.search(r'<h1>(.*?)</h1>', html, re.S|re.M).group(1)

        return name
    elif "www.yhdm2.com" in url:
        resp = req.get(url=url)
        resp.encoding = resp.apparent_encoding
        html = resp.text

        selector = etree.HTML(html)
        name = selector.xpath('//*/dt[@class="name"]/text()')[0]

        return name
