import numpy as np

from pkg.suggestion.cma_unit_test.db_init import connect_db
from pkg.suggestion.cma_unit_test.internface import get_worker, get_trials, store_worker_logs


def func(x1, x2):
    return 0.75 * np.exp(-(9 * x1 - 2) ** 2 / 4 - (9 * x2 - 2) ** 2 / 4) + 0.75 * np.exp(
        -(9 * x1 + 1) ** 2 / 49 - (9 * x2 + 1) / 10) + \
           0.5 * np.exp(-(9 * x1 - 7) ** 2 / 4 - (9 * x2 - 3) ** 2 / 4) - 0.2 * np.exp(
        -(9 * x1 - 4) ** 2 - (9 * x2 - 7) ** 2)


def spawn_worker(worker_id, worker_config):
    cnx = connect_db()
    worker = get_worker(cnx, worker_id)
    trial = get_trials(cnx, worker.trial_id, "")[0]
    x1 = float(trial.parameter_set[0].value)
    x2 = float(trial.parameter_set[1].value)
    objective = func(x1, x2)
    print(objective)

    store_worker_logs(cnx, worker_id, "precision", [objective])
