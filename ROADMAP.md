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
- Support new Kubeflow operators (MXNet, XGBoost, MPI)
- Add validation for algorithms (a.k.a suggestions) [#1126](https://github.com/kubeflow/katib/issues/1126)
- Katib UI fixes and enhancements
- Investigate Kubeflow Metadata integration

### Neural Architecture Search

- Refactor structure for NAS algorithms [#1125](https://github.com/kubeflow/katib/issues/1125)
- Refactor the design for NAS model constructor [#1127](https://github.com/kubeflow/katib/issues/1127)
- ENAS enhancements such as micro mode, RNN support
