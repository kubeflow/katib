import os
import itertools
import datetime
import json
import sys


from pkg.suggestion.NASenvelopenet.generate import NACAlg

class NAC:
    def __init__(self,
                envelopenet_params = dict(),
                algorithm = "envelopenet",
                stages = 3):
                initcell = {
    			"Layer0": {"Branch0": {"block": "conv2d", "kernel_size": [1, 1], "outputs": 64}},
   				"Layer2": {"Branch0": {"block": "lrn" }}
				}
                self.alg = NACAlg(algorithm, envelopenet_params, stages, initcell)


    def __del__(self):
        pass
    
    def get_init_arch(self):

        return self.alg.generate()
    
    def get_arch(self, arch, result):
        return self.alg.construct(arch, result)

    


    
