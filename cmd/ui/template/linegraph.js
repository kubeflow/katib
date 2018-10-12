{{ define "linegraphcjs" }}
var timeFormat = 'YYYY-MM-DDTH:mm:ss';

    var config = {
        type:    'scatter',
        data:    {
            datasets: [
                {{- range .MetricsLogs}}
                {
                    label: "{{.Name}}",
                    type: "line",
                    showLine: true,
                    data: [
                        {{- range .LogValues}}
                        {
                            x: "{{.Time}}", 
                            y: {{.Value}},
                        }, 
                        {{ end }}
                    ],
                    fill: false,
                    borderColor: "{{.Color}}"
                },
                {{- end }}
            ]
        },
        options: {
            responsive: true,
            title:      {
                display: true,
                text:    "Metrics Logs",
                fontColor: "#EEE",
                fontSize: 24,
            },
            legend: {
                labels: {
                    fontColor: "#EEE",
                    fontSize: 18,
                }
            },
            scales:     {
                xAxes: [{
                    scaleLabel: {
                        display:     true,
                        labelString: 'Time (Second)',
                        fontColor: "#EEE",
                        fontSize: 18,
                    },
                    ticks: {
                        fontColor: "#EEE", 
                    },
                }],
                yAxes: [{
                    scaleLabel: {
                        display:     true,
                        labelString: 'Metrics Value',
                        fontColor: "#EEE",
                        fontSize: 18,
                    },
                    ticks: {
                        fontColor: "#EEE",
                    },
                }]
            }
        }
    };

    window.onload = function () {
        var ctx       = document.getElementById("canvas").getContext("2d");
        window.myLine = new Chart(ctx, config);
    };
{{ end }}
