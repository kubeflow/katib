from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

import os
import re
import sys
import tarfile
from six.moves import urllib

import tensorflow as tf
import tensorflow.contrib.slim as slim
import cifar10_input
import cell_init
import cell_main
import cell_classification
TOWER_NAME = 'tower'
NUM_EXAMPLES_PER_EPOCH_FOR_TRAIN = cifar10_input.NUM_EXAMPLES_PER_EPOCH_FOR_TRAIN
NUM_EXAMPLES_PER_EPOCH_FOR_EVAL = cifar10_input.NUM_EXAMPLES_PER_EPOCH_FOR_EVAL

DATA_URL = 'https://www.cs.toronto.edu/~kriz/cifar-10-binary.tar.gz'
NUM_EPOCHS_PER_DECAY = 2
LEARNING_RATE_DECAY_FACTOR = 0.999    # Learning rate decay factor.
INITIAL_LEARNING_RATE = 0.1             # Initial learning rate.
MOVING_AVERAGE_DECAY = 0.9999


class Net:
    def __init__(self, task_config,params):
        self.task_config = task_config
        self.params=params
        self.cells = []
        self.end_points = []
        self.nets = []
        self.use_fp16 = False
        self.get_task_params()

    def get_task_params(self):
        self.batch_size = self.params["batch_size"]
        self.dataset = self.params["dataset"]
        self.image_size = self.params["input_size"]
        self.arch = self.task_config
        self.data_dir = self.params["data_dir"]
    

    def distorted_inputs(self):
        """Construct distorted input for a given dataset using the Reader ops.
        Returns:
            images: Images. 4D tensor of [batch_size, IMAGE_SIZE, IMAGE_SIZE, 3] size.
            labels: Labels. 1D tensor of [batch_size] size.
        Raises:
            ValueError: If no data_dir
        """
        if not self.data_dir:
            raise ValueError('Please supply a data_dir')
        if self.dataset == 'cifar10':
            data_dir = os.path.join(self.data_dir, 'cifar-10-batches-bin')
            images, labels = cifar10_input.distorted_inputs(
                data_dir=data_dir, batch_size=self.batch_size, image_size=self.image_size)
        elif self.dataset == 'imagenet':
            images, labels = imagenet_input.distorted_inputs()
        if self.use_fp16:
            images = tf.cast(images, tf.float16)
            labels = tf.cast(labels, tf.float16)
        return images, labels

    def inputs(self, eval_data):
        """Construct input for CIFAR evaluation using the Reader ops.
        Args:
            eval_data: bool, indicating if one should use the train or eval data set.
        Returns:
            images: Images. 4D tensor of [batch_size, IMAGE_SIZE, IMAGE_SIZE, 3] size.
            labels: Labels. 1D tensor of [batch_size] size.
        Raises:
            ValueError: If no data_dir
        """
        if not self.data_dir:
            raise ValueError('Please supply a data_dir')
        if self.dataset == 'cifar10':
            data_dir = os.path.join(self.data_dir, 'cifar-10-batches-bin')
            images, labels = cifar10_input.inputs(
                eval_data=eval_data, data_dir=data_dir, batch_size=self.batch_size, image_size=self.image_size)
        elif self.dataset == 'imagenet':
            data_dir = self.data_dir
            if self.dataset_split_name == "test":
                self.dataset_split_name = "validation"
            images, labels = imagenet_input.inputs()
        if self.use_fp16:
            images = tf.cast(images, tf.float16)
            labels = tf.cast(labels, tf.float16)
        return images, labels

    def inference(self, images, arch=None,
                  initcell=None,
                  classificationcell=None,
                  log_stats=False,
                  is_training=None,
                  scope='Nacnet'):

        softmax_linear = self.gen_amlanet(
            images,
            arch,
            initcell,
            classificationcell,
            log_stats,
            is_training,
            scope)
        return softmax_linear

    def maybe_download_and_extract(self):
        """Download and extract the tarball from Alex's website."""
        if self.dataset == 'cifar10':
            dest_directory = self.data_dir
            if not os.path.exists(dest_directory):
                os.makedirs(dest_directory)
            filename = DATA_URL.split('/')[-1]
            filepath = os.path.join(dest_directory, filename)
            if not os.path.exists(filepath):
                def _progress(count, block_size, total_size):
                    sys.stdout.write('\r>> Downloading %s %.1f%%' % (filename,
                         float(count * block_size) / float(total_size) * 100.0))
                    sys.stdout.flush()
                filepath, _ = urllib.request.urlretrieve(
                    DATA_URL, filepath, _progress)
                print()
                statinfo = os.stat(filepath)
                print(
                    'Successfully downloaded',
                    filename,
                    statinfo.st_size,
                    'bytes.')
            extracted_dir_path = os.path.join(
                dest_directory, 'cifar-10-batches-bin')
            if not os.path.exists(extracted_dir_path):
                tarfile.open(filepath, 'r:gz').extractall(dest_directory)
        elif self.dataset == 'imagenet':
            """ It is assumed that if imagenet dataset is specified then it already exists
                and not supposed to be downloaded
            """
            if not os.path.exists(self.data_dir):
                print("Directory {} doesn't exist!".format(self.data_dir))
                exit(-1)
        else:
            print("Unknown dataset {}".format(self.dataset))
            exit(-1)

    def get_macro_net(self,inputs, log_stats=False, is_training=True, arch=None):
        print(self.arch)
        arch = self.arch["network"]
        net = inputs
        cellnumber = 1  # Init block is 0
        nets = [inputs]

        channelwidth = int(inputs.shape[3])
        for celltype in arch:
            if 'filters' in celltype:
                if "inputs" in celltype.keys():
                    all_inputs = [net]
                    input_dim = net.shape
                    if celltype["inputs"] == "all":
                        for index, reduced_inputs in enumerate(nets[:-1]):
                            while reduced_inputs.shape[1] != input_dim[1]:
                                reduced_inputs = slim.max_pool2d(
                                    reduced_inputs, [2, 2], padding='SAME')
                            diffrentiable_scalar = tf.get_variable(name='{}-{}'.format(index, cellnumber), shape=[1],
                                initializer=tf.initializers.random_normal(mean=0.5, stddev=0.01))
                            reduced_inputs = diffrentiable_scalar * reduced_inputs
                            if log_stats:
                                l2_norm = self.calc_l2norm(reduced_inputs)
                                reduced_inputs = tf.Print(reduced_inputs, [l2_norm],
                                    message="l2norm:source-{}dest-{}:".format(index, cellnumber))
                                reduced_inputs = tf.Print(reduced_inputs, [diffrentiable_scalar],
                                    message="scalar:source-{}dest-{}:".format(index, cellnumber))
                            all_inputs.append(reduced_inputs)
                    else:
                        for input_conn in celltype["inputs"]:
                            reduced_inputs = nets[input_conn]
                            while reduced_inputs.shape[1] != input_dim[1]:
                                reduced_inputs = slim.max_pool2d(
                                    reduced_inputs, [2, 2], padding='SAME')
                            diffrentiable_scalar = tf.get_variable(name='{}-{}'.format(input_conn, cellnumber), shape=[1],
                                initializer=tf.initializers.random_normal(mean=0.5, stddev=0.01))
                            reduced_inputs = diffrentiable_scalar * reduced_inputs
                            if log_stats:
                                l2_norm = self.calc_l2norm(reduced_inputs)
                                reduced_inputs = tf.Print(reduced_inputs, [l2_norm],
                                    message="l2norm:source-{}dest-{}:".format(input_conn, cellnumber))
                                reduced_inputs = tf.Print(reduced_inputs, [diffrentiable_scalar],
                                    message="scalar:source-{}dest-{}:".format(input_conn, cellnumber))
                            all_inputs.append(reduced_inputs)
                    net = tf.concat(axis=3, values=all_inputs)
                    num_channels = int(input_dim[3])
                    net = slim.conv2d(
                        net, num_channels, [
                            1, 1], scope='BottleneckLayer_1x1_Envelope_' + str(cellnumber))
                outputs = int(celltype["outputs"] /
                              len(celltype["filters"]))
                envelope = cell_main.CellEnvelope(
                    cellnumber,
                    channelwidth,
                    net,
                    filters=celltype["filters"],
                    log_stats=log_stats,
                    outputs=outputs)
                net = envelope.cell(
                    net, channelwidth, is_training, filters=celltype["filters"])

            elif 'widener' in celltype:
                nscope = 'Widener_' + str(cellnumber) + '_MaxPool_2x2'
                net1 = slim.max_pool2d(
                    net, [2, 2], scope=nscope, padding='SAME')
                nscope = 'Widener_' + str(cellnumber) + '_conv_3x3'
                net2 = slim.conv2d(
                    net, channelwidth, [
                        3, 3], stride=2, scope=nscope, padding='SAME')
                net = tf.concat(axis=3, values=[net1, net2])
                channelwidth *= 2
            elif 'widener2' in celltype:
                for input_conn in celltype["inputs"]:
                    reduced_inputs = nets[input_conn]
                    while(reduced_inputs.shape[1] != input_dim[1]):
                        reduced_inputs = slim.max_pool2d(reduced_inputs, [2,2], padding='SAME')
                    all_inputs.append(reduced_inputs)
                net = tf.concat(axis=3, values=all_inputs)
                num_channels = int(input_dim[3])
                nscope='Widener_'+str(cellnumber)+'_MaxPool_2x2'
                print("Initial #channels={}, after skip={}".format(num_channels, int(net.shape[3])))
                net = slim.max_pool2d(net, [2,2], scope=nscope, padding='SAME')
                channelwidth *= 2
            elif 'outputs' in celltype:
                pass
            else:
                print("Error: Invalid cell definition" + str(celltype))
                exit(-1)

            nets.append(net)
            cellnumber += 1
        return net

    def calc_l2norm(self,tensor):
        return tf.norm(tensor, ord=2)

    def gen_network(self,inputs, log_stats=False, is_training=True,arch=None):
            print(arch)
            return self.get_macro_net(inputs, log_stats, is_training,arch)


    def add_init(self,inputs, arch, is_training):
        init = cell_init.Init(0)
        net = init.cell(inputs, arch, is_training)
        return net

    def add_net(self,net, log_stats, is_training,
                    arch=None):

        net = self.gen_network(net, log_stats, is_training,
                               arch)
        return net

    def add_classification(self,net, arch, is_training):
        classification = cell_classification.Classification()
        logits = classification.cell(net, arch, is_training)
        return logits


    def gen_amlanet(self,
                inputs,
                arch=None,
                initcell=None,
                classificationcell=None,
                log_stats=False,
                is_training=True,
                scope='Nacnet'):
        net = self.add_init(inputs, initcell, is_training)
        end_points = {}
        net = self.add_net(net, log_stats, is_training,arch)
        linear_softmax = self.add_classification(
            net, classificationcell, is_training)


        # return logits, end_points
        return linear_softmax

    def loss(self, logits, labels):
        """Add L2Loss to all the trainable variables.
        Add summary for "Loss" and "Loss/avg".
        Args:
            logits: Logits from inference().
            labels: Labels from distorted_inputs or inputs(). 1-D tensor
                            of shape [batch_size]
        Returns:
            Loss tensor of type float.
        """
        # Calculate the average cross entropy loss across the batch.
        labels = tf.cast(labels, tf.int64)
        cross_entropy = tf.nn.sparse_softmax_cross_entropy_with_logits(
            labels=labels, logits=logits, name='cross_entropy_per_example')
        cross_entropy_mean = tf.reduce_mean(
            cross_entropy, name='cross_entropy')
        tf.add_to_collection('losses', cross_entropy_mean)

        # Get auxiliary loss, TODO remove hardcoded weight of 0.4
        aux_logits = tf.get_collection('auxiliary_loss')
        weight = 0.4
        for num, logits in enumerate(aux_logits):
            cross_entropy = tf.nn.sparse_softmax_cross_entropy_with_logits(
                labels=labels, logits=logits, name='aux_loss_{}'.format(num))
            cross_entropy_mean = weight * tf.reduce_mean(cross_entropy)
            tf.add_to_collection('losses', cross_entropy_mean)

        # The total loss is defined as the cross entropy loss plus all of the weight
        # decay terms (L2 loss).
        return tf.add_n(tf.get_collection('losses'), name='total_loss')

    def get_train_op(self, total_loss, global_step, child_training):
        """Train CIFAR-10 model.
        Create an optimizer and apply to all trainable variables. Add moving
        average for all trainable variables.
        Args:
            total_loss: Total loss from loss().
            global_step: Integer Variable counting the number of training steps
                processed.
        Returns:
            train_op: op for training.
        """
        learning_rate = self.get_learning_rate(global_step, child_training)
        tf.summary.scalar('learning_rate', learning_rate)

        if "regularization" in child_training.keys():
            total_loss = self.get_regularization_loss(total_loss, child_training)

        # Generate moving averages of all losses and associated summaries.
        loss_averages_op = self._add_loss_summaries(total_loss)

        # Compute gradients.
        with tf.control_dependencies([loss_averages_op]):
            opt = self.get_opt(learning_rate, child_training)
            grads = opt.compute_gradients(total_loss)

        grads = self.clip_gradients(grads, child_training)

        # Apply gradients.
        apply_gradient_op = opt.apply_gradients(grads, global_step=global_step)

        # Add histograms for trainable variables.
        for var in tf.trainable_variables():
            tf.summary.histogram(var.op.name, var)

        # Add histograms for gradients.
        for grad, var in grads:
            if grad is not None:
                tf.summary.histogram(var.op.name + '/gradients', grad)

        # Track the moving averages of all trainable variables.
        variable_averages = tf.train.ExponentialMovingAverage(
            MOVING_AVERAGE_DECAY, global_step)
        variables_averages_op = variable_averages.apply(
            tf.trainable_variables())

        with tf.control_dependencies([apply_gradient_op, variables_averages_op]):
            train_op = tf.no_op(name='train')

        return train_op

    def get_regularization_loss(self, total_loss, child_training):
        if child_training["regularization"]["type"] == "l2":
            # l2-regularization
            l2_reg = child_training["regularization"]["value"]
            tf_variables = [var for var in tf.trainable_variables()]
            l2_losses = []
            for var in tf_variables:
               l2_losses.append(tf.reduce_sum(var**2))
            l2_loss = tf.add_n(l2_losses)
            total_loss += l2_reg * l2_loss
        return total_loss

    def _add_loss_summaries(self, total_loss):
        """Add summaries for losses in CIFAR-10 model.
        Generates moving average for all losses and associated summaries for
        visualizing the performance of the network.
        Args:
            total_loss: Total loss from loss().
        Returns:
            loss_averages_op: op for generating moving averages of losses.
        """
        # Compute the moving average of all individual losses and the total
        # loss.
        loss_averages = tf.train.ExponentialMovingAverage(0.9, name='avg')
        losses = tf.get_collection('losses')
        loss_averages_op = loss_averages.apply(losses + [total_loss])

        # Attach a scalar summary to all individual losses and the total loss; do the
        # same for the averaged version of the losses.
        for loss in losses + [total_loss]:
            # Name each loss as '(raw)' and name the moving average version of the loss
            # as the original loss name.
            tf.summary.scalar(loss.op.name + ' (raw)', loss)
            tf.summary.scalar(loss.op.name, loss_averages.average(loss))

        return loss_averages_op

    def get_opt(self, learning_rate, child_training):
        if "optimizer" not in child_training.keys() or child_training["optimizer"]["type"] == "sgd":
            opt = tf.train.GradientDescentOptimizer(learning_rate)
        elif child_training["optimizer"]["type"] == "rms":
            opt = tf.train.RMSPropOptimizer(lr, 0.9, 0.9, 1.0)
        elif child_training["optimizer"]["type"] == "momentum":
            momentum = child_training["optimizer"]["momentum"]
            opt = tf.train.MomentumOptimizer(learning_rate,
              momentum, use_locking=True, use_nesterov=True)

        return opt

    def clip_gradients(self, grads, child_training):
        if "gradient_clipping" in child_training.keys():
            if child_training["gradient_clipping"]["type"] == "norm":
                # Gradient clipping based on norm
                clipped = []
                grad_bound = child_training["gradient_clipping"]["value"]
                for grad, var in grads:
                   if isinstance(grad, tf.IndexedSlices):
                       c_g = tf.clip_by_norm(grad.values, grad_bound)
                       c_g = tf.IndexedSlices(grad.indices, c_g)
                   else:
                       c_g = tf.clip_by_norm(grad, grad_bound)
                   clipped.append((c_g, var))
                grads = clipped
        return grads

    def get_learning_rate(self, global_step, child_training):
        # Variables that affect learning rate.
        num_batches_per_epoch = NUM_EXAMPLES_PER_EPOCH_FOR_TRAIN // self.batch_size
        curr_epoch = global_step // num_batches_per_epoch

        if "lr" not in child_training.keys() or child_training["lr"]["type"] == "exponential_decay":
            initial_lr = child_training["lr"].get("initial", INITIAL_LEARNING_RATE)
            lr_decay = child_training["lr"].get("decay", LEARNING_RATE_DECAY_FACTOR)
            epochs_per_decay = child_training["lr"].get("epochs_per_decay", NUM_EPOCHS_PER_DECAY)

            decay_steps = int(num_batches_per_epoch * epochs_per_decay)
            learning_rate = tf.train.exponential_decay(initial_lr,
                                                       global_step,
                                                       decay_steps,
                                                       lr_decay,
                                                       staircase=True)            

        elif child_training["lr"]["type"] == "cosine_decay":
            curr_epoch = tf.to_int32(curr_epoch)
            last_reset = tf.Variable(0, dtype=tf.int32, trainable=False,
                                                 name="last_reset")
            lr_max = child_training["lr"]["max"]
            lr_min = child_training["lr"]["min"]
            lr_T_0 = child_training["lr"]["T_0"]
            lr_T_mul = child_training["lr"]["T_mul"]
            T_curr = curr_epoch - last_reset
            T_i = tf.Variable(lr_T_0, dtype=tf.int32,trainable=False, name="T_i")

            def _update():
                update_last_reset = tf.assign(last_reset, curr_epoch, use_locking=True)
                update_T_i = tf.assign(T_i, T_i * lr_T_mul, use_locking=True)
                with tf.control_dependencies([update_last_reset, update_T_i]):
                  rate = tf.to_float(T_curr) / tf.to_float(T_i) * 3.1415926
                  learning_rate = lr_min + 0.5 * (lr_max - lr_min) * (1.0 + tf.cos(rate))
                return learning_rate

            def _no_update():
                rate = tf.to_float(T_curr) / tf.to_float(T_i) * 3.1415926
                learning_rate = lr_min + 0.5 * (lr_max - lr_min) * (1.0 + tf.cos(rate))
                return learning_rate

            learning_rate = tf.cond(
                tf.greater_equal(T_curr, T_i), _update, _no_update)

        return learning_rate
        
    def tower_loss(self, scope, logits, labels):
        """Calculate the total loss on a single tower running the CIFAR model.
            Args:
            scope: unique prefix string identifying the CIFAR tower, e.g. 'tower_0'
            images: Images. 4D tensor of shape [batch_size, height, width, 3].
            labels: Labels. 1D tensor of shape [batch_size].
            Returns:
            Tensor of shape [] containing the total loss for a batch of data
        """
        # Build the portion of the Graph calculating the losses. Note that we will
        # assemble the total_loss using a custom function below.
        _ = self.loss(logits, labels)

        # Assemble all of the losses for the current tower only.
        losses = tf.get_collection('losses', scope)

        # Calculate the total loss for the current tower.
        total_loss = tf.add_n(losses, name='total_loss')

        # Attach a scalar summary to all individual losses and the total loss; do the
        # same for the averaged version of the losses.
        for loss in losses + [total_loss]:
            loss_name = re.sub('%s_[0-9]*/' % 'tower', '', loss.op.name)
            tf.summary.scalar(loss_name, loss)

        return total_loss

    def average_gradients(self, tower_grads):
        """Calculate the average gradient for each shared variable across all
           towers.
           Note that this function provides a synchronization point across all
           towers.
        Args:
            tower_grads: List of lists of (gradient, variable) tuples. The outer
            list is over individual gradients. The inner list is over the
            gradient calculation for each tower.
        Returns:
            List of pairs of (gradient, variable) where the gradient has been 
            average across all towers.
        """
        average_grads = []
        for grad_and_vars in zip(*tower_grads):
            grads = []
            for grad, _ in grad_and_vars:

                if True:
                    expanded_g = tf.expand_dims(grad, 0)
                    grads.append(expanded_g)
            """
            This break also to avoid adding empty list because of variance and mean calculation
            Average over the 'tower' dimension. """
            grad = tf.concat(axis=0, values=grads)
            grad = tf.reduce_mean(grad, 0)

            """ Keep in mind that the Variables are redundant because they are shared
            across towers. So .. we will just return the first tower's
            pointer to # the Variable. """
            var = grad_and_vars[0][1]
            grad_and_var = (grad, var)
            average_grads.append(grad_and_var)
        return average_grads
