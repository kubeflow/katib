from .bayesian_optimization_algorithm import BOAlgorithm
from .random_search import RandomSearch

ALGORITHM_REGISTER = {"random_search": RandomSearch,
                      "bayesian_optimization": BOAlgorithm}
