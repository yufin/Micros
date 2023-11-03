import asyncio
import copy

from api.dwdata.v2 import dw_data_pb2
from typing import Union
from internal.data.data import DataRepo
from pandas import Series


import pprint


class ContentPipelineLatest:

    def __init__(self, doc: dict, repo: DataRepo):
        # assert doc["version"] == "V2.0"
        self.doc = doc
        self.content: dict = self.doc["content"]['impExpEntReport']
        self.repo = repo

    async def process(self):
        # pprint.pprint(self.content)
        tag_map = await self.gather_tag_dict()
        for row in self.content["customerDetail_12"]:
            row["QYBQ"] = tag_map[row["PURCHASER_NAME"]]

        for row in self.content["customerDetail_24"]:
            row["QYBQ"] = tag_map[row["PURCHASER_NAME"]]

        for row in self.content["supplierRanking_12"]:
            row["QYBQ"] = tag_map[row["SALES_NAME"]]

        for row in self.content["supplierRanking_24"]:
            row["QYBQ"] = tag_map[row["SALES_NAME"]]

        # pprint.pprint(self.content["customerDetail_12"][0])

    async def gather_tag_dict(self) -> dict:
        companies: list[str] = []
        companies += [row["PURCHASER_NAME"] for row in self.content["customerDetail_12"]]
        companies += [row["PURCHASER_NAME"] for row in self.content["customerDetail_24"]]
        companies += [row["SALES_NAME"] for row in self.content["supplierRanking_12"]]
        companies += [row["SALES_NAME"] for row in self.content["supplierRanking_24"]]
        s = Series(companies)
        s.drop_duplicates(inplace=True)
        tags = await asyncio.gather(*(self.repo.dw_data.get_tags_by_name(row) for row in s))
        return {name: tag for name, tag in tags}




