import pandas as pd
from pydantic import BaseModel, Field
from typing import Union, Any, Optional


class NewCases(BaseModel):
    case_reason: str = Field(alias='CaseReason')
    illegal_type: str = Field(alias='IllegalType')


class RiskLevelCard(BaseModel):
    threshold: str
    risk_level: str


class ScoreCardL1(BaseModel):
    threshold: str
    risk_level: str
    score: int


def eval_cases_l1(cases_records: list[NewCases], score_card_l1: list[ScoreCardL1]) -> Optional[tuple[str, int]]:

    def _match_score(reason: str, df_score_card_l1_: pd.DataFrame, res_col: str) -> Optional[Any]:
        df_loc = df_score_card_l1_.loc[df_score_card_l1_['threshold'] == reason]
        if df_loc.empty:
            return None
        return df_score_card_l1_.loc[df_loc['score'].idxmax(), res_col]

    df = pd.DataFrame([c.model_dump() for c in cases_records])
    df_score_card_l1 = pd.DataFrame([s.model_dump() for s in score_card_l1])
    df['l1_score'] = df['case_reason'].apply(lambda x: _match_score(x, df_score_card_l1, 'score'))
    df['l1_risk_level'] = df['case_reason'].apply(lambda x: _match_score(x, df_score_card_l1, 'risk_level'))
    df_l1 = df[df['l1_score'].notnull()]
    if not df_l1.empty:
        s_l1_min = df_l1.loc[df_l1['l1_score'].idxmin()]
        return s_l1_min['l1_score'], s_l1_min['l1_risk_level']
    return None


def eval_cases_l2(case_records: list[NewCases], risk_level_card: list[RiskLevelCard]):

    def _eval_risk_level(reason: str, df_risk_level_card_: pd.DataFrame) -> Optional[str]:
        df_loc = df_risk_level_card_.loc[df_risk_level_card_['threshold'] == reason]
        if df_loc.empty:
            return None
        return df_loc.iloc[0]['risk_level']

    df = pd.DataFrame([c.model_dump() for c in case_records])
    df_risk_level_card = pd.DataFrame([s.model_dump() for s in risk_level_card])

    df['l2_risk_level'] = df['case_reason'].apply(lambda x: _eval_risk_level(x, df_risk_level_card))













if __name__ == '__main__':
    score_card = {
        "l1": [
            {
                "threshold": "司法拍卖",
                "risk_level": "高",
                "score": -10
            },
            {
                "threshold": "司法拍卖",
                "risk_level": "高",
                "score": -10
            },
            {
                "threshold": "买卖合同纠纷",
                "risk_level": "高",
                "score": -10
            },
        ]
    }
    cases_records = [
        {
            "IllegalType": "原告",
            "CaseReason": "买卖合同纠纷",
            "PunishDate": 1672675200
        },
        {
            "IllegalType": "原告",
            "CourtYear": "2023",
            "ExecuteGov": "上海市浦东新区人民法院",
            "CaseReason": "承揽合同纠纷",
        }
    ]

    risk_level_card = [
        {
            "threshold": "借款合同纠纷",
            "risk_level": "高"
        },
        {
            "threshold": "融资租赁合同纠纷",
            "risk_level": "高"
        }
    ]

    r1 = eval_cases_l1(
        [NewCases(**c) for c in cases_records],
        [ScoreCardL1(**d) for d in score_card['l1']],
        # [RiskLevelCard(**d) for d in risk_level_card]
    )
    print(r1)




    sc = [ScoreCardL1(**d) for d in score_card['l1']]
    df = pd.DataFrame([s.model_dump() for s in sc])
    df_ = df.loc[df['threshold'] == '司法拍卖']
    print(
        df_['score'].idxmin()
    )
    # print(
    #     df.loc[df.loc[df['threshold'] == '司法拍卖']['score'].idxmin(), 'scores']
    # )




    # d = ScoreCardL1(**{
    #             "threshold": "司法拍卖",
    #             "risk_level": "高",
    #             "score": -10
    #         })
    # print(d)
    # d = ScoreCardL1(threshold="w", risk_level="1", score=1)
    # print(d)