from pydantic import BaseModel


class GRPC(BaseModel):
    addr: str = '0.0.0.0:50051'
    timeout: int = 10


class Server(BaseModel):
    GRPC: GRPC


class MongoDb(BaseModel):
    uri: str = "mongodb://root:dev-password@192.168.44.169:27020"


class Data(BaseModel):
    MongoDb: MongoDb


class Bootstrap(BaseModel):
    Server: Server
    Data: Data


def new_config() -> Bootstrap:
    return Bootstrap(Server=Server(GRPC=GRPC()), Data=Data(MongoDb=MongoDb()))


