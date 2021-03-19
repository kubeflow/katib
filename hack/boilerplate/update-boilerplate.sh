#!/bin/bash

# This script is used to update or add the boilerplate in Go and Python files.

# Generate boilerplate for the Go files
# in  ./cmd ./pkg ./hack and ./test dirs
# Exclude client, gRPC manager, swagger, deepcopy and mock dirs from the search.
find_go_files=$(
  find ./cmd ./pkg ./hack ./test -name "*.go" \
    ! -path "./pkg/client/*" \
    ! -path "./pkg/apis/manager/*" \
    ! -path "./pkg/apis/v1beta1/*" \
    ! -path "./pkg/apis/controller/*.deepcopy.go" \
    ! -path "./pkg/mock/*"
)

for i in ${find_go_files}; do
  # If the 2nd line starts with "Copyright" remove the current boilerplate.
  if [[ $(sed -n 2p $i) =~ "Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    tail -n +17 $i >$i.temp
  # Otherwise, add the new boilerplate to the file.
  else
    echo "Add the new boilerplate to $i"
    cat $i >$i.temp
  fi
  cat ./hack/boilerplate/boilerplate.go.txt $i.temp >$i && rm $i.temp
done

# Generate boilerplate for the Python files
# in ./pkg ./cmd ./hack and ./test
# Exclude gRPC manager from the search
find_python_files=$(
  find ./cmd ./pkg ./hack ./test -name "*.py" \
    ! -path "./pkg/apis/manager/*"
)

for i in ${find_python_files}; do
  # If the 1st line starts with "# Copyright" remove the boilerplate.
  if [[ $(sed -n 1p $i) =~ "# Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    tail -n +15 $i >$i.temp
  # Otherwise, add the new boilerplate to the file.
  else
    echo "Add the new boilerplate to $i"
    cat $i >$i.temp
  fi
  # Add new boilerplate to the file.
  cat ./hack/boilerplate/boilerplate.python.txt $i.temp >$i && rm $i.temp
done
