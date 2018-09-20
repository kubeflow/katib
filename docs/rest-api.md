## REST API 

For each RPC service, we define an equivalent HTTP REST API method using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). The mapping between service and REST API method can be found in this [file](https://github.com/kubeflow/katib/blob/master/pkg/api/api.proto). The mapping includes the URL path, query parameters and request body. You can read more details on the supported methods and binding options at this [link](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http)

## Examples
Assuming `{HOST}` as the ingress host for vizier-core-rest, 

### Create a new Study

Sample JSON file (study_config.json)

```json
{
	"name": "study1",
	"owner": "mayankjuneja",
	"optimization_type": 2,
	"optimization_goal": null,
	"parameter_configs": {
		"configs": [{
			"name": "dropout",
			"parameter_type": 1,
			"feasible": {
				"max": "0.5",
				"min": "0.1",
				"list": []
			}
		}, {
			"name": "activation",
			"parameter_type": 4,
			"feasible": {
				"max": "0",
				"min": "0",
				"list": ["tanh", "relu"]
			}
		}]
	},
	"objective_value_name": "accuracy"
}
```

Request:

```shell
curl -X POST {HOST}/api/Manager/CreateStudy
```

Response:

```json
{"study_id": "k350afb1ad0e580e"}
```

### Get List of Studies

Request:
```shell
curl -X GET {HOST}/api/Manager/GetStudyList
```

Response:
```json
{
  "study_overviews": [
    {
      "name": "study1",
      "owner": "mayankjuneja",
      "id": "k350afb1ad0e580e"
    },
    {
      "name": "study2",
      "owner": "mayankjuneja",
      "id": "pae3c470c9584cc4"
    }
  ]
}
```

