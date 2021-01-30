#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-15 17:24:24
@LastEditTime : 2021-01-30 23:47:56
'''

import sys
sys.dont_write_bytecode = True  # 不生成pyc文件

from flask import Flask
from flask import request
from flask import jsonify
from flask import make_response
from flask import render_template
from flask import send_from_directory
from werkzeug.exceptions import HTTPException

from blue.set import set_blueprint
from blue.user import user_blueprint
from blue.article import article_blueprint
from blue.video import video_blueprint
from blue.data import data_blueprint
from blue.message import message_blueprint
from blue.upload import upload_blueprint

from setting import screct_key
from setting import logger


app = Flask(__name__, template_folder='../frontend/start')
app.secret_key = screct_key
app.register_blueprint(blueprint=set_blueprint, url_prefix='/api/set')
app.register_blueprint(blueprint=user_blueprint, url_prefix='/api/user')
app.register_blueprint(blueprint=article_blueprint, url_prefix='/api/article')
app.register_blueprint(blueprint=video_blueprint, url_prefix='/api/video')
app.register_blueprint(blueprint=data_blueprint, url_prefix='/api/data')
app.register_blueprint(blueprint=message_blueprint, url_prefix='/api/message')
app.register_blueprint(blueprint=upload_blueprint, url_prefix='/api/upload')


@app.route("/", methods=['GET'])
def index():
    return render_template('index.html')

# 请求 /face/xxx.png 可直接定位到 static/userface/xxx.png
@app.route('/face/<path:filename>')
def route_face_path(filename):
    return send_from_directory(app.root_path + '/static/userface/', filename)

# 请求 /src/index.js 可直接定位到 frontend/src/index.js
@app.route('/src/<path:filename>')
def route_src_path(filename):
    return send_from_directory(app.root_path + '/../frontend/src/', filename)
# 请求 /layui/xxx 可直接定位到 frontend/start/layui/xxx
@app.route('/layui/<path:filename>')
def route_layui_path(filename):
    return send_from_directory(app.root_path + '/../frontend/start/layui/', filename)

# # 请求 /json/xxx 可直接定位到 frontend/start/json/xxx
# @app.route('/json/<path:filename>')
# def route_json_path(filename):
#     return send_from_directory(app.root_path + '/../frontend/start/json/', filename, conditional=True)

@app.errorhandler(HTTPException)
def handle_http_error(e):
    data = {'status': 'error', 'description': e.description}
    response = make_response(jsonify(data), e.code)
    return response

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8083, use_reloader=True, debug=True)
