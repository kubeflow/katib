#!/bin/bash

# This script is used to update or add the boilerplate
# to Python, Go files in ./pkg and ./cmd dirs

# for i in $(find ./cmd ./pkg -name "*.go"); do
#   echo $i
# done

for i in $(find ./hack -name "*.go"); do
  # If 2nd file line starts with "Copyright" remove the boilerplate.
  if [[ $(head -n 2 $i | sed -n 2p) =~ "Copyright" ]]; then
    echo "Remove boilerplate from $i"
    tail -n +17 $i >$i.temp
  else
    cat $i >$i.temp
  fi
  # Add new boilerplate to the file.
  cat ./hack/boilerplate/boilerplate.go.txt $i.temp >$i && rm $i.temp
done
