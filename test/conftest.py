import os
from sys import path

root = os.path.join(os.path.dirname(__file__), "..")
path.extend(
    [
        os.path.join(root, "pkg/apis/manager/v1beta1/python"),
        os.path.join(root, "pkg/apis/manager/health/python"),
        os.path.join(root, "pkg/metricscollector/v1beta1/common"),
        os.path.join(root, "pkg/metricscollector/v1beta1/tfevent-metricscollector"),
    ]
)
