from dependency_injector import containers, providers
from internal.server import grpc_server

import pathlib

config_path = pathlib.Path(__file__).parents[2].joinpath("configs").joinpath("config.yml")


class Container(containers.DeclarativeContainer):
    config = providers.Configuration(yaml_files=[str("./config.yml")])

    # config.get("Server.GRPC.addr")
    print(config.get_json_files())
    # server = providers.Factory(
    #     grpc_server.GrpcServer,
    #     addr=config.Server.GRPC.addr,
    # )


if __name__ == '__main__':
    # print(config_path)
    # config_path = pathlib.Path(__file__).parents[2].joinpath("configs").joinpath("config.yml")
    # print(str(config_path))
    #
    container = Container()
    container.init_resources()
    # container.wire()