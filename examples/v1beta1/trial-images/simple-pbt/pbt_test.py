#!/usr/bin/env python

# Implementation based on:
#   https://github.com/ray-project/ray/blob/master/python/ray/tune/examples/pbt_example.py

import argparse
import numpy as np
import os
import pickle
import random
import tensorflow as tf
import time


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

    def __init__(self, lr, log_dir: str, log_interval: int, checkpoint: str):
        # Allow lazy creation of tfevent file
        self._log_dir = log_dir
        self._writer = None
        self._log_interval = log_interval
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
        midpoint = 100  # lr starts decreasing after acc > midpoint
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
        q_err = max(self._lr, optimal_lr) / min(self._lr, optimal_lr)
        if q_err < q_tolerance:
            self._accuracy += (1.0 / q_err) * random.random()
        elif self._lr > optimal_lr:
            self._accuracy -= (q_err - q_tolerance) * random.random()
        self._accuracy += noise_level * np.random.normal()
        self._accuracy = max(0, self._accuracy)

        if self._step == 1 or self._step % self._log_interval == 0:
            self.save_checkpoint()
            if not self._writer:
                self._writer = tf.summary.create_file_writer(self._log_dir)
            with self._writer.as_default():
                tf.summary.scalar(
                    "Validation-accuracy", self._accuracy, step=self._step
                )
                tf.summary.scalar("lr", self._lr, step=self._step)
                self._writer.flush()

        self._step += 1

    def __repr__(self):
        return "epoch {}:\nlr={:0.4f}\nValidation-accuracy={:0.4f}".format(
            self._step, self._lr, self._accuracy
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
        "--log-interval",
        type=int,
        default=10,
        metavar="N",
        help="how many batches to wait before logging training status (default: 1)",
    )
    parser.add_argument(
        "--log-path",
        type=str,
        default="/var/log/katib/tfevent/",
        help="tfevent output path (default: /var/log/katib/tfevent/)",
    )
    parser.add_argument(
        "--checkpoint",
        type=str,
        default="/var/log/katib/checkpoints/",
        help="checkpoint directory (resume and save)",
    )
    opt = parser.parse_args()

    benchmark = PBTBenchmarkExample(
        opt.lr, opt.log_path, opt.log_interval, opt.checkpoint
    )
    for i in range(opt.epochs):
        benchmark.step()
        time.sleep(0.2)
    print(benchmark)
