import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Divider from '@material-ui/core/Divider';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';

import {
  addParameter,
  editParameter,
  deleteParameter,
  addListParameter,
  editListParameter,
  deleteListParameter,
} from '../../../../actions/hpCreateActions';

import { HP_CREATE_MODULE } from '../../../../constants/constants';

const useStyles = makeStyles({
  textField: {
    width: '80%',
  },
  parameter: {
    padding: 2,
    marginBottom: 10,
  },
  formControl: {
    width: '100%',
  },
  selectEmpty: {
    marginTop: 10,
  },
  group: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  addButton: {
    margin: 10,
  },
  fab: {
    margin: 2,
  },
});

const Parameters = props => {
  const classes = useStyles();

  const onDelete = index => event => {
    props.deleteParameter(index);
  };

  const onGeneralEdit = (index, field) => event => {
    props.editParameter(index, field, event.target.value);
  };

  const onParamAdd = index => event => {
    props.addListParameter(index);
  };

  const onParamEdit = (paramIndex, index) => event => {
    props.editListParameter(paramIndex, index, event.target.value);
  };

  const onParamDelete = (paramIndex, index) => event => {
    props.deleteListParameter(paramIndex, index);
  };

  return (
    <div>
      <Button
        variant={'contained'}
        color={'primary'}
        className={classes.addButton}
        onClick={props.addParameter}
      >
        Add parameter
      </Button>
      {props.parameters.map((param, i) => {
        return (
          <div className={classes.parameter} key={i}>
            <Grid container alignItems={'center'}>
              <Grid item xs={1}>
                <TextField
                  label={'Name'}
                  className={classes.textField}
                  value={param.name}
                  onChange={onGeneralEdit(i, 'name')}
                />
              </Grid>
              <Grid item xs={2}>
                <FormControl variant="outlined" className={classes.formControl}>
                  <InputLabel>Parameter Type</InputLabel>
                  <Select
                    value={param.parameterType}
                    onChange={onGeneralEdit(i, 'parameterType')}
                    label="Parameter Type"
                    className={classes.select}
                  >
                    {props.allParameterTypes.map((type, i) => {
                      return (
                        <MenuItem value={type} key={i}>
                          {type}
                        </MenuItem>
                      );
                    })}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={3}>
                <RadioGroup
                  aria-label="Gender"
                  name="gender1"
                  className={classes.group}
                  value={param.feasibleSpace}
                  onChange={onGeneralEdit(i, 'feasibleSpace')}
                >
                  <FormControlLabel
                    value="feasibleSpace"
                    control={<Radio color={'primary'} />}
                    label="FeasibleSpace"
                  />
                  <FormControlLabel
                    value="list"
                    control={<Radio color={'primary'} />}
                    label="List"
                  />
                </RadioGroup>
              </Grid>
              <Grid item xs={4}>
                {param.feasibleSpace === 'list' &&
                  param.list.map((element, elIndex) => {
                    return (
                      <div key={elIndex}>
                        <TextField
                          className={classes.textField}
                          value={element.value}
                          onChange={onParamEdit(i, elIndex)}
                        />
                        <IconButton
                          key="close"
                          aria-label="Close"
                          color={'primary'}
                          className={classes.icon}
                          onClick={onParamDelete(i, elIndex)}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </div>
                    );
                  })}
                {param.feasibleSpace === 'feasibleSpace' && (
                  <div>
                    <TextField
                      label={'Min'}
                      className={classes.textField}
                      value={param.min}
                      onChange={onGeneralEdit(i, 'min')}
                    />
                    <TextField
                      label={'Max'}
                      className={classes.textField}
                      value={param.max}
                      onChange={onGeneralEdit(i, 'max')}
                    />
                    {props.algorithmName === 'grid' && param.parameterType === 'double' && (
                      <TextField
                        label={'Step'}
                        className={classes.textField}
                        value={param.step}
                        onChange={onGeneralEdit(i, 'step')}
                      />
                    )}
                  </div>
                )}
              </Grid>
              <Grid item xs={1}>
                {param.feasibleSpace === 'list' && (
                  <Fab color={'primary'} className={classes.fab} onClick={onParamAdd(i)}>
                    <AddIcon />
                  </Fab>
                )}
              </Grid>
              <Grid item xs={1}>
                <IconButton
                  key="close"
                  aria-label="Close"
                  color={'primary'}
                  className={classes.fab}
                  onClick={onDelete(i)}
                >
                  <DeleteIcon />
                </IconButton>
              </Grid>
            </Grid>
            <Divider className={classes.divider} />
          </div>
        );
      })}
    </div>
  );
};

const mapStateToProps = state => {
  return {
    parameters: state[HP_CREATE_MODULE].parameters,
    allParameterTypes: state[HP_CREATE_MODULE].allParameterTypes,
    algorithmName: state[HP_CREATE_MODULE].algorithmName,
  };
};

export default connect(mapStateToProps, {
  addParameter,
  editParameter,
  deleteParameter,
  addListParameter,
  editListParameter,
  deleteListParameter,
})(Parameters);
