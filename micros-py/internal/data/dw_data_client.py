from api.dwdata.v2 import dw_data_pb2, dw_data_pb2_grpc
import grpc
from typing import Sequence, Optional, Tuple


class DwDataClient:
    def __init__(self, target: str, options: Optional[Sequence[Tuple[str, any]]]):
        self.target = target
        self.options = options

    @property
    def channel(self):
        return grpc.aio.insecure_channel(target=self.target, options=self.options)

    async def get_ent_branches(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.BranchesResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntBranches(req)

    async def get_ent_investment(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.InvestmentResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntInvestment(req)

    async def get_ent_shareholders(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.ShareholdersResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntShareholders(req)

    async def get_ent_equity_transparency(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.EquityTransparencyResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEntEquityTransparency(req)

    async def get_ent_product(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.EntStrArrayResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseProduct(req)

    async def get_ent_industry(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.EntStrArrayResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseIndustry(req)

    async def get_ent_ranking_list(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.EntRankingListResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseRankingList(req)

    async def get_ent_credential(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.EntCredentialResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseCredential(req)

    async def get_ent_info(self, req: dw_data_pb2.GetEntInfoReq) -> dw_data_pb2.EntInfoResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseInfo(req)

    async def get_ent_ident(self, req: dw_data_pb2.GetEntIdentReq) -> dw_data_pb2.EntIdentResp:
        async with self.channel as channel:
            stub = dw_data_pb2_grpc.DwdataServiceStub(channel)
            return await stub.GetEnterpriseIdent(req)


if __name__ == '__main__':
    import asyncio
    dw = DwDataClient(target="192.168.44.150:50052", options=None)

    # res1 = asyncio.run(dw.get_ent_investment(dw_data_pb2.GetEntInfoReq(usc_id="91440300071131008W")))
    # res2 = asyncio.run(dw.get_ent_investment(dw_data_pb2.GetEntInfoReq(usc_id="91440300071131008W")))
    res3 = asyncio.run(dw.get_ent_equity_transparency(dw_data_pb2.GetEntInfoReq(usc_id="91310105134638405A")))
    print(res3)

    from google.protobuf import json_format
    print(type(res3))
    rd = json_format.MessageToDict(res3)
    print(rd)
    print(type(rd))

    # print(res1)

