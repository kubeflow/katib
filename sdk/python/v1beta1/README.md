Katib is a Kubernetes-native project designed for automated machine learning. It focuses on hyperparameter tuning, neural architecture search (NAS), and early stopping of machine learning experiments. Katib supports various machine learning frameworks and languages, making it versatile for different applications.

## Requirements
**Python**: >= 3.8

The Katib Python SDK follows the [Python release cycle](https://devguide.python.org/versions/#python-release-cycle) for supported Python versions.

## Installation
You can install the Katib Python SDK using either pip or setuptools.

### Using pip
To install the SDK via pip, run:

```sh
pip install kubeflow-katib
```

After installation, you can import the package in your Python code:

```python
from kubeflow import katib
```

### Using Setuptools
To install via Setuptools, clone the repository and run:

Install via [Setuptools](http://pypi.python.org/pypi/setuptools).

```sh
python setup.py install --user
```
Alternatively, use sudo to install the package for all users:

```sh
sudo python setup.py install
```

### Publish new SDK version to PyPi

The Katib Python SDK is released as part of Katib's patch releases. For each patch release, a new version of the SDK is uploaded to PyPi. The SDK version corresponds directly to the Katib version.

Our SDK is located in [`kubeflow-katib` PyPi package](https://pypi.org/project/kubeflow-katib/).
Katib Python SDK is published as part of the Katib patch releases.
You can check the release process [here](https://github.com/kubeflow/katib/blob/master/scripts/v1beta1/release.sh).
For each Katib patch release, we upload a new SDK version to the PyPi.
The SDK version is equal to the Katib version.


## Getting Started

Please follow the [examples](https://github.com/kubeflow/katib/tree/master/examples/v1beta1/sdk) to learn more about Katib SDK.

## Authorization Details

All endpoints do not require authorization.

## Author

prem0912