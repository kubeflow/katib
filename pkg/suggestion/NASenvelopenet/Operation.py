import itertools
import numpy as np
from pkg.api.python import api_pb2

class SearchSpace(object):
    def __init__(self, operations):
        self.operation_list = list(operations.operation)
        self.search_space = dict()
        self._parse_operations()
    
    def _parse_operations(self):

        for operation in self.operation_list:
            opt_spec = list(operation.parameter_configs.configs)
            # avail_space is dict with the format {"spec_nam": [spec feasible values]}
            avail_space = dict()
            num_spec = len(opt_spec)

            for ispec in opt_spec:
                spec_name = ispec.name
                if ispec.parameter_type == api_pb2.CATEGORICAL:
                    avail_space[spec_name] = list(ispec.feasible.list)
                    if len(avail_space[spec_name])==1:
                        avail_space[spec_name]=int(avail_space[spec_name][0])
                elif ispec.parameter_type == api_pb2.INT:
                    avail_space[spec_name] = int(ispec.feasible.value)
        self.search_space = avail_space
