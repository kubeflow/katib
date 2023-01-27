# Copyright 2021 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import shutil
import setuptools

REQUIRES = [
    "certifi>=14.05.14",
    "six>=1.10",
    "setuptools>=21.0.0",
    "urllib3>=1.15.1",
    "kubernetes>=23.6.0",
    "grpcio==1.41.1",
    "protobuf==3.19.5",
]

katib_grpc_api_file = "../../../pkg/apis/manager/v1beta1/python/api_pb2.py"

# Copy Katib gRPC Python APIs to use it in the Katib SDK Client.
# We need to always copy this file only on the SDK building stage, not on SDK installation stage.
if os.path.exists(katib_grpc_api_file):
    shutil.copy(
        katib_grpc_api_file, "kubeflow/katib/katib_api_pb2.py",
    )

setuptools.setup(
    name="kubeflow-katib",
    version="0.15.0rc0",
    author="Kubeflow Authors",
    author_email="premnath.vel@gmail.com",
    license="Apache License Version 2.0",
    url="https://github.com/kubeflow/katib/tree/master/sdk/python/v1beta1",
    description="Katib Python SDK for APIVersion v1beta1",
    long_description="Katib Python SDK for APIVersion v1beta1",
    packages=setuptools.find_packages(include=("kubeflow*")),
    package_data={},
    include_package_data=False,
    zip_safe=False,
    classifiers=[
        "Intended Audience :: Developers",
        "Intended Audience :: Education",
        "Intended Audience :: Science/Research",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3 :: Only",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: OS Independent",
        "Topic :: Scientific/Engineering",
        "Topic :: Scientific/Engineering :: Artificial Intelligence",
        "Topic :: Software Development",
        "Topic :: Software Development :: Libraries",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
    install_requires=REQUIRES,
)
