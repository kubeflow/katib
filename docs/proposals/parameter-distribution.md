# Proposal for Supporting various parameter distributions in Katib

## Summary
The goal of this project is to enhance the existing Katib Experiment APIs to support various parameter distributions such as uniform, log-uniform, and qlog-uniform. Then extend the suggestion services to be able to configure distributions for search space using libraries provided in each framework.

## Motivation
Currently, Katib is limited to supporting only uniform distribution for integer, float, and categorical hyperparameters. By introducing additional distributions, Katib will become more flexible and powerful in conducting hyperparameter optimization tasks.

## Proposal

### Maintaining two versions of CRD APIs
- Introduce `v1beta2` API version
- The (conversion webhook)[https://book.kubebuilder.io/multiversion-tutorial/conversion.html] will serve as a critical component in managing transitions between different API versions (`v1beta1` to `v1beta2`).
- TODO 
> I will create a correspondence table between v1beta1 and v1beta2. Maybe, we only need to create a table for the `ParameterType` and the `FeasibleSpace`.
Look into this.
- A specific field will be added to the katib-config `katib/pkg/util/v1beta1/katibconfig/config.go` to handle different gRPC client versions required by the suggestion controller. (One approach could be this for maintaining two versions of gRPC APIs).
- TODO Based on this comment here:
https://github.com/kubeflow/katib/pull/2059#discussion_r1049329229
Look how to maintain two versions of gRPC APIs

### Experiment API changes
Scope: `pkg/apis/controller/experiments/v1beta1/experiment_types.go`
- Modifying the `ParameterSpec` struct to include a new field called Distribution.

```go
type ParameterSpec struct {
    Name          string       `json:"name,omitempty"`
    Distribution  Distribution `json:"distribution,omitempty"`
    FeasibleSpace FeasibleSpace `json:"feasibleSpace,omitempty"`
}
```

- Renaming Parameters

```go
type Distribution string

const (
 	CategoricalDistribution     Distribution = "categorical"
 	IntUniformDistribution      Distribution = "intUniform"
 	IntLogUniformDistribution   Distribution = "intLogUniform"
 	FloatUniformDistribution    Distribution = "floatUniform"
)
```

### Suggestion Service Logic
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

#### Chocolate
ref: https://chocolate.readthedocs.io/en/latest/api/space.html

#### Goptuna
ref: https://github.com/c-bata/goptuna/blob/2245ddd9e8d1edba750839893c8a618f852bc1cf/distribution.go

#### Hyperopt
ref: http://hyperopt.github.io/hyperopt/getting-started/search_spaces/#parameter-expressions

#### Scikit Optimize
ref: https://scikit-optimize.github.io/stable/modules/classes.html#module-skopt.space.space

### Testing
- Write unit tests for the conversion webhook to ensure all fields are correctly converted.
- Write unit tests for the new parameter distributions to ensure they are processed correctly by the controller.

