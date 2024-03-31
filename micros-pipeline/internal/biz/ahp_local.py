import pandas as pd
from typing import Optional
import pathlib
import math

df_ahp_scorecard: Optional[pd.DataFrame] = None


class AhpLocal:
    def __init__(self):
        self.__score_card: Optional[pd.DataFrame] = None

    def get_score_card(self, ent_type: int, level: int) -> Optional[pd.DataFrame]:
        global df_ahp_scorecard
        if self.__score_card is None:
            if df_ahp_scorecard is None:
                path_rules = pathlib.Path(__file__).parents[2].joinpath('static').joinpath('rc_scorecard_rule.csv')
                df_ = pd.read_csv(path_rules, encoding='utf8')
                df_['expr'] = df_['expr'].replace(' ', '')
                df_['score'] = df_["score"].astype(float)
                df_ahp_scorecard = df_
            self.__score_card = df_ahp_scorecard

        type_dict = {1: 470269276289576086, 2: 470269276289576086}
        scorecard_id = type_dict.get(ent_type, 470269276289576086)
        df_r = self.__score_card[(self.__score_card['scorecard_id'] == scorecard_id) & (self.__score_card['level'] == level)]
        return df_r.reset_index(drop=True)

    def scorecard_eval(self, factor: dict, ent_type: int, level: int) -> dict:
        df_sc = self.get_score_card(ent_type, level)
        df_f = pd.DataFrame.from_dict(factor, orient='index', columns=['value'])
        card = pd.merge(df_sc, df_f, left_on='field', right_index=True, how='left')
        result = card.apply(lambda x: self._pattern_eval(x), axis=1)
        card['result'] = result
        res = card[['agg_field', 'result']].groupby('agg_field').sum().to_dict()['result']
        return res

    @staticmethod
    def _pattern_eval(s: pd.Series) -> float:
        pattern = s['pattern']
        expr = s['expr']
        val = s['value']
        score = s['score']
        try:
            val = float(val)
        except (ValueError, TypeError):
            pass
        if pattern == 'inRange':
            # pattern = "(122,444]"
            min_ = float(expr[1:expr.find(',')])
            max_ = float(expr[expr.find(',') + 1:-1])
            try:
                if min_ < val <= max_:
                    return 1. * score
                else:
                    return 0
            except TypeError:
                return 0
        elif pattern == 'exact':
            try:
                expr = float(expr)
            except (ValueError, TypeError):
                pass

            if type(expr) is float:
                if type(val) == float:
                    if val == expr:
                        return 1.0 * score
                    else:
                        return 0
                else:
                    return 0
            elif type(expr) is str:
                if type(val) == str:
                    if val == expr:
                        return 1. * score
                    else:
                        return 0
                else:
                    return 0
            else:
                return 0
        elif pattern == 'default':
            if val is None:
                return 1. * score
            elif type(val) is not str:
                if math.isnan(val):
                    return 1. * score
                else:
                    return 0
            else:
                return 0
        elif pattern == 'weightAgg':
            if type(val) == float:
                weight = float(expr)
                return weight * val
            else:
                return 0

if __name__ == '__main__':
    from pprint import pprint
    # import json
    # from internal.biz.model.ahp_params import AhpParams
    # with open('../../static/input_data_sample.json', encoding='utf8') as f:
    #     factor = json.load(f)
    factor =  {'sw_sdsnb_cyrs': None, 'gs_gdct': 1, 'gs_gdwdx': 4, 'gs_frwdx': 4, 'lh_cylwz': 1, 'lh_md_ppjzl': 0, 'md_qybq': 0, 'sw_cwbb_yyzjzzts': None, 'sf_fh_ssqk_qy': 4, 'sw_jcxx_nsrxypj': 'B', 'zx_yhsxqk': 1, 'zx_dsfsxqk': 1, 'lh_qylx': 1, 'nsrsbh': '91440300350006672M', 'sw_sb_nsze_zzsqysds_12m': 277934.76, 'sw_sb_nszezzl_zzsqysds_12m_a': -0.6449, 'sw_sdsnb_gzxjzzjezzl': 0.8312, 'sw_sbzs_sflhypld_12m': 0.0008, 'sw_sdsnb_yjfy': 0.0, 'fp_jx_lxfy_12m': None, 'sw_cwbb_sszb': 1199000.0, 'fp_jy_sychje_zb_12m_lh': 0.1101, 'fp_jx_zyjyjezb_12m_lh': None, 'fp_xx_xychje_zb_12m_lh': 0.1207, 'fp_xx_zyjyjezb_12m_lh': 0.3999, 'sw_sb_qbxse_12m': 40189080.97, 'sw_sb_qbxsezzl_12m': -0.7098, 'sw_sb_lsxs_12m': 67.3291, 'sw_cwbb_chzzts_cb': 18.7082, 'sw_cwbb_zcfzl': 0.9205, 'sw_cwbb_mlrzzlv': 0.0779, 'sw_cwbb_jlrzzlv': 0.7281, 'sw_cwbb_jzcszlv': 0.1821, 'sw_jcxx_clnx': 8.18, 'apply_time': '2023-12-11 16:10:40', 'order_no': '45a6fb28-3b36-4152-8d1e-1ab20f404506'}

    # ahp_params = AhpParams(**factor)
    # print(ahp_params)
    se = AhpLocal()
    r = se.scorecard_eval(factor, 1, 1)
    pprint(r)
    # # print(len(r))
    #
    r2 = se.scorecard_eval(r, 1, 2)

    r3 = se.scorecard_eval(r2, 1, 3)
    #
    r4 = se.scorecard_eval(r3, 1, 4)
    #
    print(r2)
    print(r3)
    print(r4)
