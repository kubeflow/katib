from __future__ import absolute_import, division, print_function, unicode_literals
import click
import os
import logging
import shutil
import ast
import numpy as np
import pathlib
import json
import tensorflow as tf
from storage import Storage
from tensorflow import keras
from tensorflow.keras import layers
from tensorflow.keras.models import Sequential
from tensorflow.python.keras.backend import dropout


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





@click.command()
@click.option('--tensorboard-gcs-logs', default='gs://kubeflowusecases/imageclassfication/logs')
@click.option('--gcs-path', default="gs://kubeflowusecases/imageclassfication/model")
@click.option('--artifacts', default="")
@click.option('--mode', default="local")
def train_model(tensorboard_gcs_logs,gcs_path,artifacts,mode):

    model_output_base_path='/mnt/saved_model'
    tensorboard_logs='/mnt/logs/'

    
    artifacts=ast.literal_eval(artifacts)

    learning_rate=artifacts['learning-rate']
    optimizer=artifacts['optimizer']
    dropout=artifacts['drop-out']

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
    learning_rate = float(learning_rate)
    logging.info("learning rate : {0}".format(learning_rate))
    model = build_model(learning_rate,float(dropout),optimizer)


    tensorboard_callback = tf.keras.callbacks.TensorBoard(log_dir=tensorboard_logs, histogram_freq=1)

    logging.info("Training starting...")

    history=model.fit(train_ds, 
            epochs=15,
            validation_data=val_ds,
            validation_steps=1,
            callbacks=[tensorboard_callback])

    model.save(model_output_base_path)

    Storage.upload(tensorboard_logs,tensorboard_gcs_logs)
        
    metadata = {
                'outputs': [{
                        'type': 'tensorboard',
                        'source': tensorboard_gcs_logs,        
                }]
                }
    with open("/mlpipeline-ui-metadata.json", 'w') as f:
                json.dump(metadata,f)

    if mode!= 'local':
        print("uploading to {0}".format(gcs_path))
        Storage.upload(model_output_base_path,gcs_path)

    else:
        print("Model will not be uploaded")
        pass

    
    
    return


if __name__ == "__main__":
    train_model()