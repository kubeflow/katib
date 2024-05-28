# Proposal for Supporting various parameter distributions in Katib

## Summary
The goal of this project is to enhance the existing Katib Experiment APIs to support various parameter distributions such as uniform, log-uniform, and qlog-uniform. Then extend the suggestion services to be able to configure distributions for search space using libraries provided in each framework.

## Motivation
Currently, [Katib](https://github.com/kubeflow/katib) is limited to supporting only uniform distribution for integer, float, and categorical hyperparameters. By introducing additional distributions, Katib will become more flexible and powerful in conducting hyperparameter optimization tasks.

A Data Scientist requires Katib to support multiple hyperparameter distributions, such as log-uniform, normal, and log-normal, in addition to the existing uniform distribution. This enhancement is crucial for more flexible and precise hyperparameter optimization. For instance, learning rates often benefit from a log-uniform distribution because small values can significantly impact performance. Similarly, normal distributions are useful for parameters that are expected to vary around a central value.

### Goals
- Add `Distribution` field to `FeasibleSpace` alongside `ParameterType`.
- Support for the log-uniform, normal, and log-normal Distributions.
- Update the Experiment and gRPC API to support `Distribution`.
- Update logic to handle the new parameter distributions for each suggestion service (e.g., Optuna, Hyperopt).
- Extend the Python SDK to support the new `Distribution` field.
### Non-Goals
- This proposal do not aim to create new version for CRD APIs.
- This proposal do not aim to make the necessary Katib UI changes.
- No changes will be made to the core optimization algorithms beyond supporting new distributions.

## Proposal

### Parameter Distribution Comparison Table

| Distribution Type             | Hyperopt              | Optuna                                          | Ray Tune              | Nevergrad                                    |
|-------------------------------|-----------------------|-------------------------------------------------|-----------------------|---------------------------------------------|
| **Uniform Continuous**        | `hp.uniform`          | `FloatDistribution`                             | `tune.uniform`        | `p.Scalar` with uniform transformation      |
| **Quantized Uniform**         | `hp.quniform`         | `DiscreteUniformDistribution` (deprecated)      | `tune.quniform`       | `p.Scalar` with uniform and step specified  |
| **Log Uniform**               | `hp.loguniform`       | `LogUniformDistribution` (deprecated)           | `tune.loguniform`     | `p.Log` with uniform transformation         |
| **Uniform Integer**           | `hp.randint` or quantized distributions with step size `q` set to 1 | `IntDistribution`                    | `tune.randint`        | `p.Scalar` with integer transformation     |
| **Categorical**               | `hp.choice`           | `CategoricalDistribution`                       | `tune.choice`         | `p.Choice`                                  |
| **Quantized Log Uniform**     | `hp.qloguniform`      | Custom Implementation                           | `tune.qloguniform`    | `p.Log` with uniform and step specified    |
| **Normal**                    | `hp.normal`           | (Not directly supported)                        | `tune.randn`          | (Not directly supported)                    |
| **Quantized Normal**          | `hp.qnormal`          | (Not directly supported)                        | `tune.qrandn`         | (Not directly supported)                    |
| **Log Normal**                | `hp.lognormal`        | (Not directly supported)                        | (Use custom transformation in `tune.randn`) | (Not directly supported)                    |
| **Quantized Log Normal**      | `hp.qlognormal`       | (Not directly supported)                        | (Use custom transformation in `tune.qrandn`) | (Not directly supported)                    |
| **Quantized Integer**         | `hp.quniformint`      | `IntUniformDistribution` (deprecated)           |                       | `p.Scalar` with integer and step specified  |
| **Log Integer**               |                       | `IntLogUniformDistribution` (deprecated)        | `tune.lograndint`     | `p.Scalar` with log-integer transformation |


- Note:
In `Nevergrad`, parameter types like `p.Scalar`, `p.Log`, and `p.Choice` are mapped to corresponding `Hyperopt` search space definitions like `hp.uniform`, `hp.loguniform`, and `hp.choice` using internal functions to convert parameter bounds and distributions.

## API Design
### FeasibleSpace
Feasible space for optimization.
Int and Double type use Max/Min.
Discrete and Categorical type use List.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| max | [string](#string) |  | Max Value |
| min | [string](#string) |  | Minimum Value |
| list | [string](#string) | repeated | List of Values. |
| step | [string](#string) |  | Step for double or int parameter or q for quantization|
| distribution | [Distribution](#api-v1-beta1-Distribution) |  | Type of the Distribution. |


<a name="api-v1-beta1-Distribution"></a>

### Distribution
- Types of value for HyperParameter Distributions.
- We add the `distribution` field to represent the hyperparameters search space rather than [`ParameterType`](https://github.com/kubeflow/katib/blob/2c575227586ff1c03cf6b5190d066e2f3061a404/pkg/apis/controller/experiments/v1beta1/experiment_types.go#L199-L207).
- The `distribution` allows users to configure more granular search space customizations.
- In this enhancement, we would propose the following 4 distributions:

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNIFORM | 0 | Continuous uniform distribution. Samples values evenly between a minimum and maximum value. Use &#34;Max/Min&#34;. Use &#34;Step&#34; for `q`. |
| LOGUNIFORM | 1 | Samples values such that their logarithm is uniformly distributed. Use &#34;Max/Min&#34;. Use &#34;Step&#34; for `q`. |
| NORMAL | 2 | Normal (Gaussian) distribution type. Samples values according to a normal distribution characterized by a mean and standard deviation. Use &#34;Max/Min&#34;. Use &#34;Step&#34; for `q`. |
| LOGNORMAL | 3 | Log-normal distribution type. Samples values such that their logarithm is normally distributed. Use &#34;Max/Min&#34;. Use &#34;Step&#34; for `q`. |


## Experiment API changes
Scope: `pkg/apis/controller/experiments/v1beta1/experiment_types.go`

```go
type ParameterSpec struct {
	Name          string        `json:"name,omitempty"`
	ParameterType ParameterType `json:"parameterType,omitempty"`
	FeasibleSpace FeasibleSpace `json:"feasibleSpace,omitempty"`
}
```
- Adding new field `Distribution` to `FeasibleSpace`

- The `Step` field can be used to define quantization steps for uniform or log-uniform distributions, effectively covering q-quantization requirements.

Updated `FeasibleSpace` struct
```diff
type FeasibleSpace struct {
	Max           string        `json:"max,omitempty"`
	Min           string        `json:"min,omitempty"`
	List          []string      `json:"list,omitempty"`
	Step          string        `json:"step,omitempty"` // Step can be used to define q-quantization
+       Distribution  Distribution  `json:"distribution,omitempty"` // Added Distribution field
}
```
 - New Field Description: `Distribution`
  - Type: `Distribution`
  - Description: The Distribution field specifies the type of statistical distribution to be applied to the parameter. This allows the definition of various distributions, such as uniform, log-uniform, or other supported types.

- Defining `Distribution` type
```go
type Distribution string

const (
	DistributionUniform    Distribution = "uniform"
	DistributionLogUniform Distribution = "logUniform"
	DistributionNormal     Distribution = "normal"
	DistributionLogNormal  Distribution = "logNormal"
)
```

## gRPC API changes
Scope: `pkg/apis/manager/v1beta1/api.proto`
- Add the `Distribution` field to the `FeasibleSpace` message
```diff
/**
 * Feasible space for optimization.
 * Int and Double type use Max/Min.
 * Discrete and Categorical type use List.
 */
message FeasibleSpace {
    string max = 1; /// Max Value
    string min = 2; /// Minimum Value
    repeated string list = 3; /// List of Values.
    string step = 4; /// Step for double or int parameter
+   Distribution distribution = 4; // Distribution of the parameter.
}
```
- Define the `Distribution` enum
```
/**
 * Distribution types for HyperParameter.
 */
enum Distribution {
    UNIFORM = 0;
    LOG_UNIFORM = 1;
    NORMAL = 2;
    LOG_NORMAL = 3;
}
```

## Suggestion Service Logic
- For each suggestion service (e.g., Optuna, Hyperopt), the logic will be updated to handle the new parameter distributions.
- This involves modifying the conversion functions to map Katib distributions to the corresponding framework-specific distributions.

#### Optuna
ref: https://optuna.readthedocs.io/en/stable/reference/distributions.html

For example:
- Update the `_get_optuna_search_space` for new Distributions.
scope: `pkg/suggestion/v1beta1/optuna/base_service.py`

#### Goptuna
ref: https://github.com/c-bata/goptuna/blob/2245ddd9e8d1edba750839893c8a618f852bc1cf/distribution.go

#### Hyperopt
ref: http://hyperopt.github.io/hyperopt/getting-started/search_spaces/#parameter-expressions

#### Ray-tune
ref: https://docs.ray.io/en/latest/tune/api/search_space.html

## Python SDK
Extend the Python SDK to support the new `Distribution` field.

