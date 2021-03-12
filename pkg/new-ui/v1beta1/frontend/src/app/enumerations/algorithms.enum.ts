export enum AlgorithmsEnum {
  GRID = 'grid',
  RANDOM = 'random',
  HYPERBAND = 'hyperband',
  BAYESIAN_OPTIMIZATION = 'bayesianoptimization',
  TPE = 'tpe',
  CMAES = 'cmaes',
  ENAS = 'enas',
  DARTS = 'darts',
}

export enum EarlyStoppingAlgorithmsEnum {
  NONE = 'none',
  MEDIAN = 'medianstop',
}

export enum AlgorithmSettingType {
  STRING = 'string',
  INTEGER = 'integer',
  FLOAT = 'float',
}
