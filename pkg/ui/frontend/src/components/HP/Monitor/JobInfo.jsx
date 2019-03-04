import React from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';

import Plot from 'react-plotly.js';


const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
        paddingTop: 20,
    }
})

class HPJobInfo extends React.Component {

    render () {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <Link to="/hp_monitor">
                    <Button variant={"contained"} color={"primary"}>
                        Back
                    </Button>
                </Link>
                <Typography variant={"h5"}>
                    JOB INFO for {this.props.match.params.id}
                </Typography>
                <br />
                <Plot
                    data={[{
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
                        }]}]}
                    layout={ {title: `Job id: ${this.props.match.params.id}`} }
                />
            </div>
        )
    }
}

const mapStateToProps = (state) => ({
  
})


export default connect(mapStateToProps, null)(withStyles(styles)(HPJobInfo));
