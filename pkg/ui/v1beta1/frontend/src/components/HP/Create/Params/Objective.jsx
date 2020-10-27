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
import MenuItem from '@material-ui/core/MenuItem';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';

import {
  changeObjective,
  addMetrics,
  editMetrics,
  deleteMetrics,
  metricStrategyChange,
} from '../../../../actions/hpCreateActions';

import { HP_CREATE_MODULE } from '../../../../constants/constants';

const useStyles = makeStyles({
  textField: {
    width: '100%',
  },
  textFieldStrategy: {
    width: '80%',
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
  checkBox: {
    textAlign: 'center',
  },
});

const Objective = props => {
  const classes = useStyles();

  const [checkedSetStrategies, setCheckedSetStrategies] = React.useState(false);

  const onCheckBoxChange = event => {
    setCheckedSetStrategies(event.target.checked);
  };

  const onObjectiveChange = name => event => {
    props.changeObjective(name, event.target.value);
  };

  const onMetricsEdit = index => event => {
    props.editMetrics(index, event.target.value);
  };

  const onMetricsDelete = index => event => {
    props.deleteMetrics(index);
  };

  const onMetricStrategyChange = index => event => {
    props.metricStrategyChange(index, event.target.value);
  };

  return (
    <div>
      {props.objective.map((param, i) => {
        return param.name === 'Type' ? (
          <Grid container alignItems={'center'} key={i} className={classes.parameter}>
            <Grid item xs={12} sm={3}>
              <Typography>
                <Tooltip title={param.description}>
                  <HelpOutlineIcon className={classes.help} color={'primary'} />
                </Tooltip>
                {param.name}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={8}>
              <FormControl variant="outlined">
                <InputLabel>Objective Type</InputLabel>
                <Select
                  value={param.value}
                  onChange={onObjectiveChange(param.name)}
                  className={classes.selectBox}
                  label="Objective Type"
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
        ) : (
          <Grid container alignItems={'center'} key={i} className={classes.parameter}>
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
        );
      })}
      <Grid container alignItems={'center'} className={classes.parameter}>
        <Grid item xs={12} sm={3}>
          <Typography variant={'subtitle1'}>
            <Tooltip title={'Additional metrics that you want to collect'}>
              <HelpOutlineIcon className={classes.help} color={'primary'} />
            </Tooltip>
            AdditionalMetricNames
          </Typography>
        </Grid>
        <Grid item xs={12} sm={8}>
          {props.additionalMetricNames.map((metric, mIndex) => {
            return (
              <Grid container key={mIndex}>
                <Grid item xs={10}>
                  <TextField
                    className={classes.textField}
                    value={metric}
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
      <Grid container alignItems={'center'} className={classes.parameter}>
        <Grid item sm={2}>
          <Typography variant={'subtitle1'}>
            <Tooltip title={'Strategy for extracting metrics to calculate objective'}>
              <HelpOutlineIcon className={classes.help} color={'primary'} />
            </Tooltip>
            MetricStrategies (optional)
          </Typography>
        </Grid>
        <Grid item sm={1} className={classes.checkBox}>
          <FormControlLabel
            control={
              <Checkbox
                checked={checkedSetStrategies}
                onChange={onCheckBoxChange}
                color="primary"
              />
            }
            label="Set"
          />
        </Grid>
        {checkedSetStrategies && (
          <Grid item sm={9}>
            {props.metricStrategies.map((metric, mIndex) => {
              return (
                <Grid container key={mIndex} className={classes.parameter}>
                  <Grid item xs={3}>
                    <TextField
                      className={classes.textFieldStrategy}
                      value={metric.name}
                      InputProps={{
                        readOnly: true,
                      }}
                    />
                  </Grid>
                  <Grid item xs={3}>
                    <FormControl variant="outlined">
                      <InputLabel>Strategy Type</InputLabel>
                      <Select
                        value={metric.strategy}
                        onChange={onMetricStrategyChange(mIndex)}
                        className={classes.selectBox}
                        label="Strategy Type"
                      >
                        {props.metricStrategiesList.map((strategy, i) => {
                          return (
                            <MenuItem value={strategy} key={i}>
                              {strategy}
                            </MenuItem>
                          );
                        })}
                      </Select>
                    </FormControl>
                  </Grid>
                </Grid>
              );
            })}
          </Grid>
        )}
      </Grid>
    </div>
  );
};

const mapStateToProps = state => {
  return {
    allObjectiveTypes: state[HP_CREATE_MODULE].allObjectiveTypes,
    objective: state[HP_CREATE_MODULE].objective,
    additionalMetricNames: state[HP_CREATE_MODULE].additionalMetricNames,
    metricStrategiesList: state[HP_CREATE_MODULE].metricStrategiesList,
    metricStrategies: state[HP_CREATE_MODULE].metricStrategies,
  };
};

export default connect(mapStateToProps, {
  changeObjective,
  addMetrics,
  editMetrics,
  deleteMetrics,
  metricStrategyChange,
})(Objective);
