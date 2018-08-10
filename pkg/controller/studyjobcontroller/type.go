package studyjobcontroller

const (
	DefaultWorkerTemplate = `apiVersion: batch/v1
kind: Job
metadata:
  name: {{.WorkerId}}
spec:
  template:
    spec:
      containers:
      - name: {{.WorkerId}}
        image: {{.Image}}
        command: 
		{{- range .Command}}
		  - "{{.}}"
		{{- end}}
		{{- with .HyperParameters}}
		{{- range .}}
		  - "{{.Key}}={{.Value}}"
		{{- end}}
		{{- end}}
		{{- with .WorkerParameters}}
	    {{- if .RestartPolicy}}
        restartPolicy: {{.RestartPolicy}}
	    {{- else}}
        restartPolicy: Never
	    {{- end}}
		{{- end}}
	    volumeMounts:
	    {{- with .VolumeConfigs}}
	    {{- range .}}
	      - name: {{.Name}}
	        mountPath: {{.MountPath}}
	    {{- end}}
	    {{- end}}
	  volumes:
	  {{- with .VolumeConfigs}}
	  {{- range .}}
      - name: {{.Name}}
	    persistentVolumeClaim:
		claimName: {{.PvcName}}
	  {{- end}}
	  {{- end}}
`
)
