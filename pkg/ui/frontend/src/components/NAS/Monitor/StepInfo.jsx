import React from 'react';
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';

import * as d3             from 'd3'
import * as d3Graphviz     from 'd3-graphviz'


const styles = theme => ({
})

const StepInfo = (props) => {
    const { step, classes } = props;
    var dotSrc = 'digraph  {a -> b}';

    d3.select(".graph").graphviz().renderDot(dotSrc);
    return (
            <div>
                <Typography variant={"h6"}>
                    Architecture
                </Typography>
                <div className="graph">
                </div>
                <br />
                {step.metricsname.map((metrics, index) => {
                    return (
                        <Typography variant={"h6"}>
                            {step.metricsname[index]}: {step.metricsvalue[index]}.
                        </Typography>
                    )
                })}
                <br />
                <a href={`${step.link}`}>
                    <Button variant={"contained"} color={"primary"}>
                        Download
                    </Button>
                </a>
            </div>
    )
}



export default withStyles(styles)(StepInfo);