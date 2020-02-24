# Katib 2020 Roadmap

This document provides a high level view of where Katib will grow in 2020. 
You can find 2019 Katib's Critical User Journey (CUJ) [here](https://bit.ly/2QNKMwt).

The original Katib design document can be found [here](https://docs.google.com/document/d/1ZEKhou4z1utFTOgjzhSsnvysJFNEJmygllgDCBnYvm8/edit#heading=h.7fzqir88ovr).

# Katib 1.0 Readiness

* Stabilize APIs for Experiments
	* Reconsider the design of Trial Template [#906](https://github.com/kubeflow/katib/issues/906)
	* Early Stopping [#692](https://github.com/kubeflow/katib/issues/692)
	* Resuming Experiment [#1061](https://github.com/kubeflow/katib/issues/1061), [#1062](https://github.com/kubeflow/katib/issues/1062)
* Fully integrate Katib with existing E2E examples:
	* Xgboost
	* Mnist
	* GitHub issue summarization
* Publish API documentation, best practices, tutorials
* [Issues list](https://github.com/kubeflow/katib/issues)

# Enhance HP Tuning Experience

The objectives here are organized around the three stages defined in the CUJ:

## 1. Defining Model and Parameters

Integration with KF distributed training components
* TFJob
* PyTorch
* Allow Katib to support other operator types generically [#341](https://github.com/kubeflow/katib/issues/341)

## 2. Configuring a Experiment
* Supporting additional suggestion algorithms [#15](https://github.com/kubeflow/katib/issues/15)

## 3. Tracking Model Performance
* UI enhancements: allowing data scientists to visualize results easier
* Support for persistent model and metadata storage
	* Ideally users should be able to export and reuse trained models from a common storage

# Test and Release Infrastructure

* Improve e2e test coverage
* Improve test harness
* Enhance release process; adding automation (see https://bit.ly/2F7o4gM) 
