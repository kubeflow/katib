#!/bin/sh

for i in 1 2 3 4 5 6; do
    "$@"
    echo restarting
    sleep 10
done
