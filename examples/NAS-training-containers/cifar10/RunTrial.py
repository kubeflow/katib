import keras
import numpy as np
from keras.datasets import cifar10
from ModelConstructor import ModelConstructor
from keras.utils import to_categorical
import argparse
import time

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='TrainingContainer')
    parser.add_argument('--architecture', type=str, default="", metavar='N',
                        help='architecture of the neural network')
    parser.add_argument('--nn_config', type=str, default="", metavar='N',
                        help='configurations and search space embeddings')
    parser.add_argument('--num_epochs', type=int, default=10, metavar='N',
                        help='number of epoches that each child will be trained')
    args = parser.parse_args()

    arch = args.architecture.replace("\'", "\"")
    print(">>> arch received by trial")
    print(arch)

    nn_config = args.nn_config.replace("\'", "\"")
    print(">>> nn_config received by trial")
    print(nn_config)

    num_epochs = args.num_epochs
    print(">>> num_epochs received by trial")
    print(num_epochs)

    print(">>> Constructing Model...")
    constructor = ModelConstructor(arch, nn_config)
    test_model = constructor.build_model()
    print(">>> Model Constructed Successfully")

    test_model.summary()
    test_model.compile(loss=keras.losses.categorical_crossentropy,
                       optimizer=keras.optimizers.Adam(lr=1e-3, decay=1e-4),
                       metrics=['accuracy'])

    (x_train, y_train), (x_test, y_test) = cifar10.load_data()
    x_train = x_train.astype('float32')
    x_test = x_test.astype('float32')
    x_train /= 255
    x_test /= 255
    y_train = to_categorical(y_train)
    y_test = to_categorical(y_test)

    print(">>> Data Loaded. Training start.")
    for e in range(num_epochs):
        print("\nTotal Epoch {}/{}".format(e+1, num_epochs))
        history = test_model.fit(x=x_train, y=y_train,
                                 shuffle=True, batch_size=128,
                                 epochs=1, verbose=1,
                                 validation_data=(x_test, y_test))
        print("Training-Accuracy={}".format(history.history['acc'][-1]))
        print("Training-Loss={}".format(history.history['loss'][-1]))
        print("Validation-Accuracy={}".format(history.history['val_acc'][-1]))
        print("Validation-Loss={}".format(history.history['val_loss'][-1]))
