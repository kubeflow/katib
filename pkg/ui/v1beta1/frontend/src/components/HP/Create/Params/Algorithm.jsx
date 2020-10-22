import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import TextField from '@material-ui/core/TextField';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';

import {
  changeAlgorithmName,
  addAlgorithmSetting,
  changeAlgorithmSetting,
  deleteAlgorithmSetting,
} from '../../../../actions/hpCreateActions';

import { HP_CREATE_MODULE } from '../../../../constants/constants';

const useStyles = makeStyles({
  textField: {
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
  icon: {
    padding: 4,
    margin: '0 auto',
    verticalAlign: 'middle !important',
  },
  formControl: {
    width: '100%',
  },
  addButton: {
    margin: 10,
  },
});

const Algorithm = props => {
  const classes = useStyles();

  const onAlgorithmNameChange = event => {
    props.changeAlgorithmName(event.target.value);
  };

  const onAddAlgorithmSetting = () => {
    props.addAlgorithmSetting();
  };

  const onChangeAlgorithmSetting = (name, index) => event => {
    props.changeAlgorithmSetting(index, name, event.target.value);
  };

  const onDeleteAlgorithmSetting = index => event => {
    props.deleteAlgorithmSetting(index);
  };
  return (
    <div>
      <Button
        variant={'contained'}
        color={'primary'}
        className={classes.addButton}
        onClick={onAddAlgorithmSetting}
      >
        Add algorithm setting
      </Button>
      <div className={classes.parameter}>
        <Grid container alignItems={'center'}>
          <Grid item xs={12} sm={3}>
            <Typography>
              <Tooltip title={'Name for the HP Algorithm'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'Algorithm Name'}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={8}>
            <FormControl variant="outlined" className={classes.formControl}>
              <InputLabel>Algorithm Name</InputLabel>
              <Select
                value={props.algorithmName}
                onChange={onAlgorithmNameChange}
                label="Algorithm Name"
              >
                {props.allAlgorithms.map((algorithm, i) => {
                  return (
                    <MenuItem value={algorithm} key={i}>
                      {algorithm}
                    </MenuItem>
                  );
                })}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
      </div>
      <br />
      {props.algorithmSettings.map((param, i) => {
        return (
          <div key={i} className={classes.parameter}>
            <Grid container alignItems={'center'}>
              <Grid item xs={3} />
              <Grid item xs={4}>
                <TextField
                  label={'Name'}
                  className={classes.textField}
                  value={param.name}
                  onChange={onChangeAlgorithmSetting('name', i)}
                />
              </Grid>
              <Grid item xs={4}>
                <TextField
                  label={'Value'}
                  className={classes.textField}
                  value={param.value}
                  onChange={onChangeAlgorithmSetting('value', i)}
                />
              </Grid>
              <Grid item xs={1}>
                <IconButton
                  key="close"
                  aria-label="Close"
                  color={'primary'}
                  className={classes.icon}
                  onClick={onDeleteAlgorithmSetting(i)}
                >
                  <DeleteIcon />
                </IconButton>
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
    algorithmName: state[HP_CREATE_MODULE].algorithmName,
    allAlgorithms: state[HP_CREATE_MODULE].allAlgorithms,
    algorithmSettings: state[HP_CREATE_MODULE].algorithmSettings,
  };
};

export default connect(mapStateToProps, {
  changeAlgorithmName,
  addAlgorithmSetting,
  changeAlgorithmSetting,
  deleteAlgorithmSetting,
})(Algorithm);
