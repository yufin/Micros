import asyncio
from pprint import pprint
from typing import Optional

import numpy as np
import pandas as pd
from pandas import Series
import time
from internal.data.data import DataRepo
from typing import Union

pd.set_option('display.max_columns', None)
pd.set_option('display.max_rows', None)


class ContentPipelineLatest:
    def __init__(self, doc: dict, repo: DataRepo):
        # assert doc["version"] == "V2.0"
        self.doc = doc
        self.content: dict = self.doc["content"]['impExpEntReport']
        self.repo = repo

        self.__buff_df = {}

    def _get_buff_df(self, key: str):
        if self.__buff_df.get(key) is None:
            self.__buff_df[key] = pd.DataFrame(self.content[key])
        return self.__buff_df[key]

    def _quick_parse(self, key: str, cond: tuple, to_float: bool = False) -> Union[str, float, None]:
        keys = key.split(".")
        assert len(keys) == 2
        assert len(cond) == 2
        _df = self._get_buff_df(keys[0])
        _df_filtered = _df[_df[cond[0]] == cond[1]]
        if _df_filtered.empty:
            return None
        _df_filtered = _df_filtered.iloc[0]
        if to_float:
            if _df_filtered[keys[1]] is None:
                return None
            return self._to_numeric(_df_filtered[keys[1]])
        else:
            return _df_filtered[keys[1]]

    async def process(self):
        pprint(self.content)
        print("===================================== \n")

        # 1.add calculate majorCommodityProportion for trades
        self._major_commodity_proportion("customerDetail_12")
        self._major_commodity_proportion("customerDetail_24")
        self._major_commodity_proportion("supplierRanking_12")
        self._major_commodity_proportion("supplierRanking_24")

        # 2.add tag for companies
        tag_map = await self._gather_tag_dict()
        for row in self.content["customerDetail_12"]:
            row["QYBQ"] = tag_map[row["PURCHASER_NAME"]]
        for row in self.content["customerDetail_24"]:
            row["QYBQ"] = tag_map[row["PURCHASER_NAME"]]
        for row in self.content["supplierRanking_12"]:
            row["QYBQ"] = tag_map[row["SALES_NAME"]]
        for row in self.content["supplierRanking_24"]:
            row["QYBQ"] = tag_map[row["SALES_NAME"]]

        # 3.add capitalPaidIn field to businessInfo
        try:
            subj_d = tag_map[self.content["businessInfo"]["QYMC"]]
            self.content["businessInfo"]["capitalPaidIn"] = None
            if subj_d is not None:
                self.content["businessInfo"]["capitalPaidIn"] = subj_d["companyInfo"]["paidInCapital"]
            else:
                pass
        except KeyError:
            pass

        # 4.add fincal summary
        self._financial_summary()

        self._trading_summary()

    def _trading_summary(self):
        _df_sales = pd.DataFrame(self.content["ydxsqkDetail"])
        _df_purchase = pd.DataFrame(self.content["ydcgqkSTA"])
        _df_sales['year'] = _df_sales['MONTH'].str.slice(0, 4)
        _df_purchase['year'] = _df_purchase['LAST_24M'].str.slice(0, 4)
        _df_sales["SBKJZSR"] = _df_sales["SBKJZSR"].apply(lambda x: self._to_numeric(x))
        _df_sales["FPKJSR"] = _df_sales["FPKJSR"].apply(lambda x: self._to_numeric(x))
        _df_purchase['HJ_M'] = _df_purchase['HJ_M'].apply(lambda x: self._to_numeric(x))
        # Group by content_id and year and calculate the sum
        grouped_sales = _df_sales.groupby(['year']).sum().reset_index()
        grouped_purchase = _df_purchase.groupby(['year']).sum().reset_index()
        # Merge the two dataframes
        merged = pd.merge(grouped_sales, grouped_purchase, on=['year'], suffixes=('_sales', '_purchase'))
        # Calculate the result column using the condition
        merged['result'] = np.where(
            merged['SBKJZSR'] == 0,
            (merged['FPKJSR'] - merged['HJ_M']) / merged['FPKJSR'] * 100,
            (merged['SBKJZSR'] - merged['HJ_M']) / merged['SBKJZSR'] * 100
        )
        # Select the necessary columns
        result = merged[[ 'year', 'result']]
        # Sort by year in descending order
        result = result.sort_values(by='year', ascending=False)
        print(result)

    def _financial_summary(self):
        self.content['lrbAnalysisSummary'] = {}

        _summary1 = ""
        _df_lrb = pd.DataFrame(self.content["lrbDetail"])
        _df_income = _df_lrb[_df_lrb["XM"] == "营业收入"]
        if not _df_income.empty:
            _s = _df_income.iloc[0]
            t_last = time.strptime(_s['RQ'], "%Y-%m-%d")
            _summary1 += f"{_s['SSNRQ']}-{time.strftime('%Y年%m月%d日', t_last)}营业收入是{_s['Y2013']}万元, {_s['Y2014']}万元, {_s['Y2015']}万元；"
            # t_last.
            last_income = self._to_numeric(_s['Y2015'])
            senior_income = self._to_numeric(_s['Y2014'])
            junior_income = self._to_numeric(_s['Y2013'])
            avg_income_last = last_income / t_last.tm_yday * 30 if last_income is not None else None
            avg_income_senior = senior_income / 12 if senior_income is not None else None
            avg_income_junior = junior_income / 12 if junior_income is not None else None
            avg_ratio_latest = avg_income_last / avg_income_senior - 1. if avg_income_last is not None and avg_income_senior is not None else None
            avg_ratio_senior = avg_income_senior / avg_income_junior - 1. if avg_income_senior is not None and avg_income_junior is not None else None
            avg_ratio_latest_str = f'{avg_ratio_latest * 100:.2f}%' if avg_ratio_latest is not None else 'N/A'
            avg_ratio_senior_str = f'{avg_ratio_senior * 100:.2f}%' if avg_ratio_senior is not None else 'N/A'
            # Add descriptions
            if avg_ratio_latest is not None:
                avg_ratio_latest_str = f'上升{avg_ratio_latest_str}' if avg_ratio_latest >= 0 else f'下降{avg_ratio_latest_str}'
            if avg_ratio_senior is not None:
                avg_ratio_senior_str = f'上升{avg_ratio_senior_str}' if avg_ratio_senior >= 0 else f'上升{avg_ratio_senior_str}'
            _summary1 += f"其中{t_last.tm_year}年年初至{time.strftime('%m月%d日', t_last)}月均收入{_s['Y2015']}，同比{avg_ratio_latest_str}，{self._pct_range_desc(avg_ratio_latest * 100)}。"
            _summary1 += f"{_s['SNRQ']}全年收入同比{avg_ratio_senior_str}, {self._pct_range_desc(avg_ratio_senior * 100)}."
        self.content['lrbAnalysisSummary']['summary1'] = _summary1

        _summary2 = ""
        _summary3 = ""

        t1 = self._quick_parse("zcfzbDetail.SSNRQ", ('XM', '所有者权益（或股东权益）合计'))
        t2 = self._quick_parse("zcfzbDetail.SNRQ", ('XM', '所有者权益（或股东权益）合计'))
        t3 = self._quick_parse("zcfzbDetail.RQ", ('XM', '所有者权益（或股东权益）合计'))
        equity_1 = self._quick_parse("zcfzbDetail.NUM2", ('XM', '所有者权益（或股东权益）合计'), True)
        equity_2 = self._quick_parse("zcfzbDetail.NUM1", ('XM', '所有者权益（或股东权益）合计'), True)
        equity_3 = self._quick_parse("zcfzbDetail.NUM0", ('XM', '所有者权益（或股东权益）合计'), True)
        debt_ratio_1 = self._quick_parse("finIndexes.INDEX_VALUE", ('INDEX_TITLE', '资产负债率'), True)
        debt_ratio_2 = self._quick_parse("finIndexes.INDEX_VALUE_1", ('INDEX_TITLE', '资产负债率'), True)
        debt_ratio_3 = self._quick_parse("finIndexes.INDEX_VALUE_2", ('INDEX_TITLE', '资产负债率'), True)
        liquidity_ratio_1 = self._quick_parse("finIndexes.INDEX_VALUE", ('INDEX_TITLE', '流动比例'), True)
        liquidity_ratio_2 = self._quick_parse("finIndexes.INDEX_VALUE_1", ('INDEX_TITLE', '流动比例'), True)
        liquidity_ratio_3 = self._quick_parse("finIndexes.INDEX_VALUE_2", ('INDEX_TITLE', '流动比例'), True)
        quick_ratio_1 = self._quick_parse("finIndexes.INDEX_VALUE", ('INDEX_TITLE', '速动比例'), True)
        quick_ratio_2 = self._quick_parse("finIndexes.INDEX_VALUE_1", ('INDEX_TITLE', '速动比例'), True)
        quick_ratio_3 = self._quick_parse("finIndexes.INDEX_VALUE_2", ('INDEX_TITLE', '速动比例'), True)
        receivable_turnover_days_1 = self._quick_parse("finIndexes.INDEX_VALUE", ('INDEX_TITLE', '应收款周转日'), True)
        receivable_turnover_days_2 = self._quick_parse("finIndexes.INDEX_VALUE_1", ('INDEX_TITLE', '应收款周转日'),True)
        receivable_turnover_days_3 = self._quick_parse("finIndexes.INDEX_VALUE_2", ('INDEX_TITLE', '应收款周转日'),True)
        asset_total_1 = self._quick_parse("zcfzbDetail.NUM2", ('XM', '资产合计'), True)
        # asset_total_2 = self._quick_parse("zcfzbDetail.NUM1", ('XM', '资产合计'), True)
        # asset_total_3 = self._quick_parse("zcfzbDetail.NUM0", ('XM', '资产合计'), True)
        receivable_1 = self._quick_parse("zcfzbDetail.NUM2", ('XM', '应收账款'), True)
        # receivable_2 = self._quick_parse("zcfzbDetail.NUM1", ('XM', '应收账款'), True)
        # receivable_3 = self._quick_parse("zcfzbDetail.NUM0", ('XM', '应收账款'), True)
        storage_turnover_days_1 = self._quick_parse("finIndexes.INDEX_VALUE", ('INDEX_TITLE', '库存周转日'), True)
        storage_turnover_days_2 = self._quick_parse("finIndexes.INDEX_VALUE_1", ('INDEX_TITLE', '库存周转日'), True)
        storage_turnover_days_3 = self._quick_parse("finIndexes.INDEX_VALUE_2", ('INDEX_TITLE', '库存周转日'), True)
        storage_1 = self._quick_parse("zcfzbDetail.NUM2", ('XM', '存货'), True)
        # storage_2 = self._quick_parse("zcfzbDetail.NUM1", ('XM', '存货'), True)
        storage_3 = self._quick_parse("zcfzbDetail.NUM0", ('XM', '存货'), True)


        _summary2 += f"{t1}至{t3}净资产是{equity_1}万元，{equity_2}万元，{equity_3}万元; "

        _summary2 += (
            f"{t1}至{t3}资产负债率为{debt_ratio_1}，{debt_ratio_2}，{debt_ratio_3}，"
            f"负债率{self.rate_summary(debt_ratio_1, '负债率')}，"
            f"{self.rate_summary(debt_ratio_2, '负债率')}，"
            f"{self.rate_summary(debt_ratio_3, '负债率')}；,")
        #
        _summary2 += f"{t1}至{t3}流动比例为{liquidity_ratio_1}，{liquidity_ratio_2}，{liquidity_ratio_3}，"
        #
        _summary2 += f"{t1}至{t3}速动比例为{quick_ratio_2}，{quick_ratio_2}，{quick_ratio_3}，"
        _summary2 += f"{self.rate_summary(quick_ratio_3, '速动比例')}；"

        _summary3 += f"{t1}至{t3}应收账款周转日为{receivable_turnover_days_1}，{receivable_turnover_days_2}，{receivable_turnover_days_3}，"

        try:
            rate_rec = str(round(receivable_1 / asset_total_1 * 100, 1)) + "%"
        except TypeError:
            rate_rec = "-"

        try:
            rate_stor = str(round(storage_1 / asset_total_1 * 100, 1)) + "%"
        except TypeError:
            rate_stor = "-"

        _summary3 += f"{t1}应收款周转速度{self.rate_summary(receivable_turnover_days_1, '应收款周转日')}，应收款金额为{receivable_1}占总资产约{rate_rec}；"
        # %s - %s存货周转天数分别是 % s， %s， %s； %s存货周转速度 % s，存货金额为 % s，占总资产约 % .2f %%；近期公司运营能力 % s。
        _summary3 += f"{t1}至{t3}存货周转天数分别是{storage_turnover_days_1}，{storage_turnover_days_2}，{storage_turnover_days_3}；"
        _summary3 += f"{t3}存货周转速度{self.rate_summary(storage_turnover_days_1, '库存周转日')}，存货金额为{storage_3}，占总资产约{rate_stor}；"

        try:
            _temp_ysk = self.rate_summary(receivable_turnover_days_1, '应收款周转日')
            _temp_ch = self.rate_summary(storage_turnover_days_1, '库存周转日')
            if _temp_ysk == "-" and _temp_ch == "-":
                op_desc = "-"
            elif _temp_ysk == "较快" and _temp_ch == "较快":
                op_desc = "较强"
            elif (_temp_ysk == "较慢" and _temp_ch == "较慢") or (_temp_ysk == "很慢" and _temp_ch == "很慢"):
                op_desc = "较弱"
            else:
                op_desc = "一般"
        except Exception:
            op_desc = "-"

        _summary3 += f"近期公司运营能力{op_desc}。"

        self.content['zcfzbAnalysisSummary'] = {
            "summary1": _summary1,
            "summary2": _summary2,
            "summary3": _summary3
        }

    @staticmethod
    def rate_summary(s: str, summary_type: str) -> str:
        try:
            s = float(s)
        except ValueError:
            return "-"
        except TypeError:
            return "-"

        if summary_type == "负债率":
            s *= 100
            if 70 < s <= 100:
                return "较高"
            elif 50 < s <= 70:
                return "居中"
            elif s < 50:
                return "较低"

        elif summary_type == "速动比例":
            if s <= 1:
                return "近期公司短期偿债能力较弱"
            elif s > 1:
                return "近期公司短期偿债能力较强"

        elif summary_type in ["应收款周转日", "库存周转日"]:
            if s > 90:
                return "很慢"
            elif 40 < s <= 90:
                return "较慢"
            elif 0 <= s <= 40:
                return "较快"

        return "-"

    @staticmethod
    def _pct_range_desc(pct):
        if 10 < pct <= 30:
            return "小幅度增长"
        elif 30 < pct <= 100:
            return "大幅度增长"
        elif 100 < pct <= 170:
            return "约增长1倍"
        elif 170 < pct:
            x = int((pct - 70) / 100) + 1
            return f"约增长{x}倍"
        elif -10 > pct >= -40:
            return "小幅度下降"
        elif -40 > pct >= -60:
            return "约下降一半"
        elif -60 > pct > -100:
            return "大幅度下滑"
        else:
            return "几乎持平"

    def _major_commodity_proportion(self, key_trades: str):
        _df_trades = pd.DataFrame(self.content[key_trades])
        _df_trades['_major_commodity'] = _df_trades["COMMODITY_NAME"].apply(
            lambda x: x.split("*")[1] if str(x).__contains__("*") else None)
        _df_trades['_sum_amount_tax']: float = _df_trades["SUM_AMOUNT_TAX"].apply(lambda x: self._to_numeric(x))
        _df_trades['_commodity_ratio']: float = _df_trades["COMMODITY_RATIO"].apply(lambda x: self._to_numeric(x))

        _df_grouped = _df_trades.groupby(['_major_commodity'])['_sum_amount_tax'].sum().reset_index()
        _df_grouped = _df_grouped.loc[_df_grouped['_sum_amount_tax'].idxmax()].to_frame().T
        major_commodity, sum_amount_tax = _df_grouped['_major_commodity'].values[0], \
        _df_grouped['_sum_amount_tax'].values[0]

        _df_trades['major_commodity_proportion'] = _df_trades.apply(
            lambda x: x['_commodity_ratio'] / 100 * x['_sum_amount_tax'] / sum_amount_tax
            if x['_major_commodity'] == major_commodity else None, axis=1
        )

        del _df_trades['_major_commodity']
        del _df_trades['_sum_amount_tax']
        del _df_trades['_commodity_ratio']
        self.content[key_trades] = _df_trades.to_dict(orient='records')

    @staticmethod
    def _to_numeric(x: str) -> Optional[float]:
        try:
            return float(x.replace(",", "").replace("%", ""))
        except ValueError:
            return None

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


if __name__ == '__main__':
    td = {}

    print(td.get('a'))
