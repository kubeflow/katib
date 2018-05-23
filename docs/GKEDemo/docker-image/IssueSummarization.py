"""Generates predictions using a stored model.

Uses trained model files to generate a prediction.
"""

from __future__ import print_function

import numpy as np
import dill as dpickle
from keras.models import load_model
from seq2seq_utils import Seq2Seq_Inference

class IssueSummarization(object):

  def __init__(self):
    with open('body_pp.dpkl', 'rb') as body_file:
      body_pp = dpickle.load(body_file)
    with open('title_pp.dpkl', 'rb') as title_file:
      title_pp = dpickle.load(title_file)
    self.model = Seq2Seq_Inference(encoder_preprocessor=body_pp,
                                   decoder_preprocessor=title_pp,
                                   seq2seq_model=load_model('seq2seq_model_tutorial.h5'))

  def predict(self, input_text, feature_names): # pylint: disable=unused-argument
    return np.asarray([[self.model.generate_issue_title(body[0])[1]] for body in input_text])
