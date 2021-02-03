import { ParameterSpec } from 'src/app/shared/params-list/types';

export interface NasOperation {
  type: string;
  params: ParameterSpec[];
}
