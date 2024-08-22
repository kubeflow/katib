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

from typing import List

from kubeflow.katib import models


def double(min: float, max: float, step: float = None):
    """Sample a float value uniformly between `min` and `max`.

    Args:
        min: Lower boundary for the float value.
        max: Upper boundary for the float value.
        step: Step between float values.
    """

    parameter = models.V1beta1ParameterSpec(
        parameter_type="double",
        feasible_space=models.V1beta1FeasibleSpace(min=str(min), max=str(max)),
    )
    if step is not None:
        parameter.feasible_space.step = str(step)
    return parameter


def int(min: int, max: int, step: int = None):
    """Sample an integer value uniformly between `min` and `max`.

    Args:
        min: Lower boundary for the integer value.
        max: Upper boundary for the integer value.
        step: Step between integer values.
    """

    parameter = models.V1beta1ParameterSpec(
        parameter_type="int",
        feasible_space=models.V1beta1FeasibleSpace(min=str(min), max=str(max)),
    )
    if step is not None:
        parameter.feasible_space.step = str(step)
    return parameter


def categorical(list: List):
    """Sample a categorical value from the `list`.

    Args:
        list: List of categorical values.
    """

    return models.V1beta1ParameterSpec(
        parameter_type="categorical",
        feasible_space=models.V1beta1FeasibleSpace(list),
    )
