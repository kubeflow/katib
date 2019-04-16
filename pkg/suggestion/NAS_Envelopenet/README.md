# About the Neural Architecture Search with Envelopenet Suggestion

The algorithm follows the idea proposed in *Fast Neural Architecture Construction using EnvelopeNets* by Kamath et al. It is not a Reinforcement Learning or evolution based NAS, rather a method to construct deep network
architectures by pruning and expansion of a base network. This approach directly compares the utility of different filters using statistics derived from filter featuremaps reach a state where the utility of different filters
within a network can be compared and hence can be used to construct networks. 

# Envelopenets

The EnvelopeCell is a set of M convolution blocks connected in parallel. E.g. one of the EnvelopeCells used in this work has 6 convolution blocks connected in parallel: 1x1 convolution,
3x3 convolution, 3x3 separable convolution, 5x5 convolution, 5x5 separable convolution and 7x7 separable convolution. The EnvelopeNet consists of a number of the EnvelopeCells stacked in series organized into stages of
n layers, separated by wideners.

# Output of `GetSuggestion()`

The output of `GetSuggestion()` is the `architecture`

`architecture` is a json string of the definition of a neural architecture. The format is as stated above. One example is:
```
{'type': 'macro', 'network': [{'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 26, 'inputs': []}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 26, 'inputs': []}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 64, 'inputs': [1]}, {'widener': {}}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 53, 'inputs': [3, 2, 1]}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 53, 'inputs': [4, 3, 2, 1]}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 128, 'inputs': [5, 4, 3, 2, 1]}, {'widener': {}}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 106, 'inputs': [7, 6, 5, 4, 3, 2, 1]}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 106, 'inputs': [8, 7, 6, 5, 4, 3, 2, 1]}, {'filters': ['1', '3', '5', '3sep', '5sep', '7sep'], 'outputs': 256, 'inputs': [9, 8, 7, 6, 5, 4, 3, 2, 1]}]}
```
# Flow of program

The parameters and initial architecture from yaml file is passed to ModelConstructor, which constructs the model and trains. Then it sends the featuremap statistics from this truncated training and refines the architecture. This process
continues till the max_iterations parameter.
