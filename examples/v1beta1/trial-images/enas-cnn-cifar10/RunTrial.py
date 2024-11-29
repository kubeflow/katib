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

import argparse

import tensorflow as tf
from keras.datasets import cifar10
from ModelConstructor import ModelConstructor
from tensorflow import keras
from tensorflow.keras.layers import RandomFlip, RandomTranslation, Rescaling
from tensorflow.keras.utils import to_categorical

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="TrainingContainer")
    parser.add_argument(
        "--architecture",
        type=str,
        default="",
        metavar="N",
        help="architecture of the neural network",
    )
    parser.add_argument(
        "--nn_config",
        type=str,
        default="",
        metavar="N",
        help="configurations and search space embeddings",
    )
    parser.add_argument(
        "--num_epochs",
        type=int,
        default=10,
        metavar="N",
        help="number of epoches that each child will be trained",
    )
    parser.add_argument(
        "--num_gpus",
        type=int,
        default=1,
        metavar="N",
        help="number of GPU that used for training",
    )
    args = parser.parse_args()

    arch = args.architecture.replace("'", '"')
    print(">>> arch received by trial")
    print(arch)

    nn_config = args.nn_config.replace("'", '"')
    print(">>> nn_config received by trial")
    print(nn_config)

    num_epochs = args.num_epochs
    print(">>> num_epochs received by trial")
    print(num_epochs)

    num_gpus = args.num_gpus
    print(">>> num_gpus received by trial:")
    print(num_gpus)

    print("\n>>> Constructing Model...")
    constructor = ModelConstructor(arch, nn_config)

    num_physical_gpus = len(tf.config.experimental.list_physical_devices("GPU"))
    if 1 <= num_gpus <= num_physical_gpus:
        devices = ["/gpu:" + str(i) for i in range(num_physical_gpus)]
    else:
        num_physical_cpu = len(tf.config.experimental.list_physical_devices("CPU"))
        devices = ["/cpu:" + str(j) for j in range(num_physical_cpu)]

    print(f">>> Using devices: {devices}")

    strategy = tf.distribute.MirroredStrategy(devices)
    with strategy.scope():
        print("Setup TensorFlow distributed training")
        test_model = constructor.build_model()
        test_model.summary()
        test_model.compile(
            loss=keras.losses.categorical_crossentropy,
            optimizer=keras.optimizers.Adam(learning_rate=1e-3),
            metrics=["accuracy"],
        )

    print(">>> Model Constructed Successfully\n")

    (x_train, y_train), (x_test, y_test) = cifar10.load_data()
    x_train = x_train.astype("float32")
    x_test = x_test.astype("float32")
    x_train /= 255
    x_test /= 255
    y_train = to_categorical(y_train)
    y_test = to_categorical(y_test)

    augmentation = tf.keras.Sequential(
        [
            Rescaling(1.0 / 255),
            RandomFlip("horizontal"),
            RandomTranslation(height_factor=0.1, width_factor=0.1),
        ]
    )

    train_dataset = tf.data.Dataset.from_tensor_slices((x_train, y_train))
    train_dataset = train_dataset.map(lambda x, y: (augmentation(x, training=True), y))
    # TODO: Add batch size to args
    train_dataset = train_dataset.batch(128)

    print(">>> Data Loaded. Training starts.")
    for e in range(num_epochs):
        print("\nTotal Epoch {}/{}".format(e + 1, num_epochs))
        history = test_model.fit(
            train_dataset,
            steps_per_epoch=int(len(x_train) / 128) + 1,
            epochs=1,
            verbose=1,
            validation_data=(x_test, y_test),
        )
        print("Training-Accuracy={}".format(history.history["accuracy"][-1]))
        print("Training-Loss={}".format(history.history["loss"][-1]))
        print("Validation-Accuracy={}".format(history.history["val_accuracy"][-1]))
        print("Validation-Loss={}".format(history.history["val_loss"][-1]))
