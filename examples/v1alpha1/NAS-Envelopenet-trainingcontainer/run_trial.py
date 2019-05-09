from model_constructor import ModelConstructor
import argparse
import json
        
if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='TrainingContainer')
    parser.add_argument('--architecture', type=str, default="", metavar='N',
                        help='architecture of the neural network')
    parser.add_argument('--parameters', type=str, default="", metavar='N',
                        help='configurations')
    parser.add_argument('--current_itr', type=str, default="0", metavar='N',
                        help='Current restructuring iteration')
    args = parser.parse_args()

    arch=args.architecture.replace("\'", "\"")
    arch=json.loads(arch)
    config=args.parameters.replace("\'", "\"")
    config=json.loads(config)
    current_itr=int(args.current_itr)
    constructor = ModelConstructor(arch,config,current_itr)
    constructor.build_model()

    max_iterations=config["iterations"]
    if(current_itr>max_iterations):
        constructor.evaluate()
