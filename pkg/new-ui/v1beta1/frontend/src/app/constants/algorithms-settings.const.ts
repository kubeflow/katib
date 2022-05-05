import {
  AlgorithmsEnum,
  AlgorithmSettingType,
  EarlyStoppingAlgorithmsEnum,
} from '../enumerations/algorithms.enum';

export interface AlgorithmSetting {
  name: string;
  value: any;
  values?: any[];
  type: AlgorithmSettingType;
}

export const GridSettings: AlgorithmSetting[] = [];

export const RandomSearchSettings: AlgorithmSetting[] = [
  {
    name: 'random_state',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const BayesianOptimizationSettings: AlgorithmSetting[] = [
  {
    name: 'base_estimator',
    value: 'GP',
    values: ['GP', 'RF', 'ET', 'GBRT'],
    type: AlgorithmSettingType.STRING,
  },
  {
    name: 'n_initial_points',
    value: 10,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'acq_func',
    value: 'gp_hedge',
    values: ['gp_hedge', 'LCB', 'EI', 'PI', 'EIps', 'PIps'],
    type: AlgorithmSettingType.STRING,
  },
  {
    name: 'acq_optimizer',
    value: 'auto',
    values: ['auto', 'sampling', 'lbfgs'],
    type: AlgorithmSettingType.STRING,
  },
  {
    name: 'random_state',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const HyperbandSettings: AlgorithmSetting[] = [];

export const TPESettings: AlgorithmSetting[] = [
  {
    name: 'gamma',
    value: null,
    type: AlgorithmSettingType.STRING,
  },
  {
    name: 'prior_weight',
    value: null,
    type: AlgorithmSettingType.STRING,
  },
  {
    name: 'n_EI_candidates',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'random_state',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const MultivariateTPESettings: AlgorithmSetting[] = [
  {
    name: 'n_startup_trials',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'n_ei_candidates',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'random_state',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const CMAESSettings: AlgorithmSetting[] = [
  {
    name: 'random_state',
    value: null,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'sigma',
    value: 0.001,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'restart_strategy',
    value: 'none',
    values: ['none', 'ipop', 'bipop'],
    type: AlgorithmSettingType.STRING,
  },
];

export const SOBOLSettings: AlgorithmSetting[] = [];

export const ENASSettings: AlgorithmSetting[] = [
  {
    name: 'controller_hidden_size',
    value: 64,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'controller_temperature',
    value: 5.0,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_tanh_const',
    value: 2.25,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_entropy_weight',
    value: 1e-5,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_baseline_decay',
    value: 0.999,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_learning_rate',
    value: 5e-5,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_skip_target',
    value: 0.4,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_skip_weight',
    value: 0.8,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'controller_train_steps',
    value: 50,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'controller_log_every_steps',
    value: 10,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const DartsSettings: AlgorithmSetting[] = [
  {
    name: 'num_epochs',
    value: 50,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'w_lr',
    value: 0.025,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'w_lr_min',
    value: 0.001,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'w_momentum',
    value: 0.9,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'w_weight_decay',
    value: 3e-4,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'w_grad_clip',
    value: 5.0,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'alpha_lr',
    value: 3e-4,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'alpha_weight_decay',
    value: 1e-3,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'batch_size',
    value: 128,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'num_workers',
    value: 4,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'init_channels',
    value: 16,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'print_step',
    value: 50,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'num_nodes',
    value: 4,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'stem_multiplier',
    value: 3,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const PbtSettings: AlgorithmSetting[] = [
  {
    name: 'suggestion_trial_dir',
    value: '/var/log/katib/checkpoints/',
    type: AlgorithmSettingType.STRING,
  },
  {
    name: 'n_population',
    value: 40,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'resample_probability',
    value: null,
    type: AlgorithmSettingType.FLOAT,
  },
  {
    name: 'truncation_threshold',
    value: 0.2,
    type: AlgorithmSettingType.FLOAT,
  },
];

export const EarlyStoppingSettings: AlgorithmSetting[] = [
  {
    name: 'min_trials_required',
    value: 3,
    type: AlgorithmSettingType.INTEGER,
  },
  {
    name: 'start_step',
    value: 4,
    type: AlgorithmSettingType.INTEGER,
  },
];

export const EarlyStoppingSettingsMap: { [key: string]: AlgorithmSetting[] } = {
  [EarlyStoppingAlgorithmsEnum.NONE]: [],
  [EarlyStoppingAlgorithmsEnum.MEDIAN]: EarlyStoppingSettings,
};

export const AlgorithmSettingsMap: { [key: string]: AlgorithmSetting[] } = {
  [AlgorithmsEnum.GRID]: GridSettings,
  [AlgorithmsEnum.RANDOM]: RandomSearchSettings,
  [AlgorithmsEnum.HYPERBAND]: HyperbandSettings,
  [AlgorithmsEnum.BAYESIAN_OPTIMIZATION]: BayesianOptimizationSettings,
  [AlgorithmsEnum.TPE]: TPESettings,
  [AlgorithmsEnum.MULTIVARIATE_TPE]: MultivariateTPESettings,
  [AlgorithmsEnum.CMAES]: CMAESSettings,
  [AlgorithmsEnum.SOBOL]: SOBOLSettings,
  [AlgorithmsEnum.ENAS]: ENASSettings,
  [AlgorithmsEnum.DARTS]: DartsSettings,
  [AlgorithmsEnum.PBT]: PbtSettings,
};
