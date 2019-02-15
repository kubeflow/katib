#!/usr/bin/env python

import os
import torch
import torch.distributed as dist
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim

from math import ceil
from random import Random
from torch.autograd import Variable
from torchvision import datasets, transforms

# Set hyperparameters
BATCH_SIZE = 64
NUM_EPOCHS = 3
LEARNING_RATE = 0.01
MOMENTUM = 0.9

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
    """ Partitions a dataset into different chunks. """

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
    def __init__(self):
        super(Net, self).__init__()
        self.conv1 = nn.Conv2d(3, 6, 5)
        self.pool = nn.MaxPool2d(2, 2)
        self.conv2 = nn.Conv2d(6, 16, 5)
        self.fc1 = nn.Linear(16 * 5 * 5, 120)
        self.fc2 = nn.Linear(120, 84)
        self.fc3 = nn.Linear(84, 10)

    def forward(self, x):
        x = self.pool(F.relu(self.conv1(x)))
        x = self.pool(F.relu(self.conv2(x)))
        x = x.view(-1, 16 * 5 * 5)
        x = F.relu(self.fc1(x))
        x = F.relu(self.fc2(x))
        x = self.fc3(x)
        return x


def partition_dataset():
    """ Partitioning CIFAR10. """
    train_set_full = datasets.CIFAR10(
        '/data',
        train=True,
        download=True,
        transform=transforms.Compose([
            transforms.RandomHorizontalFlip(),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.4914, 0.4822, 0.4465],
                                 std=[0.2023, 0.1994, 0.2010]),
        ]))
    test_set_full = datasets.CIFAR10(
        '/data',
        train=False,
        transform=transforms.Compose([
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.4914, 0.4822, 0.4465],
                                 std=[0.2023, 0.1994, 0.2010]),
        ]))

    size = dist.get_world_size()
    bsz = int(BATCH_SIZE / float(size))
    partition_sizes = [1.0 / size for _ in range(size)]
    partition = DataPartitioner(train_set_full, partition_sizes)
    partition = partition.use(dist.get_rank())
    train_set = torch.utils.data.DataLoader(
            partition, batch_size=bsz, shuffle=True)
    test_set = torch.utils.data.DataLoader(test_set_full, batch_size=BATCH_SIZE, shuffle=False)
    return train_set, test_set, bsz


def average_gradients(model):
    """ Gradient averaging. """
    size = float(dist.get_world_size())
    for param in model.parameters():
        dist.all_reduce(param.grad.data, op=dist.reduce_op.SUM, group=0)
        param.grad.data /= size


def run():
    """ Distributed SGD Example. """
    rank = dist.get_rank()
    torch.manual_seed(1234)
    train_set, test_set, bsz = partition_dataset()
    model = Net()
    criterion = nn.CrossEntropyLoss()
    num_batches = ceil(len(train_set.dataset) / float(bsz))

    optimizer = optim.SGD(model.parameters(), lr=LEARNING_RATE, momentum=MOMENTUM)
    print('\nTraining for %i epochs with learning rate %.3f and momentum %.3f:' %
          (NUM_EPOCHS, LEARNING_RATE, MOMENTUM))

    for epoch in range(NUM_EPOCHS):
        running_loss = 0.0
        model.train()
        for i, data in enumerate(train_set, 0):
            inputs, labels = data
            inputs, labels = Variable(inputs), Variable(labels)
            optimizer.zero_grad()
            outputs = model(inputs)
            loss = criterion(outputs, labels)
            loss.backward()
            average_gradients(model)
            optimizer.step()

            # Print metrics
            running_loss += loss.item()
            if i % 100 == 99:  # Print every 100 mini-batches
                print('[Rank %i, Epoch %d, Batch %d] training loss: %.3f' %
                        (rank, epoch+1, i+1, running_loss/100))
                running_loss = 0.0

        if dist.get_rank() == 0:
            correct = 0
            total = 0
            model.eval()
            for data in test_set:
                images, labels = data
                outputs = model(Variable(images))
                _, predicted = torch.max(outputs.data, 1)
                total += labels.size(0)
                correct += (predicted == labels).sum()

            print('Validation set accuracy after epoch %i: %d %%' % (
                epoch+1, 100 * correct / total))


def init_processes(fn, backend='tcp'):
    """ Initialize the distributed environment. """
    dist.init_process_group(backend)
    fn()


if __name__ == "__main__":
    init_processes(run)
