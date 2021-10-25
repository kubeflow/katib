import lowerCase from 'lodash-es/lowerCase';
import {
  FormControl,
  FormGroup,
  Validators,
  FormArray,
  ValidatorFn,
  AbstractControl,
  ValidationErrors,
} from '@angular/forms';
import {
  ParameterSpec,
  FeasibleSpaceMinMax,
  FeasibleSpaceList,
  NasOperation,
  FeasibleSpace,
  ParameterType,
} from '../models/experiment.k8s.model';

export function createNasOperationGroup(op: NasOperation): FormGroup {
  const array = op.parameters.map(param => createParameterGroup(param));

  return new FormGroup({
    operationType: new FormControl(op.operationType, Validators.required),
    parameters: new FormArray(array),
  });
}

export function createFeasibleSpaceGroup(
  parameterType: ParameterType,
  feasibleSpace: FeasibleSpace,
) {
  let fs: FeasibleSpace;

  // min-max-step
  if (parameterType === 'int' || parameterType === 'double') {
    fs = feasibleSpace as FeasibleSpaceMinMax;

    return new FormGroup({
      min: new FormControl(fs.min, Validators.required),
      max: new FormControl(fs.max, Validators.required),
      step: new FormControl(fs.step, checkIfZero()),
    });
  }

  // list values
  fs = feasibleSpace as FeasibleSpaceList;

  const ctrls = fs.list.map(v => new FormControl(v, Validators.required));
  return new FormGroup({
    list: new FormArray(ctrls, Validators.required),
  });
}

export function createParameterGroup(param: ParameterSpec): FormGroup {
  return new FormGroup({
    name: new FormControl(param.name, Validators.required),
    parameterType: new FormControl(param.parameterType, Validators.required),
    feasibleSpace: createFeasibleSpaceGroup(
      param.parameterType,
      param.feasibleSpace,
    ),
  });
}

/*
 * Arithmetics
 **/
export const numberToExponential = (num: number, digits: number): string => {
  if (isNaN(Number(num))) {
    return '';
  }

  if (num.toString().replace(/[.]/g, '').length <= digits) {
    return num.toString();
  }

  if (
    num.toExponential().search(/e[+]1/) > -1 ||
    num.toExponential().search(/e[-]1/) > -1
  ) {
    const slicedNumber = num.toString().slice(0, digits + 1);

    return slicedNumber.replace(/[.]*0+$/, '');
  }

  let exponentialNumber = num.toExponential(digits - 1);

  // If toExponential added e+0 in the end of the string remove it
  exponentialNumber = exponentialNumber.replace(/[.]*0*e[+]0$/, '');

  // If the number is e.g. 2.1000e-3, the zeros must to be removed
  if (/[.]*0+e[+-][1-9]$/.test(exponentialNumber)) {
    // Split the number and the exponent in order to remove from
    // the number the zeros
    const [numberToFix, exponentNumber] = exponentialNumber.split('e');
    const fixed = numberToFix.replace(/[0]*$/g, '');

    // Build again the exponential number
    exponentialNumber = `${fixed}e${exponentNumber}`;

    /*If the number was e.g. 2.000e-3 after the above replacement
    it would be 2.e-3 so a zero between . and e has to be added*/
    return exponentialNumber.replace(/[.]e/g, '.0e');
  }

  return exponentialNumber;
};

export const transformStringResponses = (
  response: string,
): { types: string[]; details: string[][] } => {
  response = response.replace(/\n$/, '');

  // Separate each line
  const lines = response.split('\n');
  let types = [];
  let details = [];

  // The first line is the names of the types
  types = lines[0].split(',');
  // Transform them for consistency
  types = types.map(column => lowerCase(column));

  // Separate types from details of every item
  lines.splice(0, 1);

  // Change the first letter of the column to upper case
  types = types.map(column => column.charAt(0).toUpperCase() + column.slice(1));

  // Transform the details of each item from string to an array
  details = lines.map(detail => detail.split(','));

  return { types, details };
};

export const safeDivision = (divided: number, divider: number): number =>
  Math.round((divided * 10000.0) / divider) / 10000;

export const safeMultiplication = (
  multiplicand: number,
  multiplier: number,
): number => Math.round(multiplicand * 10000.0 * multiplier) / 10000;

export const checkIfZero = (): ValidatorFn => {
  return (control: AbstractControl): ValidationErrors | null => {
    if (control.value === null || control.value === '') return null;

    const isZero = !/[^0.]/g.test(control.value.toString());
    return isZero ? { mustNotBeZero: { value: control.value } } : null;
  };
};
