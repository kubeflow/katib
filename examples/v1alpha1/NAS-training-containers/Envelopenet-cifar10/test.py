import yaml
import json
import os
import sys
import time
from nac_gen import NAC
from RunTrial import Trial
from Operation import SearchSpace
class EnvelopenetService:
    def __init__(self):
        self.manager_addr = "vizier-core"
        self.manager_port = 6789
        self.current_study_id = ""
        self.current_trial_id = ""
        self.ctrl_cache_file = ""
        self.is_first_run = True
        self.max_steps=5
        self.res=sys.argv[2]
    def generate_arch(self, doc):
        self._get_suggestion_param(doc)
        self._get_search_space(doc)

        self.generator = NAC(
            self.search_space,
            algorithm = self.suggestion_config["algorithm"],
            stages=self.stages)
    
    
    def GetSuggestion(self):
            #result=self.GetEvaluationResult()
        narch=self.generator.get_init_arch()
        i=0
        while i<self.max_steps:
            self.trial(narch, self.suggestion_config, i)

            narch=self.generator.get_arch(narch, self.res)
            print(narch)
            i+=1
            open(self.res,"w").close()

    
    def _get_search_space(self, doc):

        all_params = doc["spec"]["nasConfig"]["operations"]
        graph_config = list(all_params)
        #print(graph_config)
        for operation_dict in graph_config:
            opt_spec = list(operation_dict["parameterconfigs"])
            avail_space = dict()
            for ispec in opt_spec:
                if ispec["parametertype"]=="categorical":
                    spec_name = ispec["name"]
                    avail_space[spec_name] = list(ispec["feasible"]["list"])
                elif ispec["parametertype"]=="int":
                    spec_name = ispec["name"]
                    avail_space[spec_name] = int(ispec["feasible"]["value"])
        self.stages = 3
        self.input_size = 32
        self.output_size = 10
        #search_space_object = SearchSpace(search_space_raw)
        self.search_space = avail_space
        #print(self.search_space)
        
    def _get_suggestion_param(self, doc):  
            suggestion_d = dict()
            suggestion_d = doc["spec"]["suggestionSpec"]
            suggestion_list=list()
            suggestion_list = list(suggestion_d["suggestionParameters"])
            self.suggestion_config=dict()
            #print(suggestion_list)
            for attr in suggestion_list:
                   self.suggestion_config[attr["name"]]=attr["value"]
            #print(self.suggestion_config)
    
    def trial(self,arch,config,i):
        print(arch)
        arch=json.dumps(arch)
        config=json.dumps(config)
        arch=str(arch).replace('\"', '\'')
        config=str(config).replace('\"', '\'')
        cmd = "python RunTrial.py --architecture=\"{}\" --nn_config=\"{}\" --num_epochs=\"{}\" &>{}".format(arch, config, i, self.res)
        print(cmd)
        os.system(cmd)
