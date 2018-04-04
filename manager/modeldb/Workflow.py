from modeldb.basic.Structs import (
    Model, ModelConfig, ModelMetrics, Dataset)
from modeldb.basic.ModelDbSyncerBase import Syncer
import sys
import json
import argparse

parser = argparse.ArgumentParser(description='model db tiny client')
parser.add_argument('request') 
parser.add_argument('-s', '--server', default="modeldb-backend") 
parser.add_argument('-p', '--port', default=6543) 
args = parser.parse_args()

req_j = args.request
req = json.loads(req_j)
owner = req["owner"]
study = req["study"]
train = req["train"]
modelpath = req["modelpath"]
hyp = req["hyperparameter"]
met = req["metrics"]

syncer_obj = Syncer.create_syncer(study, owner, "", host=args.server, port=args.port)
datasets = {
    "train": Dataset("", {}),
    "test": Dataset("", {}),
}

model = train
mdb_model = Model(study, model, modelpath)
model_config = ModelConfig("NN", hyp)
model_metrics = ModelMetrics(met)

syncer_obj.sync_datasets(datasets)
syncer_obj.sync_model("train", model_config, mdb_model)
syncer_obj.sync_metrics("test", mdb_model, model_metrics)
syncer_obj.sync()
