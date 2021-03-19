#!/usr/bin/env bash

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

cd ${SCRIPT_ROOT}
doctoc --github ./README.md --title "
# Table of Contents"

cd - >/dev/null
