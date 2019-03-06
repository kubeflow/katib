import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';

import { connect } from 'react-redux';
import Plot from 'react-plotly.js';


const useStyles = makeStyles({
  root: {
    textAlign: 'center',
  }
})

const HPPlot = (props) => {
  const classes = useStyles();
    
  return (
      <div className={classes.root}>
        <Plot
             data={
                [{
                type: 'parcoords',
                line: {
                  color: 'blue'
                },
                
                dimensions: [{
                  range: [1, 5],
                  constraintrange: [1, 2],
                  label: 'A',
                  values: [1,4]
                }, {    
                  range: [1,5],
                  label: 'B',
                  values: [3,1.5],
                  tickvals: [1.5,3,4.5]
                }, {
                  range: [1, 5],
                  label: 'C',
                  values: [2,4],
                  tickvals: [1,2,4,5],
                  ticktext: ['text 1','text 2','text 4','text 5']
                }, {
                  range: [1, 5],
                  label: 'D',
                  values: [4,2]
                }]
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

export default connect(null, null)(HPPlot)