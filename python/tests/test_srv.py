import pytest
import asyncio

from jointrpc.aio import Client

@pytest.mark.asyncio
async def cancel_handler(handler, sleep):
    await asyncio.sleep(sleep)
    await handler.close()

@pytest.mark.asyncio
async def run_call_client():
    await asyncio.sleep(0.1)
    c = Client('h2c://127.0.0.1:50055')
    r = await c.call('add2num', 5, 6)
    assert r.result == 12

@pytest.mark.asyncio
async def test_call_fn():
    c = Client('h2c://127.0.0.1:50055')

    async def add2num(msg, a, b):
        return a, b

    c.on('add2num', add2num)

    handler = c.live_stream()
    asyncio.ensure_future(cancel_handler(handler, 5))
    asyncio.ensure_future(run_call_client())
    #asyncio.ensure_future(handler.handle())

    #await run_call_client()

    await handler.handle()
