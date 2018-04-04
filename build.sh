#/bin/bash
set -x
set -e
PREFIX="katib/"
docker build -t ${PREFIX}vizier-core -f manager/Dockerfile .
docker build -t ${PREFIX}suggestion-random -f suggestion/random/Dockerfile .
docker build -t ${PREFIX}suggestion-grid -f suggestion/grid/Dockerfile .
docker build -t ${PREFIX}suggestion-hyperband -f suggestion/hyperband/Dockerfile .
docker build -t ${PREFIX}dlk-manager -f vendor/github.com/osrg/dlk/build/Dockerfile vendor/github.com/osrg/dlk
docker build -t ${PREFIX}katib-frontend -f manager/modeldb/Dockerfile .
docker build -t ${PREFIX}katib-cli -f cli/Dockerfile .
mkdir -p bin
docker run --name katib-cli -itd ${PREFIX}katib-cli sh
docker cp katib-cli:/go/src/github.com/mlkube/katib/cli/katib-cli bin/katib-cli 
docker rm -f katib-cli
