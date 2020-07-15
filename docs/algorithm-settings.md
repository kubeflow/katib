# Hyperparameter Tuning Algorithms

Table of Contents
=================

   * [Hyperparameter Tuning Algorithms](#hyperparameter-tuning-algorithms)
   * [Table of Contents](#table-of-contents)
      * [Grid Search](#grid-search)
         * [<a href="https://chocolate.readthedocs.io" rel="nofollow">Chocolate</a>](#chocolate)
      * [Random Search](#random-search)
         * [<a href="http://hyperopt.github.io/hyperopt/" rel="nofollow">Hyperopt</a>](#hyperopt)
      * [TPE](#tpe)
         * [<a href="http://hyperopt.github.io/hyperopt/" rel="nofollow">Hyperopt</a>](#hyperopt-1)
      * [Bayesian Optimization](#bayesian-optimization)
         * [<a href="https://github.com/scikit-optimize/scikit-optimize">scikit-optimize</a>](#scikit-optimize)
      * [References](#references)

Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

<!-- ## Quasi Random Search

QuasiRandom sampling ensures a much more uniform exploration of the search space than traditional pseudo random. Thus, quasi random sampling is preferable when not all variables are discrete, the number of dimensions is high and the time required to evaluate a solution is high.

### [Chocolate](https://chocolate.readthedocs.io)

Algorithm name in katib is `chocolate-quasirandom`. -->

<!-- ## CMAES

CMAES search is one of the most powerful black-box optimization algorithm. However, it requires a significant number of model evaluation (in the order of 10 to 50 times the number of dimensions) to converge to an optimal solution. This search method is more suitable when the time required for a model evaluation is relatively low.

###  [Chocolate](https://chocolate.readthedocs.io)

Algorithm name in katib is `chocolate-CMAES`. -->

For information about the hyperparameter tuning algorithms and neural
architecture search implemented or integrated in Katib, see the detailed guide
to [configuring and running a Katib 
experiment](https://kubeflow.org/docs/components/hyperparameter-tuning/experiment/).
For information about supported algorithms in Katib, see the [Katib configuration settings](https://kubeflow.org/docs/components/hyperparameter-tuning/katib-config/#suggestion-settings).
