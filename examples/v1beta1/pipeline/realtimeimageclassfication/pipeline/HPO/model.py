from __future__ import absolute_import, division, print_function, unicode_literals
import tensorflow as tf
import numpy as np
import logging
from datetime import datetime
import os
import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers
from tensorflow.keras.models import Sequential
import pathlib
logger = tf.get_logger()
logging.basicConfig(
    format="%(asctime)s %(levelname)-8s %(message)s",
    datefmt="%Y-%m-%dT%H:%M:%SZ",
    level=logging.INFO)
print('Tensorflow-version: {0}'.format(tf.__version__))

import os
import argparse
import json



    

# build model
def build_model(learning_rate,data_augmentation,opt):
    img_height = 180
    img_width = 180
    num_classes=5

    if opt == "SGD":
        optimizer = tf.keras.optimizers.SGD(learning_rate=learning_rate)
    elif opt =="ADAM":
        optimizer = tf.keras.optimizers.Adam(learning_rate=learning_rate),
    else:
        optimizer=tf.keras.optimizers.RMSprop(learning_rate=learning_rate)

    data_augmentation = keras.Sequential(
            [
                layers.experimental.preprocessing.RandomFlip("horizontal", 
                                                            input_shape=(img_height, 
                                                                        img_width,
                                                                        3)),
                layers.experimental.preprocessing.RandomRotation(0.1),
                layers.experimental.preprocessing.RandomZoom(0.1),
            ]
            )
    model = Sequential([
                data_augmentation,
                layers.experimental.preprocessing.Rescaling(1./255),
                layers.Conv2D(16, 3, padding='same', activation='relu'),
                layers.MaxPooling2D(),
                layers.Conv2D(32, 3, padding='same', activation='relu'),
                layers.MaxPooling2D(),
                layers.Conv2D(64, 3, padding='same', activation='relu'),
                layers.MaxPooling2D(),
                layers.Dropout(0.2),
                layers.Flatten(),
                layers.Dense(128, activation='relu'),
                layers.Dense(num_classes)
                ])

    model.compile(optimizer=optimizer,
              loss=tf.keras.losses.SparseCategoricalCrossentropy(from_logits=True),
              metrics=['accuracy'])
    

      
   
    return model

    
# callbacks
def get_callbacks():
    # callbacks 
    # checkpoint directory
    checkpointdir = '/tmp/model-ckpt'

    class customLog(tf.keras.callbacks.Callback):
        def on_epoch_end(self, epoch, logs={}):
            logging.info('epoch: {:.4f}'.format(epoch + 1))
            logging.info('loss={:.4f}'.format(logs['loss']))
            logging.info('accuracy={:.4f}'.format(logs['accuracy']))
            logging.info('val_accuracy={:.4f}'.format(logs['val_accuracy']))
    #logging.info("{{metricName: accuracy, metricValue: {:.4f}}};{{metricName: loss, metricValue: {:.4f}}}\n".format())

    callbacks = [
        #tf.keras.callbacks.TensorBoard(logdir),
        tf.keras.callbacks.ModelCheckpoint(filepath=checkpointdir),
        customLog()
    ]
    return callbacks


# parse arguments
def parse_arguments():
    parser = argparse.ArgumentParser()
  
    parser.add_argument('--log-path',
                        type=str,
                        default="",
                        help='The number of training steps to perform.')
    parser.add_argument('--learning-rate',
                        type=float,
                        default=0.001,
                        help='Learning rate for training.')
    parser.add_argument('--drop-out',
                        type=float,
                        default=0.2,
                        help='Drop out rate for training.')
    parser.add_argument('--optimizer',
                        type=str,
                        default="sgd",
                        help='optimizer for training.')

    args = parser.parse_known_args()[0]
    return args


def main():

    # parse arguments
    args = parse_arguments()

   
    if args.log_path == "":
        logging.basicConfig(
            format="%(asctime)s %(levelname)-8s %(message)s",
            datefmt="%Y-%m-%dT%H:%M:%SZ",
            level=logging.DEBUG)
    else:
        logging.basicConfig(
            format="%(asctime)s %(levelname)-8s %(message)s",
            datefmt="%Y-%m-%dT%H:%M:%SZ",
            level=logging.DEBUG,
            filename=args.log_path)
    

    dataset_url = "https://storage.googleapis.com/download.tensorflow.org/example_images/flower_photos.tgz"
    data_dir = tf.keras.utils.get_file('flower_photos', origin=dataset_url, untar=True)
    data_dir = pathlib.Path(data_dir)

    image_count = len(list(data_dir.glob('*/*.jpg')))
    batch_size = 32
    img_height = 180
    img_width = 180

    train_ds = tf.keras.preprocessing.image_dataset_from_directory(
                data_dir,
                validation_split=0.2,
                subset="training",
                seed=123,
                image_size=(img_height, img_width),
                batch_size=batch_size)


    val_ds = tf.keras.preprocessing.image_dataset_from_directory(
                        data_dir,
                        validation_split=0.2,
                        subset="validation",
                        seed=123,
                        image_size=(img_height, img_width),
                        batch_size=batch_size)
    class_names = train_ds.class_names
    print(class_names)


    AUTOTUNE = tf.data.AUTOTUNE

    train_ds = train_ds.cache().shuffle(3000).prefetch(buffer_size=AUTOTUNE)
    val_ds = val_ds.cache().prefetch(buffer_size=AUTOTUNE)

    normalization_layer = layers.experimental.preprocessing.Rescaling(1./255)
    normalized_ds = train_ds.map(lambda x, y: (normalization_layer(x), y))
    image_batch, labels_batch = next(iter(normalized_ds))
    first_image = image_batch[0]
    
    # Notice the pixels values are now in `[0,1]`.
    print(np.min(first_image), np.max(first_image)) 
    num_classes = 5
   

     # build and compile model
    learning_rate = float(args.learning_rate)
    logging.info("learning rate : {0}".format(learning_rate))
    model = build_model(learning_rate,float(args.drop_out),args.optimizer)


  
    history=model.fit(train_ds, 
            epochs=15,
            #steps_per_epoch=TF_STEPS_PER_EPOCHS, 
            validation_data=val_ds,
            validation_steps=1,
            callbacks=get_callbacks())

    logging.info("Training completed.")
    # successful completion
    exit(0)
  

if __name__ == "__main__":
    main()