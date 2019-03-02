class ParameterConfig:


    def __init__(self, name_ids, dim, lower_bounds, upper_bounds,
                 parameter_types, names, discrete_info, categorical_info):
        self.name_ids = name_ids
        self.dim = dim
        self.lower_bounds = lower_bounds
        self.upper_bounds = upper_bounds
        self.parameter_types = parameter_types
        self.names = names
        self.discrete_info = discrete_info
        self.categorical_info = categorical_info
