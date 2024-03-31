import sys
import pathlib
import argparse
proj_path = pathlib.Path(__file__).parents[2]
sys.path.insert(0, proj_path.__str__())
from provider import Runners


# parser = argparse.ArgumentParser(description='parse command line arguments for config the app')
# parser.add_argument('-conf', type=str, help='configuration file path', default="../../configs/config.yml")
# args = parser.parse_args()
# config_path = args.conf


def main():
    # print(config_path)
    Runners.server_runner()


if __name__ == '__main__':
    main()


"""
#  python -m grpc_tools.protoc -I api/pipeline/v1 --python_out=api/pipeline/v1 --grpc_python_out=api/pipeline/v1 --pyi_out=api/pipeline/v1 api/pipeline/v1/pipeline.proto
#  python -m grpc_tools.protoc -I api/dwdata/v2 --python_out=api/dwdata/v2 --grpc_python_out=api/dwdata/v2 --pyi_out=api/dwdata/v2 api/dwdata/v2/dw_data.proto
#  python -m grpc_tools.protoc -I api/dwdata/v3 --python_out=api/dwdata/v3 --grpc_python_out=api/dwdata/v3 --pyi_out=api/dwdata/v3 api/dwdata/v3/dw_data_v3.proto
"""
