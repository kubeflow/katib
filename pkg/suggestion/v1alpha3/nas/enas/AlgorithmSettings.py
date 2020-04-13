def parseAlgorithmSettings(params_raw, logger):

    param_standard = {
        "controller_hidden_size":       ['value', int, [1, 'inf']],
        "controller_temperature":       ['value', float, [0, 'inf']],
        "controller_tanh_const":        ['value', float, [0, 'inf']],
        "controller_entropy_weight":    ['value', float, [0.0, 'inf']],
        "controller_baseline_dec":      ['value', float, [0.0, 1.0]],
        "controller_learning_rate":     ['value', float, [0.0, 1.0]],
        "controller_skip_target":       ['value', float, [0.0, 1.0]],
        "controller_skip_weight":       ['value', float, [0.0, 'inf']],
        "controller_train_steps":       ['value', int, [1, 'inf']],
        "controller_log_every_steps":   ['value', int, [1, 'inf']],
    }

    algorithm_settings = {
        "controller_hidden_size":       64,
        "controller_temperature":       5.,
        "controller_tanh_const":        2.25,
        "controller_entropy_weight":    1e-5,
        "controller_baseline_decay":    0.999,
        "controller_learning_rate":     5e-5,
        "controller_skip_target":       0.4,
        "controller_skip_weight":       0.8,
        "controller_train_steps":       50,
        "controller_log_every_steps":   10,
    }

    # TODO: Enable to add None values, e.g in controller_temperature parameter
    # TODO: Delete it and add to the Validation part
    def checktype(param_name, param_value, check_mode, supposed_type, supposed_range=None, logger=None):
        correct = True

        try:
            converted_value = supposed_type(param_value)
        except:
            correct = False
            logger.info("Parameter {} is of wrong type. Set back to default value {}"
                        .format(param_name, algorithm_settings[param_name]))

        if correct and check_mode == 'value':
            if (
                (supposed_range[0] != '-inf' and
                    ((supposed_type == float and converted_value <= supposed_range[0]) or
                        converted_value < supposed_range[0])
                 ) or
                (supposed_range[1] != 'inf' and converted_value > supposed_range[1])
            ):
                correct = False
                logger.info("Parameter {} out of range. Set back to default value {}"
                            .format(param_name, algorithm_settings[param_name]))

        elif correct and check_mode == 'categorical':
            if converted_value not in supposed_range:
                correct = False
                logger.info("Parameter {} out of range. Set back to default value {}"
                            .format(param_name, algorithm_settings[param_name]))

        if correct:
            algorithm_settings[param_name] = converted_value

    for param in params_raw:
        if param.name in algorithm_settings.keys():
            checktype(param.name,
                      param.value,
                      param_standard[param.name][0],  # mode
                      param_standard[param.name][1],  # type
                      param_standard[param.name][2],  # range
                      logger)
        else:
            logger.info("Unknown Parameter name: {}".format(param.name))

    return algorithm_settings
