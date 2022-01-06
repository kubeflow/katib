import numpy as np  
import pandas as pd  
import argparse
import os
import cv2
from tensorflow.keras.models import Sequential  
from tensorflow.keras.layers import Dense,Dropout,Flatten,Conv2D,MaxPooling2D,BatchNormalization
from tensorflow.keras import callbacks as keras_callbacks
from tensorflow.keras.callbacks import ModelCheckpoint
import time
import keras

dirdata1='/app/4class_2d/0_img/'
dirdata2='/app/4class_2d/1_img/'
dirdata3='/app/4class_2d/2_img/'
dirdata4='/app/4class_2d/3_img/'

class LossAndErrorPrintingCallback(keras_callbacks.Callback):
    def on_epoch_end(self, epoch, logs=None):
        print("epoch={}".format(epoch))
        print("Training-Accuracy={:7.6f}".format(logs["accuracy"]))
        print("Training-Loss={:7.6f}".format(logs["loss"]))
        print("Validation-Accuracy={:7.6f}".format(logs["val_accuracy"]))
        print("Validation-Loss={:7.6f}".format(logs["val_loss"]))

def data_img(dirdata1,dirdata2,dirdata3,dirdata4):
    data = [] 
    label = []
    height=120
    BLACK = [0,0,0]
    for i in os.listdir(dirdata1): 
        img = cv2.imread(dirdata1+i)
        x,y,z = img.shape
        if x<=height:
            constant = cv2.copyMakeBorder(img,0,height-x,0,0,cv2.BORDER_CONSTANT,value=BLACK)
        data.append(constant) 
        label.append(0) 
    for i in os.listdir(dirdata2):
        img = cv2.imread(dirdata2+i)
        x,y,z = img.shape
        if x<=height:
            constant = cv2.copyMakeBorder(img,0,height-x,0,0,cv2.BORDER_CONSTANT,value=BLACK)
        data.append(constant) 
        label.append(1) 
        
    for i in os.listdir(dirdata3):
        img = cv2.imread(dirdata3+i)
        x,y,z = img.shape
        if x<=height:
            constant = cv2.copyMakeBorder(img,0,height-x,0,0,cv2.BORDER_CONSTANT,value=BLACK)
        data.append(constant) 
        label.append(2) 
        
    for i in os.listdir(dirdata4):
        img = cv2.imread(dirdata4+i)
        x,y,z = img.shape
        if x<=height:
            constant = cv2.copyMakeBorder(img,0,height-x,0,0,cv2.BORDER_CONSTANT,value=BLACK)
        data.append(constant) 
        label.append(3) 
        
    data = np.array(data) 
    label = np.array(label) 

    label_one_hot = np.eye(4)[label] # 4 one-hot encoding

    data_norm = data / 255 
    return data_norm,label_one_hot

def run(args):

    data, label=data_img(dirdata1,dirdata2,dirdata3,dirdata4)
    
    from sklearn.model_selection import train_test_split
    X_train, X_test, Y_train, Y_test = train_test_split(data, label, test_size=0.5, random_state=42)
    
    from keras.applications import vgg16
    vgg16=vgg16.VGG16(weights='imagenet', include_top=False, input_shape=(120,50,3))
    vgg16.summary()
    
    model = Sequential()
    for layer in vgg16.layers:
        model.add(layer)

    for layer in model.layers:
        layer.trainable = False
    model.add(Flatten())
    model.add(Dense(128, activation='relu'))
    model.add(Dropout(0.25))
    model.add(Dense(4, activation='softmax'))

    model.compile(loss='categorical_crossentropy', optimizer=args.optimizer, metrics=['accuracy'])  
    
    checkpoint = ModelCheckpoint(filepath='./', monitor='val_accuracy', verbose=1, save_best_only=True)
    callbacks = [checkpoint, LossAndErrorPrintingCallback()]
    start_time = time.time()
 
    train_history = model.fit(X_train,Y_train,validation_split=0.2,epochs=args.epochs, batch_size=30, verbose=1, callbacks=callbacks)
    end_time = time.time()
  

if __name__ == '__main__':

    parser = argparse.ArgumentParser(description="train transfer learning",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument('--optimizer', type=str, default='adam', dest='optimizer',
                        help='optimizer: Adadelta/Adagrad/Adam/Adamax/Ftrl/SGD/RMSprop')
    parser.add_argument('--epochs', type=int, default=100, dest='epochs',
                        help='epoch: number of iterations run')



    args = parser.parse_args()
    run(args)


