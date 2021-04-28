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

# History

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
