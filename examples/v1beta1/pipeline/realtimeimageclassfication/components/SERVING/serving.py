import kfserving
from typing import Dict
import os
import cv2
import numpy as np
import tensorflow as tf
import sys
import json
from io import BytesIO
import base64
import logging
from PIL import Image
import glob
import boto3
from tensorflow import keras


session = boto3.Session()
client = session.client('s3', endpoint_url='http://minio-service:9000', aws_access_key_id='minio', aws_secret_access_key='minio123')


class KFServingSampleModel(kfserving.KFModel):

    def __init__(self, name: str):
        super().__init__(name)
        self.name = name
        self.ready = False
        self.model_path=os.environ["STORAGE_URI"]
       
    
  

    def load(self):
        self.model = tf.keras.models.load_model(self.model_path)
        self.ready = True 
    
    def predict(self, request: Dict) -> Dict:
            
            imagepath=  r"source.jpg"

            if "instances" in request:
                self._prediction_type="request"
                logging.info("KF Serving with Curl request")
                inputs = request['instances'][0]
                originalimage = base64.b64decode(inputs)
                jpg_as_np = np.frombuffer(originalimage, dtype=np.uint8)
                image = cv2.imdecode(jpg_as_np, flags=1)
                logging.info(f"Image transform {image}")
                cv2.imwrite(imagepath,image)
            

            if "EventName" in request:
                logging.info("KF Serving with Minio")
                self._prediction_type="minio"
                # through minio event
                if request['EventName'] == 's3:ObjectCreated:Put':
                    bucket = request['Records'][0]['s3']['bucket']['name']
                    logging.info(f"Bucket name is {bucket}")
                    key = request['Records'][0]['s3']['object']['key']
                    self._key = key
                    self._bucket = bucket
                    logging.info(f"key name is {key}")
                    client.download_file(bucket, key, '/tmp/' + key)
                    image_nm='/tmp/' + key
                    image = cv2.imread(image_nm)
                    cv2.imwrite(imagepath,image)

            class_names= ['daisy', 'dandelion', 'roses', 'sunflowers', 'tulips']
            img_height = 180
            img_width = 180
            img = keras.preprocessing.image.load_img(imagepath, target_size=(img_height, img_width))
            img_array = keras.preprocessing.image.img_to_array(img)
            img_array = tf.expand_dims(img_array, 0) 
            predictions = self.model.predict(img_array)
            score = tf.nn.softmax(predictions[0])
            logging.info("This image most likely belongs to {} with a {:.2f} percent confidence.".format(class_names[np.argmax(score)], 100 * np.max(score)))
                        
            prediction_class=class_names[np.argmax(score)]


            image = cv2.imread(imagepath,cv2.IMREAD_UNCHANGED)

            position = (10,50)
            cv2.putText(
                image, #numpy array on which text is written
                prediction_class, #text
                position, #position at which writing has to start
                cv2.FONT_HERSHEY_SIMPLEX, #font family
                1, #font size
                (209, 80, 0, 255), #font color
                3) #font stroke
            output_file_name = '/tmp/' + r"Result.jpg"

            cv2.imwrite(output_file_name, image)
            
            finaloutput=[]
            with open(output_file_name, "rb") as image_file:
                encoded_bytes = base64.b64encode(image_file.read())
                encoded_string = encoded_bytes.decode('utf-8')
                finaloutput.append(encoded_string)


            if self._prediction_type == "minio":
                num= np.random.randint(0,10000)
                result ="String"
                con =result + "_" + str(num)
                filepath = "type-{0}/{1}".format(con, self._key)
                client.upload_file(output_file_name, 'imageprediction', filepath)

            elif self._prediction_type == "request":
                logging.info("Return prediction request")
            
            else:
                pass

            return {"predictions": finaloutput,"classname": prediction_class}



if __name__ == "__main__":
    model = KFServingSampleModel("kfserving-vision-realtime")
    model.load()
    kfserving.KFServer(workers=1).start([model])
        
