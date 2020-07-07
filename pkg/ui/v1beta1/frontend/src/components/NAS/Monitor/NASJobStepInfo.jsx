import React from 'react';

// TODO: graphviz-react requires --max_old_space_size=4096
// Think about switching to a different lib
import { Graphviz } from 'graphviz-react';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

const styles = theme => ({
  root: {
    margin: '0 auto',
    textAlign: 'center',
  },
  link: {
    textDecoration: 'none',
  },
});

class NASJobStepInfo extends React.Component {
  render() {
    const { step, classes } = this.props;
    return (
      <div className={classes.root}>
        <Typography variant={'h5'}>Architecture for Trial: {step.trialname}</Typography>
        <Graphviz dot={this.props.step.architecture} options={{ zoom: true }} />
        <br />
        {step.metricsname.map((metrics, index) => {
          return (
            <Typography variant={'h6'}>
              {step.metricsname[index]}: {step.metricsvalue[index]}.
            </Typography>
          );
        })}
        <br />
        {/* TODO: add link in backend */}
        {/* <a href={`${step.link}`} className={classes.link}>
          <Button variant={'contained'} color={'primary'}>
            Download
          </Button>
        </a> */}
      </div>
    );
  }
}

export default withStyles(styles)(NASJobStepInfo);
