from typing import Any, Union
import json
import uuid
import logging
from urllib.parse import urlparse
from grpclib.client import Channel, Stream
import jointrpc.pb.jointrpc_pb2 as pb2
import jointrpc.pb.jointrpc_grpc as grpc_srv

from .message import Request, Notify, Result, Error, parse

logger = logging.getLogger(__name__)

class Client:
    _server_url: str
    _channel: Channel
    _username: str
    _password: str

    def __init__(self, server_url: str, **kwargs):
        self._server_url = server_url
        parsed = urlparse(self._server_url)
        self._username = parsed.username
        self._password = parsed.password
        host, p = parsed.netloc.split(':')
        port = int(p)
        self._channel = Channel(host, port, **kwargs)

    def close(self):
        self._channel.close()

    @property
    def client_auth(self) -> pb2.ClientAuth:
        return pb2.ClientAuth(
            username=self._username,
            password=self._password)

    async def call(self, method: str,
                   *params: Any,
                   id: str='',
                   trace_id: str='',
                   timeout: int=10,
                   broadcast: bool=False) -> Union[Result, Error]:
        if not id:
            id = uuid.uuid4().hex

        if not trace_id:
            trace_id = uuid.uuid4().hex

        reqmsg = Request(id, method, *params)
        envolope = pb2.JSONRPCEnvolope(
            body=json.dumps(reqmsg.encode()),
            trace_id=trace_id)
        req = pb2.JSONRPCCallRequest(
            auth=self.client_auth,
            envolope=envolope,
            broadcast=broadcast,
            timeout=timeout)

        srv = grpc_srv.JointRPCStub(self._channel)
        res = await srv.Call(req)
        msg = parse(json.loads(res.envolope.body))

        if isinstance(msg, Result) or isinstance(msg, Error):
            return msg
        else:
            logger.error("invalid result for call method %s, %s", method, res.body)
            raise ValueError("invalid result")



