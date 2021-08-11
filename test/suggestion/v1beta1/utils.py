
from pkg.apis.manager.v1beta1.python import api_pb2


def call_validate(test_server, experiment_spec):
    experiment = api_pb2.Experiment(name="validation-test", spec=experiment_spec)

    request = api_pb2.ValidateAlgorithmSettingsRequest(experiment=experiment)
    validate_algorithm_settings = test_server.invoke_unary_unary(
        method_descriptor=(api_pb2.DESCRIPTOR
                           .services_by_name['Suggestion']
                           .methods_by_name['ValidateAlgorithmSettings']),
        invocation_metadata={},
        request=request, timeout=1)

    response, metadata, code, details = validate_algorithm_settings.termination()

    return response, metadata, code, details
