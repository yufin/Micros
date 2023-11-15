from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.data.data import DataRepo
from google.protobuf.struct_pb2 import Struct
from internal.biz.content_pipeline import ContentPipelineV3
import http
from logging import Logger
from loguru import logger


class PipelineService(pipeline_pb2_grpc.PipelineServiceServicer):
    def __init__(
            self,
            repo: DataRepo,
            logger: Logger
    ):

        self.repo: DataRepo = repo
        self.logger: Logger = logger

    @logger.catch
    async def GetContentProcess(self, request: pipeline_pb2.GetContentProcessReq, context) -> pipeline_pb2.GetContentProcessResp:
        self.logger.info(f"GetContentProcess req={{contentId={request.content_id}, reportVersion={request.report_version}}}")
        if request.report_version == pipeline_pb2.ReportVersion.V3:
            doc = await self.repo.mongo_db.get_content_origin(request.content_id)
            pipeline = ContentPipelineV3(logger=self.logger, doc=doc, repo=self.repo)
            await pipeline.process()
            st = Struct()
            st.update(pipeline.doc["content"])
            return pipeline_pb2.GetContentProcessResp(
                success=True,
                code=http.HTTPStatus.OK,
                msg="",
                data=st
            )

        elif request.report_version == pipeline_pb2.ReportVersion.V2:
            return pipeline_pb2.GetContentProcessResp(
                success=False,
                code=http.HTTPStatus.NOT_IMPLEMENTED,
                msg="ReportVersion corresponding method Not Implemented",
            )

        else:
            return pipeline_pb2.GetContentProcessResp(
                success=False,
                code=http.HTTPStatus.NOT_IMPLEMENTED,
                msg="ReportVersion corresponding method Not Implemented",
            )




