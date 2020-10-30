import grpc
import grpc_testing
import unittest

from pkg.apis.manager.v1beta1.python import api_pb2

from pkg.earlystopping.v1beta1.medianstop.service import MedianStopService


class TestMedianStop(unittest.TestCase):
    def setUp(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name['EarlyStopping']: MedianStopService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time())

    def test_get_earlystopping_rules(self):

        # TODO (andreyvelich): Add more informative tests.
        trials = [
            api_pb2.Trial(
                name="test-asfjh",
            ),
            api_pb2.Trial(
                name="test-234hs",
            )
        ]

        experiment = api_pb2.Experiment(
            name="test",
        )

        request = api_pb2.GetEarlyStoppingRulesRequest(
            experiment=experiment,
            trials=trials,
            db_manager_address="katib-db-manager.kubeflow:6789"
        )

        get_earlystopping_rules = self.test_server.invoke_unary_unary(
            method_descriptor=(api_pb2.DESCRIPTOR
                               .services_by_name['EarlyStopping']
                               .methods_by_name['GetEarlyStoppingRules']),
            invocation_metadata={},
            request=request, timeout=1)

        response, metadata, code, details = get_earlystopping_rules.termination()

        self.assertEqual(code, grpc.StatusCode.OK)


if __name__ == '__main__':
    unittest.main()
