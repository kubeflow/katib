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
                    data={[
                    {
                        x: [1, 2, 3],
                        y: [2, 6, 3],
                        type: 'scatter',
                        mode: 'lines+points',
                        marker: {color: 'red'},
                    },
                    {type: 'bar', x: [1, 2, 3], y: [2, 5, 3]},
                    ]}
                    layout={ {title: `Job id: ${this.props.match.params.id}`} }
                />
            </div>
        )
    }
}

const mapStateToProps = (state) => ({
  
})


export default connect(mapStateToProps, null)(withStyles(styles)(HPJobInfo));
