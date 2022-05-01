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

# This script is used to update or add the boilerplate.

# ------------------ Go files ------------------
# Exclude client, gRPC manager, swagger, deepcopy and mock from the search.
find_go_files=$(
  find ./cmd ./pkg ./hack ./test -name "*.go" \
    ! -path "./pkg/client/*" \
    ! -path "./pkg/apis/manager/*" \
    ! -path "./pkg/apis/v1beta1/*" \
    ! -path "./pkg/apis/controller/*.deepcopy.go" \
    ! -path "./pkg/mock/*"
)

for i in ${find_go_files}; do
  # If the 2nd line starts with "Copyright", remove the current boilerplate.
  if [[ $(sed -n 2p "$i") =~ "Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    tail -n +17 "$i" >"$i.tmp"
  # If the 1st line starts with "//go:build" and the 2nd line starts with "Copyright", remove the current boilerplate and copy the marker for Go.
  elif [[ $(sed -n 1p "$i") =~ "//go:build" ]] && [[ $(sed -n 4p "$i") =~ "Copyright" ]]; then
    echo "Remove the current boilerplate, copy the marker for Go, and add the new boilerplate to $i"
    sed -e "2,17d" "$i" > "$i.tmp"
  # Otherwise, copy the whole file.
  else
    echo "Add the new boilerplate to $i"
    cat "$i" >"$i.tmp"
  fi

  # If the 1st line starts with "//go:build", copy the marker for Go and add the new boilerplate to the file.
  if [[ $(sed -n 1p "$i.tmp") =~ "//go:build" ]]; then
    (head -2 "$i.tmp" && cat ./hack/boilerplate/boilerplate.go.txt && tail -n +3 "$i.tmp") >"$i" && rm "$i.tmp"
  # Otherwise, add the new boilerplate to the file.
  else
    cat ./hack/boilerplate/boilerplate.go.txt "$i.tmp" >"$i" && rm "$i.tmp"
  fi
done

# ------------------ Python files ------------------
# Exclude gRPC manager and __init__.py files from the search.
find_python_files=$(
  find ./cmd ./pkg ./hack ./test ./examples -name "*.py" \
    ! -path "./pkg/apis/manager/*" \
    ! -path "*__init__.py" \
    ! -path "./examples/v1beta1/trial-images/mxnet-mnist/*"
)

for i in ${find_python_files}; do
  # If the 1st line starts with "# Copyright", remove the boilerplate.
  if [[ $(sed -n 1p "$i") =~ "# Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    tail -n +15 "$i" >"$i.tmp"
  # Otherwise, copy the whole file.
  else
    echo "Add the new boilerplate to $i"
    cat "$i" >"$i.tmp"
  fi
  # Add the new boilerplate to the file.
  cat ./hack/boilerplate/boilerplate.py.txt "$i.tmp" >"$i" && rm "$i.tmp"
done

# ------------------ Shell files ------------------
find_shell_files=$(find ./pkg ./hack ./scripts ./test ./examples -name "*.sh")

for i in ${find_shell_files}; do
  # If the 3rd line starts with "# Copyright", remove the boilerplate.
  # In the shell files we should save the first line.
  if [[ $(sed -n 3p "$i") =~ "# Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    sed -e "2,15d" "$i" >"$i.tmp"
  # Otherwise, copy the whole file.
  else
    echo "Add the new boilerplate to $i"
    cat "$i" >"$i.tmp"
  fi
  # Add the new boilerplate to the file.
  (head -2 "$i.tmp" && cat ./hack/boilerplate/boilerplate.sh.txt && tail -n +3 "$i.tmp") >"$i" && rm "$i.tmp"
done
