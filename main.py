#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
Author: reber
Mail: reber0ask@qq.com
Date: 2021-09-15 11:39:21
LastEditTime: 2021-09-15 21:20:17
'''
import os
import time
import json
import yaml
import pathlib
import argparse
from multiprocessing import Process
from multiprocessing import Manager
from lib.myqueue import MyQueue
from lib.utils import fileGetContents
from lib.utils import initLogger


class MainApp():
    def __init__(self):
        super(MainApp, self).__init__()

    def my_parser(self):
        '''使用说明'''
        service_type_list = [
            "web", "schedule", "all"
        ]
        example = """Examples:
                            \r  python3 {shell_name}
                            \r  python3 {shell_name} -s web
                            """

        parser = argparse.ArgumentParser(
            formatter_class = argparse.RawDescriptionHelpFormatter,#使 example 可以换行
            add_help = True,
            # description = "RssSub",
            )
        parser.epilog = example.format(shell_name=parser.prog)
        parser.add_argument("-s", dest="service_type", type=str, default="all",
                            choices=service_type_list, help="the type of service to start")

        # options, args = parser.parse_known_args()
        # parser.print_help()

        return parser

    def initApp(self):
        # 检查文件夹 log 和 data 是否存在，不存在则创建
        if not os.path.exists("log"):
            os.mkdir("log")
        if not os.path.exists("data"):
            os.mkdir("data")

        # 解析配置
        yaml_data = fileGetContents("config.yaml")
        yaml_config = yaml.safe_load(yaml_data)

        # 设置路径
        root_abspath = pathlib.Path(__file__).parent.resolve()
        yaml_config["root_abspath"] = root_abspath

        # 声明传入进程的变量，设置初始化数据为 self.config
        self.config = Manager().dict()
        self.config.update(yaml_config)

        # 初始化，初始化数据库
        if not os.path.exists(self.config["sqlite_uri"].split("///")[-1]):
            from sqlmodule import init_db
            init_db(self.config["sqlite_uri"])

        self.logger = initLogger(self.config["log_file"], self.config["log_err_file"])

    def startWeb(self):
        from web.start_web import startWeb
        web = Process(target=startWeb, args=(self.config,))
        web.daemon = False
        web.start()
        web_pid = web.pid

        time.sleep(1)

    def startSchedule(self):
        from schedule.start_schedule import startSchedule
        schedule = Process(target=startSchedule, args=(self.config,))
        schedule.daemon = False
        schedule.start()
        schedule_pid = schedule.pid

        time.sleep(1)

    def start(self, service_type):
        """
        开启
        """
        if service_type == "all":
            self.startWeb()
            self.startSchedule()
        elif service_type == "web":
            self.startWeb()
        elif service_type == "schedule":
            self.startSchedule()

    def run(self):
        parser = self.my_parser()
        options, args = parser.parse_known_args()

        self.initApp()

        if options.service_type:
            self.start(options.service_type)
        else:
            parser.print_help()


if __name__ == '__main__':
    app = MainApp()
    app.run()

