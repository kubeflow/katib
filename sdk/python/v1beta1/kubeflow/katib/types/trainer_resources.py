class TrainerResources(object):
    def __init__(
        self,
        num_workers=None,
        num_procs_per_worker=None,
        resources_per_worker=None,
    ):
        self.num_workers = num_workers
        self.num_procs_per_worker = num_procs_per_worker
        self.resources_per_worker = resources_per_worker
