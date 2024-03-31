# import consul
# import socket
# import json
# from typing import Optional
# import platform
# from dns import resolver
#
# class Registry:
#     def __init__(self,
#                  addr: str,
#                  port: int,
#                  token: str,
#                  server_name: str,
#                  server_version: str,
#                  server_port: int):
#
#         # 初始化，连接consul服务器
#         self.__token = token
#         self._consul = consul.Consul(addr, str(port), token=self.__token)
#         self.server_name: str = server_name
#         self.server_version: str = server_version
#         self.server_port: int = server_port
#         self.__server_id: Optional[str] = None
#
#     @property
#     def server_id(self):
#         if self.__server_id is None:
#             self.__server_id = f"{self.server_name}-{self.server_version}-{platform.node()}-{self.get_lan_ip()}:{self.server_port}"
#         return self.__server_id
#
#     def register(self):
#         # 注册服务
#         self._consul.agent.service.register(
#             name=self.server_name,
#             service_id=self.server_id,
#             port=self.server_port,
#             tags=[f"version={self.server_version}"],
#             # 健康检查ip端口，检查时间：5,超时时间：30，注销时间：30s
#             check=consul.Check.tcp(
#                 host=self.get_lan_ip(),
#                 port=self.server_port,
#                 interval="5s",
#                 timeout="30s",
#                 deregister="30s"
#             ),
#         )
#
#     def deregister(self):
#         self._consul.agent.service.deregister(self.server_id, token=self.server_id)
#
#     def get_service(self, name):
#         services = self._consul.agent.services()
#         service = services.get(name)
#         if not service:
#             return None, None
#         addr = "{0}:{1}".format(service['Address'], service['Port'])
#         return service, addr
#
#     @staticmethod
#     def get_lan_ip():
#         try:
#             sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
#             # This doesn't even have to be a reachable address, as we're not trying to connect to it
#             sock.connect(('10.255.255.255', 1))
#             host_ip = sock.getsockname()[0]
#         except Exception:
#             host_ip = '127.0.0.1'
#         finally:
#             sock.close()
#         return host_ip
#
#
# if __name__ == '__main__':
#
#     # c = ConsulClient(
#     #     addr="192.168.44.169",
#     #     port=8500,
#     #     token="5a2af290-6fb4-23ff-66b8-1b9a8aa89b7f",
#     # )
#     # c.register_service(
#     #     name="pipeline.micros",
#     #     host=c.get_lan_ip(),
#     #     port=50053,
#     #     version="v1.0.0",
#     # )
#
#     # print(platform.node())
#
#
#     # host = "192.168.1.120"  # consul服务器的ip
#     # port = "8500"  # consul服务器对外的端口
#     # token = "xxxx"
#     # consul_client = ConsulClient(host, port, token)
#     #
#     # name = "loginservice"
#     # host = "192.168.1.10"
#     # port = 8500
#     # consul_client.RegisterService(name, host, port)
#     #
#     # check = consul.Check().tcp(host, port, "5s", "30s", "30s")
#     # print(check)
#     # res = consul_client.GetService("loginservice")
#     # print(res)
#     # print(socket.gethostname())
#     # print()
#
#     # def get_lan_ip():
#     #     try:
#     #         sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
#     #         # This doesn't even have to be a reachable address, as we're not trying to connect to it
#     #         sock.connect(('10.255.255.255', 1))
#     #         IP = sock.getsockname()[0]
#     #     except Exception:
#     #         IP = '127.0.0.1'
#     #     finally:
#     #         sock.close()
#     #     return IP
#     #
#     #
#     # print(get_lan_ip())
#
#
#
#
