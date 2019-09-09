# Suggestions

This document describes the usage of suggestion algorithms implemented or integrated in katib.

## Grid Search

Grid sampling applies when all variables are discrete (Doubles and integers need to be quantized) and the number of possibilities is low. A grid search will perform the exhaustive combinatorial search over all possibilities making the search extremely long even for medium sized problems.

### [Chocolate][]

> Chocolate is a completely asynchronous optimisation framework relying solely on a database to share information between workers. Chocolate uses no master process for distributing tasks. Every task is completely independent and only gets its information from the database. Chocolate is thus ideal in controlled computing environments where it is hard to maintain a master process for the duration of the optimisation.

Algorithm name in katib is `chocolate-grid`.

## Random Search

Random sampling is an alternative to grid search when the number of discrete parameters to optimize and the time required for each evaluation is high. When all parameters are discrete, random search will perform sampling without replacement making it an algorithm of choice when combinatorial exploration is not possible. With continuous parameters, it is preferable to use quasi random sampling.

### [Chocolate][]

Algorithm name in katib is `chocolate-random`.

## Quasi Random Search

QuasiRandom sampling ensures a much more uniform exploration of the search space than traditional pseudo random. Thus, quasi random sampling is preferable when not all variables are discrete, the number of dimensions is high and the time required to evaluate a solution is high.

### [Chocolate][]

Algorithm name in katib is `chocolate-quasirandom`.

## CMAES

CMAES search is one of the most powerful black-box optimization algorithm. However, it requires a significant number of model evaluation (in the order of 10 to 50 times the number of dimensions) to converge to an optimal solution. This search method is more suitable when the time required for a model evaluation is relatively low.

###  [Chocolate][]

Algorithm name in katib is `chocolate-CMAES`.

## Bayesian Optimization

### [scikit-optimize](https://github.com/scikit-optimize/scikit-optimize)

> Scikit-Optimize, or skopt, is a simple and efficient library to minimize (very) expensive and noisy black-box functions. It implements several methods for sequential model-based optimization. skopt aims to be accessible and easy to use in many contexts.

Algorithm name in katib is `skopt-bayesian-optimization`, and there are some algortihm settings that we support:

| Setting Name     | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | Example  |
|------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| base_estimator   |  ["GP", "RF", "ET", "GBRT" or sklearn regressor, default="GP"]:   Should inherit from `sklearn.base.RegressorMixin`. In addition, the `predict`   method, should have an optional `return_std` argument, which returns   `std(Y | x)` along with `E[Y | x]`. If base_estimator is one of   ["GP", "RF", "ET", "GBRT"], a default surrogate model of the corresponding   type is used corresponding to what is used in the minimize functions. More in [skopt document](https://scikit-optimize.github.io/#skopt.Optimizer) | GP       |
| n_initial_points |  [int, default=10]: Number of evaluations of `func` with initialization points  before approximating it with `base_estimator`. Points provided as `x0` count  as initialization points. If len(x0) < n_initial_points additional points  are sampled at random. More in [skopt document](https://scikit-optimize.github.io/#skopt.Optimizer)                                                                                                                                                                               | 10       |
| acq_func         |  [string, default=`"gp_hedge"`]: Function to minimize over the posterior distribution. More in [skopt document](https://scikit-optimize.github.io/#skopt.Optimizer)                                                                                                                                                                                                                                                                                                                                                        | gp_hedge |
| acq_optimizer    |  [string, "sampling" or "lbfgs", default="auto"]: Method to minimize the acquistion function.    The fit model is updated with the optimal value obtained by optimizing acq_func with acq_optimizer. More in [skopt document](https://scikit-optimize.github.io/#skopt.Optimizer)                                                                                                                                                                                                                                          | auto     |
| random_state     | [int, RandomState instance, or None (default)]: Set random state to something other than None for reproducible results.                                                                                                                                                                                                                                                                                                                                                                                                    | 10       |

[Chocolate]: https://chocolate.readthedocs.io