from kubeflow.katib.configuration import Configuration


class TrainerResources(object):
    def __init__(
        self,
        num_workers=None,
        num_procs_per_worker=None,
        resources_per_worker=None,
        local_vars_configuration=None,
    ):
        if local_vars_configuration is None:
            local_vars_configuration = Configuration()
        self.local_vars_configuration = local_vars_configuration

        self._num_workers = None
        self._num_procs_per_worker = None
        self._resources_per_worker = None

        if num_workers is not None:
            self.num_workers = num_workers
        if num_procs_per_worker is not None:
            self.num_procs_per_worker = num_procs_per_worker
        if resources_per_worker is not None:
            self.resources_per_worker = resources_per_worker
