import random
import string

from google.protobuf.json_format import MessageToJson, Parse

from pkg.api.python import api_pb2


def generate_randid():
    return ''.join(random.sample(string.ascii_letters + string.digits, 16))


def create_study(cnx, study_config):
    configs = MessageToJson(study_config.parameter_configs)
    tags = []
    for tag in study_config.tags:
        tags.append(MessageToJson(tag))

    # todo: duplicate study id
    study_id = generate_randid()
    add_study = ("INSERT INTO studies "
                 "VALUES (%(study_id)s,"
                 " %(name)s,"
                 " %(owner)s,"
                 " %(optimization_type)s, "
                 " %(optimization_goal)s,"
                 " %(parameter_configs)s,"
                 " %(suggest_algo)s,"
                 " %(early_stop_algo)s,"
                 " %(tags)s,"
                 " %(objective_value_name)s,"
                 " %(metrics)s)")
    data_study = {
        'study_id': study_id,
        'name': study_config.name,
        'owner': study_config.owner,
        'optimization_type': study_config.optimization_type,
        'optimization_goal': study_config.optimization_goal,
        'parameter_configs': configs,
        'suggest_algo': study_config.default_suggestion_algorithm,
        'early_stop_algo': study_config.default_early_stopping_algorithm,
        'tags': "&\n".join(tags),
        'objective_value_name': study_config.objective_value_name,
        'metrics': "&\n".join(study_config.metrics),
    }
    cnx.cursor().execute(add_study, data_study)
    cnx.commit()
    cnx.cursor().close()

    return study_id


def get_study_config(cnx, study_id):
    cursor = cnx.cursor()
    cursor.execute("SELECT * FROM studies WHERE id = '%s'" % study_id)
    row = cursor.fetchone()

    study = api_pb2.StudyConfig()
    _, study.name, study.owner, study.optimization_type, study.optimization_goal, configs, \
    study.default_suggestion_algorithm, study.default_early_stopping_algorithm, tags, study.objective_value_name, metrics = row

    Parse(configs, study.parameter_configs)

    temp_tag = api_pb2.Tag()
    tags_array = tags.split("&\n")
    for tag in tags_array:
        if tag != '':
            study.tags.extend([Parse(tag, temp_tag)])

    study.metrics.extend(metrics.split("&\n"))
    cursor.close()
    return study


def set_suggestion_param(cnx, algorithm, study_id, params):
    suggestion_params = []
    for param in params:
        suggestion_params.append(MessageToJson(param))

    # todo: duplicate id
    param_id = generate_randid()
    cursor = cnx.cursor()
    add_param = ("INSERT INTO suggestion_param "
                 "VALUES (%(id)s,"
                 " %(suggestion_algo)s,"
                 " %(parameters)s,"
                 " %(study_id)s)")
    data_param = {
        "id": param_id,
        'suggestion_algo': algorithm,
        'parameters': "&\n".join(suggestion_params),
        'study_id': study_id,
    }
    cursor.execute(add_param, data_param)
    cnx.commit()
    cursor.close()

    return param_id


def update_suggestion_param(cnx, param_id, params):
    suggestion_params = []
    for param in params:
        suggestion_params.append(MessageToJson(param))

    cursor = cnx.cursor()
    cursor.execute(
        "UPDATE suggestion_param SET parameters = '%s' WHERE id = '%s'" % ('&\n'.join(suggestion_params).replace("\\", "\\\\"), param_id))
    cnx.commit()
    cursor.close()


def get_suggestion_param(cnx, param_id):
    cursor = cnx.cursor()
    cursor.execute("SELECT * FROM suggestion_param WHERE id = '%s'" % param_id)
    row = cursor.fetchone()
    _, _, params, _ = row
    # todo: split
    params = params.split("&\n")
    ret = []
    for param in params:
        temp = api_pb2.SuggestionParameter()
        ret.append(Parse(param, temp))

    cursor.close()
    return ret


def get_suggestion_param_list(cnx, study_id):
    cursor = cnx.cursor()
    cursor.execute("SELECT * FROM suggestion_param WHERE study_id = '%s'" % study_id)
    parameter_set = []
    for (id, suggestion_algo, parameters, study_id) in cursor:
        params = []
        parameters = parameters.split("&\n")
        param_name = ""
        for param in parameters:
            temp = api_pb2.SuggestionParameter()
            params.append(Parse(param, temp))
            param_name = temp.name
        parameter_set.append(api_pb2.GetSuggestionParameterListReply.SuggestionParameterSet(
            param_id=id,
            param_name=param_name,
            suggestion_parameters=params
        ))

    return parameter_set


def get_trials(cnx, trial_id, study_id):
    cursor = cnx.cursor()
    if trial_id != "":
        cursor.execute("SELECT * FROM trials WHERE id = '%s'" % trial_id)
    else:
        cursor.execute("SELECT * FROM trials WHERE study_id = '%s'" % study_id)

    ret = []
    for (id, study_id, parameters, status, objective_value, tags) in cursor:
        params = parameters.split("&\n")
        param_list = []
        for param in params:
            temp = api_pb2.Parameter()
            param_list.append(Parse(param, temp))

        tags = tags.split("&\n")
        tag_list = []
        for tag in tags:
            if tag != "":
                temp = api_pb2.Tag()
                tag_list.append(Parse(tag, temp))

        ret.append(api_pb2.Trial(
            trial_id=id,
            study_id=study_id,
            parameter_set=param_list,
            status=status,
            objective_value=objective_value,
            tags=tag_list,
        ))

    cursor.close()
    return ret


def create_trial(cnx, trial):
    params = []
    for param in trial.parameter_set:
        params.append(MessageToJson(param))

    tags = []
    for tag in trial.tags:
        tags.append(MessageToJson(tag))

    # todo: duplicate id
    trial_id = generate_randid()
    cursor = cnx.cursor()
    add_trial = ("INSERT INTO trials "
                 "VALUES (%(id)s,"
                 " %(study_id)s,"
                 " %(parameters)s,"
                 " %(status)s, "
                 " %(objective_value)s,"
                 " %(tag)s)")
    data_trial = {
        "id": trial_id,
        "study_id": trial.study_id,
        "parameters": "&\n".join(params),
        "status": trial.status,
        "objective_value": trial.objective_value,
        "tag": "&\n".join(tags),
    }
    cursor.execute(add_trial, data_trial)
    cnx.commit()
    cursor.close()
    return trial_id


def update_trial_status(cnx, trial_id, status):
    cursor = cnx.cursor()
    cursor.execute("UPDATE trials SET status = '%s' WHERE id = '%s'" % (status, trial_id))
    cnx.commit()
    cursor.close()


def update_trial_value(cnx, trial_id, objective_value):
    cursor = cnx.cursor()
    cursor.execute("UPDATE trials SET objective_value = '%s' WHERE id = '%s'" % (objective_value, trial_id))
    cnx.commit()
    cursor.close()
