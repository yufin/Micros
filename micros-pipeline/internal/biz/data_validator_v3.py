import asyncio
from typing import Optional
import json
from logging import Logger
import pandas as pd
from internal.data.data import DataRepo
from internal.biz.model.company_data_stat import CompanyDataStat
from internal.biz.model.company_data_item import CompanyDataItem
from datetime import datetime
from pprint import pprint


class DataValidatorV3:
    def __init__(self, logger: Logger, doc: dict, repo: DataRepo, with_cache: bool = False):
        assert isinstance(doc, dict)
        assert doc['version'] == "V2.0"
        self.logger: Logger = logger
        self.repo: DataRepo = repo
        # self.version: str = doc.get("version", None)
        self.__doc: dict = doc
        self.__content: dict = self.doc["content"]['impExpEntReport']
        self.__items_internal: Optional[list[CompanyDataStat]] = None
        self.__items_external: Optional[list[CompanyDataStat]] = None
        self.__with_cache = with_cache
        self.__item_dict: Optional[dict[str, CompanyDataStat]] = None
        self.__time_point: Optional[datetime] = None

    @property
    def subj_company(self) -> CompanyDataStat:
        return self.item_dict.get(self.subj_company_name)

    @property
    def item_dict(self) -> Optional[dict[str, CompanyDataStat]]:
        if self.__item_dict is None:
            self.__item_dict = {}
            for item in self.items:
                self.__item_dict[item.company_name] = item
        return self.__item_dict

    @property
    def items(self) -> Optional[list[CompanyDataStat]]:
        if self.__items_internal is None:
            return None
        if self.__items_external is None:
            return self.__items_internal
        return self.__items_internal + self.__items_external

    @property
    def _partner_companies_required_item(self) -> list[CompanyDataItem]:
        return [
            CompanyDataItem.company_info,
            CompanyDataItem.ranking_list,
            CompanyDataItem.certification,
        ]

    @property
    def _subj_companies_required_item(self) -> list[CompanyDataItem]:
        return self._partner_companies_required_item + [
            CompanyDataItem.equity_transparency,
            CompanyDataItem.investment,
            CompanyDataItem.shareholder,
            CompanyDataItem.branches,
        ]

    @property
    def doc(self) -> dict:
        return self.__doc

    @property
    def content(self) -> dict:
        return self.__content

    @property
    def time_point(self) -> datetime:
        if self.__time_point is None:
            t_str = self.doc["attribute_month"] + "-01"
            self.__time_point = pd.to_datetime(t_str).to_pydatetime()
        return self.__time_point


    @property
    def subj_company_name(self) -> str:
        _info = self.content.get("businessInfo", {})
        return _info.get("QYMC")

    @property
    def subj_company_usc_id(self):
        _info = self.content.get("businessInfo", {})
        return _info.get("TYSHXYDM")

    @property
    def partner_company_names(self) -> list[str]:
        dataset = self.content.get("ssfmxxSum")
        assert isinstance(dataset, dict)

        _partner_item = [
            ("bngyssptj", "SALES_NAME"),
            ("bnkhsptj", "PURCHASER_NAME"),
            ("sngyssptj", "SALES_NAME"),
            ("snkhsptj", "PURCHASER_NAME"),
            ("ssngyssptj", "SALES_NAME"),
            ("ssnkhsptj", "PURCHASER_NAME"),
            ("sssngyssptj", "SALES_NAME"),
            ("sssnkhsptj", "PURCHASER_NAME"),
        ]

        partner_companies = []
        for _k, _n in _partner_item:
            try:
                for row in dataset[_k]:
                    partner_companies.append(row[_n])
            except TypeError:
                pass
        s = pd.Series(partner_companies)
        s.drop_duplicates(inplace=True)
        s = s[~s.isin(["其他", "合计"])]
        partner_companies = s.tolist()
        return partner_companies

    @property
    def data_stat_internal(self) -> Optional[list[CompanyDataStat]]:
        if self.__items_internal:
            return self.__items_internal

        # append up and down stream customer
        temp_data_stat = []
        for p in self.partner_company_names:
            temp_data_stat.append(CompanyDataStat(
                repo=self.repo,
                time_point=self.time_point,
                company_name=p,
                usc_id=None,
                data_item_required=self._partner_companies_required_item,
                data_item_status={},
                with_cache=self.__with_cache
            ))

        # add subject company
        temp_data_stat.append(CompanyDataStat(
            repo=self.repo,
            company_name=self.subj_company_name,
            time_point=self.time_point,
            _usc_id=self.subj_company_usc_id,
            data_item_required=self._subj_companies_required_item,
            data_item_status={},
            with_cache=self.__with_cache
        ))
        self.__items_internal = temp_data_stat
        return self.__items_internal

    @property
    async def data_stat_external(self) -> Optional[list[CompanyDataStat]]:
        if self.__items_external:
            return self.__items_external
        _, related = await self.repo.dw_data.get_related(self.subj_company_usc_id)
        if related is None:
            return None

        req = self.repo.dw_data_v3.new_pb_req(
            time_point=self.time_point,
            usc_id=self.subj_company_usc_id,
            page_size=1000,
            page_num=1
        )
        shareholders, branches, investment = await asyncio.gather(
            self.repo.dw_data_v3.get_ent_shareholders(req),
            self.repo.dw_data_v3.get_ent_branches(req),
            self.repo.dw_data_v3.get_ent_investment(req)
        )

        temp_names = []
        if isinstance(shareholders, dict):
            sh = shareholders.get("data", [])
            for n in sh:
                _name = n.get('StockName', "")
                if str(_name).__contains__("公司"):
                    temp_names.append(_name)


        for n in self.get_company_name(related, 'branch', 'enterpriseName'):
            if n and str(n).__contains__('公司'):
                temp_names.append(n)

        for n in self.get_company_name(related, 'investment', 'enterpriseName'):
            if n and str(n).__contains__('公司'):
                temp_names.append(n)

        s = pd.Series(temp_names)
        s.drop_duplicates(inplace=True)
        temp_names = s.tolist()
        temp_data_stats_external = []
        for p in temp_names:
            temp_data_stats_external.append(CompanyDataStat(
                repo=self.repo,
                time_point=self.time_point,
                company_name=p,
                data_item_required=self._partner_companies_required_item,
                with_cache=self.__with_cache
            ))

        self.__items_external = temp_data_stats_external
        return self.__items_external

    @staticmethod
    def get_company_name(data: dict, k_item: str, k_name: str) -> Optional[str]:
        item = data.get(k_item, None)
        if item:
            for d in item:
                yield d.get(k_name, None)
        else:
            yield None

    async def validate(self) -> list[dict]:
        t = [ds.verify() for ds in self.data_stat_internal] + [ds.verify() for ds in await self.data_stat_external]
        await asyncio.gather(*t)
        res = [json.loads(e.model_dump_json()) for e in self.items]
        return res

    def is_all_passed(self) -> bool:
        for item in self.items:
            if item.is_legal:
                for _, v in item.data_item_status.items():
                    if v is False:
                        return False
        return True


