
import pymongo


class MongoDb:
    def __init__(self, uri: str):
        self.__client = pymongo.MongoClient(uri)

    @property
    def client(self):
        return self.__client

    def __del__(self):
        self.__client.close()

    def get_content(self, content_id: int) -> dict:
        db = self.client["rc"]
        doc = db["raw_content"].find_one(
            {"content_id": content_id},
            sort=[("created_at", pymongo.DESCENDING)]
        )
        return doc


if __name__ == '__main__':
    mgo_client = MongoDb("mongodb://root:dev-password@192.168.44.169:27020")
    # db = mgo_client.client["rc"]
    # doc = db["origin_content"].find_one({"content_id": "471122972527045782"})
    # print(type(doc))
    # print(doc)
    # d = mgo_client.get_content(485507537341394070)
    d = mgo_client.get_content(485511236482509974)
    print(d)
    # import pprint
    # pprint.pprint(d)

