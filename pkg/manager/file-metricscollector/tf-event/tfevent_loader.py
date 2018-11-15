import tensorflow as tf
import os
from datetime import datetime
import rfc3339
import grpc
import api_pb2
import api_pb2_grpc
import sys
from logging import getLogger, StreamHandler, INFO
class TFEventFileParser:
    def find_all_files(self, directory):
        for root, dirs, files in os.walk(directory):
            yield root
            for f in files:
                yield os.path.join(root, f)
# TFEventFileParser parse tfevent and get metrics you specified.
# When the event file under a directory(e.g. test dir), please specify "{{dirname}}/{{metrics name}}"
# For example, in the TensorFlow official tutorial for mnist with summary (https://github.com/tensorflow/tensorflow/blob/master/tensorflow/examples/tutorials/mnist/mnist_with_summaries.py), it will save "accracy" metrics under "train" dir and "test" dir.
# Then in Katib, please specify name of metrics "train/accuracy" and "test/accuracy".

    def parse_summary(self, tfefile, metrics):
        metrics_log = {}
        for m in metrics:
            metrics_log[m] = api_pb2.MetricsLog(name=m,values=[])
        for summary in tf.train.summary_iterator(tfefile):
            paths=tfefile.split("/")
            for v in summary.summary.value:
                for m in metrics:
                    tag = str(v.tag)
                    if len(paths) >= 2 and len(m.split("/")) >= 2:
                            tag = str(paths[-2]+"/"+v.tag)
                    if tag.startswith(m):
                        mv = metrics_log[m].values.add()
                        mv.time=rfc3339.rfc3339(datetime.fromtimestamp(summary.wall_time))
                        mv.value=str(v.simple_value)
        return metrics_log

class MetricsCollector:
    def __init__(self, manager_addr, manager_port, study_id, worker_id):
        self.logger = getLogger(__name__)
        handler = StreamHandler()
        handler.setLevel(INFO)
        self.logger.setLevel(INFO)
        self.logger.addHandler(handler)
        self.logger.propagate = False
        self.manager_addr = manager_addr
        self.study_id  = study_id
        self.worker_id = worker_id
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, manager_port)

        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=study_id), 10)
            self.metrics = gsrep.study_config.metrics
        self.parser = TFEventFileParser()

    def parse_file(self, directory):
        mls = []
        for f in self.parser.find_all_files(directory):
            if os.path.isdir(f):
                continue
            try:
                self.logger.info(f+" will be parsed.")
                ml = self.parser.parse_summary(f, self.metrics)
                for m in ml:
                    mls.append(ml[m])
            except:
                self.logger.warning("Unexpected error:"+ str(sys.exc_info()[0]))
                continue
        return mls 
