from .random_search import RandomSearch
from .grid_search import GridSearch
from .bayesian_optimization_algorithm import BOAlgorithm


ALGORITHM_REGISTER = {"random_search": RandomSearch,
                      "grid_search": GridSearch,
                      "bayesian_optimization": BOAlgorithm}
