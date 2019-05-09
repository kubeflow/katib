# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import api_pb2 as api__pb2


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
        request_serializer=api__pb2.GetSuggestionsRequest.SerializeToString,
        response_deserializer=api__pb2.GetSuggestionsReply.FromString,
        )
    self.ValidateAlgorithmSettings = channel.unary_unary(
        '/api.v1.alpha2.Manager/ValidateAlgorithmSettings',
        request_serializer=api__pb2.ValidateAlgorithmSettingsRequest.SerializeToString,
        response_deserializer=api__pb2.ValidateAlgorithmSettingsReply.FromString,
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
          request_deserializer=api__pb2.GetSuggestionsRequest.FromString,
          response_serializer=api__pb2.GetSuggestionsReply.SerializeToString,
      ),
      'ValidateAlgorithmSettings': grpc.unary_unary_rpc_method_handler(
          servicer.ValidateAlgorithmSettings,
          request_deserializer=api__pb2.ValidateAlgorithmSettingsRequest.FromString,
          response_serializer=api__pb2.ValidateAlgorithmSettingsReply.SerializeToString,
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
        request_serializer=api__pb2.GetSuggestionsRequest.SerializeToString,
        response_deserializer=api__pb2.GetSuggestionsReply.FromString,
        )
    self.ValidateAlgorithmSettings = channel.unary_unary(
        '/api.v1.alpha2.Suggestion/ValidateAlgorithmSettings',
        request_serializer=api__pb2.ValidateAlgorithmSettingsRequest.SerializeToString,
        response_deserializer=api__pb2.ValidateAlgorithmSettingsReply.FromString,
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
          request_deserializer=api__pb2.GetSuggestionsRequest.FromString,
          response_serializer=api__pb2.GetSuggestionsReply.SerializeToString,
      ),
      'ValidateAlgorithmSettings': grpc.unary_unary_rpc_method_handler(
          servicer.ValidateAlgorithmSettings,
          request_deserializer=api__pb2.ValidateAlgorithmSettingsRequest.FromString,
          response_serializer=api__pb2.ValidateAlgorithmSettingsReply.SerializeToString,
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
