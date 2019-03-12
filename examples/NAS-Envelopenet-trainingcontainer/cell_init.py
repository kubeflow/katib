"""Initialization (Stem) cell"""
import tensorflow as tf
from cell import Cell

slim = tf.contrib.slim

def trunc_normal(stddev):
    return tf.truncated_normal_initializer(0.0, stddev)

class Init(Cell):
    """Initialization (Stem) cell: The first cell of a CNN"""
    def __init__(self, cellidx):
        self.cellidx = cellidx
        self.cellname = "Init"
        Cell.__init__(self)

    def cell(self, inputs, arch, is_training):
        """Create the cell by instantiating the cell blocks"""
        nscope = 'Cell_' + self.cellname + '_' + str(self.cellidx)
        reuse = None
        with tf.variable_scope(nscope, 'initial_block', [inputs], reuse=reuse) as scope:
            with slim.arg_scope([slim.conv2d, slim.max_pool2d], stride=1, padding='SAME'):
                net = inputs
                layeridx = 0
                for layer in sorted(arch.keys()):
                    cells = []
                    for branch in sorted(arch[layer].keys()):
                        block = arch[layer][branch]
                        if block["block"] == "conv2d":
                            output_filters = int(block["outputs"])
                            kernel_size = block["kernel_size"]
                            if "stride" not in block.keys():
                                stride = 1
                            else:
                                stride = block["stride"]
                            cell = slim.conv2d(
                                net,
                                output_filters,
                                kernel_size,
                                stride=stride,
                                padding='SAME') 
                        elif block["block"] == "max_pool":
                            kernel_size = block["kernel_size"]
                            cell = slim.max_pool2d(
                                net, kernel_size, padding='SAME', stride=2)
                        elif block["block"] == "lrn":
                            cell = tf.nn.lrn(
                                net, 4, bias=1.0, alpha=0.001 / 9.0, beta=0.75)
                        elif block["block"] == "dropout":
                            keep_prob = block["keep_prob"]
                            cell = slim.dropout(net, keep_prob=keep_prob, is_training=is_training)
                        elif block["block"] == "cutout":
                            if not is_training:
                              cell = net
                            else:
                              cutout_size = block["size"]
                              img_dim = int(net.shape[1])
                              batch_size = int(net.shape[0])
                              channels = int(net.shape[3])
                              """ Puts white rectange on a RGB image with random x,y coordinates """
                              mask = tf.ones([cutout_size, cutout_size], dtype=tf.int32)
                              start = tf.random_uniform([2], minval=0, maxval=img_dim, dtype=tf.int32)
                              mask = tf.pad(mask, [[cutout_size + start[0], img_dim - start[0]],
                                                   [cutout_size + start[1], img_dim - start[1]]])
                              mask = mask[cutout_size: cutout_size + img_dim,
                                          cutout_size: cutout_size + img_dim]
                              mask = tf.reshape(mask, [img_dim, img_dim, 1])
                              mask = tf.expand_dims(mask, axis=0)
                              mask = tf.tile(mask, [batch_size, 1, 1, channels])
                              cell = tf.where(tf.equal(mask, 0), x=net, y=tf.zeros_like(net))
                        else:
                            print("Invalid block")
                            exit(-1)
                        cells.append(cell)
                    net = tf.concat(cells, axis=-1)

                    layeridx += 1
        return net
