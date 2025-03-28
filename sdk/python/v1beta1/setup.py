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
    "kubernetes>=27.2.0",
    "grpcio>=1.64.1",
    "protobuf>=4.21.12,<5",
    "kubeflow-training==1.9.0",
]

katib_grpc_api_file = "../../../pkg/apis/manager/v1beta1/python/api_pb2.py"
katib_grpc_svc_file = "../../../pkg/apis/manager/v1beta1/python/api_pb2_grpc.py"

# Copy Katib gRPC Python APIs to use it in the Katib SDK Client.
# We need to always copy this file only on the SDK building stage, not on SDK installation stage.
if os.path.exists(katib_grpc_api_file):
    shutil.copy(
        katib_grpc_api_file,
        "kubeflow/katib/katib_api_pb2.py",
    )

# TODO(Electronic-Waste): Remove the import rewrite when protobuf supports `python_package` option.
# REF: https://github.com/protocolbuffers/protobuf/issues/7061
if os.path.exists(katib_grpc_svc_file):
    shutil.copy(
        katib_grpc_svc_file,
        "kubeflow/katib/katib_api_pb2_grpc.py",
    )

    with open("kubeflow/katib/katib_api_pb2_grpc.py", "r+") as file:
        content = file.read()
        new_content = content.replace("api_pb2", "kubeflow.katib.katib_api_pb2")
        file.seek(0)
        file.write(new_content)
        file.truncate()

setuptools.setup(
    name="kubeflow-katib",
    version="0.18.0rc0",
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
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: OS Independent",
        "Topic :: Scientific/Engineering",
        "Topic :: Scientific/Engineering :: Artificial Intelligence",
        "Topic :: Software Development",
        "Topic :: Software Development :: Libraries",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
    install_requires=REQUIRES,
    extras_require={
        "huggingface": ["kubeflow-training[huggingface]==1.9.0"],
    },
)
