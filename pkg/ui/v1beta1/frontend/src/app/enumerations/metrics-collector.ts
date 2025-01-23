export enum CollectorKind {
  STDOUT = 'StdOut',
  FILE = 'File',
  TFEVENT = 'TensorFlowEvent',
  PROMETHEUS = 'PrometheusMetric',
  CUSTOM = 'Custom',
  PUSH = 'Push',
}
