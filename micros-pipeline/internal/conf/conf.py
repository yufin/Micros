from pydantic import BaseModel
from pydantic.fields import Field
import pathlib
import yaml
from typing import Any, Optional


class Consul(BaseModel):
    addr: str
    port: int
    token: str


class Registry(BaseModel):
    consul: Consul = Field("consul")


class GRPC(BaseModel):
    addr: str
    timeout: str


class Server(BaseModel):
    grpc: GRPC = Field("grpc")
    registry: Registry = Field("registry")


class DwData(BaseModel):
    target: str


class MongoDb(BaseModel):
    uri: str


class SqlDb(BaseModel):
    host: str
    port: int
    user: str
    password: str
    db: str


class Tls(BaseModel):
    folder_path: str = Field(alias="folder_path")

class BaiduTranslateApi(BaseModel):
    app_id: str = Field(alias="app_id")
    secret_key: str = Field(alias="secret_key")
    url: str = Field(alias="url")

class Data(BaseModel):
    mongodb: MongoDb = Field(alias="mongodb")
    dwdata: DwData = Field(alias="dwdata")
    sqldb: SqlDb = Field(alias="sqldb")
    baidu_translate_api: BaiduTranslateApi = Field(alias="baidu_translate_api")


class Bootstrap(BaseModel):
    server: Server = Field(alias="server")
    data: Data = Field(alias="data")
    TLS: Tls = Field(alias="tls")

    def __init__(self, config_path: Optional[str] = None, **data: Any):
        if config_path is not None:
            with open(str(config_path), "r") as f:
                d = yaml.safe_load(f)
            super().__init__(**d)
        else:
            super().__init__(**data)
