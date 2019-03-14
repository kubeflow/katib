import React from 'react';
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';

import * as d3             from 'd3'
import * as d3Graphviz     from 'd3-graphviz'


const styles = theme => ({
})

class StepInfo extends React.Component {

    componentDidMount() {
        let dotSec = 'digraph  {a -> b}'

        const id = `graph${this.props.id}`
        d3.select(`#${id}`)
        .graphviz()
        .renderDot(this.props.step.architecture)
        .width(640)
        .height(480)
        .fit(true)
    }

    render() {

        const { step, classes } = this.props;
        const id = `graph${this.props.id}`
        return (
            <div>
                <Typography variant={"h6"}>
                    Architecture
                </Typography>
                <div id={id} style={{textAlign: "center"}} />
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
}
    

export default withStyles(styles)(StepInfo);