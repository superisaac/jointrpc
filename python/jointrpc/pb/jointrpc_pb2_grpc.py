# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from jointrpc.pb import jointrpc_pb2 as jointrpc__pb2


class JointRPCStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Call = channel.unary_unary(
                '/JointRPC/Call',
                request_serializer=jointrpc__pb2.JSONRPCCallRequest.SerializeToString,
                response_deserializer=jointrpc__pb2.JSONRPCCallResult.FromString,
                )
        self.Notify = channel.unary_unary(
                '/JointRPC/Notify',
                request_serializer=jointrpc__pb2.JSONRPCNotifyRequest.SerializeToString,
                response_deserializer=jointrpc__pb2.JSONRPCNotifyResponse.FromString,
                )
        self.ListMethods = channel.unary_unary(
                '/JointRPC/ListMethods',
                request_serializer=jointrpc__pb2.ListMethodsRequest.SerializeToString,
                response_deserializer=jointrpc__pb2.ListMethodsResponse.FromString,
                )
        self.ListDelegates = channel.unary_unary(
                '/JointRPC/ListDelegates',
                request_serializer=jointrpc__pb2.ListDelegatesRequest.SerializeToString,
                response_deserializer=jointrpc__pb2.ListDelegatesResponse.FromString,
                )
        self.SubscribeState = channel.unary_stream(
                '/JointRPC/SubscribeState',
                request_serializer=jointrpc__pb2.AuthRequest.SerializeToString,
                response_deserializer=jointrpc__pb2.SubscribeStateResponse.FromString,
                )
        self.Worker = channel.stream_stream(
                '/JointRPC/Worker',
                request_serializer=jointrpc__pb2.JSONRPCEnvolope.SerializeToString,
                response_deserializer=jointrpc__pb2.JSONRPCEnvolope.FromString,
                )


class JointRPCServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Call(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Notify(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListMethods(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListDelegates(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SubscribeState(self, request, context):
        """state stream
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Worker(self, request_iterator, context):
        """request/response dual streams
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_JointRPCServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Call': grpc.unary_unary_rpc_method_handler(
                    servicer.Call,
                    request_deserializer=jointrpc__pb2.JSONRPCCallRequest.FromString,
                    response_serializer=jointrpc__pb2.JSONRPCCallResult.SerializeToString,
            ),
            'Notify': grpc.unary_unary_rpc_method_handler(
                    servicer.Notify,
                    request_deserializer=jointrpc__pb2.JSONRPCNotifyRequest.FromString,
                    response_serializer=jointrpc__pb2.JSONRPCNotifyResponse.SerializeToString,
            ),
            'ListMethods': grpc.unary_unary_rpc_method_handler(
                    servicer.ListMethods,
                    request_deserializer=jointrpc__pb2.ListMethodsRequest.FromString,
                    response_serializer=jointrpc__pb2.ListMethodsResponse.SerializeToString,
            ),
            'ListDelegates': grpc.unary_unary_rpc_method_handler(
                    servicer.ListDelegates,
                    request_deserializer=jointrpc__pb2.ListDelegatesRequest.FromString,
                    response_serializer=jointrpc__pb2.ListDelegatesResponse.SerializeToString,
            ),
            'SubscribeState': grpc.unary_stream_rpc_method_handler(
                    servicer.SubscribeState,
                    request_deserializer=jointrpc__pb2.AuthRequest.FromString,
                    response_serializer=jointrpc__pb2.SubscribeStateResponse.SerializeToString,
            ),
            'Worker': grpc.stream_stream_rpc_method_handler(
                    servicer.Worker,
                    request_deserializer=jointrpc__pb2.JSONRPCEnvolope.FromString,
                    response_serializer=jointrpc__pb2.JSONRPCEnvolope.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'JointRPC', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class JointRPC(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Call(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/JointRPC/Call',
            jointrpc__pb2.JSONRPCCallRequest.SerializeToString,
            jointrpc__pb2.JSONRPCCallResult.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Notify(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/JointRPC/Notify',
            jointrpc__pb2.JSONRPCNotifyRequest.SerializeToString,
            jointrpc__pb2.JSONRPCNotifyResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListMethods(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/JointRPC/ListMethods',
            jointrpc__pb2.ListMethodsRequest.SerializeToString,
            jointrpc__pb2.ListMethodsResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListDelegates(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/JointRPC/ListDelegates',
            jointrpc__pb2.ListDelegatesRequest.SerializeToString,
            jointrpc__pb2.ListDelegatesResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SubscribeState(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(request, target, '/JointRPC/SubscribeState',
            jointrpc__pb2.AuthRequest.SerializeToString,
            jointrpc__pb2.SubscribeStateResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Worker(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/JointRPC/Worker',
            jointrpc__pb2.JSONRPCEnvolope.SerializeToString,
            jointrpc__pb2.JSONRPCEnvolope.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
