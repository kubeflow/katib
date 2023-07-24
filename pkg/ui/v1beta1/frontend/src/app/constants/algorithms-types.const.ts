import {
  AlgorithmsEnum,
  EarlyStoppingAlgorithmsEnum,
} from '../enumerations/algorithms.enum';

export const AlgorithmNames = {
  [AlgorithmsEnum.GRID]: 'Grid',
  [AlgorithmsEnum.RANDOM]: 'Random',
  [AlgorithmsEnum.HYPERBAND]: 'Hyperband',
  [AlgorithmsEnum.BAYESIAN_OPTIMIZATION]: 'Bayesian Optimization',
  [AlgorithmsEnum.TPE]: 'Tree of Parzen Estimators',
  [AlgorithmsEnum.MULTIVARIATE_TPE]: 'Multivariate Tree of Parzen Estimators',
  [AlgorithmsEnum.CMAES]: 'Covariance Matrix Adaptation: Evolution Strategy',
  [AlgorithmsEnum.SOBOL]: 'Sobol Quasirandom Sequence',
  [AlgorithmsEnum.PBT]: 'Population Based Training',
};

export const NasAlgorithmNames = {
  [AlgorithmsEnum.ENAS]: 'Efficient Neural Architecture Search',
  [AlgorithmsEnum.DARTS]: 'Differentiable Architecture Search',
};

export const EarlyStoppingAlgorithmNames = {
  [EarlyStoppingAlgorithmsEnum.NONE]: 'None',
  [EarlyStoppingAlgorithmsEnum.MEDIAN]: 'Median Stopping Rule',
};
