from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.data.data import DataRepo
from google.protobuf.struct_pb2 import Struct
from internal.data.content_pipeline import ContentPipelineLatest
import http


class PipelineService(pipeline_pb2_grpc.PipelineServiceServicer):
    def __init__(self, repo: DataRepo):
        self.repo: DataRepo = repo

    async def GetContentProcess(
            self, request: pipeline_pb2.GetContentProcessReq, context) -> pipeline_pb2.GetContentProcessResp:

        doc = await self.repo.mongo_db.get_content(request.content_id)
        pipeline = ContentPipelineLatest(doc=doc, repo=self.repo)
        await pipeline.process()

        st = Struct()
        st.update(pipeline.content)
        return pipeline_pb2.GetContentProcessResp(success=True, code=http.HTTPStatus.OK, msg="", data=st)




