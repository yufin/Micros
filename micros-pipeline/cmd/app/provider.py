from dependency_injector import providers, containers
from loguru import logger
import logging
import argparse
import pathlib
from internal.server.grpc_server import GrpcServer
from internal.service.pipeline_service import PipelineService
from internal.conf.conf import Bootstrap
from internal.data.data import DataRepo
from internal.data.mongo_db import MotorClient
from internal.data.dw_data_client import DwDataClient
from internal.data.dw_data_client_v3 import DwDataClientV3
from internal.data.sql_db import SqlDb
from internal.data.json_translator import Translator
from internal.server.server import run_server
# logger.remove()
# logger.add(sys.stderr, level='WARNING')


parser = argparse.ArgumentParser(description='parse command line arguments for config the app')
parser.add_argument(
    '-conf',
    type=str,
    help='configuration file path',
    default=str(pathlib.Path(__file__).parents[2].joinpath("configs").joinpath("config.yml"))
)
args = parser.parse_args()
config_path = args.conf


class Configs(containers.DeclarativeContainer):
    config = providers.Singleton(
        Bootstrap,
        config_path=config_path
    )


class Loggers(containers.DeclarativeContainer):
    # logging.basicConfig(level=logging.INFO)
    default_logger = logger


class DataClients(containers.DeclarativeContainer):
    motor_client = providers.Factory(
        MotorClient,
        uri=Configs.config().data.mongodb.uri,
        logger=Loggers.default_logger
    )

    dw_data_client = providers.Factory(
        DwDataClient,
        logger=Loggers.default_logger,
        target=Configs.config().data.dwdata.target,
        certs_folder=Configs.config().TLS.folder_path,
        options=None
    )

    dw_data_client_v3 = providers.Factory(
        DwDataClientV3,
        logger=Loggers.default_logger,
        target=Configs.config().data.dwdata.target,
        certs_folder=Configs.config().TLS.folder_path,
        options=None
    )

    sql_db = providers.Factory(
        SqlDb,
        logger=Loggers.default_logger,
        db=Configs.config().data.sqldb.db,
        host=Configs.config().data.sqldb.host,
        password=Configs.config().data.sqldb.password,
        user=Configs.config().data.sqldb.user,
        port=Configs.config().data.sqldb.port
    )

    translator = providers.Factory(
        Translator,
        logger=Loggers.default_logger,
        app_id=Configs.config().data.baidu_translate_api.app_id,
        secret_key=Configs.config().data.baidu_translate_api.secret_key,
        url=Configs.config().data.baidu_translate_api.url
    )


class Repos(containers.DeclarativeContainer):

    data_repo = providers.Factory(
        DataRepo,
        mongo_db=DataClients.motor_client,
        dw_data=DataClients.dw_data_client,
        dw_data_v3=DataClients.dw_data_client_v3,
        sql_db=DataClients.sql_db,
        translator=DataClients.translator,
        logger=Loggers.default_logger
    )


class Services(containers.DeclarativeContainer):
    pipeline_service = providers.Factory(
        PipelineService,
        repo=Repos.data_repo,
        logger=Loggers.default_logger
    )


class Servers(containers.DeclarativeContainer):
    grpc_server = providers.Factory(
        GrpcServer,
        logger=Loggers.default_logger,
        addr=Configs.config().server.grpc.addr,
        certs_folder=Configs.config().TLS.folder_path,
        pipeline_service=Services.pipeline_service,
    )


class Runners(containers.DeclarativeContainer):
    server_runner = providers.Callable(
        run_server,
        logger=Loggers.default_logger,
        server=Servers.grpc_server
    )

