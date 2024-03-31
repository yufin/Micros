from enum import Enum


class CompanyDataItem(Enum):
    company_info: str = "企业基本信息"
    ranking_list: str = "排行榜"
    certification: str = "政府资质认证"
    industry: str = "所属行业"
    product: str = "产品信息"
    equity_transparency: str = "股权穿透"
    investment: str = "对外投资"
    shareholder: str = "股东信息"
    branches: str = "分支机构"
    # case_regis: str = "立案信息"




