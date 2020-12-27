# Generated by the Protocol Buffers compiler. DO NOT EDIT!
# source: etcdclient/rpc.proto
# plugin: grpclib.plugin.main
import abc
import typing

import grpclib.const
import grpclib.client
if typing.TYPE_CHECKING:
    import grpclib.server

import etcdclient.kv_pb2
import etcdclient.auth_pb2
import etcdclient.rpc_pb2


class KVBase(abc.ABC):

    @abc.abstractmethod
    async def Range(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.RangeRequest, etcdclient.rpc_pb2.RangeResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Put(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.PutRequest, etcdclient.rpc_pb2.PutResponse]') -> None:
        pass

    @abc.abstractmethod
    async def DeleteRange(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.DeleteRangeRequest, etcdclient.rpc_pb2.DeleteRangeResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Txn(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.TxnRequest, etcdclient.rpc_pb2.TxnResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Compact(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.CompactionRequest, etcdclient.rpc_pb2.CompactionResponse]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/etcdclient.etcdserverpb.KV/Range': grpclib.const.Handler(
                self.Range,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.RangeRequest,
                etcdclient.rpc_pb2.RangeResponse,
            ),
            '/etcdclient.etcdserverpb.KV/Put': grpclib.const.Handler(
                self.Put,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.PutRequest,
                etcdclient.rpc_pb2.PutResponse,
            ),
            '/etcdclient.etcdserverpb.KV/DeleteRange': grpclib.const.Handler(
                self.DeleteRange,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.DeleteRangeRequest,
                etcdclient.rpc_pb2.DeleteRangeResponse,
            ),
            '/etcdclient.etcdserverpb.KV/Txn': grpclib.const.Handler(
                self.Txn,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.TxnRequest,
                etcdclient.rpc_pb2.TxnResponse,
            ),
            '/etcdclient.etcdserverpb.KV/Compact': grpclib.const.Handler(
                self.Compact,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.CompactionRequest,
                etcdclient.rpc_pb2.CompactionResponse,
            ),
        }


class KVStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.Range = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.KV/Range',
            etcdclient.rpc_pb2.RangeRequest,
            etcdclient.rpc_pb2.RangeResponse,
        )
        self.Put = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.KV/Put',
            etcdclient.rpc_pb2.PutRequest,
            etcdclient.rpc_pb2.PutResponse,
        )
        self.DeleteRange = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.KV/DeleteRange',
            etcdclient.rpc_pb2.DeleteRangeRequest,
            etcdclient.rpc_pb2.DeleteRangeResponse,
        )
        self.Txn = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.KV/Txn',
            etcdclient.rpc_pb2.TxnRequest,
            etcdclient.rpc_pb2.TxnResponse,
        )
        self.Compact = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.KV/Compact',
            etcdclient.rpc_pb2.CompactionRequest,
            etcdclient.rpc_pb2.CompactionResponse,
        )


class WatchBase(abc.ABC):

    @abc.abstractmethod
    async def Watch(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.WatchRequest, etcdclient.rpc_pb2.WatchResponse]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/etcdclient.etcdserverpb.Watch/Watch': grpclib.const.Handler(
                self.Watch,
                grpclib.const.Cardinality.STREAM_STREAM,
                etcdclient.rpc_pb2.WatchRequest,
                etcdclient.rpc_pb2.WatchResponse,
            ),
        }


class WatchStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.Watch = grpclib.client.StreamStreamMethod(
            channel,
            '/etcdclient.etcdserverpb.Watch/Watch',
            etcdclient.rpc_pb2.WatchRequest,
            etcdclient.rpc_pb2.WatchResponse,
        )


class LeaseBase(abc.ABC):

    @abc.abstractmethod
    async def LeaseGrant(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.LeaseGrantRequest, etcdclient.rpc_pb2.LeaseGrantResponse]') -> None:
        pass

    @abc.abstractmethod
    async def LeaseRevoke(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.LeaseRevokeRequest, etcdclient.rpc_pb2.LeaseRevokeResponse]') -> None:
        pass

    @abc.abstractmethod
    async def LeaseKeepAlive(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.LeaseKeepAliveRequest, etcdclient.rpc_pb2.LeaseKeepAliveResponse]') -> None:
        pass

    @abc.abstractmethod
    async def LeaseTimeToLive(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.LeaseTimeToLiveRequest, etcdclient.rpc_pb2.LeaseTimeToLiveResponse]') -> None:
        pass

    @abc.abstractmethod
    async def LeaseLeases(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.LeaseLeasesRequest, etcdclient.rpc_pb2.LeaseLeasesResponse]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/etcdclient.etcdserverpb.Lease/LeaseGrant': grpclib.const.Handler(
                self.LeaseGrant,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.LeaseGrantRequest,
                etcdclient.rpc_pb2.LeaseGrantResponse,
            ),
            '/etcdclient.etcdserverpb.Lease/LeaseRevoke': grpclib.const.Handler(
                self.LeaseRevoke,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.LeaseRevokeRequest,
                etcdclient.rpc_pb2.LeaseRevokeResponse,
            ),
            '/etcdclient.etcdserverpb.Lease/LeaseKeepAlive': grpclib.const.Handler(
                self.LeaseKeepAlive,
                grpclib.const.Cardinality.STREAM_STREAM,
                etcdclient.rpc_pb2.LeaseKeepAliveRequest,
                etcdclient.rpc_pb2.LeaseKeepAliveResponse,
            ),
            '/etcdclient.etcdserverpb.Lease/LeaseTimeToLive': grpclib.const.Handler(
                self.LeaseTimeToLive,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.LeaseTimeToLiveRequest,
                etcdclient.rpc_pb2.LeaseTimeToLiveResponse,
            ),
            '/etcdclient.etcdserverpb.Lease/LeaseLeases': grpclib.const.Handler(
                self.LeaseLeases,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.LeaseLeasesRequest,
                etcdclient.rpc_pb2.LeaseLeasesResponse,
            ),
        }


class LeaseStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.LeaseGrant = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Lease/LeaseGrant',
            etcdclient.rpc_pb2.LeaseGrantRequest,
            etcdclient.rpc_pb2.LeaseGrantResponse,
        )
        self.LeaseRevoke = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Lease/LeaseRevoke',
            etcdclient.rpc_pb2.LeaseRevokeRequest,
            etcdclient.rpc_pb2.LeaseRevokeResponse,
        )
        self.LeaseKeepAlive = grpclib.client.StreamStreamMethod(
            channel,
            '/etcdclient.etcdserverpb.Lease/LeaseKeepAlive',
            etcdclient.rpc_pb2.LeaseKeepAliveRequest,
            etcdclient.rpc_pb2.LeaseKeepAliveResponse,
        )
        self.LeaseTimeToLive = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Lease/LeaseTimeToLive',
            etcdclient.rpc_pb2.LeaseTimeToLiveRequest,
            etcdclient.rpc_pb2.LeaseTimeToLiveResponse,
        )
        self.LeaseLeases = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Lease/LeaseLeases',
            etcdclient.rpc_pb2.LeaseLeasesRequest,
            etcdclient.rpc_pb2.LeaseLeasesResponse,
        )


class ClusterBase(abc.ABC):

    @abc.abstractmethod
    async def MemberAdd(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.MemberAddRequest, etcdclient.rpc_pb2.MemberAddResponse]') -> None:
        pass

    @abc.abstractmethod
    async def MemberRemove(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.MemberRemoveRequest, etcdclient.rpc_pb2.MemberRemoveResponse]') -> None:
        pass

    @abc.abstractmethod
    async def MemberUpdate(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.MemberUpdateRequest, etcdclient.rpc_pb2.MemberUpdateResponse]') -> None:
        pass

    @abc.abstractmethod
    async def MemberList(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.MemberListRequest, etcdclient.rpc_pb2.MemberListResponse]') -> None:
        pass

    @abc.abstractmethod
    async def MemberPromote(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.MemberPromoteRequest, etcdclient.rpc_pb2.MemberPromoteResponse]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/etcdclient.etcdserverpb.Cluster/MemberAdd': grpclib.const.Handler(
                self.MemberAdd,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.MemberAddRequest,
                etcdclient.rpc_pb2.MemberAddResponse,
            ),
            '/etcdclient.etcdserverpb.Cluster/MemberRemove': grpclib.const.Handler(
                self.MemberRemove,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.MemberRemoveRequest,
                etcdclient.rpc_pb2.MemberRemoveResponse,
            ),
            '/etcdclient.etcdserverpb.Cluster/MemberUpdate': grpclib.const.Handler(
                self.MemberUpdate,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.MemberUpdateRequest,
                etcdclient.rpc_pb2.MemberUpdateResponse,
            ),
            '/etcdclient.etcdserverpb.Cluster/MemberList': grpclib.const.Handler(
                self.MemberList,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.MemberListRequest,
                etcdclient.rpc_pb2.MemberListResponse,
            ),
            '/etcdclient.etcdserverpb.Cluster/MemberPromote': grpclib.const.Handler(
                self.MemberPromote,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.MemberPromoteRequest,
                etcdclient.rpc_pb2.MemberPromoteResponse,
            ),
        }


class ClusterStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.MemberAdd = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Cluster/MemberAdd',
            etcdclient.rpc_pb2.MemberAddRequest,
            etcdclient.rpc_pb2.MemberAddResponse,
        )
        self.MemberRemove = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Cluster/MemberRemove',
            etcdclient.rpc_pb2.MemberRemoveRequest,
            etcdclient.rpc_pb2.MemberRemoveResponse,
        )
        self.MemberUpdate = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Cluster/MemberUpdate',
            etcdclient.rpc_pb2.MemberUpdateRequest,
            etcdclient.rpc_pb2.MemberUpdateResponse,
        )
        self.MemberList = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Cluster/MemberList',
            etcdclient.rpc_pb2.MemberListRequest,
            etcdclient.rpc_pb2.MemberListResponse,
        )
        self.MemberPromote = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Cluster/MemberPromote',
            etcdclient.rpc_pb2.MemberPromoteRequest,
            etcdclient.rpc_pb2.MemberPromoteResponse,
        )


class MaintenanceBase(abc.ABC):

    @abc.abstractmethod
    async def Alarm(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AlarmRequest, etcdclient.rpc_pb2.AlarmResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Status(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.StatusRequest, etcdclient.rpc_pb2.StatusResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Defragment(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.DefragmentRequest, etcdclient.rpc_pb2.DefragmentResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Hash(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.HashRequest, etcdclient.rpc_pb2.HashResponse]') -> None:
        pass

    @abc.abstractmethod
    async def HashKV(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.HashKVRequest, etcdclient.rpc_pb2.HashKVResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Snapshot(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.SnapshotRequest, etcdclient.rpc_pb2.SnapshotResponse]') -> None:
        pass

    @abc.abstractmethod
    async def MoveLeader(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.MoveLeaderRequest, etcdclient.rpc_pb2.MoveLeaderResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Downgrade(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.DowngradeRequest, etcdclient.rpc_pb2.DowngradeResponse]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/etcdclient.etcdserverpb.Maintenance/Alarm': grpclib.const.Handler(
                self.Alarm,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AlarmRequest,
                etcdclient.rpc_pb2.AlarmResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/Status': grpclib.const.Handler(
                self.Status,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.StatusRequest,
                etcdclient.rpc_pb2.StatusResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/Defragment': grpclib.const.Handler(
                self.Defragment,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.DefragmentRequest,
                etcdclient.rpc_pb2.DefragmentResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/Hash': grpclib.const.Handler(
                self.Hash,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.HashRequest,
                etcdclient.rpc_pb2.HashResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/HashKV': grpclib.const.Handler(
                self.HashKV,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.HashKVRequest,
                etcdclient.rpc_pb2.HashKVResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/Snapshot': grpclib.const.Handler(
                self.Snapshot,
                grpclib.const.Cardinality.UNARY_STREAM,
                etcdclient.rpc_pb2.SnapshotRequest,
                etcdclient.rpc_pb2.SnapshotResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/MoveLeader': grpclib.const.Handler(
                self.MoveLeader,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.MoveLeaderRequest,
                etcdclient.rpc_pb2.MoveLeaderResponse,
            ),
            '/etcdclient.etcdserverpb.Maintenance/Downgrade': grpclib.const.Handler(
                self.Downgrade,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.DowngradeRequest,
                etcdclient.rpc_pb2.DowngradeResponse,
            ),
        }


class MaintenanceStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.Alarm = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/Alarm',
            etcdclient.rpc_pb2.AlarmRequest,
            etcdclient.rpc_pb2.AlarmResponse,
        )
        self.Status = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/Status',
            etcdclient.rpc_pb2.StatusRequest,
            etcdclient.rpc_pb2.StatusResponse,
        )
        self.Defragment = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/Defragment',
            etcdclient.rpc_pb2.DefragmentRequest,
            etcdclient.rpc_pb2.DefragmentResponse,
        )
        self.Hash = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/Hash',
            etcdclient.rpc_pb2.HashRequest,
            etcdclient.rpc_pb2.HashResponse,
        )
        self.HashKV = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/HashKV',
            etcdclient.rpc_pb2.HashKVRequest,
            etcdclient.rpc_pb2.HashKVResponse,
        )
        self.Snapshot = grpclib.client.UnaryStreamMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/Snapshot',
            etcdclient.rpc_pb2.SnapshotRequest,
            etcdclient.rpc_pb2.SnapshotResponse,
        )
        self.MoveLeader = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/MoveLeader',
            etcdclient.rpc_pb2.MoveLeaderRequest,
            etcdclient.rpc_pb2.MoveLeaderResponse,
        )
        self.Downgrade = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Maintenance/Downgrade',
            etcdclient.rpc_pb2.DowngradeRequest,
            etcdclient.rpc_pb2.DowngradeResponse,
        )


class AuthBase(abc.ABC):

    @abc.abstractmethod
    async def AuthEnable(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthEnableRequest, etcdclient.rpc_pb2.AuthEnableResponse]') -> None:
        pass

    @abc.abstractmethod
    async def AuthDisable(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthDisableRequest, etcdclient.rpc_pb2.AuthDisableResponse]') -> None:
        pass

    @abc.abstractmethod
    async def AuthStatus(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthStatusRequest, etcdclient.rpc_pb2.AuthStatusResponse]') -> None:
        pass

    @abc.abstractmethod
    async def Authenticate(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthenticateRequest, etcdclient.rpc_pb2.AuthenticateResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserAdd(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserAddRequest, etcdclient.rpc_pb2.AuthUserAddResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserGet(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserGetRequest, etcdclient.rpc_pb2.AuthUserGetResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserList(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserListRequest, etcdclient.rpc_pb2.AuthUserListResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserDelete(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserDeleteRequest, etcdclient.rpc_pb2.AuthUserDeleteResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserChangePassword(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserChangePasswordRequest, etcdclient.rpc_pb2.AuthUserChangePasswordResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserGrantRole(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserGrantRoleRequest, etcdclient.rpc_pb2.AuthUserGrantRoleResponse]') -> None:
        pass

    @abc.abstractmethod
    async def UserRevokeRole(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthUserRevokeRoleRequest, etcdclient.rpc_pb2.AuthUserRevokeRoleResponse]') -> None:
        pass

    @abc.abstractmethod
    async def RoleAdd(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthRoleAddRequest, etcdclient.rpc_pb2.AuthRoleAddResponse]') -> None:
        pass

    @abc.abstractmethod
    async def RoleGet(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthRoleGetRequest, etcdclient.rpc_pb2.AuthRoleGetResponse]') -> None:
        pass

    @abc.abstractmethod
    async def RoleList(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthRoleListRequest, etcdclient.rpc_pb2.AuthRoleListResponse]') -> None:
        pass

    @abc.abstractmethod
    async def RoleDelete(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthRoleDeleteRequest, etcdclient.rpc_pb2.AuthRoleDeleteResponse]') -> None:
        pass

    @abc.abstractmethod
    async def RoleGrantPermission(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthRoleGrantPermissionRequest, etcdclient.rpc_pb2.AuthRoleGrantPermissionResponse]') -> None:
        pass

    @abc.abstractmethod
    async def RoleRevokePermission(self, stream: 'grpclib.server.Stream[etcdclient.rpc_pb2.AuthRoleRevokePermissionRequest, etcdclient.rpc_pb2.AuthRoleRevokePermissionResponse]') -> None:
        pass

    def __mapping__(self) -> typing.Dict[str, grpclib.const.Handler]:
        return {
            '/etcdclient.etcdserverpb.Auth/AuthEnable': grpclib.const.Handler(
                self.AuthEnable,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthEnableRequest,
                etcdclient.rpc_pb2.AuthEnableResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/AuthDisable': grpclib.const.Handler(
                self.AuthDisable,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthDisableRequest,
                etcdclient.rpc_pb2.AuthDisableResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/AuthStatus': grpclib.const.Handler(
                self.AuthStatus,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthStatusRequest,
                etcdclient.rpc_pb2.AuthStatusResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/Authenticate': grpclib.const.Handler(
                self.Authenticate,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthenticateRequest,
                etcdclient.rpc_pb2.AuthenticateResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserAdd': grpclib.const.Handler(
                self.UserAdd,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserAddRequest,
                etcdclient.rpc_pb2.AuthUserAddResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserGet': grpclib.const.Handler(
                self.UserGet,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserGetRequest,
                etcdclient.rpc_pb2.AuthUserGetResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserList': grpclib.const.Handler(
                self.UserList,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserListRequest,
                etcdclient.rpc_pb2.AuthUserListResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserDelete': grpclib.const.Handler(
                self.UserDelete,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserDeleteRequest,
                etcdclient.rpc_pb2.AuthUserDeleteResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserChangePassword': grpclib.const.Handler(
                self.UserChangePassword,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserChangePasswordRequest,
                etcdclient.rpc_pb2.AuthUserChangePasswordResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserGrantRole': grpclib.const.Handler(
                self.UserGrantRole,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserGrantRoleRequest,
                etcdclient.rpc_pb2.AuthUserGrantRoleResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/UserRevokeRole': grpclib.const.Handler(
                self.UserRevokeRole,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthUserRevokeRoleRequest,
                etcdclient.rpc_pb2.AuthUserRevokeRoleResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/RoleAdd': grpclib.const.Handler(
                self.RoleAdd,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthRoleAddRequest,
                etcdclient.rpc_pb2.AuthRoleAddResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/RoleGet': grpclib.const.Handler(
                self.RoleGet,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthRoleGetRequest,
                etcdclient.rpc_pb2.AuthRoleGetResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/RoleList': grpclib.const.Handler(
                self.RoleList,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthRoleListRequest,
                etcdclient.rpc_pb2.AuthRoleListResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/RoleDelete': grpclib.const.Handler(
                self.RoleDelete,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthRoleDeleteRequest,
                etcdclient.rpc_pb2.AuthRoleDeleteResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/RoleGrantPermission': grpclib.const.Handler(
                self.RoleGrantPermission,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthRoleGrantPermissionRequest,
                etcdclient.rpc_pb2.AuthRoleGrantPermissionResponse,
            ),
            '/etcdclient.etcdserverpb.Auth/RoleRevokePermission': grpclib.const.Handler(
                self.RoleRevokePermission,
                grpclib.const.Cardinality.UNARY_UNARY,
                etcdclient.rpc_pb2.AuthRoleRevokePermissionRequest,
                etcdclient.rpc_pb2.AuthRoleRevokePermissionResponse,
            ),
        }


class AuthStub:

    def __init__(self, channel: grpclib.client.Channel) -> None:
        self.AuthEnable = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/AuthEnable',
            etcdclient.rpc_pb2.AuthEnableRequest,
            etcdclient.rpc_pb2.AuthEnableResponse,
        )
        self.AuthDisable = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/AuthDisable',
            etcdclient.rpc_pb2.AuthDisableRequest,
            etcdclient.rpc_pb2.AuthDisableResponse,
        )
        self.AuthStatus = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/AuthStatus',
            etcdclient.rpc_pb2.AuthStatusRequest,
            etcdclient.rpc_pb2.AuthStatusResponse,
        )
        self.Authenticate = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/Authenticate',
            etcdclient.rpc_pb2.AuthenticateRequest,
            etcdclient.rpc_pb2.AuthenticateResponse,
        )
        self.UserAdd = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserAdd',
            etcdclient.rpc_pb2.AuthUserAddRequest,
            etcdclient.rpc_pb2.AuthUserAddResponse,
        )
        self.UserGet = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserGet',
            etcdclient.rpc_pb2.AuthUserGetRequest,
            etcdclient.rpc_pb2.AuthUserGetResponse,
        )
        self.UserList = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserList',
            etcdclient.rpc_pb2.AuthUserListRequest,
            etcdclient.rpc_pb2.AuthUserListResponse,
        )
        self.UserDelete = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserDelete',
            etcdclient.rpc_pb2.AuthUserDeleteRequest,
            etcdclient.rpc_pb2.AuthUserDeleteResponse,
        )
        self.UserChangePassword = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserChangePassword',
            etcdclient.rpc_pb2.AuthUserChangePasswordRequest,
            etcdclient.rpc_pb2.AuthUserChangePasswordResponse,
        )
        self.UserGrantRole = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserGrantRole',
            etcdclient.rpc_pb2.AuthUserGrantRoleRequest,
            etcdclient.rpc_pb2.AuthUserGrantRoleResponse,
        )
        self.UserRevokeRole = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/UserRevokeRole',
            etcdclient.rpc_pb2.AuthUserRevokeRoleRequest,
            etcdclient.rpc_pb2.AuthUserRevokeRoleResponse,
        )
        self.RoleAdd = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/RoleAdd',
            etcdclient.rpc_pb2.AuthRoleAddRequest,
            etcdclient.rpc_pb2.AuthRoleAddResponse,
        )
        self.RoleGet = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/RoleGet',
            etcdclient.rpc_pb2.AuthRoleGetRequest,
            etcdclient.rpc_pb2.AuthRoleGetResponse,
        )
        self.RoleList = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/RoleList',
            etcdclient.rpc_pb2.AuthRoleListRequest,
            etcdclient.rpc_pb2.AuthRoleListResponse,
        )
        self.RoleDelete = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/RoleDelete',
            etcdclient.rpc_pb2.AuthRoleDeleteRequest,
            etcdclient.rpc_pb2.AuthRoleDeleteResponse,
        )
        self.RoleGrantPermission = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/RoleGrantPermission',
            etcdclient.rpc_pb2.AuthRoleGrantPermissionRequest,
            etcdclient.rpc_pb2.AuthRoleGrantPermissionResponse,
        )
        self.RoleRevokePermission = grpclib.client.UnaryUnaryMethod(
            channel,
            '/etcdclient.etcdserverpb.Auth/RoleRevokePermission',
            etcdclient.rpc_pb2.AuthRoleRevokePermissionRequest,
            etcdclient.rpc_pb2.AuthRoleRevokePermissionResponse,
        )
