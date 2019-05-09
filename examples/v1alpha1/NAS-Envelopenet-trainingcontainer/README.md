# About the training container

The algorithm follows the idea proposed in *Fast Neural Architecture Construction using EnvelopeNets* by Kamath et al.(https://arxiv.org/pdf/1803.06744.pdf). It is not a Reinforcement Learning or evolution based NAS,
rather a method to construct deep network

# How this code works

Firstly the yaml file is parsed using Operation.py and suggestion_param.py. Then in nasenvelopenet_service.py suggestion, calls nac_gen.py to generate initial architecture. Then it passes this to run_trial.py.
run_trial.py is entrypoint. This is called from the suggestion. It invokes Model Constructor which constructs the model. There is a parameter in the algorithm which is max_iterations, is used as a maximum number of restructuring
iterations of the model. When this is reached, it evaluates the model. 
Based on this, suggestion calls generate_arch.py to improve the architecture from the metrics collected, and this loop runs till max_iterations.

Model Constructor uses net.py various methods to build the model, which itself uses cell_classification.py, cell_init.py and cell_main.py as a definition of the initial cell, the envelopecell and the classification cell used
to build the model. cifar10_input.py is used for various methods needed for the CIFAR-10 dataset. Evaluate.py has various methods for testing.

# How to run this code

I have attached a testing code test.py which I used to parse the yaml file and run this locally. But there have been changes in the code after I tested it on Katib. So you might need to change something.
