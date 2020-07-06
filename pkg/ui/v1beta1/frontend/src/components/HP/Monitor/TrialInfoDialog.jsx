import React from 'react';

import Plot from 'react-plotly.js';

import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import { connect } from 'react-redux';

import { closeDialogTrial } from '../../../actions/hpMonitorActions';

import { HP_MONITOR_MODULE } from '../../../constants/constants';

const styles = theme => ({
  header: {
    textAlign: 'center',
  },
});

const TrialInfoDialog = props => {
  const { classes } = props;

  let dataToPlot = [];
  if (props.trialData.length !== 0) {
    let data = props.trialData.slice(1);
    let tracks = {};
    for (let i = 0; i < data.length; i++) {
      // Data format should be ["metricName", "time", "value"]
      if (data[i].length === 3) {
        if (typeof tracks[data[i][0]] !== 'undefined') {
          tracks[data[i][0]].x.push(data[i][1]);
          tracks[data[i][0]].y.push(Number(data[i][2]));
        } else {
          tracks[data[i][0]] = {};
          tracks[data[i][0]].x = [data[i][1]];
          tracks[data[i][0]].y = [Number(data[i][2])];
        }
      }
    }

    //For plot legend
    let keys = Object.keys(tracks);
    keys.map((key, i) => {
      if (key !== '') {
        dataToPlot.push({
          x: tracks[key].x,
          y: tracks[key].y,
          type: 'scatter',
          mode: 'line',
          name: key,
          showlegend: true,
          hoverinfo: 'x+y',
        });
      }
      return dataToPlot;
    });
  }
  return (
    <Dialog
      open={props.open}
      onClose={props.closeDialogTrial}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
      maxWidth={'xl'}
    >
      <DialogTitle id="alert-dialog-title" className={classes.header}>
        {'Trial Name: ' + props.trialName}
      </DialogTitle>
      <DialogContent>
        <Plot
          data={dataToPlot}
          layout={{
            width: 800,
            height: 600,
            xaxis: {
              title: 'Datetime',
            },
            yaxis: {
              title: 'Value',
            },
          }}
        />
      </DialogContent>
    </Dialog>
  );
};

const mapStateToProps = state => {
  return {
    open: state[HP_MONITOR_MODULE].dialogTrialOpen,
    trialData: state[HP_MONITOR_MODULE].trialData,
    trialName: state[HP_MONITOR_MODULE].trialName,
  };
};

export default connect(mapStateToProps, { closeDialogTrial })(withStyles(styles)(TrialInfoDialog));
