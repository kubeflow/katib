#!/usr/bin/env bash

docker build -t footprintai/transfer-learning-thu:latest -f Dockerfile .
docker push footprintai/transfer-learning-thu:latest
