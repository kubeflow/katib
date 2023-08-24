# Katib 2022/2023 Roadmap

## AutoML Features

- Support advance HyperParameter tuning algorithms:

  - Population Based Training (PBT) - [#1382](https://github.com/kubeflow/katib/issues/1382)
  - Tree of Parzen Estimators (TPE)
  - Multivariate TPE
  - Sobol’s Quasirandom Sequence
  - Asynchronous Successive Halving - [ASHA](https://arxiv.org/pdf/1810.05934.pdf)

- Support multi-objective optimization - [#1549](https://github.com/kubeflow/katib/issues/1549)
- Support various HP distributions (log-uniform, uniform, normal) - [#1207](https://github.com/kubeflow/katib/issues/1207)
- Support Auto Model Compression - [#460](https://github.com/kubeflow/katib/issues/460)
- Support Auto Feature Engineering - [#475](https://github.com/kubeflow/katib/issues/475)
- Improve Neural Architecture Search design

## Backend and API Enhancements

- Conformance tests for Katib - [#2044](https://github.com/kubeflow/katib/issues/2044)
- Support push-based metrics collection in Katib - [#577](https://github.com/kubeflow/katib/issues/577)
- Support PostgreSQL as a Katib DB - [#915](https://github.com/kubeflow/katib/issues/915)
- Improve Katib scalability - [#1847](https://github.com/kubeflow/katib/issues/1847)
- Promote Katib APIs to the `v1` version
- Support multiple CRD versions (`v1beta1`, `v1`) with conversion webhook

## Improve Katib User Experience

- Simplify Katib Experiment creation with Katib SDK - [#1951](https://github.com/kubeflow/katib/pull/1951)
- Fully migrate to a new Katib UI - [Project 1](https://github.com/kubeflow/katib/projects/1)
- Expose Trial logs in Katib UI - [#971](https://github.com/kubeflow/katib/issues/971)
- Enhance Katib UI visualization metrics for AutoML Experiments
- Improve Katib Config UX - [#2150](https://github.com/kubeflow/katib/issues/2150)

## Integration with Kubeflow Components

- Kubeflow Pipeline as a Katib Trial target - [#1914](https://github.com/kubeflow/katib/issues/1914)
- Improve data passing when Katib Experiment is part of Kubeflow Pipeline - [#1846](https://github.com/kubeflow/katib/issues/1846)

# History

# Katib 2021 Roadmap

## New Features

### AutoML

- Support Population Based Training [#1382](https://github.com/kubeflow/katib/issues/1382)
- Support [ASHA](https://arxiv.org/pdf/1810.05934.pdf)
- Support Auto Model Compression [#460](https://github.com/kubeflow/katib/issues/460)
- Support Auto Feature Engineering [#475](https://github.com/kubeflow/katib/issues/475)
- Various CRDs for HP, NAS and other AutoML techniques.

### UI

- Migrate to the new Katib UI [Project 1](https://github.com/kubeflow/katib/projects/1)
- Hyperparameter importances visualization with fANOVA algorithm

## Enhancements

- Finish AWS CI/CD migration
- Support various parameter distribution [#1207](https://github.com/kubeflow/katib/issues/1207)
- Finish validation for Algorithms [#1126](https://github.com/kubeflow/katib/issues/1126)
- Refactor Hyperband [#1389](https://github.com/kubeflow/katib/issues/1389)
- Support multiple CRD version with conversion webhook
- MLMD integration with Katib Experiments

# Katib 2020 Roadmap

## New Features

### Hyperparameter Tuning

- Support Early Stopping [#692](https://github.com/kubeflow/katib/issues/692)

### Neural Architecture Search

- Support Advanced NAS Algorithms like DARTs, ProxylessNAS [#461](https://github.com/kubeflow/katib/issues/461)

### Other Features

- Support Auto Model Compression [#460](https://github.com/kubeflow/katib/issues/460)
- Support Auto Feature Engineering [#475](https://github.com/kubeflow/katib/issues/475)

## Enhancements

### Common

- Delete Suggestion deployment after Experiment is finished [#1061](https://github.com/kubeflow/katib/issues/1061)
- Save Suggestion state after deployment is deleted [#1062](https://github.com/kubeflow/katib/issues/1062)
- Reconsider the design of Trial Template [#906](https://github.com/kubeflow/katib/issues/906)
- Design an extensible model for integrating with custom resources.
- Add validation for algorithms (a.k.a suggestions) [#1126](https://github.com/kubeflow/katib/issues/1126)
- Katib UI fixes and enhancements
- Investigate Kubeflow Metadata integration
- Investigate the alignment with concept and implementation of "experiments" and "jobs/runs" in KFP [#4955](https://github.com/kubeflow/kubeflow/issues/4955)

### Neural Architecture Search

- Refactor structure for NAS algorithms [#1125](https://github.com/kubeflow/katib/issues/1125)
- Refactor the design for NAS model constructor [#1127](https://github.com/kubeflow/katib/issues/1127)
- ENAS enhancements such as micro mode, RNN support
