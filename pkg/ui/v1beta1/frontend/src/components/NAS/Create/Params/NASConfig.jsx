import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import IconButton from '@material-ui/core/IconButton';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import Fab from '@material-ui/core/Fab';
import Divider from '@material-ui/core/Divider';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import {
  editNumLayers,
  addSize,
  editSize,
  deleteSize,
  addOperation,
  deleteOperation,
  changeOperation,
  addParameter,
  changeParameter,
  deleteParameter,
  addListParameter,
  editListParameter,
  deleteListParameter,
} from '../../../../actions/nasCreateActions';

import { NAS_CREATE_MODULE } from '../../../../constants/constants';

const useStyles = makeStyles({
  textField: {
    width: '80%',
  },
  numLayers: {
    padding: 2,
    marginBottom: 30,
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
  formControl: {
    width: '100%',
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
  section: {
    marginTop: 20,
  },
});

const SectionInTypography = (name, classes, variant) => {
  return (
    <div className={classes.section}>
      <Grid container>
        <Grid item xs={12} sm={12}>
          <Typography variant={variant}>{name}</Typography>
          <hr />
        </Grid>
      </Grid>
    </div>
  );
};

const NASConfig = props => {
  const classes = useStyles();

  const onEditNumLayers = () => event => {
    props.editNumLayers(event.target.value);
  };
  const onAddSize = type => event => {
    props.addSize(type);
  };

  const onEditSize = (index, type) => event => {
    props.editSize(type, index, event.target.value);
  };

  const onDeleteSize = (index, type) => event => {
    props.deleteSize(type, index);
  };

  const onDeleteOperation = index => event => {
    props.deleteOperation(index);
  };

  const onChangeOperation = index => event => {
    props.changeOperation(index, event.target.value);
  };

  const onAddParameter = opIndex => event => {
    props.addParameter(opIndex);
  };

  const onChangeParameter = (opIndex, paramIndex, name) => event => {
    props.changeParameter(opIndex, paramIndex, name, event.target.value);
  };

  const onDeleteParameter = (opIndex, paramIndex) => event => {
    props.deleteParameter(opIndex, paramIndex);
  };

  const onAddListParameter = (opIndex, paramIndex) => event => {
    props.addListParameter(opIndex, paramIndex);
  };

  const onDeleteListParameter = (opIndex, paramIndex, listIndex) => event => {
    props.deleteListParameter(opIndex, paramIndex, listIndex);
  };

  const onEditListParameter = (opIndex, paramIndex, listIndex) => event => {
    props.editListParameter(opIndex, paramIndex, listIndex, event.target.value);
  };

  return (
    <div>
      {/* NUM LAYERS */}
      <div className={classes.numLayers}>
        <Grid container alignItems={'center'}>
          <Grid item xs={12} sm={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip title="Number of layers">
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'NumLayers'}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              className={classes.textField}
              value={props.numLayers}
              onChange={onEditNumLayers()}
            />
          </Grid>
        </Grid>
      </div>
      {/* INPUT SIZE */}
      <div className={classes.parameter}>
        <Grid container alignItems={'center'}>
          <Grid item xs={12} sm={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip title={'Dimensions of input sizes'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'InputSizes'}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            {props.inputSize.map((size, index) => {
              return (
                <div key={index}>
                  <TextField
                    className={classes.textField}
                    value={size}
                    onChange={onEditSize(index, 'inputSize')}
                  />
                  <IconButton
                    key="close"
                    aria-label="Close"
                    color={'primary'}
                    className={classes.fab}
                    onClick={onDeleteSize(index, 'inputSize')}
                  >
                    <DeleteIcon />
                  </IconButton>
                </div>
              );
            })}
          </Grid>
          <Grid item xs={12} sm={2}>
            <Fab color={'primary'} className={classes.fab} onClick={onAddSize('inputSize')}>
              <AddIcon />
            </Fab>
          </Grid>
        </Grid>
      </div>
      {/* OUTPUT SIZE */}
      <div className={classes.parameter}>
        <Grid container alignItems={'center'}>
          <Grid item xs={12} sm={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip title={'Dimensions of output sizes'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'OutputSizes'}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            {props.outputSize.map((size, index) => {
              return (
                <div key={index}>
                  <TextField
                    className={classes.textField}
                    value={size}
                    onChange={onEditSize(index, 'outputSize')}
                  />
                  <IconButton
                    key="close"
                    aria-label="Close"
                    color={'primary'}
                    className={classes.fab}
                    onClick={onDeleteSize(index, 'outputSize')}
                  >
                    <DeleteIcon />
                  </IconButton>
                </div>
              );
            })}
          </Grid>
          <Grid item xs={12} sm={2}>
            <Fab color={'primary'} className={classes.fab} onClick={onAddSize('outputSize')}>
              <AddIcon />
            </Fab>
          </Grid>
        </Grid>
      </div>
      {/* OPERATIONS */}
      {SectionInTypography('Operations', classes, 'h6')}
      <div>
        <Button
          variant={'contained'}
          color={'primary'}
          className={classes.addButton}
          onClick={props.addOperation}
        >
          Add operation
        </Button>
      </div>
      {props.operations.map((operation, opIndex) => {
        return (
          <div key={opIndex}>
            <div className={classes.section}>
              <Grid container spacing={3}>
                <Grid item xs={4}>
                  <Typography variant={'h5'}>OperationType</Typography>
                </Grid>
                <Grid item xs={7}>
                  <TextField
                    value={operation.operationType}
                    className={classes.textField}
                    onChange={onChangeOperation(opIndex)}
                  />
                </Grid>
                <Grid item xs={1}>
                  <IconButton
                    key="close"
                    aria-label="Close"
                    color={'primary'}
                    className={classes.fab}
                    onClick={onDeleteOperation(opIndex)}
                  >
                    <DeleteIcon />
                  </IconButton>
                </Grid>
                <hr />
              </Grid>
              <div>
                <Button
                  variant={'contained'}
                  color={'primary'}
                  className={classes.addButton}
                  onClick={onAddParameter(opIndex)}
                >
                  Add parameter
                </Button>
              </div>
            </div>
            {operation.parameters.map((param, paramIndex) => {
              return (
                <div className={classes.parameter} key={paramIndex}>
                  <Grid container alignItems={'center'}>
                    <Grid item xs={1}>
                      <TextField
                        label={'Name'}
                        className={classes.textField}
                        value={param.name}
                        onChange={onChangeParameter(opIndex, paramIndex, 'name')}
                      />
                    </Grid>
                    <Grid item xs={2}>
                      <FormControl variant="outlined" className={classes.formControl}>
                        <InputLabel>Parameter Type</InputLabel>
                        <Select
                          onChange={onChangeParameter(opIndex, paramIndex, 'parameterType')}
                          value={param.parameterType}
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
                        onChange={onChangeParameter(opIndex, paramIndex, 'feasibleSpace')}
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
                                onChange={onEditListParameter(opIndex, paramIndex, elIndex)}
                              />
                              <IconButton
                                key="close"
                                aria-label="Close"
                                color={'primary'}
                                className={classes.icon}
                                onClick={onDeleteListParameter(opIndex, paramIndex, elIndex)}
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
                            onChange={onChangeParameter(opIndex, paramIndex, 'min')}
                          />
                          <TextField
                            label={'Max'}
                            className={classes.textField}
                            value={param.max}
                            onChange={onChangeParameter(opIndex, paramIndex, 'max')}
                          />
                          <TextField
                            label={'Step'}
                            className={classes.textField}
                            value={param.step}
                            onChange={onChangeParameter(opIndex, paramIndex, 'step')}
                          />
                        </div>
                      )}
                    </Grid>
                    <Grid item xs={1}>
                      {param.feasibleSpace === 'list' && (
                        <Fab
                          color={'primary'}
                          className={classes.fab}
                          onClick={onAddListParameter(opIndex, paramIndex)}
                        >
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
                        onClick={onDeleteParameter(opIndex, paramIndex)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Grid>
                  </Grid>
                </div>
              );
            })}
            <Divider />
          </div>
        );
      })}
    </div>
  );
};

const mapStateToProps = state => {
  return {
    numLayers: state[NAS_CREATE_MODULE].numLayers,
    inputSize: state[NAS_CREATE_MODULE].inputSize,
    outputSize: state[NAS_CREATE_MODULE].outputSize,
    operations: state[NAS_CREATE_MODULE].operations,
    allParameterTypes: state[NAS_CREATE_MODULE].allParameterTypes,
  };
};

export default connect(mapStateToProps, {
  editNumLayers,
  addSize,
  editSize,
  deleteSize,
  addOperation,
  deleteOperation,
  changeOperation,
  addParameter,
  changeParameter,
  deleteParameter,
  addListParameter,
  editListParameter,
  deleteListParameter,
})(NASConfig);
