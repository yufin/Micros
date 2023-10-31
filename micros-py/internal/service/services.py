from concurrent import futures
import logging

import grpc
from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2

# from ../..api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2


class PipelineServicer(pipeline_pb2_grpc.PipelineServiceServicer):

    async def GetContentProcess(self, request, context) -> pipeline_pb2.GetContentProcessResp:
        return pipeline_pb2.GetContentProcessResp()


# def serve():
#     port = '50051'
#     server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
#     pipeline_pb2_grpc.add_PipelineServiceServicer_to_server(PipelineServicer(), server)
#     server.add_insecure_port('[::]:' + port)
#     server.start()
#     print("Server started, listening on " + port)
#     server.wait_for_termination()


