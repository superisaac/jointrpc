# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: jointrpc.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='jointrpc.proto',
  package='',
  syntax='proto3',
  serialized_options=b'Z\rintf/jointrpc',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0ejointrpc.proto\"\x07\n\x05\x45mpty\"\x14\n\x12ListMethodsRequest\"8\n\x13ListMethodsResponse\x12!\n\x0cmethod_infos\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"\x16\n\x14ListDelegatesRequest\"*\n\x15ListDelegatesResponse\x12\x11\n\tdelegates\x18\x01 \x03(\t\"\x1f\n\x0fJSONRPCEnvolope\x12\x0c\n\x04\x62ody\x18\x01 \x01(\t\"K\n\x12JSONRPCCallRequest\x12\"\n\x08\x65nvolope\x18\x01 \x01(\x0b\x32\x10.JSONRPCEnvolope\x12\x11\n\tbroadcast\x18\x02 \x01(\x08\"7\n\x11JSONRPCCallResult\x12\"\n\x08\x65nvolope\x18\x01 \x01(\x0b\x32\x10.JSONRPCEnvolope\"M\n\x14JSONRPCNotifyRequest\x12\"\n\x08\x65nvolope\x18\x01 \x01(\x0b\x32\x10.JSONRPCEnvolope\x12\x11\n\tbroadcast\x18\x02 \x01(\x08\"%\n\x15JSONRPCNotifyResponse\x12\x0c\n\x04text\x18\x01 \x01(\t\"=\n\nMethodInfo\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04help\x18\x02 \x01(\t\x12\x13\n\x0bschema_json\x18\x03 \x01(\t\"/\n\x0f\x43\x61nServeRequest\x12\x1c\n\x07methods\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\" \n\x10\x43\x61nServeResponse\x12\x0c\n\x04text\x18\x01 \x01(\t\"%\n\x12\x43\x61nDelegateRequest\x12\x0f\n\x07methods\x18\x01 \x03(\t\"#\n\x13\x43\x61nDelegateResponse\x12\x0c\n\x04text\x18\x01 \x01(\t\"\x14\n\x04PING\x12\x0c\n\x04text\x18\x01 \x01(\t\"\x14\n\x04PONG\x12\x0c\n\x04text\x18\x01 \x01(\t\")\n\tTubeState\x12\x1c\n\x07methods\x18\x01 \x03(\x0b\x32\x0b.MethodInfo\"\xe6\x01\n\x12JointRPCDownPacket\x12&\n\tcan_serve\x18\x01 \x01(\x0b\x32\x11.CanServeResponseH\x00\x12,\n\x0c\x63\x61n_delegate\x18\x02 \x01(\x0b\x32\x14.CanDelegateResponseH\x00\x12\x15\n\x04ping\x18\n \x01(\x0b\x32\x05.PINGH\x00\x12\x15\n\x04pong\x18\x0b \x01(\x0b\x32\x05.PONGH\x00\x12\x1b\n\x05state\x18\x0f \x01(\x0b\x32\n.TubeStateH\x00\x12$\n\x08\x65nvolope\x18\x14 \x01(\x0b\x32\x10.JSONRPCEnvolopeH\x00\x42\t\n\x07payload\"\xc5\x01\n\x10JointRPCUpPacket\x12%\n\tcan_serve\x18\x01 \x01(\x0b\x32\x10.CanServeRequestH\x00\x12+\n\x0c\x63\x61n_delegate\x18\x02 \x01(\x0b\x32\x13.CanDelegateRequestH\x00\x12\x15\n\x04ping\x18\n \x01(\x0b\x32\x05.PINGH\x00\x12\x15\n\x04pong\x18\x0b \x01(\x0b\x32\x05.PONGH\x00\x12$\n\x08\x65nvolope\x18\x14 \x01(\x0b\x32\x10.JSONRPCEnvolopeH\x00\x42\t\n\x07payload2\xa4\x02\n\x08JointRPC\x12\x38\n\x0bListMethods\x12\x13.ListMethodsRequest\x1a\x14.ListMethodsResponse\x12>\n\rListDelegates\x12\x15.ListDelegatesRequest\x1a\x16.ListDelegatesResponse\x12/\n\x04\x43\x61ll\x12\x13.JSONRPCCallRequest\x1a\x12.JSONRPCCallResult\x12\x37\n\x06Notify\x12\x15.JSONRPCNotifyRequest\x1a\x16.JSONRPCNotifyResponse\x12\x34\n\x06Handle\x12\x11.JointRPCUpPacket\x1a\x13.JointRPCDownPacket(\x01\x30\x01\x42\x0fZ\rintf/jointrpcb\x06proto3'
)




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
  serialized_start=18,
  serialized_end=25,
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
  serialized_start=27,
  serialized_end=47,
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
  serialized_start=49,
  serialized_end=105,
)


_LISTDELEGATESREQUEST = _descriptor.Descriptor(
  name='ListDelegatesRequest',
  full_name='ListDelegatesRequest',
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
  serialized_start=107,
  serialized_end=129,
)


_LISTDELEGATESRESPONSE = _descriptor.Descriptor(
  name='ListDelegatesResponse',
  full_name='ListDelegatesResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='delegates', full_name='ListDelegatesResponse.delegates', index=0,
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
  serialized_start=131,
  serialized_end=173,
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
  serialized_start=175,
  serialized_end=206,
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
    _descriptor.FieldDescriptor(
      name='broadcast', full_name='JSONRPCCallRequest.broadcast', index=1,
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
  serialized_start=208,
  serialized_end=283,
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
  serialized_start=285,
  serialized_end=340,
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
  serialized_start=342,
  serialized_end=419,
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
  serialized_start=421,
  serialized_end=458,
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
  serialized_start=460,
  serialized_end=521,
)


_CANSERVEREQUEST = _descriptor.Descriptor(
  name='CanServeRequest',
  full_name='CanServeRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='methods', full_name='CanServeRequest.methods', index=0,
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
  serialized_start=523,
  serialized_end=570,
)


_CANSERVERESPONSE = _descriptor.Descriptor(
  name='CanServeResponse',
  full_name='CanServeResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='CanServeResponse.text', index=0,
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
  serialized_start=572,
  serialized_end=604,
)


_CANDELEGATEREQUEST = _descriptor.Descriptor(
  name='CanDelegateRequest',
  full_name='CanDelegateRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='methods', full_name='CanDelegateRequest.methods', index=0,
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
  serialized_start=606,
  serialized_end=643,
)


_CANDELEGATERESPONSE = _descriptor.Descriptor(
  name='CanDelegateResponse',
  full_name='CanDelegateResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='CanDelegateResponse.text', index=0,
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
  serialized_start=645,
  serialized_end=680,
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
  serialized_start=682,
  serialized_end=702,
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
  serialized_start=704,
  serialized_end=724,
)


_TUBESTATE = _descriptor.Descriptor(
  name='TubeState',
  full_name='TubeState',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='methods', full_name='TubeState.methods', index=0,
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
  serialized_start=726,
  serialized_end=767,
)


_JOINTRPCDOWNPACKET = _descriptor.Descriptor(
  name='JointRPCDownPacket',
  full_name='JointRPCDownPacket',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='can_serve', full_name='JointRPCDownPacket.can_serve', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='can_delegate', full_name='JointRPCDownPacket.can_delegate', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ping', full_name='JointRPCDownPacket.ping', index=2,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='pong', full_name='JointRPCDownPacket.pong', index=3,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='state', full_name='JointRPCDownPacket.state', index=4,
      number=15, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JointRPCDownPacket.envolope', index=5,
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
      name='payload', full_name='JointRPCDownPacket.payload',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=770,
  serialized_end=1000,
)


_JOINTRPCUPPACKET = _descriptor.Descriptor(
  name='JointRPCUpPacket',
  full_name='JointRPCUpPacket',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='can_serve', full_name='JointRPCUpPacket.can_serve', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='can_delegate', full_name='JointRPCUpPacket.can_delegate', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ping', full_name='JointRPCUpPacket.ping', index=2,
      number=10, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='pong', full_name='JointRPCUpPacket.pong', index=3,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='envolope', full_name='JointRPCUpPacket.envolope', index=4,
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
      name='payload', full_name='JointRPCUpPacket.payload',
      index=0, containing_type=None,
      create_key=_descriptor._internal_create_key,
    fields=[]),
  ],
  serialized_start=1003,
  serialized_end=1200,
)

_LISTMETHODSRESPONSE.fields_by_name['method_infos'].message_type = _METHODINFO
_JSONRPCCALLREQUEST.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JSONRPCCALLRESULT.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JSONRPCNOTIFYREQUEST.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_CANSERVEREQUEST.fields_by_name['methods'].message_type = _METHODINFO
_TUBESTATE.fields_by_name['methods'].message_type = _METHODINFO
_JOINTRPCDOWNPACKET.fields_by_name['can_serve'].message_type = _CANSERVERESPONSE
_JOINTRPCDOWNPACKET.fields_by_name['can_delegate'].message_type = _CANDELEGATERESPONSE
_JOINTRPCDOWNPACKET.fields_by_name['ping'].message_type = _PING
_JOINTRPCDOWNPACKET.fields_by_name['pong'].message_type = _PONG
_JOINTRPCDOWNPACKET.fields_by_name['state'].message_type = _TUBESTATE
_JOINTRPCDOWNPACKET.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JOINTRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCDOWNPACKET.fields_by_name['can_serve'])
_JOINTRPCDOWNPACKET.fields_by_name['can_serve'].containing_oneof = _JOINTRPCDOWNPACKET.oneofs_by_name['payload']
_JOINTRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCDOWNPACKET.fields_by_name['can_delegate'])
_JOINTRPCDOWNPACKET.fields_by_name['can_delegate'].containing_oneof = _JOINTRPCDOWNPACKET.oneofs_by_name['payload']
_JOINTRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCDOWNPACKET.fields_by_name['ping'])
_JOINTRPCDOWNPACKET.fields_by_name['ping'].containing_oneof = _JOINTRPCDOWNPACKET.oneofs_by_name['payload']
_JOINTRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCDOWNPACKET.fields_by_name['pong'])
_JOINTRPCDOWNPACKET.fields_by_name['pong'].containing_oneof = _JOINTRPCDOWNPACKET.oneofs_by_name['payload']
_JOINTRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCDOWNPACKET.fields_by_name['state'])
_JOINTRPCDOWNPACKET.fields_by_name['state'].containing_oneof = _JOINTRPCDOWNPACKET.oneofs_by_name['payload']
_JOINTRPCDOWNPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCDOWNPACKET.fields_by_name['envolope'])
_JOINTRPCDOWNPACKET.fields_by_name['envolope'].containing_oneof = _JOINTRPCDOWNPACKET.oneofs_by_name['payload']
_JOINTRPCUPPACKET.fields_by_name['can_serve'].message_type = _CANSERVEREQUEST
_JOINTRPCUPPACKET.fields_by_name['can_delegate'].message_type = _CANDELEGATEREQUEST
_JOINTRPCUPPACKET.fields_by_name['ping'].message_type = _PING
_JOINTRPCUPPACKET.fields_by_name['pong'].message_type = _PONG
_JOINTRPCUPPACKET.fields_by_name['envolope'].message_type = _JSONRPCENVOLOPE
_JOINTRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCUPPACKET.fields_by_name['can_serve'])
_JOINTRPCUPPACKET.fields_by_name['can_serve'].containing_oneof = _JOINTRPCUPPACKET.oneofs_by_name['payload']
_JOINTRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCUPPACKET.fields_by_name['can_delegate'])
_JOINTRPCUPPACKET.fields_by_name['can_delegate'].containing_oneof = _JOINTRPCUPPACKET.oneofs_by_name['payload']
_JOINTRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCUPPACKET.fields_by_name['ping'])
_JOINTRPCUPPACKET.fields_by_name['ping'].containing_oneof = _JOINTRPCUPPACKET.oneofs_by_name['payload']
_JOINTRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCUPPACKET.fields_by_name['pong'])
_JOINTRPCUPPACKET.fields_by_name['pong'].containing_oneof = _JOINTRPCUPPACKET.oneofs_by_name['payload']
_JOINTRPCUPPACKET.oneofs_by_name['payload'].fields.append(
  _JOINTRPCUPPACKET.fields_by_name['envolope'])
_JOINTRPCUPPACKET.fields_by_name['envolope'].containing_oneof = _JOINTRPCUPPACKET.oneofs_by_name['payload']
DESCRIPTOR.message_types_by_name['Empty'] = _EMPTY
DESCRIPTOR.message_types_by_name['ListMethodsRequest'] = _LISTMETHODSREQUEST
DESCRIPTOR.message_types_by_name['ListMethodsResponse'] = _LISTMETHODSRESPONSE
DESCRIPTOR.message_types_by_name['ListDelegatesRequest'] = _LISTDELEGATESREQUEST
DESCRIPTOR.message_types_by_name['ListDelegatesResponse'] = _LISTDELEGATESRESPONSE
DESCRIPTOR.message_types_by_name['JSONRPCEnvolope'] = _JSONRPCENVOLOPE
DESCRIPTOR.message_types_by_name['JSONRPCCallRequest'] = _JSONRPCCALLREQUEST
DESCRIPTOR.message_types_by_name['JSONRPCCallResult'] = _JSONRPCCALLRESULT
DESCRIPTOR.message_types_by_name['JSONRPCNotifyRequest'] = _JSONRPCNOTIFYREQUEST
DESCRIPTOR.message_types_by_name['JSONRPCNotifyResponse'] = _JSONRPCNOTIFYRESPONSE
DESCRIPTOR.message_types_by_name['MethodInfo'] = _METHODINFO
DESCRIPTOR.message_types_by_name['CanServeRequest'] = _CANSERVEREQUEST
DESCRIPTOR.message_types_by_name['CanServeResponse'] = _CANSERVERESPONSE
DESCRIPTOR.message_types_by_name['CanDelegateRequest'] = _CANDELEGATEREQUEST
DESCRIPTOR.message_types_by_name['CanDelegateResponse'] = _CANDELEGATERESPONSE
DESCRIPTOR.message_types_by_name['PING'] = _PING
DESCRIPTOR.message_types_by_name['PONG'] = _PONG
DESCRIPTOR.message_types_by_name['TubeState'] = _TUBESTATE
DESCRIPTOR.message_types_by_name['JointRPCDownPacket'] = _JOINTRPCDOWNPACKET
DESCRIPTOR.message_types_by_name['JointRPCUpPacket'] = _JOINTRPCUPPACKET
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Empty = _reflection.GeneratedProtocolMessageType('Empty', (_message.Message,), {
  'DESCRIPTOR' : _EMPTY,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:Empty)
  })
_sym_db.RegisterMessage(Empty)

ListMethodsRequest = _reflection.GeneratedProtocolMessageType('ListMethodsRequest', (_message.Message,), {
  'DESCRIPTOR' : _LISTMETHODSREQUEST,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:ListMethodsRequest)
  })
_sym_db.RegisterMessage(ListMethodsRequest)

ListMethodsResponse = _reflection.GeneratedProtocolMessageType('ListMethodsResponse', (_message.Message,), {
  'DESCRIPTOR' : _LISTMETHODSRESPONSE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:ListMethodsResponse)
  })
_sym_db.RegisterMessage(ListMethodsResponse)

ListDelegatesRequest = _reflection.GeneratedProtocolMessageType('ListDelegatesRequest', (_message.Message,), {
  'DESCRIPTOR' : _LISTDELEGATESREQUEST,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:ListDelegatesRequest)
  })
_sym_db.RegisterMessage(ListDelegatesRequest)

ListDelegatesResponse = _reflection.GeneratedProtocolMessageType('ListDelegatesResponse', (_message.Message,), {
  'DESCRIPTOR' : _LISTDELEGATESRESPONSE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:ListDelegatesResponse)
  })
_sym_db.RegisterMessage(ListDelegatesResponse)

JSONRPCEnvolope = _reflection.GeneratedProtocolMessageType('JSONRPCEnvolope', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCENVOLOPE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCEnvolope)
  })
_sym_db.RegisterMessage(JSONRPCEnvolope)

JSONRPCCallRequest = _reflection.GeneratedProtocolMessageType('JSONRPCCallRequest', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCCALLREQUEST,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCCallRequest)
  })
_sym_db.RegisterMessage(JSONRPCCallRequest)

JSONRPCCallResult = _reflection.GeneratedProtocolMessageType('JSONRPCCallResult', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCCALLRESULT,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCCallResult)
  })
_sym_db.RegisterMessage(JSONRPCCallResult)

JSONRPCNotifyRequest = _reflection.GeneratedProtocolMessageType('JSONRPCNotifyRequest', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCNOTIFYREQUEST,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCNotifyRequest)
  })
_sym_db.RegisterMessage(JSONRPCNotifyRequest)

JSONRPCNotifyResponse = _reflection.GeneratedProtocolMessageType('JSONRPCNotifyResponse', (_message.Message,), {
  'DESCRIPTOR' : _JSONRPCNOTIFYRESPONSE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JSONRPCNotifyResponse)
  })
_sym_db.RegisterMessage(JSONRPCNotifyResponse)

MethodInfo = _reflection.GeneratedProtocolMessageType('MethodInfo', (_message.Message,), {
  'DESCRIPTOR' : _METHODINFO,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:MethodInfo)
  })
_sym_db.RegisterMessage(MethodInfo)

CanServeRequest = _reflection.GeneratedProtocolMessageType('CanServeRequest', (_message.Message,), {
  'DESCRIPTOR' : _CANSERVEREQUEST,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:CanServeRequest)
  })
_sym_db.RegisterMessage(CanServeRequest)

CanServeResponse = _reflection.GeneratedProtocolMessageType('CanServeResponse', (_message.Message,), {
  'DESCRIPTOR' : _CANSERVERESPONSE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:CanServeResponse)
  })
_sym_db.RegisterMessage(CanServeResponse)

CanDelegateRequest = _reflection.GeneratedProtocolMessageType('CanDelegateRequest', (_message.Message,), {
  'DESCRIPTOR' : _CANDELEGATEREQUEST,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:CanDelegateRequest)
  })
_sym_db.RegisterMessage(CanDelegateRequest)

CanDelegateResponse = _reflection.GeneratedProtocolMessageType('CanDelegateResponse', (_message.Message,), {
  'DESCRIPTOR' : _CANDELEGATERESPONSE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:CanDelegateResponse)
  })
_sym_db.RegisterMessage(CanDelegateResponse)

PING = _reflection.GeneratedProtocolMessageType('PING', (_message.Message,), {
  'DESCRIPTOR' : _PING,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:PING)
  })
_sym_db.RegisterMessage(PING)

PONG = _reflection.GeneratedProtocolMessageType('PONG', (_message.Message,), {
  'DESCRIPTOR' : _PONG,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:PONG)
  })
_sym_db.RegisterMessage(PONG)

TubeState = _reflection.GeneratedProtocolMessageType('TubeState', (_message.Message,), {
  'DESCRIPTOR' : _TUBESTATE,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:TubeState)
  })
_sym_db.RegisterMessage(TubeState)

JointRPCDownPacket = _reflection.GeneratedProtocolMessageType('JointRPCDownPacket', (_message.Message,), {
  'DESCRIPTOR' : _JOINTRPCDOWNPACKET,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JointRPCDownPacket)
  })
_sym_db.RegisterMessage(JointRPCDownPacket)

JointRPCUpPacket = _reflection.GeneratedProtocolMessageType('JointRPCUpPacket', (_message.Message,), {
  'DESCRIPTOR' : _JOINTRPCUPPACKET,
  '__module__' : 'jointrpc_pb2'
  # @@protoc_insertion_point(class_scope:JointRPCUpPacket)
  })
_sym_db.RegisterMessage(JointRPCUpPacket)


DESCRIPTOR._options = None

_JOINTRPC = _descriptor.ServiceDescriptor(
  name='JointRPC',
  full_name='JointRPC',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=1203,
  serialized_end=1495,
  methods=[
  _descriptor.MethodDescriptor(
    name='ListMethods',
    full_name='JointRPC.ListMethods',
    index=0,
    containing_service=None,
    input_type=_LISTMETHODSREQUEST,
    output_type=_LISTMETHODSRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='ListDelegates',
    full_name='JointRPC.ListDelegates',
    index=1,
    containing_service=None,
    input_type=_LISTDELEGATESREQUEST,
    output_type=_LISTDELEGATESRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Call',
    full_name='JointRPC.Call',
    index=2,
    containing_service=None,
    input_type=_JSONRPCCALLREQUEST,
    output_type=_JSONRPCCALLRESULT,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Notify',
    full_name='JointRPC.Notify',
    index=3,
    containing_service=None,
    input_type=_JSONRPCNOTIFYREQUEST,
    output_type=_JSONRPCNOTIFYRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Handle',
    full_name='JointRPC.Handle',
    index=4,
    containing_service=None,
    input_type=_JOINTRPCUPPACKET,
    output_type=_JOINTRPCDOWNPACKET,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_JOINTRPC)

DESCRIPTOR.services_by_name['JointRPC'] = _JOINTRPC

# @@protoc_insertion_point(module_scope)
