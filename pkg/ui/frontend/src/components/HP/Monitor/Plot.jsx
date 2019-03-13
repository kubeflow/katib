import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';

import { connect } from 'react-redux';
import Plot from 'react-plotly.js';

const module = "hpMonitor";

const useStyles = makeStyles({
  root: {
    textAlign: 'center',
  }
})

const HPPlot = (props) => {
  const classes = useStyles();
  let dimensions = [];

  if (props.jobData.length !== 0) {
    // everything for the third column
    let header = props.jobData[0];
    let data = props.jobData.slice(1);
    for(let i = 2; i < data[0].length; i++) {
      let track = {
        label: header[i],
      }
      let flag = "number";
      let values = [];
      for (let j = 0; j < data.length - 1; j++) {
        let number = Number(data[j][i])
        if (isNaN(number)) {
          flag = "string";
          values.push(data[j][i]);
        } else {
          values.push(number);
        }
      }
      track.values = values;
      if (flag === "number" && flag !== "string") {
        track.constraintrange = [Math.min(values), Math.max(values)];
      } else {
        // check logic
        track.ticktext = values;
        let options = new Set(values);
        options = [...options]
        let mapping = {};
        for(let k = 0; k < options.length; k++) {
          mapping[options[k]] = k;
        }
        track.tickvals = options.map((option, index) => index);
        track.values = values.map((value, index) => mapping[value])
        track.constraintrange = [0, values.length];
      }
      dimensions.push(track)
    }
    console.log(dimensions)
    // dimensions= [{
    //   constraintrange: [1, 2],
    //   label: 'A',
    //   values: [1,4]
    // }, {    
    //   label: 'B',
    //   values: [3,1.5],
    //   tickvals: [1.5,3,4.5]
    // }, {
    //   label: 'C',
    //   values: [2,4],
    //   tickvals: [1,2,4,5],
    //   ticktext: ['text 1','text 2','text 4','text 5']
    // }, {
    //   label: 'D',
    //   values: [4,2]
    // }]
  }

  return (
      <div className={classes.root}>
        <Plot
          data={
            [{
              type: 'parcoords',
              line: {
                color: "red",
              },
              dimensions: dimensions,
            }]}
          layout={ { 
            title: `Job id: ${props.id}`,
            width: 1000,
            height: 600
          } }
        /> 
      </div>
    )
}

const mapStateToProps = (state) => ({
  jobData: state[module].jobData,
})

export default connect(mapStateToProps, null)(HPPlot)