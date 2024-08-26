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
import torch.nn.functional as F
from operations import FactorizedReduce, MixedOp, StdConv


class Cell(nn.Module):
    """Cell for search
    Each edge is mixed and continuous relaxed.
    """

    def __init__(
        self,
        num_nodes,
        c_prev_prev,
        c_prev,
        c_cur,
        reduction_prev,
        reduction_cur,
        search_space,
    ):
        """
        Args:
            num_nodes: Number of intermediate cell nodes
            c_prev_prev: channels_out[k-2]
            c_prev : Channels_out[k-1]
            c_cur   : Channels_in[k] (current)
            reduction_prev: flag for whether the previous cell is reduction cell or not
            reduction_cur: flag for whether the current cell is reduction cell or not
        """

        super(Cell, self).__init__()
        self.reduction_cur = reduction_cur
        self.num_nodes = num_nodes

        # If previous cell is reduction cell, current input size does not match with
        # output size of cell[k-2]. So the output[k-2] should be reduced by preprocessing
        if reduction_prev:
            self.preprocess0 = FactorizedReduce(c_prev_prev, c_cur)
        else:
            self.preprocess0 = StdConv(
                c_prev_prev, c_cur, kernel_size=1, stride=1, padding=0
            )
        self.preprocess1 = StdConv(c_prev, c_cur, kernel_size=1, stride=1, padding=0)

        # Generate dag from mixed operations
        self.dag_ops = nn.ModuleList()

        for i in range(self.num_nodes):
            self.dag_ops.append(nn.ModuleList())
            # Include 2 input nodes
            for j in range(2 + i):
                # Reduction with stride = 2 must be only for the input node
                stride = 2 if reduction_cur and j < 2 else 1
                op = MixedOp(c_cur, stride, search_space)
                self.dag_ops[i].append(op)

    def forward(self, s0, s1, w_dag):
        s0 = self.preprocess0(s0)
        s1 = self.preprocess1(s1)

        states = [s0, s1]
        for edges, w_list in zip(self.dag_ops, w_dag):
            state_cur = sum(
                edges[i](s, w) for i, (s, w) in enumerate((zip(states, w_list)))
            )
            states.append(state_cur)

        state_out = torch.cat(states[2:], dim=1)
        return state_out


class NetworkCNN(nn.Module):

    def __init__(
        self,
        init_channels,
        input_channels,
        num_classes,
        num_layers,
        criterion,
        search_space,
        num_nodes,
        stem_multiplier,
    ):
        super(NetworkCNN, self).__init__()

        self.init_channels = init_channels
        self.num_classes = num_classes
        self.num_layers = num_layers
        self.criterion = criterion

        self.num_nodes = num_nodes
        self.stem_multiplier = stem_multiplier

        c_cur = self.stem_multiplier * self.init_channels

        self.stem = nn.Sequential(
            nn.Conv2d(input_channels, c_cur, 3, padding=1, bias=False),
            nn.BatchNorm2d(c_cur),
        )

        # In first Cell stem is used for s0 and s1
        # c_prev_prev and c_prev - output channels size
        # c_cur - init channels size
        c_prev_prev, c_prev, c_cur = c_cur, c_cur, self.init_channels

        self.cells = nn.ModuleList()

        reduction_prev = False
        for i in range(self.num_layers):
            # For Network with 1 layer: Only Normal Cell
            if self.num_layers == 1:
                reduction_cur = False
            else:
                # For Network with two layers: First layer - Normal, Second - Reduction
                # For Other Networks: [1/3, 2/3] Layers - Reduction cell with double channels
                # Others - Normal cell
                if (self.num_layers == 2 and i == 1) or (
                    self.num_layers > 2
                    and i in [self.num_layers // 3, 2 * self.num_layers // 3]
                ):
                    c_cur *= 2
                    reduction_cur = True
                else:
                    reduction_cur = False

            cell = Cell(
                self.num_nodes,
                c_prev_prev,
                c_prev,
                c_cur,
                reduction_prev,
                reduction_cur,
                search_space,
            )
            reduction_prev = reduction_cur
            self.cells.append(cell)

            c_cur_out = c_cur * self.num_nodes
            c_prev_prev, c_prev = c_prev, c_cur_out

        self.global_pooling = nn.AdaptiveAvgPool2d(1)
        self.classifier = nn.Linear(c_prev, self.num_classes)

        # Initialize alphas parameters
        num_ops = len(search_space.primitives)

        self.alpha_normal = nn.ParameterList()
        self.alpha_reduce = nn.ParameterList()

        for i in range(self.num_nodes):
            self.alpha_normal.append(nn.Parameter(1e-3 * torch.randn(i + 2, num_ops)))
            if self.num_layers > 1:
                self.alpha_reduce.append(
                    nn.Parameter(1e-3 * torch.randn(i + 2, num_ops))
                )

        # Setup alphas list
        self.alphas = []
        for name, parameter in self.named_parameters():
            if "alpha" in name:
                self.alphas.append((name, parameter))

    def forward(self, x):

        weights_normal = [F.softmax(alpha, dim=-1) for alpha in self.alpha_normal]
        weights_reduce = [F.softmax(alpha, dim=-1) for alpha in self.alpha_reduce]

        s0 = s1 = self.stem(x)

        for cell in self.cells:
            weights = weights_reduce if cell.reduction_cur else weights_normal
            s0, s1 = s1, cell(s0, s1, weights)

        out = self.global_pooling(s1)

        # Make out flatten
        out = out.view(out.size(0), -1)

        logits = self.classifier(out)
        return logits

    def print_alphas(self):

        print("\n>>> Alphas Normal <<<")
        for alpha in self.alpha_normal:
            print(F.softmax(alpha, dim=-1))

        if self.num_layers > 1:
            print("\n>>> Alpha Reduce <<<")
            for alpha in self.alpha_reduce:
                print(F.softmax(alpha, dim=-1))
        print("\n")

    def getWeights(self):
        return self.parameters()

    def getAlphas(self):
        for _, parameter in self.alphas:
            yield parameter

    def loss(self, x, y):
        logits = self.forward(x)
        return self.criterion(logits, y)

    def genotype(self, search_space):
        gene_normal = search_space.parse(self.alpha_normal, k=2)
        gene_reduce = search_space.parse(self.alpha_reduce, k=2)
        # concat all intermediate nodes
        concat = range(2, 2 + self.num_nodes)

        return search_space.genotype(
            normal=gene_normal,
            normal_concat=concat,
            reduce=gene_reduce,
            reduce_concat=concat,
        )
