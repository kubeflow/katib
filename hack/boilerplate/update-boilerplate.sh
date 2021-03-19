#!/bin/bash

# This script is used to update or add the boilerplate
# to Python, Go files in ./pkg ./cmd ./hack and ./test dirs.

# Add boilerplate to Go files.
# Exclude client, gRPC manager, swagger, deepcopy and mock dir from the search.
find_go_files=$(
  find ./pkg ./cmd ./hack ./test -name "*.go" \
    ! -path "./pkg/client/*" \
    ! -path "./pkg/apis/manager/*" \
    ! -path "./pkg/apis/v1beta1/*" \
    ! -path "./pkg/apis/controller/*.deepcopy.go" \
    ! -path "./pkg/mock/*"
)

for i in ${find_go_files}; do
  # If the 2nd line starts with "Copyright" remove the current boilerplate.
  if [[ $(head -n 2 $i | sed -n 2p) =~ "Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    tail -n +17 $i >$i.temp
  # Otherwise, add the new boilerplate to the file.
  else
    echo "Add the new boilerplate to $i"
    cat $i >$i.temp
  fi
  cat ./hack/boilerplate/boilerplate.go.txt $i.temp >$i && rm $i.temp
done
