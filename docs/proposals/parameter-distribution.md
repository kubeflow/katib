# Proposal for Supporting various parameter distributions in Katib

## Summary
The goal of this project is to enhance the existing Katib Experiment APIs to support various parameter distributions such as uniform, log-uniform, and qlog-uniform. Then extend the suggestion services to be able to configure distributions for search space using libraries provided in each framework.

## Motivation
Currently, [Katib](https://github.com/kubeflow/katib) is limited to supporting only uniform distribution for integer, float, and categorical hyperparameters. By introducing additional distributions, Katib will become more flexible and powerful in conducting hyperparameter optimization tasks.

A Data Scientist requires Katib to support multiple hyperparameter distributions, such as log-uniform, normal, and log-normal, in addition to the existing uniform distribution. This enhancement is crucial for more flexible and precise hyperparameter optimization. For instance, learning rates often benefit from a log-uniform distribution because small values can significantly impact performance. Similarly, normal distributions are useful for parameters that are expected to vary around a central value.

### Goals
- Add `Distribution` field to `ParameterSpec` alongside `ParameterType`.
- Support for the log-uniform, normal, and log-normal Distributions.
- Update the Experiment and gRPC API to support `Distribution`.
- Update logic to handle the new parameter distributions for each suggestion service (e.g., Optuna, Hyperopt).
- Extend the Python SDK to support the new `Distribution` field.
### Non-Goals
- This proposal do not aim to create new version for CRD APIs.
- No changes will be made to the core optimization algorithms beyond supporting new distributions.

## Proposal

### Parameter Distribution Comparison Table
| Distribution Type             | Hyperopt              | Optuna                               | Ray Tune              |
|-------------------------------|-----------------------|--------------------------------------|-----------------------|
| **Uniform Continuous**        | `hp.uniform`          | `FloatDistribution`                  | `tune.uniform`        |
| **Quantized Uniform**         | `hp.quniform`         | `DiscreteUniformDistribution`(deprecated)  Use `FloatDistribution` instead. | `tune.quniform`       |
| **Log Uniform**               | `hp.loguniform`       | `LogUniformDistribution`(deprecated)  Use `FloatDistribution` instead. | `tune.loguniform`     |
| **Uniform Integer**           | `hp.randint` or quantized distributions with step size `q` set to 1 | `IntDistribution`                    | `tune.randint`        |
| **Categorical**               | `hp.choice`           | `CategoricalDistribution`            | `tune.choice`         |
| **Quantized Log Uniform**     | `hp.qloguniform`      | Custom Implementation Use `FloatDistribution` instead.                 | `tune.qloguniform`    |
| **Normal**                    | `hp.normal`           | (Not directly supported)             | `tune.randn`          |
| **Quantized Normal**          | `hp.qnormal`          | (Not directly supported)             | `tune.qrandn`         |
| **Log Normal**                | `hp.lognormal`        | (Not directly supported)             | (Use custom transformation in `tune.randn`) |
| **Quantized Log Normal**      | `hp.qlognormal`       | (Not directly supported)             | (Use custom transformation in `tune.qrandn`) |
| **Quantized Integer**         | `hp.quniformint`      | `IntUniformDistribution`(deprecated)  Use `IntDistribution` instead.          |                       |
| **Log Integer**               |                       | `IntLogUniformDistribution`(deprecated)  Use `IntDistribution` instead.          | `tune.lograndint`     |

### How is Nevergrad implementing Hyperopt?
Nevergrad maps parameter types (like p.Scalar, p.Log, p.Choice, etc.) from Nevergrad to corresponding Hyperopt search space definitions (hp.uniform, hp.loguniform, hp.choice, etc.).
```python
def _get_search_space(param_name, param):
    if isinstance(param, p.Scalar):
        ...
        return hp.uniform(param_name, param.bounds[0][0], param.bounds[1][0])
    elif isinstance(param, p.Log):
        ...
        return hp.loguniform(param_name, np.log(param.bounds[0][0]), np.log(param.bounds[1][0]))
    elif isinstance(param, p.Choice):
        ...
        return hp.choice()
```
The `_get_search_space` function constructs a search space that represents the entire parameter space defined by Nevergrad. 

## Experiment API changes
Scope: `pkg/apis/controller/experiments/v1beta1/experiment_types.go`
- Adding new field `Distribution` to `ParameterSpec`

```go
type ParameterSpec struct {
	Name          string        `json:"name,omitempty"`
	ParameterType ParameterType `json:"parameterType,omitempty"`
    Distribution  Distribution  `json:"distribution,omitempty"` // Added Distribution field
	FeasibleSpace FeasibleSpace `json:"feasibleSpace,omitempty"`
}
```

- Defining `Distribution` type
```go
type Distribution string

const (
 	CategoricalDistribution     Distribution = "categorical"
 	UniformDistribution         Distribution = "uniform"
 	LogUniformDistribution      Distribution = "logUniform"
    NormalDistribution          Distribution = "normal"
    LogNormalDistribution       Distribution = "logNormal"
)
```
- The `Step` field can be used to define quantization steps for uniform or log-uniform distributions, effectively covering q-quantization requirements.
```go
type FeasibleSpace struct {
	Max  string   `json:"max,omitempty"`
	Min  string   `json:"min,omitempty"`
	List []string `json:"list,omitempty"`
	Step string   `json:"step,omitempty"` // Step can be used to define q-quantization
}
```

## gRPC API changes
Scope: `pkg/apis/manager/v1beta1/api.proto`
- Add the `Distribution` field to the `ParameterSpec` message
```
/**
 * Config for a hyperparameter.
 * Katib will create each Hyper parameter from this config.
 */
message ParameterSpec {
    string name = 1; /// Name of the parameter.
    ParameterType parameter_type = 2; /// Type of the parameter.
    FeasibleSpace feasible_space = 3; /// FeasibleSpace for the parameter.
    Distribution distribution = 4; // Distribution of the parameter.
}
```
- Define the `Distribution` enum
```
/**
 * Distribution types for HyperParameter.
 */
enum Distribution {
    CATEGORICAL = 0;
    UNIFORM = 1;
    LOG_UNIFORM = 2;
    NORMAL = 3;
    LOG_NORMAL = 4;
}
```

## Suggestion Service Logic
- For each suggestion service (e.g., Optuna, Hyperopt), the logic will be updated to handle the new parameter distributions.
- This involves modifying the conversion functions to map Katib distributions to the corresponding framework-specific distributions.

#### Optuna
ref: https://optuna.readthedocs.io/en/stable/reference/distributions.html

For example:
- Update the `_get_optuna_search_space` for new parameters.
scope: `pkg/suggestion/v1beta1/optuna/base_service.py`
```python
    def _get_optuna_search_space(self):
        search_space = {}
        for param in self.search_space.params:
            if param.type == INTEGER:
                search_space[param.name] = optuna.distributions.IntDistribution(int(param.min), int(param.max))
            elif param.type == DOUBLE:
                search_space[param.name] = optuna.distributions.FloatDistribution(float(param.min), float(param.max))
            elif param.type == CATEGORICAL or param.type == DISCRETE:
                search_space[param.name] = optuna.distributions.CategoricalDistribution(param.list)
        return search_space
```

#### Goptuna
ref: https://github.com/c-bata/goptuna/blob/2245ddd9e8d1edba750839893c8a618f852bc1cf/distribution.go

#### Hyperopt
ref: http://hyperopt.github.io/hyperopt/getting-started/search_spaces/#parameter-expressions

#### Ray-tune
ref: https://docs.ray.io/en/latest/tune/api/search_space.html

## Python SDK
Extend the Python SDK to support the new `Distribution` field.

