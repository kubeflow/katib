# Katib 2019 Roadmap

This document provides a high level view of where Katib will likely grow in 2019.

Katib's Critical User Journey (CUJ) can be found [here](https://bit.ly/2QNKMwt).

# Katib 1.0 Readiness

## Stabilize APIs for StudyJobs
* Beta by end of Q2, 1.0 by end of Q4
* Formalize naming conventions (we use different names like katib vs vizier in different places)
* Fully integrate katib with existing E2E examples:
	* Xgboost
	* Mnist
	* GitHub issue summarization
* Publish API documentation, best practices, tutorials
* [Issues list](https://github.com/kubeflow/katib/issues)


# Enhance HP Tuning Experience

The objectives here are organized around the three stages defined in the CUJ:

## Defining Model and Parameters

Integration with KF distributed training components
* TFJob
* PyTorch
* Allow Katib to support other operator types generically

## Configuring a Study
* Streamlining the StudyJob schema - providing simpler ways to write worker specs and metric collector specs.
* Integration with Jupyter notebooks and Fairing
* Allow users to start with an existing model from a notebook and do HP tuning with minimal code changes
* Generating StudyJob configurations and launching StudyJobs through UI
* Supporting additional suggestion algorithms
* Support for StudyJob deployment in a different namespace


## Tracking Model Performance
* Enhance metrics collection
	* May need to revisit the design - use a push model instead of pull model?
* UI enhancements: allowing data scientists to visualize results easier
* [TODO] What to do about Data? Need to figure out a long term story for the model DB


# Other Features

* Multi-Tenancy Support
* NAS
* Batch scheduling
* Integration with Pipelines
* Early stopping feature

# Test and Release Infrastructure

* Improve e2e test coverage
* Improve test harness
* Enhance release process; adding automation (see https://bit.ly/2F7o4gM) 
