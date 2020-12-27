# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: tube.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='tube.proto',
  package='',
  syntax='proto3',
  serialized_options=b'Z\tintf/tube',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\ntube.proto\"\x07\n\x05\x45mpty\"\x14\n\x12ListMethodsRequest\"8\n\x13ListMethodsResponse\x12!\n\x0cmethod_infos\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"\x15\n\x13WatchMethodsRequest\"1\n\x0cMethodUpdate\x12!\n\x0cmethod_infos\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"\r\n\x0bMetaRequest\"?\n\nMetaResult\x12\x1a\n\x07\x65ntries\x18\x01 \x03(\x0b\x32\t.RPCEntry\x12\x15\n\rping_interval\x18\x03 \x01(\x05\"F\n\x08RPCEntry\x12\x1e\n\x08protocol\x18\x01 \x01(\x0e\x32\x0c.RPCProtocol\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x0c\n\x04help\x18\x03 \x01(\t\"\x1f\n\x0fJSONRPCEnvolope\x12\x0c\n\x04\x62ody\x18\x01 \x01(\t\"8\n\x12JSONRPCCallRequest\x12\"\n\x08\x65nvolope\x18\x01 \x01(\x0b\x32\x10.JSONRPCEnvolope\"7\n\x11JSONRPCCallResult\x12\"\n\x08\x65nvolope\x18\x01 \x01(\x0b\x32\x10.JSONRPCEnvolope\"M\n\x14JSONRPCNotifyRequest\x12\"\n\x08\x65nvolope\x18\x01 \x01(\x0b\x32\x10.JSONRPCEnvolope\x12\x11\n\tbroadcast\x18\x02 \x01(\x08\"%\n\x15JSONRPCNotifyResponse\x12\x0c\n\x04text\x18\x01 \x01(\t\"P\n\nMethodInfo\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04help\x18\x02 \x01(\t\x12\x13\n\x0bschema_json\x18\x03 \x01(\t\x12\x11\n\tdelegated\x18\x04 \x01(\x08\"4\n\x14UpdateMethodsRequest\x12\x1c\n\x07methods\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"%\n\x15UpdateMethodsResponse\x12\x0c\n\x04text\x18\x01 \x01(\t\"\x14\n\x04PING\x12\x0c\n\x04text\x18\x01 \x01(\t\"\x14\n\x04PONG\x12\x0c\n\x04text\x18\x01 \x01(\t\"\xa4\x01\n\x11JSONRPCDownPacket\x12\x30\n\x0eupdate_methods\x18\x01 \x01(\x0b\x32\x16.UpdateMethodsResponseH\x00\x12\x15\n\x04ping\x18\n \x01(\x0b\x32\x05.PINGH\x00\x12\x15\n\x04pong\x18\x0b \x01(\x0b\x32\x05.PONGH\x00\x12$\n\x08\x65nvolope\x18\x14 \x01(\x0b\x32\x10.JSONRPCEnvolopeH\x00\x42\t\n\x07payload\"\xa1\x01\n\x0fJSONRPCUpPacket\x12/\n\x0eupdate_methods\x18\x01 \x01(\x0b\x32\x15.UpdateMethodsRequestH\x00\x12\x15\n\x04ping\x18\n \x01(\x0b\x32\x05.PINGH\x00\x12\x15\n\x04pong\x18\x0b \x01(\x0b\x32\x05.PONGH\x00\x12$\n\x08\x65nvolope\x18\x14 \x01(\x0b\x32\x10.JSONRPCEnvolopeH\x00\x42\t\n\x07payload*\x1a\n\x0bRPCProtocol\x12\x0b\n\x07JSONRPC\x10\x00\x32\x9c\x02\n\x0bJSONRPCTube\x12\x38\n\x0bListMethods\x12\x13.ListMethodsRequest\x1a\x14.ListMethodsResponse\x12\x35\n\x0cWatchMethods\x12\x14.WatchMethodsRequest\x1a\r.MethodUpdate0\x01\x12/\n\x04\x43\x61ll\x12\x13.JSONRPCCallRequest\x1a\x12.JSONRPCCallResult\x12\x37\n\x06Notify\x12\x15.JSONRPCNotifyRequest\x1a\x16.JSONRPCNotifyResponse\x12\x32\n\x06Handle\x12\x10.JSONRPCUpPacket\x1a\x12.JSONRPCDownPacket(\x01\x30\x01\x42\x0bZ\tintf/tubeb\x06proto3'
)

_RPCPROTOCOL = _descriptor.EnumDescriptor(
  name='RPCProtocol',
  full_name='RPCProtocol',
  filename=None,
  file=DESCRIPTOR,
  create_key=_descriptor._internal_create_key,
  values=[
    _descriptor.EnumValueDescriptor(
      name='JSONRPC', index=0, number=0,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1145,
  serialized_end=1171,
)
_sym_db.RegisterEnumDescriptor(_RPCPROTOCOL)

RPCProtocol = enum_type_wrapper.EnumTypeWrapper(_RPCPROTOCOL)
JSONRPC = 0



_EMPTY = _descriptor.Descriptor(
  name='Empty',
  full_name='Empty',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=14,
  serialized_end=21,
)


_LISTMETHODSREQUEST = _descriptor.Descriptor(
  name='ListMethodsRequest',
  full_name='ListMethodsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=23,
  serialized_end=43,
)


_LISTMETHODSRESPONSE = _descriptor.Descriptor(
  name='ListMethodsResponse',
  full_name='ListMethodsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='method_infos', full_name='ListMethodsResponse.method_infos', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=45,
  serialized_end=101,
)


_WATCHMETHODSREQUEST = _descriptor.Descriptor(
  name='WatchMethodsRequest',
  full_name='WatchMethodsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=103,
  serialized_end=124,
)


_METHODUPDATE = _descriptor.Descriptor(
  name='MethodUpdate',
  full_name='MethodUpdate',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='method_infos', full_name='MethodUpdate.method_infos', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=126,
  serialized_end=175,
)


_METAREQUEST = _descriptor.Descriptor(
  name='MetaRequest',
  full_name='MetaRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=177,
  serialized_end=190,
)


_METARESULT = _descriptor.Descriptor(
  name='MetaResult',
  full_name='MetaResult',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='entries', full_name='MetaResult.entries', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ping_interval', full_name='MetaResult.ping_interval', index=1,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=192,
  serialized_end=255,
)


_RPCENTRY = _descriptor.Descriptor(
  name='RPCEntry',
  full_name='RPCEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='protocol', full_name='RPCEntry.protocol', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='name', full_name='RPCEntry.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='help', full_name='RPCEntry.help', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=257,
  serialized_end=327,
)


_JSONRPCENVOLOPE = _descriptor.Descriptor(
  name='JSONRPCEnvolope',
  full_name='JSONRPCEnvolope',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='body', full_name='JSONRPCEnvolope.body', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=329,
  serialized_end=360,
)


_JSONRPCCALLREQUEST = _descriptor.Descriptor(
  name='JSONRPCCallRequest',
  full_name='JSONRPCCallRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JSONRPCCallRequest.envolope', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=362,
  serialized_end=418,
)


_JSONRPCCALLRESULT = _descriptor.Descriptor(
  name='JSONRPCCallResult',
  full_name='JSONRPCCallResult',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JSONRPCCallResult.envolope', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=420,
  serialized_end=475,
)


_JSONRPCNOTIFYREQUEST = _descriptor.Descriptor(
  name='JSONRPCNotifyRequest',
  full_name='JSONRPCNotifyRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JSONRPCNotifyRequest.envolope', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='broadcast', full_name='JSONRPCNotifyRequest.broadcast', index=1,
      number=2, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=477,
  serialized_end=554,
)


_JSONRPCNOTIFYRESPONSE = _descriptor.Descriptor(
  name='JSONRPCNotifyResponse',
  full_name='JSONRPCNotifyResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='JSONRPCNotifyResponse.text', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=556,
  serialized_end=593,
)


_METHODINFO = _descriptor.Descriptor(
  name='MethodInfo',
  full_name='MethodInfo',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='MethodInfo.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='help', full_name='MethodInfo.help', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='schema_json', full_name='MethodInfo.schema_json', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='delegated', full_name='MethodInfo.delegated', index=3,
      number=4, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=595,
  serialized_end=675,
)


_UPDATEMETHODSREQUEST = _descriptor.Descriptor(
  name='UpdateMethodsRequest',
  full_name='UpdateMethodsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='methods', full_name='UpdateMethodsRequest.methods', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=677,
  serialized_end=729,
)


_UPDATEMETHODSRESPONSE = _descriptor.Descriptor(
  name='UpdateMethodsResponse',
  full_name='UpdateMethodsResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='UpdateMethodsResponse.text', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=731,
  serialized_end=768,
)


_PING = _descriptor.Descriptor(
  name='PING',
  full_name='PING',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='PING.text', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=770,
  serialized_end=790,
)


_PONG = _descriptor.Descriptor(
  name='PONG',
  full_name='PONG',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='PONG.text', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=792,
  serialized_end=812,
)


_JSONRPCDOWNPACKET = _descriptor.Descriptor(
  name='JSONRPCDownPacket',
  full_name='JSONRPCDownPacket',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='update_methods', full_name='JSONRPCDownPacket.update_methods', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ping', full_name='JSONRPCDownPacket.ping', index=1,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='pong', full_name='JSONRPCDownPacket.pong', index=2,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JSONRPCDownPacket.envolope', index=3,
      number=20, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='payload', full_name='JSONRPCDownPacket.payload',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=815,
  serialized_end=979,
)


_JSONRPCUPPACKET = _descriptor.Descriptor(
  name='JSONRPCUpPacket',
  full_name='JSONRPCUpPacket',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='update_methods', full_name='JSONRPCUpPacket.update_methods', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ping', full_name='JSONRPCUpPacket.ping', index=1,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='pong', full_name='JSONRPCUpPacket.pong', index=2,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JSONRPCUpPacket.envolope', index=3,
      number=20, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
    _descriptor.OneofDescriptor(
      name='payload', full_name='JSONRPCUpPacket.payload',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=982,
  serialized_end=1143,
)

_LISTMETHODSRESPONSE.fields_by_name['method_infos'].message_type = _METHODINFO
_METHODUPDATE.fields_by_name['method_infos'].message_type = _METHODINFO
_METARESULT.fields_by_name['entries'].message_type = _RPCENTRY
_RPCENTRY.fields_by_name['protocol'].enum_type = _RPCPROTOCOL
_JSONRPCCALLREQUEST.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JSONRPCCALLRESULT.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JSONRPCNOTIFYREQUEST.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_UPDATEMETHODSREQUEST.fields_by_name['methods'].message_type = _METHODINFO
_JSONRPCDOWNPACKET.fields_by_name['update_methods'].message_type = _UPDATEMETHODSRESPONSE
_JSONRPCDOWNPACKET.fields_by_name['ping'].message_type = _PING
_JSONRPCDOWNPACKET.fields_by_name['pong'].message_type = _PONG
_JSONRPCDOWNPACKET.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JSONRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCDOWNPACKET.fields_by_name['update_methods'])
_JSONRPCDOWNPACKET.fields_by_name['update_methods'].containing_oneof = _JSONRPCDOWNPACKET.oneofs_by_name['payload']
_JSONRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCDOWNPACKET.fields_by_name['ping'])
_JSONRPCDOWNPACKET.fields_by_name['ping'].containing_oneof = _JSONRPCDOWNPACKET.oneofs_by_name['payload']
_JSONRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCDOWNPACKET.fields_by_name['pong'])
_JSONRPCDOWNPACKET.fields_by_name['pong'].containing_oneof = _JSONRPCDOWNPACKET.oneofs_by_name['payload']
_JSONRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCDOWNPACKET.fields_by_name['envolope'])
_JSONRPCDOWNPACKET.fields_by_name['envolope'].containing_oneof = _JSONRPCDOWNPACKET.oneofs_by_name['payload']
_JSONRPCUPPACKET.fields_by_name['update_methods'].message_type = _UPDATEMETHODSREQUEST
_JSONRPCUPPACKET.fields_by_name['ping'].message_type = _PING
_JSONRPCUPPACKET.fields_by_name['pong'].message_type = _PONG
_JSONRPCUPPACKET.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JSONRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCUPPACKET.fields_by_name['update_methods'])
_JSONRPCUPPACKET.fields_by_name['update_methods'].containing_oneof = _JSONRPCUPPACKET.oneofs_by_name['payload']
_JSONRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCUPPACKET.fields_by_name['ping'])
_JSONRPCUPPACKET.fields_by_name['ping'].containing_oneof = _JSONRPCUPPACKET.oneofs_by_name['payload']
_JSONRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCUPPACKET.fields_by_name['pong'])
_JSONRPCUPPACKET.fields_by_name['pong'].containing_oneof = _JSONRPCUPPACKET.oneofs_by_name['payload']
_JSONRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCUPPACKET.fields_by_name['envolope'])
_JSONRPCUPPACKET.fields_by_name['envolope'].containing_oneof = _JSONRPCUPPACKET.oneofs_by_name['payload']
DESCRIPTOR.message_types_by_name['Empty'] = _EMPTY
DESCRIPTOR.message_types_by_name['ListMethodsRequest'] = _LISTMETHODSREQUEST
DESCRIPTOR.message_types_by_name['ListMethodsResponse'] = _LISTMETHODSRESPONSE
DESCRIPTOR.message_types_by_name['WatchMethodsRequest'] = _WATCHMETHODSREQUEST
DESCRIPTOR.message_types_by_name['MethodUpdate'] = _METHODUPDATE
DESCRIPTOR.message_types_by_name['MetaRequest'] = _METAREQUEST
DESCRIPTOR.message_types_by_name['MetaResult'] = _METARESULT
DESCRIPTOR.message_types_by_name['RPCEntry'] = _RPCENTRY
DESCRIPTOR.message_types_by_name['JSONRPCEnvolope'] = _JSONRPCENVOLOPE
DESCRIPTOR.message_types_by_name['JSONRPCCallRequest'] = _JSONRPCCALLREQUEST
DESCRIPTOR.message_types_by_name['JSONRPCCallResult'] = _JSONRPCCALLRESULT
DESCRIPTOR.message_types_by_name['JSONRPCNotifyRequest'] = _JSONRPCNOTIFYREQUEST
DESCRIPTOR.message_types_by_name['JSONRPCNotifyResponse'] = _JSONRPCNOTIFYRESPONSE
DESCRIPTOR.message_types_by_name['MethodInfo'] = _METHODINFO
DESCRIPTOR.message_types_by_name['UpdateMethodsRequest'] = _UPDATEMETHODSREQUEST
DESCRIPTOR.message_types_by_name['UpdateMethodsResponse'] = _UPDATEMETHODSRESPONSE
DESCRIPTOR.message_types_by_name['PING'] = _PING
DESCRIPTOR.message_types_by_name['PONG'] = _PONG
DESCRIPTOR.message_types_by_name['JSONRPCDownPacket'] = _JSONRPCDOWNPACKET
DESCRIPTOR.message_types_by_name['JSONRPCUpPacket'] = _JSONRPCUPPACKET
DESCRIPTOR.enum_types_by_name['RPCProtocol'] = _RPCPROTOCOL
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Empty = _reflection.GeneratedProtocolMessageType('Empty', (_message.Message,), {
  'DESCRIPTOR' : _EMPTY,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:Empty)
  })
_sym_db.RegisterMessage(Empty)

ListMethodsRequest = _reflection.GeneratedProtocolMessageType('ListMethodsRequest', (_message.Message,), {
  'DESCRIPTOR' : _LISTMETHODSREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:ListMethodsRequest)
  })
_sym_db.RegisterMessage(ListMethodsRequest)

ListMethodsResponse = _reflection.GeneratedProtocolMessageType('ListMethodsResponse', (_message.Message,), {
  'DESCRIPTOR' : _LISTMETHODSRESPONSE,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:ListMethodsResponse)
  })
_sym_db.RegisterMessage(ListMethodsResponse)

WatchMethodsRequest = _reflection.GeneratedProtocolMessageType('WatchMethodsRequest', (_message.Message,), {
  'DESCRIPTOR' : _WATCHMETHODSREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:WatchMethodsRequest)
  })
_sym_db.RegisterMessage(WatchMethodsRequest)

MethodUpdate = _reflection.GeneratedProtocolMessageType('MethodUpdate', (_message.Message,), {
  'DESCRIPTOR' : _METHODUPDATE,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:MethodUpdate)
  })
_sym_db.RegisterMessage(MethodUpdate)

MetaRequest = _reflection.GeneratedProtocolMessageType('MetaRequest', (_message.Message,), {
  'DESCRIPTOR' : _METAREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:MetaRequest)
  })
_sym_db.RegisterMessage(MetaRequest)

MetaResult = _reflection.GeneratedProtocolMessageType('MetaResult', (_message.Message,), {
  'DESCRIPTOR' : _METARESULT,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:MetaResult)
  })
_sym_db.RegisterMessage(MetaResult)

RPCEntry = _reflection.GeneratedProtocolMessageType('RPCEntry', (_message.Message,), {
  'DESCRIPTOR' : _RPCENTRY,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:RPCEntry)
  })
_sym_db.RegisterMessage(RPCEntry)

JSONRPCEnvolope = _reflection.GeneratedProtocolMessageType('JSONRPCEnvolope', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCENVOLOPE,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCEnvolope)
  })
_sym_db.RegisterMessage(JSONRPCEnvolope)

JSONRPCCallRequest = _reflection.GeneratedProtocolMessageType('JSONRPCCallRequest', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCCALLREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCCallRequest)
  })
_sym_db.RegisterMessage(JSONRPCCallRequest)

JSONRPCCallResult = _reflection.GeneratedProtocolMessageType('JSONRPCCallResult', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCCALLRESULT,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCCallResult)
  })
_sym_db.RegisterMessage(JSONRPCCallResult)

JSONRPCNotifyRequest = _reflection.GeneratedProtocolMessageType('JSONRPCNotifyRequest', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCNOTIFYREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCNotifyRequest)
  })
_sym_db.RegisterMessage(JSONRPCNotifyRequest)

JSONRPCNotifyResponse = _reflection.GeneratedProtocolMessageType('JSONRPCNotifyResponse', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCNOTIFYRESPONSE,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCNotifyResponse)
  })
_sym_db.RegisterMessage(JSONRPCNotifyResponse)

MethodInfo = _reflection.GeneratedProtocolMessageType('MethodInfo', (_message.Message,), {
  'DESCRIPTOR' : _METHODINFO,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:MethodInfo)
  })
_sym_db.RegisterMessage(MethodInfo)

UpdateMethodsRequest = _reflection.GeneratedProtocolMessageType('UpdateMethodsRequest', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEMETHODSREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:UpdateMethodsRequest)
  })
_sym_db.RegisterMessage(UpdateMethodsRequest)

UpdateMethodsResponse = _reflection.GeneratedProtocolMessageType('UpdateMethodsResponse', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEMETHODSRESPONSE,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:UpdateMethodsResponse)
  })
_sym_db.RegisterMessage(UpdateMethodsResponse)

PING = _reflection.GeneratedProtocolMessageType('PING', (_message.Message,), {
  'DESCRIPTOR' : _PING,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:PING)
  })
_sym_db.RegisterMessage(PING)

PONG = _reflection.GeneratedProtocolMessageType('PONG', (_message.Message,), {
  'DESCRIPTOR' : _PONG,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:PONG)
  })
_sym_db.RegisterMessage(PONG)

JSONRPCDownPacket = _reflection.GeneratedProtocolMessageType('JSONRPCDownPacket', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCDOWNPACKET,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCDownPacket)
  })
_sym_db.RegisterMessage(JSONRPCDownPacket)

JSONRPCUpPacket = _reflection.GeneratedProtocolMessageType('JSONRPCUpPacket', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCUPPACKET,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCUpPacket)
  })
_sym_db.RegisterMessage(JSONRPCUpPacket)


DESCRIPTOR._options = None

_JSONRPCTUBE = _descriptor.ServiceDescriptor(
  name='JSONRPCTube',
  full_name='JSONRPCTube',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=1174,
  serialized_end=1458,
  methods=[
  _descriptor.MethodDescriptor(
    name='ListMethods',
    full_name='JSONRPCTube.ListMethods',
    index=0,
    containing_service=None,
    input_type=_LISTMETHODSREQUEST,
    output_type=_LISTMETHODSRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='WatchMethods',
    full_name='JSONRPCTube.WatchMethods',
    index=1,
    containing_service=None,
    input_type=_WATCHMETHODSREQUEST,
    output_type=_METHODUPDATE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Call',
    full_name='JSONRPCTube.Call',
    index=2,
    containing_service=None,
    input_type=_JSONRPCCALLREQUEST,
    output_type=_JSONRPCCALLRESULT,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Notify',
    full_name='JSONRPCTube.Notify',
    index=3,
    containing_service=None,
    input_type=_JSONRPCNOTIFYREQUEST,
    output_type=_JSONRPCNOTIFYRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Handle',
    full_name='JSONRPCTube.Handle',
    index=4,
    containing_service=None,
    input_type=_JSONRPCUPPACKET,
    output_type=_JSONRPCDOWNPACKET,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_JSONRPCTUBE)

DESCRIPTOR.services_by_name['JSONRPCTube'] = _JSONRPCTUBE

# @@protoc_insertion_point(module_scope)
