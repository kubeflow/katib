# V1alpha3SourceSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**file_system_path** | [**V1alpha3FileSystemPath**](V1alpha3FileSystemPath.md) | During training model, metrics may be persisted into local file in source code, such as tfEvent use case | [optional] 
**filter** | [**V1alpha3FilterSpec**](V1alpha3FilterSpec.md) | Default metric output format is {\&quot;metric\&quot;: \&quot;&lt;metric_name&gt;\&quot;, \&quot;value\&quot;: &lt;int_or_float&gt;, \&quot;epoch\&quot;: &lt;int&gt;, \&quot;step\&quot;: &lt;int&gt;}, but if the output doesn&#39;t follow default format, please extend it here | [optional] 
**http_get** | [**V1HTTPGetAction**](V1HTTPGetAction.md) | Model-train source code can expose metrics by http, such as HTTP endpoint in prometheus metric format | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


