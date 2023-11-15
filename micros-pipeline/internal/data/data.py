from internal.data.mongo_db import MotorClient
from internal.data.dw_data_client import DwDataClient
from logging import Logger

class DataRepo:
    def __init__(self,
                 logger: Logger,
                 mongo_db: MotorClient,
                 dw_data: DwDataClient
                 ):
        self.logger: Logger = logger
        self.mongo_db: MotorClient = mongo_db
        self.dw_data: DwDataClient = dw_data
