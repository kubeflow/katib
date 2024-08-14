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

import torch
import torch.nn as nn

OPS = {
    "none": lambda channels, stride: Zero(stride),
    "avg_pooling_3x3": lambda channels, stride: PoolBN(
        "avg", channels, kernel_size=3, stride=stride, padding=1
    ),
    "max_pooling_3x3": lambda channels, stride: PoolBN(
        "max", channels, kernel_size=3, stride=stride, padding=1
    ),
    "skip_connection": lambda channels, stride: (
        Identity() if stride == 1 else FactorizedReduce(channels, channels)
    ),
    "separable_convolution_3x3": lambda channels, stride: SepConv(
        channels, kernel_size=3, stride=stride, padding=1
    ),
    "separable_convolution_5x5": lambda channels, stride: SepConv(
        channels, kernel_size=5, stride=stride, padding=2
    ),
    # 3x3 -> 5x5
    "dilated_convolution_3x3": lambda channels, stride: DilConv(
        channels, kernel_size=3, stride=stride, padding=2, dilation=2
    ),
    # 5x5 -> 9x9
    "dilated_convolution_5x5": lambda channels, stride: DilConv(
        channels, kernel_size=5, stride=stride, padding=4, dilation=2
    ),
}


class Zero(nn.Module):
    """
    Zero operation
    """

    def __init__(self, stride):
        super(Zero, self).__init__()
        self.stride = stride

    def forward(self, x):
        if self.stride == 1:
            return x * 0.0
        # Resize by stride
        return x[:, :, :: self.stride, :: self.stride] * 0.0


class PoolBN(nn.Module):
    """
    Avg or Max pooling - BN
    """

    def __init__(self, pool_type, channels, kernel_size, stride, padding):
        super(PoolBN, self).__init__()
        if pool_type == "avg":
            self.pool = nn.AvgPool2d(
                kernel_size, stride, padding, count_include_pad=False
            )
        elif pool_type == "max":
            self.pool = nn.MaxPool2d(kernel_size, stride, padding)

        self.bn = nn.BatchNorm2d(channels, affine=False)
        self.net = nn.Sequential(self.pool, self.bn)

    def forward(self, x):
        # out = self.pool(x),
        # print(out)
        # out = self.bn(out)
        # print(out)
        return self.net(x)


class Identity(nn.Module):

    def __init__(self):
        super(Identity, self).__init__()

    def forward(self, x):
        return x


class FactorizedReduce(nn.Module):
    """
    Reduce feature map size by factorized pointwise (stride=2)
    ReLU - Conv1 - Conv2 - BN
    """

    def __init__(self, c_in, c_out):
        super(FactorizedReduce, self).__init__()
        self.relu = nn.ReLU()
        self.conv1 = nn.Conv2d(
            c_in, c_out // 2, kernel_size=1, stride=2, padding=0, bias=False
        )
        self.conv2 = nn.Conv2d(
            c_in, c_out // 2, kernel_size=1, stride=2, padding=0, bias=False
        )
        self.bn = nn.BatchNorm2d(c_out, affine=False)

    def forward(self, x):

        x = self.relu(x)
        out = torch.cat([self.conv1(x), self.conv2(x[:, :, 1:, 1:])], dim=1)
        out = self.bn(out)

        return out


class StdConv(nn.Module):
    """Standard convolition
    ReLU - Conv - BN
    """

    def __init__(self, c_in, c_out, kernel_size, stride, padding):
        super(StdConv, self).__init__()
        self.net = nn.Sequential(
            nn.ReLU(),
            nn.Conv2d(
                c_in,
                c_out,
                kernel_size=kernel_size,
                stride=stride,
                padding=padding,
                bias=False,
            ),
            nn.BatchNorm2d(c_out, affine=False),
        )

    def forward(self, x):
        return self.net(x)


class DilConv(nn.Module):
    """(Dilated) depthwise separable conv
    ReLU - (Dilated) depthwise separable - Pointwise - BN

    If dilation == 2, 3x3 conv => 5x5 receptive field
                      5x5 conv => 9x9 receptive field
    """

    def __init__(self, channels, kernel_size, stride, padding, dilation):
        super(DilConv, self).__init__()

        self.net = nn.Sequential(
            nn.ReLU(),
            nn.Conv2d(
                channels,
                channels,
                kernel_size,
                stride,
                padding,
                dilation=dilation,
                groups=channels,
                bias=False,
            ),
            nn.Conv2d(
                channels, channels, kernel_size=1, stride=1, padding=0, bias=False
            ),
            nn.BatchNorm2d(channels, affine=False),
        )

    def forward(self, x):
        return self.net(x)


class SepConv(nn.Module):
    """Depthwise separable conv
    DilConv (dilation=1) * 2
    """

    def __init__(self, channels, kernel_size, stride, padding):
        super(SepConv, self).__init__()
        self.net = nn.Sequential(
            DilConv(channels, kernel_size, stride=stride, padding=padding, dilation=1),
            DilConv(channels, kernel_size, stride=1, padding=padding, dilation=1),
        )

    def forward(self, x):
        return self.net(x)


class MixedOp(nn.Module):
    """Mixed operation"""

    def __init__(self, channels, stride, search_space):
        super(MixedOp, self).__init__()
        self.ops = nn.ModuleList()

        for primitive in search_space.primitives:
            op = OPS[primitive](channels, stride)
            self.ops.append(op)

    def forward(self, x, weights):
        """
        Args:
            x: input
            weights: weight for each operation
        """
        return sum(w * op(x) for w, op in zip(weights, self.ops))
