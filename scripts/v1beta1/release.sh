#!/bin/bash

# Copyright 2021 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

usage() {
  echo "Usage: $0 [-b <BRANCH>] [-t <TAG>]" 1>&2
  echo "You must follow this format, Branch: release-X.Y, Tag: vX.Y.Z"
  exit 1
}

while getopts ":b::t:" opt; do
  case $opt in
  b)
    BRANCH=${OPTARG}
    ;;
  t)
    TAG=${OPTARG}
    ;;
  *)
    usage
    ;;
  esac
done

if [[ -z "$BRANCH" || -z "$TAG" ]]; then
  echo "Branch and Tag must be set"
  echo "Usage: $0 [-b <BRANCH>] [-t <TAG>]" 1>&2
  echo "You must follow this format, Branch: release-X.Y, Tag: vX.Y.Z"
  exit 1
fi

# Clone Katib repo to temp dir.
temp_dir=$(mktemp -d)
git clone "git@github.com:andreyvelich/test-argocd.git" ${temp_dir}""
cd $temp_dir

# Check if tag exists.
if [[ ! -z $(git tag --list ${TAG}) ]]; then
  echo "Tag: ${TAG} exists. Release can't be published"
  exit 1
fi

echo -e "\nCreating new release. Branch: ${BRANCH}, TAG: ${TAG}\n"

# Create or use the branch.
if [[ -z $(git branch -r -l origin/${BRANCH}) ]]; then
  echo "Branch: ${BRANCH} does not exist. Creating a new minor release"
  git checkout -b ${BRANCH}
else
  echo "Branch: ${BRANCH} exists. Creating a new patch release"
  git checkout ${BRANCH}
  read -p "Did you cherry pick all commits from the master to the ${BRANCH} branch? [y|n] "
  if [ "$REPLY" != "y" ]; then
    exit 1
  fi
fi

# Change Katib image tags to the release ${TAG}.
# Get current image tag.
current_tag=$(cat ./manifests/v1beta1/installs/katib-standalone/kustomization.yaml | grep -m 1 "newTag:" | awk '{print $2}')
echo -e "\nUpdating Katib image tags from ${current_tag} to ${TAG}"

# For MacOS we should set -i '' to avoid temp files from sed.
if [[ $(uname) == "Darwin" ]]; then
  find ./manifests/v1beta1/installs -regex ".*\.yaml" -exec sed -i '' -e "s@${current_tag}@${TAG}@" {} \;
else
  find ./manifests/v1beta1/installs -regex ".*\.yaml" -exec sed -i -e "s@${current_tag}@${TAG}@" {} \;
fi
echo -e "Katib images have been updated\n"

git commit -a -m "Katib official release ${TAG}"

# Create new tag.
git tag ${TAG}

# Publish images to the registry with 2 tags: ${TAG} and v1beta1-<commit-sha>.
# ---------------------------------

read -p "Do you want to push Katib ${TAG} version to upstream? [y|n] "
if [ "$REPLY" != "y" ]; then
  exit 1
fi

# Push a new Branch and Tag.
git push -u origin ${BRANCH}
git push -u origin ${TAG}

echo -e "\nKatib ${TAG} has been released"
