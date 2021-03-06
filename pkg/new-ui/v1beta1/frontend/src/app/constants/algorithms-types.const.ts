import { AlgorithmsEnum } from '../enumerations/algorithms.enum';

export const AlgorithmNames = {
  [AlgorithmsEnum.GRID]: 'Grid',
  [AlgorithmsEnum.RANDOM]: 'Random',
  [AlgorithmsEnum.HYPERBAND]: 'Hyperband',
  [AlgorithmsEnum.BAYESIAN_OPTIMIZATION]: 'Bayesian Optimization',
  [AlgorithmsEnum.TPE]: 'Tree of Parzen Estimators',
  [AlgorithmsEnum.CMAES]: 'Covariance Matrix Adaptation: Evolution Strategy',
};

export const NasAlgorithmNames = {
  [AlgorithmsEnum.ENAS]: 'Efficient Neural Architecture Search',
  [AlgorithmsEnum.DARTS]: 'Differentiable Architecture Search',
};
