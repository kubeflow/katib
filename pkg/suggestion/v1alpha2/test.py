#from . import parsing_utils
import parsing_utils
import api_pb2

class Feasible:
  def __init__(self, mi=[], ma=[] , li=[], step=[]):
    self.min =mi
    self.max = ma
    self.list =li
    self.step =step

class ParameterSpec:
  def __init__(self, name, ptype, feasibleSpace):
    self.name = name
    self.parameter_type = ptype
    self.feasible = feasibleSpace

def main():
  print("Hello World!")
  request_num =3
  f = Feasible(0.01,0.05)
  p =ParameterSpec("--lr", api_pb2.DOUBLE, f)
  f1 = Feasible(0.01,0.05)
  p1 =ParameterSpec("--lr1", api_pb2.DOUBLE, f1)
  f2 = Feasible(0.01,0.05)
  p2 =ParameterSpec("--lr2", api_pb2.DOUBLE, f2)
  f3 = Feasible(1,9)
  p3 =ParameterSpec("--lr3", api_pb2.INT, f3)
  f4 = Feasible(10,50,"",2)
  p4 =ParameterSpec("--lr4", api_pb2.INT, f4)
  f5 = Feasible("","",['first','second','third'],"")
  p5 =ParameterSpec("--lr5", api_pb2.CATEGORICAL, f5)
  parameters = [ p,p1,p2,p3,p4,p5]
  parameter_config = parsing_utils.parse_parameter_configs(parameters)
  trial_specs =[]
  for _ in range(request_num):
      sample = parameter_config.random_sample()
      suggestion = parsing_utils.parse_x_next_vector(sample,
                            parameter_config.parameter_types,
                            parameter_config.names,
                            parameter_config.discrete_info,
                            parameter_config.categorical_info)
      trial_spec = api_pb2.TrialSpec()
      trial_spec.experiment_name = "abc"
      for param in suggestion:
        trial_spec.parameter_assignments.assignments.add(name = param['name'], value = str(param['value']))
      trial_specs.append(trial_spec)
  reply = api_pb2.GetSuggestionsReply()
  for trial_spec in trial_specs:
      reply.trials.add(spec=trial_spec)
  print(reply)
  
if __name__== "__main__":
  main()
