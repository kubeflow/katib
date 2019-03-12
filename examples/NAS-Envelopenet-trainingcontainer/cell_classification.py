"""Classification cell"""

import tensorflow as tf
from cell import Cell

slim = tf.contrib.slim

def trunc_normal(stddev):
    return tf.truncated_normal_initializer(0.0, stddev)

class Classification(Cell):
    """Classification cell: The final classification block of a CNN"""
    def __init__(self):
        self.cellname = "Classification"
        Cell.__init__(self)
    def __del__(self):
        pass
    def cell(self, inputs, arch, is_training):
        """Create the cell by instantiating the cell blocks"""
        nscope = 'Cell_' + self.cellname
        net = inputs
        reuse = None
        with tf.variable_scope(nscope, 'classification_block', [inputs], reuse=reuse) as scope:
            for layer in sorted(arch.keys()):
                for branch in sorted(arch[layer].keys()):
                    block = arch[layer][branch]
                    if block["block"] == "reduce_mean":
                        net = tf.reduce_mean(net, [1, 2])
                    elif block["block"] == "flatten":
                        net = slim.flatten(net)
                    elif block["block"] == "fc":
                        outputs = block["outputs"]
                        net = slim.fully_connected(net, outputs)
                    elif block["block"] == "fc-final":
                        outputs = block["outputs"]
                        inputs = block["inputs"]
                        weights_initializer = trunc_normal(1 / float(inputs))
                        biases_initializer = tf.zeros_initializer()
                        net = slim.fully_connected(
                            net,
                            outputs,
                            biases_initializer=biases_initializer,
                            weights_initializer=weights_initializer,
                            weights_regularizer=None,
                            activation_fn=None)
                    elif block["block"] == "dropout":
                        keep_prob = block["keep_prob"]
                        net = slim.dropout(
                            net, keep_prob=keep_prob, is_training=is_training)
                    else:
                        print("Invalid block")
                        exit(-1)
        return net
