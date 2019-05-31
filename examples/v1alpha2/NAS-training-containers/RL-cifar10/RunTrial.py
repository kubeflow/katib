import keras
import numpy as np
from keras.datasets import cifar10
from ModelConstructor import ModelConstructor
from keras.utils import to_categorical
from keras.utils import multi_gpu_model
from keras.preprocessing.image import ImageDataGenerator
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
    parser.add_argument('--num_gpus', type=int, default=1, metavar='N',
                        help='number of GPU that used for training')
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

    num_gpus = args.num_gpus
    print(">>> num_gpus received by trial:")
    print(num_gpus)

    print("\n>>> Constructing Model...")
    constructor = ModelConstructor(arch, nn_config)
    test_model = constructor.build_model()
    print(">>> Model Constructed Successfully\n")

    if num_gpus > 1:
        test_model = multi_gpu_model(test_model, gpus=num_gpus)
    
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

    augmentation = ImageDataGenerator(
        width_shift_range=0.1,
        height_shift_range=0.1,
        horizontal_flip=True)
    
    aug_data_flow = augmentation.flow(x_train, y_train, batch_size=128)

    print(">>> Data Loaded. Training starts.")
    for e in range(num_epochs):
        print("\nTotal Epoch {}/{}".format(e+1, num_epochs))
        history = test_model.fit_generator(generator=aug_data_flow,
                                           steps_per_epoch=int(len(x_train)/128)+1,
                                           epochs=1, verbose=1,
                                           validation_data=(x_test, y_test))
        print("Training-Accuracy={}".format(history.history['acc'][-1]))
        print("Training-Loss={}".format(history.history['loss'][-1]))
        print("Validation-Accuracy={}".format(history.history['val_acc'][-1]))
        print("Validation-Loss={}".format(history.history['val_loss'][-1]))
