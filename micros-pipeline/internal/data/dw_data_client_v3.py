from api.dwdata.v3 import dw_data_v3_pb2, dw_data_v3_pb2_grpc
from api.dwdata.v2 import dw_data_pb2, dw_data_pb2_grpc
import grpc
from typing import Sequence, Optional, Tuple
import asyncio
from google.protobuf import json_format
from logging import Logger
from grpc.aio._channel import Channel
from datetime import datetime
from google.protobuf.timestamp_pb2 import Timestamp
import datetime


class DwDataClientV3:
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

        async def close_channel():
            if self._channel is not None:
                await self._channel.close()

        loop = asyncio.get_event_loop()
        try:
            loop.run_until_complete(close_channel())
            # new_loop.run_until_complete(self._channel.close())
        except Exception as e:
            print("grpc conn close error:", e)
        finally:
            # new_loop.close()
            return

    @property
    def channel(self) -> Channel:
        if self._channel is None or self._channel._channel.check_connectivity_state(True) == grpc.ChannelConnectivity.SHUTDOWN:
            self._channel = grpc.aio.secure_channel(
                target=self.target,
                credentials=self.__client_credentials,
                options=self.options,
            )
        return self._channel

    async def get_usc_id_by_name(self, enterprise_name: str) -> dw_data_v3_pb2.GetUscIdByEnterpriseNameResp:
        self.logger.info("get_usc_id_by_name enterprise_name=" + enterprise_name)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        return await stub.GetUscIdByEnterpriseName(dw_data_v3_pb2.GetUscIdByEnterpriseNameReq(enterprise_name=enterprise_name))

    async def get_ent_info(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_info usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataDictResp = await stub.GetEnterpriseInfo(req)
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_ent_branches(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_branches usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetEntBranches(req)
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_ent_investment(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_investment usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetEntInvestment(req)
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_ent_shareholders(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_shareholders usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetEntShareholders(req)
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_ent_equity_transparency(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_equity_transparency usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetEntEquityTransparency(req)
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_ent_equity_transparency_conclusion(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
        self.logger.info("get_ent_equity_transparency_conclusion usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataDictResp = await stub.GetEntEquityTransparencyConclusion(req)
        if not res.success:
            return None
        d = json_format.MessageToDict(res)
        return d.get("data", {}).get("conclusion")

    async def get_ent_credential(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_ranking_list usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetEnterpriseCredential(req)
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_court_announcement(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_court_announcement usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetCourtAnnouncement(req)
        if not res.success:
            return None
        d = json_format.MessageToDict(res)
        if d.get("data") is None:
            d['data'] = []
        return d

    async def get_executive(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_executive usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetExecutive(req)
        if not res.success:
            return None
        d = json_format.MessageToDict(res)
        if d.get("data") is None:
            d['data'] = []
        return d

    async def get_historical_executive(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_executive usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetHistoricalExecutive(req)
        if not res.success:
            return None
        d = json_format.MessageToDict(res)
        if d.get("data") is None:
            d['data'] = []
        return d

    async def get_case_registration_info(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_case_registration_info usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetCaseRegistrationInfo(req)
        if not res.success:
            return None
        d = json_format.MessageToDict(res)
        if d.get("data") is None:
            d['data'] = []
        return d

    async def get_ent_financing(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_financing usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        res: dw_data_v3_pb2.GetDataListResp = await stub.GetEntFinancing(req)
        if not res.success:
            return None
        d = json_format.MessageToDict(res)
        if d.get("data") is None:
            d['data'] = []
        return d

    async def get_ent_ranking_list(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[dict]:
        self.logger.info("get_ent_credential usc_id=" + req.usc_id)
        stub = dw_data_v3_pb2_grpc.DwdataServiceStub(self.channel)
        try:
            res: dw_data_v3_pb2.GetEntRankingListResp = await stub.GetEnterpriseRankingList(req)
        except Exception:
            return None
        if not res.success:
            return None
        return json_format.MessageToDict(res)

    async def get_tags_by_name(self, name: str, tp: datetime.datetime) -> tuple[str, Optional[dict]]:
        self.logger.info("get_tags_by_name name=" + name)

        ident = await self.get_usc_id_by_name(name)

        if ident and ident.success:
            req = self.new_pb_req(
                usc_id=ident.data.usc_id,
                time_point=tp
            )
            info, ranking, credential, equity_conclusion = await asyncio.gather(
                self.get_ent_info(req),
                self.get_ent_ranking_list(req),
                self.get_ent_credential(req),
                self.get_ent_equity_transparency_conclusion(req)
            )

            resp = {
                "companyInfo": info.get("data") if info else None,
                "authorizedTag": credential.get("data") if credential else None,
                "rankingTag": ranking.get("data") if ranking else None,
                "equityConclusion": equity_conclusion if equity_conclusion else None
            }
            return name, resp
        else:
            return name, None

    @staticmethod
    def to_pb_timestamp(dt: datetime) -> Timestamp:
        ts = Timestamp()
        ts.FromDatetime(dt)
        return ts

    @staticmethod
    def new_pb_req(time_point: datetime.datetime, usc_id: str, page_size: int = 30, page_num: int = 1) -> dw_data_v3_pb2.GetDataBeforeTimePointReq:
        ts = Timestamp()
        ts.FromDatetime(time_point)
        return dw_data_v3_pb2.GetDataBeforeTimePointReq(
            usc_id=usc_id,
            time_point=ts,
            page_size=page_size,
            page_num=page_num,
        )


if __name__ == '__main__':
    import asyncio
    from loguru import logger
    dw = DwDataClientV3(logger=logger, target="192.168.44.150:50052", certs_folder="../../../certs")

    dt = datetime.datetime.strptime("2023-12-08 12:30:45", "%Y-%m-%d %H:%M:%S")

    req = dw_data_v3_pb2.GetDataBeforeTimePointReq(
        usc_id="91440300576351268T",
        time_point=dw.to_pb_timestamp(dt),
        page_size=20,
        page_num=1,
    )

    async def test():
        _res = await asyncio.gather(
            # dw.get_ent_equity_transparency(req),
            dw.get_ent_ranking_list(req),
        )
        return _res

    res1 = asyncio.run(test())

    from pprint import pprint
    pprint(res1)
    print(type(res1))


    t = '2023-12-08 10:00:12'
    import pandas as pd

    dt = pd.to_datetime(t).to_pydatetime()
    print(dt)
    print(type(dt))
    pbt = dw.to_pb_timestamp(dt)
    print(pbt)
    # print(type(dt))
    # print(dt.to_datetime64())
    # print(type(dt.to_datetime64()))



    # print(res1)

