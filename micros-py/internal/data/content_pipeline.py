from api.pipeline.v1 import pipeline_pb2
from typing import Union

class ContentPipeline:
    def __init__(self,
                 content: dict,
                 content_version: str,
                 report_version: Union[pipeline_pb2.REPORT_V2, pipeline_pb2.REPORT_V3, pipeline_pb2.REPORT_LATEST],
                 ):
        self.content = content
        self.content_version = content_version
        self.report_version = report_version

    def process(self):
        pass
