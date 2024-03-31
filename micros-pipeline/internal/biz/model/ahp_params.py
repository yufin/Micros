from pydantic import BaseModel
from typing import Optional, Union, get_type_hints
import numbers
import time
import uuid


class AhpParams(BaseModel):
    sw_sdsnb_cyrs: Optional[int] = None
    gs_gdct: Optional[int] = None
    gs_gdwdx: Optional[int] = None
    gs_frwdx: Optional[int] = None
    lh_cylwz: Optional[int] = None
    lh_md_ppjzl: Optional[int] = None
    md_qybq: Optional[int] = None
    sw_cwbb_yyzjzzts: Optional[int] = None
    sf_fh_ssqk_qy: Optional[int] = None
    sw_jcxx_nsrxypj: Optional[str] = None
    zx_yhsxqk: Optional[int] = None
    zx_dsfsxqk: Optional[int] = None
    lh_qylx: Optional[int] = None
    nsrsbh: Optional[str] = None
    sw_sb_nsze_zzsqysds_12m: Optional[float] = None
    sw_sb_nszezzl_zzsqysds_12m_a: Optional[float] = None
    sw_sdsnb_gzxjzzjezzl: Optional[float] = None
    sw_sbzs_sflhypld_12m: Optional[float] = None
    sw_sdsnb_yjfy: Optional[float] = None
    fp_jx_lxfy_12m: Optional[float] = None
    sw_cwbb_sszb: Optional[float] = None
    fp_jy_sychje_zb_12m_lh: Optional[float] = None
    fp_jx_zyjyjezb_12m_lh: Optional[float] = None
    fp_xx_xychje_zb_12m_lh: Optional[float] = None
    fp_xx_zyjyjezb_12m_lh: Optional[float] = None
    sw_sb_qbxse_12m: Optional[float] = None
    sw_sb_qbxsezzl_12m: Optional[float] = None
    sw_sb_lsxs_12m: Optional[float] = None
    sw_cwbb_chzzts_cb: Optional[float] = None
    sw_cwbb_zcfzl: Optional[float] = None
    sw_cwbb_mlrzzlv: Optional[float] = None
    sw_cwbb_jlrzzlv: Optional[float] = None
    sw_cwbb_jzcszlv: Optional[float] = None
    sw_jcxx_clnx: Optional[float] = None
    apply_time: Optional[str] = None
    order_no: Optional[str] = None
    lh_lxfybqbxse_12m: Optional[float] = None

    def __init__(self, **data):
        # Get the type hints from the class
        hints = get_type_hints(self.__class__)
        for name, value in data.items():
            # Get the type hint for this field
            hint = hints.get(name, None)
            # Unpack the hint if it's a Union
            if getattr(hint, '__origin__', None) is Union:
                hint = hint.__args__
            # If the field is supposed to be a numeric type and the value is a string
            if isinstance(hint, tuple) and any(issubclass(t, numbers.Number) for t in hint) and isinstance(value, str):
                try:
                    # Convert strings that represent integers to integers
                    if int in hint:
                        data[name] = int(value)
                    # Convert strings that represent floats to floats
                    elif float in hint:
                        data[name] = float(value)
                except ValueError:
                    data[name] = None
        super().__init__(**data)
        # current time zone
        self.apply_time = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
        # gen uuid
        self.order_no = str(uuid.uuid4())

