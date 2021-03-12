#!/usr/bin/env python3

import logging

from ops.charm import CharmBase
from ops.main import main
from ops.model import ActiveStatus, MaintenanceStatus

from oci_image import OCIImageResource, OCIImageResourceError

logger = logging.getLogger(__name__)


class Operator(CharmBase):
    """Deploys the katib-ui service."""

    def __init__(self, framework):
        super().__init__(framework)

        if not self.model.unit.is_leader():
            logger.info("Not a leader, skipping any work")
            self.model.unit.status = ActiveStatus()
            return

        self.image = OCIImageResource(self, "oci-image")
        self.framework.observe(self.on.install, self.set_pod_spec)
        self.framework.observe(self.on.upgrade_charm, self.set_pod_spec)

    def set_pod_spec(self, event):
        self.model.unit.status = MaintenanceStatus("Setting pod spec")

        try:
            image_details = self.image.fetch()
        except OCIImageResourceError as e:
            self.model.unit.status = e.status
            return

        self.model.pod.set_spec(
            {
                "version": 3,
                "serviceAccount": {
                    "roles": [
                        {
                            "global": True,
                            "rules": [
                                {
                                    "apiGroups": [""],
                                    "resources": [
                                        "configmaps",
                                        "namespaces",
                                    ],
                                    "verbs": ["*"],
                                },
                                {
                                    "apiGroups": ["kubeflow.org"],
                                    "resources": [
                                        "experiments",
                                        "trials",
                                        "suggestions",
                                    ],
                                    "verbs": ["*"],
                                },
                            ],
                        }
                    ]
                },
                "containers": [
                    {
                        "name": "katib-ui",
                        "command": ["./katib-ui"],
                        "args": [f"--port={self.model.config['port']}"],
                        "imageDetails": image_details,
                        "ports": [
                            {
                                "name": "http",
                                "containerPort": self.model.config["port"],
                            }
                        ],
                        "envConfig": {
                            "KATIB_CORE_NAMESPACE": self.model.name,
                        },
                    }
                ],
            },
        )

        self.model.unit.status = ActiveStatus()


if __name__ == "__main__":
    main(Operator)
