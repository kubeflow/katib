import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';

import { changeMeta } from '../../../../actions/nasCreateActions';

import { GENERAL_MODULE, NAS_CREATE_MODULE } from '../../../../constants/constants';

const styles = () => ({
  textField: {
    width: '100%',
  },
  help: {
    padding: 4 / 2,
    verticalAlign: 'middle',
    marginRight: 5,
  },
  parameter: {
    padding: 2,
    marginBottom: 10,
  },
});

class CommonParametersMeta extends React.Component {
  componentDidMount() {
    if (this.props.globalNamespace !== '') {
      this.props.changeMeta('Namespace', this.props.globalNamespace);
    }
  }

  onMetaChange = param => event => {
    this.props.changeMeta(param, event.target.value);
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        {this.props.commonParametersMetadata.map((param, i) => {
          return (
            <div key={i} className={classes.parameter}>
              <Grid container alignItems={'center'}>
                <Grid item xs={12} sm={3}>
                  <Typography variant={'subtitle1'}>
                    <Tooltip title={param.description}>
                      <HelpOutlineIcon className={classes.help} color={'primary'} />
                    </Tooltip>
                    {param.name}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={8}>
                  {param.name === 'Namespace' && this.props.globalNamespace === '' && (
                    <TextField
                      className={classes.textField}
                      value={param.value}
                      onChange={this.onMetaChange(param.name)}
                    />
                  )}
                  {param.name === 'Namespace' && this.props.globalNamespace !== '' && (
                    <TextField className={classes.textField} value={param.value} disabled />
                  )}
                  {param.name !== 'Namespace' && (
                    <TextField
                      className={classes.textField}
                      value={param.value}
                      onChange={this.onMetaChange(param.name)}
                    />
                  )}
                </Grid>
              </Grid>
            </div>
          );
        })}
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    commonParametersMetadata: state[NAS_CREATE_MODULE].commonParametersMetadata,
    globalNamespace: state[GENERAL_MODULE].globalNamespace,
  };
};

export default connect(mapStateToProps, { changeMeta })(withStyles(styles)(CommonParametersMeta));
