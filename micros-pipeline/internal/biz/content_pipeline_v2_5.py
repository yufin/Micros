import asyncio
from typing import Optional
import math
from logging import Logger
import numpy as np
import pandas as pd
from pandas import Series
import time
from internal.data.data import DataRepo
from typing import Union
from loguru import logger

pd.set_option('display.max_columns', None)
pd.set_option('display.max_rows', None)
pd.set_option('display.width', 1000)


class ContentPipelineV25:
    # 用于报文v1版本

    def __init__(
            self,
            logger: Logger,
            doc: dict,
            repo: DataRepo
    ):
        assert isinstance(doc, dict)
        assert doc["version"] == "V1.0"
        self.logger: Logger = logger
        self.doc = doc
        self.content: dict = self.doc["content"]['impExpEntReport']
        self.repo = repo
        self.__buff_df: dict[str, pd.DataFrame] = {}

    def _get_buff_df(self, key: str, use_buf: bool = True) -> pd.DataFrame:
        if not use_buf:
            return pd.DataFrame(self.content[key])
        if self.__buff_df.get(key) is None:
            self.__buff_df[key] = pd.DataFrame(self.content[key])
        return self.__buff_df[key].copy()

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

    def _summary_supplier_invoice(self):
        def __parse_established_date(d):
            try:
                dt = d.get("companyInfo", {}).get("establishDate", None)
                return dt
            except Exception:
                return None

        def __eval_decades(dt):
            try:
                return int(dt.year // 10 * 10)
            except Exception:
                return None

        _summary = []
        _df_sup12 = self._get_buff_df("supplierRanking_12", use_buf=False)
        _df_sup12 = _df_sup12[~_df_sup12['SALES_NAME'].isin(['合计', '其他'])]
        _df_sup12["RATIO_AMOUNT_TAX"] = _df_sup12["RATIO_AMOUNT_TAX"].apply(lambda x: self._to_numeric(x) / 100)
        sup12_desc_tax_ratio = _df_sup12.sort_values(by="RATIO_AMOUNT_TAX", ascending=False)
        top5_sum_ratio = sup12_desc_tax_ratio.head(5)["RATIO_AMOUNT_TAX"].sum()
        top20_sum_ratio = sup12_desc_tax_ratio.head(20)["RATIO_AMOUNT_TAX"].sum()

        if top5_sum_ratio < 0.3:
            desc12_top5 = "很低"
        elif 0.3 <= top5_sum_ratio < 0.5:
            desc12_top5 = "较低"
        elif 0.5 <= top5_sum_ratio < 0.7:
            desc12_top5 = "比较平均"
        elif 0.7 <= top5_sum_ratio < 0.9:
            desc12_top5 = "较高"
        else:
            desc12_top5 = "很高"

        _summary.append(
            f"近12个月前五大供应商合计占比{round(top5_sum_ratio * 100, 2)}%，前20供应商合计占比{round(top20_sum_ratio * 100, 2)}%，供应商集中度{desc12_top5}。"
        )

        stability = self._quick_parse("riskIndexes.INDEX_VALUE", ('INDEX_DEC', '供应商稳定性'), True)
        if stability is None:
            stability = 0
        stability *= 100
        if stability <= 10:
            stab_desc = "很低，新供应商数量很多"
        elif 10 < stability <= 30:
            stab_desc = "较低，新供应商数量较多"
        elif 30 < stability <= 50:
            stab_desc = "一般，新老供应商约一半"
        elif 50 < stability <= 70:
            stab_desc = "较高，新供应商数量较少"
        else:
            stab_desc = "很高，新供应商数量很少"
        _summary.append(
            f"近12个月与前12个月重要供应商金额占比重合程度比值为{round(stability, 2)}%,供应商稳定性{stab_desc}"
        )

        _df_cust24 = self._get_buff_df("supplierRanking_24")
        _df_cust24 = _df_cust24[~_df_cust24['SALES_NAME'].isin(['合计', '其他'])]
        _df_cust24["RATIO_AMOUNT_TAX"] = _df_cust24["RATIO_AMOUNT_TAX"].apply(lambda x: self._to_numeric(x) / 100)
        _df_cust24 = _df_cust24.sort_values(by="RATIO_AMOUNT_TAX", ascending=False)
        _df_cust24['rank'] = _df_cust24['RATIO_AMOUNT_TAX'].rank(method='first', ascending=False)
        if not _df_cust24.empty:
            for i in range(0, min(len(_df_cust24) - 1, 5)):
                s_cus12_top = sup12_desc_tax_ratio.iloc[i]
                s_cus12_top_at24 = _df_cust24[_df_cust24["SALES_NAME"] == s_cus12_top["SALES_NAME"]].iloc[0]
                annual_changed_ratio = (s_cus12_top["RATIO_AMOUNT_TAX"] - s_cus12_top_at24["RATIO_AMOUNT_TAX"]) / \
                                       s_cus12_top_at24["RATIO_AMOUNT_TAX"] * 100
                annual_changed_desc = ""
                if annual_changed_ratio > 0:
                    if 10 < annual_changed_ratio <= 30:
                        annual_changed_desc = "近期采购量小幅度提高"
                    elif 30 < annual_changed_ratio <= 100:
                        annual_changed_desc = "近期采购量大幅度提高"
                    elif 100 < annual_changed_ratio <= 170:
                        annual_changed_desc = "近期采购量约提高了1倍"
                    elif 170 < annual_changed_ratio:
                        x = int((annual_changed_ratio - 100) / 70) + 1
                        annual_changed_desc = f"近期采购量约提高了{x}倍"
                elif annual_changed_ratio < 0:
                    if -10 > annual_changed_ratio >= -40:
                        annual_changed_desc = "近期采购量小幅度下降"
                    elif -40 > annual_changed_ratio >= -60:
                        annual_changed_desc = "近期采购量减少约一半"
                    elif -60 > annual_changed_ratio >= -90:
                        annual_changed_desc = "近期采购量约大幅减少"
                    elif -90 > annual_changed_ratio >= -100:
                        annual_changed_desc = "近期基本没有采购"
                if annual_changed_desc == "":
                    annual_changed_desc = "近期采购量基本稳定"

                _summary.append(
                    f'近12个月第{i + 1}大供应商为{s_cus12_top["SALES_NAME"]}，采购占比为'
                    f'{round(s_cus12_top["RATIO_AMOUNT_TAX"] * 100, 2)}%，近24个月此供应商采购占比为'
                    f'{round(s_cus12_top_at24["RATIO_AMOUNT_TAX"] * 100, 2)}%，采购排名第'
                    f'{int(s_cus12_top_at24["rank"])}，近12个月内采购占比相较于24个月内变化'
                    f'{"+" if annual_changed_ratio > 0 else ""}{round(annual_changed_ratio, 2)}%。{annual_changed_desc}。'
                )
        _df_sup12 = _df_sup12.iloc[0:20]
        _df_sup12['established_at'] = _df_sup12["QYBQ"].apply(lambda x: __parse_established_date(x))
        _df_sup12['established_at'] = pd.to_datetime(_df_sup12['established_at'], errors="ignore")
        estab_stat = _df_sup12["established_at"].apply(lambda x: __eval_decades(x)).value_counts().sort_values(
            ascending=False)
        prop_stat = estab_stat / estab_stat.sum()
        if not estab_stat.empty:
            for i in range(0, len(estab_stat)):
                _summary.append(
                    f"近12个月前20大供应商中，成立时间在{int(estab_stat.index[i])}年代的供应商数量为{estab_stat.iloc[i]}家，占比{round(prop_stat.iloc[i] * 100, 1)}%。"
                )
        now_year = time.localtime().tm_year
        recent5_count = len(_df_sup12[_df_sup12["established_at"] > pd.to_datetime(f"{now_year - 5}-01")]),
        try:
            recent5_prop = recent5_count / estab_stat.values.sum()
        except Exception:
            recent5_prop = 0
        _summary[
            -1] += f"{'大部分供应商成立时间较久，经营稳定。' if 0 <= recent5_prop < 0.4 else '大部分供应商是近5年成立，经营稳定性偏弱。'}"

        self.content['purchaseDetailSummary'] = _summary

    def _summary_customer_invoice(self):
        def __parse_established_date(d):
            try:
                dt = d.get("companyInfo", {}).get("establishDate", None)
                return dt
            except Exception:
                return None

        def __eval_decades(dt):
            try:
                return int(dt.year // 10 * 10)
            except Exception:
                return None

        _summary = []
        _df_cust12 = self._get_buff_df("customerDetail_12")
        _df_cust12 = _df_cust12[~_df_cust12['PURCHASER_NAME'].isin(['合计', '其他'])]
        _df_cust12["RATIO_AMOUNT_TAX"] = _df_cust12["RATIO_AMOUNT_TAX"].apply(lambda x: self._to_numeric(x) / 100)
        cus12_desc_tax_ratio = _df_cust12.sort_values(by="RATIO_AMOUNT_TAX", ascending=False)
        top5_sum_ratio = cus12_desc_tax_ratio.head(5)["RATIO_AMOUNT_TAX"].sum()
        top20_sum_ratio = cus12_desc_tax_ratio.head(20)["RATIO_AMOUNT_TAX"].sum()

        if top5_sum_ratio < 0.3:
            desc12_top5 = "很低"
        elif 0.3 <= top5_sum_ratio < 0.5:
            desc12_top5 = "较低"
        elif 0.5 <= top5_sum_ratio < 0.7:
            desc12_top5 = "比较平均"
        elif 0.7 <= top5_sum_ratio < 0.9:
            desc12_top5 = "较高"
        else:
            desc12_top5 = "很高"

        _summary.append(
            f"近12个月前五大客户合计占比{round(top5_sum_ratio * 100, 2)}%，前20下游客户合计占比{round(top20_sum_ratio * 100, 2)}%，客户集中度{desc12_top5}。"
        )

        stability = self._quick_parse("riskIndexes.INDEX_VALUE", ('INDEX_DEC', '客户稳定性'), True)
        if stability is None:
            stability = 0
        stability *= 100
        if stability <= 10:
            stab_desc = "很低，新客户数量很多"
        elif 10 < stability <= 30:
            stab_desc = "较低，新客户数量较多"
        elif 30 < stability <= 50:
            stab_desc = "一般，新老客户约一半"
        elif 50 < stability <= 70:
            stab_desc = "较高，新客户数量较少"
        else:
            stab_desc = "很高，新客户数量很少"
        _summary.append(
            f"近12个月与前12个月重要客户金额占比重合程度比值为{round(stability, 2)}%,客户稳定性{stab_desc}"
        )

        _df_cust24 = self._get_buff_df("customerDetail_24")
        _df_cust24 = _df_cust24[~_df_cust24['PURCHASER_NAME'].isin(['合计', '其他'])]
        _df_cust24["RATIO_AMOUNT_TAX"] = _df_cust24["RATIO_AMOUNT_TAX"].apply(lambda x: self._to_numeric(x) / 100)
        _df_cust24 = _df_cust24.sort_values(by="RATIO_AMOUNT_TAX", ascending=False)
        _df_cust24['rank'] = _df_cust24['RATIO_AMOUNT_TAX'].rank(method='first', ascending=False)
        if not _df_cust24.empty:
            for i in range(0, min(len(_df_cust24) - 1, 5)):
                s_cus12_top = cus12_desc_tax_ratio.iloc[i]
                df_cus12_top_at24 = _df_cust24[_df_cust24["PURCHASER_NAME"] == s_cus12_top["PURCHASER_NAME"]]
                if df_cus12_top_at24.empty:
                    continue
                s_cus12_top_at24 = df_cus12_top_at24.iloc[0]
                annual_changed_ratio = (s_cus12_top["RATIO_AMOUNT_TAX"] - s_cus12_top_at24["RATIO_AMOUNT_TAX"]) / \
                                       s_cus12_top_at24["RATIO_AMOUNT_TAX"] * 100
                annual_changed_desc = ""
                if annual_changed_ratio > 0:
                    if 10 < annual_changed_ratio <= 30:
                        annual_changed_desc = "近期订单小幅度提高"
                    elif 30 < annual_changed_ratio <= 100:
                        annual_changed_desc = "近期订单大幅度提高"
                    elif 100 < annual_changed_ratio <= 170:
                        annual_changed_desc = "近期订单约提高了1倍"
                    elif 170 < annual_changed_ratio:
                        x = int((annual_changed_ratio - 100) / 70) + 1
                        annual_changed_desc = f"近期订单约提高了{x}倍"
                elif annual_changed_ratio < 0:
                    if -10 > annual_changed_ratio >= -40:
                        annual_changed_desc = "近期订单小幅度下降"
                    elif -40 > annual_changed_ratio >= -60:
                        annual_changed_desc = "近期订单减少约一半"
                    elif -60 > annual_changed_ratio >= -90:
                        annual_changed_desc = "近期订单约大幅减少"
                    elif -90 > annual_changed_ratio >= -100:
                        annual_changed_desc = "近期基本没有订单"
                if annual_changed_desc == "":
                    annual_changed_desc = "近期订单基本稳定"

                _summary.append(
                    f'近12个月第{i + 1}大客户为{s_cus12_top["PURCHASER_NAME"]}，销售占比为'
                    f'{round(s_cus12_top["RATIO_AMOUNT_TAX"]* 100, 2)}%，近24个月此客户销售占比为'
                    f'{round(s_cus12_top_at24["RATIO_AMOUNT_TAX"] * 100, 2)}%，销售排名第'
                    f'{int(s_cus12_top_at24["rank"])}，近12个月内销售占比相较于24个月内变化'
                    f'{"+" if annual_changed_ratio > 0 else ""}{round(annual_changed_ratio, 2)}%。{annual_changed_desc}。'
                )

        _df_cust12 = _df_cust12.iloc[0:20]
        _df_cust12['established_at'] = _df_cust12["QYBQ"].apply(lambda x: __parse_established_date(x))
        _df_cust12['established_at'] = pd.to_datetime(_df_cust12['established_at'], errors="ignore")
        estab_stat = _df_cust12["established_at"].apply(lambda x: __eval_decades(x)).value_counts().sort_values(
            ascending=False)
        prop_stat = estab_stat / estab_stat.sum()
        if not estab_stat.empty:
            for i in range(0, len(estab_stat)):
                _summary.append(
                    f"近12个月前20大客户中，成立时间在{int(estab_stat.index[i])}年代的客户数量为{estab_stat.iloc[i]}家，占比{round(prop_stat.iloc[i] * 100, 1)}%。"
                )
        now_year = time.localtime().tm_year
        recent5_count = len(_df_cust12[_df_cust12["established_at"] > pd.to_datetime(f"{now_year - 5}-01")]),
        try:
            if not estab_stat.empty:
                recent5_prop = recent5_count / estab_stat.values.sum()
            else:
                recent5_prop = 0
        except Exception:
            recent5_prop = 0
        _summary[
            -1] += f"{'大部分下游客户成立时间较久，经营稳定。' if 0 <= recent5_prop < 0.4 else '大部分下游客户是近5年成立，经营稳定性偏弱。'}"

        self.content['sellingDetailSummary'] = _summary

    def _selling_sta(self):
        # 4.5 主要销售商品分类
        _df_selling_sta = self._get_buff_df('sellingSTA', True)

        _selling_sta_total = self._quick_parse("sellingSTA.CGJE", ("SSSPDL", "合计"), True)

        _df_cus_detail_24 = self._get_buff_df('customerDetail_24')
        # 筛选'PURCHASER_NAME'长度小于等于4，或包含”税“这个字
        filtered_rows = (_df_cus_detail_24['PURCHASER_NAME'].str.len() <= 4) | _df_cus_detail_24[
            'PURCHASER_NAME'].str.contains("税")
        # 最后两行”其他“”合计“去除
        _person_or_tax_obj = _df_cus_detail_24.loc[filtered_rows, "SUM_AMOUNT_TAX"][:-2]
        _person_or_tax_total = sum([float(i.replace(',', '')) for i in list(_person_or_tax_obj)])

        # 把_cus_24_total重新赋值给合计
        _cus_24_total = self._quick_parse("customerDetail_24.SUM_AMOUNT_TAX", ("PURCHASER_NAME", "合计"), True)
        _df_selling_sta.loc[_df_selling_sta["SSSPDL"] == "合计", "CGJE"] = _cus_24_total
        # 新“其他”金额 = 原有的“其他”金额 + “个人”或"税局"金额
        _selling_others_total = self._quick_parse("sellingSTA.CGJE", ("SSSPDL", "其他"), True)
        _df_selling_sta.loc[_df_selling_sta["SSSPDL"] == "其他", "CGJE"] = _selling_others_total + _person_or_tax_total

        # 创建新行的数据
        new_row_data = {
            'SSSPDL': '出口',
            'CGJE': _cus_24_total - _selling_sta_total - _person_or_tax_total,
            'JEZB': '--',
            'SSSPZL': '未开票',
            'SSSPXL': '--'
        }

        # 创建包含新行数据的数据框
        new_row_df = pd.DataFrame([new_row_data])

        # 插入新行到指定位置
        insert_index = len(_df_selling_sta) - 2
        _df_selling_sta = pd.concat(
            [_df_selling_sta.iloc[:insert_index], new_row_df, _df_selling_sta.iloc[insert_index:]], ignore_index=True)

        # 重置索引
        _df_selling_sta.reset_index(drop=True, inplace=True)

        # 定义一个函数，将数字格式化为字符串并添加千位符
        def __format_str(money) -> str:
            if type(money) is str:
                return money
            else:
                return "{:,.2f}".format(money)  # 格式化为带千位符的浮点数字符串

        _df_selling_sta['CGJE'] = _df_selling_sta['CGJE'].apply(__format_str)

        def __count_radio(money, total) -> str:
            try:
                format_money = float(money.replace(",", ""))
                format_total = float(str(total).replace(",", ""))
                result = format_money / format_total * 100
                if result == 100:
                    return "100%"
                else:
                    return f"{result:.2f}%"
            except Exception:
                return "0%"

        _df_selling_sta['JEZB'] = _df_selling_sta.apply(lambda x: __count_radio(x['CGJE'], _cus_24_total), axis=1)

        self.content["sellingSTA"] = _df_selling_sta.to_dict(orient="records")

        # 5.5 主要采购商品
        _df_purchase_sta = self._get_buff_df('purchaseSTA')

        _sup_24_total = self._quick_parse("supplierRanking_24.SUM_AMOUNT_TAX", ("SALES_NAME", "合计"), True)
        _purchase_sta_total = self._quick_parse("purchaseSTA.CGJE", ("SSSPDL", "合计"), True)
        # 把_sup_24_total重新赋值给合计
        _df_purchase_sta.loc[_df_purchase_sta["SSSPDL"] == "合计", "CGJE"] = __format_str(_sup_24_total)
        # 新“其他”金额 = 原有的“其他”金额 + （_sup_24_total — 原来地合计数字）
        _sup_others_total = self._quick_parse("purchaseSTA.CGJE", ("SSSPDL", "其他"), True)
        _df_purchase_sta.loc[_df_purchase_sta["SSSPDL"] == "其他", "CGJE"] = __format_str(_sup_others_total + (
                _sup_24_total - _purchase_sta_total))

        # 修改“其他”的金额占比
        _df_purchase_sta['JEZB'] = _df_purchase_sta.apply(lambda x: __count_radio(x['CGJE'], _sup_24_total), axis=1)

        self.content["purchaseSTA"] = _df_purchase_sta.to_dict(orient="records")

    @logger.catch
    async def process(self):
        self.content["impJsonDate"] = self.doc["attribute_month"] + "-01"

        # 1.add calculate majorCommodityProportion for trades
        self._major_commodity_proportion("customerDetail_12")
        self._major_commodity_proportion("customerDetail_24")
        self._major_commodity_proportion("supplierRanking_12")
        self._major_commodity_proportion("supplierRanking_24")

        tag_map = await self._gather_tag_dict()
        iter_list = [
            ("supplierRanking_12", "SALES_NAME"),
            ("supplierRanking_24", "SALES_NAME"),
            ("customerDetail_12", "PURCHASER_NAME"),
            ("customerDetail_24", "PURCHASER_NAME")
        ]
        for k_item, k_name in iter_list:
            for row in self.content[k_item]:
                try:
                    row["QYBQ"] = tag_map.get(row[k_name])
                except Exception:
                    pass

        # 3.add capitalPaidIn field to businessInfo
        try:
            subj_d = tag_map.get(self.content["businessInfo"]["QYMC"])
            self.content["businessInfo"]["capitalPaidIn"] = None
            if subj_d is not None:
                ci = subj_d.get("companyInfo", {})
                if isinstance(ci, dict):
                    self.content["businessInfo"]["capitalPaidIn"] = ci.get("paidInCapital")
                self.content['subjectCompanyTags'] = {_k: _v for _k, _v in subj_d.items() if _k != "companyInfo"}
            else:
                pass
        except KeyError:
            pass

        if self.content.get('subjectCompanyTags') is None:
            self.content['subjectCompanyTags'] = {}
        self._subj_product_proportion()

        # 4.add fincal summary
        self._financial_summary()

        self._trading_summary()

        self._revenue_detail_summary()

        await self._subject_company_related_info()

        self._risk_indexes()

        self._summary_customer_invoice()

        self._summary_supplier_invoice()

        self._selling_sta()

    async def _subject_company_related_info(self):
        try:
            usc_id = self.content["businessInfo"]["TYSHXYDM"]
            assert type(usc_id) is str
        except Exception:
            return
        _, _related = await self.repo.dw_data.get_related(usc_id)
        self.content["shareholderData"] = _related.get("shareholder")
        self.content["branchesData"] = _related.get("branch")
        self.content["investmentData"] = _related.get("investment")
        self.content["equityTransparency"] = _related.get("equityTransparency")
        self.content["equityConclusion"] = _related.get("equityConclusion")

    def _subj_product_proportion(self):
        _df_selling_sta = self._get_buff_df("sellingSTA")
        _df_selling_sta = _df_selling_sta[~_df_selling_sta['SSSPDL'].isin(['合计', '其他'])]
        _df_selling_sta['category'] = _df_selling_sta['SSSPXL'].apply(lambda x: x.split('*')[1])
        _df_selling_sta['JEZB'] = _df_selling_sta['JEZB'].apply(lambda x: self._to_numeric(x)).astype(float)
        _df_selling_sta['proportion'] = _df_selling_sta.groupby('category')['JEZB'].transform('sum').astype(str) + '%'

        _df_selling_sta['ssspxl_last'] = _df_selling_sta['SSSPXL'].apply(lambda x: x.split('*')[-1])
        _df_selling_sta['jezb_bracket'] = '(' + _df_selling_sta['JEZB'].astype(str) + '%' + ')'
        _df_selling_sta = _df_selling_sta.sort_values(by='JEZB', ascending=False)
        _df_selling_sta['category_detail'] = _df_selling_sta.groupby('category')['ssspxl_last'].transform(
            lambda x: ','.join(x + _df_selling_sta.loc[x.index, 'jezb_bracket']))

        result = _df_selling_sta[['category', 'proportion', 'category_detail']].drop_duplicates()

        result['proportion_sort'] = result['proportion'].str.replace('%', '').astype(float)
        result = result.sort_values(by='proportion_sort', ascending=False).drop('proportion_sort', axis=1)
        result['proportion'] = result['proportion'].apply(lambda x: str(round(self._to_numeric(x), 2)) + '%')
        self.content['subjectCompanyTags']['productProportion'] = result.to_dict("records")

        # print(self.content['subjectCompanyTags'])

    def _revenue_detail_summary(self):
        # lrbDetail
        lx_year = [
            self._quick_parse("lrbDetail.SSNRQ", ('XM', '财务费用'), False),
            self._quick_parse("lrbDetail.SNRQ", ('XM', '财务费用'), False),
            self._quick_parse("lrbDetail.RQ", ('XM', '财务费用'), False),
        ]
        lx_value = [
            self._quick_parse("lrbDetail.Y2015", ('XM', '财务费用'), False),
            self._quick_parse("lrbDetail.Y2014", ('XM', '财务费用'), False),
            self._quick_parse("lrbDetail.Y2013", ('XM', '财务费用'), False),
        ]
        self.content['interestObj'] = {
            'year': lx_year,
            'value': lx_value,
        }

        _df_lrb = self._get_buff_df("lrbDetail")
        _df_lrb_interst = _df_lrb[_df_lrb["XM"] == "其中：利息费用"]
        _df_lrb_fincost = _df_lrb[_df_lrb["XM"] == "财务费用"]

        replace = False
        if not _df_lrb_interst.empty and not _df_lrb_fincost.empty:
            if _df_lrb_interst.iloc[0]['Y2015'] == '0':
                _df_lrb.loc[_df_lrb["XM"] == "其中：利息费用", 'Y2015'] = _df_lrb_fincost.iloc[0]['Y2015']
                replace = True
            if _df_lrb_interst.iloc[0]['Y2014'] == '0':
                _df_lrb.loc[_df_lrb["XM"] == "其中：利息费用", 'Y2014'] = _df_lrb_fincost.iloc[0]['Y2014']
                replace = True
            if _df_lrb_interst.iloc[0]['Y2013'] == '0':
                _df_lrb.loc[_df_lrb["XM"] == "其中：利息费用", 'Y2014'] = _df_lrb_fincost.iloc[0]['Y2013']
                replace = True

        if replace:
            self.content['lrbDetail'] = _df_lrb.to_dict(orient='records')
        # print(_df_lrb.to_dict(orient='records'))

    def _trading_summary(self):
        def __sub_rule(value1: float, value2: float) -> bool:
            value_ranges = [
                {"min": -math.inf, "max": -30},
                {"min": -30, "max": -10},
                {"min": -10, "max": 0},
                {"min": 0, "max": 10},
                {"min": 10, "max": 30},
                {"min": 30, "max": math.inf}
            ]
            for r in value_ranges:
                if r["min"] <= value1 < r["max"] and r["min"] <= value2 < r["max"]:
                    return True  # same range
            return False

        def __sub_summary_rule(value: float, time1: str, time2: str) -> str:
            _s = "-"
            if value > 30:
                _s = f"{time1}年-{time2}年利润较高，备货很少。"
            elif 10 < value <= 30:
                _s = f"{time1}年-{time2}年有部分利润，备货较少。"
            elif -10 <= value <= 10:
                _s = f"{time1}年-{time2}年基本上按销定采，供应链管理能力较好。"
            elif -30 <= value < -10:
                _s = f"{time1}年-{time2}年小幅度备货或者部分原材料呆滞，需关注库存风险。"
            elif value < -30:
                _s = f"{time1}年-{time2}年大幅度备货或大部分原材料呆滞，需重点关注库存风险。"
            return _s

        def _sub_summary_4(res, index, special_rule) -> str:
            _s4 = ''
            reset_val = res[index]['result']
            if special_rule == 1:
                _s4 = __sub_summary_rule(reset_val, res[2]['year'], res[0]['year'])
            elif special_rule == 2:
                _s4 = __sub_summary_rule(reset_val, res[2]['year'], res[1]['year'])
            elif special_rule == 3:
                _s4 = __sub_summary_rule(reset_val, res[1]['year'], res[0]['year'])
            elif special_rule == 4:
                if reset_val > 30:
                    _s4 = f"{res[index]['year']}年利润较高，备货很少。"
                elif 10 < reset_val <= 30:
                    _s4 = f"{res[index]['year']}年有部分利润，备货较少。"
                elif -10 <= reset_val <= 10:
                    _s4 = f"{res[index]['year']}年基本上按销定采，供应链管理能力较好。"
                elif -30 <= reset_val < -10:
                    _s4 = f"{res[index]['year']}年小幅度备货或者部分原材料呆滞，需关注库存风险。"
                elif reset_val < -30:
                    _s4 = f"{res[index]['year']}年大幅度备货或大部分原材料呆滞，需重点关注库存风险。"
            return _s4

        def _sub_summary_3(value) -> str:
            _s3 = ''
            if value > 30:
                _s3 = f"增加{value:.2f}%，快速增长"
            elif 10 < value <= 30:
                _s3 = f"增加{value:.2f}%，小幅增长"
            elif 0 <= value <= 10:
                _s3 = f"增加{value:.2f}%，趋势平稳"
            elif -10 <= value < 0:
                _s3 = f"减少{-value:.2f}%，趋势平稳"
            elif -30 <= value < -10:
                _s3 = f"减少{-value:.2f}%，小幅下降"
            elif value < -30:
                _s3 = f"减少{-value:.2f}%，大幅下降"
            return _s3

        def _sub_summary_2() -> str:
            info = self.content["ydcgqkInfo"]
            last_two_year_cg_time, last_year_cg_time, this_year_cg_time = "", "", ""
            last_two_year_cg_val, last_year_cg_val, this_year_cg_val = 0.0, 0.0, 0.0
            if 'TITLE3' in info and info['TITLE3'] != "" and info["TITLE3"] is not None:
                last_two_year_cg_time = info['TITLE1'][:4]
                last_year_cg_time = info['TITLE2'][:4]
                this_year_cg_time = info['TITLE3'][:4]
                last_two_year_cg_val = float(info['HJ_AVG1'].replace(",", "") if info['HJ_AVG1'] is not None else 0)
                last_year_cg_val = float(info['HJ_AVG2'].replace(",", "") if info['HJ_AVG2'] is not None else 0)
                this_year_cg_val = float(info['HJ_AVG3'].replace(",", "") if info['HJ_AVG3'] is not None else 0)
            else:
                last_year_cg_time = info['TITLE1'][:4]
                this_year_cg_time = info['TITLE2'][:4]
                last_year_cg_val = float(info['HJ_AVG1'].replace(",", "") if info['HJ_AVG1'] is not None else 0)
                this_year_cg_val = float(info['HJ_AVG2'].replace(",", "") if info['HJ_AVG2'] is not None else 0)
            cg_avg1, cg_avg2, cg_avg3 = 0.0, 0.0, 0.0
            if last_two_year_cg_val != 0:
                cg_avg1 = (last_year_cg_val - last_two_year_cg_val) / last_two_year_cg_val * 100
            cg_avg_desc1 = _sub_summary_3(cg_avg1)
            if last_year_cg_val != 0:
                cg_avg2 = (this_year_cg_val - last_year_cg_val) / last_year_cg_val * 100
            cg_avg_desc2 = _sub_summary_3(cg_avg2)
            if last_two_year_cg_val != 0:
                cg_avg3 = (this_year_cg_val - last_two_year_cg_val) / last_two_year_cg_val * 100
            cg_avg_desc3 = _sub_summary_3(cg_avg3)
            last_two_year_cg_val = round(last_two_year_cg_val)
            last_year_cg_val = round(last_year_cg_val)
            this_year_cg_val = round(this_year_cg_val)
            if last_two_year_cg_time != "":
                _s2 = f"{last_two_year_cg_time}年-{this_year_cg_time}年每月平均采购额分别是{last_two_year_cg_val}、{last_year_cg_val}、{this_year_cg_val}；其中{last_two_year_cg_time}年-{last_year_cg_time}年{cg_avg_desc1}，{last_year_cg_time}年-{this_year_cg_time}年{cg_avg_desc2}，整体来看，{last_two_year_cg_time}年-{this_year_cg_time}年{cg_avg_desc3}。"
            else:
                _s2 = f"{last_year_cg_time}年-{this_year_cg_time}年每月平均采购额分别是{last_year_cg_val}、{this_year_cg_val}；整体来看，{last_year_cg_time}年-{this_year_cg_time}年{cg_avg_desc2}。"
            return _s2

        def _sub_summary_1() -> str:
            info = self.content["ydxsqkInfo"]
            if 'TITLE_2' in info and info['TITLE_2'] == "年月-月均销售额":
                last_two_year_xs_time = ""
            else:
                last_two_year_xs_time = info['TITLE_2'][:4]
            last_year_xs_time = info['TITLE_1'][:4]
            this_year_xs_time = info['TITLE'][:4]
            sbavg2 = float(info['SB_AVG_2'].replace(",", "") if info["SB_AVG_2"] is not None else 0)
            sbavg1 = float(info['SB_AVG_1'].replace(",", "") if info["SB_AVG_1"] is not None else 0)
            sbavg = float(info['SB_AVG'].replace(",", "") if info["SB_AVG"] is not None else 0)
            if sbavg2 != 0:
                last_two_year_xs_val = sbavg2
            else:
                last_two_year_xs_val = float(
                    info['FP_XX_AVG_2'].replace(",", "") if info['FP_XX_AVG_2'] is not None else 0)
            if sbavg1 != 0:
                last_year_xs_val = sbavg1
            else:
                last_year_xs_val = float(info['FP_XX_AVG_1'].replace(",", "") if info['FP_XX_AVG_1'] is not None else 0)
            if sbavg != 0:
                this_year_xs_val = sbavg
            else:
                this_year_xs_val = float(info['FP_XX_AVG'].replace(",", "") if info['FP_XX_AVG'] is not None else 0)
            xs_avg1, xs_avg2, xs_avg3 = 0.0, 0.0, 0.0
            if last_two_year_xs_val != 0:
                xs_avg1 = (last_year_xs_val - last_two_year_xs_val) / last_two_year_xs_val * 100
            xs_avg_desc1 = _sub_summary_3(xs_avg1)
            if last_year_xs_val != 0:
                xs_avg2 = (this_year_xs_val - last_year_xs_val) / last_year_xs_val * 100
            xs_avg_desc2 = _sub_summary_3(xs_avg2)
            if last_two_year_xs_val != 0:
                xs_avg3 = (this_year_xs_val - last_two_year_xs_val) / last_two_year_xs_val * 100
            xs_avg_desc3 = _sub_summary_3(xs_avg3)
            last_two_year_xs_val = round(last_two_year_xs_val)
            last_year_xs_val = round(last_year_xs_val)
            this_year_xs_val = round(this_year_xs_val)
            if last_two_year_xs_time != "":
                _s1 = f"{last_two_year_xs_time}年-{this_year_xs_time}年每月平均销售额分别是{last_two_year_xs_val}、{last_year_xs_val}、{this_year_xs_val}；其中{last_two_year_xs_time}年-{last_year_xs_time}年{xs_avg_desc1}，{last_year_xs_time}年-{this_year_xs_time}年{xs_avg_desc2}，整体来看，{last_two_year_xs_time}年-{this_year_xs_time}年{xs_avg_desc3}。"
            else:
                _s1 = f"{last_year_xs_time}年-{this_year_xs_time}年每月平均销售额分别是{last_year_xs_val}、{this_year_xs_val}；整体来看，{last_year_xs_time}年-{this_year_xs_time}年{xs_avg_desc2}。"
            return _s1

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
        result = merged[['year', 'result']]
        # Sort by year in descending order
        result = result.sort_values(by='year', ascending=False)
        res = result.to_dict(orient='records')
        self.content["ydxsqkInfo"]["czxsPercentage"] = res

        xs_summary, cg_summary, last_two_year_summary, last_year_summary, this_year_summary = "", "", "", "", ""
        try:
            xs_summary = _sub_summary_1()
            cg_summary = _sub_summary_2()
        except Exception as e:
            raise e

        summaries1 = [
            xs_summary,
            cg_summary,
        ]

        if len(res) == 2:
            summary = f"{res[1]['year']}-{res[0]['year']}年销售增长率分别为{res[1]['result']:.2f}%, {res[0]['result']:.2f}%;"
        else:
            summary = f"{res[2]['year']}-{res[0]['year']}年销售增长率分别为{res[2]['result']:.2f}%, {res[1]['result']:.2f}%, {res[0]['result']:.2f}%;"

        summaries2 = [
            summary,
        ]

        if len(res) == 2:
            try:
                last_year_summary = _sub_summary_4(res, 1, 4)
                this_year_summary = _sub_summary_4(res, 0, 4)
            except Exception as e:
                raise e

            summaries3 = [
                last_two_year_summary,
                last_year_summary,
                this_year_summary,
            ]
        else:
            if __sub_rule(res[2]['result'], res[0]['result']):
                try:
                    sepcial_rule1 = _sub_summary_4(res, 0, 1)
                except Exception as e:
                    raise e

                summaries3 = [
                    sepcial_rule1,
                ]
            elif __sub_rule(res[2]['result'], res[1]['result']):
                try:
                    sepcial_rule2 = _sub_summary_4(res, 1, 2)
                    sepcial_rule4 = _sub_summary_4(res, 0, 4)
                except Exception as e:
                    raise e

                summaries3 = [
                    sepcial_rule2,
                    sepcial_rule4,
                ]
            elif __sub_rule(res[1]['result'], res[0]['result']):
                try:
                    sepcialRule3 = _sub_summary_4(res, 0, 3)
                    sepcial_rule4 = _sub_summary_4(res, 2, 4)
                except Exception as e:
                    raise e

                summaries3 = [
                    sepcial_rule4,
                    sepcialRule3,
                ]
            else:
                try:
                    last_two_year_summary = _sub_summary_4(res, 2, 4)
                    last_year_summary = _sub_summary_4(res, 1, 4)
                    this_year_summary = _sub_summary_4(res, 0, 4)
                except Exception as e:
                    raise e

                summaries3 = [
                    last_two_year_summary,
                    last_year_summary,
                    this_year_summary,
                ]

        self.content['ydxscgSummary'] = {
            "S1": summaries1,
            "S2": summaries2,
            "S3": summaries3,
        }

    def _financial_summary(self):
        self.content['lrbAnalysisSummary'] = {}
        _summary1 = ""
        _df_lrb = pd.DataFrame(self.content["lrbDetail"])
        _df_income = _df_lrb[_df_lrb["XM"] == "营业收入"]
        if not _df_income.empty:
            _s = _df_income.iloc[0]
            t_last = time.strptime(_s['RQ'], "%Y-%m-%d")
            _summary1 += (f"{_s['SSNRQ']}"
                          f"-{time.strftime('%Y年%m月%d日', t_last)}营业收入是"
                          f"{_s['Y2013'] if _s['Y2013'] is not None else '-'}万元, "
                          f"{_s['Y2014'] if _s['Y2014'] is not None else '-'}万元, "
                          f"{_s['Y2015'] if _s['Y2015'] is not None else '-'}万元；")
            # t_last.
            last_income = self._to_numeric(_s['Y2015'])
            senior_income = self._to_numeric(_s['Y2014'])
            junior_income = self._to_numeric(_s['Y2013'])
            avg_income_last = last_income / t_last.tm_yday * 30 if last_income is not None else None
            avg_income_senior = senior_income / 12 if senior_income is not None else None
            avg_income_junior = junior_income / 12 if junior_income is not None else None
            avg_ratio_latest = avg_income_last / avg_income_senior - 1. if avg_income_last is not None and avg_income_senior is not None and avg_income_senior != 0 else None
            avg_ratio_senior = avg_income_senior / avg_income_junior - 1. if avg_income_senior is not None and avg_income_junior is not None and avg_income_junior != 0 else None
            avg_ratio_latest_str = f'{avg_ratio_latest * 100:.2f}%' if avg_ratio_latest is not None else 'N/A'
            avg_ratio_senior_str = f'{avg_ratio_senior * 100:.2f}%' if avg_ratio_senior is not None else 'N/A'
            # Add descriptions
            if avg_ratio_latest is not None:
                avg_ratio_latest_str = f'上升{avg_ratio_latest_str}' if avg_ratio_latest >= 0 else f'下降{avg_ratio_latest_str}'
            if avg_ratio_senior is not None:
                avg_ratio_senior_str = f'上升{avg_ratio_senior_str}' if avg_ratio_senior >= 0 else f'上升{avg_ratio_senior_str}'
            if avg_ratio_latest is not None:
                _summary1 += (f"其中{t_last.tm_year}年年初至{time.strftime('%m月%d日', t_last)}月均收入"
                              f"{_s['Y2015']}，同比"
                              f"{avg_ratio_latest_str}，"
                              f"{self._pct_range_desc(avg_ratio_latest * 100)}。")
            if avg_ratio_senior is not None:
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
        receivable_turnover_days_2 = self._quick_parse("finIndexes.INDEX_VALUE_1", ('INDEX_TITLE', '应收款周转日'),
                                                       True)
        receivable_turnover_days_3 = self._quick_parse("finIndexes.INDEX_VALUE_2", ('INDEX_TITLE', '应收款周转日'),
                                                       True)
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

        _summary1 = _summary1.replace("None", "-")
        _summary2 = _summary2.replace("None", "-")
        _summary3 = _summary3.replace("None", "-")
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
        if _df_grouped.empty:
            return
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
        _df_trades.fillna(0, inplace=True)
        self.content[key_trades] = _df_trades.to_dict(orient='records')

    @staticmethod
    def _to_numeric(x: str) -> Optional[float]:
        try:
            return float(x.replace(",", "").replace("%", ""))
        except Exception:
            return None

    async def _gather_tag_dict(self) -> dict:
        companies: list[str] = [self.content["businessInfo"]["QYMC"]]

        k_list = [
            ("supplierRanking_12", "SALES_NAME"),
            ("supplierRanking_24", "SALES_NAME"),
            ("customerDetail_12", "PURCHASER_NAME"),
            ("customerDetail_24", "PURCHASER_NAME")
        ]
        for name_k, item_k in k_list:
            try:
                for row in self.content[name_k]:
                    companies.append(row[item_k])
            except TypeError:
                pass
        s = Series(companies)
        s.drop_duplicates(inplace=True)

        tags = await asyncio.gather(*(self.repo.dw_data.get_tags_by_name(row, True) for row in s))
        return {name: tag for name, tag in tags}

    def _risk_indexes(self):
        def __lambda_reformat_index_value(template, value, value_type) -> str:
            if value == "0" or value == "-":
                return "无"
            try:
                new_value = float(value)
            except ValueError:
                return value

            new_value = round(new_value * 100)
            if value_type == "percent":
                return template % new_value
            elif value_type == "radio":
                return template % __radio_summary(new_value)
            else:
                return template % value

        def __radio_summary(value) -> str:
            if value <= 10:
                return "很低"
            elif 10 < value <= 30:
                return "较低"
            elif 30 < value <= 50:
                return "一般"
            elif 50 < value <= 70:
                return "较高"
            else:
                return "很高"

        def __day_format(value) -> str:
            try:
                new_value = float(value)
            except ValueError:
                return value

            return f"{int(round(new_value))}天"

        def __lambda_summary(desc: str, val: str):
            # desc = item.get("INDEX_DEC", "")
            # val = item.get("INDEX_VALUE", "")
            reset_val = ""

            # 先定义一个名称映射字典
            name_mapping = {
                "历史变更-企业名称变更": "企业名称变更",
                "历史变更-地址变更": "地址变更",
                "变更的风险提示\n（减资）": "是否减资",
            }

            # 如果 desc 在映射字典中，将其映射为新名称
            if desc in name_mapping:
                desc = name_mapping[desc]

            if desc == "企业名称变更":
                reset_val = __lambda_reformat_index_value("变更企业名称%s次", val, "nopercent")
            elif desc == "地址变更":
                reset_val = __lambda_reformat_index_value("变更地址%s次", val, "nopercent")
            elif desc == "是否减资":
                try:
                    v = int(val)
                    if v > 0:
                        reset_val = __lambda_reformat_index_value("增资注册资本%s万", val, "nopercent")
                    elif v < 0:
                        reset_val = __lambda_reformat_index_value("减资注册资本%s万", val, "nopercent")

                except ValueError:
                    reset_val = "无"

            elif desc == "变更的风险提示\n（变更频繁）":
                reset_val = __lambda_reformat_index_value("变更%s次", val, "nopercent")
            elif desc == "金融欠款纠纷":
                reset_val = __lambda_reformat_index_value("近5年内企业和法人作为被告%s次", val, "nopercent")
            elif desc == "作为被告裁判文书涉案金额":
                reset_val = __lambda_reformat_index_value(
                    "%s（近5年企业民事裁判文书中作为被告涉案金额加总+近5年法人民事裁判文书中作为被告涉案金额加总）/（上年）营业收入",
                    val, "nopercent")
            elif desc == "作为原告裁判文书涉案金额":
                reset_val = __lambda_reformat_index_value(
                    "%s（近5年企业民事裁判文书中作为原告涉案金额加总+近5年法人民事裁判文书中作为原告涉案金额加总）/（上年）营业收入",
                    val, "nopercent")
            elif desc == "法人历史失信记录":
                reset_val = __lambda_reformat_index_value("近5年法人失信%s次", val, "nopercent")
            elif desc == "历史被执行人记录":
                reset_val = __lambda_reformat_index_value("近5年企业被执行%s次", val, "nopercent")
            elif desc == "工商处罚记录":
                reset_val = __lambda_reformat_index_value("近5年企业工商处罚%s次", val, "nopercent")
            elif desc == "收入成长率":
                reset_val = __lambda_reformat_index_value("%.f%%", val, "percent")
            elif desc == "毛利润成长率":
                reset_val = __lambda_reformat_index_value("%.f%%", val, "percent")
            elif desc == "供应商评价":
                reset_val = __lambda_reformat_index_value("%.f%%", val, "percent")
            elif desc == "供应商稳定性":
                reset_val = __lambda_reformat_index_value("%s", val, "radio")
            elif desc == "客户评价":
                reset_val = __lambda_reformat_index_value("%.f%%", val, "percent")
            elif desc == "客户稳定性":
                reset_val = __lambda_reformat_index_value("%s", val, "radio")
            elif desc == "主营业务专注度":
                reset_val = __lambda_reformat_index_value("%.f%%", val, "percent")
            elif desc == "净利与毛利波动差异":
                reset_val = __lambda_reformat_index_value("%s%%", val, "nopercent")
            elif desc == "现金比率":
                reset_val = __lambda_reformat_index_value("%.f%%", val, "percent")
            elif desc == "应收运营周转天数":
                reset_val = __day_format(val)
            elif desc == "应付运营周转天数":
                reset_val = __day_format(val)
            elif desc == "存货周转天数":
                reset_val = __day_format(val)
            else:
                reset_val = '无'
            return desc, reset_val

        _df_risk_idx = self._get_buff_df('riskIndexes')
        _df_risk_idx['INDEX_DEC'], _df_risk_idx['INDEX_VALUE'] = zip(
            *_df_risk_idx.apply(lambda x: __lambda_summary(x['INDEX_DEC'], x['INDEX_VALUE']), axis=1))

        self.content["riskIndexes"] = _df_risk_idx.to_dict("records")

        count = _df_risk_idx['INDEX_FLAG'].value_counts()
        total = count.sum()

        # 使用 get 方法获取各个值的计数，如果计数为 None 则设置为 0
        normal = count.get("正常", 0)
        attention = count.get("关注", 0)
        ordinary = count.get("普通", 0)
        abnormal = count.get("异常", 0)

        # 使用列表推导获取异常的 INDEX_DEC 值，如果没有异常则为空列表
        abnormal_tags = [item['INDEX_DEC'] for item in _df_risk_idx.to_dict(orient="records") if
                         item['INDEX_FLAG'] == "异常"]

        self.content['riskIndexesSummary'] = {
            "total": int(total),
            "normal": int(normal),
            "attention": int(attention),
            "ordinary": int(ordinary),
            "abnormal": int(abnormal),
            "abnormalTags": abnormal_tags
        }
