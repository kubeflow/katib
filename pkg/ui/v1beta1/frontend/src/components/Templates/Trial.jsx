import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';

import TemplateList from './Common/TemplateList';

import { fetchTrialTemplates } from '../../actions/templateActions';
import { fetchNamespaces } from '../../actions/generalActions';

const styles = () => ({
  root: {
    width: '90%',
    margin: '0 auto',
    marginTop: 10,
  },
  text: {
    marginBottom: 20,
  },
});

class Trial extends React.Component {
  componentDidMount() {
    this.props.fetchTrialTemplates();
    this.props.fetchNamespaces();
  }

  render() {
    const { classes } = this.props;

    return (
      <div className={classes.root}>
        <Typography variant={'h4'} className={classes.text}>
          {'Trial Templates'}
        </Typography>

        <TemplateList />
      </div>
    );
  }
}
export default connect(null, { fetchTrialTemplates, fetchNamespaces })(withStyles(styles)(Trial));
