export type ParameterType = 'int' | 'double' | 'discrete' | 'categorical';

export interface Range {
  min: string;
  max: string;
  step: string;
}

export type FeasibleSpace = Range | any[];

export interface ParameterSpec {
  name: string;
  type: ParameterType;
  value: FeasibleSpace;
}
