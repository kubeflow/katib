import os
import argparse

import logging
logging.basicConfig(levle=logging.DEBUG)


a = '''
    lr = 0.01, exponentional decay

'''
def model_constructor():
    model = "SHIT"
    return model

def load_data():
    data = "MORE SHIT"
    return data

def get_data_iter():
    data = "ITER SHIT"
    return data

if __name__ == "__main__":
    
    parser = argparse.ArgumentParser(description="Training container",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument('--num-classes', type=int, default=10,
                        help='The number of classes')
    parser.add_argument('--num-examples', type=int, default=60000,
                        help='The number of training examples')
    parser.add_argument('--image_shape', default='1, 28, 28', help='shape of training images')

    args = parser.parse_args()