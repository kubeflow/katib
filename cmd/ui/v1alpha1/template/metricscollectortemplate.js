{{ define "metricscollectortemplatescript" }}
    <script type="text/javascript">
		Vue.component('modal', {
		  template: '#modal-template'
		})
        var app = new Vue({
            el: '#app',
            data: {
                MetricsCollectorTemplates: [
                    {{- range $k,$v := .MetricsCollectorTemplate}}
                    {
                        Name: '{{$k}}',
                        Value: '{{$v}}',
                    },
                    {{- end}}
                ],
                TmpName: '',
                EditTemplate: {
                    Name: '',
                    Value: '',
                },
                edit: false,
                edited: false,
            },
            mounted: function() {
                editor = CodeMirror.fromTextArea(document.getElementById("code"), {
                    lineNumbers: true,
                    autoRefresh: true,
                });
            },
            computed: {
            },
			methods: {
				addMetricsCollectorTemplate: function() {
                    var newMCT = {
                        Name: 'newMetricsCollectorTemplate',
                        Value: '',
                    };
                    this.MetricsCollectorTemplates.push(newMCT);
				},
				copyMetricsCollectorTemplate: function(mt) {
                    var newMCT = {
                        Name: 'copy-'+mt.Name,
                        Value: mt.Value,
                    };
                    this.MetricsCollectorTemplates.push(newMCT);
                    this.edited = true;
                },
				deleteMetricsCollectorTemplate: function(mt) {
					var mtIndex = this.MetricsCollectorTemplates.indexOf(mt);
					if (mtIndex > -1) {
						this.MetricsCollectorTemplates.splice(mtIndex, 1);
					}
                    this.edited = true;
				},
                openEditor: function(wt) {
                    this.EditTemplate = wt;
                    this.TmpName = wt.Name;
                    this.edit = true;
                    editor.getDoc().setValue(wt.Value);
                    editor.refresh();
                },
                saveModal() {
                    this.EditTemplate.Value = editor.getDoc().getValue();
                    this.EditTemplate.Name = this.TmpName;
                    this.edit = false;
                    this.edited = true;
                },
                cancelModal() {
                    editor.getDoc().setValue(this.EditTemplate.Value);
                    this.edit = false;
                },
                uploadMetricsCollectorTemplates(){
                    var params = [];
                    for (var i = 0; i < this.MetricsCollectorTemplates.length; i++){
                        var mt = this.MetricsCollectorTemplates[i]
                        var param = encodeURIComponent( mt.Name ) + '=' + encodeURIComponent( wt.Value );
                        params.push( param );
                    }
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
                    xmlHttpRequest.open( 'POST', '/katib/metricscollectortemplate' );
                    xmlHttpRequest.setRequestHeader( 'Content-Type', 'application/x-www-form-urlencoded' );
                    xmlHttpRequest.send( params.join( '&' ).replace( /%20/g, '+' ) );
                    this.edited = false;
                }
			}
        });
    </script>
{{ end }}
