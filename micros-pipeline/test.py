from pydantic import BaseModel, Field


class FundRiskDetail(BaseModel):
    usc_id: str     # 企业的社会统一信用代码
    fund_risk_level: enumerate  # 企业资金风险加权等级
    shareholder_bg_score: int   # 企业股东背景加权得分
    shareholder_bg_type: enumerate   # 企业股东背景加权类型
    litigation_score: int   # 诉讼加权得分
    litigation_level: enumerate     # 诉讼风险加权等级
    honour_score: int       # 政府资质加权得分
    honour_level: enumerate     # 政府资质加权等级
    industry_rank: enumerate    # 行业地位
    industry_rank_score: enumerate    # 行业地位分数
    funding_score: int      # 融资加权得分
    funding_type: enumerate   # 融资加权类型
    ranking_score: int      # 榜单加权得分
    ranking_level: enumerate    # 榜单加权等级


class Result(BaseModel):
    obj_type: enumerate    # supplier,customer
    content_id: int     # 报告id
    usc_id: str         # 报告主体社会统一信用代码
    quality_res: enumerate  # 上/下游企业综合优质度结果
    cus_fund_risk_dist: dict[str, float]    # 上/下游综合企业资金风险分布占比
    cust_fund_risk_detail: list[FundRiskDetail]   # 上/下游每个企业资金风险详情






