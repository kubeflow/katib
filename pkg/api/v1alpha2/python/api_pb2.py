# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: api.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from dbif import dbif_pb2 as dbif_dot_dbif__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='api.proto',
  package='api.v1.alpha2',
  syntax='proto3',
  serialized_pb=_b('\n\tapi.proto\x12\rapi.v1.alpha2\x1a\x1cgoogle/api/annotations.proto\x1a\x0f\x64\x62if/dbif.proto\"`\n\x15GetSuggestionsRequest\x12\x17\n\x0f\x65xperiment_name\x18\x01 \x01(\t\x12\x16\n\x0e\x61lgorithm_name\x18\x02 \x01(\t\x12\x16\n\x0erequest_number\x18\x03 \x01(\x05\"<\n\x13GetSuggestionsReply\x12%\n\x06trials\x18\x01 \x03(\x0b\x32\x15.dbif.v1.alpha2.Trial\"s\n ValidateAlgorithmSettingsRequest\x12\x37\n\x0f\x65xperiment_spec\x18\x01 \x01(\x0b\x32\x1e.dbif.v1.alpha2.ExperimentSpec\x12\x16\n\x0e\x61lgorithm_name\x18\x02 \x01(\t\" \n\x1eValidateAlgorithmSettingsReply2\xbf\x02\n\x07Manager\x12\x82\x01\n\x0eGetSuggestions\x12$.api.v1.alpha2.GetSuggestionsRequest\x1a\".api.v1.alpha2.GetSuggestionsReply\"&\x82\xd3\xe4\x93\x02 \"\x1b/api/Manager/GetSuggestions:\x01*\x12\xae\x01\n\x19ValidateAlgorithmSettings\x12/.api.v1.alpha2.ValidateAlgorithmSettingsRequest\x1a-.api.v1.alpha2.ValidateAlgorithmSettingsReply\"1\x82\xd3\xe4\x93\x02+\"&/api/Manager/ValidateAlgorithmSettings:\x01*2\xe5\x01\n\nSuggestion\x12Z\n\x0eGetSuggestions\x12$.api.v1.alpha2.GetSuggestionsRequest\x1a\".api.v1.alpha2.GetSuggestionsReply\x12{\n\x19ValidateAlgorithmSettings\x12/.api.v1.alpha2.ValidateAlgorithmSettingsRequest\x1a-.api.v1.alpha2.ValidateAlgorithmSettingsReply2\x0f\n\rEarlyStoppingb\x06proto3')
  ,
  dependencies=[google_dot_api_dot_annotations__pb2.DESCRIPTOR,dbif_dot_dbif__pb2.DESCRIPTOR,])




_GETSUGGESTIONSREQUEST = _descriptor.Descriptor(
  name='GetSuggestionsRequest',
  full_name='api.v1.alpha2.GetSuggestionsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='experiment_name', full_name='api.v1.alpha2.GetSuggestionsRequest.experiment_name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='algorithm_name', full_name='api.v1.alpha2.GetSuggestionsRequest.algorithm_name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='request_number', full_name='api.v1.alpha2.GetSuggestionsRequest.request_number', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=75,
  serialized_end=171,
)


_GETSUGGESTIONSREPLY = _descriptor.Descriptor(
  name='GetSuggestionsReply',
  full_name='api.v1.alpha2.GetSuggestionsReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='trials', full_name='api.v1.alpha2.GetSuggestionsReply.trials', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=173,
  serialized_end=233,
)


_VALIDATEALGORITHMSETTINGSREQUEST = _descriptor.Descriptor(
  name='ValidateAlgorithmSettingsRequest',
  full_name='api.v1.alpha2.ValidateAlgorithmSettingsRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='experiment_spec', full_name='api.v1.alpha2.ValidateAlgorithmSettingsRequest.experiment_spec', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='algorithm_name', full_name='api.v1.alpha2.ValidateAlgorithmSettingsRequest.algorithm_name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=235,
  serialized_end=350,
)


_VALIDATEALGORITHMSETTINGSREPLY = _descriptor.Descriptor(
  name='ValidateAlgorithmSettingsReply',
  full_name='api.v1.alpha2.ValidateAlgorithmSettingsReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=352,
  serialized_end=384,
)

_GETSUGGESTIONSREPLY.fields_by_name['trials'].message_type = dbif_dot_dbif__pb2._TRIAL
_VALIDATEALGORITHMSETTINGSREQUEST.fields_by_name['experiment_spec'].message_type = dbif_dot_dbif__pb2._EXPERIMENTSPEC
DESCRIPTOR.message_types_by_name['GetSuggestionsRequest'] = _GETSUGGESTIONSREQUEST
DESCRIPTOR.message_types_by_name['GetSuggestionsReply'] = _GETSUGGESTIONSREPLY
DESCRIPTOR.message_types_by_name['ValidateAlgorithmSettingsRequest'] = _VALIDATEALGORITHMSETTINGSREQUEST
DESCRIPTOR.message_types_by_name['ValidateAlgorithmSettingsReply'] = _VALIDATEALGORITHMSETTINGSREPLY
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

GetSuggestionsRequest = _reflection.GeneratedProtocolMessageType('GetSuggestionsRequest', (_message.Message,), dict(
  DESCRIPTOR = _GETSUGGESTIONSREQUEST,
  __module__ = 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.v1.alpha2.GetSuggestionsRequest)
  ))
_sym_db.RegisterMessage(GetSuggestionsRequest)

GetSuggestionsReply = _reflection.GeneratedProtocolMessageType('GetSuggestionsReply', (_message.Message,), dict(
  DESCRIPTOR = _GETSUGGESTIONSREPLY,
  __module__ = 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.v1.alpha2.GetSuggestionsReply)
  ))
_sym_db.RegisterMessage(GetSuggestionsReply)

ValidateAlgorithmSettingsRequest = _reflection.GeneratedProtocolMessageType('ValidateAlgorithmSettingsRequest', (_message.Message,), dict(
  DESCRIPTOR = _VALIDATEALGORITHMSETTINGSREQUEST,
  __module__ = 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.v1.alpha2.ValidateAlgorithmSettingsRequest)
  ))
_sym_db.RegisterMessage(ValidateAlgorithmSettingsRequest)

ValidateAlgorithmSettingsReply = _reflection.GeneratedProtocolMessageType('ValidateAlgorithmSettingsReply', (_message.Message,), dict(
  DESCRIPTOR = _VALIDATEALGORITHMSETTINGSREPLY,
  __module__ = 'api_pb2'
  # @@protoc_insertion_point(class_scope:api.v1.alpha2.ValidateAlgorithmSettingsReply)
  ))
_sym_db.RegisterMessage(ValidateAlgorithmSettingsReply)



_MANAGER = _descriptor.ServiceDescriptor(
  name='Manager',
  full_name='api.v1.alpha2.Manager',
  file=DESCRIPTOR,
  index=0,
  options=None,
  serialized_start=387,
  serialized_end=706,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetSuggestions',
    full_name='api.v1.alpha2.Manager.GetSuggestions',
    index=0,
    containing_service=None,
    input_type=_GETSUGGESTIONSREQUEST,
    output_type=_GETSUGGESTIONSREPLY,
    options=_descriptor._ParseOptions(descriptor_pb2.MethodOptions(), _b('\202\323\344\223\002 \"\033/api/Manager/GetSuggestions:\001*')),
  ),
  _descriptor.MethodDescriptor(
    name='ValidateAlgorithmSettings',
    full_name='api.v1.alpha2.Manager.ValidateAlgorithmSettings',
    index=1,
    containing_service=None,
    input_type=_VALIDATEALGORITHMSETTINGSREQUEST,
    output_type=_VALIDATEALGORITHMSETTINGSREPLY,
    options=_descriptor._ParseOptions(descriptor_pb2.MethodOptions(), _b('\202\323\344\223\002+\"&/api/Manager/ValidateAlgorithmSettings:\001*')),
  ),
])
_sym_db.RegisterServiceDescriptor(_MANAGER)

DESCRIPTOR.services_by_name['Manager'] = _MANAGER


_SUGGESTION = _descriptor.ServiceDescriptor(
  name='Suggestion',
  full_name='api.v1.alpha2.Suggestion',
  file=DESCRIPTOR,
  index=1,
  options=None,
  serialized_start=709,
  serialized_end=938,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetSuggestions',
    full_name='api.v1.alpha2.Suggestion.GetSuggestions',
    index=0,
    containing_service=None,
    input_type=_GETSUGGESTIONSREQUEST,
    output_type=_GETSUGGESTIONSREPLY,
    options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ValidateAlgorithmSettings',
    full_name='api.v1.alpha2.Suggestion.ValidateAlgorithmSettings',
    index=1,
    containing_service=None,
    input_type=_VALIDATEALGORITHMSETTINGSREQUEST,
    output_type=_VALIDATEALGORITHMSETTINGSREPLY,
    options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_SUGGESTION)

DESCRIPTOR.services_by_name['Suggestion'] = _SUGGESTION


_EARLYSTOPPING = _descriptor.ServiceDescriptor(
  name='EarlyStopping',
  full_name='api.v1.alpha2.EarlyStopping',
  file=DESCRIPTOR,
  index=2,
  options=None,
  serialized_start=940,
  serialized_end=955,
  methods=[
])
_sym_db.RegisterServiceDescriptor(_EARLYSTOPPING)

DESCRIPTOR.services_by_name['EarlyStopping'] = _EARLYSTOPPING

try:
  # THESE ELEMENTS WILL BE DEPRECATED.
  # Please use the generated *_pb2_grpc.py files instead.
  import grpc
  from grpc.beta import implementations as beta_implementations
  from grpc.beta import interfaces as beta_interfaces
  from grpc.framework.common import cardinality
  from grpc.framework.interfaces.face import utilities as face_utilities


  class ManagerStub(object):
    """*
    Service for Main API for Katib
    For each RPC service, we define mapping to HTTP REST API method.
    The mapping includes the URL path, query parameters and request body.
    https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http
    *
    Register a Experiment to DB.
    """

    def __init__(self, channel):
      """Constructor.

      Args:
        channel: A grpc.Channel.
      """
      self.GetSuggestions = channel.unary_unary(
          '/api.v1.alpha2.Manager/GetSuggestions',
          request_serializer=GetSuggestionsRequest.SerializeToString,
          response_deserializer=GetSuggestionsReply.FromString,
          )
      self.ValidateAlgorithmSettings = channel.unary_unary(
          '/api.v1.alpha2.Manager/ValidateAlgorithmSettings',
          request_serializer=ValidateAlgorithmSettingsRequest.SerializeToString,
          response_deserializer=ValidateAlgorithmSettingsReply.FromString,
          )


  class ManagerServicer(object):
    """*
    Service for Main API for Katib
    For each RPC service, we define mapping to HTTP REST API method.
    The mapping includes the URL path, query parameters and request body.
    https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http
    *
    Register a Experiment to DB.
    """

    def GetSuggestions(self, request, context):
      """* 
      Get Suggestions from a Suggestion service.
      """
      context.set_code(grpc.StatusCode.UNIMPLEMENTED)
      context.set_details('Method not implemented!')
      raise NotImplementedError('Method not implemented!')

    def ValidateAlgorithmSettings(self, request, context):
      """* 
      Validate AlgorithmSettings in an Experiment.
      Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid
      """
      context.set_code(grpc.StatusCode.UNIMPLEMENTED)
      context.set_details('Method not implemented!')
      raise NotImplementedError('Method not implemented!')


  def add_ManagerServicer_to_server(servicer, server):
    rpc_method_handlers = {
        'GetSuggestions': grpc.unary_unary_rpc_method_handler(
            servicer.GetSuggestions,
            request_deserializer=GetSuggestionsRequest.FromString,
            response_serializer=GetSuggestionsReply.SerializeToString,
        ),
        'ValidateAlgorithmSettings': grpc.unary_unary_rpc_method_handler(
            servicer.ValidateAlgorithmSettings,
            request_deserializer=ValidateAlgorithmSettingsRequest.FromString,
            response_serializer=ValidateAlgorithmSettingsReply.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        'api.v1.alpha2.Manager', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


  class SuggestionStub(object):
    # missing associated documentation comment in .proto file
    pass

    def __init__(self, channel):
      """Constructor.

      Args:
        channel: A grpc.Channel.
      """
      self.GetSuggestions = channel.unary_unary(
          '/api.v1.alpha2.Suggestion/GetSuggestions',
          request_serializer=GetSuggestionsRequest.SerializeToString,
          response_deserializer=GetSuggestionsReply.FromString,
          )
      self.ValidateAlgorithmSettings = channel.unary_unary(
          '/api.v1.alpha2.Suggestion/ValidateAlgorithmSettings',
          request_serializer=ValidateAlgorithmSettingsRequest.SerializeToString,
          response_deserializer=ValidateAlgorithmSettingsReply.FromString,
          )


  class SuggestionServicer(object):
    # missing associated documentation comment in .proto file
    pass

    def GetSuggestions(self, request, context):
      # missing associated documentation comment in .proto file
      pass
      context.set_code(grpc.StatusCode.UNIMPLEMENTED)
      context.set_details('Method not implemented!')
      raise NotImplementedError('Method not implemented!')

    def ValidateAlgorithmSettings(self, request, context):
      # missing associated documentation comment in .proto file
      pass
      context.set_code(grpc.StatusCode.UNIMPLEMENTED)
      context.set_details('Method not implemented!')
      raise NotImplementedError('Method not implemented!')


  def add_SuggestionServicer_to_server(servicer, server):
    rpc_method_handlers = {
        'GetSuggestions': grpc.unary_unary_rpc_method_handler(
            servicer.GetSuggestions,
            request_deserializer=GetSuggestionsRequest.FromString,
            response_serializer=GetSuggestionsReply.SerializeToString,
        ),
        'ValidateAlgorithmSettings': grpc.unary_unary_rpc_method_handler(
            servicer.ValidateAlgorithmSettings,
            request_deserializer=ValidateAlgorithmSettingsRequest.FromString,
            response_serializer=ValidateAlgorithmSettingsReply.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        'api.v1.alpha2.Suggestion', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


  class EarlyStoppingStub(object):
    """TODO
    """

    def __init__(self, channel):
      """Constructor.

      Args:
        channel: A grpc.Channel.
      """


  class EarlyStoppingServicer(object):
    """TODO
    """


  def add_EarlyStoppingServicer_to_server(servicer, server):
    rpc_method_handlers = {
    }
    generic_handler = grpc.method_handlers_generic_handler(
        'api.v1.alpha2.EarlyStopping', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


  class BetaManagerServicer(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    """*
    Service for Main API for Katib
    For each RPC service, we define mapping to HTTP REST API method.
    The mapping includes the URL path, query parameters and request body.
    https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http
    *
    Register a Experiment to DB.
    """
    def GetSuggestions(self, request, context):
      """* 
      Get Suggestions from a Suggestion service.
      """
      context.code(beta_interfaces.StatusCode.UNIMPLEMENTED)
    def ValidateAlgorithmSettings(self, request, context):
      """* 
      Validate AlgorithmSettings in an Experiment.
      Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid
      """
      context.code(beta_interfaces.StatusCode.UNIMPLEMENTED)


  class BetaManagerStub(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    """*
    Service for Main API for Katib
    For each RPC service, we define mapping to HTTP REST API method.
    The mapping includes the URL path, query parameters and request body.
    https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http
    *
    Register a Experiment to DB.
    """
    def GetSuggestions(self, request, timeout, metadata=None, with_call=False, protocol_options=None):
      """* 
      Get Suggestions from a Suggestion service.
      """
      raise NotImplementedError()
    GetSuggestions.future = None
    def ValidateAlgorithmSettings(self, request, timeout, metadata=None, with_call=False, protocol_options=None):
      """* 
      Validate AlgorithmSettings in an Experiment.
      Suggestion service should return INVALID_ARGUMENT Error when the parameter is invalid
      """
      raise NotImplementedError()
    ValidateAlgorithmSettings.future = None


  def beta_create_Manager_server(servicer, pool=None, pool_size=None, default_timeout=None, maximum_timeout=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_deserializers = {
      ('api.v1.alpha2.Manager', 'GetSuggestions'): GetSuggestionsRequest.FromString,
      ('api.v1.alpha2.Manager', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsRequest.FromString,
    }
    response_serializers = {
      ('api.v1.alpha2.Manager', 'GetSuggestions'): GetSuggestionsReply.SerializeToString,
      ('api.v1.alpha2.Manager', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsReply.SerializeToString,
    }
    method_implementations = {
      ('api.v1.alpha2.Manager', 'GetSuggestions'): face_utilities.unary_unary_inline(servicer.GetSuggestions),
      ('api.v1.alpha2.Manager', 'ValidateAlgorithmSettings'): face_utilities.unary_unary_inline(servicer.ValidateAlgorithmSettings),
    }
    server_options = beta_implementations.server_options(request_deserializers=request_deserializers, response_serializers=response_serializers, thread_pool=pool, thread_pool_size=pool_size, default_timeout=default_timeout, maximum_timeout=maximum_timeout)
    return beta_implementations.server(method_implementations, options=server_options)


  def beta_create_Manager_stub(channel, host=None, metadata_transformer=None, pool=None, pool_size=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_serializers = {
      ('api.v1.alpha2.Manager', 'GetSuggestions'): GetSuggestionsRequest.SerializeToString,
      ('api.v1.alpha2.Manager', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsRequest.SerializeToString,
    }
    response_deserializers = {
      ('api.v1.alpha2.Manager', 'GetSuggestions'): GetSuggestionsReply.FromString,
      ('api.v1.alpha2.Manager', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsReply.FromString,
    }
    cardinalities = {
      'GetSuggestions': cardinality.Cardinality.UNARY_UNARY,
      'ValidateAlgorithmSettings': cardinality.Cardinality.UNARY_UNARY,
    }
    stub_options = beta_implementations.stub_options(host=host, metadata_transformer=metadata_transformer, request_serializers=request_serializers, response_deserializers=response_deserializers, thread_pool=pool, thread_pool_size=pool_size)
    return beta_implementations.dynamic_stub(channel, 'api.v1.alpha2.Manager', cardinalities, options=stub_options)


  class BetaSuggestionServicer(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    # missing associated documentation comment in .proto file
    pass
    def GetSuggestions(self, request, context):
      # missing associated documentation comment in .proto file
      pass
      context.code(beta_interfaces.StatusCode.UNIMPLEMENTED)
    def ValidateAlgorithmSettings(self, request, context):
      # missing associated documentation comment in .proto file
      pass
      context.code(beta_interfaces.StatusCode.UNIMPLEMENTED)


  class BetaSuggestionStub(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    # missing associated documentation comment in .proto file
    pass
    def GetSuggestions(self, request, timeout, metadata=None, with_call=False, protocol_options=None):
      # missing associated documentation comment in .proto file
      pass
      raise NotImplementedError()
    GetSuggestions.future = None
    def ValidateAlgorithmSettings(self, request, timeout, metadata=None, with_call=False, protocol_options=None):
      # missing associated documentation comment in .proto file
      pass
      raise NotImplementedError()
    ValidateAlgorithmSettings.future = None


  def beta_create_Suggestion_server(servicer, pool=None, pool_size=None, default_timeout=None, maximum_timeout=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_deserializers = {
      ('api.v1.alpha2.Suggestion', 'GetSuggestions'): GetSuggestionsRequest.FromString,
      ('api.v1.alpha2.Suggestion', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsRequest.FromString,
    }
    response_serializers = {
      ('api.v1.alpha2.Suggestion', 'GetSuggestions'): GetSuggestionsReply.SerializeToString,
      ('api.v1.alpha2.Suggestion', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsReply.SerializeToString,
    }
    method_implementations = {
      ('api.v1.alpha2.Suggestion', 'GetSuggestions'): face_utilities.unary_unary_inline(servicer.GetSuggestions),
      ('api.v1.alpha2.Suggestion', 'ValidateAlgorithmSettings'): face_utilities.unary_unary_inline(servicer.ValidateAlgorithmSettings),
    }
    server_options = beta_implementations.server_options(request_deserializers=request_deserializers, response_serializers=response_serializers, thread_pool=pool, thread_pool_size=pool_size, default_timeout=default_timeout, maximum_timeout=maximum_timeout)
    return beta_implementations.server(method_implementations, options=server_options)


  def beta_create_Suggestion_stub(channel, host=None, metadata_transformer=None, pool=None, pool_size=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_serializers = {
      ('api.v1.alpha2.Suggestion', 'GetSuggestions'): GetSuggestionsRequest.SerializeToString,
      ('api.v1.alpha2.Suggestion', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsRequest.SerializeToString,
    }
    response_deserializers = {
      ('api.v1.alpha2.Suggestion', 'GetSuggestions'): GetSuggestionsReply.FromString,
      ('api.v1.alpha2.Suggestion', 'ValidateAlgorithmSettings'): ValidateAlgorithmSettingsReply.FromString,
    }
    cardinalities = {
      'GetSuggestions': cardinality.Cardinality.UNARY_UNARY,
      'ValidateAlgorithmSettings': cardinality.Cardinality.UNARY_UNARY,
    }
    stub_options = beta_implementations.stub_options(host=host, metadata_transformer=metadata_transformer, request_serializers=request_serializers, response_deserializers=response_deserializers, thread_pool=pool, thread_pool_size=pool_size)
    return beta_implementations.dynamic_stub(channel, 'api.v1.alpha2.Suggestion', cardinalities, options=stub_options)


  class BetaEarlyStoppingServicer(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    """TODO
    """


  class BetaEarlyStoppingStub(object):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This class was generated
    only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0."""
    """TODO
    """


  def beta_create_EarlyStopping_server(servicer, pool=None, pool_size=None, default_timeout=None, maximum_timeout=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_deserializers = {
    }
    response_serializers = {
    }
    method_implementations = {
    }
    server_options = beta_implementations.server_options(request_deserializers=request_deserializers, response_serializers=response_serializers, thread_pool=pool, thread_pool_size=pool_size, default_timeout=default_timeout, maximum_timeout=maximum_timeout)
    return beta_implementations.server(method_implementations, options=server_options)


  def beta_create_EarlyStopping_stub(channel, host=None, metadata_transformer=None, pool=None, pool_size=None):
    """The Beta API is deprecated for 0.15.0 and later.

    It is recommended to use the GA API (classes and functions in this
    file not marked beta) for all further purposes. This function was
    generated only to ease transition from grpcio<0.15.0 to grpcio>=0.15.0"""
    request_serializers = {
    }
    response_deserializers = {
    }
    cardinalities = {
    }
    stub_options = beta_implementations.stub_options(host=host, metadata_transformer=metadata_transformer, request_serializers=request_serializers, response_deserializers=response_deserializers, thread_pool=pool, thread_pool_size=pool_size)
    return beta_implementations.dynamic_stub(channel, 'api.v1.alpha2.EarlyStopping', cardinalities, options=stub_options)
except ImportError:
  pass
# @@protoc_insertion_point(module_scope)
