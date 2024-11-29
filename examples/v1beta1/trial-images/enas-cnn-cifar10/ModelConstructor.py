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

import json

from keras.layers import Dense, GlobalAveragePooling2D, Input
from keras.models import Model
from op_library import concat, conv, dw_conv, reduction, sp_conv


class ModelConstructor(object):
    def __init__(self, arc_json, nn_json):
        self.arch = json.loads(arc_json)
        nn_config = json.loads(nn_json)
        self.num_layers = nn_config["num_layers"]
        self.input_sizes = nn_config["input_sizes"]
        self.output_size = nn_config["output_sizes"][-1]
        self.embedding = nn_config["embedding"]

    def build_model(self):
        # a list of the data all layers
        all_layers = [0 for _ in range(self.num_layers + 1)]

        # ================= Stacking layers =================
        # Input Layer. Layer 0
        input_layer = Input(shape=self.input_sizes)
        all_layers[0] = input_layer

        # Intermediate Layers. Starting from layer 1.
        for l_index in range(1, self.num_layers + 1):
            input_layers = list()
            opt = self.arch[l_index - 1][0]
            opt_config = self.embedding[str(opt)]
            skip = self.arch[l_index - 1][1 : l_index + 1]

            # set up the connection to the previous layer first
            input_layers.append(all_layers[l_index - 1])

            # then add skip connections
            for i in range(l_index - 1):
                if l_index > 1 and skip[i] == 1:
                    input_layers.append(all_layers[i])

            layer_input = concat(input_layers)
            if opt_config["opt_type"] == "convolution":
                layer_output = conv(layer_input, opt_config)
            if opt_config["opt_type"] == "separable_convolution":
                layer_output = sp_conv(layer_input, opt_config)
            if opt_config["opt_type"] == "depthwise_convolution":
                layer_output = dw_conv(layer_input, opt_config)
            elif opt_config["opt_type"] == "reduction":
                layer_output = reduction(layer_input, opt_config)

            all_layers[l_index] = layer_output

        # Final Layer
        # Global Average Pooling, then Fully connected with softmax.
        avgpooled = GlobalAveragePooling2D()(all_layers[self.num_layers])

        # TODO (andreyvelich): Currently, Dropout layer fails in distributed training.
        # Error: creating distributed tf.Variable with aggregation=MEAN
        # and a non-floating dtype is not supported, please use a different aggregation or dtype
        # dropped = Dropout(0.4)(avgpooled)

        logits = Dense(units=self.output_size, activation="softmax")(avgpooled)

        # Encapsulate the model
        self.model = Model(inputs=input_layer, outputs=logits)

        return self.model
