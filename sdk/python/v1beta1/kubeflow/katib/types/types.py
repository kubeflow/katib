from dataclasses import dataclass
from typing import Dict


# Trainer resources for distributed training.
@dataclass
class TrainerResources:
    num_workers: int
    num_procs_per_worker: int
    resources_per_worker: Dict[str, str]
