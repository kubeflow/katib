import tensorflow as tf
import os
from datetime import datetime
import rfc3339
import grpc
import api_pb2
import api_pb2_grpc
import sys
from logging import getLogger, StreamHandler, INFO

# TFEventFileParser parses tfevent files and returns an ObservationLog of the metrics specified.
# When the event file is under a directory(e.g. test dir), please specify "{{dirname}}/{{metrics name}}"
# For example, in the TensorFlow official tutorial for mnist with summary (https://github.com/tensorflow/tensorflow/blob/master/tensorflow/examples/tutorials/mnist/mnist_with_summaries.py),
# the "accuracy" metric is saved under "train" and "test" directories. So in Katib, please specify name of metrics as "train/accuracy" and "test/accuracy".
class TFEventFileParser:
    def find_all_files(self, directory):
        for root, dirs, files in os.walk(directory):
            yield root
            for f in files:
                yield os.path.join(root, f)

    def parse_summary(self, tfefile, metrics):
        metric_logs = []
        for summary in tf.train.summary_iterator(tfefile):
            paths=tfefile.split("/")
            for v in summary.summary.value:
                for m in metrics:
                    tag = str(v.tag)
                    if len(paths) >= 2 and len(m.split("/")) >= 2:
                            tag = str(paths[-2]+"/" + v.tag)
                    if tag.startswith(m):
                        ml = api_pb2.MetricLog(
                                time_stamp=rfc3339.rfc3339(datetime.fromtimestamp(summary.wall_time)),
                                metric=api_pb2.Metric(
                                    name=m,
                                    value=str(v.simple_value)
                                )
                            )
                        metric_logs.append(ml)
        return metric_logs

class MetricsCollector:
    def __init__(self, metric_names):
        self.logger = getLogger(__name__)
        handler = StreamHandler()
        handler.setLevel(INFO)
        self.logger.setLevel(INFO)
        self.logger.addHandler(handler)
        self.logger.propagate = False
        self.metrics = metric_names
        self.parser = TFEventFileParser()

    def parse_file(self, directory):
        mls = []
        for f in self.parser.find_all_files(directory):
            if os.path.isdir(f):
                continue
            try:
                self.logger.info(f + " will be parsed.")
                mls = self.parser.parse_summary(f, self.metrics)
            except Exception, e:
                self.logger.warning("Unexpected error: "+ str(e))
                continue
        print(str(mls))
        return api_pb2.ObservationLog(metric_logs=mls)
 
