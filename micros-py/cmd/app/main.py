import logging
from internal.server.grpc_server import GrpcServer
from internal.server.serve import run_server
from internal.service.pipeline_service import PipelineService
from configs.conf import new_config
from internal.data.data import DataRepo
from internal.data.mongo_db import MongoDb
from internal.data.dw_data_client import DwDataClient


def main():
    logging.basicConfig(level=logging.INFO)
    config = new_config()

    mgo = MongoDb(
        uri=config.data.mongodb.uri
    )

    dw = DwDataClient(
        target=config.data.dwdata.target,
        options=None
    )

    data_repo = DataRepo(
        mongo_db=mgo,
        dw_data=dw
    )

    pipeline_service = PipelineService(
        repo=data_repo
    )

    s = GrpcServer(
        addr=config.server.grpc.addr,
        pipeline_service=pipeline_service,
    )

    run_server(s)


if __name__ == '__main__':
    main()




"""
#  python -m grpc_tools.protoc -I api/pipeline/v1 --python_out=api/pipeline/v1 --grpc_python_out=api/pipeline/v1 --pyi_out=api/pipeline/v1 api/pipeline/v1/pipeline.proto
#  python -m grpc_tools.protoc -I api/dwdata/v2 --python_out=api/dwdata/v2 --grpc_python_out=api/dwdata/v2 --pyi_out=api/dwdata/v2 api/dwdata/v2/dw_data.proto
"""


