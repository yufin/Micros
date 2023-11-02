from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.data.data import DataRepo
from google.protobuf.struct_pb2 import Struct


class PipelineService(pipeline_pb2_grpc.PipelineServiceServicer):
    def __init__(self, repo: DataRepo):
        self.repo: DataRepo = repo

    async def GetContentProcess(self, request: pipeline_pb2.GetContentProcessReq, context) -> pipeline_pb2.GetContentProcessResp:
        doc = self.repo.mongo_db.get_content(request.content_id)



        st = Struct()
        st.update(doc['content'])
        return pipeline_pb2.GetContentProcessResp(success=True, msg="1", data=st)



