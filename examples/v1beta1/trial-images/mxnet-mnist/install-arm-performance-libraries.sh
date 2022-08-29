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

# Download Arm Performance Libraries for Ubuntu 20.04
# Ref: https://developer.arm.com/downloads/-/arm-performance-libraries
echo "Downloading Arm Performance Libraries for Ubuntu 20.04..."
wget -qO - \
  "https://developer.arm.com/-/media/Files/downloads/hpc/arm-performance-libraries/22-0-2/Ubuntu20.04/arm-performance-libraries_22.0.2_Ubuntu-20.04_gcc-11.2.tar?rev=577d3dbcff7847b9af57399b2978f9a6&revision=577d3dbc-ff78-47b9-af57-399b2978f9a6" \
  | tar -xf -

# Install Arm Performance Libraries
echo "Installing Arm Performance Libraries for Ubuntu 20.04..."
./arm-performance-libraries_22.0.2_Ubuntu-20.04/arm-performance-libraries_22.0.2_Ubuntu-20.04.sh -a

# Clean up
echo "Removing installer..."
rm -rf ./arm-performance-libraries_22.0.2_Ubuntu-20.04
