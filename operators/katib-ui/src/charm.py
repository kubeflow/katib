#!/usr/bin/env python3

import logging

from oci_image import OCIImageResource, OCIImageResourceError
from ops.charm import CharmBase
from ops.main import main
from ops.model import ActiveStatus, BlockedStatus, MaintenanceStatus, WaitingStatus
from serialized_data_interface import (
    NoCompatibleVersions,
    NoVersionsListed,
    get_interfaces,
)

logger = logging.getLogger(__name__)


class Operator(CharmBase):
    """Deploys the katib-ui service."""

    def __init__(self, framework):
        super().__init__(framework)

        if not self.model.unit.is_leader():
            logger.info("Not a leader, skipping any work")
            self.model.unit.status = ActiveStatus()
            return

        try:
            self.interfaces = get_interfaces(self)
        except NoVersionsListed as err:
            self.model.unit.status = WaitingStatus(str(err))
            return
        except NoCompatibleVersions as err:
            self.model.unit.status = BlockedStatus(str(err))
            return
        else:
            self.model.unit.status = ActiveStatus()

        self.image = OCIImageResource(self, "oci-image")
        self.framework.observe(self.on.install, self.set_pod_spec)
        self.framework.observe(self.on.upgrade_charm, self.set_pod_spec)
        self.framework.observe(
            self.on["ingress"].relation_changed,
            self.configure_ingress,
        )

    def configure_ingress(self, event):
        if self.interfaces["ingress"]:
            self.interfaces["ingress"].send_data(
                {
                    "prefix": "/katib/",
                    "service": self.model.app.name,
                    "port": self.model.config["port"],
                }
            )

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
