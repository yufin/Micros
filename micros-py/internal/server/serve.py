import logging
import asyncio
from internal.server import grpc_server


def run_server(s: grpc_server.GrpcServer):
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(s.serve_async())
    finally:
        loop.run_until_complete(*s._cleanup_coroutines)
        loop.close()