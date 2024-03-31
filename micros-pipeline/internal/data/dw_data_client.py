from api.dwdata.v2 import dw_data_pb2, dw_data_pb2_grpc
import grpc
from typing import Sequence, Optional, Tuple
import asyncio
from google.protobuf import json_format
from logging import Logger
from grpc.aio._channel import Channel


class DwDataClient:
    def __init__(self,
                 logger: Logger,
                 target: str,
                 certs_folder: str,
                 options: Optional[Sequence[tuple]] = None,
                 ):
        if options is None:
            options = [
                ('grpc.max_receive_message_length', 524288000),
                ('grpc.max_send_message_length', 524288000),
                ('grpc.max_message_length', 524288000),
                ('grpc.max_metadata_size', 524288000),
            ]
        self.options = options
        self.target = target
        self.logger = logger
        self._channel: Optional[Channel] = None


        with open(certs_folder+"/client_key.pem", "rb") as f:
            private_key = f.read()
        with open(certs_folder+"/client_cert.pem", "rb") as f:
            certificate_chain = f.read()
        with open(certs_folder+"/ca_cert.pem", "rb") as f:
            ca_certificates = f.read()

        self.__client_credentials = grpc.ssl_channel_credentials(
            root_certificates=ca_certificates,
            private_key=private_key,
            certificate_chain=certificate_chain,
        )

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self._channel is not None:
            await self._channel.close()

    def __del__(self):
        if self._channel is not None:
            new_loop = asyncio.new_event_loop()
            try:
                new_loop.run_until_complete(self._channel.close())
            finally:
                new_loop.close()

    @property
    def channel(self) -> Channel:
        if self._channel is None or self._channel._channel.check_connectivity_state(True) == grpc.ChannelConnectivity.SHUTDOWN:
            self._channel = grpc.aio.secure_channel(
                target=self.target,
                credentials=self.__client_credentials,
                options=self.options,
            )
        return self._channel

    async def get_ent_branches(self, usc_id: str) -> dw_data_pb2.GetBranchesResp:
        self.logger.info("get_ent_branches usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEntBranches(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_investment(self, usc_id: str) -> dw_data_pb2.GetInvestmentResp:
        self.logger.info("get_ent_investment usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEntInvestment(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_shareholders(self, usc_id: str) -> dw_data_pb2.GetShareholdersResp:
        self.logger.info("get_ent_shareholders usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEntShareholders(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_equity_transparency(self, usc_id: str) -> dw_data_pb2.GetEquityTransparencyResp:
        self.logger.info("get_ent_equity_transparency usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEntEquityTransparency(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_product(self, usc_id: str) -> dw_data_pb2.GetEntStrArrayResp:
        self.logger.info("get_ent_product usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEnterpriseProduct(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_industry(self, usc_id: str) -> dw_data_pb2.GetEntStrArrayResp:
        self.logger.info("get_ent_industry usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEnterpriseIndustry(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_ranking_list(self, usc_id: str) -> dw_data_pb2.GetEntRankingListResp:
        self.logger.info("get_ent_ranking_list usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEnterpriseRankingList(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_credential(self, usc_id: str) -> dw_data_pb2.GetEntCredentialResp:
        self.logger.info("get_ent_credential usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEnterpriseCredential(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_info(self, usc_id: str) -> dw_data_pb2.GetEntInfoResp:
        self.logger.info("get_ent_info usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEnterpriseInfo(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_ident(self, enterprise_name: str) -> dw_data_pb2.GetEntIdentResp:
        self.logger.info("get_ent_ident enterprise_name=" + enterprise_name)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetEnterpriseIdent(dw_data_pb2.GetEntIdentReq(enterprise_name=enterprise_name))

    async def get_tags(self, usc_id: str, with_info: bool = False) -> (str, dict):
        self.logger.info("get_tags usc_id=" + usc_id)
        if not usc_id:
            if with_info:
                return usc_id, {
                    "companyInfo": None,
                    "industryTag": None,
                    "productTag": None,
                    "authorizedTag": None,
                    "rankingTag": None,
                }
            else:
                return usc_id, {
                    "industryTag": None,
                    "productTag": None,
                    "authorizedTag": None,
                    "rankingTag": None,
                }

        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        if with_info:
            info, industry_resp, product_resp, credential_resp, ranking_list_resp, equity_resp = await asyncio.gather(
                stub.GetEnterpriseInfo(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseIndustry(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseProduct(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseCredential(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseRankingList(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEntEquityTransparency(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
            )

            company_info_d = json_format.MessageToDict(info.data) if info.success else None

            equity_conclusion = None
            if equity_resp:
                equity_conclusion = equity_resp.conclusion if equity_resp.success else None
            if company_info_d:
                company_info_d.update({"equityConclusion": equity_conclusion})

            return usc_id, {
                "companyInfo": company_info_d if company_info_d else None,
                "industryTag": [_i for _i in industry_resp.data] if industry_resp.success else None,
                "productTag": [_p for _p in product_resp.data] if product_resp.success else None,
                "authorizedTag": [json_format.MessageToDict(cred) for cred in credential_resp.data[:20]] if credential_resp.success else None,
                "rankingTag": [json_format.MessageToDict(rank) for rank in ranking_list_resp.data[:20]] if ranking_list_resp.success else None,
            }
        else:
            industry_resp, product_resp, credential_resp, ranking_list_resp = await asyncio.gather(
                stub.GetEnterpriseIndustry(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseProduct(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseCredential(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                stub.GetEnterpriseRankingList(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
            )
            return usc_id, {
                "industryTag": [_i for _i in industry_resp.data] if industry_resp.success else None,
                "productTag": [_p for _p in product_resp.data] if product_resp.success else None,
                "authorizedTag": [json_format.MessageToDict(cred) for cred in credential_resp.data[:20]] if credential_resp.success else None,
                "rankingTag": [json_format.MessageToDict(rank) for rank in ranking_list_resp.data[:20]] if ranking_list_resp.success else None,
            }

    async def get_tags_by_name(self, name: str, with_info: bool = False) -> (str, Optional[dict]):
        self.logger.info("get_tags_by_name name=" + name)
        ident = await self.get_ent_ident(name)
        if ident.success:
            res = await self.get_tags(ident.data.usc_id, with_info)
            return name, res[1]
        else:
            return name, None

    async def get_related(self, usc_id: str) -> (str, Optional[dict]):
        self.logger.info("get_related usc_id=" + usc_id)
        stub = dw_data_pb2_grpc.DwdataServiceStub(self.channel)
        shareholders, branches, investment, equity = await asyncio.gather(
            stub.GetEntShareholders(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
            stub.GetEntBranches(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
            stub.GetEntInvestment(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
            stub.GetEntEquityTransparency(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
        )
        return usc_id, {
            'shareholder': [json_format.MessageToDict(shareholder) for shareholder in shareholders.data] if shareholders.success else None,
            'branch': [json_format.MessageToDict(branch) for branch in branches.data] if branches.success else None,
            'investment': [json_format.MessageToDict(invest) for invest in investment.data] if investment.success else None,
            'equityTransparency': [json_format.MessageToDict(eq) for eq in equity.data] if equity.success else None,
            'equityConclusion': equity.conclusion if equity.success else None
        }




# if __name__ == '__main__':
#     import asyncio
#     from loguru import logger
#     dw = DwDataClient(logger=logger, target="192.168.44.150:50052")
#     # res1 = asyncio.run(dw.get_ent_info(usc_id="913100006073515956"))
#     # res1 = asyncio.run(dw.get_ent_info(usc_id="913100006073515956"))
#
#     async def test():
#         res1 = await asyncio.gather(
#             dw.get_ent_branches(usc_id="91310000MA1GU5JD80"),
#             dw.get_ent_branches(usc_id="91310000MA1GU5JD80"),
#             dw.get_ent_branches(usc_id="91310000MA1GU5JD80"),
#             dw.get_ent_branches(usc_id="91310000MA1GU5JD80"),
#         )
#         return res1
#
#     res1 = asyncio.run(test())
#     print(res1)

    # print(dw._trans_json(res1) )
#
#     res1 = dw.get_ent_equity_transparency_sync(usc_id="91310000MA1GU5JD80")
#
#     from pprint import pprint
#     pprint(res1)


    # res2 = asyncio.run(dw.get_ent_investment(dw_data_pb2.GetEntInfoReq(usc_id="91440300071131008W")))
#     res3 = asyncio.run(dw.get_ent_equity_transparency(dw_data_pb2.GetEntInfoReq(usc_id="91310105134638405A")))
#     print(res3)
#
#     from google.protobuf import json_format
#     print(type(res3))
#     rd = json_format.MessageToDict(res3)
#     print(rd)
#     print(type(rd))
#
#     # print(res1)
#
