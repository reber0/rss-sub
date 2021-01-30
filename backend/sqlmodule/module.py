#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-15 17:24:24
@LastEditTime : 2021-01-31 04:27:05
'''

from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column
from sqlalchemy import Integer  # 普通整数，32位
from sqlalchemy import String  # 变长字符串
from sqlalchemy import Text  # 变长字符串，对较长字符串做了优化
from sqlalchemy import DateTime
from sqlalchemy.sql import func


# 创建基类
Base = declarative_base()


class User(Base):
    """定义数据模型"""
    __tablename__ = 'user'
    id = Column(Integer, primary_key=True, autoincrement=True)
    user_id = Column(String(32), nullable=False, comment="用户唯一 id(uuid)")
    username = Column(String(50), nullable=False, comment="用户名")
    password = Column(String(32), nullable=False, comment="密码")
    role = Column(String(10), nullable=False, comment="用户身份，root/user/e.g.")
    email = Column(String(100), nullable=False, comment="用户邮箱")
    face = Column(String(40), nullable=False, comment="头像的图片的id")
    add_time = Column(DateTime(timezone=True), server_default=func.now(), comment="添加时间")

class Article(Base):
    """
    文章相关的信息
    """
    __tablename__ = 'article'
    id = Column(Integer, primary_key=True, autoincrement=True)
    user_id = Column(String(32), comment="用户唯一识别 id")
    name = Column(String(100), comment="博客名字")
    link = Column(Text, comment="文章网站的网址")
    regex = Column(Text, comment="正则")
    rss = Column(String(100), comment="RSS 地址")
    add_time = Column(DateTime(timezone=True), server_default=func.now(), comment="添加时间")

class Video(Base):
    """
    番剧等其他的信息
    """
    __tablename__ = 'video'
    id = Column(Integer, primary_key=True, autoincrement=True)
    user_id = Column(String(32), comment="用户唯一识别 id")
    name = Column(String(100), comment="系列名字，比如番剧名")
    link = Column(Text, comment="主页目录，比如番剧主页、UP 主主页的 URL")
    status = Column(String(10), comment="连载状态")
    rss = Column(String(100), comment="RSS 地址")
    add_time = Column(DateTime(timezone=True), server_default=func.now(), comment="添加时间")

class Data(Base):
    """
    存放具体的文章和链接、番剧和链接等信息
    """
    __tablename__ = 'data'
    id = Column(Integer, primary_key=True, autoincrement=True)
    ref_id = Column(Integer, comment="article、video 表中站点对应的 id")
    category = Column(String(10), comment="表名，article、video，用于区分该条记录对应那个表里的数据")
    title = Column(Text, comment="标题，文章名字、番剧每集名字")
    date = Column(DateTime(timezone=True), comment="时间")
    url = Column(Text, comment="网址，文章链接、番剧每集 URL")
    description = Column(Text, comment="简介")
    status = Column(String(10), comment="状态，是否已读、已看")
    add_time = Column(DateTime(timezone=True), server_default=func.now(), comment="添加时间")

class Message(Base):
    """
    存放日志、消息等信息
    """
    __tablename__ = 'message'
    id = Column(Integer, primary_key=True, autoincrement=True)
    msg_type = Column(String(10), comment="消息类型，user/system")
    username = Column(String(50), nullable=False, comment="msg_type 为 user 时是用户名，为 system 时是 schedule/error")
    action = Column(Text, comment="执行的操作触发的 URI、计划任务动作")
    data = Column(Text, comment="POST 的数据、得到的数据")
    status = Column(String(10), comment="状态，是否已读、已看")
    add_time = Column(DateTime(timezone=True), server_default=func.now(), comment="添加时间")

class Config(Base):
    """
    配置信息
    """
    __tablename__ = 'config'
    id = Column(Integer, primary_key=True, autoincrement=True)
    key = Column(String(30), comment="配置的名字")
    value = Column(Text, comment="配置的值")


def init_db():
    from sqlalchemy import create_engine
    from sqlalchemy_utils import database_exists, create_database
    from sqlalchemy.orm import sessionmaker

    sqlite_uri = "sqlite:///../rss.sqlite3"
    engine = create_engine(sqlite_uri, echo=True)
    db_session = sessionmaker(bind=engine)

    if database_exists(engine.url):
        Base.metadata.drop_all(engine) # 删除所有表
        Base.metadata.create_all(engine) # 建表，如果存在同名表则跳过，不会覆盖
    else:
        create_database(engine.url)
        Base.metadata.create_all(engine)



if __name__ == "__main__":
    init_db()
