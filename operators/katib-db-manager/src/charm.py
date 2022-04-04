#!/usr/bin/env python3

import logging

from oci_image import OCIImageResource, OCIImageResourceError
from ops.charm import CharmBase
from ops.main import main
from ops.model import ActiveStatus, MaintenanceStatus, WaitingStatus

logger = logging.getLogger(__name__)


class CheckFailed(Exception):
    """Raise this exception if one of the checks in main fails."""

    def __init__(self, msg, status_type=None):
        super().__init__()

        self.msg = msg
        self.status_type = status_type
        self.status = status_type(msg)


class Operator(CharmBase):
    """Deploys the katib-db-manager service."""

    def __init__(self, framework):
        super().__init__(framework)

        self.image = OCIImageResource(self, "oci-image")
        self.framework.observe(self.on.install, self.set_pod_spec)
        self.framework.observe(self.on.config_changed, self.set_pod_spec)
        self.framework.observe(self.on.upgrade_charm, self.set_pod_spec)
        self.framework.observe(self.on["mysql"].relation_joined, self.set_pod_spec)
        self.framework.observe(self.on["mysql"].relation_changed, self.set_pod_spec)

    def set_pod_spec(self, event):
        try:
            self._check_leader()

            image_details = self._check_image_details()

            mysql_data = self._check_mysql()
        except CheckFailed as check_failed:
            self.model.unit.status = check_failed.status
            return

        self.model.unit.status = MaintenanceStatus("Setting pod spec")

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

    def _check_leader(self):
        if not self.unit.is_leader():
            # We can't do anything useful when not the leader, so do nothing.
            raise CheckFailed("Waiting for leadership", WaitingStatus)

    def _check_image_details(self):
        try:
            image_details = self.image.fetch()
        except OCIImageResourceError as e:
            raise CheckFailed(f"{e.status.message}", e.status_type)
        return image_details

    def _check_mysql(self):
        try:
            relation = self.model.relations["mysql"][0]
            unit = next(iter(relation.units))
            mysql_data = relation.data[unit]
            # Ensure we've got some data sent over the relation
            mysql_data["root_password"]
        except (IndexError, StopIteration, KeyError):
            raise CheckFailed("Waiting for mysql connection information", WaitingStatus)

        return mysql_data


if __name__ == "__main__":
    main(Operator)
