import asyncio
import datetime
import math
import threading
import time
from logging import Logger
from typing import Optional
from typing import Union

import numpy as np
import pandas as pd
from dateutil.relativedelta import relativedelta

from api.pipeline.v1 import pipeline_pb2
from internal.data.data import DataRepo
from internal.biz.dw_stat import DwStat

from pandas._libs.tslibs.parsing import DateParseError
pd.set_option('display.max_columns', None)
pd.set_option('display.max_rows', None)
pd.set_option('display.width', 1000)


class ContentPipelineV3:
    # 用于报文v2.0版本
    def __init__(
            self,
            logger: Logger,
            doc: dict,
            repo: DataRepo
    ):
        assert isinstance(doc, dict)
        assert doc["version"] == "V2.0"
        self.logger: Logger = logger
        self.doc = doc
        self.content: dict = self.doc["content"]['impExpEntReport']
        self.repo = repo
        self.__buff_df: dict[str, pd.DataFrame] = {}
        self.__time_point: Optional[datetime.datetime] = None
        self.__lock: threading.Lock = threading.Lock()
        self.__dw_stat: DwStat = DwStat(
            motor_client=repo.mongo_db,
            dw_data_client_v3=repo.dw_data_v3,
            logger=logger
        )

    def reformat(self):
        try:
            _df = pd.DataFrame(self.content["ydxsqkDetail"])
        except Exception:
            return
        try:
            _ = pd.to_datetime(_df["MONTH"])
        except DateParseError:
            _df["MONTH"] = _df["MONTH"].apply(lambda x: x[:7] if isinstance(x, str) and len(x) == 9 else x)
            self.content['ydxsqkDetail'] = _df.to_dict("records")
        return


    @property
    def tag_t0(self) -> Optional[str]:
        try:
            t = self.content['ssfmxxSum']['bnsj'][0]['NY']
            return t
        except Exception:
            return None

    @property
    def tag_t1(self) -> Optional[str]:
        try:
            t = self.content['ssfmxxSum']['snsj'][0]['NY']
            return t
        except Exception:
            return None

    @property
    def tag_t2(self) -> Optional[str]:
        try:
            t = self.content['ssfmxxSum']['ssnsj'][0]['NY']
            return t
        except Exception:
            return None

    @property
    def tag_t3(self) -> Optional[str]:
        try:
            t = self.content['ssfmxxSum']['sssnsj'][0]['NY']
            return t
        except Exception:
            return None

    @property
    def time_point(self) -> datetime.datetime:
        if self.__time_point is None:
            t_str = self.doc["attribute_month"] + "-01"
            self.__time_point = pd.to_datetime(t_str).to_pydatetime()
        return self.__time_point

    async def get_tags_by_name(self, name: str, tp: datetime.datetime) -> tuple[str, Optional[dict]]:
        self.logger.info("get_tags_by_name name=" + name)

        ident = await self.repo.dw_data_v3.get_usc_id_by_name(name)

        if ident and ident.success:
            req = self.repo.dw_data_v3.new_pb_req(
                usc_id=ident.data.usc_id,
                time_point=tp
            )
            (info, ranking, credential, equity_conclusion, case_stat,
             case_regis_stat, executive_stat, his_executive_stat, financing_stat) \
                = await asyncio.gather(
                self.repo.dw_data_v3.get_ent_info(req),
                self.repo.dw_data_v3.get_ent_ranking_list(req),
                self.repo.dw_data_v3.get_ent_credential(req),
                self.repo.dw_data_v3.get_ent_equity_transparency_conclusion(req),
                self.__dw_stat.stat_case(req),
                self.__dw_stat.stat_case_registration(req),
                self.__dw_stat.stat_executive(req),
                self.__dw_stat.stat_his_executive(req),
                self.__dw_stat.stat_financing(req),
            )

            resp = {
                "companyInfo": info.get("data") if info else None,
                "authorizedTag": credential.get("data") if credential else None,
                "rankingTag": ranking.get("data") if ranking else None,
                "equityConclusion": equity_conclusion if equity_conclusion else None,
                "caseStat": case_stat,
                "caseRegisStat": case_regis_stat,
                "executiveStat": executive_stat,
                "hisExecutiveStat": his_executive_stat,
                "financingStat": financing_stat,
            }
            return name, resp
        else:
            return name, None

    async def sales_detail_chart(self,
                                 option_time_period: list[pipeline_pb2.GetTradeDetailReq.TimePeriodOption],
                                 option_top_cus: pipeline_pb2.GetTradeDetailReq.TopCusOption,
                                 option_trade_frequency: list[pipeline_pb2.GetTradeDetailReq.TradeFrequencyOption],
                                 trade_type: pipeline_pb2.GetTradeDetailReq.TradeType) -> dict:

        y0 = int(self.doc["attribute_month"][:4])
        if trade_type == pipeline_pb2.GetTradeDetailReq.TradeType.CUSTOMER:
            _df1 = self._transpose_sales(y0, "customerDetail_24")
            _df2 = self._transpose_sales(y0 - 1, "customerDetail_25")
            _df3 = self._transpose_sales(y0 - 2, "customerDetail_26")
            _df4 = self._transpose_sales(y0 - 3, "customerDetail_27")
        elif trade_type == pipeline_pb2.GetTradeDetailReq.TradeType.SUPPLIER:
            _df1 = self._transpose_sales(y0, "supplierRanking_24", name_col='SALES_NAME')
            _df2 = self._transpose_sales(y0 - 1, "supplierRanking_25", name_col='SALES_NAME')
            _df3 = self._transpose_sales(y0 - 2, "supplierRanking_26", name_col='SALES_NAME')
            _df4 = self._transpose_sales(y0 - 3, "supplierRanking_27", name_col='SALES_NAME')
        else:
            raise ValueError("invalid tradeType")

        df_sales_detail = pd.concat([_df1, _df2, _df3, _df4])
        df_sales_detail.sort_values(by='month', inplace=True)
        df_sales_detail.reset_index(drop=True, inplace=True)

        _to_concat = []
        if pipeline_pb2.GetTradeDetailReq.TimePeriodOption.PERIOD_1ST in option_time_period:
            if _df1 is not None and not _df1.empty:
                _to_concat.append(_df1)
        if pipeline_pb2.GetTradeDetailReq.TimePeriodOption.PERIOD_2ND in option_time_period:
            if _df2 is not None and not _df2.empty:
                _to_concat.append(_df2)
        if pipeline_pb2.GetTradeDetailReq.TimePeriodOption.PERIOD_3RD in option_time_period:
            if _df3 is not None and not _df3.empty:
                _to_concat.append(_df3)
        if pipeline_pb2.GetTradeDetailReq.TimePeriodOption.PERIOD_4TH in option_time_period:
            if _df4 is not None and not _df4.empty:
                _to_concat.append(_df4)

        if len(_to_concat) == 0:
            return {
                "optionalSalesDetailLineChart": [],
                "optionalSalesDetailPieChart": {}
            }

        df_filtered = pd.concat(_to_concat, ignore_index=True)
        df_filtered.sort_values(by='month', inplace=True)
        df_filtered.reset_index(drop=True, inplace=True)
        ser_sum: pd.Series = df_filtered.drop(columns=['month']).sum(axis=0)
        ser_sum.sort_values(ascending=False, inplace=True)

        top_cus_dict = {
            pipeline_pb2.GetTradeDetailReq.TopCusOption.TOP_1: 1,
            pipeline_pb2.GetTradeDetailReq.TopCusOption.TOP_5: 5,
            pipeline_pb2.GetTradeDetailReq.TopCusOption.TOP_10: 10,
            pipeline_pb2.GetTradeDetailReq.TopCusOption.TOP_20: 20,
            pipeline_pb2.GetTradeDetailReq.TopCusOption.ALL: len(ser_sum)
        }
        head_n = top_cus_dict.get(option_top_cus, len(ser_sum))

        cus_list = ser_sum.drop(index=['其他']).head(head_n).index.tolist()

        ser_trade_frequency: pd.Series = df_filtered[cus_list].apply(lambda x: x.dropna().ne(0).sum() / len(x), axis=0)

        freq_dict = {
            pipeline_pb2.GetTradeDetailReq.TradeFrequencyOption.FREQUENCY_LOW: (0, 0.35),
            pipeline_pb2.GetTradeDetailReq.TradeFrequencyOption.FREQUENCY_MID: (0.35, 0.68),
            pipeline_pb2.GetTradeDetailReq.TradeFrequencyOption.FREQUENCY_HIGH: (0.68, 1.01),
        }
        cus_list_next = []
        for freq in option_trade_frequency:
            bound = freq_dict.get(freq, None)
            if bound is not None:
                ser_trade_frequency_matched = ser_trade_frequency[
                    (ser_trade_frequency >= bound[0]) & (ser_trade_frequency < bound[1])]
                cus_list_next.extend(ser_trade_frequency_matched.index.tolist())

        cus_list_next = list(set(cus_list_next))

        ser_sum_matched: pd.Series = ser_sum.loc[cus_list_next]

        rest_value = ser_sum.sum() - ser_sum_matched.sum()

        pie_data = []
        for c, v in ser_sum_matched.items():
            pie_data.append(
                {
                    "name": c,
                    "value": round(v, 2)
                }
            )
        pie_data.append(
            {
                "name": "其他",
                "value": round(rest_value, 2)
            }
        )

        line_data = []
        line_data.append(
            ['年月'] + df_sales_detail['month'].apply(lambda x: x.strftime("%Y-%m")).tolist()
        )
        for col in cus_list_next:
            line_data.append(
                [col] + df_sales_detail[col].tolist()
            )

        return {
            "optionalSalesDetailLineChart": line_data,
            "optionalSalesDetailPieChart": pie_data
        }

    async def process(self):

        # pprint(self.content)
        # # save self.content as json file
        self.reformat()
        self.content["impJsonDate"] = self.doc["attribute_month"] + "-01"
        self._major_commodity_proportion("ssfmxxSum", "bngyssptj")
        self._major_commodity_proportion("ssfmxxSum", "bnkhsptj")
        self._major_commodity_proportion("ssfmxxSum", "sngyssptj")
        self._major_commodity_proportion("ssfmxxSum", "snkhsptj")
        self._major_commodity_proportion("ssfmxxSum", "ssngyssptj")
        self._major_commodity_proportion("ssfmxxSum", "ssnkhsptj")
        self._major_commodity_proportion("ssfmxxSum", "sssngyssptj")
        self._major_commodity_proportion("ssfmxxSum", "sssnkhsptj")

        names = self.partner_company_names
        names.append(self.subj_company_name)

        csq = await self.repo.mongo_db.get_customer_supplier_quality(self.doc['content_id'])
        if csq:
            self._customer_supplier_quality(csq['data'])
        else:
            self.content['customerQuality'] = None
            self.content['supplierQuality'] = None

        # tags_tuple = await asyncio.gather(*(self.repo.dw_data.get_tags_by_name(row, True) for row in names))
        tp = self.time_point
        tags_tuple = await asyncio.gather(*(self.get_tags_by_name(row, tp) for row in names))
        tag_dict = {name: tag for name, tag in tags_tuple}

        # supplier
        self._assign_tag('ssfmxxSum', "bngyssptj", "SALES_NAME", tag_dict)
        self._assign_tag('ssfmxxSum', "sngyssptj", "SALES_NAME", tag_dict)
        self._assign_tag('ssfmxxSum', "ssngyssptj", "SALES_NAME", tag_dict)
        self._assign_tag('ssfmxxSum', "sssngyssptj", "SALES_NAME", tag_dict)

        # customer
        self._assign_tag('ssfmxxSum', "bnkhsptj", "PURCHASER_NAME", tag_dict)
        self._assign_tag('ssfmxxSum', "snkhsptj", "PURCHASER_NAME", tag_dict)
        self._assign_tag('ssfmxxSum', "ssnkhsptj", "PURCHASER_NAME", tag_dict)
        self._assign_tag('ssfmxxSum', "sssnkhsptj", "PURCHASER_NAME", tag_dict)

        # 3.add capitalPaidIn field to businessInfo
        try:
            subj_d = tag_dict.get(self.content["businessInfo"]["QYMC"])
            self.content["businessInfo"]["capitalPaidIn"] = None
            if subj_d is not None:
                self.content["businessInfo"]["capitalPaidIn"] = (
                    {} if subj_d.get("companyInfo", {}) is None else subj_d.get("companyInfo", {})).get(
                    "paid_in_capital")
                self.content['subjectCompanyTags'] = {_k: _v for _k, _v in subj_d.items() if _k != "companyInfo"}
            else:
                pass
        except KeyError:
            pass

        if self.content.get('subjectCompanyTags') is None:
            self.content['subjectCompanyTags'] = {}

        await self._subject_company_related_info()

        thread1 = threading.Thread(target=self._customer_ratio_chart_rewrite)
        thread2 = threading.Thread(target=self._supplier_ratio_chart_rewrite)
        thread1.start()
        thread2.start()

        self._subj_product_proportion()

        self._subj_purchase_proportion()

        self._financial_summary()

        self._trading_summary()

        self._revenue_detail_summary()

        self._risk_indexes()

        self._costumer_analysis_last3years()

        self._subj_product_proportion()

        self._supplier_affiliate_transaction_detection()
        self._customer_affiliate_transaction_detection()

        thread1.join()
        thread2.join()

    def _customer_supplier_quality(self, data: dict):
        shareholder_type_key = {
            '国企控股': 'stateCentral',
            '上市公司主体': 'listedChild',
            '上市公司控股': 'listedChild',
            '自然人控股、国企参股': 'jointVenture',
            '自然人控股、上市公司参股': 'jointVenture',
            '自然人控股、境外公司参股': 'jointVenture',
            '自然人控股、国企、上市公司参股': 'jointVenture',
            '自然人控股、国企、境外公司参股': 'jointVenture',
            '自然人控股、上市公司、境外公司参股': 'jointVenture',
            '自然人控股、国企、上市公司、境外公司参股': 'jointVenture',
            '自然人、国企、上市公司合资': 'jointVenture',
            '自然人、国企、境外公司合资': 'jointVenture',
            '自然人、上市公司、境外公司合资': 'jointVenture',
            '国企、上市公司、境外公司合资': 'jointVenture',
            '自然人、国企、上市公司、境外公司合资': 'jointVenture',
            "自然人全资": 'ownPerson',
            '境外公司全资': 'oversea',
            '境外公司控股、自然人参股': 'oversea',
            '境外公司控股、国企参股': 'oversea',
            '境外公司控股、上市公司参股': 'oversea',
            '境外公司控股、国企、上市公司参股': 'oversea',
            '其他': 'other'
        }

        financing_type_key = {'无融资': 'nothing', '知名机构投资': 'famous', '非知名机构投资': 'unknown'}

        credential_type_key = {'一档': 'first', '二档': 'second', '三档': 'third', '四档': 'fourth', '无': 'nothing'}

        ranking_type_key = {'7个以上': 'overSeven', '4-7个': 'fourSeven', '1-4个': 'oneThree', '0-1个': 'zeroOne'}

        law_type_key = {'无': 'nothing', '低': 'low', '中': 'middle', '高': 'high'}

        capital_type_key = {'无': 'nothing', '低': 'low', '中': 'middle', '高': 'high'}

        def get_summary(shareholder, law, financing, credential, customer_supplier, capital, ranking):
            shareholder_percent = (shareholder['ownPerson']['percent'] + shareholder['other']['percent'] +
                                   shareholder['oversea']['percent']) * 100
            if 0 <= shareholder_percent < 40:
                s1 = "股东实力很强"
            elif 40 <= shareholder_percent < 70:
                s1 = "股东实力较强"
            else:
                s1 = ""

            law_percent = (law['nothing']['percent'] + law['low']['percent']) * 100
            if 80 <= law_percent <= 100 and (law['low']['percent']) * 100 <= 1:
                s2 = "无诉讼风险"
            elif 70 <= law_percent <= 100 and (law['low']['percent']) * 100 <= 3:
                s2 = "诉讼风险较小"
            else:
                s2 = ""

            financing_nothing = financing['nothing']['percent'] * 100
            famous = financing['famous']['percent']
            unknown = financing['unknown']['percent']
            if 0 <= financing_nothing < 50:
                if unknown == 0:
                    s3 = "融资多，全部是知名机构投资"
                elif famous == 0:
                    s3 = "融资多，全部是非知名机构投资"
                elif famous > unknown:
                    s3 = "融资多，大部分是知名机构投资"
                elif famous < unknown:
                    s3 = "融资多，小部分是知名机构投资"
            else:
                s3 = ""

            credential_percent = credential['first']['percent'] * 100
            if 50 <= credential_percent <= 100:
                s4 = "市场份额、潜力很高"
            elif 30 <= credential_percent < 50:
                s4 = "市场份额、潜力较高"
            else:
                s4 = ""

            capital_high_percent = capital['high']['percent'] * 100
            capital_nothing_percent = capital['nothing']['percent'] * 100
            capital_low_percent = capital['low']['percent'] * 100
            s5 = ""
            if 0 <= capital_high_percent <= 2:
                if 80 <= capital_nothing_percent <= 100 and 80 <= capital_low_percent <= 100:
                    s5 = f"{'下游客户' if customer_supplier == 'customer' else '上游供应商'}很优，资金风险很小"
                elif 60 <= capital_nothing_percent < 80 and 60 <= capital_low_percent < 80:
                    s5 = f"{'下游客户' if customer_supplier == 'customer' else '上游供应商'}较优，资金风险较小"

            elif 2 < capital_high_percent <= 5:
                if 80 <= capital_nothing_percent < 98 and 80 <= capital_low_percent < 98:
                    s5 = f"{'下游客户' if customer_supplier == 'customer' else '上游供应商'}较优，资金风险较小"

            elif 5 <= capital_high_percent <= 8:
                if 85 <= capital_nothing_percent < 95 and 85 <= capital_low_percent < 95:
                    s5 = f"{'下游客户' if customer_supplier == 'customer' else '上游供应商'}较优，资金风险较小"

            ranking_percent = (ranking['overSeven']['percent'] + ranking['fourSeven']['percent']) * 100
            if ranking_percent >= 50:
                s6 = "入选榜单很多"
            elif 30 <= ranking_percent < 50:
                s6 = "入选榜单较多"
            else:
                s6 = ""

            return [s1, s2, s3, s4, s5, s6]

        # 客户
        def process_customer_detail(arr, shareholder_type_key, financing_type_key, credential_type_key, ranking_type_key, law_type_key, capital_type_key):
            shareholder_customer = {key: {"percent": 0, "info": []} for key in shareholder_type_key.values()}
            financing_customer = {key: {"percent": 0, "info": []} for key in financing_type_key.values()}
            credential_customer = {key: {"percent": 0, "info": []} for key in credential_type_key.values()}
            ranking_customer = {key: {"percent": 0, "info": []} for key in ranking_type_key.values()}
            law_customer = {key: {"percent": 0, "info": []} for key in law_type_key.values()}
            capital_customer = {key: {"percent": 0, "info": []} for key in capital_type_key.values()}

            for i in arr:
                for k in i['data']:
                    if k['category'] == '股东背景' and k['type'] in shareholder_type_key:
                        key = shareholder_type_key[k['type']]
                        shareholder_customer[key]['percent'] += i['transaction_percent']
                        shareholder_customer[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '融资信息' and k['type'] in financing_type_key:
                        key = financing_type_key[k['type']]
                        financing_customer[key]['percent'] += i['transaction_percent']
                        financing_customer[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '政府资质' and k['type'] in credential_type_key:
                        key = credential_type_key[k['type']]
                        credential_customer[key]['percent'] += i['transaction_percent']
                        credential_customer[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '榜单' and k['type'] in ranking_type_key:
                        key = ranking_type_key[k['type']]
                        ranking_customer[key]['percent'] += i['transaction_percent']
                        ranking_customer[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '诉讼' and k['type'] in law_type_key:
                        key = law_type_key[k['type']]
                        law_customer[key]['percent'] += i['transaction_percent']
                        law_customer[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '资金风险' and k['type'] in capital_type_key:
                        key = capital_type_key[k['type']]
                        capital_customer[key]['percent'] += i['transaction_percent']
                        capital_customer[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

            return shareholder_customer, financing_customer, credential_customer, ranking_customer, law_customer, capital_customer

        def create_customer_data(arr, year_data):
            shareholder_customer, financing_customer, credential_customer, ranking_customer, law_customer, capital_customer = process_customer_detail(
                arr, shareholder_type_key, financing_type_key, credential_type_key, ranking_type_key, law_type_key, capital_type_key
            )

            year_arr = [self.tag_t0, self.tag_t1, self.tag_t2, self.tag_t3]
            year = ""
            for item in year_arr:
                if str(year_data['year']) in item:
                    year = item
                    break

            summary = get_summary(shareholder_customer, law_customer, financing_customer, credential_customer, "customer", capital_customer, ranking_customer)

            return {
                "year": year,
                "quality": year_data['quality'],
                "detail": [
                    {"name": "资金风险", "data": capital_customer},
                    {"name": "企业性质", "data": shareholder_customer},
                    {"name": "司法风险", "data": law_customer},
                    {"name": "融资信息", "data": financing_customer},
                    {"name": "政府资质", "data": credential_customer},
                    {"name": "榜单数量", "data": ranking_customer},
                ],
                "summary": summary
            }

        if data['customerDetail_24']:
            latest_arr_customer = data['customerDetail_24']['detail']
            latest_year_customer = create_customer_data(latest_arr_customer, data['customerDetail_24'])
        else:
            latest_year_customer = None

        second_arr_customer = data['customerDetail_25']['detail']
        second_year_customer = create_customer_data(second_arr_customer, data['customerDetail_25'])

        third_arr_customer = data['customerDetail_26']['detail']
        third_year_customer = create_customer_data(third_arr_customer, data['customerDetail_26'])

        fourth_arr_customer = data['customerDetail_27']['detail']
        fourth_year_customer = create_customer_data(fourth_arr_customer, data['customerDetail_27'])

        customer_result = [latest_year_customer, second_year_customer, third_year_customer, fourth_year_customer] if latest_year_customer else [second_year_customer, third_year_customer, fourth_year_customer]
        customer_conclusion = f"{latest_year_customer['year']}、{second_year_customer['year']}、{third_year_customer['year']}、{fourth_year_customer['year']}下游客户资质：{latest_year_customer['quality']} -> {second_year_customer['quality']} -> {third_year_customer['quality']} -> {fourth_year_customer['quality']}" \
            if latest_year_customer else f"{second_year_customer['year']}、{third_year_customer['year']}、{fourth_year_customer['year']}下游客户资质：{second_year_customer['quality']} -> {third_year_customer['quality']} -> {fourth_year_customer['quality']}"

        # 供应商
        def process_supplier_detail(arr, shareholder_type_key, financing_type_key, credential_type_key, ranking_type_key, law_type_key, capital_type_key):
            shareholder_supplier = {key: {"percent": 0, "info": []} for key in shareholder_type_key.values()}
            financing_supplier = {key: {"percent": 0, "info": []} for key in financing_type_key.values()}
            credential_supplier = {key: {"percent": 0, "info": []} for key in credential_type_key.values()}
            ranking_supplier = {key: {"percent": 0, "info": []} for key in ranking_type_key.values()}
            law_supplier = {key: {"percent": 0, "info": []} for key in law_type_key.values()}
            capital_supplier = {key: {"percent": 0, "info": []} for key in capital_type_key.values()}

            for i in arr:
                for k in i['data']:
                    if k['category'] == '股东背景' and k['type'] in shareholder_type_key:
                        key = shareholder_type_key[k['type']]
                        shareholder_supplier[key]['percent'] += i['transaction_percent']
                        shareholder_supplier[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '融资信息' and k['type'] in financing_type_key:
                        key = financing_type_key[k['type']]
                        financing_supplier[key]['percent'] += i['transaction_percent']
                        financing_supplier[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '政府资质' and k['type'] in credential_type_key:
                        key = credential_type_key[k['type']]
                        credential_supplier[key]['percent'] += i['transaction_percent']
                        credential_supplier[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '榜单' and k['type'] in ranking_type_key:
                        key = ranking_type_key[k['type']]
                        ranking_supplier[key]['percent'] += i['transaction_percent']
                        ranking_supplier[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '诉讼' and k['type'] in law_type_key:
                        key = law_type_key[k['type']]
                        law_supplier[key]['percent'] += i['transaction_percent']
                        law_supplier[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

                    if k['category'] == '资金风险' and k['type'] in capital_type_key:
                        key = capital_type_key[k['type']]
                        capital_supplier[key]['percent'] += i['transaction_percent']
                        capital_supplier[key]['info'].append({
                            "companyName": i['enterprise_name'],
                            "percent": i['transaction_percent'],
                        })

            return shareholder_supplier, financing_supplier, credential_supplier, ranking_supplier, law_supplier, capital_supplier

        def create_supplier_data(arr, year_data):
            shareholder_supplier, financing_supplier, credential_supplier, ranking_supplier, law_supplier, capital_supplier = process_supplier_detail(
                arr, shareholder_type_key, financing_type_key, credential_type_key, ranking_type_key, law_type_key, capital_type_key
            )

            year_arr = [self.tag_t0, self.tag_t1, self.tag_t2, self.tag_t3]
            year = ""
            for item in year_arr:
                if str(year_data['year']) in item:
                    year = item
                    break

            summary = get_summary(shareholder_supplier, law_supplier, financing_supplier, credential_supplier, "supplier", capital_supplier, ranking_supplier)

            return {
                "year": year,
                "quality": year_data['quality'],
                "detail": [
                    {"name": "资金风险", "data": capital_supplier},
                    {"name": "企业性质", "data": shareholder_supplier},
                    {"name": "司法风险", "data": law_supplier},
                    {"name": "融资信息", "data": financing_supplier},
                    {"name": "政府资质", "data": credential_supplier},
                    {"name": "榜单数量", "data": ranking_supplier},
                ],
                "summary": summary
            }

        if data['supplierRanking_24']:
            latest_arr_supplier = data['supplierRanking_24']['detail']
            latest_year_supplier = create_supplier_data(latest_arr_supplier, data['supplierRanking_24'])
        else:
            latest_year_supplier = None

        second_arr_supplier = data['supplierRanking_25']['detail']
        second_year_supplier = create_supplier_data(second_arr_supplier, data['supplierRanking_25'])

        third_arr_supplier = data['supplierRanking_26']['detail']
        third_year_supplier = create_supplier_data(third_arr_supplier, data['supplierRanking_26'])

        fourth_arr_supplier = data['supplierRanking_27']['detail']
        fourth_year_supplier = create_supplier_data(fourth_arr_supplier, data['supplierRanking_27'])

        supplier_result = [latest_year_supplier, second_year_supplier, third_year_supplier, fourth_year_supplier] if latest_year_supplier else [second_year_supplier, third_year_supplier, fourth_year_supplier]
        supplier_conclusion = f"{latest_year_supplier['year']}、{second_year_supplier['year']}、{third_year_supplier['year']}、{fourth_year_supplier['year']}上游供应商资质：{latest_year_supplier['quality']} -> {second_year_supplier['quality']} -> {third_year_supplier['quality']} -> {fourth_year_supplier['quality']}" \
            if latest_year_supplier else f"{second_year_supplier['year']}、{third_year_supplier['year']}、{fourth_year_supplier['year']}上游供应商资质：{second_year_supplier['quality']} -> {third_year_supplier['quality']} -> {fourth_year_supplier['quality']}"

        self.content['customerQuality'] = customer_result
        self.content['supplierQuality'] = supplier_result

        self.content['customerQualityConclusion'] = customer_conclusion
        self.content['supplierQualityConclusion'] = supplier_conclusion

    def _customer_affiliate_transaction_detection(self):
        def _sup_parse_trade_detail(key: str):
            pie_data = []
            d = self.content.get('ssfmxxSum', {})
            accum = 0
            if isinstance(d, dict):
                detail = d.get(key, None)
                if isinstance(detail, list):
                    for _item in detail:
                        if isinstance(_item, dict):
                            com_name = _item.get("PURCHASER_NAME", None)
                            if isinstance(com_name, str):
                                if com_name in affiliates:
                                    ratio_str = _item.get("RATIO_AMOUNT_TAX", None)
                                    if isinstance(ratio_str, str):
                                        try:
                                            ratio_str = ratio_str.replace("%", "")
                                            r = float(ratio_str) / 100
                                        except Exception:
                                            r = None
                                        if isinstance(r, float):
                                            pie_data.append(
                                                {"name": com_name, "value": round(r, 2)}
                                            )
                                            accum += r
                else:
                    return None
            remain = 1 - accum
            pie_data.append(
                {"name": "非关联交易", "value": round(remain, 2)}
            )
            return pie_data

        affiliates = []
        shareholders = self.content.get("shareholderData", None)
        if isinstance(shareholders, list):
            for item in shareholders:
                if isinstance(item, dict):
                    name = item.get("StockName", None)
                    if isinstance(name, str) and name != "":
                        affiliates.append(name)

        branches = self.content.get("branchesData", None)
        if isinstance(branches, list):
            for item in branches:
                if isinstance(item, dict):
                    name = item.get("Name", None)
                    if isinstance(name, str) and name != "":
                        affiliates.append(name)

        investment = self.content.get("investmentData", None)
        if isinstance(investment, list):
            for item in investment:
                if isinstance(item, dict):
                    name = item.get("Name", None)
                    if isinstance(name, str) and name != "":
                        affiliates.append(name)

        t0_data = _sup_parse_trade_detail('bnkhsptj')
        t1_data = _sup_parse_trade_detail('snkhsptj')
        t2_data = _sup_parse_trade_detail('ssnkhsptj')
        t3_data = _sup_parse_trade_detail('sssnkhsptj')

        dataset = [t0_data, t1_data, t2_data, t3_data]
        titles = [self.tag_t0, self.tag_t1, self.tag_t2, self.tag_t3]

        res = []
        for t, dt in zip(titles, dataset):
            if dt is not None:
                res.append(
                    {
                        "title": f"{t}关联交易占比",
                        "data": dt
                    }
                )

        self.content['affiliateTradePieCharts'] = res

    def _supplier_affiliate_transaction_detection(self):
        def _parse_trade_detail(key: str):
            pie_data = []
            d = self.content.get('ssfmxxSum', {})
            accum = 0
            if isinstance(d, dict):
                detail = d.get(key, None)
                if isinstance(detail, list):
                    for _item in detail:
                        if isinstance(_item, dict):
                            com_name = _item.get("SALES_NAME", None)
                            if isinstance(com_name, str):
                                if com_name in affiliates:
                                    ratio_str = _item.get("RATIO_AMOUNT_TAX", None)
                                    if isinstance(ratio_str, str):
                                        try:
                                            ratio_str = ratio_str.replace("%", "")
                                            r = float(ratio_str) / 100
                                        except Exception:
                                            r = None
                                        if isinstance(r, float):
                                            pie_data.append(
                                                {"name": com_name, "value": round(r, 2)}
                                            )
                                            accum += r
                else:
                    return None
            remain = 1 - accum
            pie_data.append(
                {"name": "非关联交易", "value": round(remain, 2)}
            )
            return pie_data

        affiliates = []
        shareholders = self.content.get("shareholderData", None)
        if isinstance(shareholders, list):
            for item in shareholders:
                if isinstance(item, dict):
                    name = item.get("StockName", None)
                    if isinstance(name, str) and name != "":
                        affiliates.append(name)

        branches = self.content.get("branchesData", None)
        if isinstance(branches, list):
            for item in branches:
                if isinstance(item, dict):
                    name = item.get("Name", None)
                    if isinstance(name, str) and name != "":
                        affiliates.append(name)

        investment = self.content.get("investmentData", None)
        if isinstance(investment, list):
            for item in investment:
                if isinstance(item, dict):
                    name = item.get("Name", None)
                    if isinstance(name, str) and name != "":
                        affiliates.append(name)
        # bngyssptj sngyssptj ssngyssptj sssngyssptj
        t0_data = _parse_trade_detail('bngyssptj')
        t1_data = _parse_trade_detail('sngyssptj')
        t2_data = _parse_trade_detail('ssngyssptj')
        t3_data = _parse_trade_detail('sssngyssptj')

        dataset = [t0_data, t1_data, t2_data, t3_data]
        titles = [self.tag_t0, self.tag_t1, self.tag_t2, self.tag_t3]

        res = []
        for t, dt in zip(titles, dataset):
            if dt is not None:
                res.append(
                    {
                        "title": f"{t}关联交易占比",
                        "data": dt
                    }
                )

        self.content['supplierAffiliateTradePieCharts'] = res

    @property
    def subj_company_name(self) -> str:
        _info = self.content.get("businessInfo", {})
        return _info.get("QYMC")

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

    def _customer_ratio_chart_rewrite(self):
        y1 = int(self.doc["attribute_month"][:4])
        _df1 = self._transpose_sales(y1, "customerDetail_24")
        _df2 = self._transpose_sales(y1 - 1, "customerDetail_25")
        _df3 = self._transpose_sales(y1 - 2, "customerDetail_26")
        _df4 = self._transpose_sales(y1 - 3, "customerDetail_27")

        df_sales_detail = pd.concat([_df1, _df2, _df3, _df4])
        df_sales_detail.replace(0, np.nan, inplace=True)
        df_sales_detail.sort_values(by='month', inplace=True)
        df_sales_detail.reset_index(drop=True, inplace=True)

        def _recent_deactivate(df):
            start_m = self.time_point - relativedelta(months=12)
            df = df[df['month'] >= start_m]
            df.reset_index(drop=True, inplace=True)

            purchase_y = []
            df.replace(0, np.nan, inplace=True)
            cols = df.columns.tolist()[1:]
            for col in cols:
                last_idx_valid = df[col].last_valid_index()
                first_idx_valid = df[col].first_valid_index()
                if last_idx_valid:
                    purchase_y.append({
                        "name": col,
                        "last_valid_month": df.loc[last_idx_valid, 'month'],
                        "first_valid_month": df.loc[first_idx_valid, 'month']
                    })

            df_valid_idx = pd.DataFrame(purchase_y)
            end_m = self.time_point - relativedelta(months=3)
            df_valid_idx = df_valid_idx[
                (df_valid_idx['last_valid_month'] > start_m) & (df_valid_idx['last_valid_month'] < end_m)]

            recent3_deactivate_cus = df_valid_idx['name'].tolist()

            pie_data = []
            for cus in recent3_deactivate_cus:
                value = df[cus].sum()
                pie_data.append(
                    {
                        "name": cus,
                        "value": round(value, 2)
                    }
                )

            to_drop = recent3_deactivate_cus + ["month"]
            df['其他'] = df.drop(columns=to_drop).sum(axis=1).round(decimals=2)
            pie_data.append({
                "name": "其他",
                "value": round(df['其他'].sum(), 2)
            })

            res = {
                "name": "近三个月未合作客户近12个月销售占比",
                "data": pie_data
            }
            with self.__lock:
                self.content["recent3MInactivePieCharts"] = res

            line_data = {}
            months = df['month'].apply(lambda x: x.strftime("%Y-%m")).tolist()
            line_data["x_axis"] = months
            series = []

            cus_add = recent3_deactivate_cus
            for col in cus_add:
                series.append(df[col].tolist())
            line_data["series"] = series
            line_data["legend"] = cus_add
            with self.__lock:
                self.content['recent3MInactiveLineCharts'] = {
                    "data": line_data,
                    "name": "近三个月未合作客户近12个月销售占比"
                }

        _recent_deactivate(df_sales_detail.copy())

        def _calculate_cr_values(df):
            df['year'] = df['month'].dt.year
            annual_sales = df.drop(columns=['month']).groupby('year').sum()
            cr = pd.DataFrame()

            cr['CR1'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(1).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            cr['CR5'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(5).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            cr['CR10'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(10).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            cr['CR20'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(20).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            max_m = df['month'].max().strftime("%m")
            min_m = df['month'].max().strftime("%m")
            m_list = annual_sales.index.tolist()
            m_list[0] = f"{m_list[0]}年{min_m}月-12月"
            m_list[-1] = f"{m_list[-1]}年1月-{max_m}月"
            cr_pct = cr.pct_change()
            summary_cr = []
            summary_cr.append(
                f"{m_list[-1]}第1大客户销售占比{cr.iloc[-1]['CR1']:.2%}，同比{self._desc(r1 := cr_pct.iloc[-1]['CR1'])}{abs(r1):.2%};\n"
                f"前5大客户销售占比{cr.iloc[-1]['CR5']:.2%}，同比{self._desc(r2 := cr_pct.iloc[-1]['CR5'])}{abs(r2):.2%};\n"
                f"前10大客户销售占比{cr.iloc[-1]['CR10']:.2%}，同比{self._desc(r3 := cr_pct.iloc[-1]['CR10'])}{abs(r3):.2%};\n"
                f"前20大客户销售占比{cr.iloc[-1]['CR20']:.2%}, 同比{self._desc(r4 := cr_pct.iloc[-1]['CR20'])}{abs(r4):.2%}。\n")
            summary_cr.append(
                f"{m_list[1]}第1大客户销售占比{cr.iloc[-2]['CR1']:.2%}，同比{self._desc(r1 := cr_pct.iloc[-2]['CR1'])}{abs(r1):.2%};\n"
                f"前5大客户销售占比{cr.iloc[-2]['CR5']:.2%}，同比{self._desc(r2 := cr_pct.iloc[-2]['CR5'])}{abs(r2):.2%};\n"
                f"前10大客户销售占比{cr.iloc[-2]['CR10']:.2%}，同比{self._desc(r3 := cr_pct.iloc[-2]['CR10'])}{abs(r3):.2%};\n"
                f"前20大客户销售占比{cr.iloc[-2]['CR20']:.2%}, 同比{self._desc(r4 := cr_pct.iloc[-2]['CR20'])}{abs(r4):.2%}。")

            cr[['CR1', 'CR5', 'CR10', 'CR20']] *= 100
            # cr_pct *= 100
            return {
                'x_axis': m_list,
                'series': cr[['CR1', 'CR5', 'CR10', 'CR20']].values.round(2).T.tolist(),
                'legend': ['第1大客户销售占比', '前5大客户销售占比', '前10大客户销售占比',
                           '前20大客户销售占比'],
                'summary': summary_cr[::-1],
                'comment': "前5大客户的销售占比是衡量市场集中度的重要指标。可参考如下分级：\n0-25%显示市场高度分散；25%-40%表明客户分布较广；40%-50%反映市场均衡；50-70%指市场较集中；而70%-100%则意味着市场高度集中。"
            }

        cr_res = _calculate_cr_values(df_sales_detail)
        with self.__lock:
            self.content['concentrationRateChart'] = cr_res

        def _customer_repurchase_ratio_chart(_dfs: list, df, x_axis):
            dfs = []
            for df_ in _dfs:
                if isinstance(df_, pd.DataFrame):
                    if df_.shape[1] > 2:
                        dfs.append(df_)

            x_axis_reverse = x_axis.copy()
            x_axis_reverse.reverse()
            # df['year'] = df['month'].dt.year
            df_gp_year_sum: pd.DataFrame = df.drop(columns=['month']).groupby('year').sum()
            # y3_repeat_ratio = df_gp_year_sum[y3_repeat.index.tolist()].loc[y3] / df_gp_year_sum.loc[y3].sum()
            year = df_gp_year_sum.sort_index(ascending=False).index.tolist()
            bar_chart = []
            pie_chart = []
            repurchase_amount = []
            repurchase_amount_ratio = []
            cus_repeat_list = []
            for i in range(len(year)):
                if i > len(year) - 2:
                    break
                cus_list_current = dfs[i].drop(columns=['其他', 'month']).columns.tolist()
                cus_list_prev = dfs[i + 1].drop(columns=['其他', 'month']).columns.tolist()

                cus_repeat = list(set(cus_list_current).intersection(cus_list_prev))
                cus_new = list(set(cus_list_current).difference(cus_list_prev))
                cus_dropped = list(set(cus_list_prev).difference(cus_list_current))

                y = year[i]

                # repurchase sum ratio
                cus_repeat_list.append(cus_repeat)
                repurchase_amount.append(
                    round(df_gp_year_sum.loc[y, cus_repeat].sum(), 2)
                )
                repurchase_amount_ratio.append(
                    round(df_gp_year_sum.loc[y, cus_repeat].sum() / df_gp_year_sum.loc[y].sum() * 100, 1)
                )

                y1_new_ratio = df_gp_year_sum[cus_new].loc[y] / df_gp_year_sum.loc[y].sum()
                y1_unknown_ratio = df_gp_year_sum[["其他"]].loc[y] / df_gp_year_sum.loc[y].sum()
                y1_drop_ratio = df_gp_year_sum[cus_dropped].loc[y2 := year[i + 1]] / df_gp_year_sum.loc[
                    y2].sum()
                y1_repeat_ratio = df_gp_year_sum[cus_repeat].loc[y] / df_gp_year_sum.loc[y].sum()
                y1_repeat_ratio_aty2 = df_gp_year_sum[cus_repeat].loc[y2] / df_gp_year_sum.loc[y2].sum()
                y1_repurchase = pd.concat([y1_repeat_ratio, y1_repeat_ratio_aty2], axis=1, keys=['curr', 'prev'])
                y1_repurchase['deltaRatio'] = y1_repurchase['curr'] / y1_repurchase['prev'].replace(0, np.nan) - 1
                y1_repurchase.sort_values(by='curr', ascending=False, inplace=True)

                curr_list = y1_repurchase['curr'].tolist()
                delta_list = y1_repurchase['deltaRatio'].tolist()
                bar_chart.append(
                    {
                        "title": f"{x_axis_reverse[i]}老客户供货占比",
                        'legend': ['供货占比', '同比'],
                        'x_axis': y1_repurchase.index.tolist(),
                        'series': [
                            [round(x1 * 100, 2) for x1 in curr_list],
                            [round(x2 * 100, 0) for x2 in delta_list],
                        ]
                    },
                )

                y1_new_ratio.loc['老客户'] = 1 - y1_new_ratio.sum() - y1_unknown_ratio.sum()
                y1_new_ratio.loc['未知'] = y1_unknown_ratio.sum()
                new_ratio_dict = y1_new_ratio.to_dict()

                y1_drop_ratio.loc['未退出'] = 1 - y1_drop_ratio.sum() - y1_unknown_ratio.sum()
                y1_drop_ratio.loc['未知'] = y1_unknown_ratio.sum()
                drop_ratio_dict = y1_drop_ratio.to_dict()
                pie_chart.append(
                    {
                        "name": f"{x_axis_reverse[i]}新晋前20名的客户（相比去年）",
                        "data": [{"name": k, "value": round(v * 100, 2)} for k, v in new_ratio_dict.items()],
                    },
                )
                pie_chart.append(
                    {
                        "name": f"{x_axis_reverse[i]}退出前20名的客户（相比去年）",
                        "data": [{"name": k, "value": round(v * 100, 2)} for k, v in drop_ratio_dict.items()],
                    },
                )

            stability_summary = ''
            for k in range(len(cus_repeat_list) - 1):
                try:
                    count_repeat = round(((len(cus_repeat_list[k]) / len(cus_repeat_list[k + 1])) - 1) * 100, 0)
                except ZeroDivisionError:
                    count_repeat = np.nan
                try:
                    ratio_repeat = round(((repurchase_amount_ratio[k] / repurchase_amount_ratio[k + 1]) - 1) * 100, 0)
                except ZeroDivisionError:
                    ratio_repeat = np.nan

                stability_summary += (f"{x_axis_reverse[k]}有老客户{len(cus_repeat_list[k])}名,"
                                      f"同比上升{count_repeat}%,"
                                      f"复购销售额占比{repurchase_amount_ratio[k]}%, "
                                      f"同比上升{ratio_repeat}%.\n")

            with self.__lock:
                stability_bar_chart = {
                    'summary': stability_summary,
                    'comment': "复购金额占比是衡量客户稳定性的重要指标。可参考如下分级：0-30%低忠诚、不稳定；30%-50%中等忠诚、适度稳定；50%-65%良好稳定性；65%-75%高忠诚、稳定；75%-100%极高忠诚、非常稳定。",
                    'legend': ['复购金额占比', '复购金额'],
                    'x_axis': x_axis_reverse[:len(repurchase_amount)][::-1],
                    'series': [repurchase_amount_ratio[::-1], repurchase_amount[::-1]]
                }
                self.content['customerStability'] = stability_bar_chart
                self.content['repurchaseRatioPieChart'] = pie_chart
                self.content['repurchaseRatioBarChart'] = bar_chart

        _customer_repurchase_ratio_chart([_df1, _df2, _df3, _df4], df_sales_detail.copy(), cr_res['x_axis'])


        def _mega_drop_customers(df: pd.DataFrame):
            def _calculate_yoy(s: pd.Series):
                s_last = s.shift(periods=1)
                return s / s_last - 1

            df_gp_year_sum: pd.DataFrame = df.drop(columns=['month']).groupby('year').sum()
            df_gp_year_sum.replace(0, np.nan, inplace=True)
            df_yoy = df_gp_year_sum.apply(lambda x: _calculate_yoy(x))
            total = df_gp_year_sum.sum(axis=1)
            df_prop = df_gp_year_sum.div(total, axis=0)

            sales_prop, ratio, x_axis = [], [], []
            for col in df_prop.columns:
                prop = df_prop.iloc[-1][col]
                yoy = df_yoy.iloc[-1][col]
                if prop >= 0.01 and yoy <= -0.2:
                    x_axis.append(col)
                    sales_prop.append(round(prop * 100, 2))
                    ratio.append(round(yoy * 100, 2))

            bar_res = {
                'title': f"{self.time_point.year}年订单同比大幅下滑客户",
                'x_axis': x_axis,
                'series': [sales_prop, ratio],
                'legend': ['销售占比', '同比']
            }
            with self.__lock:
                self.content['salesSharplyDropCustomersBarChart'] = bar_res

        _mega_drop_customers(df_sales_detail)

        def _sharply_drop_customers_this_month(df: pd.DataFrame):
            # set month as index
            df.set_index('month', inplace=True)
            df.drop(columns=['year', '其他'], inplace=True)
            df.replace(0, np.nan, inplace=True)

            df_mom = df.apply(lambda x: (x / x.shift(periods=1)) - 1, axis=0)

            cus = []
            for col in df_mom.columns:
                if df_mom.iloc[-1][col] < -0.2:
                    cus.append(col)

            sales = []
            for col in cus:
                sales.append(df[col].tolist())

            the_month = df.index.map(lambda x: x.strftime("%Y-%m")).tolist()
            line_chart = {
                'title': f"{the_month[-1]}订单环比大幅度下滑的客户",
                'series-list': sales,
                'legend': cus,
                'x-axis': the_month
            }
            with self.__lock:
                self.content['sharplyDropCustomerRecentMonthLineChart'] = line_chart

        _sharply_drop_customers_this_month(df_sales_detail.copy())

    def _supplier_ratio_chart_rewrite(self):
        y1 = int(self.doc["attribute_month"][:4])
        _df1 = self._transpose_sales(y1, "supplierRanking_24", name_col="SALES_NAME")
        _df2 = self._transpose_sales(y1 - 1, "supplierRanking_25", name_col="SALES_NAME")
        _df3 = self._transpose_sales(y1 - 2, "supplierRanking_26", name_col="SALES_NAME")
        _df4 = self._transpose_sales(y1 - 3, "supplierRanking_27", name_col="SALES_NAME")

        df_sales_detail = pd.concat([_df1, _df2, _df3, _df4])
        df_sales_detail.replace(0, np.nan, inplace=True)
        df_sales_detail.sort_values(by='month', inplace=True)
        df_sales_detail.reset_index(drop=True, inplace=True)

        def _recent_deactivate(df):
            start_m = self.time_point - relativedelta(months=12)
            df = df[df['month'] >= start_m]
            df.reset_index(drop=True, inplace=True)

            purchase_y = []
            df.replace(0, np.nan, inplace=True)
            cols = df.columns.tolist()[1:]
            for col in cols:
                last_idx_valid = df[col].last_valid_index()
                first_idx_valid = df[col].first_valid_index()
                if last_idx_valid:
                    purchase_y.append({
                        "name": col,
                        "last_valid_month": df.loc[last_idx_valid, 'month'],
                        "first_valid_month": df.loc[first_idx_valid, 'month']
                    })

            df_valid_idx = pd.DataFrame(purchase_y)
            end_m = self.time_point - relativedelta(months=3)
            df_valid_idx = df_valid_idx[
                (df_valid_idx['last_valid_month'] > start_m) & (df_valid_idx['last_valid_month'] < end_m)]

            recent3_deactivate_cus = df_valid_idx['name'].tolist()

            pie_data = []
            for cus in recent3_deactivate_cus:
                value = df[cus].sum()
                pie_data.append(
                    {
                        "name": cus,
                        "value": round(value, 2)
                    }
                )

            to_drop = recent3_deactivate_cus + ["month"]
            df['其他'] = df.drop(columns=to_drop).sum(axis=1).round(decimals=2)
            pie_data.append({
                "name": "其他",
                "value": round(df['其他'].sum(), 2)
            })

            res = {
                "name": "近三个月未合作供应商近12个月销售占比",
                "data": pie_data
            }
            with self.__lock:
                self.content["supplierRecent3MInactivePieCharts"] = res

            line_data = {}
            months = df['month'].apply(lambda x: x.strftime("%Y-%m")).tolist()
            line_data["x_axis"] = months
            series = []

            cus_add = recent3_deactivate_cus
            for col in cus_add:
                series.append(df[col].tolist())
            line_data["series"] = series
            line_data["legend"] = cus_add
            with self.__lock:
                self.content['supplierRecent3MInactiveLineCharts'] = {
                    "data": line_data,
                    "name": "近三个月未合作供应商近12个月销售占比"
                }

        _recent_deactivate(df_sales_detail.copy())

        def _calculate_cr_values(df):
            df['year'] = df['month'].dt.year
            annual_sales = df.drop(columns=['month']).groupby('year').sum()
            cr = pd.DataFrame()

            cr['CR1'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(1).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            cr['CR5'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(5).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            cr['CR10'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(10).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            cr['CR20'] = annual_sales.apply(
                lambda x: x[x.index != '其他'].sort_values(ascending=False).head(20).sum() / x.sort_values(
                    ascending=False).sum(), axis=1)
            max_m = df['month'].max().strftime("%m")
            min_m = df['month'].max().strftime("%m")
            m_list = annual_sales.index.tolist()
            m_list[0] = f"{m_list[0]}年{min_m}月-12月"
            m_list[-1] = f"{m_list[-1]}年1月-{max_m}月"
            cr_pct = cr.pct_change()
            summary_cr = []
            summary_cr.append(
                f"{m_list[-1]}第1大供应商采购占比{cr.iloc[-1]['CR1']:.2%}，同比{self._desc(r1 := cr_pct.iloc[-1]['CR1'])}{abs(r1):.2%};\n"
                f"前5大供应商采购占比{cr.iloc[-1]['CR5']:.2%}，同比{self._desc(r2 := cr_pct.iloc[-1]['CR5'])}{abs(r2):.2%};\n"
                f"前10大供应商采购占比{cr.iloc[-1]['CR10']:.2%}，同比{self._desc(r3 := cr_pct.iloc[-1]['CR10'])}{abs(r3):.2%};\n"
                f"前20大供应商采购占比{cr.iloc[-1]['CR20']:.2%}, 同比{self._desc(r4 := cr_pct.iloc[-1]['CR20'])}{abs(r4):.2%}。\n")
            summary_cr.append(
                f"{m_list[1]}第1大供应商采购占比{cr.iloc[-2]['CR1']:.2%}，同比{self._desc(r1 := cr_pct.iloc[-2]['CR1'])}{abs(r1):.2%};\n"
                f"前5大供应商采购占比{cr.iloc[-2]['CR5']:.2%}，同比{self._desc(r2 := cr_pct.iloc[-2]['CR5'])}{abs(r2):.2%};\n"
                f"前10大供应商采购占比{cr.iloc[-2]['CR10']:.2%}，同比{self._desc(r3 := cr_pct.iloc[-2]['CR10'])}{abs(r3):.2%};\n"
                f"前20大供应商采购占比{cr.iloc[-2]['CR20']:.2%}, 同比{self._desc(r4 := cr_pct.iloc[-2]['CR20'])}{abs(r4):.2%}。")

            cr[['CR1', 'CR5', 'CR10', 'CR20']] *= 100
            # cr_pct *= 100
            return {
                'x_axis': m_list,
                'series': cr[['CR1', 'CR5', 'CR10', 'CR20']].values.round(2).T.tolist(),
                'legend': ['第1大供应商采购占比', '前5大供应商采购占比', '前10大供应商采购占比',
                           '前20大供应商采购占比'],
                'summary': summary_cr,
                'comment': "前五大供应商的采购占比是衡量市场集中度的重要指标。可参考如下分级：\n0-25%显示市场高度分散；25%-40%表明供应商分布较广；40%-50%反映市场均衡；50-70%指市场较集中；而70%-100%则意味着市场高度集中。"
            }

        cr_res = _calculate_cr_values(df_sales_detail)
        with self.__lock:
            self.content['supplierConcentrationRateChart'] = cr_res

        def _supplier_repurchase_ratio_chart(_dfs: list, df, x_axis):
            dfs = []
            for df_ in _dfs:
                if isinstance(df_, pd.DataFrame):
                    if df_.shape[1] > 2:
                        dfs.append(df_)

            x_axis_reverse = x_axis.copy()
            x_axis_reverse.reverse()
            # df['year'] = df['month'].dt.year
            df_gp_year_sum: pd.DataFrame = df.drop(columns=['month']).groupby('year').sum()
            # y3_repeat_ratio = df_gp_year_sum[y3_repeat.index.tolist()].loc[y3] / df_gp_year_sum.loc[y3].sum()
            year = df_gp_year_sum.sort_index(ascending=False).index.tolist()
            bar_chart = []
            pie_chart = []
            repurchase_amount = []
            repurchase_amount_ratio = []
            cus_repeat_list = []
            for i in range(len(year)):
                if i > len(year) - 2:
                    break
                cus_list_current = dfs[i].drop(columns=['其他', 'month']).columns.tolist()
                cus_list_prev = dfs[i + 1].drop(columns=['其他', 'month']).columns.tolist()

                cus_repeat = list(set(cus_list_current).intersection(cus_list_prev))
                cus_new = list(set(cus_list_current).difference(cus_list_prev))
                cus_dropped = list(set(cus_list_prev).difference(cus_list_current))

                y = year[i]

                # repurchase sum ratio
                cus_repeat_list.append(cus_repeat)
                repurchase_amount.append(
                    round(df_gp_year_sum.loc[y, cus_repeat].sum(), 2)
                )
                repurchase_amount_ratio.append(
                    round(df_gp_year_sum.loc[y, cus_repeat].sum() / df_gp_year_sum.loc[y].sum() * 100, 1)
                )

                y1_new_ratio = df_gp_year_sum[cus_new].loc[y] / df_gp_year_sum.loc[y].sum()
                y1_unknown_ratio = df_gp_year_sum[["其他"]].loc[y] / df_gp_year_sum.loc[y].sum()
                y1_drop_ratio = df_gp_year_sum[cus_dropped].loc[y2 := year[i + 1]] / df_gp_year_sum.loc[
                    y2].sum()
                y1_repeat_ratio = df_gp_year_sum[cus_repeat].loc[y] / df_gp_year_sum.loc[y].sum()
                y1_repeat_ratio_aty2 = df_gp_year_sum[cus_repeat].loc[y2] / df_gp_year_sum.loc[y2].sum()
                y1_repurchase = pd.concat([y1_repeat_ratio, y1_repeat_ratio_aty2], axis=1, keys=['curr', 'prev'])
                y1_repurchase['deltaRatio'] = y1_repurchase['curr'] / y1_repurchase['prev'].replace(0, np.nan) - 1
                y1_repurchase.sort_values(by='curr', ascending=False, inplace=True)

                curr_list = y1_repurchase['curr'].tolist()
                delta_list = y1_repurchase['deltaRatio'].tolist()
                bar_chart.append(
                    {
                        "title": f"{x_axis_reverse[i]}老供应商供货占比",
                        'legend': ['供货占比', '同比'],
                        'x_axis': y1_repurchase.index.tolist(),
                        'series': [
                            [round(x1 * 100, 2) for x1 in curr_list],
                            [round(x2 * 100, 0) for x2 in delta_list],
                        ]
                    },
                )

                y1_new_ratio.loc['老供应商'] = 1 - y1_new_ratio.sum() - y1_unknown_ratio.sum()
                y1_new_ratio.loc['未知'] = y1_unknown_ratio.sum()
                new_ratio_dict = y1_new_ratio.to_dict()

                y1_drop_ratio.loc['未退出'] = 1 - y1_drop_ratio.sum() - y1_unknown_ratio.sum()
                y1_drop_ratio.loc['未知'] = y1_unknown_ratio.sum()
                drop_ratio_dict = y1_drop_ratio.to_dict()
                pie_chart.append(
                    {
                        "name": f"{x_axis_reverse[i]}新晋前20名的供应商（相比去年）",
                        "data": [{"name": k, "value": round(v * 100, 2)} for k, v in new_ratio_dict.items()],
                    },
                )
                pie_chart.append(
                    {
                        "name": f"{x_axis_reverse[i]}退出前20名的供应商（相比去年）",
                        "data": [{"name": k, "value": round(v * 100, 2)} for k, v in drop_ratio_dict.items()],
                    },
                )

            stability_summary = ''
            for k in range(len(cus_repeat_list) - 1):
                stability_summary += (f"{x_axis_reverse[k]}有供应商{len(cus_repeat_list[k])}名,"
                                      f"同比上升{round(((len(cus_repeat_list[k]) / len(cus_repeat_list[k + 1])) - 1) * 100, 0)}%,"
                                      f"复购金额占比{repurchase_amount_ratio[k]}%, "
                                      f"同比上升{round(((repurchase_amount_ratio[k] / repurchase_amount_ratio[k + 1]) - 1) * 100, 0)}%.\n")

            with self.__lock:
                self.content['supplierStability'] = {
                    'summary': stability_summary,
                    'comment': "复购金额占比是衡量客户稳定性的重要指标。可参考如下分级：0-30%低忠诚、不稳定；30%-50%中等忠诚、适度稳定；50%-65%良好稳定性；65%-75%高忠诚、稳定；75%-100%极高忠诚、非常稳定。",
                    'legend': ['复购金额占比', '复购金额'],
                    'x_axis': x_axis_reverse[:len(repurchase_amount)][::-1],
                    'series': [repurchase_amount_ratio[::-1], repurchase_amount[::-1]]
                }
                self.content['supplierRepurchaseRatioPieChart'] = pie_chart
                self.content['supplierRepurchaseRatioBarChart'] = bar_chart

        _supplier_repurchase_ratio_chart([_df1, _df2, _df3, _df4], df_sales_detail.copy(), cr_res['x_axis'])

        def _mega_drop_customers(df: pd.DataFrame):
            def _calculate_yoy(s: pd.Series):
                s_last = s.shift(periods=1)
                return s / s_last - 1

            df_gp_year_sum: pd.DataFrame = df.drop(columns=['month']).groupby('year').sum()
            df_gp_year_sum.replace(0, np.nan, inplace=True)
            df_yoy = df_gp_year_sum.apply(lambda x: _calculate_yoy(x))
            total = df_gp_year_sum.sum(axis=1)
            df_prop = df_gp_year_sum.div(total, axis=0)

            sales_prop, ratio, x_axis = [], [], []
            for col in df_prop.columns:
                prop = df_prop.iloc[-1][col]
                yoy = df_yoy.iloc[-1][col]
                if prop >= 0.01 and yoy <= -0.2:
                    x_axis.append(col)
                    sales_prop.append(round(prop * 100, 2))
                    ratio.append(round(yoy * 100, 2))

            bar_res = {
                'title': f"{self.time_point.year}年订单同比大幅下滑供应商",
                'x_axis': x_axis,
                'series': [sales_prop, ratio],
                'legend': ['销售占比', '同比']
            }
            with self.__lock:
                self.content['purchaseSharplyDropCustomersBarChart'] = bar_res

        _mega_drop_customers(df_sales_detail)

        def _sharply_drop_customers_this_month(df: pd.DataFrame):
            # set month as index
            df.set_index('month', inplace=True)
            df.drop(columns=['year', '其他'], inplace=True)
            df.replace(0, np.nan, inplace=True)

            df_mom = df.apply(lambda x: (x / x.shift(periods=1)) - 1, axis=0)

            cus = []
            for col in df_mom.columns:
                if df_mom.iloc[-1][col] < -0.2:
                    cus.append(col)

            sales = []
            for col in cus:
                sales.append(df[col].tolist())

            the_month = df.index.map(lambda x: x.strftime("%Y-%m")).tolist()
            line_chart = {
                'title': f"{the_month[-1]}订单环比大幅度下滑的供应商",
                'series-list': sales,
                'legend': cus,
                'x-axis': the_month
            }
            with self.__lock:
                self.content['sharplyDropSupplierRecentMonthLineChart'] = line_chart

        _sharply_drop_customers_this_month(df_sales_detail.copy())

    @staticmethod
    def _desc(_x):
        if _x >= 0:
            return '增长'
        else:
            return '减少'

    def _inactive_customer_recent(self):
        y = int(self.doc["attribute_month"][:4])
        df_recent = self._transpose_sales(y, "customerDetail_24")
        if df_recent is None or df_recent.empty:
            return None
        names = []
        for col in df_recent.columns[1:]:
            s = df_recent.head(3)[col].sum()
            if s == 0:
                names.append(col)

        summary = "近三个月不活跃的客户有" + ",".join(names)
        self.content['inactiveCustomerRecent3Month'] = {
            "summary": summary,
            "data": names
        }

    def _costumer_analysis_last3years(self):
        year = self.doc["attribute_month"][:4]
        _df = pd.DataFrame(self.content['customerDetail_24'])
        if _df.empty:
            return
        _df = _df.loc[:, ~_df.apply(lambda x: (x == '-').all())]
        _df = _df[~_df["PURCHASER_NAME"].isin(["其他", "总计"])]
        _df = self._convert_to_numeric(_df)

        months = sorted([int(str(m).replace("YF_", "")) for m in _df.columns.tolist() if str(m).startswith("YF")])
        cols = ["YF_" + str(m) for m in months[max(-4, -len(months)):]]

        for col in cols:
            _df[col + "_pct"] = _df[col] / _df[col].sum()
        _df['ZJE_pct'] = _df['ZJE'] / _df['ZJE'].sum()

        for i in range(1, len(cols)):
            last_col = cols[i - 1]
            col = cols[i]
            _df[col + "_pct_ratio"] = _df[col + "_pct"] / _df[last_col + "_pct"].replace(0, np.nan) - 1

        _df = _df.replace([np.inf, -np.Inf, np.nan], None)

        res = {}
        _df: pd.DataFrame
        for i, s in _df.iterrows():
            cat = s["PURCHASER_NAME"]
            item = [s[k + "_pct"] for k in cols[1:]]
            rate = [s[k + "_pct_ratio"] for k in cols[1:]]
            rate.insert(0, None)
            item.insert(0, s['ZJE_pct'])
            col = [(str(year) + "-" + k.replace("YF_", "")) for k in cols[1:]]
            col.insert(0, str(year))

            res[cat] = {
                'bar': item,
                'trend': rate,
                'col': col
            }
        self.content['customerGraph1'] = res

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

    # async def _get_tags_by_name(self, name: str, tp: datetime.datetime) -> tuple[str, Optional[dict]]:
    #     ident = await self.repo.dw_data_v3.get_usc_id_by_name(name)
    #     if ident.success and ident.data.usc_id:
    #         req = self.repo.dw_data_v3.new_pb_req(
    #             usc_id=ident.data.usc_id,
    #             time_point=self.time_point
    #         )
    #
    #
    #
    #
    #     else:
    #         return name, None

    async def _subject_company_related_info(self):
        try:
            usc_id = self.content["businessInfo"]["TYSHXYDM"]
            assert type(usc_id) is str
        except Exception:
            return
        # _, _related = await self.repo.dw_data.get_related(usc_id)

        req = self.repo.dw_data_v3.new_pb_req(
            usc_id=usc_id,
            time_point=self.time_point
        )

        shareholder, branches, investment, equity, equity_conclusion = await asyncio.gather(
            self.repo.dw_data_v3.get_ent_shareholders(req),
            self.repo.dw_data_v3.get_ent_branches(req),
            self.repo.dw_data_v3.get_ent_investment(req),
            self.repo.dw_data_v3.get_ent_equity_transparency(req),
            self.repo.dw_data_v3.get_ent_equity_transparency_conclusion(req),
        )

        self.content["shareholderData"] = shareholder.get("data") if shareholder else None
        self.content["branchesData"] = branches.get("data") if branches else None
        self.content["investmentData"] = investment.get("data") if investment else None
        self.content["equityTransparency"] = equity.get("data") if equity else None
        self.content["equityConclusion"] = equity_conclusion

    def _assign_tag(self, k1: str, k2: str, k_name: str, tag_dict: dict[str, str]):
        if (items := self.content.get(k1, {}).get(k2)) is not None:
            df_ = pd.DataFrame(items)
            df_["QYBQ"] = df_[k_name].apply(lambda x: tag_dict.get(x) if x not in ["合计", "总计"] else None)
            self.content[k1][k2] = df_.to_dict(orient="records")

    def _major_commodity_proportion(self, key_item: str, key_trades: str):
        record = self.content.get(key_item, {}).get(key_trades)
        if isinstance(record, dict):
            _df_trades = pd.DataFrame([record])
        else:
            _df_trades = pd.DataFrame(record)
        if _df_trades.empty:
            return
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
        self.content[key_item][key_trades] = _df_trades.to_dict(orient='records')

    def _subj_product_proportion(self):
        _df_selling_sta = self._get_buff_df("sellingSTA36")
        _df_selling_sta = _df_selling_sta[~_df_selling_sta['SSSPDL'].isin(['合计', '其他'])]
        _df_selling_sta['category'] = _df_selling_sta['SSSPXL'].apply(lambda x: x.split('*')[1] if isinstance(x, str) else None)
        _df_selling_sta['JEZB'] = _df_selling_sta['JEZB'].apply(lambda x: self._to_numeric(x)).astype(float)
        _df_selling_sta['proportion'] = _df_selling_sta.groupby('category')['JEZB'].transform('sum').astype(str) + '%'

        _df_selling_sta['ssspxl_last'] = _df_selling_sta['SSSPXL'].apply(lambda x: x.split('*')[-1] if isinstance(x, str) else None)
        _df_selling_sta['jezb_bracket'] = '(' + _df_selling_sta['JEZB'].astype(str) + '%' + ')'
        _df_selling_sta = _df_selling_sta.sort_values(by='JEZB', ascending=False)
        _df_selling_sta['category_detail'] = _df_selling_sta.groupby('category')['ssspxl_last'].transform(
            lambda x: ','.join(x + _df_selling_sta.loc[x.index, 'jezb_bracket']))

        result = _df_selling_sta[['category', 'proportion', 'category_detail']].drop_duplicates()

        result['proportion_sort'] = result['proportion'].str.replace('%', '').astype(float)
        result = result.sort_values(by='proportion_sort', ascending=False).drop('proportion_sort', axis=1)
        result['proportion'] = result['proportion'].apply(lambda x: str(round(self._to_numeric(x), 2)) + '%')
        self.content['subjectCompanyTags']['productProportion'] = result.to_dict("records")

    def _subj_purchase_proportion(self):
        def _safe_split(cat: str):
            try:
                return cat.split('*')[1]
            except Exception:
                return ""

        _df_purchase_sta = self._get_buff_df("purchaseSTA36")
        _df_purchase_sta = _df_purchase_sta[~_df_purchase_sta['SSSPDL'].isin(['合计', '其他'])]
        _df_purchase_sta['category'] = _df_purchase_sta['SSSPXL'].apply(lambda x: _safe_split(x))
        _df_purchase_sta['JEZB'] = _df_purchase_sta['JEZB'].apply(lambda x: self._to_numeric(x)).astype(float)
        _df_purchase_sta['proportion'] = _df_purchase_sta.groupby('category')['JEZB'].transform('sum').astype(str) + '%'

        _df_purchase_sta['ssspxl_last'] = _df_purchase_sta['SSSPXL'].apply(lambda x: x.split('*')[-1] if isinstance(x, str) else "")
        _df_purchase_sta['jezb_bracket'] = '(' + _df_purchase_sta['JEZB'].astype(str) + '%' + ')'
        _df_purchase_sta = _df_purchase_sta.sort_values(by='JEZB', ascending=False)
        _df_purchase_sta['category_detail'] = _df_purchase_sta.groupby('category')['ssspxl_last'].transform(
            lambda x: ','.join(x + _df_purchase_sta.loc[x.index, 'jezb_bracket']))

        result = _df_purchase_sta[['category', 'proportion', 'category_detail']].drop_duplicates()

        result['proportion_sort'] = result['proportion'].str.replace('%', '').astype(float)
        result = result.sort_values(by='proportion_sort', ascending=False).drop('proportion_sort', axis=1)
        result['proportion'] = result['proportion'].apply(lambda x: str(round(self._to_numeric(x), 2)) + '%')
        self.content['subjectCompanyTags']['purchaseProportion'] = result.to_dict("records")

    def _financial_summary(self):
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
                              f"{_pct_range_desc(avg_ratio_latest * 100)}。")
            if avg_ratio_senior is not None:
                _summary1 += f"{_s['SNRQ']}全年收入同比{avg_ratio_senior_str}, {_pct_range_desc(avg_ratio_senior * 100)}."
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
            f"负债率{rate_summary(debt_ratio_1, '负债率')}，"
            f"{rate_summary(debt_ratio_2, '负债率')}，"
            f"{rate_summary(debt_ratio_3, '负债率')}；,")
        _summary2 += f"{t1}至{t3}流动比例为{liquidity_ratio_1}，{liquidity_ratio_2}，{liquidity_ratio_3}，"
        _summary2 += f"{t1}至{t3}速动比例为{quick_ratio_2}，{quick_ratio_2}，{quick_ratio_3}，"
        _summary2 += f"{rate_summary(quick_ratio_3, '速动比例')}；"
        _summary3 += f"{t1}至{t3}应收账款周转日为{receivable_turnover_days_1}，{receivable_turnover_days_2}，{receivable_turnover_days_3}，"

        try:
            rate_rec = str(round(receivable_1 / asset_total_1 * 100, 1)) + "%"
        except (TypeError, ZeroDivisionError):
            rate_rec = "-"

        try:
            rate_stor = str(round(storage_1 / asset_total_1 * 100, 1)) + "%"
        except (TypeError, ZeroDivisionError):
            rate_stor = "-"

        _summary3 += f"{t1}应收款周转速度{rate_summary(receivable_turnover_days_1, '应收款周转日')}，应收款金额为{receivable_1}占总资产约{rate_rec}；"
        # %s - %s存货周转天数分别是 % s， %s， %s； %s存货周转速度 % s，存货金额为 % s，占总资产约 % .2f %%；近期公司运营能力 % s。
        _summary3 += f"{t1}至{t3}存货周转天数分别是{storage_turnover_days_1}，{storage_turnover_days_2}，{storage_turnover_days_3}；"
        _summary3 += f"{t3}存货周转速度{rate_summary(storage_turnover_days_1, '库存周转日')}，存货金额为{storage_3}，占总资产约{rate_stor}；"

        try:
            _temp_ysk = rate_summary(receivable_turnover_days_1, '应收款周转日')
            _temp_ch = rate_summary(storage_turnover_days_1, '库存周转日')
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

    def _trading_summary(self):
        def __sub_rule(value1: float, value2: float) -> bool:
            value_ranges = [
                {"min": -math.inf, "max": -30},
                {"min": -30, "max": -10},
                {"min": -10, "max": 10},
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
                _s = f"{time1}年-{time2}年利润较高，备货很少"
            elif 10 < value <= 30:
                _s = f"{time1}年-{time2}年有部分利润，备货较少"
            elif -10 <= value <= 10:
                _s = f"{time1}年-{time2}年基本上按销定采"
            elif -30 <= value < -10:
                _s = f"{time1}年-{time2}年小幅度备货或者部分原材料呆滞，需关注库存风险"
            elif value < -30:
                _s = f"{time1}年-{time2}年大幅度备货或大部分原材料呆滞，需重点关注库存风险"
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
                    _s4 = f"{res[index]['year']}年利润较高，备货很少"
                elif 10 < reset_val <= 30:
                    _s4 = f"{res[index]['year']}年有部分利润，备货较少"
                elif -10 <= reset_val <= 10:
                    _s4 = f"{res[index]['year']}年基本上按销定采"
                elif -30 <= reset_val < -10:
                    _s4 = f"{res[index]['year']}年小幅度备货或者部分原材料呆滞，需关注库存风险"
                elif reset_val < -30:
                    _s4 = f"{res[index]['year']}年大幅度备货或大部分原材料呆滞，需重点关注库存风险"
            return _s4

        def _sub_summary_3(value) -> str:
            _s3 = ''
            if value > 20:
                _s3 = "大幅增长"
            elif 5 < value <= 20:
                _s3 = "小幅增长"
            elif -5 <= value <= 5:
                _s3 = "趋势平稳"
            elif -20 < value < -5:
                _s3 = "小幅下降"
            elif value < -20:
                _s3 = "大幅下降"
            return _s3

        def _sub_summary_2() -> str:
            # 求采购额的合计
            df_ = pd.DataFrame(self.content['ydcgqkSTA'])
            df_['month'] = pd.to_datetime(df_['LAST_24M'])
            df_['year'] = df_['month'].dt.year
            years = list(set(df_['year'].tolist()))
            years.sort(reverse=True)

            _summary_cg_sum = []

            for y in years:
                dfy = df_[df_['year'] == y]
                data_float = dfy["HJ_M"].astype(float)  # 字符串转为int
                cg_sum = round(data_float.sum()) if round(data_float.sum()) != 0 else '-'  # 采购额求和
                if cg_sum == '-':
                    continue

                if dfy['month'].dt.month.max() == 12:  # 判断最大月份是否为12月
                    _summary_cg_sum.append(f"{y}年采购金额为{cg_sum}万元")  # 最大月份为12月则显示全年，不显示具体月份
                else:
                    _max_month = dfy['month'].dt.month.max()  # 获取最大月份
                    _min_month = dfy['month'].dt.month.min()  # 获取最小月份
                    if _max_month == _min_month:  # 判断最大月份是否等于最小月份，等于则只显示1个月份，不等于则显示月份区间
                        _summary_cg_sum.append(f"{y}年{_max_month}月采购金额为{cg_sum}万元")
                    else:
                        _summary_cg_sum.append(f"{y}年{_min_month}月-{_max_month}月采购金额为{cg_sum}万元")

            _summary_cg_sum = ', '.join(_summary_cg_sum)

            # 找到最新年份的截止月份
            latest_year = df_[df_['month'].dt.year == df_['month'].dt.year.max()]
            latest_year_max_month = latest_year['month'].dt.month.max()

            # 对每年的数据进行截断，使其与最新年份的截止月份相匹配
            cut_df_month = df_[df_['month'].dt.month <= latest_year_max_month]

            # 计算最新年份和去年的数据
            latest_year_data = cut_df_month[cut_df_month['month'].dt.year == cut_df_month['month'].dt.year.max()]["HJ_M"].astype(float).sum()
            cut_last_year_data = cut_df_month[cut_df_month['month'].dt.year == cut_df_month['month'].dt.year.max() - 1]["HJ_M"].astype(float).sum()
            # 计算去去年的数据
            cut_last_two_year_data = cut_df_month[cut_df_month['month'].dt.year == cut_df_month['month'].dt.year.max() - 2]["HJ_M"].astype(float).sum()

            # 计算最新年份和去年的同比
            latest_on_last_growth = round(((latest_year_data - cut_last_year_data) / cut_last_year_data) * 100, 2)
            # 计算最新年份和去去年的对比
            latest_on_two_last_growth = round(
                ((latest_year_data - cut_last_two_year_data) / cut_last_two_year_data) * 100, 2)

            # 计算去年和去去年的数据
            last_year_data = df_[df_['month'].dt.year == df_['month'].dt.year.max() - 1]["HJ_M"].astype(float).sum()
            last_two_year_data = df_[df_['month'].dt.year == df_['month'].dt.year.max() - 2]["HJ_M"].astype(float).sum()
            # 计算去年和去去年的同比
            last_on_two_last_growth = round(((last_year_data - last_two_year_data) / last_two_year_data) * 100, 2)

            if latest_year_max_month == 12:  # 判断如果月份为12月，则不显示月份了
                _summary_cg_growth = f"其中{cut_df_month['month'].dt.year.max()}年同比{str(cut_last_year_data) + '万元' if cut_last_year_data <= 0 else str(latest_on_last_growth) + '%'}，{'大幅增长' if cut_last_year_data <= 0 else _sub_summary_3(latest_on_last_growth)}，相比{cut_df_month['month'].dt.year.max() - 2}年{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(latest_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(latest_on_two_last_growth)}，{cut_df_month['month'].dt.year.max() - 1}年同比{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(last_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(last_on_two_last_growth)}"
            else:
                _summary_cg_growth = f"其中{cut_df_month['month'].dt.year.max()}年{latest_year_max_month}月同比{str(cut_last_year_data) + '万元' if cut_last_year_data <= 0 else str(latest_on_last_growth) + '%'}，{'大幅增长' if cut_last_year_data <= 0 else _sub_summary_3(latest_on_last_growth)}，相比{cut_df_month['month'].dt.year.max() - 2}年{latest_year_max_month}月{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(latest_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(latest_on_two_last_growth)}，{cut_df_month['month'].dt.year.max() - 1}年同比{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(last_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(last_on_two_last_growth)}"

            _s2 = f"{_summary_cg_sum}；{_summary_cg_growth}。"
            return _s2

        def _sub_summary_1() -> str:
            # 求销售额的合计
            df_ = pd.DataFrame(self.content['ydxsqkDetail'])
            # ts = df_['MONTH'].apply(lambda x: x[:7])
            df_['month'] = pd.to_datetime(df_['MONTH'].apply(lambda x: x[:7]))
            df_['year'] = df_['month'].dt.year
            years = list(set(df_['year'].tolist()))
            years.sort(reverse=True)

            field_to_sum = 'FPKJSR'
            if all(value == '0' for value in df_[field_to_sum]):  # 如果"FPKJSR"全部为0，则使用"SBKJZSR"，反之使用"FPKJSR"
                field_to_sum = 'SBKJZSR'
            _summary_xs_sum = []

            for y in years:
                dfy = df_[df_['year'] == y]
                data_float = dfy[field_to_sum].astype(float)  # 字符串转为int
                xs_sum = round(data_float.sum()) if round(data_float.sum()) != 0 else '-'  # 销售额求和
                if xs_sum == '-':
                    continue

                if dfy['month'].dt.month.max() == 12:  # 判断最大月份是否为12月
                    _summary_xs_sum.append(f"{y}年销售金额为{xs_sum}万元")  # 最大月份为12月则显示全年，不显示具体月份
                else:
                    _max_month = dfy['month'].dt.month.max()  # 获取最大月份
                    _min_month = dfy['month'].dt.month.min()  # 获取最小月份
                    if _max_month == _min_month:  # 判断最大月份是否等于最小月份，等于则只显示1个月份，不等于则显示月份区间
                        _summary_xs_sum.append(f"{y}年{_max_month}月销售金额为{xs_sum}万元")
                    else:
                        _summary_xs_sum.append(f"{y}年{_min_month}月-{_max_month}月销售金额为{xs_sum}万元")

            _summary_xs_sum = ', '.join(_summary_xs_sum)

            # 找到最新年份的截止月份
            latest_year = df_[df_['month'].dt.year == df_['month'].dt.year.max()]
            latest_year_max_month = latest_year['month'].dt.month.max()

            # 对每年的数据进行截断，使其与最新年份的截止月份相匹配
            cut_df_month = df_[df_['month'].dt.month <= latest_year_max_month]

            # 计算最新年份和去年的数据
            latest_year_data = cut_df_month[cut_df_month['month'].dt.year == cut_df_month['month'].dt.year.max()][field_to_sum].astype(float).sum()
            cut_last_year_data = cut_df_month[cut_df_month['month'].dt.year == cut_df_month['month'].dt.year.max() - 1][field_to_sum].astype(float).sum()
            # 计算去去年的数据
            cut_last_two_year_data = cut_df_month[cut_df_month['month'].dt.year == cut_df_month['month'].dt.year.max() - 2][field_to_sum].astype(float).sum()

            # 计算最新年份和去年的同比
            latest_on_last_growth = round(((latest_year_data - cut_last_year_data) / cut_last_year_data) * 100, 2)
            # 计算最新年份和去去年的对比
            latest_on_two_last_growth = round(((latest_year_data - cut_last_two_year_data) / cut_last_two_year_data) * 100, 2)

            # 计算去年和去去年的数据
            last_year_data = df_[df_['month'].dt.year == df_['month'].dt.year.max() - 1][field_to_sum].astype(float).sum()
            last_two_year_data = df_[df_['month'].dt.year == df_['month'].dt.year.max() - 2][field_to_sum].astype(float).sum()
            # 计算去年和去去年的同比
            last_on_two_last_growth = round(((last_year_data - last_two_year_data) / last_two_year_data) * 100, 2)

            if latest_year_max_month == 12:  # 判断如果月份为12月，则不显示月份了
                _summary_xs_growth = f"其中{cut_df_month['month'].dt.year.max()}年同比{str(cut_last_year_data) + '万元' if cut_last_year_data <= 0 else str(latest_on_last_growth) + '%'}，{'大幅增长' if cut_last_year_data <= 0 else _sub_summary_3(latest_on_last_growth)}，相比{cut_df_month['month'].dt.year.max() - 2}年{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(latest_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(latest_on_two_last_growth)}，{cut_df_month['month'].dt.year.max() - 1}年同比{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(last_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(last_on_two_last_growth)}"
            else:
                _summary_xs_growth = f"其中{cut_df_month['month'].dt.year.max()}年{latest_year_max_month}月同比{str(cut_last_year_data) + '万元' if cut_last_year_data <= 0 else str(latest_on_last_growth) + '%' }，{'大幅增长' if cut_last_year_data <= 0 else _sub_summary_3(latest_on_last_growth)}，相比{cut_df_month['month'].dt.year.max() - 2}年{latest_year_max_month}月{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(latest_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(latest_on_two_last_growth)}，{cut_df_month['month'].dt.year.max() - 1}年同比{str(cut_last_two_year_data) + '万元' if cut_last_two_year_data <= 0 else str(last_on_two_last_growth) + '%'}，{'大幅增长' if cut_last_two_year_data <= 0 else _sub_summary_3(last_on_two_last_growth)}"

            _s1 = f"{_summary_xs_sum}；{_summary_xs_growth}。"
            return _s1

        # 求 差值/销售金额
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
            merged['FPKJSR'] == 0,
            (merged['SBKJZSR'] - merged['HJ_M']) / merged['SBKJZSR'] * 100,
            (merged['FPKJSR'] - merged['HJ_M']) / merged['FPKJSR'] * 100

        )
        # Select the necessary columns
        result = merged[['year', 'result']]
        # Sort by year in descending order
        result = result.sort_values(by='year', ascending=False)
        res = result.to_dict(orient='records')
        self.content["ydxsqkInfo"]["czxsPercentage"] = res

        try:
            xs_summary = _sub_summary_1()
            cg_summary = _sub_summary_2()
        except Exception as e:
            raise e

        summaries1 = [
            xs_summary,
            cg_summary,
        ]

        if __sub_rule(res[0]['result'], res[1]['result']):
            if __sub_rule(res[1]['result'], res[2]['result']):
                try:
                    sepcial_rule1 = f"{_sub_summary_4(res, 0, 1)}。"
                except Exception as e:
                    raise e
                summaries2 = [
                    sepcial_rule1,
                ]
            else:
                try:
                    sepcial_rule1 = f"{_sub_summary_4(res, 0, 3)}，{_sub_summary_4(res, 2, 4)}。"
                except Exception as e:
                    raise e
                summaries2 = [
                    sepcial_rule1,
                ]
        elif __sub_rule(res[1]['result'], res[2]['result']):
                try:
                    sepcial_rule1 = f"{_sub_summary_4(res, 0, 4)}，{_sub_summary_4(res, 1, 2)}。"
                except Exception as e:
                    raise e
                summaries2 = [
                    sepcial_rule1,
                ]
        elif __sub_rule(res[0]['result'], res[2]['result']):
                try:
                    sepcial_rule1 = f"{_sub_summary_4(res, 0, 1)}，{_sub_summary_4(res, 1, 4)}。"
                except Exception as e:
                    raise e
                summaries2 = [
                    sepcial_rule1,
                ]
        else:
            try:
                sepcial_rule1 = f"{_sub_summary_4(res, 0, 4)}，{_sub_summary_4(res, 1, 4)}，{_sub_summary_4(res, 2, 4)}。"
            except Exception as e:
                raise e
            summaries2 = [
                sepcial_rule1,
            ]

        self.content['ydxscgSummary'] = {
            "S1": summaries1,
            "S2": summaries2,
        }

    def _revenue_detail_summary(self):
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

    @staticmethod
    def _to_numeric(x: str) -> Optional[float]:
        try:
            return float(x.replace(",", "").replace("%", ""))
        except Exception:
            return None

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

    @staticmethod
    def _convert_to_numeric(_df):
        for column in _df.columns:
            if not pd.api.types.is_numeric_dtype(_df[column]):
                try:
                    _df[column] = pd.to_numeric(_df[column].str.replace(',', ''), errors='raise')
                except ValueError:
                    pass
        return _df

    def _transpose_sales(self, year: int, k: str, name_col: str = 'PURCHASER_NAME') -> Optional[pd.DataFrame]:
        def _to_numeric(x):
            if isinstance(x, str):
                try:
                    return float(x.replace(',', ''))
                except ValueError:
                    return np.nan

        _df = pd.DataFrame(self.content[k])
        if _df.empty:
            return
        _df = _df.loc[:, ~_df.apply(lambda x: (x == '-').all())]
        _df = _df.drop(columns=['XH', 'ZJE'])
        _df.set_index(name_col, inplace=True)
        rename_col = {}
        for col in _df.columns:
            if col.startswith('YF_'):
                rename_col[col] = str(year) + '-' + col.replace('YF_', '')
        _df.rename(columns=rename_col, inplace=True)
        _df = _df.transpose()
        _df = _df.map(_to_numeric)
        _df.reset_index(inplace=True, drop=False)
        _df.columns.name = None
        _df.rename(columns={'index': 'month'}, inplace=True)
        _df['month'] = pd.to_datetime(_df['month'])
        _df.sort_values(by=['month'], inplace=True, ascending=False)
        _df.reset_index(inplace=True, drop=True)
        _df.drop(columns=['总计'], inplace=True)
        return _df

# if __name__ == '__main__':
#     ts = "1234"
#     print(ts[:6])