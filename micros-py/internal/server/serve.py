import grpc
import logging

_cleanup_coroutines = []


async def serve(grpc_server: grpc.Server):
    await grpc_server.start()
    print("Server started")
    # grpc_server.wait_for_termination()

    async def server_graceful_shutdown():
        logging.info("Starting graceful shutdown...")
        await grpc_server.stop(5)

    _cleanup_coroutines.append(server_graceful_shutdown())
    await grpc_server.wait_for_termination()