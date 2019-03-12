""" Envelope Cell"""
import tensorflow as tf
import sys
from cell import Cell

slim = tf.contrib.slim

class CellEnvelope(Cell):
    """ Defintion of an envelope cell"""
    def __init__(
            self,
            cellidx,
            channelwidth,
            net,
            filters,
            log_stats,
            outputs):
        self.cellidx = cellidx
        self.log_stats = log_stats
        self.res=sys.argv[2]
        self.cellname = "Envelope"
        self.numbranches = 4
        self.numbins = 100
        self.batchsize = int(net.shape[0])
        numfilters = len(filters)
        self.output_per_filter = outputs
        img_dims = int(net.shape[1])
        self.imagesize = [img_dims, img_dims]
        Cell.__init__(self)
        scope = 'Cell' + str(self.cellidx)
        if self.log_stats:
            with tf.variable_scope(scope, reuse=False):
                for branch in filters:
                    with tf.variable_scope(branch, reuse=False):
                        self.init_stats()

    def cell(self, inputs, channelwidth, is_training=True, filters=None):
        """
        Args:
          inputs: a tensor of size [batch_size, height, width, channels].
          By default use stride=1 and SAME padding
        """
        dropout_keep_prob = 0.8
        nscope = 'Cell_' + self.cellname + '_' + str(self.cellidx)

        scope = 'Cell' + str(self.cellidx)
        nets = []
        with tf.variable_scope(scope):
            for branch in sorted(filters):
                with tf.variable_scope(branch):
                    conv_h, conv_w = branch[0], branch[0]
                    outchannels = self.output_per_filter
                    if branch.endswith("sep"):
                        net = slim.separable_conv2d(
                            inputs, outchannels, [
                                conv_h, conv_w], 1, normalizer_fn=slim.batch_norm)
                    else:
                        net = slim.conv2d(
                            inputs, outchannels, [
                                conv_h, conv_w], normalizer_fn=slim.batch_norm)
                if self.log_stats:
                    msss = self.calc_stats(net, branch)
                    net = tf.Print(
                        net,
                        [msss],
                        message="MeanSSS=:" +
                        scope +
                        "/" +
                        branch +
                        ":")
                net = slim.dropout(
                    net,
                    keep_prob=dropout_keep_prob,
                    scope='dropout',
                    is_training=is_training)
                nets.append(net)
            net = tf.concat(axis=3, values=nets)
        return net

    def init_stats(self):
        size = [
            self.batchsize,
            self.imagesize[0],
            self.imagesize[1],
            self.output_per_filter]
        sumsquaredsamples = tf.contrib.framework.model_variable(
            "sumsquaredsamples", size, initializer=tf.zeros_initializer)
        sumsamples = tf.contrib.framework.model_variable(
            "sumsamples", size, initializer=tf.zeros_initializer)
        samplecount = tf.contrib.framework.model_variable(
            "samplecount", [1], initializer=tf.zeros_initializer)

    def calc_stats(self, inputs, scope):
        with tf.variable_scope(scope, reuse=True):
            size = [
                self.batchsize,
                self.imagesize[0],
                self.imagesize[1],
                self.output_per_filter]
            sumsquaredsamples = tf.get_variable("sumsquaredsamples", size)
            sumsamples = tf.get_variable("sumsamples", size)

            samplecount = tf.get_variable("samplecount", [1])
            tsamplecount = tf.add(samplecount, tf.to_float(tf.constant(1)))
            samplecount = samplecount.assign(tsamplecount)

            """ input is N*H*W*C. We need to calcualte running variance over 
            time (i.e over the N Images in this batch and in all batches.
             Hence need to reduce across the N dimension """
            sum_across_batch = tf.reduce_sum(inputs, axis=0)
            tsumsamples = tf.add(sumsamples, sum_across_batch)
            sumsamples = sumsamples.assign(tsumsamples)
            squared_inputs = tf.square(inputs)
            squared_sum_across_batch = tf.reduce_sum(squared_inputs, axis=0)
            tsumsquaredsamples = tf.add(
                sumsquaredsamples, squared_sum_across_batch)
            sumsquaredsamples = sumsquaredsamples.assign(tsumsquaredsamples)

            msss = (1 / samplecount) * (sumsquaredsamples)
            msss = tf.reduce_mean(msss)

            return msss

    def init_entropy(self):
        bincount = tf.contrib.framework.model_variable(
            "bincount", [self.numbins], initializer=tf.zeros_initializer)
        featuremapsum = tf.contrib.framework.model_variable(
            "featuremapsum", [1], initializer=tf.zeros_initializer)
        featuremapcount = tf.contrib.framework.model_variable(
            "featuremapcount", [1], initializer=tf.zeros_initializer)

    def calc_entropy(self, inputs, scope):
        with tf.variable_scope(scope, reuse=True):
            maxtensor = tf.to_float(tf.size(inputs))

            bincount = tf.get_variable("bincount", [self.numbins])
            featuremapsum = tf.get_variable("featuremapsum", [1])
            featuremapcount = tf.get_variable("featuremapcount", [1])
            inputs = tf.Print(inputs, [inputs, tf.shape(
                inputs)], message="Framemap:", summarize=100)
            binnum = tf.to_int32(
                tf.floor((tf.reduce_sum(inputs) / maxtensor) * (self.numbins - 1)))
            tbincount = tf.scatter_add(
                bincount, binnum, tf.to_float(
                    tf.constant(1)))
            bincount = bincount.assign(tbincount)
            bincount = tf.Print(bincount,
                                [tf.count_nonzero(bincount)],
                                message="Non zero bins count:")

            tfeaturemapsum = tf.add(featuremapsum, tf.reduce_sum(inputs))
            featuremapsum = featuremapsum.assign(tfeaturemapsum)

            tfeaturemapcount = tf.add(featuremapcount, tf.to_float(tf.constant(1)))
            featuremapcount = featuremapcount.assign(tfeaturemapcount)

            meanactivation = tf.divide(featuremapsum, featuremapcount)
            pbin = tf.divide(tf.to_float(bincount), tf.to_float(featuremapcount))
            entropy = tf.multiply(pbin, tf.log(pbin))
            entropy = tf.where(
                tf.is_nan(entropy),
                tf.zeros_like(entropy),
                entropy)
            entropy = tf.reduce_sum(entropy)
            entropy = tf.Print(entropy, [entropy], message=": raw entropy: ")
            entropy = tf.multiply(entropy, tf.multiply(
                meanactivation, tf.constant(-1.0)))
            entropy = tf.Print(
                entropy, [
                    scope, entropy], message=": scaled entropy: ")
            return entropy
