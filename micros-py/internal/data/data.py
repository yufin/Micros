from internal.data.mongo_db import MongoDb
from internal.data.dw_data_client import DwDataClient


class DataRepo:
    def __init__(self,
                 mongo_db: MongoDb,
                 dw_data: DwDataClient
                 ):
        self.mongo_db: MongoDb = mongo_db
        self.dw_data: DwDataClient = dw_data
