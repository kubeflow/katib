import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import MenuItem from '@material-ui/core/MenuItem';

import {
  changeObjective,
  addMetrics,
  editMetrics,
  deleteMetrics,
} from '../../../../actions/nasCreateActions';

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
  selectBox: {
    width: 150,
  },
});

const Objective = props => {
  const classes = useStyles();

  const onObjectiveChange = name => event => {
    props.changeObjective(name, event.target.value);
  };

  const onMetricsEdit = index => event => {
    props.editMetrics(index, event.target.value);
  };

  const onMetricsDelete = index => event => {
    props.deleteMetrics(index);
  };

  return (
    <div>
      {props.objective.map((param, i) => {
        return param.name === 'Type' ? (
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
                  <InputLabel>Objective Type</InputLabel>
                  <Select
                    value={param.value}
                    onChange={onObjectiveChange(param.name)}
                    input={<OutlinedInput labelWidth={160} />}
                    className={classes.selectBox}
                  >
                    {props.allObjectiveTypes.map((type, i) => {
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
                  onChange={onObjectiveChange(param.name)}
                />
              </Grid>
            </Grid>
          </div>
        );
      })}
      <div className={classes.parameter}>
        <Grid container alignItems={'center'}>
          <Grid item xs={12} sm={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip title={'Additional metrics that you want to collect'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              AdditionalMetricNames
            </Typography>
          </Grid>
          <Grid item xs={12} sm={8}>
            {props.additionalMetricNames.map((metrics, mIndex) => {
              return (
                <Grid container key={mIndex}>
                  <Grid item xs={10}>
                    <TextField
                      className={classes.textField}
                      value={metrics.value}
                      onChange={onMetricsEdit(mIndex)}
                    />
                  </Grid>
                  <Grid item xs={2}>
                    <IconButton
                      key="close"
                      aria-label="Close"
                      color={'primary'}
                      className={classes.icon}
                      onClick={onMetricsDelete(mIndex)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Grid>
                </Grid>
              );
            })}
          </Grid>
          <Grid item xs={12} sm={1}>
            <Fab color={'primary'} className={classes.fab} onClick={props.addMetrics}>
              <AddIcon />
            </Fab>
          </Grid>
        </Grid>
      </div>
    </div>
  );
};

const mapStateToProps = state => {
  return {
    allObjectiveTypes: state[NAS_CREATE_MODULE].allObjectiveTypes,
    objective: state[NAS_CREATE_MODULE].objective,
    additionalMetricNames: state[NAS_CREATE_MODULE].additionalMetricNames,
  };
};

export default connect(mapStateToProps, {
  changeObjective,
  addMetrics,
  editMetrics,
  deleteMetrics,
})(Objective);
