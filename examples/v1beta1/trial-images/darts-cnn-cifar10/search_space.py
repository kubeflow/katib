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

from collections import namedtuple

import torch


class SearchSpace:
    def __init__(self, search_space):
        self.primitives = search_space
        self.primitives.append("none")

        print(">>> All Primitives")
        print("{}\n".format(self.primitives))
        self.genotype = namedtuple(
            "Genotype", "normal normal_concat reduce reduce_concat"
        )

    def parse(self, alpha, k):
        """
        Parse continuous alpha to discrete gene.
        alpha is ParameterList:
        ParameterList [
            Parameter(n_edges1, n_ops),
            Parameter(n_edges2, n_ops),
            ...
        ]

        gene is list:
        [
            [('node1_ops_1', node_idx), ..., ('node1_ops_k', node_idx)],
            [('node2_ops_1', node_idx), ..., ('node2_ops_k', node_idx)],
            ...
        ]
        each node has two edges (k=2) in CNN.
        """

        gene = []
        assert self.primitives[-1] == "none"  # assume last PRIMITIVE is 'none'

        # 1) Convert the mixed op to discrete edge (single op) by choosing top-1 weight edge
        # 2) Choose top-k edges per node by edge score (top-1 weight in edge)
        for edges in alpha:
            # edges: Tensor(n_edges, n_ops)
            edge_max, primitive_indices = torch.topk(edges[:, :-1], 1)  # ignore 'none'
            topk_edge_values, topk_edge_indices = torch.topk(edge_max.view(-1), k)
            node_gene = []
            for edge_idx in topk_edge_indices:
                prim_idx = primitive_indices[edge_idx]
                prim = self.primitives[prim_idx]
                node_gene.append((prim, edge_idx.item()))

            gene.append(node_gene)

        return gene
