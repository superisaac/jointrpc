from typing import Any, Union, Dict, Callable, Optional, Type, Tuple
import asyncio
import json
import uuid
import logging
from urllib.parse import urlparse
from grpclib.client import Channel, Stream
from grpclib.exceptions import StreamTerminatedError
import jointrpc.pb.jointrpc_pb2 as pb2    # type: ignore
import jointrpc.pb.jointrpc_grpc as grpc_srv  # type: ignore

from jointrpc.message import RPCError, Message, Request, Notify, Result, Error, parse, IdType

logger = logging.getLogger(__name__)

class RPCHandler:
    fn: Callable
    schema_json: str
    help_text: str

    def __init__(self, fn: Callable):
        self.fn = fn                 # type: ignore
        self.schema_json = ''
        self.help_text = ''

    def help(self, help_text: str) -> 'RPCHandler':
        self.help_text = help_text
        return self

    def schema(self, schema_json: str) -> 'RPCHandler':
        self.schema_json = schema_json
        return self

class Client:
    _server_url: str
    _channel: Channel
    _username: str
    _password: str
    handlers: Dict[str, 'RPCHandler']


    def __init__(self, server_url: str, **kwargs):
        self._server_url = server_url
        parsed = urlparse(self._server_url)
        self._username = parsed.username or ''
        self._password = parsed.password or ''
        host, p = parsed.netloc.split(':')
        port = int(p)
        self._channel = Channel(host, port, **kwargs)
        self.handlers = {}

    def close(self):
        self._channel.close()

    # handler methods
    def on(self, method: str, fn: Callable) -> 'RPCHandler':
        assert method not in self.handlers
        h = RPCHandler(fn)
        self.handlers[method] = h
        return h

    @property
    def stub(self):
        return grpc_srv.JointRPCStub(self._channel)

    @property
    def client_auth(self) -> Any:
        return pb2.ClientAuth(
            username=self._username,
            password=self._password)

    @property
    def auth(self) -> Tuple[str, str]:
        return (self._username, self._password)

    async def call(self, method: str,
                   *params: Any,
                   id: str='',
                   traceid: str='',
                   timeout: int=10,
                   broadcast: bool=False) -> Union[Result, Error]:
        if not id:
            id = uuid.uuid4().hex

        if not traceid:
            traceid = uuid.uuid4().hex

        reqmsg = Request(id, method,
                         *params,
                         traceid=traceid)
        envolope = pb2.JSONRPCEnvolope(
            body=json.dumps(reqmsg.encode()))

        req = pb2.JSONRPCCallRequest(
            auth=self.client_auth,
            envolope=envolope,
            broadcast=broadcast,
            timeout=timeout)

        res = await self.stub.Call(req)
        if res.status.code != 0:
            logger.error("error status of calling %s, status %s", method, res.status.code)
            raise ValueError("error status %s" % res.status.code)
        msg = parse(json.loads(res.envolope.body))
        if isinstance(msg, (Result, Error)):
            return msg
        else:
            logger.error("invalid result for call method %s, %s", method, res.body)
            raise ValueError("invalid result")

    async def notify(self, method: str,
                   *params: Any,
                   traceid: str='',
                   broadcast: bool=False) -> None:
        if not traceid:
            traceid = uuid.uuid4().hex

        reqmsg = Request(id, method, *params, traceid=traceid)
        envolope = pb2.JSONRPCEnvolope(
            body=json.dumps(reqmsg.encode()))
        req = pb2.JSONRPCNotifyRequest(
            auth=self.client_auth,
            envolope=envolope,
            broadcast=broadcast)

        res = await self.stub.Notify(req)
        if res.status.code != 0:
            logger.debug('error notify %s, status=%s', method, res.status.code)

    def methods(self):
        return [dict(
            name=m,
            help=h.help_text,
            schema=h.schema_json)
                for m, h in self.handlers.items()]

    async def declare_methods(self, stream):
        methods = [pb2.MethodInfo(
            name=m,
            help=h.help_text,
            schema_json=h.schema_json) for m, h in self.handlers.items()]

        req = pb2.DeclareMethodsRequest(
            request_id=uuid.uuid4().hex,
            methods=methods)
        uppac = pb2.JointRPCUpPacket(methods_request=req)
        await stream.send_message(uppac)

    def live_stream(self) -> 'LiveStream':
        return LiveStream(self)

class LiveStream:
    client: 'Client'
    stream: Optional[Stream]
    pending_requests: Dict[IdType, Callable]

    def __init__(self, client: 'Client'):
        self.client = client
        self.pending_requests = {}
        self.stream = None
        self._ready_cb = None

    def set_ready_cb(self, cb: Callable):
        self._ready_cb = cb

    async def close(self) -> None:
        logger.info('close stream %s', self.stream)
        if self.stream:
            await self.stream.cancel()
            self.stream = None

    async def authorize(self):
        username, password = self.client.auth
        await self.live_call('_stream.authorize',
                             username, password, self.declare_methods)

    async def ready(self):
        if self._ready_cb:
            await self._ready_cb()

    async def declare_methodfs(self):
        methods = self.client.methods()
        await self.live_call('_stream.declareMethods',
                             methods, self.ready)

    async def live_call(self, method: str,
                  *params: Any,
                  cb: Callable,
                  traceid: str='',
                  timeout: int=10) -> None:

        id = uuid.uuid4().hex

        if not traceid:
            traceid = uuid.uuid4().hex

        reqmsg = Request(id, method,
                         *params,
                         traceid=traceid)
        envo = pb2.JSONRPCEnvolope(
            body=json.dumps(reqmsg.encode()))

        assert self.stream is not None
        await self.stream.send_message(envo)

        self.pending_requests[id] = cb

    async def handle(self) -> None:
        async with self.client.stub.Live.open() as stream:
            self.stream = stream

            asyncio.ensure_future(self.authorize())

            while self.stream:
                try:
                    envo = await asyncio.wait_for(
                        self.stream.recv_message(),
                        timeout=1)
                except asyncio.TimeoutError:
                    continue
                except StreamTerminatedError:
                    logger.info('stream terminated')
                    break
                except RuntimeError as e:
                    logger.warning("wait for recv messages", exc_info=True)
                    return

                await self._handle_message(envo)

    async def _send_message(self, msg: Message):
        assert self.stream
        envo = pb2.JSONRPCEnvolope(
            body=json.dumps(msg.encode()))
        await self.stream.send_message(envo)

    async def _handle_message(self, envolope: pb2.JSONRPCEnvolope):
        assert self.stream
        msg = parse(json.loads(envolope.body))

        if isinstance(msg, (Request, Notify)):
            self._handle_request(msg)
        else:
            assert isinstance(msg, (Result, Error))
            self._handle_result(msg)

    async def _handle_request(self, reqmsg: Union[Request, Notify]):
        if reqmsg.method in self.client.handlers:
            h = self.client.handlers[reqmsg.method]
            try:
                result = await h.fn(reqmsg, *reqmsg.params)
            except RPCError as e:
                logger.warning("trace: %s",
                               reqmsg.traceid, exc_info=True)
                if isinstance(reqmsg, Request):
                    errmsg = e.error_message(
                        reqmsg.id, traceid=reqmsg.traceid)
                    await self._send_message(errmsg)
                    return
            except Exception as e:
                logger.error("trace: %s, error on request %s",
                             reqmsg.traceid, reqmsg.encode(),
                             exc_info=True)
                if isinstance(reqmsg, Request):
                    errmsg = Error(reqmsg.id,
                                   100,
                                   'server error',
                                   data=str(e),
                                   traceid=reqmsg.traceid)
                    await self._send_message(errmsg)
                await self.close()
                return

            if isinstance(reqmsg, Notify):
                return
            # send result back to stream
            resmsg = Result(reqmsg.id, result, traceid=reqmsg.traceid)
            return await self._send_message(resmsg)
        else:
            logger.warning("trace: %s, method not found %s ", reqmsg.traceid, reqmsg.encode())
            if isinstance(reqmsg, Request):
                errmsg = Error(reqmsg.id, -32601, 'method not found', traceid=reqmsg.traceid)
                return await self._send_message(errmsg)

    async def _handle_result(self, res: Union[Result, Error]):
        cb = self.pending_requests.get(res.id)
        if cb:
            del self.pending_requests[res.id]
            asyncio.ensure_future(cb(res))
        else:
            logger.warning("fail to find request for %s", res.encode())
