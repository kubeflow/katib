# Proposal for Parameter Distribution

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
- [Design Details](#design-details)
  - [Experiment API changes](#experiment-api-changes)
  - [Correspondence for Katib Distributions and Framework Distributions](#correspondence-for-katib-distributions-and-framework-distributions)
    - [Chocolate](#chocolate)
    - [Goptuna](#goptuna)
    - [Hyperopt](#hyperopt)
    - [Optuna](#optuna)
    - [Scikit-Optimize](#scikit-optimize)
<!-- /toc -->

## Summary
This enhancement introduces `Distribution` to tuning parameters and remove redundantly `ParameterType`. 

API field in the Experiment spec determine parameter type with distribution.  

## Motivation
Currently, Katib does not support determining a distribution for search space that samplers pick up parameters by users.

Katib should be able to determine it by users since
almost hyperparameter tuning algorithms (framework) can determine it by users.

## Proposal
We introduce a mechanism to determine a distribution for search space by users.
That also means we introduce a mechanism to propagate distributions to suggestion-services and
set them samplers of each framework.     

## Design Details
The proposal consists of a new Experiment API field and
correspondence between Katib Distributions and Framework Distributions.

### Experiment API changes
We extend the Experiment API to introduce the new fields `Distribution` to configure the distribution for search space and
remove the redundant fields `ParameterType`.

```diff
type ParameterSpec struct {
	Name          string        `json:"name,omitempty"`
- 	ParameterType ParameterType `json:"parameterType,omitempty"`
+ 	Distribution  Distribution  `json:"distribution,omitempty"`
	FeasibleSpace FeasibleSpace `json:"feasibleSpace,omitempty"`
}

- type ParameterType string
+ type Distribution string

const (
- 	ParameterTypeUnknown     ParameterType = "unknown"
- 	ParameterTypeDouble      ParameterType = "double"
- 	ParameterTypeInt         ParameterType = "int"
- 	ParameterTypeDiscrete    ParameterType = "discrete"
- 	ParameterTypeCategorical ParameterType = "categorical"
+ 	UnknownDistribution         Distribution = "unknown"
+ 	CategoricalDistribution     Distribution = "categorical"
+ 	IntUniformDistribution      Distribution = "intUniform"
+ 	IntLogUniformDistribution   Distribution = "intLogUniform"
+ 	FloatUniformDistribution    Distribution = "floatUniform"
+ 	FloatLogUniformDistribution Distribution = "floatLogUniform"
)
```

### Correspondence between Katib Distributions and Framework Distributions
We extend suggestion services to be able to configure distributions for
search space using libraries provided in each framework.

#### Chocolate
TODO

#### Goptuna
We can extend Goptuna Suggestion Service using Goptuna libraries shown in the below correspondence table for
Katib Distributions and Goptuna Distributions.

ref: https://github.com/c-bata/goptuna/blob/2245ddd9e8d1edba750839893c8a618f852bc1cf/distribution.go

| Katib Distribution          | Goptuna Distribution                                 |
|-----------------------------|------------------------------------------------------|
| CategoricalDistribution     | CategoricalDistribution                              |
| IntUniformDistribution      | IntUniformDistribution or StepIntUniformDistribution |
| IntLogUniformDistribution   | IntUniformDistribution or StepIntUniformDistribution |
| FloatUniformDistribution    | UniformDistribution                                  |
| FloatLogUniformDistribution | LogUniformDistribution                               |


#### Hyperopt
We can extend Hyperopt Suggestion Service using Hyperopt libraries shown in the below correspondence table for
Katib Distributions and Hyperopt Distributions.

ref: http://hyperopt.github.io/hyperopt/getting-started/search_spaces/#parameter-expressions

| Katib Distribution          | Hyperopt Distribution |
|-----------------------------|-----------------------|
| CategoricalDistribution     | hp.choice             |
| IntUniformDistribution      | hp.quniform           |
| IntLogUniformDistribution   | hp.qloguniform        |
| FloatUniformDistribution    | hp.quniform           |
| FloatLogUniformDistribution | hp.qloguniform        |

#### Optuna
We can extend Optuna Suggestion Service using Optuna libraries shown in the below correspondence table for
Katib Distributions and Optuna Distributions.

ref: https://optuna.readthedocs.io/en/stable/reference/distributions.html

| Katib Distribution          | Optuna Distribution                   |
|-----------------------------|---------------------------------------|
| CategoricalDistribution     | distributions.CategoricalDistribution |
| IntUniformDistribution      | distributions.IntDistribution         |
| IntLogUniformDistribution   | distributions.IntDistribution         |
| FloatUniformDistribution    | distributions.FloatDistribution       |
| FloatLogUniformDistribution | distributions.FloatDistribution       |

#### Scikit Optimize
We can extend Scikit-Optimize Suggestion Service using Scikit-Optimize libraries shown in the below correspondence table for
Katib Distributions and Scikit-Optimize Distributions.

ref: https://scikit-optimize.github.io/stable/modules/classes.html#module-skopt.space.space

| Katib Distribution          | Scikit-Optimize Distribution |
|-----------------------------|------------------------------|
| CategoricalDistribution     | space.Categorical            |
| IntUniformDistribution      | space.Integer                |
| IntLogUniformDistribution   | space.Integer                |
| FloatUniformDistribution    | space.Real                   |
| FloatLogUniformDistribution | space.Real                   |
