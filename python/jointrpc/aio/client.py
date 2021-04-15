from typing import Any, Union, Dict, Callable, Optional
import asyncio
import json
import uuid
import logging
from urllib.parse import urlparse
from grpclib.client import Channel, Stream
from grpclib.exceptions import StreamTerminatedError
import jointrpc.pb.jointrpc_pb2 as pb2
import jointrpc.pb.jointrpc_grpc as grpc_srv

from jointrpc.message import RPCError, Message, Request, Notify, Result, Error, parse

logger = logging.getLogger(__name__)

class RPCHandler:
    fn: Callable
    schema_json: str
    help_text: str

    def __init__(self, fn: Callable):
        self.fn = fn
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
        self._username = parsed.username
        self._password = parsed.password
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
                   trace_id: str='',
                   broadcast: bool=False) -> None:
        if not trace_id:
            trace_id = uuid.uuid4().hex

        reqmsg = Request(id, method, *params)
        envolope = pb2.JSONRPCEnvolope(
            body=json.dumps(reqmsg.encode()),
            trace_id=trace_id)
        req = pb2.JSONRPCNotifyRequest(
            auth=self.client_auth,
            envolope=envolope,
            broadcast=broadcast)

        res = await self.stub.Notify(req)
        if res.status.code != 0:
            logger.debug('error notify %s, status=%s', method, res.status.code)

    async def declare_methods(self, conn_public_id: str):
        methods = [pb2.MethodInfo(
            name=m,
            help=h.help_text,
            schema_json=h.schema_json) for m, h in self.handlers.items()]

        req = pb2.DeclareMethodsRequest(
            auth=self.client_auth,
            conn_public_id=conn_public_id,
            methods=methods)
        resp = await self.stub.DeclareMethods(req)
        logger.info('declare methods response %s', resp)
        assert resp.status.code == 0

    def handler_stream(self) -> 'HandlerStream':
        return HandlerStream(self)

class HandlerStream:
    client: 'Client'
    stream: Optional[Stream]
    conn_public_id: str

    def __init__(self, client: 'Client'):
        self.client = client
        self.stream = None
        self.conn_public_id = ''


    async def close(self):
        logger.info('close stream %s', self.stream)
        if self.stream:
            await self.stream.cancel()
            self.stream = None

    async def handle(self) -> None:
        async with self.client.stub.Worker.open() as stream:
            self.stream = stream
            uppac = pb2.JointRPCUpPacket(auth=self.client.client_auth)
            await stream.send_message(uppac)

            while self.stream:
            #async for downpac in stream:
                try:
                    downpac = await asyncio.wait_for(self.stream.recv_message(), timeout=1)
                except asyncio.TimeoutError:
                    continue
                except StreamTerminatedError:
                    logger.info('stream terminated')
                    break
                except RuntimeError as e:
                    logger.warning("wait for recv messages", exc_info=True)
                    return

                if downpac.echo.conn_public_id:
                    if downpac.echo.status.code == 0:
                        self.conn_public_id = downpac.echo.conn_public_id
                        await self.client.declare_methods(self.conn_public_id)
                    else:
                        logger.warning("different status %s of payload echo", downpac.echo.status)
                elif downpac.ping.text:
                    logger.info("ping received")
                elif downpac.pong.text:
                    logger.info("pong received")
                elif downpac.state.methods:
                    logger.info("state change")
                elif downpac.envolope.body:
                    await self._handle_message(
                        downpac.envolope,
                        trace_id=downpac.envolope.trace_id)

    async def _send_message(self, msg: Message, trace_id:str=''):
        assert self.stream
        envo = pb2.JSONRPCEnvolope(
            body=json.dumps(msg.encode()),
            trace_id=trace_id)
        uppac = pb2.JointRPCUpPacket(envolope=envo)
        await self.stream.send_message(uppac)

    async def _handle_message(self, envolope: pb2.JSONRPCEnvolope, trace_id: str=''):
        assert self.stream
        reqmsg = parse(json.loads(envolope.body))
        assert isinstance(reqmsg, (Request, Notify))

        if reqmsg.method in self.client.handlers:
            h = self.client.handlers[reqmsg.method]
            try:
                result = await h.fn(reqmsg, *reqmsg.params)
            except RPCError as e:
                logger.warning("trace: %s", trace_id, exc_info=True)
                if isinstance(reqmsg, Request):
                    errmsg = e.error_message(reqmsg.id)
                    await self._send_message(errmsg, trace_id=trace_id)
                    return
            except Exception as e:
                logger.error("trace: %s, error on request %s", trace_id, reqmsg.encode(), exc_info=True)
                if isinstance(reqmsg, Request):
                    errmsg = Error.server_error(reqmsg.id, reason=str(e))
                    await self._send_message(errmsg, trace_id=trace_id)
                await self.close()
                return

            if isinstance(reqmsg, Notify):
                return
            # send result back to stream
            resmsg = Result(reqmsg.id, result)
            return await self._send_message(resmsg, trace_id=trace_id)
        else:
            logger.warning("trace: %s, method not found %s ", trace_id, reqmsg.encode())
            if isinstance(reqmsg, Request):
                errmsg = Error.method_not_found(reqmsg.id)
                return await self._send_message(errmsg, trace_id=trace_id)
