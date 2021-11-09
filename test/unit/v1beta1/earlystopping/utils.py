from pkg.apis.manager.v1beta1.python import api_pb2


def call_validate(test_server, early_stopping):
    request = api_pb2.ValidateEarlyStoppingSettingsRequest(early_stopping=early_stopping)
    validate_early_stopping_settings = test_server.invoke_unary_unary(
        method_descriptor=(api_pb2.DESCRIPTOR
                           .services_by_name['EarlyStopping']
                           .methods_by_name['ValidateEarlyStoppingSettings']),
        invocation_metadata={},
        request=request, timeout=1)

    response, metadata, code, details = validate_early_stopping_settings.termination()

    return response, metadata, code, details
