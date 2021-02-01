#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
@Author: reber
@Mail: reber0ask@qq.com
@Date: 2019-11-14 12:19:21
@LastEditTime : 2021-01-11 14:53:27
'''
from urllib.parse import quote
import urllib3
import requests
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


# HTTP 头设置
def get_headers(random_ua=False, random_xff=False):
    USER_AGENTS = [
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
        "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
        "Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
        "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
        "Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
        "Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
        "Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
        "Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
        "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
        "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
        "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
    ]
    if random_ua:
        user_agent = random.choice(USER_AGENTS)
    else:
        user_agent = USER_AGENTS[0]

    if random_xff:
        tmp_num1 = random.randint(1, 254)
        tmp_num2 = random.randint(1, 254)
        tmp_num3 = random.randint(1, 254)
        tmp_num4 = random.randint(1, 254)
        x_forwarded_for = "{a}.{b}.{c}.{d}".format(a=tmp_num1, b=tmp_num2, c=tmp_num3, d=tmp_num4)
    else:
        x_forwarded_for = '8.8.8.8'

    return {
        "Content-Type": "application/x-www-form-urlencoded",
        "Accept": "application/json, text/html, text/plain, */*",
        'User-Agent': user_agent,
        'X_FORWARDED_FOR': x_forwarded_for
    }


class ReqExceptin(Exception):
    """
    捕获错误
    """
    def __init__(self, error_msg):
        self.error_msg = error_msg

    def __str__(self):
        return self.error_msg


class MySession(requests.Session):
    """
    重写 requsts 的 request 方法，对超时时间等进行设置、处理错误信息
    """

    def request(self, method, url,
                params=None, data=None, headers=get_headers(), cookies=None, files=None,
                auth=None, timeout=20, allow_redirects=True, proxies=None,
                hooks=None, stream=None, verify=False, cert=None, json=None):

        # Create the Request.
        req = requests.Request(
            method=method.upper(),
            url=url,
            headers=headers,
            files=files,
            data=data or {},
            json=json,
            params=params or {},
            auth=auth,
            cookies=cookies,
            hooks=hooks,
        )
        prep = self.prepare_request(req)

        proxies = proxies or {}

        settings = self.merge_environment_settings(
            prep.url, proxies, stream, verify, cert
        )

        # Send the request.
        send_kwargs = {
            'timeout': timeout,
            'allow_redirects': allow_redirects,
        }
        send_kwargs.update(settings)

        try:
            resp = None
            error_msg = ""
            try:
                resp = self.send(prep, **send_kwargs)
            except requests.exceptions.ConnectTimeout:
                error_msg = "ConnectTimeout"
            except requests.exceptions.ReadTimeout:
                error_msg = "ReadTimeout"
            except requests.exceptions.Timeout:
                error_msg = "Timeout"
            except requests.exceptions.ProxyError:
                error_msg = "ProxyError"
            except requests.exceptions.SSLError:
                error_msg = "SSLError"
            except ConnectionResetError:
                error_msg = "ConnectionResetError"
            except requests.exceptions.ConnectionError:
                error_msg = "ConnectionError"
            except Exception as e:
                # raise e
                error_msg = str(e)
            else:
                return resp
            raise ReqExceptin("{} {}".format(url, error_msg))
        except KeyboardInterrupt as e:
            raise e


req = MySession()

if __name__ == "__main__":
    proxies = {"http": "http://127.0.0.1:8080",
               "https": "http://127.0.0.1:8080"}
    resp = req.get(url="https://google.com",
                        proxies=proxies, verify=False, timeout=3)
    print(resp.status_code)
