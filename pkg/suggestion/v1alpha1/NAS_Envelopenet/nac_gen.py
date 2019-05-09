from pkg.suggestion.NAS_Envelopenet.generate_arch import NACAlg

class NAC:
    def __init__(self,
                envelopenet_params = dict()):
                initcell = {
    			"Layer0": {"Branch0": {"block": "conv2d", "kernel_size": [1, 1], "outputs": 64}},
   				"Layer2": {"Branch0": {"block": "lrn" }}
				}
                self.alg = NACAlg(envelopenet_params, initcell)
    
    def get_init_arch(self):
        return self.alg.generate()
    
    def get_arch(self, arch, result):
        return self.alg.construct(arch, result)
