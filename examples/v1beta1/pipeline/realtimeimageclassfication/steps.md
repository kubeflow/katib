### Real time Image Classification Use-Case : End to End Pipeline

### Pre-Requisite
- Follow the **setup-guides/steps-gcp.md** to setup kubeflow and required components. 

### Steps for End-To-End Pipeline


1. set working directory 

Connect to GCP via Local
```bash
gcloud init
gcloud auth application-default login
```

2. Build the container for Data preprocessing

```bash
PROJECT_ID=$(gcloud config get-value core/project)
IMAGE_NAME=kubeflow/hpo
IMAGE_VERSION=v4
IMAGE_NAME=gcr.io/$PROJECT_ID/$IMAGE_NAME
```
Building docker image. 
```
docker build -t $IMAGE_NAME:$IMAGE_VERSION .
```
Push training image to GCR
```
docker push $IMAGE_NAME:$IMAGE_VERSION
```

3. Build the container for Training Model

```bash
cd $WORKDIR/2_HPO_train/
PROJECT_ID=$(gcloud config get-value core/project)
IMAGE_NAME=kubeflow/train
IMAGE_VERSION=v5
IMAGE_NAME=gcr.io/$PROJECT_ID/$IMAGE_NAME
```

Building docker image. 
```
docker build -t $IMAGE_NAME:$IMAGE_VERSION .
```
Push training image to GCR
```
docker push $IMAGE_NAME:$IMAGE_VERSION
```

4. Build the container for Serving Model

```bash

PROJECT_ID=$(gcloud config get-value core/project)
IMAGE_NAME=kubeflow/serving
IMAGE_VERSION=v1 
IMAGE_NAME=gcr.io/$PROJECT_ID/$IMAGE_NAME

```

Building docker image. 
```
docker build -t $IMAGE_NAME:$IMAGE_VERSION .
```
Push training image to GCR
```
docker push $IMAGE_NAME:$IMAGE_VERSION
```

