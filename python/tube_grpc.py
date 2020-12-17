# Generated by the Protocol Buffers compiler. DO NOT EDIT!
# source: tube.proto
# plugin: grpclib.plugin.main
import abc
import typing

import grpclib.const
import grpclib.client
if typing.TYPE_CHECKING:
    import grpclib.server

import tube_pb2


class JSONRPCTubeBase(abc.ABC):

    @abc.abstractmethod
    async def ListMethods(self, stream: 'grpclib.server.Stream[tube_pb2.ListMethodsRequest, tube_pb2.ListMethodsResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Call(self, stream: 'grpclib.server.Stream[tube_pb2.JSONRPCRequest, tube_pb2.JSONRPCResult]') -> None:
        pass

    @abc.abstractmethod
    async def Notify(self, stream: 'grpclib.server.Stream[tube_pb2.JSONRPCNotifyRequest, tube_pb2.JSONRPCNotifyResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Handle(self, stream: 'grpclib.server.Stream[tube_pb2.JSONRPCUpPacket, tube_pb2.JSONRPCDownPacket]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/JSONRPCTube/ListMethods': grpclib.const.Handler(
                self.ListMethods,
                grpclib.const.Cardinality.UNARY_UNARY,
                tube_pb2.ListMethodsRequest,
                tube_pb2.ListMethodsResponse,
            ),
            '/JSONRPCTube/Call': grpclib.const.Handler(
                self.Call,
                grpclib.const.Cardinality.UNARY_UNARY,
                tube_pb2.JSONRPCRequest,
                tube_pb2.JSONRPCResult,
            ),
            '/JSONRPCTube/Notify': grpclib.const.Handler(
                self.Notify,
                grpclib.const.Cardinality.UNARY_UNARY,
                tube_pb2.JSONRPCNotifyRequest,
                tube_pb2.JSONRPCNotifyResponse,
            ),
            '/JSONRPCTube/Handle': grpclib.const.Handler(
                self.Handle,
                grpclib.const.Cardinality.STREAM_STREAM,
                tube_pb2.JSONRPCUpPacket,
                tube_pb2.JSONRPCDownPacket,
            ),
        }


class JSONRPCTubeStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.ListMethods = grpclib.client.UnaryUnaryMethod(
            channel,
            '/JSONRPCTube/ListMethods',
            tube_pb2.ListMethodsRequest,
            tube_pb2.ListMethodsResponse,
        )
        self.Call = grpclib.client.UnaryUnaryMethod(
            channel,
            '/JSONRPCTube/Call',
            tube_pb2.JSONRPCRequest,
            tube_pb2.JSONRPCResult,
        )
        self.Notify = grpclib.client.UnaryUnaryMethod(
            channel,
            '/JSONRPCTube/Notify',
            tube_pb2.JSONRPCNotifyRequest,
            tube_pb2.JSONRPCNotifyResponse,
        )
        self.Handle = grpclib.client.StreamStreamMethod(
            channel,
            '/JSONRPCTube/Handle',
            tube_pb2.JSONRPCUpPacket,
            tube_pb2.JSONRPCDownPacket,
        )
