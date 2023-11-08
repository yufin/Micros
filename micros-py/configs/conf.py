from pydantic import BaseModel
from pydantic.fields import Field
import pathlib
import yaml
from typing import Any, Optional

class GRPC(BaseModel):
    addr: str
    timeout: str


class Server(BaseModel):
    grpc: GRPC = Field("grpc")


class DwData(BaseModel):
    target: str


class MongoDb(BaseModel):
    uri: str


class Data(BaseModel):
    mongodb: MongoDb = Field(alias="mongodb")
    dwdata: DwData = Field(alias="dwdata")


class Bootstrap(BaseModel):
    server: Server = Field(alias="server")
    data: Data = Field(alias="data")

    def __init__(self, config_path: Optional[str] = None, **data: Any):
        if config_path is not None:
            with open(str(config_path), "r") as f:
                d = yaml.safe_load(f)
            super().__init__(**d)
        else:
            super().__init__(**data)

