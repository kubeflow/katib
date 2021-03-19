#!/bin/bash

# This script is used to update or add the boilerplate
# to Python, Go files in ./pkg ./cmd and ./hack dirs.

# Add boilerplate to Go files.
# Exclude client, gRPC manager and mock dir from the search.
find_go_files=$(
  find ./pkg ./cmd ./hack -name "*.go" \
    ! -path "./pkg/client/*" \
    ! -path "./pkg/apis/manager*" \
    ! -path "./pkg/mock/*"
)
for i in ${find_go_files}; do
  # If the 2nd line starts with "Copyright" remove the boilerplate.
  if [[ $(head -n 2 $i | sed -n 2p) =~ "Copyright" ]]; then
    echo "Remove the current boilerplate and add the new boilerplate to $i"
    tail -n +17 $i >$i.temp
  # Only add boilerplate if the file doesn't have other license.
  elif ! grep -q "http://www.apache.org/licenses/LICENSE-2.0" $i; then
    echo "Add the new boilerplate to $i"
    cat $i >$i.temp
  else
    echo "Skip file $i"
  fi
  # Add new boilerplate to the file.
  cat ./hack/boilerplate/boilerplate.go.txt $i.temp >$i && rm $i.temp
done
