# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import numpy as np
from sklearn.preprocessing import MinMaxScaler


class ParameterConfig:
    """
    Class to hold the parameter configuration for an experiment.

    Attributes:
        name_ids (dict): Mapping from a parameter name to the index of that
            parameter in the other fields.
        dim (int): Dimension of the vectors created when parameter assignments
            are mapped to a vector. Each int, double, or discrete parameter
            adds one to the dimension, and each categorical parameter adds
            the number of feasible values for that parameter due to one-hot
            encoding.
        lower_bounds (ndarray): The lower bounds for each parameter in the
            search space.
        upper_bounds (ndarray): The lower bounds for each parameter in the
            search space.
        parameter_types (list): The type of each parameter.
        names (list): The name of each parameter.
        discrete_info (list): A list of dicts where each dict contains the
            information for a single discrete parameter. An example of a dict
            is {"name": "discrete_parameter, "values": [2, 3, 5]}]
        categorical_info (list): A list of dicts where each dict contains the
            information for a single categorical parameter. An example dict is
            {"name": "cat_param", "values": ["true", "false"], "number": 2}.
    """

    def __init__(self, name_ids, dim, lower_bounds, upper_bounds,
                 parameter_types, names, discrete_info, categorical_info):
        self.name_ids = name_ids
        self.dim = dim
        self.lower_bounds = np.array(lower_bounds).reshape((1, dim))
        self.upper_bounds = np.array(upper_bounds).reshape((1, dim))
        self.parameter_types = parameter_types
        self.names = names
        self.discrete_info = discrete_info
        self.categorical_info = categorical_info
        if len(self.names) != len(set(self.names)):
            raise Exception("Parameter names are not unique.")

    def create_scaler(self):
        search_space = np.append(self.lower_bounds, self.upper_bounds, axis=0)
        scaler = MinMaxScaler()
        scaler.fit(search_space)
        return scaler

    def random_sample(self):
        new_sample = np.random.uniform(self.lower_bounds, self.upper_bounds,
                                       size=(1, self.dim))
        return new_sample
