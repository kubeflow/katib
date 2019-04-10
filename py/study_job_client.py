# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
import datetime
import json
import logging
import multiprocessing
import time

from kubernetes import client as k8s_client
from kubernetes.client import rest

STUDY_JOB_GROUP = "kubeflow.org"
STUDY_JOB_PLURAL = "studyjobs"
STUDY_JOB_KIND = "StudyJob"

TIMEOUT = 120

def wait_for_condition(client,
                       namespace,
                       name,
                       expected_condition=[],
                       version="v1alpha1",
                       timeout=datetime.timedelta(minutes=10),
                       polling_interval=datetime.timedelta(seconds=30),
                       status_callback=None):
  """Waits until any of the specified conditions occur.

  Args:
    client: K8s api client.
    namespace: namespace for the studyjob.
    name: Name of the studyjob.
    expected_condition: A list of conditions. Function waits until any of the
      supplied conditions is reached.
    timeout: How long to wait for the job.
    polling_interval: How often to poll for the status of the job.
    status_callback: (Optional): Callable. If supplied this callable is
      invoked after we poll the job. Callable takes a single argument which
      is the job.
  """
  crd_api = k8s_client.CustomObjectsApi(client)
  end_time = datetime.datetime.now() + timeout
  while True:
    # By setting async_req=True ApiClient returns multiprocessing.pool.AsyncResult
    # If we don't set async_req=True then it could potentially block forever.
    thread = crd_api.get_namespaced_custom_object(
      STUDY_JOB_GROUP, version, namespace, STUDY_JOB_PLURAL, name, async_req=True)

    # Try to get the result but timeout.
    results = None
    try:
      results = thread.get(TIMEOUT)
    except multiprocessing.TimeoutError:
      logging.error("Timeout trying to get studyJob %s/%s.", namespace, name)
    except Exception as e:
      logging.error("There was a problem waiting for studyJob %s/%s; Exception: %s",
                    namespace, name, e)
      raise

    if results:
      if status_callback:
        status_callback(results)

      condition = results.get("status", {}).get("condition")
      if condition in expected_condition:
          return results

    if datetime.datetime.now() + polling_interval > end_time:
      raise Exception(
        "Timeout waiting for studyJob {0} in namespace {1} to enter one of the "
        "conditions {2}.".format(name, namespace, expected_condition), results)

    time.sleep(polling_interval.seconds)

def create_study_job(client, spec, version="v1alpha1"):
  """Create a studyJob.

  Args:
    client: A K8s api client.
    spec: The spec for the job.
  """
  crd_api = k8s_client.CustomObjectsApi(client)
  try:
    # Create a Resource
    namespace = spec["metadata"].get("namespace", "default")
    thread = crd_api.create_namespaced_custom_object(
      STUDY_JOB_GROUP, version, namespace, STUDY_JOB_PLURAL, spec, async_req=True)
    api_response = thread.get(TIMEOUT)
    logging.info("Created studyJob %s", api_response["metadata"]["name"])
    return api_response
  except rest.ApiException as e:
    message = ""
    if e.message:
      message = e.message
    if e.body:
      try:
        body = json.loads(e.body)
      except ValueError:
        # There was a problem parsing the body of the response as json.
        logging.error(
          ("Exception when calling DefaultApi->"
           "apis_fqdn_v1_namespaces_namespace_resource_post. body: %s"), e.body)
        raise
      message = body.get("message")

    logging.error(("Exception when calling DefaultApi->"
                   "apis_fqdn_v1_namespaces_namespace_resource_post: %s"),
                  message)
    raise e

def delete_study_job(client, name, namespace, version="v1alpha1"):
  crd_api = k8s_client.CustomObjectsApi(client)
  try:
    body = {
      # Set garbage collection so that job won't be deleted until all
      # owned references are deleted.
      "propagationPolicy": "Foreground",
    }
    logging.info("Deleting studyJob %s/%s", namespace, name)
    thread = crd_api.delete_namespaced_custom_object(
      STUDY_JOB_GROUP,
      version,
      namespace,
      STUDY_JOB_PLURAL,
      name,
      body,
      async_req=True)
    api_response = thread.get(TIMEOUT)
    logging.info("Deleting studyJob %s/%s returned: %s", namespace, name,
                 api_response)
    return api_response
  except rest.ApiException as e:
    message = ""
    if e.message:
      message = e.message
    if e.body:
      try:
        body = json.loads(e.body)
      except ValueError:
        # There was a problem parsing the body of the response as json.
        logging.error(
          ("Exception when calling DefaultApi->"
           "apis_fqdn_v1_namespaces_namespace_resource_delete. body: %s"), e.body)
        raise
      message = body.get("message")

    logging.error(("Exception when calling DefaultApi->"
                   "apis_fqdn_v1_namespaces_namespace_resource_delete: %s"),
                  message)
    raise e
