import asyncio
from internal.server import grpc_server
from logging import Logger


def run_server(logger: Logger,server: grpc_server.GrpcServer):
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(server.serve_async())
    finally:
        loop.run_until_complete(*server._cleanup_coroutines)
        loop.close()