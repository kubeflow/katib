#!/usr/bin/env python3

import logging

from ops.charm import CharmBase
from ops.main import main
from ops.model import ActiveStatus, MaintenanceStatus, WaitingStatus

from oci_image import OCIImageResource, OCIImageResourceError

logger = logging.getLogger(__name__)


class Operator(CharmBase):
    """Deploys the katib-db-manager service."""

    def __init__(self, framework):
        super().__init__(framework)

        if not self.model.unit.is_leader():
            logger.info("Not a leader, skipping any work")
            self.model.unit.status = ActiveStatus()
            return

        self.image = OCIImageResource(self, "oci-image")
        self.framework.observe(self.on.install, self.set_pod_spec)
        self.framework.observe(self.on.config_changed, self.set_pod_spec)
        self.framework.observe(self.on.upgrade_charm, self.set_pod_spec)
        self.framework.observe(self.on["mysql"].relation_joined, self.set_pod_spec)
        self.framework.observe(self.on["mysql"].relation_changed, self.set_pod_spec)

    def set_pod_spec(self, event):
        self.model.unit.status = MaintenanceStatus("Setting pod spec")

        try:
            image_details = self.image.fetch()
        except OCIImageResourceError as e:
            self.model.unit.status = e.status
            return

        try:
            relation = self.model.relations["mysql"][0]
            unit = next(iter(relation.units))
            mysql_data = relation.data[unit]
            # Ensure we've got some data sent over the relation
            mysql_data["root_password"]
        except (IndexError, StopIteration, KeyError):
            self.model.unit.status = WaitingStatus(
                "Waiting for mysql connection information"
            )
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
                        "name": "katib-db-manager",
                        "command": ["./katib-db-manager"],
                        "imageDetails": image_details,
                        "ports": [
                            {
                                "name": "api",
                                "containerPort": self.model.config["port"],
                            }
                        ],
                        "envConfig": {
                            "DB_NAME": "mysql",
                            "DB_USER": "root",
                            "DB_PASSWORD": mysql_data["root_password"],
                            "KATIB_MYSQL_DB_HOST": mysql_data["host"],
                            "KATIB_MYSQL_DB_PORT": mysql_data["port"],
                            "KATIB_MYSQL_DB_DATABASE": mysql_data["database"],
                        },
                        "kubernetes": {
                            "readinessProbe": {
                                "exec": {
                                    "command": [
                                        "/bin/grpc_health_probe",
                                        f"-addr=:{self.model.config['port']}",
                                    ]
                                },
                                "initialDelaySeconds": 5,
                            },
                            "livenessProbe": {
                                "exec": {
                                    "command": [
                                        "/bin/grpc_health_probe",
                                        f"-addr=:{self.model.config['port']}",
                                    ]
                                },
                                "initialDelaySeconds": 10,
                                "periodSeconds": 60,
                                "failureThreshold": 5,
                            },
                        },
                    }
                ],
            },
        )

        self.model.unit.status = ActiveStatus()


if __name__ == "__main__":
    main(Operator)
