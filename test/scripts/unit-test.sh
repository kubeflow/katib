set -o errexit
set -o nounset
set -o pipefail

export PATH=${GOPATH}/bin:/usr/local/go/bin:${PATH}
REGISTRY="${GCP_REGISTRY}"
PROJECT="${GCP_PROJECT}"
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}
VERSION=$(git describe --tags --always --dirty)

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Create symlink to GOPATH"
mkdir -p ${GOPATH}/src/github.com/${REPO_OWNER}
ln -s ${PWD} ${GO_DIR}

echo "Run unit test cases"
cd ${GO_DIR}
go test ./...
cd - > /dev/null
