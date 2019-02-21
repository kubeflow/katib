import itertools
import numpy as np
from pkg.api.python import api_pb2


class Operation(object):
    def __init__(self, opt_id, opt_type, opt_params):
        self.opt_id = opt_id
        self.opt_type = opt_type
        self.opt_params = opt_params
    
    def get_dict(self):
        opt_dict = dict()
        opt_dict['opt_id'] = self.opt_id
        opt_dict['opt_type'] = self.opt_type
        opt_dict['opt_params'] = self.opt_params
        return opt_dict

    def print_op(self, logger):
        logger.info("Operation ID: \n\t{}".format(self.opt_id))
        logger.info("Operation Type: \n\t{}".format(self.opt_type))
        logger.info("Operations Parameters:")
        for ikey in self.opt_params:
            logger.info("\t{}: {}".format(ikey, self.opt_params[ikey]))
        logger.info("")


class SearchSpace(object):
    def __init__(self, operations):
        self.operation_list = list(operations.operation)
        self.search_space = list()
        self._parse_operations()
        print()
        self.num_operations = len(self.search_space)
    
    def _parse_operations(self):
        # search_sapce is a list of Operation class

        operation_id = 0

        for operation_dict in self.operation_list:
            opt_type = operation_dict.operationType
            opt_spec = list(operation_dict.parameter_configs.configs)
            # avail_space is dict with the format {"spec_nam": [spec feasible values]}
            avail_space = dict()
            num_spec = len(opt_spec)

            for ispec in opt_spec:
                spec_name = ispec.name
                if ispec.parameter_type == api_pb2.CATEGORICAL:
                    avail_space[spec_name] = list(ispec.feasible.list)
                elif ispec.parameter_type == api_pb2.INT:
                    spec_min = int(ispec.feasible.min)
                    spec_max = int(ispec.feasible.max)
                    spec_step = int(ispec.feasible.step)
                    avail_space[spec_name] = range(spec_min, spec_max+1, spec_step)
                elif ispec.parameter_type == api_pb2.DOUBLE:
                    spec_min = float(ispec.feasible.min)
                    spec_max = float(ispec.feasible.max)
                    spec_step = float(ispec.feasible.step)
                    if spec_step == 0:
                        print("Error, NAS Reinforcement Learning algorithm cannot accept continuous search space!")
                        exit(999)
                    double_list = np.arange(spec_min, spec_max+spec_step, spec_step)
                    if double_list[-1] > spec_max:
                        del double_list[-1]
                    avail_space[spec_name] = double_list

            # generate all the combinations of possible operations
            key_avail_space = list(avail_space.keys())
            val_avail_space = list(avail_space.values())

            for this_opt_vector in itertools.product(*val_avail_space):
                opt_params = dict()
                for i in range(num_spec):
                    opt_params[key_avail_space[i]] = this_opt_vector[i]
                this_opt_class = Operation(operation_id, opt_type, opt_params)
                self.search_space.append(this_opt_class)
                operation_id += 1
