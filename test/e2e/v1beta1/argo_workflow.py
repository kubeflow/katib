# Copyright 2021 The Kubeflow Authors.
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


def build_init_workflow(name, registry):
    pass


def create_workflow(name=None, **kwargs):
    """Main function which returns Argo Workflow.

    :param name: Argo Workflow name.
    :param kwargs: Argo Workflow additional arguments.

    :return: Created Argo Workflow.
    :rtype: dict
    """

    workflow = build_init_workflow(name, kwargs[""])

    return workflow
