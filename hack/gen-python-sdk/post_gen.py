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

import os
import shutil
import sys

IGNORE_LINES = [
    "from kubeflow.katib.models.v1_unstructured_unstructured import V1UnstructuredUnstructured",
    "from kubeflow.katib.models.v1_time import V1Time",
]


def _rewrite_helper(input_file, output_file, rewrite_rules):
    rules = rewrite_rules or []
    lines = []
    with open(input_file, "r") as f:
        while True:
            line = f.readline()
            if not line:
                break
            # Apply rewrite rules to the line.
            for rule in rules:
                line = rule(line)
            # Remove ignored lines.
            if not any(li in line for li in IGNORE_LINES):
                lines.append(line)

    # Add Katib APIs to the init file.
    if output_file == "sdk/python/v1beta1/kubeflow/katib/__init__.py":
        lines.append("# Import Katib API client.\n")
        lines.append("from kubeflow.katib.api.katib_client import KatibClient\n")
        lines.append("# Import Katib TrainerResources class.\n")
        lines.append("from kubeflow.katib.types.types import TrainerResources\n")
        lines.append("# Import Katib report metrics functions\n")
        lines.append("from kubeflow.katib.api.report_metrics import report_metrics\n")
        lines.append("# Import Katib helper functions.\n")
        lines.append("import kubeflow.katib.api.search as search\n")
        lines.append("# Import Katib helper constants.\n")
        lines.append(
            "from kubeflow.katib.constants.constants import BASE_IMAGE_TENSORFLOW\n"
        )
        lines.append(
            "from kubeflow.katib.constants.constants import BASE_IMAGE_TENSORFLOW_GPU\n"
        )
        lines.append(
            "from kubeflow.katib.constants.constants import BASE_IMAGE_PYTORCH\n"
        )

    # Add Kubernetes models to proper deserialization of Katib models.
    if output_file == "sdk/python/v1beta1/kubeflow/katib/models/__init__.py":
        lines.append("\n")
        lines.append("# Import Kubernetes models.\n")
        lines.append("from kubernetes.client import *\n")

    with open(output_file, "w") as f:
        f.writelines(lines)


def update_python_sdk(src, dest, versions=("v1beta1")):
    # tiny transformers to refine generated codes
    rewrite_rules = [
        # Models rules.
        lambda line: line.replace("import katib", "import kubeflow.katib"),
        lambda line: line.replace("from katib", "from kubeflow.katib"),
        # For the api_client.py.
        lambda line: line.replace(
            "klass = getattr(katib.models, klass)",
            "klass = getattr(kubeflow.katib.models, klass)",
        ),
        # Doc rules.
        lambda line: line.replace("[**datetime**](V1Time.md)", "**datetime**"),
        lambda line: line.replace(
            "[**object**](V1UnstructuredUnstructured.md)", "**object**"
        ),
        lambda line: line.replace(
            "[**V1Container**](V1Container.md)",
            "[**V1Container**](https://github.com/kubernetes-client/"
            "python/blob/master/kubernetes/docs/V1Container.md)",
        ),
        lambda line: line.replace(
            "[**V1ObjectMeta**](V1ObjectMeta.md)",
            "[**V1ObjectMeta**](https://github.com/kubernetes-client/"
            "python/blob/master/kubernetes/docs/V1ObjectMeta.md)",
        ),
        lambda line: line.replace(
            "[**V1ListMeta**](V1ListMeta.md)",
            "[**V1ListMeta**](https://github.com/kubernetes-client/"
            "python/blob/master/kubernetes/docs/V1ListMeta.md)",
        ),
        lambda line: line.replace(
            "[**V1HTTPGetAction**](V1HTTPGetAction.md)",
            "[**V1HTTPGetAction**](https://github.com/kubernetes-client/"
            "python/blob/master/kubernetes/docs/V1HTTPGetAction.md)",
        ),
    ]

    # TODO (andreyvelich): Currently test can't be generated properly.
    src_dirs = [
        os.path.join(src, "katib"),
        os.path.join(src, "katib", "models"),
        # os.path.join(src, 'test'),
        os.path.join(src, "docs"),
    ]
    dest_dirs = [
        os.path.join(dest, "kubeflow", "katib"),
        os.path.join(dest, "kubeflow", "katib", "models"),
        # os.path.join(dest, 'test'),
        os.path.join(dest, "docs"),
    ]

    for src_dir, dest_dir in zip(src_dirs, dest_dirs):
        # Remove previous generated files explicitly, in case of deprecated instances.
        for file in os.listdir(dest_dir):
            path = os.path.join(dest_dir, file)
            # We should not remove KatibClient doc.
            if not os.path.isfile(path) or "/docs/KatibClient.md" in path:
                continue
            for v in versions:
                if v in file.lower():
                    os.remove(path)
                    break
        # fill latest generated files
        for file in os.listdir(src_dir):
            in_file = os.path.join(src_dir, file)
            out_file = os.path.join(dest_dir, file)
            if not os.path.isfile(in_file):
                continue
            _rewrite_helper(in_file, out_file, rewrite_rules)

    # Update doc for Models README.md
    buffer = []
    update_buffer = []

    # Get data from generated doc
    with open(os.path.join(src, "README.md"), "r") as src_f:
        anchor = 0
        for line in src_f.readlines():
            if line.startswith("## Documentation For Models"):
                if anchor == 0:
                    anchor = 1
            elif line.startswith("##") and anchor == 1:
                anchor = 2
            if anchor == 0:
                continue
            if anchor == 2:
                break
            # Remove leading space from the list
            if len(line) > 0:
                line = line.lstrip(" ")
            update_buffer.append(line)
    # Remove latest redundant newline
    update_buffer = update_buffer[:-1]

    # Update README with new models
    with open(os.path.join(dest, "README.md"), "r") as dest_f:
        anchor = 0
        for line in dest_f.readlines():
            if line.startswith("## Documentation For Models"):
                if anchor == 0:
                    buffer.extend(update_buffer)
                    anchor = 1
            elif line.startswith("##") and anchor == 1:
                anchor = 2
            if anchor == 1:
                continue
            buffer.append(line)
    with open(os.path.join(dest, "README.md"), "w") as dest_f:
        dest_f.writelines(buffer)

    # Clear working dictionary
    shutil.rmtree(src)


if __name__ == "__main__":
    update_python_sdk(src=sys.argv[1], dest=sys.argv[2])
