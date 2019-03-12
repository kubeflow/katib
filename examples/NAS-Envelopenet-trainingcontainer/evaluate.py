from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

from datetime import datetime
import math
import numpy as np
import tensorflow as tf
import net

class Evaluate(object):

    def __init__(self, arch, config, checkpoint_dir):
        self.task_config = config
        self.arch = arch
        self.checkpoint_dir = checkpoint_dir
        self.get_task_params()

    def get_task_params(self):
        self.log_stats = self.task_config["log_stats"]
        self.dataset = self.task_config["dataset"]
        self.image_size = self.task_config["input_size"]
        self.init_cell={
        "Layer0": {"Branch0": {"block": "conv2d", "kernel_size": [1, 1], "outputs": 64}},
        "Layer2": {"Branch0": {"block": "lrn" }}
                }
        self.classification_cell={
        "Layer0": {"Branch0": {"block": "reduce_mean", "size": [1, 2]}},
        "Layer1": {"Branch0": {"block": "flatten", "size": [3, 3]}},
        "Layer2": {"Branch0": {"block": "dropout", "keep_prob": 0.8}},
        "Layer3": {"Branch0": {"block": "fc-final", "inputs": 192, "outputs": 10}}
                }
        self.batch_size = self.task_config["batch_size"]
        self.num_examples = 10000
        self.run_once = True
        self.eval_dir = self.task_config["data_dir"] + "/results/" + \
            "nac_envelopenet" + "/" + "/evaluate"
        self.evaluate()
        

    def eval_once(self, saver, summary_writer, top_k_op, summary_op, k=1):
        """Run Eval once.
        Args:
          saver: Saver.
          summary_writer: Summary writer.
          top_k_op: Top K op.
          summary_op: Summary op.
        """
        with tf.Session() as sess:
            ckpt = tf.train.get_checkpoint_state(self.checkpoint_dir)
            if ckpt and ckpt.model_checkpoint_path:
                # Restores from checkpoint
                saver.restore(sess, ckpt.model_checkpoint_path)
                # Assuming model_checkpoint_path looks something like:
                #   /my-favorite-path/cifar10_train/model.ckpt-0,
                # extract global_step from it.
                global_step = ckpt.model_checkpoint_path.split(
                    '/')[-1].split('-')[-1]
            else:
                print('No checkpoint file found')
                return

            # Start the queue runners.
            coord = tf.train.Coordinator()
            try:
                threads = []
                for q_runner in tf.get_collection(tf.GraphKeys.QUEUE_RUNNERS):
                    threads.extend(
                        q_runner.create_threads(
                            sess,
                            coord=coord,
                            daemon=True,
                            start=True))

                num_iter = int(math.ceil(self.num_examples / self.batch_size))
                true_count = 0  # Counts the number of correct predictions.
                total_sample_count = num_iter * self.batch_size
                step = 0
                while step < num_iter and not coord.should_stop():
                    predictions = sess.run([top_k_op])
                    true_count += np.sum(predictions)
                    step += 1

                if k == 1:
                    # Compute precision @ 1.
                    precision = true_count / total_sample_count
                    print(
                        '%s: precision @ 1 = %.3f' %
                        (datetime.now(), precision))
                elif k == 5:
                    # Compute precision @ 5.
                    precision = true_count / total_sample_count
                    print(
                        '%s: precision @ 5 = %.3f' %
                        (datetime.now(), precision))

                summary = tf.Summary()
                summary.ParseFromString(sess.run(summary_op))
                summary.value.add(
                    tag='Precision @ %d' %
                    (k), simple_value=precision)
                summary_writer.add_summary(summary, global_step)
            except Exception as excpn:  # pylint: disable=broad-except
                coord.request_stop(excpn)

            coord.request_stop()
            coord.join(threads, stop_grace_period_secs=10)

    def evaluate(self):
        network = net.Net(self.arch, self.task_config)
        """Eval a network for a number of steps."""
        with tf.Graph().as_default() as grph:
            # Get images and labels for CIFAR-10.
            eval_data = True
            images, labels = network.inputs(eval_data=eval_data)

            # Build a Graph that computes the logits predictions from the
            # inference model.
            # TODO: Clean up all args
            arch = self.arch
            init_cell = self.init_cell
            classification_cell = self.classification_cell
            log_stats = self.log_stats
            scope = "Nacnet"
            is_training = False
            logits = network.inference(images,
                                       arch,
                                       init_cell,
                                       classification_cell,
                                       log_stats,
                                       is_training,
                                       scope)

            # Calculate predictions.
            # if imagenet is running then run precision@1,5
            top_k_op = tf.nn.in_top_k(logits, labels, 1)
            if self.dataset == "imagenet":
                    # Quick dirty fixes to incorporate changes brought by
                    # imagenet
                self.num_examples = 50000
                top_5_op = tf.nn.in_top_k(logits, labels, 5)

            # Restore the moving average version of the learned variables for
            # eval.
            variable_averages = tf.train.ExponentialMovingAverage(
                net.MOVING_AVERAGE_DECAY)
            variables_to_restore = variable_averages.variables_to_restore()
            saver = tf.train.Saver(variables_to_restore)

            # Build the summary operation based on the TF collection of
            # Summaries.
            summary_op = tf.summary.merge_all()

            summary_writer = tf.summary.FileWriter(self.eval_dir, grph)

            while True:
                self.eval_once(saver, summary_writer, top_k_op, summary_op)
                if self.dataset == "imagenet":
                    self.eval_once(saver, summary_writer, top_5_op,
                                   summary_op,
                                   k=5)
                if self.run_once:
                    break

