#!/usr/bin/env python3

import logging
from pathlib import Path
from subprocess import check_call

import yaml
from oci_image import OCIImageResource, OCIImageResourceError
from ops.charm import CharmBase
from ops.framework import StoredState
from ops.main import main
from ops.model import ActiveStatus, MaintenanceStatus, WaitingStatus
from charms.prometheus_k8s.v0.prometheus_scrape import MetricsEndpointProvider
from charms.grafana_k8s.v0.grafana_dashboard import GrafanaDashboardProvider

logger = logging.getLogger(__name__)


class CheckFailed(Exception):
    """Raise this exception if one of the checks in main fails."""

    def __init__(self, msg, status_type=None):
        super().__init__()

        self.msg = msg
        self.status_type = status_type
        self.status = status_type(msg)


class Operator(CharmBase):
    """Deploys the katib-controller service."""

    _stored = StoredState()

    def __init__(self, framework):
        super().__init__(framework)

        self._stored.set_default(**self.gen_certs())
        self.image = OCIImageResource(self, "oci-image")

        self.prometheus_provider = MetricsEndpointProvider(
            charm=self,
            jobs=[
                {
                    "job_name": "katib_controller_metrics",
                    "static_configs": [
                        {"targets": [f"*:{self.config['metrics-port']}"]}
                    ],
                }
            ],
        )
        self.dashboard_provider = GrafanaDashboardProvider(self)

        for event in [
            self.on.config_changed,
            self.on.install,
            self.on.leader_elected,
            self.on.upgrade_charm,
        ]:
            self.framework.observe(event, self.set_pod_spec)

    def set_pod_spec(self, event):
        self.model.unit.status = MaintenanceStatus("Setting pod spec")

        try:
            self._check_leader()

            image_details = self._check_image_details()
        except CheckFailed as check_failed:
            self.model.unit.status = check_failed.status
            return

        validating, mutating = yaml.safe_load_all(Path("src/webhooks.yaml").read_text())

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
                                        "serviceaccounts",
                                        "services",
                                        "events",
                                        "namespaces",
                                        "persistentvolumes",
                                        "persistentvolumeclaims",
                                        "pods",
                                        "pods/log",
                                        "pods/status",
                                    ],
                                    "verbs": ["*"],
                                },
                                {
                                    "apiGroups": ["apps"],
                                    "resources": ["deployments"],
                                    "verbs": ["*"],
                                },
                                {
                                    "apiGroups": ["rbac.authorization.k8s.io"],
                                    "resources": [
                                        "roles",
                                        "rolebindings",
                                    ],
                                    "verbs": ["*"],
                                },
                                {
                                    "apiGroups": ["batch"],
                                    "resources": ["jobs", "cronjobs"],
                                    "verbs": ["*"],
                                },
                                {
                                    "apiGroups": ["kubeflow.org"],
                                    "resources": [
                                        "experiments",
                                        "experiments/status",
                                        "experiments/finalizers",
                                        "trials",
                                        "trials/status",
                                        "trials/finalizers",
                                        "suggestions",
                                        "suggestions/status",
                                        "suggestions/finalizers",
                                        "tfjobs",
                                        "pytorchjobs",
                                        "mpijobs",
                                        "xgboostjobs",
                                        "mxjobs",
                                    ],
                                    "verbs": ["*"],
                                },
                            ],
                        }
                    ],
                },
                "containers": [
                    {
                        "name": "katib-controller",
                        "imageDetails": image_details,
                        "command": ["./katib-controller"],
                        "args": [
                            f"--webhook-port={self.model.config['webhook-port']}",
                            "--trial-resources=Job.v1.batch",
                            "--trial-resources=TFJob.v1.kubeflow.org",
                            "--trial-resources=PyTorchJob.v1.kubeflow.org",
                            "--trial-resources=MPIJob.v1.kubeflow.org",
                            "--trial-resources=PipelineRun.v1beta1.tekton.dev",
                        ],
                        "ports": [
                            {
                                "name": "webhook",
                                "containerPort": self.model.config["webhook-port"],
                            },
                            {
                                "name": "metrics",
                                "containerPort": self.model.config["metrics-port"],
                            },
                        ],
                        "envConfig": {
                            "KATIB_CORE_NAMESPACE": self.model.name,
                        },
                        "volumeConfig": [
                            {
                                "name": "certs",
                                "mountPath": "/tmp/cert",
                                "files": [
                                    {
                                        "path": "tls.crt",
                                        "content": self._stored.cert,
                                    },
                                    {
                                        "path": "tls.key",
                                        "content": self._stored.key,
                                    },
                                ],
                            }
                        ],
                    }
                ],
            },
            k8s_resources={
                "kubernetesResources": {
                    "customResourceDefinitions": [
                        {"name": crd["metadata"]["name"], "spec": crd["spec"]}
                        for crd in yaml.safe_load_all(Path("src/crds.yaml").read_text())
                    ],
                    "mutatingWebhookConfigurations": [
                        {
                            "name": mutating["metadata"]["name"],
                            "webhooks": mutating["webhooks"],
                        }
                    ],
                    "validatingWebhookConfigurations": [
                        {
                            "name": validating["metadata"]["name"],
                            "webhooks": validating["webhooks"],
                        }
                    ],
                },
                "configMaps": {
                    "katib-config": {
                        f: Path(f"src/{f}.json").read_text()
                        for f in (
                            "metrics-collector-sidecar",
                            "suggestion",
                            "early-stopping",
                        )
                    },
                    "trial-template": {
                        f + suffix: Path(f"src/{f}.yaml").read_text()
                        for f, suffix in (
                            ("defaultTrialTemplate", ".yaml"),
                            ("enasCPUTemplate", ""),
                            ("pytorchJobTemplate", ""),
                        )
                    },
                },
            },
        )

        self.model.unit.status = ActiveStatus()

    def gen_certs(self):
        model = self.model.name
        app = self.model.app.name
        Path("/run/ssl.conf").write_text(
            f"""[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn
[ dn ]
C = GB
ST = Canonical
L = Canonical
O = Canonical
OU = Canonical
CN = 127.0.0.1
[ req_ext ]
subjectAltName = @alt_names
[ alt_names ]
DNS.1 = {app}
DNS.2 = {app}.{model}
DNS.3 = {app}.{model}.svc
DNS.4 = {app}.{model}.svc.cluster
DNS.5 = {app}.{model}.svc.cluster.local
IP.1 = 127.0.0.1
[ v3_ext ]
authorityKeyIdentifier=keyid,issuer:always
basicConstraints=CA:FALSE
keyUsage=keyEncipherment,dataEncipherment,digitalSignature
extendedKeyUsage=serverAuth,clientAuth
subjectAltName=@alt_names"""
        )

        check_call(["openssl", "genrsa", "-out", "/run/ca.key", "2048"])
        check_call(["openssl", "genrsa", "-out", "/run/server.key", "2048"])
        check_call(
            [
                "openssl",
                "req",
                "-x509",
                "-new",
                "-sha256",
                "-nodes",
                "-days",
                "3650",
                "-key",
                "/run/ca.key",
                "-subj",
                "/CN=127.0.0.1",
                "-out",
                "/run/ca.crt",
            ]
        )
        check_call(
            [
                "openssl",
                "req",
                "-new",
                "-sha256",
                "-key",
                "/run/server.key",
                "-out",
                "/run/server.csr",
                "-config",
                "/run/ssl.conf",
            ]
        )
        check_call(
            [
                "openssl",
                "x509",
                "-req",
                "-sha256",
                "-in",
                "/run/server.csr",
                "-CA",
                "/run/ca.crt",
                "-CAkey",
                "/run/ca.key",
                "-CAcreateserial",
                "-out",
                "/run/cert.pem",
                "-days",
                "365",
                "-extensions",
                "v3_ext",
                "-extfile",
                "/run/ssl.conf",
            ]
        )

        return {
            "cert": Path("/run/cert.pem").read_text(),
            "key": Path("/run/server.key").read_text(),
            "ca": Path("/run/ca.crt").read_text(),
        }

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


if __name__ == "__main__":
    main(Operator)
