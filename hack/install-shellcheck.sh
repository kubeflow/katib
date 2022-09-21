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

set -o errexit
set -o pipefail
set -o nounset
cd "$(dirname "$0")"

ARCH=$(uname -m)
OS=$(uname)
SHELLCHECK_VERSION=v0.8.0

if [ "$ARCH" = "arm64" ] && [ "$OS" = "Darwin" ]; then
  echo "Please install the shellcheck via Homebrew."
  exit 1
fi

curl -sSL "https://github.com/koalaman/shellcheck/releases/download/${SHELLCHECK_VERSION}/shellcheck-${SHELLCHECK_VERSION}.${OS,,}.${ARCH}.tar.xz" \
  | tar Jxf - -C /tmp
mv /tmp/shellcheck-$SHELLCHECK_VERSION/shellcheck /usr/local/bin/shellcheck
chmod +x /usr/local/bin/shellcheck

rm -rf /tmp/shellcheck-$SHELLCHECK_VERSION
