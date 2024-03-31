from api.pipeline.v1 import pipeline_pb2_grpc, pipeline_pb2
from internal.data.data import DataRepo
from google.protobuf.struct_pb2 import Struct
from internal.biz.content_pipeline_v2_5 import ContentPipelineV25
from internal.biz.content_pipeline_v3 import ContentPipelineV3
from internal.biz.data_validator_v2_5 import DataValidatorV25
from internal.biz.data_validator_v3 import DataValidatorV3
from internal.biz.ahp_request import AhpReqeust
import http
from logging import Logger
from loguru import logger
import json
from google.protobuf import json_format

class PipelineService(pipeline_pb2_grpc.PipelineServiceServicer):
    def __init__(self, repo: DataRepo, logger: Logger):
        self.repo: DataRepo = repo
        self.logger: Logger = logger

    @logger.catch
    async def GetTradeDetail(self, request: pipeline_pb2.GetTradeDetailReq, context) -> pipeline_pb2.GetTradeDetailResp:
        self.logger.info(
            f"GetTradeDetail req={{contentId={request.content_id}, reportVersion={request.report_version}}}")

        doc = await self.repo.mongo_db.get_content(request.content_id)
        if not isinstance(doc, dict):
            return pipeline_pb2.GetTradeDetailResp(
                success=False,
                code=http.HTTPStatus.NO_CONTENT,
                msg="Content not found",
            )
        version = doc['version']

        if request.report_version == pipeline_pb2.ReportVersion.V3:
            if version == 'V2.0':
                pipeline = ContentPipelineV3(logger=self.logger, doc=doc, repo=self.repo)
                res = await pipeline.sales_detail_chart(
                    list(request.option_time_period),
                    request.option_top_cus,
                    list(request.option_trade_frequency),
                    request.trade_type
                )
                st = Struct()
                st.update(res)

                return pipeline_pb2.GetTradeDetailResp(
                    success=True,
                    code=http.HTTPStatus.OK,
                    msg="",
                    data=st
                )

        return pipeline_pb2.GetTradeDetailResp(
            success=False,
            code=http.HTTPStatus.NOT_IMPLEMENTED,
            msg="ReportVersion not support for this service",
        )

    @logger.catch
    async def GetContentProcess(self, request: pipeline_pb2.GetContentProcessReq, context) -> pipeline_pb2.GetContentProcessResp:
        self.logger.info(f"GetContentProcess req={{contentId={request.content_id}, reportVersion={request.report_version}}}")

        doc = await self.repo.mongo_db.get_content(request.content_id)
        if not isinstance(doc, dict):
            return pipeline_pb2.GetContentProcessResp(
                success=False,
                code=http.HTTPStatus.NO_CONTENT,
                msg="Content not found",
            )
        version = doc['version']

        if request.report_version == pipeline_pb2.ReportVersion.V2_5:
            if version == 'V1.0':
                pipeline = ContentPipelineV25(logger=self.logger, doc=doc, repo=self.repo)
                await pipeline.process()
                st = Struct()
                st.update(pipeline.doc["content"])

                return pipeline_pb2.GetContentProcessResp(
                    success=True,
                    code=http.HTTPStatus.OK,
                    msg="",
                    data=st
                )
        if request.report_version == pipeline_pb2.ReportVersion.V3:
            if version == 'V2.0':

                if request.lang == "en-US":
                    res = await self.repo.mongo_db.get_processed_i18n_content(
                        content_id=request.content_id,
                        lang=request.lang,
                    )
                    if not res:
                        return pipeline_pb2.GetContentProcessResp(
                            success=False,
                            code=http.HTTPStatus.NO_CONTENT,
                            msg="Content not found",
                        )
                    st = Struct()
                    st.update(res["content"])
                    return pipeline_pb2.GetContentProcessResp(
                        success=True,
                        code=http.HTTPStatus.OK,
                        msg="",
                        data=st
                    )

                pipeline = ContentPipelineV3(logger=self.logger, doc=doc, repo=self.repo)
                await pipeline.process()
                st = Struct()
                st.update(pipeline.doc["content"])

                # save pipeline.doc['content'] as json file
                # with open(f"../../static/{request.content_id}.json", "w") as f:
                #     json.dump(pipeline.content, f)

                return pipeline_pb2.GetContentProcessResp(
                    success=True,
                    code=http.HTTPStatus.OK,
                    msg="",
                    data=st
                )

        return pipeline_pb2.GetContentProcessResp(
            success=False,
            code=http.HTTPStatus.NOT_IMPLEMENTED,
            msg="ReportVersion corresponding method Not Implemented",
        )

    @logger.catch
    async def GetContentValidate(self, request: pipeline_pb2.GetContentProcessReq, context) -> pipeline_pb2.GetContentValidateResp:
        self.logger.info(f"GetContentValidate req={{contentId={request.content_id}, reportVersion={request.report_version}}}")

        doc = await self.repo.mongo_db.get_content(request.content_id)
        if not isinstance(doc, dict):
            return pipeline_pb2.GetContentValidateResp(
                success=False,
                code=http.HTTPStatus.NO_CONTENT,
                msg="Content not found",
            )
        version = doc['version']

        if request.report_version == pipeline_pb2.ReportVersion.V2_5:
            if version == 'V1.0':
                validator = DataValidatorV25(logger=self.logger, doc=doc, repo=self.repo)
                res = await validator.validate()
                repeat_st = []
                for d in res:
                    st = Struct()
                    st.update(d)
                    repeat_st.append(st)

                return pipeline_pb2.GetContentValidateResp(
                    success=True,
                    code=http.HTTPStatus.OK,
                    msg="",
                    data=repeat_st
                )

        elif request.report_version == pipeline_pb2.ReportVersion.V3:
            if version == 'V2.0':
                validator = DataValidatorV3(logger=self.logger, doc=doc, repo=self.repo)
                res = await validator.validate()
                repeat_st = []
                for d in res:
                    st = Struct()
                    st.update(d)
                    repeat_st.append(st)

                return pipeline_pb2.GetContentValidateResp(
                    success=True,
                    code=http.HTTPStatus.OK,
                    msg="",
                    data=repeat_st
                )

        return pipeline_pb2.GetContentValidateResp(
            success=False,
            code=http.HTTPStatus.NOT_IMPLEMENTED,
            msg="Internal Server Error",
        )

    @logger.catch
    async def GetAhpScore(self, request: pipeline_pb2.GetAhpScoreReq, context) -> pipeline_pb2.GetAhpScoreResp:
        use_local = True
        self.logger.info(f"GetAhpScore req={{claimId={request.claim_id}}}")
        try:
            ahp = AhpReqeust(logger=self.logger, repo=self.repo, claim_id=request.claim_id)
            if use_local:
                res = await ahp.local_calculate()
                if res is None:
                    raise Exception("Local calculate failed")
                result = {"data": res.model_dump()}
            else:
                result, msg = await ahp.request()
                if msg != "":
                    raise Exception(msg)
                assert isinstance(result, dict)
        except Exception as e:
            self.logger.error(f"GetAhpScore failed, error: {e}")
            return pipeline_pb2.GetAhpScoreResp(
                success=False,
                code=http.HTTPStatus.INTERNAL_SERVER_ERROR,
                msg=str(e),
            )
        data = result.get("data", {})
        st = Struct()
        st.update(data)
        return pipeline_pb2.GetAhpScoreResp(
            success=True,
            code=http.HTTPStatus.OK,
            msg="",
            data=st
        )

    @logger.catch
    async def GetJsonTranslate(self, request: pipeline_pb2.GetJsonTranslateReq, context) -> pipeline_pb2.GetJsonTranslateResp:
        self.logger.info("GetJsonTranslate called")
        st_data = request.data
        # replace_nan_in_struct(st_data)
        data = json_format.MessageToDict(st_data)

        res_str, err = self.repo.translator.trans_json(data)
        res = json.loads(res_str)
        st = Struct()
        st.update(res)
        return pipeline_pb2.GetJsonTranslateResp(
            success=True,
            code=http.HTTPStatus.OK,
            msg="",
            data=st
        )


