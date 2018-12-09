{{ define "studyjobgenscript" }}
    <script type="text/javascript">
        var app = new Vue({
            el: '#app',
            data: {
                StudyName: {{.StudyName}},
                Owner: {{.Owner}},
                OptType: {{.OptimizationType}},
                OptGoal: {{.OptimizationGoal}},
                RequestCount: 0,
                NameSpace: 'kubeflow',
                ObjName: {{.ObjectiveValueName}},
                Metrics: {{.Metrics}},
                StudyJobNameOverwrite: '',
				ParameterConfigs: [
                    {{- range .ParamConf}}
                    {
                        Name: "{{.Name}}",
                        Type: "{{.Type}}",
                        Min: {{.Min}},
                        Max: {{.Max}},
                        List: '{{.List}}',
                    },
                    {{- end}}
                ],
                WorkerTemplates: {
                    {{- range $k,$v := .WorkerTemplates}}
                        {{$k}}: {{$v}},
                    {{- end}}
                },
                WorkerTemplateName: "scratch",
                WorkerTemplateScratch: "",
                {{ if eq .SuggestionAlgorithm "" }}
                SuggestAlgoSelect: "random",
                {{else}}
                SuggestAlgoSelect: "{{.SuggestionAlgorithm}}",
                {{end}}
                SuggestAlgoCustom: "",
                SuggestAlgoTmp: "",
                SuggestionReqNum: {{.RequestNumber}},
                SuggestionParameters: [
                    {{- range $k, $v := .SuggestionParams}}
                    {
                        NameDefault: "{{$k}}",
                        ValueDefault: "{{$v}}",
                        Name: "",
                        Value: "",
                    },
                    {{- end}}
                ],
                created: false,
            },
            computed: {
				MetricsList: function() {
                    var ml = []
					if (this.Metrics.length > 0) {
						ml = this.Metrics.split(" ");
					}
                    return ml
                    
				},
                StudyJobName: function() {
                    if (this.StudyJobNameOverwrite.length > 0) {
                        return this.StudyJobNameOverwrite
                    }
                    return this.StudyName + "-job"
                },
                WorkerTemplateValue: function() {
                    if (this.WorkerTemplateName == "scratch"){
                        return this.WorkerTemplateScratch;
                    }else{
                        return this.WorkerTemplates[this.WorkerTemplateName];
                    }
                },
                SuggestionAlgo: function() {
                    if (this.SuggestAlgoSelect != "custom"){
                        return this.SuggestAlgoSelect
                    }else{
                        return this.SuggestAlgoCustom
                    }
                },
                SuggestionParametersFormated: function() {
					var sparamConfigs = [];
					for (var i = 0; i < this.SuggestionParameters.length; i++) {
                        var name = ""
                        if (this.SuggestionParameters[i].Name != "") {
                            name = this.SuggestionParameters[i].Name
                        }else if (this.SuggestionParameters[i].NameDefault != 'Input Parameter Name') {
                            name = this.SuggestionParameters[i].NameDefault
                        }
                        if (name != "") {
                            var value = ""
                            if (this.SuggestionParameters[i].Value != "") {
                                value = this.SuggestionParameters[i].Value
                            }else{
                                value = this.SuggestionParameters[i].ValueDefault
                            }
                            sparamConfigs.push({
                                "name": name,
                                "value": value,
                            });
                        }
                    }
                    return sparamConfigs
                },
				StudyJobYaml: function() {
					var paramConfigs = [];
					for (var i = 0; i < this.ParameterConfigs.length; i++) {
                        if ( this.ParameterConfigs[i].Type === 'int' || this.ParameterConfigs[i].Type === 'double'){
                            paramConfigs.push({
                                'name': this.ParameterConfigs[i].Name,
                                'parametertype': this.ParameterConfigs[i].Type,
                                'feasible': {
                                    'max': this.ParameterConfigs[i].Max,
                                    'min': this.ParameterConfigs[i].Min,
                                },
                            });
                        }else{
                            paramConfigs.push({
                                'name': this.ParameterConfigs[i].Name,
                                'parametertype': this.ParameterConfigs[i].Type,
                                'feasible': {
                                    'list': this.ParameterConfigs[i].List.split(" "),
                                },
                            });
                        }
					}
					var studyjobObj = {
                        'apiVersion': 'kubeflow.org/v1alpha1',
                        'kind': 'StudyJob',
                        'metadata': {
                            'namespace': 'kubeflow',
                            'labels': {
                                'controller-tools.k8s.io': "1.0",
                            },
							'name': this.StudyJobName,
						},
                       
						'spec': {
                            'studyName': this.StudyName,
                            'owner': this.Owner,
                            'optimizationtype': this.OptType,
                            'objectivevaluename': this.ObjName,
                            'optimizationgoal': parseFloat(this.OptGoal),
                            'metricsnames': this.MetricsList,
                            'parameterconfigs': paramConfigs,
                            'requestcount': parseInt(this.RequestCount),
                            'suggestionSpec': {
                                'suggestionAlgorithm': this.SuggestionAlgo,
                                'requestNumber': parseInt(this.SuggestionReqNum),
                                'suggestionParameters': this.SuggestionParametersFormated,
                            },
                            'workerspec': {
                                'goTemplate': {
                                    'rawTemplate': this.WorkerTemplateValue,
                                },
                            },
						} // end spec
					} // end studyjobObj
					return jsyaml.dump(studyjobObj);
				}
            },
			methods: {
				addParameterConfig: function() {
					var newParam = {
                        Name: 'NewParameter',
                        Type: 'int',
                        Max: '0',
                        Min:'0',
                        List: '',
					};
					this.ParameterConfigs.push(newParam);
				},
				deleteParam: function(param) {
					var paramIndex = this.ParameterConfigs.indexOf(param);
					if (paramIndex > -1) {
						this.ParameterConfigs.splice(paramIndex, 1);
					}
				},
                addSuggestionParameter: function() {
                    var newParam = {
                        Name: '',
                        Value: '',
                        NameDefault: 'Input Parameter Name',
                        ValueDefault: 'Value',
                    };
                    this.SuggestionParameters.push(newParam);
                },
				deleteSuggestionParameter: function(sparam) {
					var sparamIndex = this.SuggestionParameters.indexOf(sparam);
					if (sparamIndex > -1) {
						this.SuggestionParameters.splice(sparamIndex, 1);
					}
				},
				SetSuggestionParameterDefault: function(sparam) {
                    this.SuggestionParameters = []
                    if (this.SuggestAlgoSelect == "grid"){
                        this.SuggestionParameters = [
                            {
                                Name: '',
                                Value: '',
                                NameDefault: 'DefaultGrid',
                                ValueDefault: '1',
                            },
                        ];
                    }else if (this.SuggestAlgoSelect == "bayesianoptimization"){
                        this.SuggestionParameters = [
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "N",
                                ValueDefault: "100",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "model_type",
                                ValueDefault: "gp",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "max_features",
                                ValueDefault: "auto",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "length_scale",
                                ValueDefault: "0.5",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "noise",
                                ValueDefault: "0.0005",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "nu",
                                ValueDefault: "1.5",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "kernel_type",
                                ValueDefault: "matern",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "n_estimators",
                                ValueDefault: "50",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "mode",
                                ValueDefault: "pi",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "trade_off",
                                ValueDefault: "0.01",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "burn_in",
                                ValueDefault: "10",
                            },
                        ];
                    }else if (this.SuggestAlgoSelect == "hyperband"){
                        this.SuggestionParameters = [
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "eta",
                                ValueDefault: "3",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "r_l",
                                ValueDefault: "float",
                            },
                            {
                                Name: '',
                                Value: '',
                                NameDefault: "ResourceName",
                                ValueDefault: "string",
                            },
                        ];
                    }
				},
                unsetCreated: function() {
                    this.created =  false
                },
                uploadStudyJob: function() {
                    var xmlHttpRequest = new XMLHttpRequest();
                    xmlHttpRequest.onreadystatechange = function()
                    {
                        var READYSTATE_COMPLETED = 4;
                        var HTTP_STATUS_OK = 200;
                        if( this.readyState == READYSTATE_COMPLETED
                         && this.status != HTTP_STATUS_OK )
                        {
                            alert( "Fail to upload templates "+this.responseText );
                        }
                    }
                    xmlHttpRequest.open( 'POST', '/katib/studyjob' );
                    xmlHttpRequest.setRequestHeader( 'Content-Type', 'application/x-www-form-urlencoded' );
                    xmlHttpRequest.send( "StudyJobManifest=" + encodeURIComponent(this.StudyJobYaml) )
                    if (xmlHttpRequest.status != 200){
                        this.created = true;
                        setTimeout(this.unsetCreated, 6000);
                    }
                },
			}
        });
    </script>
{{ end }}
