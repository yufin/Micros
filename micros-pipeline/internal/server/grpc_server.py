import grpc
from grpc_reflection.v1alpha import reflection
from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.service.pipeline_service import PipelineService
from logging import Logger


class GrpcServer:
    def __init__(
            self,
            logger: Logger,
            addr: str,
            pipeline_service: PipelineService
    ):
        self.logger: Logger = logger
        self._cleanup_coroutines = []
        self._server: grpc.Server = grpc.aio.server()
        self._server.add_insecure_port(addr)
        pipeline_pb2_grpc.add_PipelineServiceServicer_to_server(pipeline_service, self._server)

        service_names = (
            pipeline_pb2.DESCRIPTOR.services_by_name['PipelineService'].full_name,
            reflection.SERVICE_NAME,
        )

        reflection.enable_server_reflection(service_names, self._server)

    async def serve_async(self):
        async def server_graceful_shutdown():
            self.logger.info("Starting graceful shutdown...")
            await self._server.stop(5)

        await self._server.start()
        self.logger.info("Server started")

        self._cleanup_coroutines.append(server_graceful_shutdown())
        await self._server.wait_for_termination()
        self.logger.info("Server stopped")




