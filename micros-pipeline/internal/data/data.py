from internal.data.mongo_db import MotorClient
from internal.data.dw_data_client import DwDataClient
from internal.data.dw_data_client_v3 import DwDataClientV3
from internal.data.sql_db import SqlDb
from internal.data.json_translator import Translator
from logging import Logger


class DataRepo:
    def __init__(self,
                 logger: Logger,
                 mongo_db: MotorClient,
                 dw_data: DwDataClient,
                 dw_data_v3: DwDataClientV3,
                 sql_db: SqlDb,
                 translator: Translator
                 ):
        self.logger: Logger = logger
        self.mongo_db: MotorClient = mongo_db
        self.dw_data: DwDataClient = dw_data
        self.dw_data_v3: DwDataClientV3 = dw_data_v3
        self.sql_db: SqlDb = sql_db
        self.translator: Translator = translator
