
import pymongo
import motor.motor_asyncio
from typing import Optional


class MotorClient:
    def __init__(self, uri: str):
        # self.__client = pymongo.MongoClient(uri)
        self._async_client = motor.motor_asyncio.AsyncIOMotorClient(uri)

    @property
    def client(self):
        return self._async_client

    async def __aenter__(self):
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        self.client.close()

    async def get_content(self, content_id: int) -> Optional[dict]:
        db = self.client["rc"]
        return await db["raw_content"].find_one(
            {"content_id": content_id},
            sort=[("created_at", -1)]
        )


if __name__ == '__main__':
    mgo_client = MotorClient("mongodb://root:dev-password@192.168.44.169:27020")
    # db = mgo_client.client["rc"]
    # doc = db["origin_content"].find_one({"content_id": "471122972527045782"})
    # print(type(doc))
    # print(doc)
    # d = mgo_client.get_content(485507537341394070)

    import asyncio
    loop = mgo_client.client.get_io_loop()

    d = loop.run_until_complete(mgo_client.get_content(485511236482509974))

    # d = mgo_client.get_content(485511236482509974)
    print(d)
    # import pprint
    # pprint.pprint(d)

