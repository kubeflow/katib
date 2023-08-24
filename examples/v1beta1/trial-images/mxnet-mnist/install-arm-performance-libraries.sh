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
set -o nounset
set -o pipefail
cd "$(dirname "$0")"

# Download Arm Performance Libraries for Ubuntu 22.04
# Ref: https://developer.arm.com/downloads/-/arm-performance-libraries
echo "Downloading Arm Performance Libraries for Ubuntu 22.04..."
wget -qO - \
  "https://developer.arm.com/-/media/Files/downloads/hpc/arm-performance-libraries/23-04-1/ubuntu-22/arm-performance-libraries_23.04.1_Ubuntu-22.04_gcc-11.3.tar?rev=207c1f7aaa16400e94eb9a980494a6eb&revision=207c1f7a-aa16-400e-94eb-9a980494a6eb" \
  | tar -xf -

# Install Arm Performance Libraries
echo "Installing Arm Performance Libraries for Ubuntu 22.04..."
./arm-performance-libraries_23.04.1_Ubuntu-22.04/arm-performance-libraries_23.04.1_Ubuntu-22.04.sh -a

# Clean up
echo "Removing installer..."
rm -rf ./arm-performance-libraries_23.04.1_Ubuntu-22.04
