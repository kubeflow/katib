from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

from datetime import datetime
import ast
import sys
import time
import json

import tensorflow as tf
import net
from evaluate import Evaluate
global_batch_size = 128
global_log_frequency = 10

class _LoggerHook(tf.train.SessionRunHook):
    """Logs loss and runtime."""

    def __init__(self, global_step_init, loss):
        self.global_step_init = global_step_init
        self.loss = loss

    def begin(self):
        self._step = self.global_step_init
        self._start_time = time.time()

    def before_run(self, run_context):
        self._step += 1
        # Asks for loss value.
        return tf.train.SessionRunArgs(self.loss)

    def after_run(self, run_context, run_values):
        if self._step % global_log_frequency == 0:
            current_time = time.time()
            duration = current_time - self._start_time
            self._start_time = current_time

            loss_value = run_values.results
            examples_per_sec = global_log_frequency * global_batch_size / duration
            sec_per_batch = float(duration / global_log_frequency)

            format_str = (
                '%s: step %d, loss = %.2f (%.1f examples/sec; %.3f '
                'sec/batch)')
            print(
                format_str %
                (datetime.now(),
                 self._step,
                 loss_value,
                 examples_per_sec,
                 sec_per_batch))

class ModelConstructor(object):
    """Trainingtask
    """
    def __init__(self,arch, params,i=0):
        
        self.name = 'train'
        self.task_config = arch
        self.iteration=i
        self.max_iterations = params["iterations"]
        self.params=params
        self.train_dir=params["data_dir"]+ "/results/"+str(self.iteration)+"/train"
        if self.iteration==0 or self.iteration==self.max_iterations+1:
            self.max_steps=10*params["steps"]
        else:
            self.max_steps=params["steps"]
        self.get_params(params)  
        
    def get_params(self, params):
        global global_batch_size
        global global_log_frequency
        self.count_params = False
        self.log_device_placement = False
        global_log_frequency = 10
        self.log_stats = params["log_stats"]
        self.gpus = params["gpus"]
        gpu_fraction = params["gpu_usage"]
        self.gpu_options = tf.GPUOptions(per_process_gpu_memory_fraction=gpu_fraction)
        self.batch_size = params["batch_size"]
        self.dataset = params["dataset"]
        self.image_size = params["input_size"]
        self.arch = self.task_config
        global_batch_size = self.batch_size
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
        self.child_training={
        "optimizer": {"type": "momentum", "momentum": 0.9},
        "lr": {"type":"exponential_decay", "initial":0.040},
        "gradient_clipping": {"type": "norm", "value":5.0},
        "regularization": {"type": "l2", "value": 3e-4}
        }
        self.params=params

    def build_model(self):

        self.network = net.Net(self.task_config, self.params)
        #self.network.maybe_download_and_extract()
        with tf.Graph().as_default():
            ckpt = tf.train.get_checkpoint_state(self.train_dir)
            global_step_init = -1
            if ckpt and ckpt.model_checkpoint_path:
                global_step_init = int(
                    ckpt.model_checkpoint_path.split('/')[-1].split('-')[-1])
                global_step = tf.Variable(
                    global_step_init,
                    name='global_step',
                    dtype=tf.int64,
                    trainable=False)
            else:
                global_step = tf.train.get_or_create_global_step()

            images, labels = self.network.distorted_inputs()
            if self.gpus:
                batch_queue = tf.contrib.slim.prefetch_queue.prefetch_queue(
                        [images, labels], capacity=2 * len(self.gpus))
                tower_grads = []

            arch = self.arch
            init_cell = self.init_cell
            classification_cell = self.classification_cell
            log_stats = self.log_stats
            scope = "Nacnet"
            is_training = True

            if self.gpus:
                # Multi-gpu setting
                learning_rate = self.network.get_learning_rate(global_step, self.child_training)
                tf.summary.scalar('learning_rate', learning_rate)
                opt = self.network.get_opt(learning_rate, self.child_training)
                with tf.variable_scope(tf.get_variable_scope()):
                    for i in self.gpus:
                        i=int(i)
                        with tf.device('/gpu:%d' % i):
                            with tf.name_scope('%s_%d' % ('tower', i)) as scope:
                                # Dequeues one batch for the GPU
                                image_batch, label_batch = batch_queue.dequeue()
                                logits = self.network.inference(image_batch,
                                                           arch,
                                                           init_cell,
                                                           classification_cell,
                                                           log_stats,
                                                           is_training,
                                                           scope)
                                # Calculate the loss for one tower of the CIFAR model. This function
                                # constructs the entire CIFAR model but shares the variables across
                                # all towers.
                                loss = self.network.tower_loss(scope, logits, label_batch)
                                loss = self.network.get_regularization_loss(loss, self.child_training)
                                tf.get_variable_scope().reuse_variables()
                                # Retain the summaries from the final tower. TODO:
                                # not a nice way to use the last iteration of the
                                # loop
                                summaries = tf.get_collection(
                                    tf.GraphKeys.SUMMARIES, scope)
                                grads = opt.compute_gradients(loss)
                                grads = self.network.clip_gradients(grads, self.child_training)

                                tower_grads.append(grads)
                grads = self.network.average_gradients(tower_grads)
                for grad, var in grads:
                    if grad is not None:
                        summaries.append(
                            tf.summary.histogram(
                                var.op.name + '/gradients', grad))
        
                apply_gradient_op = opt.apply_gradients(
                    grads, global_step=global_step)

                variable_averages = tf.train.ExponentialMovingAverage(
                    net.MOVING_AVERAGE_DECAY, global_step)
                variables_averages_op = variable_averages.apply(
                    tf.trainable_variables())

                train_op = tf.group(apply_gradient_op, variables_averages_op)
            else:
                logits = self.network.inference(images,
                                       arch,
                                       init_cell,
                                       classification_cell,
                                       log_stats,
                                       is_training,
                                       scope)
                loss = self.network.loss(logits, labels)
                train_op = self.network.get_train_op(loss, global_step, self.child_training)

            if self.count_params:
                # For counting parameters
                param_stats = tf.profiler.profile(
                    tf.get_default_graph(),
                    options=tf.profiler.ProfileOptionBuilder()
                    .with_max_depth(2)
                    .with_accounted_types(['_trainable_variables'])
                    .select(['params'])
                    .build())
                # For counting flops
                flop_stats = tf.profiler.profile(
                    tf.get_default_graph(),
                    options=tf.profiler.ProfileOptionBuilder() .with_max_depth(1) .select(
                        ['float_ops']).build())
                print(param_stats)
                print(flop_stats)
                exit()

            saver = tf.train.Saver()
            with tf.train.MonitoredTrainingSession(
                        checkpoint_dir=self.train_dir,
                        hooks=[tf.train.StopAtStepHook(last_step=self.max_steps),
                               tf.train.NanTensorHook(loss),
                               _LoggerHook(global_step_init, loss)],
                        save_checkpoint_secs=300,
                        save_summaries_steps=100,
                        config=tf.ConfigProto(
                            log_device_placement=self.log_device_placement,
                            allow_soft_placement=True,
                            gpu_options=self.gpu_options
                            )) as mon_sess:

                ckpt = tf.train.get_checkpoint_state(self.train_dir)
                if ckpt and ckpt.model_checkpoint_path:
                    print("Restoring existing model")
                    saver.restore(mon_sess, ckpt.model_checkpoint_path)

                while not mon_sess.should_stop():
                    mon_sess.run(train_op)

    def evaluate():
        eval=Evaluate(self.arch, self.params, self.train_dir)
        
