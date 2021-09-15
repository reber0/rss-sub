#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
Author: reber
Mail: reber0ask@qq.com
Date: 2021-09-15 11:38:35
LastEditTime: 2021-09-15 11:38:37
'''

from sqlalchemy import create_engine
from sqlalchemy_utils import database_exists, create_database
from sqlalchemy.orm import sessionmaker

from .database import Base

def init_db(sqlite_uri):
    engine = create_engine(sqlite_uri, echo=False)
    db_session = sessionmaker(bind=engine)

    if database_exists(engine.url):
        # Base.metadata.drop_all(engine) # 删除所有表
        Base.metadata.create_all(engine) # 建表，如果存在同名表则跳过，不会覆盖
    else:
        create_database(engine.url)
        Base.metadata.create_all(engine)


if __name__ == "__main__":
    sqlite_uri = "sqlite:///../data/kone.db"
    init_db(sqlite_uri)
