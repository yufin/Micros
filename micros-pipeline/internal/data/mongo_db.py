
import pymongo
import motor.motor_asyncio
from typing import Optional
from logging import Logger


class MotorClient:
    def __init__(self, uri: str, logger: Logger):
        # self.__client = pymongo.MongoClient(uri)
        self._async_client = motor.motor_asyncio.AsyncIOMotorClient(uri)
        self.logger: Logger = logger

    @property
    def client(self):
        return self._async_client

    async def __aenter__(self):
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        self.client.close()

    async def _get_content_raw(self, content_id: int) -> Optional[dict]:
        db = self.client["rc"]
        return await db["raw_content"].find_one(
            {"content_id": content_id},
            sort=[("created_at", -1)]
        )

    async def get_content(self, content_id: int) -> Optional[dict]:
        db = self.client["rc"]
        doc = await db["origin_content"].find_one(
            {"content_id": str(content_id)},
            sort=[("created_at", -1)]
        )
        if isinstance(doc, dict):
            doc['version'] = 'V1.0'
        if doc is None:
            doc = await self._get_content_raw(content_id)
        return doc

    async def get_processed_i18n_content(self, content_id: int, lang: str) -> Optional[dict]:
        db = self.client["rc"]
        return await db['processed_content_i18n'].find_one(
            {
                'content_id': content_id,
                'lang': lang
            },
            sort=[("created_at", -1)]
        )

    async def save_data(self,db:str, table: str, data: dict):
        db = self.client[db]
        return await db[table].insert_one(data)

    async def get_customer_supplier_quality(self, content_id: int) -> Optional[dict]:
        db = self.client["dw2"]
        return await db["customer_supplier_quality"].find_one(
            {"content_id": content_id}
        )

if __name__ == '__main__':

    from loguru import logger
    mgo_client = MotorClient(logger=logger,uri="mongodb://root:dev-password@192.168.44.169:27020")
    # db = mgo_client.client["rc"]
    # doc = db["origin_content"].find_one({"content_id": "471122972527045782"})
    # print(type(doc))
    # print(doc)
    # d = mgo_client.get_content(485507537341394070)

    import asyncio
    loop = mgo_client.client.get_io_loop()

    d = loop.run_until_complete(mgo_client._get_content_raw(485507528617110678))

    # d = mgo_client.get_content(485511236482509974)
    print(d)
    # import pprint
    # pprint.pprint(d)

