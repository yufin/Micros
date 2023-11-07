from api.dwdata.v2 import dw_data_pb2, dw_data_pb2_grpc
import grpc
from typing import Sequence, Optional, Tuple
import asyncio
from google.protobuf import json_format


class DwDataClient:
    def __init__(self, target: str, options: Optional[Sequence[Tuple[str, any]]]):
        self.target = target
        self.options = options

    @property
    def channel(self):
        return grpc.aio.insecure_channel(target=self.target, options=self.options)

    async def get_ent_branches(self, usc_id: str) -> dw_data_pb2.GetBranchesResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntBranches(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_investment(self, usc_id: str) -> dw_data_pb2.GetInvestmentResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntInvestment(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_shareholders(self, usc_id: str) -> dw_data_pb2.GetShareholdersResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntShareholders(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_equity_transparency(self, usc_id: str) -> dw_data_pb2.GetEquityTransparencyResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntEquityTransparency(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_product(self, usc_id: str) -> dw_data_pb2.GetEntStrArrayResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseProduct(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_industry(self, usc_id: str) -> dw_data_pb2.GetEntStrArrayResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseIndustry(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_ranking_list(self, usc_id: str) -> dw_data_pb2.GetEntRankingListResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseRankingList(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_credential(self, usc_id: str) -> dw_data_pb2.GetEntCredentialResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseCredential(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_info(self, usc_id: str) -> dw_data_pb2.GetEntInfoResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseInfo(dw_data_pb2.GetEntInfoReq(usc_id=usc_id))

    async def get_ent_ident(self, enterprise_name: str) -> dw_data_pb2.GetEntIdentResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseIdent(dw_data_pb2.GetEntIdentReq(enterprise_name=enterprise_name))

    async def get_tags(self, usc_id: str, with_info: bool = False) -> (str, dict):
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            if with_info:
                info, industry_resp, product_resp, credential_resp, ranking_list_resp = await asyncio.gather(
                    stub.GetEnterpriseInfo(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                    stub.GetEnterpriseIndustry(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                    stub.GetEnterpriseProduct(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                    stub.GetEnterpriseCredential(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                    stub.GetEnterpriseRankingList(dw_data_pb2.GetEntInfoReq(usc_id=usc_id)),
                )
                return usc_id, {
                    "companyInfo": json_format.MessageToDict(info.data) if info.success else None,
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
        ident = await self.get_ent_ident(name)
        if ident.success:
            res = await self.get_tags(ident.data.usc_id, with_info)
            return name, res[1]
        else:
            return name, None






# if __name__ == '__main__':
#     import asyncio
#     dw = DwDataClient(target="192.168.44.150:50052", options=None)
#
#     # res1 = asyncio.run(dw.get_ent_investment(dw_data_pb2.GetEntInfoReq(usc_id="91440300071131008W")))
#     # res2 = asyncio.run(dw.get_ent_investment(dw_data_pb2.GetEntInfoReq(usc_id="91440300071131008W")))
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
