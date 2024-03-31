# -*- coding:utf-8 -*-
import json
from varname import nameof
import ast
import aiohttp
import asyncio
import urllib
import hashlib
import sys
import time
import random
from tqdm import tqdm
from typing import Any, Coroutine, Iterable, List, Tuple, Optional
from pprint import pprint
from logging import Logger
from internal.data.mongo_db import MotorClient
import nest_asyncio
nest_asyncio.apply()


class Translator:
    def __init__(self,logger: Logger, app_id: str, secret_key: str, url: str):
        self.logger: Logger = logger
        self.app_id: str = app_id
        self.secret_key: str = secret_key
        self.url: str = url
        self.salt: int = random.randint(32768, 65536)
        self.sleep_interval_s: float = 0.099
        self.show_pbar: bool = False
        self.pbar_kws = {}

    def _baidu_translation_url_builder(self, source_text: str, to_lang: str = "en", domain: str = 'finance', from_lang: str = 'auto'):
        sign = self.app_id + source_text + str(self.salt) + domain + self.secret_key
        sign = hashlib.md5(sign.encode()).hexdigest()
        request_url = self.url + '?appid=' + self.app_id + '&q=' + urllib.parse.quote(
            source_text) + '&from=' + from_lang + '&to=' + to_lang + '&salt=' + str(
            self.salt) + '&domain=' + domain + '&sign=' + sign
        return request_url

    async def _fetch(self, session, url: str, interval: float):
        await asyncio.sleep(interval)
        async with session.get(url) as response:
            return await response.json(encoding='utf-8')

    async def _async_request(self, session, urls: list) -> List[Any]:
        async def tup(idx: int, _task: Coroutine) -> Tuple[int, Any]:
            return idx, await _task

        async with session as session:
            tasks = [self._fetch(session, url, i*self.sleep_interval_s) for i, url in enumerate(urls)]

            wrapped_tasks = [tup(idx, t) for idx, t in enumerate(tasks)]
            if self.show_pbar:
                pbar = tqdm(asyncio.as_completed(wrapped_tasks), total=len(wrapped_tasks), **self.pbar_kws)
                results = [await f for f in pbar]
            else:
                results = [await f for f in asyncio.as_completed(wrapped_tasks)]
            return [r[1] for r in sorted(results, key=lambda x: x[0])]

    async def _run_loop(self, urls: list, limit: int = 1):
        connector = aiohttp.TCPConnector(limit=limit, force_close=True)
        async with aiohttp.ClientSession(connector=connector) as session:
            data = await self._async_request(session, urls)
            return data

    def _translate_process_async_baidu(self, source_text_list: list, limit: int = 4) -> list:
        urls = []
        for source_text in source_text_list:
            urls.append(self._baidu_translation_url_builder(source_text))

        new_loop = asyncio.new_event_loop()
        asyncio.set_event_loop(new_loop)
        loop = asyncio.get_event_loop()
        results = loop.run_until_complete(self._run_loop(urls, limit))
        return results

    def trans_json(self, report_data: dict) -> (Optional[str], list):
        self.logger.info("start translating json data with size: {}bytes".format(sys.getsizeof(report_data)))
        """
        get field which need to translate from json data,
        and translate by translation api, then recover and return json data
        :param report_data:
        :return: translated Json data
        """

        # variables required for processing
        to_translate_fields = []
        translated_fields = []
        holder_prefix = '__TRANSLATE__'
        translation_done = False

        def zh_field_detect(field: str) -> bool:
            """
            detect if the field need to translate
            :param field: field to detect
            :return: whether field contains chinese
            """
            if any('\u4e00' <= c <= '\u9fff' for c in field):
                return True
            else:
                return False

        def get_ground_fields(data):
            """
            Nested traversal of json data, get fields which need to translate,
            append field to to_translate_fields
            :param data: input data, dict, list or str, int, float
            :return: input data with fields covered by translate holder
            """
            if type(data) is dict:      # dict
                for k, v in data.items():   # traverse dict
                    data[k] = get_ground_fields(v)  # recursive call
                return data    # return dict
            elif type(data) is list:
                for idx, item in enumerate(data):
                    data[idx] = get_ground_fields(item)
                return data
            else:
                if translation_done:
                    if type(data) is str and data.__contains__(holder_prefix):
                        return eval(data.replace(holder_prefix, ''))
                    else:
                        return data
                else:
                    # in case iterable data in type of str,
                    # which will cause error when translate api called cause split list was too long
                    try:
                        temp_data = ast.literal_eval(data)
                        if type(temp_data) in [dict, list]:
                            return get_ground_fields(temp_data)
                    except Exception:
                        pass

                    if type(data) is str and zh_field_detect(data):
                        to_translate_fields.append(data)
                        return holder_prefix + nameof(translated_fields) + "[{}]".format(len(to_translate_fields) - 1)
                    else:
                        return data

        def translate_result_handler(result: dict) -> str or None:
            try:
                dst_str = result["trans_result"][0]["dst"]
                return dst_str
            except Exception as e:
                return None

        recovered_report_data = get_ground_fields(report_data)                      # jsonData替换为占位符的变量
        err_returned: list = []     # err returned by translation api
        results = self._translate_process_async_baidu(to_translate_fields, limit=5)                # 调用翻译接口
        for d in results:
            field = translate_result_handler(d)      # 处理翻译结果

            if field is None:
                err_returned.append(d)
                err_returned = [i for n, i in enumerate(err_returned) if i not in err_returned[n + 1:]]   # drop duplicate dict

            translated_fields.append(field)
        tqdm.write(f"err returned: {err_returned}")

        _retry_counter = 0
        while len(none_index := [i for i, x in enumerate(translated_fields) if x is None]) > 0:
            if _retry_counter > 6:
                return None
            _retry_counter += 1

            retrans_fields = [to_translate_fields[i] for i in none_index]
            retrans_results = self._translate_process_async_baidu(retrans_fields, limit=3)
            for i, d in zip(none_index, retrans_results):
                field = translate_result_handler(d)

                if field is None:
                    err_returned.append(d)
                    err_returned = [i for n, i in enumerate(err_returned) if i not in err_returned[n + 1:]]

                translated_fields[i] = field

            tqdm.write(f"err returned: {err_returned}\n")

        translation_done = True                                                      # 翻译完成
        translated_report_data = get_ground_fields(recovered_report_data)            # 将占位符替换为翻译结果
        return json.dumps(translated_report_data, ensure_ascii=False), err_returned                # 返回翻译后的json string

def process(content_id_dev: int, content_id_prod: int, file_path: str):
    import loguru
    t = Translator(
        logger=loguru.logger,
        app_id='20221108001442272',
        secret_key='Afkrj0NefMvyNNiTw8qi',
        url='https://fanyi-api.baidu.com/api/trans/vip/fieldtranslate'
    )
    t.show_pbar = True

    with open(file_path, 'r', encoding='utf-8') as f:
        report_data = json.load(f)

    res, errs = t.trans_json(report_data)

    # res(json string ) to dict
    res = json.loads(res)
    # pprint(res)


    data = {
        "content": res,
        "created_at": time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()),
        "content_id": content_id_dev,
        "updated_at": time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()),
        "lang": "en-US"
    }

    data_online = {
        "content": res,
        "created_at": time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()),
        "content_id": content_id_prod,
        "updated_at": time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()),
        "lang": "en-US"
    }

    if res is not None:
        mgo_client = MotorClient("mongodb://root:dev-password@192.168.44.169:27020", logger=loguru.logger)
        loop = mgo_client.client.get_io_loop()
        insert_res = loop.run_until_complete(mgo_client.save_data("rc", "processed_content_i18n", data))
        print(insert_res)

        mgo_client_online = MotorClient("mongodb://root:brillink8818mgo@10.0.232.121:27017", logger=loguru.logger)
        loop = mgo_client_online.client.get_io_loop()
        insert_res = loop.run_until_complete(
            mgo_client_online.save_data("rc", "processed_content_i18n", data_online))
        print(insert_res)


if __name__ == '__main__':

    process(
        content_id_dev=502435341588769942,
        content_id_prod=490699493877016598,
        file_path='../../static/502435341588769942.json'
    )