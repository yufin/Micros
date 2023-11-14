from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Shareholders(_message.Message):
    __slots__ = ["shareholder_name", "shareholder_type", "capital_amount", "real_amount", "capital_type", "percent"]
    SHAREHOLDER_NAME_FIELD_NUMBER: _ClassVar[int]
    SHAREHOLDER_TYPE_FIELD_NUMBER: _ClassVar[int]
    CAPITAL_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    REAL_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    CAPITAL_TYPE_FIELD_NUMBER: _ClassVar[int]
    PERCENT_FIELD_NUMBER: _ClassVar[int]
    shareholder_name: str
    shareholder_type: str
    capital_amount: str
    real_amount: str
    capital_type: str
    percent: str
    def __init__(self, shareholder_name: _Optional[str] = ..., shareholder_type: _Optional[str] = ..., capital_amount: _Optional[str] = ..., real_amount: _Optional[str] = ..., capital_type: _Optional[str] = ..., percent: _Optional[str] = ...) -> None: ...

class GetShareholdersResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[Shareholders]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[Shareholders, _Mapping]]] = ...) -> None: ...

class Investment(_message.Message):
    __slots__ = ["enterprise_name", "operator", "shareholding_ratio", "invested_amount", "start_data", "status"]
    ENTERPRISE_NAME_FIELD_NUMBER: _ClassVar[int]
    OPERATOR_FIELD_NUMBER: _ClassVar[int]
    SHAREHOLDING_RATIO_FIELD_NUMBER: _ClassVar[int]
    INVESTED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    START_DATA_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    enterprise_name: str
    operator: str
    shareholding_ratio: str
    invested_amount: str
    start_data: str
    status: str
    def __init__(self, enterprise_name: _Optional[str] = ..., operator: _Optional[str] = ..., shareholding_ratio: _Optional[str] = ..., invested_amount: _Optional[str] = ..., start_data: _Optional[str] = ..., status: _Optional[str] = ...) -> None: ...

class GetInvestmentResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[Investment]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[Investment, _Mapping]]] = ...) -> None: ...

class Branches(_message.Message):
    __slots__ = ["enterprise_name", "operator", "area", "start_date", "status"]
    ENTERPRISE_NAME_FIELD_NUMBER: _ClassVar[int]
    OPERATOR_FIELD_NUMBER: _ClassVar[int]
    AREA_FIELD_NUMBER: _ClassVar[int]
    START_DATE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    enterprise_name: str
    operator: str
    area: str
    start_date: str
    status: str
    def __init__(self, enterprise_name: _Optional[str] = ..., operator: _Optional[str] = ..., area: _Optional[str] = ..., start_date: _Optional[str] = ..., status: _Optional[str] = ...) -> None: ...

class GetBranchesResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[Branches]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[Branches, _Mapping]]] = ...) -> None: ...

class GetEquityTransparencyResp(_message.Message):
    __slots__ = ["success", "code", "msg", "conclusion", "data", "usc_id"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    CONCLUSION_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    conclusion: str
    data: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    usc_id: str
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., conclusion: _Optional[str] = ..., data: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ..., usc_id: _Optional[str] = ...) -> None: ...

class GetEntIdentReq(_message.Message):
    __slots__ = ["enterprise_name"]
    ENTERPRISE_NAME_FIELD_NUMBER: _ClassVar[int]
    enterprise_name: str
    def __init__(self, enterprise_name: _Optional[str] = ...) -> None: ...

class EntIdent(_message.Message):
    __slots__ = ["exists", "isLegal", "usc_id"]
    EXISTS_FIELD_NUMBER: _ClassVar[int]
    ISLEGAL_FIELD_NUMBER: _ClassVar[int]
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    exists: bool
    isLegal: bool
    usc_id: str
    def __init__(self, exists: bool = ..., isLegal: bool = ..., usc_id: _Optional[str] = ...) -> None: ...

class GetEntIdentResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: EntIdent
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[EntIdent, _Mapping]] = ...) -> None: ...

class GetEntInfoResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: EntInfo
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Union[EntInfo, _Mapping]] = ...) -> None: ...

class EntInfo(_message.Message):
    __slots__ = ["usc_id", "enterprise_title", "enterprise_title_en", "business_registration_number", "establish_date", "region", "approved_date", "registered_address", "registered_capital", "paid_in_capital", "enterprise_type", "stuff_size", "stuff_insured_number", "business_scope", "import_export_qualification_code", "legal_representative", "registration_authority", "registration_status", "taxpayer_qualification", "organization_code", "url_qcc", "url_homepage", "business_term_start", "business_term_end", "id", "created_at", "updated_at"]
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    ENTERPRISE_TITLE_FIELD_NUMBER: _ClassVar[int]
    ENTERPRISE_TITLE_EN_FIELD_NUMBER: _ClassVar[int]
    BUSINESS_REGISTRATION_NUMBER_FIELD_NUMBER: _ClassVar[int]
    ESTABLISH_DATE_FIELD_NUMBER: _ClassVar[int]
    REGION_FIELD_NUMBER: _ClassVar[int]
    APPROVED_DATE_FIELD_NUMBER: _ClassVar[int]
    REGISTERED_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    REGISTERED_CAPITAL_FIELD_NUMBER: _ClassVar[int]
    PAID_IN_CAPITAL_FIELD_NUMBER: _ClassVar[int]
    ENTERPRISE_TYPE_FIELD_NUMBER: _ClassVar[int]
    STUFF_SIZE_FIELD_NUMBER: _ClassVar[int]
    STUFF_INSURED_NUMBER_FIELD_NUMBER: _ClassVar[int]
    BUSINESS_SCOPE_FIELD_NUMBER: _ClassVar[int]
    IMPORT_EXPORT_QUALIFICATION_CODE_FIELD_NUMBER: _ClassVar[int]
    LEGAL_REPRESENTATIVE_FIELD_NUMBER: _ClassVar[int]
    REGISTRATION_AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    REGISTRATION_STATUS_FIELD_NUMBER: _ClassVar[int]
    TAXPAYER_QUALIFICATION_FIELD_NUMBER: _ClassVar[int]
    ORGANIZATION_CODE_FIELD_NUMBER: _ClassVar[int]
    URL_QCC_FIELD_NUMBER: _ClassVar[int]
    URL_HOMEPAGE_FIELD_NUMBER: _ClassVar[int]
    BUSINESS_TERM_START_FIELD_NUMBER: _ClassVar[int]
    BUSINESS_TERM_END_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    usc_id: str
    enterprise_title: str
    enterprise_title_en: str
    business_registration_number: str
    establish_date: str
    region: str
    approved_date: str
    registered_address: str
    registered_capital: str
    paid_in_capital: str
    enterprise_type: str
    stuff_size: str
    stuff_insured_number: int
    business_scope: str
    import_export_qualification_code: str
    legal_representative: str
    registration_authority: str
    registration_status: str
    taxpayer_qualification: str
    organization_code: str
    url_qcc: str
    url_homepage: str
    business_term_start: str
    business_term_end: str
    id: int
    created_at: str
    updated_at: str
    def __init__(self, usc_id: _Optional[str] = ..., enterprise_title: _Optional[str] = ..., enterprise_title_en: _Optional[str] = ..., business_registration_number: _Optional[str] = ..., establish_date: _Optional[str] = ..., region: _Optional[str] = ..., approved_date: _Optional[str] = ..., registered_address: _Optional[str] = ..., registered_capital: _Optional[str] = ..., paid_in_capital: _Optional[str] = ..., enterprise_type: _Optional[str] = ..., stuff_size: _Optional[str] = ..., stuff_insured_number: _Optional[int] = ..., business_scope: _Optional[str] = ..., import_export_qualification_code: _Optional[str] = ..., legal_representative: _Optional[str] = ..., registration_authority: _Optional[str] = ..., registration_status: _Optional[str] = ..., taxpayer_qualification: _Optional[str] = ..., organization_code: _Optional[str] = ..., url_qcc: _Optional[str] = ..., url_homepage: _Optional[str] = ..., business_term_start: _Optional[str] = ..., business_term_end: _Optional[str] = ..., id: _Optional[int] = ..., created_at: _Optional[str] = ..., updated_at: _Optional[str] = ...) -> None: ...

class EntCredential(_message.Message):
    __slots__ = ["id", "usc_id", "certification_title", "certification_code", "certification_level", "certification_type", "certification_source", "certification_date", "certification_term_start", "certification_term_end", "certification_authority", "created_at", "updated_at"]
    ID_FIELD_NUMBER: _ClassVar[int]
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_TITLE_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_CODE_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_LEVEL_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_TYPE_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_SOURCE_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_DATE_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_TERM_START_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_TERM_END_FIELD_NUMBER: _ClassVar[int]
    CERTIFICATION_AUTHORITY_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    id: int
    usc_id: str
    certification_title: str
    certification_code: str
    certification_level: str
    certification_type: str
    certification_source: str
    certification_date: str
    certification_term_start: str
    certification_term_end: str
    certification_authority: str
    created_at: str
    updated_at: str
    def __init__(self, id: _Optional[int] = ..., usc_id: _Optional[str] = ..., certification_title: _Optional[str] = ..., certification_code: _Optional[str] = ..., certification_level: _Optional[str] = ..., certification_type: _Optional[str] = ..., certification_source: _Optional[str] = ..., certification_date: _Optional[str] = ..., certification_term_start: _Optional[str] = ..., certification_term_end: _Optional[str] = ..., certification_authority: _Optional[str] = ..., created_at: _Optional[str] = ..., updated_at: _Optional[str] = ...) -> None: ...

class GetEntCredentialResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[EntCredential]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[EntCredential, _Mapping]]] = ...) -> None: ...

class GetEntInfoReq(_message.Message):
    __slots__ = ["usc_id"]
    USC_ID_FIELD_NUMBER: _ClassVar[int]
    usc_id: str
    def __init__(self, usc_id: _Optional[str] = ...) -> None: ...

class GetEntStrArrayResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[str]] = ...) -> None: ...

class GetEntRankingListResp(_message.Message):
    __slots__ = ["success", "code", "msg", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MSG_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    code: int
    msg: str
    data: _containers.RepeatedCompositeFieldContainer[EnterpriseRankingList]
    def __init__(self, success: bool = ..., code: _Optional[int] = ..., msg: _Optional[str] = ..., data: _Optional[_Iterable[_Union[EnterpriseRankingList, _Mapping]]] = ...) -> None: ...

class EnterpriseRankingList(_message.Message):
    __slots__ = ["usc_id", "ranking_position", "list_title", "list_type", "list_source", "list_participants_total", "list_published_date", "list_url_qcc", "list_url_origin"]
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
