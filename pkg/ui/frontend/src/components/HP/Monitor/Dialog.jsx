import React from 'react';
import { makeStyles } from '@material-ui/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import { connect } from 'react-redux';

import { closeDialog } from '../../../actions/hpMonitorActions';
import Plot from 'react-plotly.js';


const module = "hpMonitor";


const useStyles = makeStyles({
});

const PlotDialog = (props) => {
    const classes = useStyles();
    const trace1 = {
        x: [1, 2, 3, 4],
        y: [10, 15, 13, 17],
        mode: 'markers',
        type: 'scatter'
      };
      
    const trace2 = {
        x: [2, 3, 4, 5],
        y: [16, 5, 11, 9],
        mode: 'lines',
        type: 'scatter'
    };
      
    const trace3 = {
        x: [1, 2, 3, 4],
        y: [12, 9, 15, 12],
        mode: 'lines+markers',
        type: 'scatter'
    };
      
    const data = [trace1, trace2, trace3];

    return (
        <Dialog
                open={props.open}
                onClose={props.closeDialog}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
                maxWidth={"xl"}
            >
            <DialogTitle id="alert-dialog-title">{"TriadID data"}</DialogTitle>
            <DialogContent>
                <Plot 
                    data={data}
                    layout={{
                        width: 640,
                        height: 480,
                    }}
                    />
            </DialogContent>
      </Dialog>
    );
}

const mapStateToProps = state => {
    return {
        open: state[module].dialogOpen,
    }
}

export default connect(mapStateToProps, { closeDialog })(PlotDialog)