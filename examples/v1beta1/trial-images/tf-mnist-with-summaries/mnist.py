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
import os

import tensorflow as tf
from tensorflow.keras import Model
from tensorflow.keras.layers import Conv2D, Dense, Flatten


class MyModel(Model):
    def __init__(self):
        super(MyModel, self).__init__()
        self.conv1 = Conv2D(32, 3, activation="relu")
        self.flatten = Flatten()
        self.d1 = Dense(128, activation="relu")
        self.d2 = Dense(10)

    def call(self, x):
        x = self.conv1(x)
        x = self.flatten(x)
        x = self.d1(x)
        return self.d2(x)


def train_step(
    args,
    model,
    optimizer,
    train_ds,
    epoch,
    loss_object,
    train_summary_writer,
    train_loss,
    train_accuracy,
):
    for step, (images, labels) in enumerate(train_ds):
        with tf.GradientTape() as tape:
            # training=True is only needed if there are layers with different
            # behavior during training versus inference (e.g. Dropout).
            predictions = model(images, training=True)
            loss = loss_object(labels, predictions)
            gradients = tape.gradient(loss, model.trainable_variables)
            optimizer.apply_gradients(zip(gradients, model.trainable_variables))

            train_loss(loss)
            train_accuracy(labels, predictions)

        if step % args.log_interval == 0:
            print(
                "Train Epoch: {} [{}/60000 ({:.0f}%)]\tloss={:.4f}, accuracy={:.4f}".format(
                    epoch + 1,
                    step * args.batch_size,
                    100.0 * step * args.batch_size / 60000,
                    train_loss.result(),
                    train_accuracy.result() * 100,
                )
            )

    with train_summary_writer.as_default():
        tf.summary.scalar("loss", train_loss.result(), step=epoch)
        tf.summary.scalar("accuracy", train_accuracy.result(), step=epoch)


def test_step(
    model, test_ds, epoch, loss_object, test_summary_writer, test_loss, test_accuracy
):
    for images, labels in test_ds:
        # training=False is only needed if there are layers with different
        # behavior during training versus inference (e.g. Dropout).
        predictions = model(images, training=False)
        t_loss = loss_object(labels, predictions)

        test_loss(t_loss)
        test_accuracy(labels, predictions)

    with test_summary_writer.as_default():
        tf.summary.scalar("loss", test_loss.result(), step=epoch)
        tf.summary.scalar("accuracy", test_accuracy.result(), step=epoch)

    print(
        "Test Loss: {:.4f}, Test Accuracy: {:.4f}\n".format(
            test_loss.result(), test_accuracy.result() * 100
        )
    )


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--batch-size",
        type=int,
        default=64,
        help="input batch size for training (default: 64)",
    )
    parser.add_argument(
        "--learning-rate",
        type=float,
        default=0.001,
        help="learning rate (default: 0.001)",
    )
    parser.add_argument(
        "--epochs",
        type=int,
        default=10,
        metavar="N",
        help="number of epochs to train (default: 10)",
    )
    parser.add_argument(
        "--log-interval",
        type=int,
        default=100,
        metavar="N",
        help="how many batches to wait before logging training status (default: 100)",
    )
    parser.add_argument(
        "--log-path",
        type=str,
        default=os.path.join(
            os.getenv("TEST_TMPDIR", "/tmp"),
            "tensorflow/mnist/logs/mnist_with_summaries",
        ),
        help="Summaries log PATH",
    )
    args = parser.parse_args()

    # Setup dataset
    mnist = tf.keras.datasets.mnist
    (x_train, y_train), (x_test, y_test) = mnist.load_data()
    x_train, x_test = x_train / 255.0, x_test / 255.0
    # Add a channels dimension
    x_train = x_train[..., tf.newaxis].astype("float32")
    x_test = x_test[..., tf.newaxis].astype("float32")
    train_ds = (
        tf.data.Dataset.from_tensor_slices((x_train, y_train))
        .shuffle(10000)
        .batch(args.batch_size)
    )
    test_ds = tf.data.Dataset.from_tensor_slices((x_test, y_test)).batch(
        args.batch_size
    )

    # Setup tensorflow summaries
    train_log_dir = os.path.join(args.log_path, "train")
    test_log_dir = os.path.join(args.log_path, "test")
    train_summary_writer = tf.summary.create_file_writer(train_log_dir)
    test_summary_writer = tf.summary.create_file_writer(test_log_dir)

    # Create an instance of the model
    model = MyModel()
    loss_object = tf.keras.losses.SparseCategoricalCrossentropy(from_logits=True)
    optimizer = tf.keras.optimizers.Adam(learning_rate=args.learning_rate)

    train_loss = tf.keras.metrics.Mean(name="train_loss")
    train_accuracy = tf.keras.metrics.SparseCategoricalAccuracy(name="train_accuracy")

    test_loss = tf.keras.metrics.Mean(name="test_loss")
    test_accuracy = tf.keras.metrics.SparseCategoricalAccuracy(name="test_accuracy")

    for epoch in range(args.epochs):
        # Reset the metrics at the start of the next epoch
        train_summary_writer.flush()
        test_summary_writer.flush()

        train_step(
            args,
            model,
            optimizer,
            train_ds,
            epoch,
            loss_object,
            train_summary_writer,
            train_loss,
            train_accuracy,
        )
        test_step(
            model,
            test_ds,
            epoch,
            loss_object,
            test_summary_writer,
            test_loss,
            test_accuracy,
        )


if __name__ == "__main__":
    main()
