from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ContentVersion(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    CONTENT_V1: _ClassVar[ContentVersion]
    CONTENT_V2: _ClassVar[ContentVersion]

class ReportVersion(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    V2: _ClassVar[ReportVersion]
    V2_5: _ClassVar[ReportVersion]
    V3: _ClassVar[ReportVersion]
    LATEST: _ClassVar[ReportVersion]
CONTENT_V1: ContentVersion
CONTENT_V2: ContentVersion
V2: ReportVersion
V2_5: ReportVersion
V3: ReportVersion
LATEST: ReportVersion

class GetTradeDetailReq(_message.Message):
    __slots__ = ("content_id", "report_version", "option_time_period", "option_top_cus", "option_trade_frequency", "trade_type")
    class TimePeriodOption(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        PERIOD_1ST: _ClassVar[GetTradeDetailReq.TimePeriodOption]
        PERIOD_2ND: _ClassVar[GetTradeDetailReq.TimePeriodOption]
        PERIOD_3RD: _ClassVar[GetTradeDetailReq.TimePeriodOption]
        PERIOD_4TH: _ClassVar[GetTradeDetailReq.TimePeriodOption]
    PERIOD_1ST: GetTradeDetailReq.TimePeriodOption
    PERIOD_2ND: GetTradeDetailReq.TimePeriodOption
    PERIOD_3RD: GetTradeDetailReq.TimePeriodOption
    PERIOD_4TH: GetTradeDetailReq.TimePeriodOption
    class TradeFrequencyOption(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        FREQUENCY_LOW: _ClassVar[GetTradeDetailReq.TradeFrequencyOption]
        FREQUENCY_MID: _ClassVar[GetTradeDetailReq.TradeFrequencyOption]
        FREQUENCY_HIGH: _ClassVar[GetTradeDetailReq.TradeFrequencyOption]
    FREQUENCY_LOW: GetTradeDetailReq.TradeFrequencyOption
    FREQUENCY_MID: GetTradeDetailReq.TradeFrequencyOption
    FREQUENCY_HIGH: GetTradeDetailReq.TradeFrequencyOption
    class TopCusOption(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        TOP_1: _ClassVar[GetTradeDetailReq.TopCusOption]
        TOP_5: _ClassVar[GetTradeDetailReq.TopCusOption]
        TOP_10: _ClassVar[GetTradeDetailReq.TopCusOption]
        TOP_20: _ClassVar[GetTradeDetailReq.TopCusOption]
        ALL: _ClassVar[GetTradeDetailReq.TopCusOption]
    TOP_1: GetTradeDetailReq.TopCusOption
    TOP_5: GetTradeDetailReq.TopCusOption
    TOP_10: GetTradeDetailReq.TopCusOption
    TOP_20: GetTradeDetailReq.TopCusOption
    ALL: GetTradeDetailReq.TopCusOption
    class TradeType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        CUSTOMER: _ClassVar[GetTradeDetailReq.TradeType]
        SUPPLIER: _ClassVar[GetTradeDetailReq.TradeType]
    CUSTOMER: GetTradeDetailReq.TradeType
    SUPPLIER: GetTradeDetailReq.TradeType
    CONTENT_ID_FIELD_NUMBER: _ClassVar[int]
    REPORT_VERSION_FIELD_NUMBER: _ClassVar[int]
    OPTION_TIME_PERIOD_FIELD_NUMBER: _ClassVar[int]
    OPTION_TOP_CUS_FIELD_NUMBER: _ClassVar[int]
    OPTION_TRADE_FREQUENCY_FIELD_NUMBER: _ClassVar[int]
    TRADE_TYPE_FIELD_NUMBER: _ClassVar[int]
    content_id: int
    report_version: ReportVersion
    option_time_period: _containers.RepeatedScalarFieldContainer[GetTradeDetailReq.TimePeriodOption]
    option_top_cus: GetTradeDetailReq.TopCusOption
    option_trade_frequency: _containers.RepeatedScalarFieldContainer[GetTradeDetailReq.TradeFrequencyOption]
    trade_type: GetTradeDetailReq.TradeType
    def __init__(self, content_id: _Optional[int] = ..., report_version: _Optional[_Union[ReportVersion, str]] = ..., option_time_period: _Optional[_Iterable[_Union[GetTradeDetailReq.TimePeriodOption, str]]] = ..., option_top_cus: _Optional[_Union[GetTradeDetailReq.TopCusOption, str]] = ..., option_trade_frequency: _Optional[_Iterable[_Union[GetTradeDetailReq.TradeFrequencyOption, str]]] = ..., trade_type: _Optional[_Union[GetTradeDetailReq.TradeType, str]] = ...) -> None: ...

class GetTradeDetailResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _struct_pb2.Struct
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class GetAhpScoreReq(_message.Message):
    __slots__ = ("claim_id",)
    CLAIM_ID_FIELD_NUMBER: _ClassVar[int]
    claim_id: int
    def __init__(self, claim_id: _Optional[int] = ...) -> None: ...

class GetAhpScoreResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _struct_pb2.Struct
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class GetContentProcessReq(_message.Message):
    __slots__ = ("content_id", "report_version", "lang")
    CONTENT_ID_FIELD_NUMBER: _ClassVar[int]
    REPORT_VERSION_FIELD_NUMBER: _ClassVar[int]
    LANG_FIELD_NUMBER: _ClassVar[int]
    content_id: int
    report_version: ReportVersion
    lang: str
    def __init__(self, content_id: _Optional[int] = ..., report_version: _Optional[_Union[ReportVersion, str]] = ..., lang: _Optional[str] = ...) -> None: ...

class GetContentProcessResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _struct_pb2.Struct
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class GetContentValidateResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ...) -> None: ...

class GetJsonTranslateReq(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: _struct_pb2.Struct
    def __init__(self, data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class GetJsonTranslateResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _struct_pb2.Struct
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...
