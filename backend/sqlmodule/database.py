#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2020-12-15 17:24:24
@LastEditTime : 2021-01-31 04:04:00
'''

from contextlib import contextmanager
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm import scoped_session


def get_db_session(sql_uri, echo=False):
    """
    获取 db session
    """

    # sqlite3 创建 engine 时不能使用编码、连接池等参数
    engine = create_engine(
        sql_uri,
        echo=echo, # 不显示具体执行过程

        # client_encoding='utf8',
        # max_overflow=0,  # 超过连接池大小外最多创建的连接
        # pool_size=10, # 连接池大小
        # pool_timeout=30, # 获取连接的超时阈值
        # pool_recycle=3600, # 3600之后对线程池中的线程进行一次连接的回收（重置）
    )
    session_factory = sessionmaker(bind=engine)

    # session 不是线程安全的，并且我们一般 session 对象都是全局的
    # 当多个线程共享一个session时，数据处理就会发生错误
    # 用 scoped_session 使线程安全
    session = scoped_session(session_factory)
    # session = session_factory()

    return session


@contextmanager
def session_maker(sql_uri, echo=False):
    db_session = get_db_session(sql_uri, echo=echo)
    try:
        yield db_session
        db_session.commit()
    except:
        db_session.rollback()
        raise
    finally:
        db_session.close()


def to_dict(vendors):
    """
    将结果转为字典列表: [{}, {}]
    """
    def _tmp(vendor):
        return {c.name: getattr(vendor, c.name, None) for c in vendor.__table__.columns}

    if vendors:
        if isinstance(vendors, list):
            return [_tmp(vendor) for vendor in vendors]
        return [_tmp(vendors)]
    else:
        return []
