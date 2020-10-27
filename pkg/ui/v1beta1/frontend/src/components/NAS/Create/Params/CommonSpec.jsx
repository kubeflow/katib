import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';

import { changeSpec } from '../../../../actions/nasCreateActions';

import { NAS_CREATE_MODULE } from '../../../../constants/constants';

const useStyles = makeStyles({
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
  selectBox: {
    width: 150,
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
        return param.name === 'ResumePolicy' ? (
          <div key={i} className={classes.parameter}>
            <Grid container alignItems={'center'}>
              <Grid item xs={12} sm={3}>
                <Typography>
                  <Tooltip title={param.description}>
                    <HelpOutlineIcon className={classes.help} color={'primary'} />
                  </Tooltip>
                  {param.name}
                </Typography>
              </Grid>
              <Grid item xs={12} sm={8}>
                <FormControl variant="outlined" className={classes.formControl}>
                  <InputLabel>Resume Policy</InputLabel>
                  <Select
                    value={param.value}
                    onChange={onSpecChange(param.name)}
                    className={classes.selectBox}
                    label="Resume Policy"
                  >
                    {props.allResumePolicyTypes.map((type, i) => {
                      return (
                        <MenuItem value={type} key={i}>
                          {type}
                        </MenuItem>
                      );
                    })}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </div>
        ) : (
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
    allResumePolicyTypes: state[NAS_CREATE_MODULE].allResumePolicyTypes,
  };
};

export default connect(mapStateToProps, { changeSpec })(CommonParametersSpec);
