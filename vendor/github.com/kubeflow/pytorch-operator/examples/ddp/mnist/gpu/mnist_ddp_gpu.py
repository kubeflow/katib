#!/usr/bin/env python

import os
import sys
import torch
import torch.utils.data
import torch.utils.data.distributed
import torch.distributed as dist
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
from torch._utils import _flatten_dense_tensors, _unflatten_dense_tensors
from torch.nn.modules import Module

from math import ceil
from random import Random
from torch.autograd import Variable
from torchvision import datasets, transforms

import datetime

gbatch_size = 128

class DistributedDataParallel(Module):
    def __init__(self, module):
        super(DistributedDataParallel, self).__init__()
        self.module = module
        self.first_call = True

        def allreduce_params():
            if (self.needs_reduction):
                self.needs_reduction = False
                buckets = {}
                for param in self.module.parameters():
                    if param.requires_grad and param.grad is not None:
                        tp = type(param.data)
                        if tp not in buckets:
                            buckets[tp] = []
                        buckets[tp].append(param)
                for tp in buckets:
                    bucket = buckets[tp]
                    grads = [param.grad.data for param in bucket]
                    coalesced = _flatten_dense_tensors(grads)
                    dist.all_reduce(coalesced)
                    coalesced /= dist.get_world_size()
                    for buf, synced in zip(grads, _unflatten_dense_tensors(coalesced, grads)):
                        buf.copy_(synced)
        for param in list(self.module.parameters()):
            def allreduce_hook(*unused):
                Variable._execution_engine.queue_callback(allreduce_params)

            if param.requires_grad:
                param.register_hook(allreduce_hook)

    def weight_broadcast(self):
        for param in self.module.parameters():
            dist.broadcast(param.data, 0)

    def forward(self, *inputs, **kwargs):
        if self.first_call:
            print("first broadcast start")
            self.weight_broadcast()
            self.first_call = False
            print("first broadcast done")
        self.needs_reduction = True
        return self.module(*inputs, **kwargs)

class Partition(object):
    """ Dataset-like object, but only access a subset of it. """

    def __init__(self, data, index):
        self.data = data
        self.index = index

    def __len__(self):
        return len(self.index)

    def __getitem__(self, index):
        data_idx = self.index[index]
        return self.data[data_idx]


class DataPartitioner(object):
    """ Partitions a dataset into different chuncks. """

    def __init__(self, data, sizes=[0.7, 0.2, 0.1], seed=1234):
        self.data = data
        self.partitions = []
        rng = Random()
        rng.seed(seed)
        data_len = len(data)
        indexes = [x for x in range(0, data_len)]
        rng.shuffle(indexes)

        for frac in sizes:
            part_len = int(frac * data_len)
            self.partitions.append(indexes[0:part_len])
            indexes = indexes[part_len:]

    def use(self, partition):
        return Partition(self.data, self.partitions[partition])


class Net(nn.Module):
    """ Network architecture. """

    def __init__(self):
        super(Net, self).__init__()
        self.conv1 = nn.Conv2d(1, 10, kernel_size=5)
        self.conv2 = nn.Conv2d(10, 20, kernel_size=5)
        self.conv2_drop = nn.Dropout2d()
        self.fc1 = nn.Linear(320, 50)
        self.fc2 = nn.Linear(50, 10)

    def forward(self, x):
        x = F.relu(F.max_pool2d(self.conv1(x), 2))
        x = F.relu(F.max_pool2d(self.conv2_drop(self.conv2(x)), 2))
        x = x.view(-1, 320)
        x = F.relu(self.fc1(x))
        x = F.dropout(x, training=self.training)
        x = self.fc2(x)
        return F.log_softmax(x)


def partition_dataset(rank):
    """ Partitioning MNIST """
    dataset = datasets.MNIST(
        './data{}'.format(rank),
        train=True,
        download=True,
        transform=transforms.Compose([
            transforms.ToTensor(),
            transforms.Normalize((0.1307, ), (0.3081, ))
        ]))
    size = dist.get_world_size()
    bsz = int(gbatch_size / float(size))
    train_sampler = torch.utils.data.distributed.DistributedSampler(dataset)
    train_set = torch.utils.data.DataLoader(
        dataset, batch_size=bsz, shuffle=(train_sampler is None), sampler=train_sampler)
    return train_set, bsz

def average_gradients(model):
    """ Gradient averaging. """
    size = float(dist.get_world_size())
    for param in model.parameters():
        dist.all_reduce(param.grad.data, op=dist.reduce_op.SUM, group=0)
        param.grad.data /= size


def run(rank, size):
    """ Distributed Synchronous SGD Example """
    torch.manual_seed(1234)
    train_set, bsz = partition_dataset(rank)
    model = Net()
    model = model.cuda()
    model = DistributedDataParallel(model)
    optimizer = optim.SGD(model.parameters(), lr=0.01, momentum=0.5)

    num_batches = ceil(len(train_set.dataset) / float(bsz))
    print("num_batches = ", num_batches)
    time_start = datetime.datetime.now()
    for epoch in range(3):
        epoch_loss = 0.0
        for data, target in train_set:
            data, target = Variable(data).cuda(), Variable(target).cuda()
            optimizer.zero_grad()
            output = model(data)
            loss = F.nll_loss(output, target)
            epoch_loss += loss.item()
            loss.backward()
            average_gradients(model)
            optimizer.step()
        print('Epoch {} Loss {:.6f} Global batch size {} on {} ranks'.format(
                  epoch, epoch_loss / num_batches, gbatch_size, dist.get_world_size()))
    print("training time=",datetime.datetime.now() - time_start)

def init_print(rank, size, debug_print=True):
    if not debug_print:
        """ In case run on hundreds of nodes, you may want to mute all the nodes except master """
        if rank > 0:
            sys.stdout = open(os.devnull, 'w')
            sys.stderr = open(os.devnull, 'w')
    else:
        # labelled print with info of [rank/size]
        old_out = sys.stdout
        class LabeledStdout:
            def __init__(self, rank, size):
                self._r = rank
                self._s = size
                self.flush = sys.stdout.flush

            def write(self, x):
                if x == '\n':
                    old_out.write(x)
                else:
                    old_out.write('[%d/%d] %s' % (self._r, self._s, x))

        sys.stdout = LabeledStdout(rank, size)

if __name__ == "__main__":
    print("\n======= CUDA INFO =======")
    print("CUDA Availibility:", torch.cuda.is_available())
    if(torch.cuda.is_available()):
        print("CUDA Device Name:", torch.cuda.get_device_name(0))
        print("CUDA Version:", torch.version.cuda)
    print("=========================\n")
    dist.init_process_group(backend='mpi')
    size = dist.get_world_size()
    rank = dist.get_rank()
    init_print(rank, size)
    run(rank, size)
