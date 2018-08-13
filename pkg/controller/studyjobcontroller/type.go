package studyjobcontroller

const (
	DefaultWorkerTemplate = `apiVersion: batch/v1
kind: Job
metadata:
  name: {{.WorkerId}}
  namespace: {{.NameSpace}}
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
          - "{{.Name}}={{.Value}}"
        {{- end}}
        {{- end}}
        {{- with .VolumeConfigs}}
        volumeMounts:
        {{- range .}}
          - name: {{.Name}}
            mountPath: {{.MountPath}}
        {{- end}}
        {{- end}}
      {{- if .WorkerParameters.RestartPolicy}}
      restartPolicy: {{.WorkerParameters.RestartPolicy}}
      {{- else}}
      restartPolicy: Never
      {{- end}}
      {{- with .VolumeConfigs}}
      volumes:
      {{- range .}}
      - name: {{.Name}}
        persistentVolumeClaim:
        claimName: {{.PvcName}}
      {{- end}}
      {{- end}}
`
	DefaultMetricsCollectorTemplate = `apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: {{.WorkerId}}
  namespace: {{.NameSpace}}
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          {{- if .ServiceAccount}}
          serviceAccountName: {{.ServiceAccount}}
          {{- else}}
          serviceAccountName: metrics-collector
          {{- end}}
          containers:
          - name: {{.WorkerId}}
            image: katib/metrics-collector
            args:
            - "./metricscollector"
            - "-s"
            - "{{.StudyId}}"
            - "-w"
            - "{{.WorkerId}}"
            - "-n"
            - "{{.NameSpace}}"
          restartPolicy: Never
`
)
