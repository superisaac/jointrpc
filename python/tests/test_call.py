import pytest

from jointrpc.aio import Client
from jointrpc.message import Result, Error

@pytest.mark.asyncio
async def test_call_echo():
    c = Client('h2c://127.0.0.1:50055')
    res = await c.call('_echo', 'hello007')
    assert isinstance(res, Result)
    assert res.result['echo'] == 'hello007'

    res = await c.call('_echo', 2, 3)
    assert isinstance(res, Error)
    assert res.error['code'] == 10901
    assert res.error['reason'] == 'Validation Error: .params length of params mismatch'
