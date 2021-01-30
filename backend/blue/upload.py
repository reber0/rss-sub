#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2021-01-25 23:46:25
@LastEditTime : 2021-01-29 23:35:11
'''
import os
import demjson
from pathlib import Path
from flask import Blueprint
from flask import request
from flask import jsonify
from flask import make_response
from werkzeug.utils import secure_filename

from sqlmodule import session_maker
from sqlmodule import User
from sqlmodule import Config

from libs.auth import is_login
from libs.auth import get_user_id
from libs.utils import get_uuid

from setting import rss_sqlite_uri
from setting import logger

upload_blueprint = Blueprint('upload', __name__, static_folder="../static")


def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ['png', 'jpg', 'jpeg']


@upload_blueprint.route("/face", methods=['POST'])
@is_login
def upload_face():
    def save_tmp_file(file):
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file_ext = Path(filename).suffix
            _uid = get_uuid()
            file_path = "/tmp/"+_uid+file_ext
            file.save(file_path)

            file_size = os.stat(file_path).st_size # 单位是字节
            return file_path, file_size/1024 # 返回大小为 KB
        else:
            r_data = {"code": 1, "msg": "非法文件"}
            return make_response(jsonify(r_data), 400)

    def get_upload_max_size():
        upload_max_size = 0
        with session_maker(rss_sqlite_uri) as db_session:
            result = db_session.query(Config.value).filter_by().first()
            website = demjson.decode(result.website)
            upload_max_size = website.get("upload_max_size", "0")
        return upload_max_size

    file = request.files['file']
    access_token = request.headers.get('access_token')
    user_id = get_user_id(access_token)

    tmp_file_path, tmp_file_size = save_tmp_file(file)
    upload_max_size = get_upload_max_size()

    if tmp_file_size < upload_max_size:
        filename = tmp_file_path.split("/")[-1]
        try:
            with session_maker(rss_sqlite_uri) as db_session:
                affect_num = db_session.query(User).filter_by(user_id=user_id).update({"face": filename})
        except Exception as e:
            logger.error(e)
            r_data = {"code": 1, "msg": "更新头像失败"}
            return make_response(jsonify(r_data), 500)
        else:
            if affect_num:
                file_path = Path.cwd().joinpath('static', 'userface', filename)
                os.rename(tmp_file_path, file_path)

                r_data = {"code": 0, "msg": "更新头像成功", "data": {"filename": filename}}
                return make_response(jsonify(r_data))
            else:
                r_data = {"code": 1, "msg": "更新头像失败"}
                return make_response(jsonify(r_data), 500)
    else:
        r_data = {"code": 1, "msg": "文件过大"}
        return make_response(jsonify(r_data), 200) # 实际该返回 403.3，为了 layui 弹窗提示消息，设为200
