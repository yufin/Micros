from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ContentVersion(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    CONTENT_V1: _ClassVar[ContentVersion]
    CONTENT_V2: _ClassVar[ContentVersion]

class ReportVersion(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    V2: _ClassVar[ReportVersion]
    V3: _ClassVar[ReportVersion]
    LATEST: _ClassVar[ReportVersion]
CONTENT_V1: ContentVersion
CONTENT_V2: ContentVersion
V2: ReportVersion
V3: ReportVersion
LATEST: ReportVersion

class GetContentProcessReq(_message.Message):
    __slots__ = ["content_id", "report_version"]
    CONTENT_ID_FIELD_NUMBER: _ClassVar[int]
    REPORT_VERSION_FIELD_NUMBER: _ClassVar[int]
    content_id: int
    report_version: ReportVersion
    def __init__(self, content_id: _Optional[int] = ..., report_version: _Optional[_Union[ReportVersion, str]] = ...) -> None: ...

class GetContentProcessResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _struct_pb2.Struct
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...
