#!/usr/bin/env python
# Implementation based on:
#   https://github.com/ray-project/ray/blob/7f1bacc7dc9caf6d0ec042e39499bbf1d9a7d065/python/ray/tune/examples/pbt_example.py

# Copyright 2022 The Kubeflow Authors.
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

import argparse
import os
import pickle
import random
import time

import numpy as np

# Ensure job runs for at least this long (secs) to allow metrics collector to
# read PID correctly before cleanup
_METRICS_COLLECTOR_SPAWN_LATENCY = 7


class PBTBenchmarkExample:
    """Toy PBT problem for benchmarking adaptive learning rate.
    The goal is to optimize this trainable's accuracy. The accuracy increases
    fastest at the optimal lr, which is a function of the current accuracy.
    The optimal lr schedule for this problem is the triangle wave as follows.
    Note that many lr schedules for real models also follow this shape:
     best lr
      ^
      |    /\
      |   /  \
      |  /    \
      | /      \
      ------------> accuracy
    In this problem, using PBT with a population of 2-4 is sufficient to
    roughly approximate this lr schedule. Higher population sizes will yield
    faster convergence. Training will not converge without PBT.
    """

    def __init__(self, lr, checkpoint: str):
        self._lr = lr

        self._checkpoint_file = os.path.join(checkpoint, "training.ckpt")
        if os.path.exists(self._checkpoint_file):
            with open(self._checkpoint_file, "rb") as fin:
                checkpoint_data = pickle.load(fin)
            self._accuracy = checkpoint_data["accuracy"]
            self._step = checkpoint_data["step"]
        else:
            os.makedirs(checkpoint, exist_ok=True)
            self._step = 1
            self._accuracy = 0.0

    def save_checkpoint(self):
        with open(self._checkpoint_file, "wb") as fout:
            pickle.dump({"step": self._step, "accuracy": self._accuracy}, fout)

    def step(self):
        midpoint = 50  # lr starts decreasing after acc > midpoint
        q_tolerance = 3  # penalize exceeding lr by more than this multiple
        noise_level = 2  # add gaussian noise to the acc increase
        # triangle wave:
        #  - start at 0.001 @ t=0,
        #  - peak at 0.01 @ t=midpoint,
        #  - end at 0.001 @ t=midpoint * 2,
        if self._accuracy < midpoint:
            optimal_lr = 0.01 * self._accuracy / midpoint
        else:
            optimal_lr = 0.01 - 0.01 * (self._accuracy - midpoint) / midpoint
        optimal_lr = min(0.01, max(0.001, optimal_lr))

        # compute accuracy increase
        q_err = max(self._lr, optimal_lr) / (
            min(self._lr, optimal_lr) + np.finfo(float).eps
        )
        if q_err < q_tolerance:
            self._accuracy += (1.0 / q_err) * random.random()
        elif self._lr > optimal_lr:
            self._accuracy -= (q_err - q_tolerance) * random.random()
        self._accuracy += noise_level * np.random.normal()
        self._accuracy = max(0, min(100, self._accuracy))

        self._step += 1

    def __repr__(self):
        return "epoch {}:\nlr={:0.4f}\nValidation-accuracy={:0.4f}".format(
            self._step, self._lr, self._accuracy / 100
        )


if __name__ == "__main__":
    # Parse CLI arguments
    parser = argparse.ArgumentParser(description="PBT Basic Test")
    parser.add_argument(
        "--lr", type=float, default=0.0001, help="learning rate (default: 0.0001)"
    )
    parser.add_argument(
        "--epochs", type=int, default=20, help="number of epochs to train (default: 20)"
    )
    parser.add_argument(
        "--checkpoint",
        type=str,
        default="/var/log/katib/checkpoints/",
        help="checkpoint directory (resume and save)",
    )
    opt = parser.parse_args()

    benchmark = PBTBenchmarkExample(opt.lr, opt.checkpoint)

    start_time = time.time()
    for i in range(opt.epochs):
        benchmark.step()
    exec_time_thresh = time.time() - start_time - _METRICS_COLLECTOR_SPAWN_LATENCY
    if exec_time_thresh < 0:
        time.sleep(abs(exec_time_thresh))
    benchmark.save_checkpoint()

    print(benchmark)
