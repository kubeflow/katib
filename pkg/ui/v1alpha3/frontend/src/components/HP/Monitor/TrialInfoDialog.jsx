import React from 'react';
import { withStyles } from  '@material-ui/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import { connect } from 'react-redux';

import { closeDialog } from '../../../actions/hpMonitorActions';
import Plot from 'react-plotly.js';


const module = "hpMonitor";

const styles = theme => ({
    header: {
        textAlign: "center"
    }
})


const TrialInfoDialog = (props) => {
    const { classes } = props;
      
    let dataToPlot = [];
    if (props.trialData.length !== 0) { 
        let data = props.trialData.slice(1);   
        let tracks = {};
        for(let i = 0; i < data.length; i++) {
            // Data format should be ["metricName", "time", "value"]
            if (data[i].length == 3) {
                if (typeof tracks[data[i][0]] !== "undefined") {
                    // Formatted date if seconds < 10. Length of date should be the same
                    if (data[i][1].length == 18) {
                        let formattedDate = data[i][1].slice(0, 17) + "0" + data[i][1][17]
                        tracks[data[i][0]].x.push(formattedDate);
                    } else {
                        tracks[data[i][0]].x.push(data[i][1]);
                    }
                    tracks[data[i][0]].y.push(Number(data[i][2]));
                } else {
                    tracks[data[i][0]] = {};
                    if (data[i][1].length == 18) {
                        let formattedDate = data[i][1].slice(0, 17) + "0" + data[i][1][17]
                        tracks[data[i][0]].x = [formattedDate];
                    } else {
                        tracks[data[i][0]].x = [data[i][1]];
                    }
                    tracks[data[i][0]].y = [Number(data[i][2])];
                }
            }
        }

        //For plot legend
        let keys = Object.keys(tracks);
        keys.map((key, i) => {
            if (key !== "") {
                dataToPlot.push({
                    x: tracks[key].x,
                    y: tracks[key].y,
                    type: "scatter",
                    mode: "line",
                    name: key,
                    showlegend: true,
                    hoverinfo: "x+y"
                })
            }
        })
    }
    return (
        <Dialog
                open={props.open}
                onClose={props.closeDialog}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
                maxWidth={"xl"}
            >
            <DialogTitle id="alert-dialog-title" className = {classes.header}>{"Trial data"}</DialogTitle>
            <DialogContent>
                <Plot 
                    data={dataToPlot}
                    layout={{
                        width: 800,
                        height: 600,
                        xaxis: {
                            title: "Datetime",
                        },
                        yaxis: {
                            title: "Value",
                        }
                    }}
                    />
            </DialogContent>
      </Dialog>
    );
}

const mapStateToProps = state => {
    return {
        open: state[module].dialogOpen,
        trialData: state[module].trialData,
    }
}

export default connect(mapStateToProps, { closeDialog })(withStyles(styles)(TrialInfoDialog));
