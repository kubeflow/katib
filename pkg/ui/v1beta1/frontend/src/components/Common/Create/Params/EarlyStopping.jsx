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
  changeEarlyStoppingAlgorithm,
  addEarlyStoppingSetting,
  changeEarlyStoppingSetting,
  deleteEarlyStoppingSetting,
} from '../../../../actions/generalActions';

import { GENERAL_MODULE } from '../../../../constants/constants';

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

const EarlyStopping = props => {
  const classes = useStyles();

  const onEarlyStoppingAlgorithmChange = event => {
    props.changeEarlyStoppingAlgorithm(event.target.value);
  };

  const onAddEarlyStoppingSetting = () => {
    props.addEarlyStoppingSetting();
  };

  const onChangeEarlyStoppingSetting = (field, index) => event => {
    props.changeEarlyStoppingSetting(index, field, event.target.value);
  };

  const onDeleteEarlyStoppingSetting = index => event => {
    props.deleteEarlyStoppingSetting(index);
  };
  return (
    <div>
      <Button
        variant={'contained'}
        color={'primary'}
        className={classes.addButton}
        onClick={onAddEarlyStoppingSetting}
      >
        Add early stopping setting
      </Button>
      <div className={classes.parameter}>
        <Grid container alignItems={'center'}>
          <Grid item xs={12} sm={3}>
            <Typography>
              <Tooltip title={'Name for the Early Stopping Algorithm'}>
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              {'Early Stopping Algorithm Name'}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={8}>
            <FormControl variant="outlined" className={classes.formControl}>
              <InputLabel>Algorithm Name</InputLabel>
              <Select
                value={props.earlyStoppingAlgorithm}
                onChange={onEarlyStoppingAlgorithmChange}
                label="Algorithm Name"
              >
                {props.allEarlyStoppingAlgorithms.map((algorithm, i) => {
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
      {props.earlyStoppingSettings.map((setting, i) => {
        return (
          <div key={i} className={classes.parameter}>
            <Grid container alignItems={'center'}>
              <Grid item xs={3} />
              <Grid item xs={4}>
                <TextField
                  label={'Name'}
                  className={classes.textField}
                  value={setting.name}
                  onChange={onChangeEarlyStoppingSetting('name', i)}
                />
              </Grid>
              <Grid item xs={4}>
                <TextField
                  label={'Value'}
                  className={classes.textField}
                  value={setting.value}
                  onChange={onChangeEarlyStoppingSetting('value', i)}
                />
              </Grid>
              <Grid item xs={1}>
                <IconButton
                  key="close"
                  aria-label="Close"
                  color={'primary'}
                  className={classes.icon}
                  onClick={onDeleteEarlyStoppingSetting(i)}
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
    earlyStoppingAlgorithm: state[GENERAL_MODULE].earlyStoppingAlgorithm,
    allEarlyStoppingAlgorithms: state[GENERAL_MODULE].allEarlyStoppingAlgorithms,
    earlyStoppingSettings: state[GENERAL_MODULE].earlyStoppingSettings,
  };
};

export default connect(mapStateToProps, {
  changeEarlyStoppingAlgorithm,
  addEarlyStoppingSetting,
  changeEarlyStoppingSetting,
  deleteEarlyStoppingSetting,
})(EarlyStopping);
