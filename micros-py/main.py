import logging
import asyncio
from internal.server import grpc_server, serve
from configs.conf import new_config
from internal.service import services

#  python -m grpc_tools.protoc -I api/pipeline/v1 --python_out=api/pipeline/v1 --grpc_python_out=api/pipeline/v1 --pyi_out=api/pipeline/v1 api/pipeline/v1/pipeline.proto


def async_run():
    logging.basicConfig(level=logging.INFO)
    loop = asyncio.get_event_loop()

    c = new_config()
    s = grpc_server.new_grpc_server(c)

    try:
        loop.run_until_complete(serve.serve(s))
    finally:
        loop.run_until_complete(*serve._cleanup_coroutines)
        loop.close()


if __name__ == '__main__':
    async_run()
    # c = new_config()
    # s = grpc_server.new_grpc_server(c)
    # serve.serve(s)
    # services.serve()
    # import sys
    # print(sys.path)
    # sys.path.insert()

