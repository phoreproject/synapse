# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: relayerrpc.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='relayerrpc.proto',
  package='pb',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x10relayerrpc.proto\x12\x02pb2\x0c\n\nRelayerRPCb\x06proto3')
)



_sym_db.RegisterFileDescriptor(DESCRIPTOR)



_RELAYERRPC = _descriptor.ServiceDescriptor(
  name='RelayerRPC',
  full_name='pb.RelayerRPC',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=24,
  serialized_end=36,
  methods=[
])
_sym_db.RegisterServiceDescriptor(_RELAYERRPC)

DESCRIPTOR.services_by_name['RelayerRPC'] = _RELAYERRPC

# @@protoc_insertion_point(module_scope)