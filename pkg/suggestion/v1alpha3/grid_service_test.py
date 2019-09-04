from pkg.apis.manager.v1alpha3.python import api_pb2
import grid_service
import unittest

class TestGridSearch(unittest.TestCase):
    def setUp(self):
        pass

    def test_combinations(self):
        grid_ser = grid_service.GridService()
        feasible_space = api_pb2.FeasibleSpace(min="1", max="5")
        parameter_specs = api_pb2.ExperimentSpec.ParameterSpecs()
        parameter_specs.parameters.add(name="lr", parameter_type=api_pb2.DOUBLE, feasible_space=feasible_space)
        comb, _ = grid_ser._create_all_combinations(parameter_specs.parameters, {})
        self.assertEqual(len(comb), 10)
        comb, _ = grid_ser._create_all_combinations(parameter_specs.parameters, {'DefaultGrid':5})
        self.assertEqual(len(comb), 5)
        self.assertEqual(len(comb[0]), 1)
        parameter_specs.parameters.add(name="iterations", parameter_type=api_pb2.INT, feasible_space=feasible_space)
        comb, _ = grid_ser._create_all_combinations(parameter_specs.parameters, {'DefaultGrid':5})
        self.assertEqual(len(comb), 25)
        self.assertEqual(len(comb[0]),2)
        comb, _ = grid_ser._create_all_combinations(parameter_specs.parameters, {'DefaultGrid':5, "lr":4})
        self.assertEqual(len(comb), 20)
        self.assertEqual(len(comb[0]),2)

if __name__ == '__main__':
    unittest.main()

