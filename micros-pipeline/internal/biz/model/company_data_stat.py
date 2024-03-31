from pydantic import BaseModel, Field
from typing import Optional
from internal.data.data import DataRepo
from internal.biz.model.company_data_item import CompanyDataItem
from datetime import datetime


class CompanyDataStat(BaseModel):
    # repo: DataRepo
    company_name: str
    is_legal: Optional[bool] = None
    usc_id: Optional[str] = None
    data_item_required: list[CompanyDataItem]
    data_item_status: Optional[dict[CompanyDataItem, Optional[bool]]] = None

    def __init__(self, repo: DataRepo, time_point: datetime,with_cache: bool = False, **data):
        super().__init__(**data)
        self.__repo: DataRepo = repo
        self.__data_cached: Optional[dict[CompanyDataItem, Optional[dict]]] = None
        self._with_cache = with_cache
        self.__time_point = time_point


    @property
    def data_cached(self) -> Optional[dict[CompanyDataItem, Optional[dict]]]:
        return self.__data_cached

    @property
    def repo(self):
        return self.__repo

    async def verify(self):
        self.data_item_status = {}
        for item in self.data_item_required:
            self.data_item_status[item] = await self._validate_data_item(item)
        return self.data_item_status

    @property
    async def _usc_id(self) -> Optional[str]:
        if self.is_legal:
            return self.usc_id
        if self.is_legal is None:
            resp = await self.repo.dw_data_v3.get_usc_id_by_name(self.company_name)
            # resp = await self.repo.dw_data.get_ent_ident(self.company_name)
            # jsonschema.
            if resp:
                if resp.success:
                    self.usc_id = resp.data.usc_id
                    self.is_legal = resp.data.isLegal
                    return self.usc_id
        return None

    async def _validate_data_item(self, item: CompanyDataItem) -> Optional[bool]:
        methods_collection = {
            CompanyDataItem.company_info: self.repo.dw_data_v3.get_ent_info,
            CompanyDataItem.ranking_list: self.repo.dw_data_v3.get_ent_ranking_list,
            CompanyDataItem.certification: self.repo.dw_data_v3.get_ent_credential,
            CompanyDataItem.shareholder: self.repo.dw_data_v3.get_ent_shareholders,
            CompanyDataItem.branches: self.repo.dw_data_v3.get_ent_branches,
            CompanyDataItem.investment: self.repo.dw_data_v3.get_ent_investment,
            CompanyDataItem.equity_transparency: self.repo.dw_data_v3.get_ent_equity_transparency,
        }

        usc_id = await self._usc_id
        if self.is_legal is True:
            func = methods_collection.get(item)
            assert isinstance(usc_id, str)

            req = self.repo.dw_data_v3.new_pb_req(
                usc_id=usc_id,
                time_point=self.__time_point,
                page_size=1
            )

            resp = await func(req)

            if self._with_cache:
                if self.__data_cached is None:
                    self.__data_cached = {}
                self.__data_cached[item] = resp

            if resp is None:
                return False
            return resp.get('success')
        else:
            return None
