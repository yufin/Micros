from dependency_injector import containers, providers

import pathlib

config_path = pathlib.Path(__file__).parents[2].joinpath("configs").joinpath("config.yml")


class Container(containers.DeclarativeContainer):
    config = providers.Configuration(yaml_files=[str(config_path)])




if __name__ == '__main__':
    # print(config_path)
    config_path = pathlib.Path(__file__).parents[2].joinpath("configs").joinpath("config.yml")
    print(str(config_path))