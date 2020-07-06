import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';

import { changeSpec } from '../../../../actions/nasCreateActions';

import { NAS_CREATE_MODULE } from '../../../../constants/constants';

const useStyles = makeStyles({
  textField: {
    marginLeft: 4,
    marginRight: 4,
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

const CommonParametersSpec = props => {
  const classes = useStyles();

  const onSpecChange = name => event => {
    props.changeSpec(name, event.target.value);
  };

  return (
    <div>
      {props.commonParametersSpec.map((param, i) => {
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
                <TextField
                  className={classes.textField}
                  value={param.value}
                  onChange={onSpecChange(param.name)}
                />
              </Grid>
            </Grid>
          </div>
        );
      })}
    </div>
  );
};

const mapStateToProps = state => {
  return {
    commonParametersSpec: state[NAS_CREATE_MODULE].commonParametersSpec,
  };
};

export default connect(mapStateToProps, { changeSpec })(CommonParametersSpec);
