from concurrent import futures
import logging
import grpc
from grpc_reflection.v1alpha import reflection
from configs.conf import Bootstrap
from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.service.services import PipelineServicer


def new_grpc_server(c: Bootstrap) -> grpc.Server:
    # server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    server = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))
    server.add_insecure_port(c.Server.GRPC.addr)
    pipeline_pb2_grpc.add_PipelineServiceServicer_to_server(PipelineServicer(), server)

    SERVICE_NAMES = (
        pipeline_pb2.DESCRIPTOR.services_by_name['PipelineService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    return server


