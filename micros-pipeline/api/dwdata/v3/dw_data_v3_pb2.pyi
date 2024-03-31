from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetDataDictResp(_message.Message):
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

class GetDataListResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data", "total")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    total: int
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class GetDataWithDurationReq(_message.Message):
    __slots__ = ("usc_id", "time_point", "validate_extend_date")
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    TIME_POINT_FIELD_NUMBER: _ClassVar[int]
    VALIDATE_EXTEND_DATE_FIELD_NUMBER: _ClassVar[int]
    usc_id: str
    time_point: _timestamp_pb2.Timestamp
    validate_extend_date: int
    def __init__(self, usc_id: _Optional[str] = ..., time_point: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., validate_extend_date: _Optional[int] = ...) -> None: ...

class GetDataBeforeTimePointReq(_message.Message):
    __slots__ = ("usc_id", "time_point", "page_size", "page_num")
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    TIME_POINT_FIELD_NUMBER: _ClassVar[int]
    PAGE_SIZE_FIELD_NUMBER: _ClassVar[int]
    PAGE_NUM_FIELD_NUMBER: _ClassVar[int]
    usc_id: str
    time_point: _timestamp_pb2.Timestamp
    page_size: int
    page_num: int
    def __init__(self, usc_id: _Optional[str] = ..., time_point: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ..., page_size: _Optional[int] = ..., page_num: _Optional[int] = ...) -> None: ...

class GetUscIdByEnterpriseNameReq(_message.Message):
    __slots__ = ("enterprise_name",)
    ENTERPRISE_NAME_FIELD_NUMBER: _ClassVar[int]
    enterprise_name: str
    def __init__(self, enterprise_name: _Optional[str] = ...) -> None: ...

class GetUscIdByEnterpriseNameResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data")
    class EntIdent(_message.Message):
        __slots__ = ("exists", "isLegal", "usc_id")
        EXISTS_FIELD_NUMBER: _ClassVar[int]
        ISLEGAL_FIELD_NUMBER: _ClassVar[int]
        USC_ID_FIELD_NUMBER: _ClassVar[int]
        exists: bool
        isLegal: bool
        usc_id: str
        def __init__(self, exists: bool = ..., isLegal: bool = ..., usc_id: _Optional[str] = ...) -> None: ...
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: GetUscIdByEnterpriseNameResp.EntIdent
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[GetUscIdByEnterpriseNameResp.EntIdent, _Mapping]] = ...) -> None: ...

class GetEntRankingListResp(_message.Message):
    __slots__ = ("success", "code", "msg", "data", "total")
    class EnterpriseRankingList(_message.Message):
        __slots__ = ("usc_id", "ranking_position", "list_title", "list_type", "list_source", "list_participants_total", "list_published_date", "list_url_qcc", "list_url_origin")
        USC_ID_FIELD_NUMBER: _ClassVar[int]
        RANKING_POSITION_FIELD_NUMBER: _ClassVar[int]
        LIST_TITLE_FIELD_NUMBER: _ClassVar[int]
        LIST_TYPE_FIELD_NUMBER: _ClassVar[int]
        LIST_SOURCE_FIELD_NUMBER: _ClassVar[int]
        LIST_PARTICIPANTS_TOTAL_FIELD_NUMBER: _ClassVar[int]
        LIST_PUBLISHED_DATE_FIELD_NUMBER: _ClassVar[int]
        LIST_URL_QCC_FIELD_NUMBER: _ClassVar[int]
        LIST_URL_ORIGIN_FIELD_NUMBER: _ClassVar[int]
        usc_id: str
        ranking_position: int
        list_title: str
        list_type: str
        list_source: str
        list_participants_total: int
        list_published_date: str
        list_url_qcc: str
        list_url_origin: str
        def __init__(self, usc_id: _Optional[str] = ..., ranking_position: _Optional[int] = ..., list_title: _Optional[str] = ..., list_type: _Optional[str] = ..., list_source: _Optional[str] = ..., list_participants_total: _Optional[int] = ..., list_published_date: _Optional[str] = ..., list_url_qcc: _Optional[str] = ..., list_url_origin: _Optional[str] = ...) -> None: ...
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[GetEntRankingListResp.EnterpriseRankingList]
    total: int
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[GetEntRankingListResp.EnterpriseRankingList, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...
