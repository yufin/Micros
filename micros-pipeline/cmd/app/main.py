import sys
import pathlib
proj_path = pathlib.Path(__file__).parents[2]
sys.path.insert(0, proj_path.__str__())

from provider import Runners


def main():
    Runners.server_runner()


if __name__ == '__main__':
    main()


"""
#  python -m grpc_tools.protoc -I api/pipeline/v1 --python_out=api/pipeline/v1 --grpc_python_out=api/pipeline/v1 --pyi_out=api/pipeline/v1 api/pipeline/v1/pipeline.proto
#  python -m grpc_tools.protoc -I api/dwdata/v2 --python_out=api/dwdata/v2 --grpc_python_out=api/dwdata/v2 --pyi_out=api/dwdata/v2 api/dwdata/v2/dw_data.proto
"""


# import logging
# import pathlib
# from internal.server.grpc_server import GrpcServer
# from internal.server.server import run_server
# from internal.service.pipeline_service import PipelineService
# from internal.conf.conf import Bootstrap
# from internal.data.data import DataRepo
# from internal.data.mongo_db import MotorClient
# from internal.data.dw_data_client import DwDataClient
# from loguru import logger
#
#
# def main():
#     # logging.basicConfig(level=logging.INFO)
#
#     cfg_path = pathlib.Path(__file__).parents[2].joinpath("configs").joinpath("config.yml")
#     config = Bootstrap(config_path=str(cfg_path))
#
#     mgo = MotorClient(
#         uri=config.data.mongodb.uri,
#         logger=logger
#     )
#
#     dw = DwDataClient(
#         target=config.data.dwdata.target,
#         options=None,
#         logger=logger
#     )
#
#     data_repo = DataRepo(
#         mongo_db=mgo,
#         dw_data=dw,
#         logger=logger,
#     )
#
#     pipeline_service = PipelineService(
#         repo=data_repo,
#         logger=logger,
#     )
#
#     s = GrpcServer(
#         addr=config.server.grpc.addr,
#         pipeline_service=pipeline_service,
#         logger=logger,
#     )
#
#     run_server(logger=logger, server=s)
#
#
# if __name__ == '__main__':
#     main()
#
#
#
#
# """
# #  python -m grpc_tools.protoc -I api/pipeline/v1 --python_out=api/pipeline/v1 --grpc_python_out=api/pipeline/v1 --pyi_out=api/pipeline/v1 api/pipeline/v1/pipeline.proto
# #  python -m grpc_tools.protoc -I api/dwdata/v2 --python_out=api/dwdata/v2 --grpc_python_out=api/dwdata/v2 --pyi_out=api/dwdata/v2 api/dwdata/v2/dw_data.proto
# """
#
