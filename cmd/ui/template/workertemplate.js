{{ define "workertemplatescript" }}
    <script type="text/javascript">
		Vue.component('modal', {
		  template: '#modal-template'
		})
        var app = new Vue({
            el: '#app',
            data: {
                WorkerTemplates: [
                    {{- range $k,$v := .WorkerTemplate}}
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
				addWorkerTemplate: function() {
                    var newWT = {
                        Name: 'newWorkerTemplate',
                        Value: '',
                    };
                    this.WorkerTemplates.push(newWT);
				},
				copyWorkerTemplate: function(wt) {
                    var newWT = {
                        Name: 'copy-'+wt.Name,
                        Value: wt.Value,
                    };
                    this.WorkerTemplates.push(newWT);
                },
				deleteWorkerTemplate: function(wt) {
					var wtIndex = this.WorkerTemplates.indexOf(wt);
					if (wtIndex > -1) {
						this.WorkerTemplates.splice(wtIndex, 1);
					}
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
                uploadWorkerTemplates(){
                    var params = [];
                    for (var i = 0; i < this.WorkerTemplates.length; i++){
                        var wt = this.WorkerTemplates[i]
                        var param = encodeURIComponent( wt.Name ) + '=' + encodeURIComponent( wt.Value );
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
                    xmlHttpRequest.open( 'POST', '/katib/workertemplates' );
                    xmlHttpRequest.setRequestHeader( 'Content-Type', 'application/x-www-form-urlencoded' );
                    xmlHttpRequest.send( params.join( '&' ).replace( /%20/g, '+' ) );
                    this.edited = false;
                }
			}
        });
    </script>
{{ end }}
