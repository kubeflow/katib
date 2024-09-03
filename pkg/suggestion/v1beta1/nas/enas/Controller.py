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

import numpy as np
import tensorflow as tf


class Controller(object):
    def __init__(
        self,
        num_layers=12,
        num_operations=16,
        controller_hidden_size=64,
        controller_temperature=5.0,
        controller_tanh_const=2.25,
        controller_entropy_weight=1e-5,
        controller_baseline_decay=0.999,
        controller_learning_rate=5e-5,
        controller_skip_target=0.4,
        controller_skip_weight=0.8,
        controller_name="controller",
        logger=None,
    ):

        self.logger = logger
        self.logger.info(">>> Building Controller\n")

        self.num_layers = num_layers
        self.num_operations = num_operations

        self.controller_hidden_size = controller_hidden_size
        self.controller_temperature = controller_temperature
        self.controller_tanh_const = controller_tanh_const
        self.controller_entropy_weight = controller_entropy_weight
        self.controller_baseline_decay = controller_baseline_decay
        self.controller_learning_rate = controller_learning_rate
        self.controller_skip_target = controller_skip_target
        self.controller_skip_weight = controller_skip_weight

        self.controller_name = controller_name

        self._build_params()
        self._build_sampler()

    def _build_params(self):
        """Create TF parameters"""
        self.logger.info(">>> Building Controller Parameters\n")
        initializer = tf.compat.v1.random_uniform_initializer(minval=-0.01, maxval=0.01)
        hidden_size = self.controller_hidden_size

        with tf.compat.v1.variable_scope(self.controller_name, initializer=initializer):
            with tf.compat.v1.variable_scope("lstm"):
                self.w_lstm = tf.compat.v1.get_variable(
                    "w", [2 * hidden_size, 4 * hidden_size]
                )

            self.g_emb = tf.compat.v1.get_variable("g_emb", [1, hidden_size])

            with tf.compat.v1.variable_scope("embedding"):
                self.w_emb = tf.compat.v1.get_variable(
                    "w", [self.num_operations, hidden_size]
                )

            with tf.compat.v1.variable_scope("softmax"):
                self.w_soft = tf.compat.v1.get_variable(
                    "w", [hidden_size, self.num_operations]
                )

            with tf.compat.v1.variable_scope("attention"):
                self.attn_w_1 = tf.compat.v1.get_variable(
                    "w_1", [hidden_size, hidden_size]
                )
                self.attn_w_2 = tf.compat.v1.get_variable(
                    "w_2", [hidden_size, hidden_size]
                )
                self.attn_v = tf.compat.v1.get_variable("v", [hidden_size, 1])

        num_params = sum(
            [
                np.prod(v.shape)
                for v in tf.compat.v1.trainable_variables()
                if v.name.startswith(self.controller_name)
            ]
        )
        self.logger.info(">>> Controller has {} Trainable params\n".format(num_params))

    def _build_sampler(self):
        """Build the sampler ops and the log_prob ops."""
        self.logger.info(">>> Building Controller Sampler\n")

        hidden_size = self.controller_hidden_size

        arc_seq = []
        sample_log_probs = []
        sample_entropies = []

        skip_penalties = []
        skip_count = []

        all_h = []
        all_h_w = []

        prev_c = tf.zeros([1, hidden_size], tf.float32)
        prev_h = tf.zeros([1, hidden_size], tf.float32)

        skip_targets = tf.constant(
            [1.0 - self.controller_skip_target, self.controller_skip_target],
            dtype=tf.float32,
        )

        inputs = self.g_emb

        for layer_id in range(self.num_layers):

            next_c, next_h = _lstm(inputs, prev_c, prev_h, self.w_lstm)
            prev_c, prev_h = next_c, next_h

            logits = tf.matmul(next_h, self.w_soft)

            if self.controller_temperature is not None:
                logits /= self.controller_temperature
            if self.controller_tanh_const is not None:
                logits = self.controller_tanh_const * tf.tanh(logits)

            func = tf.random.categorical(logits, 1)
            func = tf.dtypes.cast(func, tf.int32)
            func = tf.reshape(func, [1])

            arc_seq.append(func)

            log_prob = tf.nn.sparse_softmax_cross_entropy_with_logits(
                logits=logits, labels=func
            )

            sample_log_probs.append(log_prob)
            entropy = log_prob * tf.exp(-log_prob)
            entropy = tf.stop_gradient(entropy)
            sample_entropies.append(entropy)
            inputs = tf.nn.embedding_lookup(params=self.w_emb, ids=func)

            next_c, next_h = _lstm(inputs, prev_c, prev_h, self.w_lstm)
            prev_c, prev_h = next_c, next_h

            if layer_id > 0:

                query = tf.matmul(next_h, self.attn_w_2)
                query = query + tf.concat(all_h_w, axis=0)
                query = tf.tanh(query)
                query = tf.matmul(query, self.attn_v)

                logits = tf.concat([-query, query], axis=1)

                if self.controller_temperature is not None:
                    logits /= self.controller_temperature
                if self.controller_tanh_const is not None:
                    logits = self.controller_tanh_const * tf.tanh(logits)

                skip_index = tf.random.categorical(logits, 1)
                skip_index = tf.dtypes.cast(skip_index, tf.int32)

                skip_index = tf.reshape(skip_index, [layer_id])
                arc_seq.append(skip_index)

                skip_prob = tf.sigmoid(logits)
                kl = skip_prob * tf.math.log(skip_prob / skip_targets)
                kl = tf.reduce_sum(input_tensor=kl)
                skip_penalties.append(kl)

                log_prob = tf.nn.sparse_softmax_cross_entropy_with_logits(
                    logits=logits, labels=skip_index
                )

                sample_log_probs.append(
                    tf.reduce_sum(input_tensor=log_prob, keepdims=True)
                )

                entropy = tf.stop_gradient(
                    tf.reduce_sum(
                        input_tensor=log_prob * tf.exp(-log_prob), keepdims=True
                    )
                )
                sample_entropies.append(entropy)

                skip_index = tf.dtypes.cast(skip_index, tf.float32)
                skip_index = tf.reshape(skip_index, [1, layer_id])

                skip_count.append(tf.reduce_sum(input_tensor=skip_index))

                inputs = tf.matmul(skip_index, tf.concat(all_h, axis=0))

                inputs /= 1.0 + tf.reduce_sum(input_tensor=skip_index)
            else:
                inputs = self.g_emb

            all_h.append(next_h)
            all_h_w.append(tf.matmul(next_h, self.attn_w_1))

        arc_seq = tf.concat(arc_seq, axis=0)
        self.sample_arc = tf.reshape(arc_seq, [-1])

        sample_entropies = tf.stack(sample_entropies)
        self.sample_entropy = tf.reduce_sum(input_tensor=sample_entropies)

        sample_log_probs = tf.stack(sample_log_probs, axis=0)
        self.sample_log_probs = tf.reduce_sum(input_tensor=sample_log_probs)

        skip_penalties = tf.stack(skip_penalties)
        self.skip_penalties = tf.reduce_mean(input_tensor=skip_penalties)

        skip_count = tf.stack(skip_count)
        self.skip_count = tf.reduce_sum(input_tensor=skip_count)

    def build_trainer(self):
        """Build the train ops by connecting Controller with candidate."""
        self.child_val_accuracy = tf.compat.v1.placeholder(tf.float32, shape=())

        self.reward = self.child_val_accuracy

        normalize = tf.dtypes.cast(
            (self.num_layers * (self.num_layers - 1) / 2), tf.float32
        )
        self.skip_rate = tf.dtypes.cast((self.skip_count / normalize), tf.float32)

        if self.controller_entropy_weight is not None:
            self.reward += self.controller_entropy_weight * self.sample_entropy

        self.sample_log_probs = tf.reduce_sum(input_tensor=self.sample_log_probs)
        self.baseline = tf.Variable(0.0, dtype=tf.float32, trainable=False)
        baseline_update = tf.compat.v1.assign_sub(
            self.baseline,
            (1 - self.controller_baseline_decay) * (self.baseline - self.reward),
        )

        with tf.control_dependencies([baseline_update]):
            self.reward = tf.identity(self.reward)

        self.loss = self.sample_log_probs * (self.reward - self.baseline)

        if self.controller_skip_weight is not None:
            self.loss += self.controller_skip_weight * self.skip_penalties

        self.train_step = tf.Variable(
            0,
            dtype=tf.int32,
            trainable=False,
            name=self.controller_name + "_train_step",
        )

        tf_variables = [
            var
            for var in tf.compat.v1.trainable_variables()
            if var.name.startswith(self.controller_name)
        ]

        self.train_op, self.grad_norm = _build_train_op(
            loss=self.loss,
            tf_variables=tf_variables,
            train_step=self.train_step,
            learning_rate=self.controller_learning_rate,
        )


# TODO: will remove this function and use tf.nn.LSTMCell instead
def _lstm(x, prev_c, prev_h, w_lstm):
    ifog = tf.matmul(tf.concat([x, prev_h], axis=1), w_lstm)
    i, f, o, g = tf.split(ifog, 4, axis=1)
    i = tf.sigmoid(i)
    f = tf.sigmoid(f)
    o = tf.sigmoid(o)
    g = tf.tanh(g)
    next_c = i * g + f * prev_c
    next_h = o * tf.tanh(next_c)
    return next_c, next_h


def _build_train_op(loss, tf_variables, train_step, learning_rate):
    """Build training ops from `loss` tensor."""
    optimizer = tf.compat.v1.train.AdamOptimizer(learning_rate)
    grads = tf.gradients(ys=loss, xs=tf_variables)

    grad_norm = tf.linalg.global_norm(grads)
    train_op = optimizer.apply_gradients(
        zip(grads, tf_variables), global_step=train_step
    )

    return train_op, grad_norm
