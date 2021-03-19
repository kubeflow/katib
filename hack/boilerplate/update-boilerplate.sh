#!/bin/bash

# This script is used to update or add the boilerplate
# to Python, Go files in ./pkg and ./cmd dir

# Add boilerplate to Go files
for i in $(find ./hack ./cmd ./pkg -name "*.go"); do
  # If 2nd file line starts with "Copyright" remove the boilerplate.
  if [[ $(head -n 2 $i | sed -n 2p) =~ "Copyright" ]]; then
    echo "Remove and add boilerplate to $i"
    tail -n +17 $i >$i.temp
  else
    echo "Add boilerplate to $i"
    cat $i >$i.temp
  fi
  # Add new boilerplate to the file.
  cat ./hack/boilerplate/boilerplate.go.txt $i.temp >$i && rm $i.temp
done
