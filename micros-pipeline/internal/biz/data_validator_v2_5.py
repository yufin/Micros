import asyncio
from typing import Optional
import json
from logging import Logger
import pandas as pd
from internal.data.data import DataRepo
from internal.biz.model.company_data_stat import CompanyDataStat
from internal.biz.model.company_data_item import CompanyDataItem


class DataValidatorV25:
    def __init__(self, logger: Logger, doc: dict, repo: DataRepo):
        assert isinstance(doc, dict)
        assert doc['version'] == "V1.0"
        self.logger: Logger = logger
        self.repo: DataRepo = repo
        # self.version: str = doc.get("version", None)
        self.__doc: dict = doc
        self.__content: dict = self.doc["content"]['impExpEntReport']
        self.__data_stats_internal: Optional[list[CompanyDataStat]] = None
        self.__data_stats_external: Optional[list[CompanyDataStat]] = None

    @property
    def data_stats(self) -> Optional[list[CompanyDataStat]]:
        if self.__data_stats_internal is None:
            return None
        if self.__data_stats_external is None:
            return self.__data_stats_internal
        return self.__data_stats_internal + self.__data_stats_external

    @property
    def _partner_companies_required_item(self) -> list[CompanyDataItem]:
        return [
            CompanyDataItem.company_info,
            CompanyDataItem.ranking_list,
            CompanyDataItem.certification,
        ]

    @property
    def _subj_companies_required_item(self) ->list[CompanyDataItem]:
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
    def subj_company_name(self) -> str:
        _info = self.content.get("businessInfo", {})
        return _info.get("QYMC")

    @property
    def subj_company_usc_id(self):
        _info = self.content.get("businessInfo", {})
        return _info.get("TYSHXYDM")

    @property
    def partner_company_names(self) -> list[str]:
        _partner_item = [
            ("supplierRanking_12", "SALES_NAME"),
            ("supplierRanking_24", "SALES_NAME"),
            ("customerDetail_12", "PURCHASER_NAME"),
            ("customerDetail_24", "PURCHASER_NAME")
        ]
        partner_companies = []
        for _k, _n in _partner_item:
            try:
                for row in self.content[_k]:
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
        if self.__data_stats_internal:
            return self.__data_stats_internal

        # append up and down stream customer
        temp_data_stat = []
        for p in self.partner_company_names:
            temp_data_stat.append(CompanyDataStat(
                repo=self.repo,
                company_name=p,
                usc_id=None,
                data_item_required=self._partner_companies_required_item,
                data_item_status={}
            ))

        # add subject company
        temp_data_stat.append(CompanyDataStat(
            repo=self.repo,
            company_name=self.subj_company_name,
            _usc_id=self.subj_company_usc_id,
            data_item_required=self._subj_companies_required_item,
            data_item_status={}
        ))
        self.__data_stats_internal = temp_data_stat
        return self.__data_stats_internal

    @property
    async def data_stat_external(self) -> Optional[list[CompanyDataStat]]:
        if self.__data_stats_external:
            return self.__data_stats_external
        _, related = await self.repo.dw_data.get_related(self.subj_company_usc_id)
        if related is None:
            return None
        temp_names = []
        for n in self.get_company_name(related,'shareholder', 'shareholderName'):
            if n and str(n).__contains__('公司'):
                temp_names.append(n)

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
                company_name=p,
                data_item_required=self._partner_companies_required_item,
            ))

        self.__data_stats_external = temp_data_stats_external
        return self.__data_stats_external

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
        return [json.loads(e.model_dump_json()) for e in self.data_stats]

    