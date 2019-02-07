# Katib 2019 Roadmap

This document provides a high level view of where Katib will grow in 2019. These objectives are based on Katib's Critical User Journey (CUJ),
which can be found [here](https://bit.ly/2QNKMwt).

The original Katib design document can be found [here](https://docs.google.com/document/d/1ZEKhou4z1utFTOgjzhSsnvysJFNEJmygllgDCBnYvm8/edit#heading=h.7fzqir88ovr).

# Katib 1.0 Readiness

* Stabilize APIs for StudyJobs
	* Beta by end of Q2, 1.0 by end of Q4
	* Formalize naming conventions (we use different names like katib vs vizier in different places)
	* Refactor studyjob field names [#351](https://github.com/kubeflow/katib/issues/351)
	* Rename fields so their names are more meaningful (e.g. requestCount vs requestNumber) [#161](https://github.com/kubeflow/katib/issues/161)
* Fully integrate katib with existing E2E examples:
	* Xgboost
	* Mnist
	* GitHub issue summarization
* Publish API documentation, best practices, tutorials
* [Issues list](https://github.com/kubeflow/katib/issues)
* [Issues for 0.5.0 release](https://github.com/kubeflow/katib/labels/area%2F0.5.0)


# Enhance HP Tuning Experience

The objectives here are organized around the three stages defined in the CUJ:

## 1. Defining Model and Parameters

Integration with KF distributed training components
* TFJob
* PyTorch
* Allow Katib to support other operator types generically [#341](https://github.com/kubeflow/katib/issues/341)

## 2. Configuring a Study
* Streamlining the StudyJob schema - providing simpler ways to write worker specs and metric collector specs.
* Expose more information in StudyJob status fields
	* List all job conditions with details [#344](https://github.com/kubeflow/katib/issues/344)
	* Returning study metadata such as number of trials and best hyperparameter values so far [#356](https://github.com/kubeflow/katib/issues/356)
* Integration with Jupyter notebooks and Fairing [#355](https://github.com/kubeflow/katib/issues/355)
	* Allow users to start with an existing model from a notebook and do HP tuning with minimal code changes
* Allowing a StudyJob to be resumed with additional trials [#346](https://github.com/kubeflow/katib/issues/346)
* Generating StudyJob configurations and launching StudyJobs through UI
* Supporting additional suggestion algorithms [#15](https://github.com/kubeflow/katib/issues/15)
* Support for StudyJob deployment in a different namespace [#343](https://github.com/kubeflow/katib/issues/343)


## 3. Tracking Model Performance
* Enhance metrics collection
	* May need to revisit the design - use a push model instead of pull model?
* UI enhancements: allowing data scientists to visualize results easier
* Support for persistent model and metadata storage
	* Ideally users should be able to export and reuse trained models from a common storage


# Other Features

Designs are pending for the following new features:
* Multi-Tenancy Support
* [NAS](https://docs.google.com/document/d/1qGWy-C5XSQmh82XYoMcJ_JWLHwmyvdMRjCkFMfkO0vE/edit)
* Batch scheduling
* [Integration with Pipelines](https://github.com/kubeflow/katib/issues/331)
* Early stopping feature

# Test and Release Infrastructure

* Improve e2e test coverage
* Improve test harness
* Enhance release process; adding automation (see https://bit.ly/2F7o4gM) 
