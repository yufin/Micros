from internal.data.mongo_db import MotorClient
from internal.data.dw_data_client_v3 import DwDataClientV3
from api.dwdata.v3 import dw_data_v3_pb2
from logging import Logger
import datetime
import pandas as pd
import json
from typing import Optional
import re


class DwStat:
    def __init__(self, motor_client: MotorClient, dw_data_client_v3: DwDataClientV3, logger: Logger):
        self.__motor_client = motor_client
        self._dw_data_v3 = dw_data_client_v3
        self.logger = logger

    async def stat_executive(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
        self.logger.info(f"stat_executive: {req.usc_id}")
        resp = await self._dw_data_v3.get_executive(req)
        if resp:
            s = self.__execution_summary(resp.get('data'))
            if isinstance(s, str):
                return s
        return None

    async def stat_his_executive(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
        self.logger.info(f"stat_his_executive: {req.usc_id}")
        resp = await self._dw_data_v3.get_historical_executive(req)
        if resp:
            s = self.__his_execution_summary(resp.get('data'))
            if isinstance(s, str):
                return s
        return None

    # async def stat_case(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
    #     self.logger.info(f"stat_case: {req.usc_id}")
    #     resp = await self._dw_data_v3.get_court_announcement(req)
    #     if resp:
    #         p = await self.__get_sub_id(req, "law_records_cases")
    #         if p:
    #             s = self.__case_summary(resp.get('data'), req.usc_id, p)
    #             if isinstance(s, str):
    #                 return s
    #     return None

    async def stat_case(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
        self.logger.info(f"stat_case: {req.usc_id}")
        resp = await self._dw_data_v3.get_court_announcement(req)
        if resp:
            s = await self.__case_summary(resp.get('data'), req.usc_id)
            if isinstance(s, str):
                return s
        return None

    async def stat_case_registration(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
        self.logger.info(f"stat_case_registration: {req.usc_id}")
        resp = await self._dw_data_v3.get_case_registration_info(req)
        if resp:
            _r = self.__new_case_summary(resp.get('data'))
            if isinstance(_r, str):
                return _r
        return None

    async def stat_financing(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq) -> Optional[str]:
        self.logger.info(f"stat_financing: {req.usc_id}")
        resp = await self._dw_data_v3.get_ent_financing(req)
        if resp:
            _r = await self.__financing_summary(resp.get("data"))
            if isinstance(_r, str):
                return _r
        return None

    async def __get_sub_id(self, req: dw_data_v3_pb2.GetDataBeforeTimePointReq, doc_name: str) -> Optional[str]:
        db = self.__motor_client.client['dw2']
        _doc = await db[doc_name].find_one(
            {
                "usc_id": req.usc_id,
                "update_date": {"$gte": req.time_point.ToDatetime()},
                "create_date": {"$lte": req.time_point.ToDatetime()},
            },
            {"_id": 0, "keyno": 1}
        )
        return _doc.get("keyno")


    async def __case_summary(self, data: list, usc_id: str) -> Optional[str]:
        """
        裁判文书总结
        """
        if len(data) > 0:

            db = self.__motor_client.client['spider']
            cur = db["basic_status_enterprise_info"].find({'usc_id': usc_id}, {'gsxx': 1, 'check_date': 1})
            res = []
            async for docs in cur:
                res.append(docs)

            if res:
                res.sort(key=lambda x: x['check_date'], reverse=True)
                enterprise_name = res[0]['gsxx']['enterprise_title']
            else:
                enterprise_name = None


            arr = []
            for j in range(len(data)):
                try:
                    v_temp = {
                        'enterprise_name': enterprise_name,
                        'usc_id': usc_id,
                        'CaseNo': data[j]['CaseNo'],
                        'CaseReason': data[j]['CaseReason'],
                        'CaseRoleJson': data[j]['CaseRoleJson'],
                        'CaseRole':json.dumps(data[j]['CaseRoleJson'],ensure_ascii=False),
                        'InvolvedAmount': data[j]['InvolvedAmount'],
                        'SubmitDate': datetime.datetime.fromtimestamp(data[j]['SubmitDate']),
                        'JudgeDate': datetime.datetime.fromtimestamp(data[j]['JudgeDate'])
                    }
                except Exception as e:
                    print(data[j])
                arr.append(v_temp)
            df = pd.DataFrame(arr)
        else:
            df = pd.DataFrame()

        if len(df) > 0:
            df = df.drop_duplicates(subset=['InvolvedAmount', 'CaseRole'], keep='first')
            dict_arr = df.to_dict(orient="records")

            for i in range(len(dict_arr)):
                case_role = []
                case_role_arr = json.loads(dict_arr[i]['CaseRole'])
                for j in range(len(case_role_arr)):
                    p_arr =[i for i in case_role_arr[j]['P'].split("，")]
                    if dict_arr[i]['enterprise_name'] in p_arr:
                        if '被' in case_role_arr[j]['R']:
                            case_role.append('被告')
                        elif '第三人' in case_role_arr[j]['R']:
                            case_role.append('第三人')
                        else:
                            case_role.append('原告')

                dict_arr[i]['role'] = '、'.join(case_role)

            df = pd.DataFrame(dict_arr)

            # 用0填充未知涉案金额
            df['InvolvedAmount'].fillna(value=0, inplace=True)

            def str_to_num(s):
                if s:
                    try:
                        num = float(s.replace(',', ''))
                    except Exception as e:
                        num = 0
                else:
                    num = 0
                return num

            # 将该列转换数据类型
            df['InvolvedAmount'] = df['InvolvedAmount'].apply(str_to_num)

            # 统计案由的分布
            result0 = df.groupby('CaseReason').groups
            # 找出前三数量级的案由
            result0_arr = []
            for k, v in result0.items():
                v_temp = {'CaseReason': k, 'count': len(v)}
                result0_arr.append(v_temp)

            result0_arr.sort(key=lambda x: x['count'], reverse=True)

            # 统计做原告、被告次数
            result1 = df.groupby('role').groups
            # 统计作为原告、被告涉及的总金额
            result2 = df.groupby('role')['InvolvedAmount'].sum()
            s1 = f'近三年案件总数量：{len(df)}件\n'
            s2 = '前三案由：'
            result0_arr = result0_arr[0:3]
            for _ in range(len(result0_arr)):
                s2 += f'{result0_arr[_]["CaseReason"]}（{result0_arr[_]["count"]}件）'

            s3 = ''
            for k, v in result1.items():
                s3 += f'{k}：{len(v)}件，金额：{"{:,}".format(float(result2[k]))}元\n'
            summary = s1 + s2 + '\n' + s3

        else:
            summary = (
                '近三年案件总数量：0件\n'
                '前三案由：其他执行（0件）其他民事（0件）劳动合同纠纷（0件）\n'
                '原告：0件，金额：0元\n'
                '被告：0件，金额：0.0元\n'
            )

        return summary

    # @staticmethod
    # def __case_summary(data: list, usc_id: str, key_no: str) -> Optional[str]:
    #     """
    #     裁判文书总结
    #     """
    #     if len(data) > 0:
    #         arr = []
    #         for j in range(len(data)):
    #             v_temp = {
    #                 'keyno': key_no,
    #                 'usc_id': usc_id,
    #                 'CaseNo': data[j]['CaseNo'],
    #                 'CaseReason': data[j]['CaseReason'],
    #                 'CaseRole': data[j]['CaseRole'],
    #                 'Court': data[j]['Court'],
    #                 'InvolvedAmount': data[j]['InvolvedAmount'],
    #                 'SubmitDate': datetime.datetime.fromtimestamp(data[j]['SubmitDate']),
    #                 'JudgeDate': datetime.datetime.fromtimestamp(data[j]['JudgeDate'])
    #             }
    #             arr.append(v_temp)
    #         df = pd.DataFrame(arr)
    #     else:
    #         df = pd.DataFrame()
    #
    #     if len(df) > 0:
    #         df = df.drop_duplicates(subset=['InvolvedAmount', 'CaseRole'], keep='first')
    #         dict_arr = df.to_dict(orient="records")
    #
    #         for i in range(len(dict_arr)):
    #             case_role = []
    #             case_role_arr = json.loads(dict_arr[i]['CaseRole'])
    #             for j in range(len(case_role_arr)):
    #                 if case_role_arr[j]['N'] == dict_arr[i]['keyno']:
    #                     if '被' in case_role_arr[j]['R']:
    #                         case_role.append('被告')
    #                     elif '第三人' in case_role_arr[j]['R']:
    #                         case_role.append('第三人')
    #                     else:
    #                         case_role.append('原告')
    #
    #             dict_arr[i]['role'] = '、'.join(case_role)
    #
    #         df = pd.DataFrame(dict_arr)
    #
    #         # 用0填充未知涉案金额
    #         df['InvolvedAmount'].fillna(value=0, inplace=True)
    #
    #         def str_to_num(s):
    #             if s:
    #                 num = float(s)
    #             else:
    #                 num = 0
    #             return num
    #
    #         # 将该列转换数据类型
    #         df['InvolvedAmount'] = df['InvolvedAmount'].apply(str_to_num)
    #
    #         # 统计案由的分布
    #         result0 = df.groupby('CaseReason').groups
    #         # 找出前三数量级的案由
    #         result0_arr = []
    #         for k, v in result0.items():
    #             v_temp = {'CaseReason': k, 'count': len(v)}
    #             result0_arr.append(v_temp)
    #
    #         result0_arr.sort(key=lambda x: x['count'], reverse=True)
    #
    #         # 统计做原告、被告次数
    #         result1 = df.groupby('role').groups
    #         # 统计作为原告、被告涉及的总金额
    #         result2 = df.groupby('role')['InvolvedAmount'].sum()
    #         s1 = f'近三年案件总数量：{len(df)}件\n'
    #         s2 = '前三案由：'
    #         result0_arr = result0_arr[0:3]
    #         for _ in range(len(result0_arr)):
    #             s2 += f'{result0_arr[_]["CaseReason"]}（{result0_arr[_]["count"]}件）'
    #
    #         s3 = ''
    #         for k, v in result1.items():
    #             s3 += f'{k}：{len(v)}件，金额：{"{:,}".format(float(result2[k]))}元\n'
    #         summary = s1 + s2 + '\n' + s3
    #
    #     else:
    #         summary = (
    #             '近三年案件总数量：0件\n'
    #             '前三案由：其他执行（0件）其他民事（0件）劳动合同纠纷（0件）\n'
    #             '原告：0件，金额：0元\n'
    #             '被告：0件，金额：0.0元\n'
    #         )
    #
    #     return summary

    @staticmethod
    def __execution_summary(data: list) -> Optional[str]:
        """
        被执行人总结
        :param json_record: mongo record
        :return: 执行案件总结
        """
        count, amount, last_exec_time, last_exec_amount = 0, 0, '', 0

        if isinstance(data, list):
            df = pd.DataFrame(data)
            if not df.empty:
                df.sort_values(by=['LiAnDate'], inplace=True)
                df.reset_index(drop=True, inplace=True)
                count = len(df)
                df['BiaoDi'] = df['BiaoDi'].apply(lambda x: str(x).replace(",", ""))
                df["BiaoDi"] = df["BiaoDi"].astype(float)
                amount = df["BiaoDi"].sum()
                last_exec_time = datetime.datetime.fromtimestamp(float(df.iloc[-1]["Liandate"])).strftime('%Y-%m-%d')
                last_exec_amount = float(df.iloc[-1]["BiaoDi"])
        else:
            return None

        return f'''当前：被执行次数：{count}，总金额：{"{:,}".format(round(amount, 0))}元，近一次被执行时间：{last_exec_time}，金额：{"{:,}".format(round(last_exec_amount, 0))}元'''


    @staticmethod
    def __his_execution_summary(data: list) -> Optional[str]:
        """
        被执行人总结
        :param json_record: mongo record
        :return: 执行案件总结
        """
        count, amount, last_exec_time, last_exec_amount = 0, 0, '', 0

        if isinstance(data, list):
            df = pd.DataFrame(data)
            if not df.empty:
                df.sort_values(by=['LiAnDate'], inplace=True)
                df.reset_index(drop=True, inplace=True)
                count = len(df)
                df["BiaoDi"] = df["BiaoDi"].astype(float)
                amount = df["BiaoDi"].sum()
                last_exec_time = datetime.datetime.fromtimestamp(float(df.iloc[-1]["Liandate"])).strftime('%Y-%m-%d')
                last_exec_amount = float(df.iloc[-1]["BiaoDi"])
        else:
            return None

        return f'''历史：被执行次数：{count}，总金额：{"{:,}".format(round(amount, 0))}元，近一次被执行时间：{last_exec_time}，金额：{"{:,}".format(round(last_exec_amount, 0))}元'''

    @staticmethod
    def __new_case_summary(data: list) -> Optional[str]:
        """
        立案信息总结
        :param json_record: mongo record
        :return: 总结
        """
        if data is None:
            return None

        if len(data) > 0:
            # 日期降序
            data.sort(key=lambda x: x['PunishDate'], reverse=True)

            # 最近一次立案信息
            # lastest_case_date = data[0]['CaseDate']

            # 最近一次立案时间
            # latest_court_time = str(common.timestamp_to_date(data[0]['PunishDate']))[:10]
            # latest_court_time = str(common.timestamp_to_date(data[0]['PunishDate']))[:10]
            latest_court_time = pd.to_datetime(data[0]["PunishDate"], unit="s")
            # pd.to_datetime()
            s = data[0]["PunishDate"]
            # 统计近三年案件
            years = [i['CourtYear'] for i in data]
            dict_ = {}
            for i in range(len(years)):
                if years[i] not in dict_:
                    dict_[years[i]] = 1
                else:
                    dict_[years[i]] += 1
            count = 0
            s = ''
            for k, v in dict_.items():
                s_ = f'-{k}：{v}件\n'
                s += s_
                count += 1

                if count > 3:
                    break
            summary = s + f'最近一次立案时间：{latest_court_time}'
            return summary


    # @staticmethod
    # def __new_case_summary(data: list) -> Optional[str]:
    #     """
    #     立案信息总结
    #     :param json_record: mongo record
    #     :return: 总结
    #     """
    #     if data is None:
    #         return None
    #
    #     if len(data) > 0:
    #         # 日期降序
    #         data.sort(key=lambda x: x['CaseDate'], reverse=True)
    #
    #         # 最近一次立案信息
    #         lastest_case_date = data[0]['CaseDate']
    #
    #         for i in range(len(data)):
    #
    #             if len(data[i]['CaseDate']) == 10:
    #                 data[i]['CaseYear'] = int(data[i]['CaseDate'].split('-')[0].strip())
    #             else:
    #                 data[i]['CaseYear'] = 1970
    #
    #         df = pd.DataFrame(data)
    #         # 立案年份分布-group by 分组
    #         result0 = df.groupby('CaseYear').groups
    #         year_now = int(str(datetime.datetime.now())[:4])
    #         years = [year_now, year_now - 1, year_now - 2]
    #         s_arr = []
    #         for y in years:
    #             try:
    #                 s = f'-{y}:{len(result0[y])}件'
    #             except Exception as e:
    #                 s = f'-{y}:0件'
    #             s_arr.append(s)
    #         summary = '\n'.join(s_arr) + '\n' + f'最近一次立案时间：{lastest_case_date}'
    #     else:
    #         year_now = int(str(datetime.datetime.now())[:4])
    #         years = [year_now, year_now - 1, year_now - 2]
    #         s_arr = []
    #         for y in years:
    #             s = f'-{y}:0件'
    #             s_arr.append(s)
    #         summary = '\n'.join(s_arr) + '\n' + f'最近一次立案时间：无'
    #     return summary

    async def __financing_summary(self, data: list):
        if len(data) == 0:
            return '无融资行为\n'

        # 汇率字典
        exchange_rate = {
            "人名币": 1,
            "美元": 7.2772,
            "日元": 0.0485,
            "英镑": 8.9771,
            "港元": 0.9305,
            "加元": 5.3089,
            "欧元": 7.7928,
            "韩元": 0.0056,
            "澳元": 4.7112
        }

        # 时间排序
        data.sort(key=lambda x: x['Date'], reverse=True)
        # 融资次数
        times = len(data)
        # 最新融资轮次
        latest_round = data[0]['Round']
        try:
            if re.search('[A-Z]+?轮', data[0]['Round']):
                s1 = f'融资总次数：{times} 次，最新融资轮次：{latest_round}\n'
            else:
                s1 = f'融资总次数：{times} 次\n'

        except Exception as e:
            s1 = f'融资总次数：{times} 次, 已上市\n'

        latest_amount = data[0]['FundedAmount']
        if re.search('[0-9]+', latest_amount):
            latest_num = float(re.search('[0-9.]+', latest_amount).group())
            rate = 1
            for k, v in exchange_rate.items():
                if k in latest_amount:
                    rate = v

            # 换成人名币
            latest_num *= rate
            # 换成万元
            if '亿' in latest_amount:
                latest_num *= 10000

            # 超过5位数换成亿
            if len(str(int(latest_num))) >= 5:
                latest_num_str = str(latest_num / 10000) + '亿人名币'
                s3 = f'''最近一次融资金额：约{latest_num_str}，日期：{data[0]['Date']}\n'''
            else:
                latest_num_str = str(latest_num) + '万人名币'
                s3 = f'''最近一次融资金额：约{latest_num_str}，日期：{data[0]['Date']}\n'''
        else:
            s3 = f'''最近一次融资金额：未披露，日期：{data[0]['Date']}\n'''

        # 融资总金额
        amount_total = 0

        # 近三年融资总金额
        amount_3 = []

        # 投资机构
        investors = []

        for i in range(len(data)):
            if isinstance(data[i]['Investors'], list):
                investors += data[i]['Investors']
            amount = data[i]['FundedAmount']
            if re.search('[0-9]+', amount):
                num = float(re.search('[0-9.]+', amount).group())
                rate = 1
                for k, v in exchange_rate.items():
                    if k in amount:
                        rate = v
                # 换算成万
                if '亿' in amount:
                    num = num * 10000
                # 换算成人名币
                num_rmb = num * rate
                amount_total += num_rmb

                # 近三年
                date_str = data[i]['Date']
                time_format = "%Y-%m-%d"
                date = datetime.datetime.strptime(date_str, time_format)
                if date >= datetime.datetime.now():
                    amount_3.append(num_rmb)
            # 未披露
            else:
                date_str = data[i]['Date']
                time_format = "%Y-%m-%d"
                date = datetime.datetime.strptime(date_str, time_format)
                if date >= datetime.datetime.now():
                    amount_3.append('未披露')

        # 融资总金额
        if len(str(int(amount_total))) >= 5:
            amount_total_str = str(amount_total / 10000) + '亿人名币'
            s2 = f'融资总金额：约{amount_total_str}以上（含未披露）\n'
        elif len(str(int(amount_total))) < 5 and len(str(int(amount_total))) > 1:
            amount_total_str = str(amount_total) + '万人名币'
            s2 = f'融资总金额：约{amount_total_str}以上（含未披露）\n'
        else:
            amount_total_str = '未披露'
            s2 = f'融资总金额：{amount_total_str}\n'

        # 近三年融资总额
        amount_3_num = [i for i in amount_3 if not isinstance(i, str)]
        amount_3_str = [i for i in amount_3 if isinstance(i, str)]
        if len(amount_3_num) > 0:
            if len(str(int(sum(amount_3_num)))) >= 5:
                amount_3_str = str(sum(amount_3_num) / 10000) + '亿人名币'
            else:
                amount_3_str = str(sum(amount_3_num)) + '万人名币'
            s4 = f'''近三年融资金额：约{amount_3_str}以上（含未披露）\n'''
        else:
            if len(amount_3_str) > 0:
                s4 = f'''近三年融资金额：未披露\n'''
            else:
                s4 = f'''近三年融资金额：无\n'''

        db = self.__motor_client.client['dw2']

        # 投资机构
        for i in range(len(investors)):
            investor_name = investors[i]['name']
            link = investors[i]['link']
            # investor_id = [i for i in link.split('/') if i][3].replace('.html', '')
            cur = db["qcc_investors"].find({'Alias': investor_name}, {})
            res = []
            async for docs in cur:
                res.append(docs)
            if len(res) > 0:
                count = res[0]['InvestCount']
            else:
                count = 1
            investors[i]['invest_count'] = count

        investors.sort(key=lambda x: x['invest_count'], reverse=True)
        investors_names = list([i['name'] for i in investors])
        new_investors_names = []
        for i in investors_names:
            if i not in new_investors_names:
                new_investors_names.append(i)

        if len(new_investors_names) >= 3:
            s5 = f'''投资机构（前三）：{new_investors_names[0]},{new_investors_names[1]},{new_investors_names[2]}'''
        elif len(new_investors_names) == 2:
            s5 = f'''投资机构（前三）：{new_investors_names[0]},{new_investors_names[1]}'''
        elif len(new_investors_names) == 1:
            s5 = f'''投资机构（前三）：{new_investors_names[0]}'''
        else:
            s5 = f'''投资机构（前三）：未披露'''

        summary = s1 + s2 + s3 + s4 + s5

        return summary


if __name__ == '__main__':
    from pprint import pprint
    from loguru import logger
    import asyncio

    dw = DwDataClientV3(logger=logger, target="192.168.44.150:50052", certs_folder="../../../certs")
    mgo_cli = MotorClient("mongodb://root:dev-password@192.168.44.169:27020", logger)

    ds = DwStat(motor_client=mgo_cli, logger=logger, dw_data_client_v3=dw)
    req = ds._dw_data_v3.new_pb_req(
        usc_id='915101007978308873',
        time_point=datetime.datetime.strptime('2023-07-01', '%Y-%m-%d')
    )

    async def test():
        return await asyncio.gather(
            ds.stat_financing(req)
        )

    r = asyncio.run(
        test()
    )
    pprint(r[0])
