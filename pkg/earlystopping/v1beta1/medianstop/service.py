import logging
from kubernetes import client, config
import multiprocessing
from datetime import datetime

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc

logger = logging.getLogger()
logging.basicConfig(level=logging.INFO)

STATUS_EARLY_STOPPED = "EarlyStopped"
KUBEFLOW_GROUP = "kubeflow.org"
KATIB_VERSION = "v1beta1"
TRIAL_PLURAL = "trials"
APISERVER_TIMEOUT = 120

DEFAULT_NAMESPACE = "default"


class MedianStopService(api_pb2_grpc.EarlyStoppingServicer):

    def __init__(self):
        super(MedianStopService, self).__init__()
        # Assume that Trial namespace = Suggestion namespace.
        try:
            with open('/var/run/secrets/kubernetes.io/serviceaccount/namespace', 'r') as f:
                self.namespace = f.readline()
                # Set config and api instance for k8s client.
                config.load_incluster_config()
        # This is used when service is not running in k8s, e.g. for unit tests.
        except Exception as e:
            logger.info("{}. Service is not running in Kubernetes Pod, \"{}\" namespace is used".format(
                e, DEFAULT_NAMESPACE
            ))
            self.namespace = DEFAULT_NAMESPACE
            # Set config and api instance for k8s client.
            config.load_kube_config()

        self.api_instance = client.CustomObjectsApi()

    def GetEarlyStoppingRules(self, request, context):
        logger.info("Get new early stopping rules")
        rules = []
        rules.append(
            api_pb2.EarlyStoppingRule(
                name="accuracy",
                value="0.8",
                comparison=api_pb2.LESS,
            )
        )
        rules.append(
            api_pb2.EarlyStoppingRule(
                name="Epoch",
                value="9",
                comparison=api_pb2.EQUAL,
            )
        )

        logger.info("New rules are\n {}".format(rules))
        return api_pb2.GetEarlyStoppingRulesReply(
            early_stopping_rules=rules
        )

    def SetTrialStatus(self, request, context):
        logger.info(request)
        trial_name = request.trial_name

        logger.info("Update status for Trial {}".format(trial_name))

        # TODO (andreyvelich): Move this part to Katib SDK ?
        # Get Trial object
        thread = self.api_instance.get_namespaced_custom_object(
            KUBEFLOW_GROUP,
            KATIB_VERSION,
            self.namespace,
            TRIAL_PLURAL,
            trial_name,
            async_req=True)

        trial = None
        try:
            trial = thread.get(APISERVER_TIMEOUT)
        except multiprocessing.TimeoutError:
            raise Exception("Timeout trying to get Katib Trial")
        except Exception as e:
            raise Exception(
                "Get Trial {} in namespace {} failed. Exception: {}".format(trial_name, self.namespace, e))

        time_now = datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ")

        early_stopped_condition = {
            "type": STATUS_EARLY_STOPPED,
            "status": "True",
            "reason": "TrialEarlyStopped",
            "message": "Trial is early stopped",
            "lastUpdateTime": time_now,
            "lastTransitionTime": time_now,
        }
        trial["status"]["conditions"].append(early_stopped_condition)

        # Update Trial object with early stopped status
        try:
            self.api_instance.patch_namespaced_custom_object_status(
                KUBEFLOW_GROUP,
                KATIB_VERSION,
                self.namespace,
                TRIAL_PLURAL,
                trial_name,
                trial,
                async_req=True)
        except Exception as e:
            raise Exception(
                "Update status for Trial {} in namespace {} failed. Exception: {}".format(
                    trial_name, self.namespace, e))

        logger.info("Changed status to {} for Trial {} in namespace {}".format(
            STATUS_EARLY_STOPPED, trial_name, self.namespace))

        return api_pb2.SetTrialStatusReply()
