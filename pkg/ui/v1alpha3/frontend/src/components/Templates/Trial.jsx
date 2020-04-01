import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import TemplateList from './Common/TemplateList';

import { fetchTrialTemplates } from '../../actions/templateActions';

const styles = theme => ({
  root: {
    width: '90%',
    margin: '0 auto',
    marginTop: 10,
  },
  text: {
    fontSize: theme.typography.pxToRem(40),
  },
});

class Trial extends React.Component {
  componentDidMount() {
    this.props.fetchTrialTemplates();
  }

  render() {
    const { classes } = this.props;

    return (
      <div className={classes.root}>
        <Typography variant={'h4'}>{'Trial Templates'}</Typography>

        <TemplateList />
      </div>
    );
  }
}
export default connect(null, { fetchTrialTemplates })(withStyles(styles)(Trial));
