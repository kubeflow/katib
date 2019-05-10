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

const HPJobPlot = (props) => {
  const classes = useStyles();
  let dimensions = [];

  if (props.jobData && props.jobData.length !== 0) {
    // everything for the third column
    let header = props.jobData[0];
    let data = props.jobData.slice(1);
    for(let i = 1; i < data[0].length; i++) {
      if (header[i] !== '') {
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
            values.push(number)
          }
        }
        track.values = values;
        if (flag === "number" && flag !== "string") {
          track.range = [Math.min.apply(null, values), Math.max.apply(null, values)];
          if (Math.min.apply(null, values) < 1) {
            track.tickformat = ".3f"
          } else {
            track.tickformat = "d"
          }
        } else {
          // check logic
          // track.ticktext = values;
          let options = new Set(values);
          options = [...options]
          let mapping = {};
          for (let k = 0; k < options.length; k++) {
            mapping[options[k]] = k;
          }
          track.tickvals = options.map((option, index) => index);
          track.ticktext = options.map((option, index) => option);
          track.values = values.map((value, index) => mapping[value])
          track.constraintrange = [0, values.length];
        }
        dimensions.push(track)
        
      }
    }
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

export default connect(mapStateToProps, null)(HPJobPlot)