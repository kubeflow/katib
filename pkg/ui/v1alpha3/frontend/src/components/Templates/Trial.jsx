import React from 'react';
import { connect } from 'react-redux';

import withStyles from '@material-ui/styles/withStyles';

import TemplateList from './Common/TemplateList';

import { fetchTrialTemplates } from '../../actions/templateActions';
import { fetchNamespaces } from '../../actions/generalActions';

const styles = theme => ({
  root: {
    width: '90%',
    margin: '0 auto',
    marginTop: 10,
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
        <h1>Trial Templates</h1>

        <TemplateList />
      </div>
    );
  }
}
export default connect(null, { fetchTrialTemplates, fetchNamespaces })(withStyles(styles)(Trial));
