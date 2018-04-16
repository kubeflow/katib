#/bin/bash
set -x
set -e
PREFIX="katib/"
docker build -t ${PREFIX}vizier-core -f manager/Dockerfile .
docker build -t ${PREFIX}suggestion-random -f suggestion/random/Dockerfile .
docker build -t ${PREFIX}suggestion-grid -f suggestion/grid/Dockerfile .
docker build -t ${PREFIX}suggestion-hyperband -f suggestion/hyperband/Dockerfile .
docker build -t ${PREFIX}earlystopping-medianstopping -f earlystopping/medianstopping/Dockerfile .
docker build -t ${PREFIX}dlk-manager -f dlk/Dockerfile .
docker build -t ${PREFIX}katib-frontend -f modeldb/Dockerfile .
docker build -t ${PREFIX}katib-cli -f cli/Dockerfile .
