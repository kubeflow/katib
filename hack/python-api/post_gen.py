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
    "from kubeflow_katib_api.models.v1_unstructured_unstructured import V1UnstructuredUnstructured",
    "from kubeflow_katib_api.models.v1_time import V1Time",
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

    # Add Kubernetes models to proper deserialization of Katib models.
    if output_file.endswith("kubeflow_katib_api/models/__init__.py"):
        lines.append("\n")
        lines.append("# Import Kubernetes models.\n")
        lines.append("from kubernetes.client import *\n")

    with open(output_file, "w") as f:
        f.writelines(lines)


def update_python_api(src, dest, versions=("v1beta1")):
    # tiny transformers to refine generated codes
    rewrite_rules = [
        # Models rules.
        lambda line: line.replace("import kubeflow_katib_api", "import kubeflow_katib_api"),
        lambda line: line.replace("from kubeflow_katib_api", "from kubeflow_katib_api"),
        # For the api_client.py.
        lambda line: line.replace(
            "klass = getattr(kubeflow_katib_api.models, klass)",
            "klass = getattr(kubeflow_katib_api.models, klass)",
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
        os.path.join(src, "kubeflow_katib_api"),
        os.path.join(src, "kubeflow_katib_api", "models"),
        # os.path.join(src, 'test'),
        os.path.join(src, "docs"),
    ]
    dest_dirs = [
        os.path.join(dest, "kubeflow_katib_api"),
        os.path.join(dest, "kubeflow_katib_api", "models"),
        # os.path.join(dest, 'test'),
        os.path.join(dest, "docs"),
    ]

    for src_dir, dest_dir in zip(src_dirs, dest_dirs):
        # Create destination directory if it doesn't exist
        os.makedirs(dest_dir, exist_ok=True)
        
        # Remove previous generated files explicitly, in case of deprecated instances.
        if os.path.exists(dest_dir):
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

    # Skip README update for API package (it's maintained separately)

    # Clear working dictionary
    shutil.rmtree(src)


if __name__ == "__main__":
    update_python_api(src=sys.argv[1], dest=sys.argv[2])
