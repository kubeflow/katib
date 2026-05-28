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

# This script is used to release Katib project locally.
# For automated releases, use the GitHub Actions "Release" workflow instead.
# Run ./scripts/v1beta1/release.sh <BRANCH> <TAG> to execute it locally.
# For example: ./scripts/v1beta1/release.sh release-0.11 v0.11.1 or
# ./scripts/v1beta1/release.sh release-0.11 v0.11.0-rc.0
# You must follow this format, Branch: release-X.Y, Tag: vX.Y.Z.

set -e

BRANCH=$1
TAG=$2

if [[ -z "$BRANCH" || -z "$TAG" ]]; then
  echo "Branch and Tag must be set"
  echo "Usage: $0 <BRANCH> <TAG>" 1>&2
  echo "You must follow this format, Branch: release-X.Y, Tag: vX.Y.Z or Tag: vX.Y.Z-rc.N"
  exit 1
fi

# Check that Branch and Tag is in correct format.
if [[ ! "$BRANCH" =~ release-[0-9]+\.[0-9]+ || ! "$TAG" =~ v[0-9]+\.[0-9]+\.([0-9]+$|[0-9]+-rc\.[0-9]+$) ]]; then
  echo "Branch or Tag format is invalid"
  echo "Usage: $0 <BRANCH> <TAG>" 1>&2
  echo "You must follow this format, Branch: release-X.Y, Tag: vX.Y.Z or Tag: vX.Y.Z-rc.N"
  exit 1
fi

# Clone Katib repo to the temp dir.
temp_dir=$(mktemp -d)
git clone "git@github.com:kubeflow/katib.git" "${temp_dir}"
cd "$temp_dir"

# Check if tag exists.
if [[ -n $(git tag --list "${TAG}") ]]; then
  echo "Tag: ${TAG} exists. Release can't be published"
  exit 1
fi

echo -e "\nCreating a new release. Branch: ${BRANCH}, Tag: ${TAG}\n"

# Create or use the branch.
if [[ -z $(git branch -r -l "origin/${BRANCH}") ]]; then
  echo "Branch: ${BRANCH} does not exist. Creating a new minor release"
  git checkout -b "${BRANCH}"
else
  echo "Branch: ${BRANCH} exists. Creating a new patch release"
  git checkout "${BRANCH}"
  read -rp "Did you cherry pick all commits from the master to the ${BRANCH} branch? [y|n] "
  if [ "$REPLY" != "y" ]; then
    exit 1
  fi
fi

# Prepare release artifacts (manifests and Python package versions).
sdk_version=${TAG:1}
if [[ ${sdk_version} == *"-rc."* ]]; then
  sdk_version=${sdk_version//-rc./rc}
fi
bash scripts/v1beta1/prepare-release.sh "${sdk_version}"

# ------------------ Publish Katib SDK ------------------
echo -e "\nPublishing Katib Python SDK, version: ${sdk_version}\n"
cd sdk/python/v1beta1
python3 setup.py sdist bdist_wheel
twine upload dist/*
rm -r dist/ build/
cd ../../..
echo -e "\nKatib Python SDK ${sdk_version} has been published\n"

# ------------------ Publish Katib Models for Kubeflow SDK ------------------
echo -e "\nPublishing Katib Python models for Kubeflow SDK, version: ${sdk_version}\n"
./hack/python-api/gen-api.sh
cd api/python_api
python -m build
twine upload dist/*
rm -r dist/
cd ../..
echo -e "\nKatib Python API models for Kubeflow SDK have been published"

# ------------------ Commit changes ------------------
git commit -a -m "Katib official release ${TAG}"
git tag "${TAG}"

# ------------------ Publish Katib images ------------------
echo -e "\nPublishing Katib images\n"
make push-tag TAG="${TAG}"
echo -e "Katib images have been published\n"

# ------------------ Push to upstream ------------------
read -rp "Do you want to push Katib ${TAG} version to the upstream? [y|n] "
if [ "$REPLY" != "y" ]; then
  exit 1
fi
git push -u origin "${BRANCH}"
git push -u origin "${TAG}"

echo -e "\nKatib ${TAG} release has been published"
echo "To finish the release process, follow RELEASE.md: https://github.com/kubeflow/katib/blob/master/RELEASE.md"
