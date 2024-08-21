import pprint

import six
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

    @property
    def num_workers(self):
        """Gets the number of workers of distributed training.
        Number of workers is setting number of workers.
        :return: The number of workers of distributed training.
        :rtype: int
        """
        return self._num_workers

    @num_workers.setter
    def num_workers(self, num_workers):
        """Sets the number of workers of distributed training.
        Number of workers is setting number of workers.
        :param num_workers: The number of workers of distributed training.
        :type: int
        """

        self._num_workers = num_workers

    @property
    def num_procs_per_worker(self):
        """Gets the number of processes per worker of distributed training.
        Number of processes per worker is the setting number of processes per worker.
        :return: The number of processed per worker of distributed training.
        :rtype: int
        """
        return self._num_procs_per_worker

    @num_procs_per_worker.setter
    def num_procs_per_worker(self, num_procs_per_worker):
        """Sets the number of processes per worker of distributed training.
        Number of processes per worker is the setting number of processes per worker.
        :param num_procs_per_worker: The number of processes per worker of distributed training.
        :type: int
        """

        self._num_procs_per_worker = num_procs_per_worker

    @property
    def resources_per_worker(self):
        """Gets the resources per worker of distributed training.
        Resources per worker is the setting resources per worker.
        :return: The resources per worker of distributed training.
        :rtype: dict or V1ResourceRequirements
        """
        return self._resources_per_worker

    @resources_per_worker.setter
    def resources_per_worker(self, resources_per_worker):
        """Sets the resources per worker of distributed training.
        Resources per worker is the setting resources per worker.
        :param resources_per_worker: The resources per worker of distributed training.
        :type: dict or V1ResourceRequirements
        """

        self._resources_per_worker = resources_per_worker

    def to_dict(self):
        """Returns the resources properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.__dict__):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(
                    map(lambda x: x.to_dict() if hasattr(x, "to_dict") else x, value)
                )
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(
                    map(
                        lambda item: (
                            (item[0], item[1].to_dict())
                            if hasattr(item[1], "to_dict")
                            else item
                        ),
                        value.items(),
                    )
                )
            else:
                result[attr] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, TrainerResources):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, TrainerResources):
            return True

        return self.to_dict() != other.to_dict()
