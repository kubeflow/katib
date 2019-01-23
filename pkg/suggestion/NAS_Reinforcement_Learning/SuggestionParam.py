
def parseSuggestionParam(params_raw):
    param_standard = {
        "lstm_num_cells": ['value', int, [1, 'inf']],
        "lstm_num_layers": ['value', int, [1, 'inf']],
        "lstm_keep_prob": ['value', float, [0.0, 1.0]],
        "optimizer": ['categorical', str, ["adam", "momentum", "momentum"]],
        "init_learning_rate": ['value', float, [1e-6, 1.0]],
        "lr_decay_start": ['value', int, [0, 'inf']],
        "lr_decay_every": ['value', int, [1, 'inf']],
        "lr_decay_rate": ['value', float, [0.0, 1.0]],
        "skip-target": ['value', float, [0.0, 1.0]],
        "skip-weight": ['value', float, [0.0, 'inf']],
        "l2_reg": ['value', float, [0.0, 'inf']],
        "entropy_weight": ['value', float, [0.0, 'inf']],
        "baseline_decay": ['value', float, [0.0, 1.0]],
    }

    suggestion_params = {
        "lstm_num_cells": 64,
        "lstm_num_layers": 1,
        "lstm_keep_prob": 1.0,
        "optimizer": "adam",
        "init_learning_rate": 1e-3,
        "lr_decay_start": 0,
        "lr_decay_every": 1000,
        "lr_decay_rate": 0.9,
        "skip-target": 0.4,
        "skip-weight": 0.8,
        "l2_reg": 0,
        "entropy_weight": 1e-4,
        "baseline_decay": 0.9999
    }

    def checktype(param_name, param_value, check_mode, supposed_type, supposed_range=None):
        correct = True

        try:
            converted_value = supposed_type(param_value)
        except:
            correct = False
            print("Parameter {} is of wrong type. Set back to default value {}"
                  .format(param_name, suggestion_params[param_name]))

        if correct and check_mode == 'value':
            if not ((supposed_range[0] == '-inf' or converted_value >= supposed_range[0]) and
                    (supposed_range[1] == 'inf' or converted_value <= supposed_range[1])):
                correct = False
                print("Parameter {} out of range. Set back to default value {}"
                      .format(param_name, suggestion_params[param_name]))
        elif correct and check_mode == 'categorical':
            if converted_value not in supposed_range:
                correct = False
                print("Parameter {} out of range. Set back to default value {}"
                      .format(param_name, suggestion_params[param_name]))

        if correct:
            suggestion_params[param_name] = converted_value


    for param in params_raw:
        if param.name in suggestion_params.keys():
            checktype(param.name,
                      param.value,
                      param_standard[param.name][0],  # mode
                      param_standard[param.name][1],  # type
                      param_standard[param.name][2])  # range
        else:
            print("Unknown Parameter name: {}".format(param.name))

    return suggestion_params
