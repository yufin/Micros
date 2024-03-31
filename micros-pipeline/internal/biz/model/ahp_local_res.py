from pydantic import BaseModel
from internal.biz.model.ahp_params import AhpParams


class L1Factor(BaseModel):
    md_qybq: float
    lh_cylwz: float
    lh_md_ppjzl: float
    gs_gdct: float
    gs_gdwdx: float
    gs_frwdx: float
    sw_sdsnb_cyrs: float
    sw_cwbb_yyzjzzts: float
    sw_jcxx_nsrxypj: float
    sf_fh_ssqk_qy: float
    zx_yhsxqk: float
    zx_dsfsxqk: float
    sw_sb_nsze_zzsqysds_12m: float
    sw_sb_nszezzl_zzsqysds_12m_a: float
    sw_sdsnb_gzxjzzjezzl: float
    sw_sbzs_sflhypld_12m: float
    sw_sdsnb_yjfy: float
    sw_cwbb_sszb: float
    # fp_jx_lxfy_12m: float
    fp_jy_sychje_zb_12m_lh: float
    fp_jx_zyjyjezb_12m_lh: float
    fp_xx_xychje_zb_12m_lh: float
    fp_xx_zyjyjezb_12m_lh: float
    sw_sb_qbxse_12m: float
    sw_sb_qbxsezzl_12m: float
    sw_sb_lsxs_12m: float
    sw_cwbb_chzzts_cb: float
    sw_cwbb_zcfzl: float
    sw_cwbb_mlrzzlv: float
    sw_cwbb_jlrzzlv: float
    sw_cwbb_jzcszlv: float
    sw_jcxx_clnx: float
    lh_lxfybqbxse_12m: float


class L2Factor(BaseModel):
    enterprise_info: float
    equity_structure: float
    industry_chain_analyse: float
    supply_chain_analyse: float
    operation_ability: float
    debt_solvency: float
    profit_ability: float
    tax_statistic: float
    operation_cost: float
    finance_cost: float
    external_publicity_credit: float
    loan_granting_credit: float


class L3Factor(BaseModel):
    represent_basic: float
    operating: float
    financial: float
    credit: float


class L4Factor(BaseModel):
    score: float


class AhpResult(BaseModel):
    l0: AhpParams
    l1: L1Factor
    l2: L2Factor
    l3: L3Factor
    l4: L4Factor



