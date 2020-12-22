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
  serialized_pb=b'\n\ntube.proto\"\x07\n\x05\x45mpty\"\x14\n\x12ListMethodsRequest\"8\n\x13ListMethodsResponse\x12!\n\x0cmethod_infos\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"\x19\n\x17SubscribeMethodsRequest\"\x1f\n\x0cMethodUpdate\x12\x0f\n\x07methods\x18\x01 \x03(\t\"\r\n\x0bMetaRequest\"?\n\nMetaResult\x12\x1a\n\x07\x65ntries\x18\x01 \x03(\x0b\x32\t.RPCEntry\x12\x15\n\rping_interval\x18\x03 \x01(\x05\"F\n\x08RPCEntry\x12\x1e\n\x08protocol\x18\x01 \x01(\x0e\x32\x0c.RPCProtocol\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x0c\n\x04help\x18\x03 \x01(\t\"<\n\x0eJSONRPCRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0e\n\x06method\x18\x02 \x01(\t\x12\x0e\n\x06params\x18\x03 \x01(\t\"D\n\rJSONRPCResult\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x02ok\x18\n \x01(\tH\x00\x12\x0f\n\x05\x65rror\x18\x0b \x01(\tH\x00\x42\x08\n\x06result\"I\n\x14JSONRPCNotifyRequest\x12\x0e\n\x06method\x18\x01 \x01(\t\x12\x0e\n\x06params\x18\x02 \x01(\t\x12\x11\n\tbroadcast\x18\x03 \x01(\x08\"\x17\n\x15JSONRPCNotifyResponse\"P\n\nMethodInfo\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04help\x18\x02 \x01(\t\x12\x13\n\x0bschema_json\x18\x03 \x01(\t\x12\x11\n\tdelegated\x18\x04 \x01(\x08\"4\n\x14UpdateMethodsRequest\x12\x1c\n\x07methods\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"%\n\x15UpdateMethodsResponse\x12\x0c\n\x04text\x18\x01 \x01(\t\"\x14\n\x04PING\x12\x0c\n\x04text\x18\x01 \x01(\t\"\x14\n\x04PONG\x12\x0c\n\x04text\x18\x01 \x01(\t\"\xc4\x01\n\x11JSONRPCDownPacket\x12\x30\n\x0eupdate_methods\x18\x01 \x01(\x0b\x32\x16.UpdateMethodsResponseH\x00\x12\x15\n\x04ping\x18\n \x01(\x0b\x32\x05.PINGH\x00\x12\x15\n\x04pong\x18\x0b \x01(\x0b\x32\x05.PONGH\x00\x12\"\n\x07request\x18\x14 \x01(\x0b\x32\x0f.JSONRPCRequestH\x00\x12 \n\x06result\x18\x15 \x01(\x0b\x32\x0e.JSONRPCResultH\x00\x42\t\n\x07payload\"\xc1\x01\n\x0fJSONRPCUpPacket\x12/\n\x0eupdate_methods\x18\x01 \x01(\x0b\x32\x15.UpdateMethodsRequestH\x00\x12\x15\n\x04ping\x18\n \x01(\x0b\x32\x05.PINGH\x00\x12\x15\n\x04pong\x18\x0b \x01(\x0b\x32\x05.PONGH\x00\x12\"\n\x07request\x18\x14 \x01(\x0b\x32\x0f.JSONRPCRequestH\x00\x12 \n\x06result\x18\x15 \x01(\x0b\x32\x0e.JSONRPCResultH\x00\x42\t\n\x07payload*\x1a\n\x0bRPCProtocol\x12\x0b\n\x07JSONRPC\x10\x00\x32\xdd\x01\n\x0bJSONRPCTube\x12\x38\n\x0bListMethods\x12\x13.ListMethodsRequest\x1a\x14.ListMethodsResponse\x12\'\n\x04\x43\x61ll\x12\x0f.JSONRPCRequest\x1a\x0e.JSONRPCResult\x12\x37\n\x06Notify\x12\x15.JSONRPCNotifyRequest\x1a\x16.JSONRPCNotifyResponse\x12\x32\n\x06Handle\x12\x10.JSONRPCUpPacket\x1a\x12.JSONRPCDownPacket(\x01\x30\x01\x42\x0bZ\tintf/tubeb\x06proto3'
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
  serialized_start=1161,
  serialized_end=1187,
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


_SUBSCRIBEMETHODSREQUEST = _descriptor.Descriptor(
  name='SubscribeMethodsRequest',
  full_name='SubscribeMethodsRequest',
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
  serialized_end=128,
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
      name='methods', full_name='MethodUpdate.methods', index=0,
      number=1, type=9, cpp_type=9, label=3,
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
  serialized_start=130,
  serialized_end=161,
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
  serialized_start=163,
  serialized_end=176,
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
  serialized_start=178,
  serialized_end=241,
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
  serialized_start=243,
  serialized_end=313,
)


_JSONRPCREQUEST = _descriptor.Descriptor(
  name='JSONRPCRequest',
  full_name='JSONRPCRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='JSONRPCRequest.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='method', full_name='JSONRPCRequest.method', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='params', full_name='JSONRPCRequest.params', index=2,
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
  serialized_start=315,
  serialized_end=375,
)


_JSONRPCRESULT = _descriptor.Descriptor(
  name='JSONRPCResult',
  full_name='JSONRPCResult',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='JSONRPCResult.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ok', full_name='JSONRPCResult.ok', index=1,
      number=10, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='error', full_name='JSONRPCResult.error', index=2,
      number=11, type=9, cpp_type=9, label=1,
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
    _descriptor.OneofDescriptor(
      name='result', full_name='JSONRPCResult.result',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=377,
  serialized_end=445,
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
      name='method', full_name='JSONRPCNotifyRequest.method', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='params', full_name='JSONRPCNotifyRequest.params', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='broadcast', full_name='JSONRPCNotifyRequest.broadcast', index=2,
      number=3, type=8, cpp_type=7, label=1,
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
  serialized_start=447,
  serialized_end=520,
)


_JSONRPCNOTIFYRESPONSE = _descriptor.Descriptor(
  name='JSONRPCNotifyResponse',
  full_name='JSONRPCNotifyResponse',
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
  serialized_start=522,
  serialized_end=545,
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
  serialized_start=547,
  serialized_end=627,
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
  serialized_start=629,
  serialized_end=681,
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
  serialized_start=683,
  serialized_end=720,
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
  serialized_start=722,
  serialized_end=742,
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
  serialized_start=744,
  serialized_end=764,
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
      name='request', full_name='JSONRPCDownPacket.request', index=3,
      number=20, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='result', full_name='JSONRPCDownPacket.result', index=4,
      number=21, type=11, cpp_type=10, label=1,
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
  serialized_start=767,
  serialized_end=963,
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
      name='request', full_name='JSONRPCUpPacket.request', index=3,
      number=20, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='result', full_name='JSONRPCUpPacket.result', index=4,
      number=21, type=11, cpp_type=10, label=1,
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
  serialized_start=966,
  serialized_end=1159,
)

_LISTMETHODSRESPONSE.fields_by_name['method_infos'].message_type = _METHODINFO
_METARESULT.fields_by_name['entries'].message_type = _RPCENTRY
_RPCENTRY.fields_by_name['protocol'].enum_type = _RPCPROTOCOL
_JSONRPCRESULT.oneofs_by_name['result'].fields.append(
  _JSONRPCRESULT.fields_by_name['ok'])
_JSONRPCRESULT.fields_by_name['ok'].containing_oneof = _JSONRPCRESULT.oneofs_by_name['result']
_JSONRPCRESULT.oneofs_by_name['result'].fields.append(
  _JSONRPCRESULT.fields_by_name['error'])
_JSONRPCRESULT.fields_by_name['error'].containing_oneof = _JSONRPCRESULT.oneofs_by_name['result']
_UPDATEMETHODSREQUEST.fields_by_name['methods'].message_type = _METHODINFO
_JSONRPCDOWNPACKET.fields_by_name['update_methods'].message_type = _UPDATEMETHODSRESPONSE
_JSONRPCDOWNPACKET.fields_by_name['ping'].message_type = _PING
_JSONRPCDOWNPACKET.fields_by_name['pong'].message_type = _PONG
_JSONRPCDOWNPACKET.fields_by_name['request'].message_type = _JSONRPCREQUEST
_JSONRPCDOWNPACKET.fields_by_name['result'].message_type = _JSONRPCRESULT
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
  _JSONRPCDOWNPACKET.fields_by_name['request'])
_JSONRPCDOWNPACKET.fields_by_name['request'].containing_oneof = _JSONRPCDOWNPACKET.oneofs_by_name['payload']
_JSONRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCDOWNPACKET.fields_by_name['result'])
_JSONRPCDOWNPACKET.fields_by_name['result'].containing_oneof = _JSONRPCDOWNPACKET.oneofs_by_name['payload']
_JSONRPCUPPACKET.fields_by_name['update_methods'].message_type = _UPDATEMETHODSREQUEST
_JSONRPCUPPACKET.fields_by_name['ping'].message_type = _PING
_JSONRPCUPPACKET.fields_by_name['pong'].message_type = _PONG
_JSONRPCUPPACKET.fields_by_name['request'].message_type = _JSONRPCREQUEST
_JSONRPCUPPACKET.fields_by_name['result'].message_type = _JSONRPCRESULT
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
  _JSONRPCUPPACKET.fields_by_name['request'])
_JSONRPCUPPACKET.fields_by_name['request'].containing_oneof = _JSONRPCUPPACKET.oneofs_by_name['payload']
_JSONRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JSONRPCUPPACKET.fields_by_name['result'])
_JSONRPCUPPACKET.fields_by_name['result'].containing_oneof = _JSONRPCUPPACKET.oneofs_by_name['payload']
DESCRIPTOR.message_types_by_name['Empty'] = _EMPTY
DESCRIPTOR.message_types_by_name['ListMethodsRequest'] = _LISTMETHODSREQUEST
DESCRIPTOR.message_types_by_name['ListMethodsResponse'] = _LISTMETHODSRESPONSE
DESCRIPTOR.message_types_by_name['SubscribeMethodsRequest'] = _SUBSCRIBEMETHODSREQUEST
DESCRIPTOR.message_types_by_name['MethodUpdate'] = _METHODUPDATE
DESCRIPTOR.message_types_by_name['MetaRequest'] = _METAREQUEST
DESCRIPTOR.message_types_by_name['MetaResult'] = _METARESULT
DESCRIPTOR.message_types_by_name['RPCEntry'] = _RPCENTRY
DESCRIPTOR.message_types_by_name['JSONRPCRequest'] = _JSONRPCREQUEST
DESCRIPTOR.message_types_by_name['JSONRPCResult'] = _JSONRPCRESULT
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

SubscribeMethodsRequest = _reflection.GeneratedProtocolMessageType('SubscribeMethodsRequest', (_message.Message,), {
  'DESCRIPTOR' : _SUBSCRIBEMETHODSREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:SubscribeMethodsRequest)
  })
_sym_db.RegisterMessage(SubscribeMethodsRequest)

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

JSONRPCRequest = _reflection.GeneratedProtocolMessageType('JSONRPCRequest', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCREQUEST,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCRequest)
  })
_sym_db.RegisterMessage(JSONRPCRequest)

JSONRPCResult = _reflection.GeneratedProtocolMessageType('JSONRPCResult', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCRESULT,
  '__module__' : 'tube_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCResult)
  })
_sym_db.RegisterMessage(JSONRPCResult)

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
  serialized_start=1190,
  serialized_end=1411,
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
    name='Call',
    full_name='JSONRPCTube.Call',
    index=1,
    containing_service=None,
    input_type=_JSONRPCREQUEST,
    output_type=_JSONRPCRESULT,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Notify',
    full_name='JSONRPCTube.Notify',
    index=2,
    containing_service=None,
    input_type=_JSONRPCNOTIFYREQUEST,
    output_type=_JSONRPCNOTIFYRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Handle',
    full_name='JSONRPCTube.Handle',
    index=3,
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
