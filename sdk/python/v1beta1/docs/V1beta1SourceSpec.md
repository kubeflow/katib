# V1beta1SourceSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**file_system_path** | [**V1beta1FileSystemPath**](V1beta1FileSystemPath.md) | During training model, metrics may be persisted into local file in source code, such as tfEvent use case | [optional] 
**filter** | [**V1beta1FilterSpec**](V1beta1FilterSpec.md) | Default metric output format is {\&quot;metric\&quot;: \&quot;&lt;metric_name&gt;\&quot;, \&quot;value\&quot;: &lt;int_or_float&gt;, \&quot;epoch\&quot;: &lt;int&gt;, \&quot;step\&quot;: &lt;int&gt;}, but if the output doesn&#39;t follow default format, please extend it here | [optional] 
**http_get** | [**V1HTTPGetAction**](https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/V1HTTPGetAction.md) | Model-train source code can expose metrics by http, such as HTTP endpoint in prometheus metric format | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


