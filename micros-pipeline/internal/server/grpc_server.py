import grpc
from grpc_reflection.v1alpha import reflection
from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.service.pipeline_service import PipelineService
from logging import Logger


class GrpcServer:
    def __init__(self, logger: Logger, addr: str, pipeline_service: PipelineService, certs_folder: str):
        self.logger: Logger = logger
        self._addr = addr
        self._cleanup_coroutines = []
        self._server: grpc.Server = grpc.aio.server()

        pipeline_pb2_grpc.add_PipelineServiceServicer_to_server(pipeline_service, self._server)

        service_names = (
            pipeline_pb2.DESCRIPTOR.services_by_name['PipelineService'].full_name,
            reflection.SERVICE_NAME,
        )
        reflection.enable_server_reflection(service_names, self._server)

        with open(certs_folder+"/server_key.pem", "rb") as f:
            private_key = f.read()
        with open(certs_folder+"/server_cert.pem", "rb") as f:
            certificate_chain = f.read()
        with open(certs_folder+"/client_ca_cert.pem", "rb") as f:
            ca_certificates = f.read()

        server_credentials = grpc.ssl_server_credentials(
            [(private_key, certificate_chain)],
            root_certificates=ca_certificates,
            require_client_auth=True,
        )

        self._server.add_secure_port(self._addr, server_credentials)
        # self._server.add_insecure_port(self._addr)

    async def serve_async(self):
        async def server_graceful_shutdown():
            self.logger.info("Starting graceful shutdown...")
            await self._server.stop(5)

        await self._server.start()
        self.logger.info("Server is listening at port: %s" % self._addr.split(":")[1])

        self._cleanup_coroutines.append(server_graceful_shutdown())
        await self._server.wait_for_termination()
        self.logger.info("Server stopped")




