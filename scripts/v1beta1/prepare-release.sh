#!/usr/bin/env bash

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

# Prepare a Katib release by updating Python package versions.
# Manifest image tags are pinned on the release branch by CI (see .github/workflows/release.yaml).
# Run from the repository root: ./scripts/v1beta1/prepare-release.sh <VERSION>
# For example: ./scripts/v1beta1/prepare-release.sh 0.19.1

set -o errexit
set -o pipefail
set -o nounset

VERSION=${1:-""}

if [[ -z "$VERSION" ]]; then
  echo "Version must be set"
  echo "Usage: $0 <VERSION>" 1>&2
  echo "You must follow this format: X.Y.Z or X.Y.ZrcN" 1>&2
  exit 1
fi

if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(rc[0-9]+)?$ ]]; then
  echo "Version format is invalid: ${VERSION}"
  echo "Usage: $0 <VERSION>" 1>&2
  echo "You must follow this format: X.Y.Z or X.Y.ZrcN" 1>&2
  exit 1
fi

if [[ "$VERSION" =~ rc[0-9]+$ ]]; then
  TAG="v${VERSION/rc/-rc.}"
else
  TAG="v${VERSION}"
fi

SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "${SCRIPT_ROOT}"

echo -e "\nPreparing Katib release. Version: ${VERSION}, Tag: ${TAG}\n"

echo -e "\nUpdating Katib Python SDK version to ${VERSION}\n"
if [[ $(uname) == "Darwin" ]]; then
  sed -i '' -e "s@version=\".*\"@version=\"${VERSION}\"@" sdk/python/v1beta1/setup.py
else
  sed -i -e "s@version=\".*\"@version=\"${VERSION}\"@" sdk/python/v1beta1/setup.py
fi

echo -e "\nUpdating Katib Python models version to ${VERSION}\n"
if [[ $(uname) == "Darwin" ]]; then
  sed -i '' -e "s@API_VERSION=\".*\"@API_VERSION=\"${VERSION}\"@" hack/python-api/gen-api.sh
  sed -i '' -e "s@__version__ = \".*\"@__version__ = \"${VERSION}\"@" api/python_api/kubeflow_katib_api/__init__.py
else
  sed -i -e "s@API_VERSION=\".*\"@API_VERSION=\"${VERSION}\"@" hack/python-api/gen-api.sh
  sed -i -e "s@__version__ = \".*\"@__version__ = \"${VERSION}\"@" api/python_api/kubeflow_katib_api/__init__.py
fi

echo -e "\nKatib release ${TAG} has been prepared. Review changes and open a PR."
