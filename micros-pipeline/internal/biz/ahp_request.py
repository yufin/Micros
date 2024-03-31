from internal.biz.model.ahp_params import AhpParams
from internal.biz.ahp_local import AhpLocal
from internal.data.data import DataRepo
from internal.biz.model.ahp_local_res import *
from logging import Logger
import aiohttp
from typing import Optional
import asyncio


class AhpReqeust:
    def __init__(self, logger: Logger, repo: DataRepo, claim_id: int, doc: Optional[dict] = None):
        self.repo: DataRepo = repo
        self.logger: Logger = logger
        self.claim_id: int = claim_id
        self._params: Optional[AhpParams] = None
        self.doc: Optional[dict] = doc

    @property
    def url(self):
        return "http://10.0.203.188:8086/decision-engine/decision/task/sync/SCENE_1/LH_AHP_SCR"

    @property
    def request_body(self) -> dict:
        assert self._params is not None
        return {'param': self._params.model_dump()}

    async def request(self) -> (Optional[dict], Optional[str]):

        await self.instant_param()
        print("ahp_params: ", self._params.model_dump())

        async with aiohttp.ClientSession() as session:
            async with session.post(self.url, json=self.request_body) as resp:
                if resp.status != 200:
                    msg = await resp.text()
                    self.logger.error(f"request failed, status: {resp.status}, text: {msg}")
                    resp.close()
                    return None, msg
                data = await resp.json()
                resp.close()
                return data, ""

    async def local_calculate(self) -> Optional[AhpResult]:
        await self.instant_param()
        assert self._params is not None
        try:
            ahp_local = AhpLocal()
            l1 = ahp_local.scorecard_eval(self._params.model_dump(), self._params.lh_qylx, 1)
            fl1 = L1Factor(**l1)
            l2 = ahp_local.scorecard_eval(fl1.model_dump(), self._params.lh_qylx, 2)
            fl2 = L2Factor(**l2)
            l3 = ahp_local.scorecard_eval(fl2.model_dump(), self._params.lh_qylx, 3)
            fl3 = L3Factor(**l3)
            l4 = ahp_local.scorecard_eval(fl3.model_dump(), self._params.lh_qylx, 4)
            fl4 = L4Factor(**l4)
            res = AhpResult(l0=self._params, l1=fl1, l2=fl2, l3=fl3, l4=fl4, )
            return res
        except Exception as e:
            self.logger.error(f"local_calculate failed, {e}")
        return None

    async def instant_param(self):
        claim = await self.repo.sql_db.query_one(
            f"select * from rc_content_factor_claim where id = {self.claim_id};")
        if not claim:
            claim = await self.repo.sql_db.query_one(
                f"select * from rc_content_factor_claim_v3 where id = {self.claim_id};")

        assert isinstance(claim, dict)
        content_id = claim.get('content_id')
        assert isinstance(content_id, int)
        factor_id = claim.get('factor_id')
        assert isinstance(factor_id, int)

        factors, doc = await asyncio.gather(
            self.repo.sql_db.query_one(f"select * from rc_decision_factor where id = {factor_id};"),
            self.repo.mongo_db.get_content(content_id)
        )
        assert isinstance(factors, dict)
        assert isinstance(doc, dict)
        models_indexes = doc.get('content',{}).get("modelIndexes")
        assert isinstance(models_indexes, list)
        usc_id = doc.get('usc_id')
        assert isinstance(usc_id, str)

        p = {}
        for item in models_indexes:
            k = item['INDEX_CODE']
            v = item['INDEX_VALUE']
            p[k] = v
        ahp_params = AhpParams(
            **p,
            lh_cylwz=factors.get('lh_cylwz'),
            gs_gdct=factors.get('lh_gdct'),
            md_qybq=factors.get('lh_qybq'),
            lh_qylx=factors.get('lh_qylx'),
            zx_dsfsxqk=factors.get('lh_sfsx'),
            zx_yhsxqk=factors.get('lh_yhsx'),
            nsrsbh=usc_id
        )
        self._params = ahp_params
