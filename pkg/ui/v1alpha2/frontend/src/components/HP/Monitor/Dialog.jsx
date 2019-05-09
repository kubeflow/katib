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

})


const PlotDialog = (props) => {
    const { classes } = props;
      
    let dataToPlot = [];
    if (props.TrialData.length !== 0) { 
        let data = props.TrialData.slice(1);   
        let tracks = {};
        for(let i = 0; i < data.length; i++) {
            if (typeof tracks[data[i][0]] !== "undefined") {
                tracks[data[i][0]].x.push(data[i][1]);
                tracks[data[i][0]].y.push(Number(data[i][2]));
            } else {
                tracks[data[i][0]] = {};
                tracks[data[i][0]].x = [data[i][1]];
                tracks[data[i][0]].y = [Number(data[i][2])];
            }
        }
        let keys = Object.keys(tracks);
        keys.map((key, i) => {
            if (key !== "") {
                dataToPlot.push({
                    x: tracks[key].x,
                    y: tracks[key].y,
                    type: "scatter",
                    mode: "line",
                    name: key,
                })
            }
        })
        console.log(dataToPlot)
    }
    return (
        <Dialog
                open={props.open}
                onClose={props.closeDialog}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
                maxWidth={"xl"}
            >
            <DialogTitle id="alert-dialog-title">{"Trial data"}</DialogTitle>
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

export default connect(mapStateToProps, { closeDialog })(withStyles(styles)(PlotDialog));