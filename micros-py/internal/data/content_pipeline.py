import asyncio

import pandas as pd

from internal.data.data import DataRepo
from pandas import Series
from pprint import pprint

pd.set_option('display.max_columns', None)
pd.set_option('display.max_rows', None)


class ContentPipelineLatest:

    def __init__(self, doc: dict, repo: DataRepo):
        # assert doc["version"] == "V2.0"
        self.doc = doc
        self.content: dict = self.doc["content"]['impExpEntReport']
        self.repo = repo

    async def process(self):
        pprint(self.content)
        print("===================================== \n")

        self._major_commodity_proportion("customerDetail_12")

        # 1.add tag for companies
        tag_map = await self._gather_tag_dict()
        for row in self.content["customerDetail_12"]:
            row["QYBQ"] = tag_map[row["PURCHASER_NAME"]]

        for row in self.content["customerDetail_24"]:
            row["QYBQ"] = tag_map[row["PURCHASER_NAME"]]

        for row in self.content["supplierRanking_12"]:
            row["QYBQ"] = tag_map[row["SALES_NAME"]]

        for row in self.content["supplierRanking_24"]:
            row["QYBQ"] = tag_map[row["SALES_NAME"]]

        # 2.add capitalPaidIn field to businessInfo
        try:
            subj_d = tag_map[self.content["businessInfo"]["QYMC"]]
            self.content["businessInfo"]["capitalPaidIn"] = None
            if subj_d is not None:
                self.content["businessInfo"]["capitalPaidIn"] = subj_d["companyInfo"]["paidInCapital"]
            else:
                pass
        except KeyError:
            pass

        # 3.apply majorCommodityProportion calculation for trades

    def _major_commodity_proportion(self, key_trades: str) -> dict:
        _df_trades = pd.DataFrame(self.content[key_trades])

        _df_trades['_major_commodity'] = _df_trades["COMMODITY_NAME"].apply(lambda x: x.split("*")[1] if str(x).__contains__("*") else None)

        _sum_amount_tax: float = _df_trades["SUM_AMOUNT_TAX"].apply(lambda x: self._to_numeric(x)).sum()


        print(_df_trades)
        print(_sum_amount_tax)

    @staticmethod
    def _to_numeric(x: str) -> float:
        try:
            return float(x.replace(",", "").replace("%", ""))
        except ValueError:
            return 0.

    async def _gather_tag_dict(self) -> dict:
        companies: list[str] = [self.content["businessInfo"]["QYMC"]]
        companies += [row["PURCHASER_NAME"] for row in self.content["customerDetail_12"]]
        companies += [row["PURCHASER_NAME"] for row in self.content["customerDetail_24"]]
        companies += [row["SALES_NAME"] for row in self.content["supplierRanking_12"]]
        companies += [row["SALES_NAME"] for row in self.content["supplierRanking_24"]]
        s = Series(companies)
        s.drop_duplicates(inplace=True)
        tags = await asyncio.gather(*(self.repo.dw_data.get_tags_by_name(row, True) for row in s))
        return {name: tag for name, tag in tags}




