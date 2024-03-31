from logging import Logger
from typing import Optional
import asyncio
import aiomysql
from aiomysql.utils import _PoolContextManager


class SqlDb:
    def __init__(self, logger: Logger, **kwargs):
        self.db_config = kwargs
        self.logger: Logger = logger
        self.__pool: Optional[_PoolContextManager] = None

    @property
    def pool(self) -> _PoolContextManager:
        if not self.__pool:
            raise Exception("pool is not initialized")
        return self.__pool

    async def __aenter__(self):
        if not self.__pool:
            self.__pool = await aiomysql.create_pool(**self.db_config, loop=asyncio.get_event_loop())
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.__pool:
            self.__pool.close()
            await self.__pool.wait_closed()

    async def query_one(self, stmt: str) -> Optional[dict]:
        await self.__aenter__()
        async with self.pool.acquire() as conn:
            async with conn.cursor() as cur:
                await cur.execute(stmt)
                resp = await cur.fetchone()
                cols = [desc[0] for desc in cur.description]
                if resp:
                    return dict(zip(cols, resp))
                return None

    async def query_all(self, stmt: str) -> Optional[list[dict]]:
        await self.__aenter__()
        async with self.pool.acquire() as conn:
            async with conn.cursor() as cur:
                await cur.execute(stmt)
                resp = await cur.fetchall()
                cols = [desc[0] for desc in cur.description]
                if resp:
                    return [dict(zip(cols, row)) for row in resp]
                return None


