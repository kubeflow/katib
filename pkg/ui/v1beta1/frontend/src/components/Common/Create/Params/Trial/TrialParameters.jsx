import React from 'react';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';

import { changeTrialParameters } from '../../../../../actions/generalActions';
import { GENERAL_MODULE } from '../../../../../constants/constants';

const useStyles = makeStyles({
  help: {
    padding: 4 / 2,
    verticalAlign: 'middle',
    marginRight: 5,
  },
  textField: {
    width: '80%',
  },
  parametersGrid: {
    marginBottom: 20,
  },
  div: {
    marginTop: 10,
    marginBottom: 10,
  },
});

const TrialParameters = props => {
  const classes = useStyles();

  const onNameChange = index => event => {
    let param = props.trialParameters[index];
    let reference = param.reference;
    let description = param.description;
    props.changeTrialParameters(index, event.target.value, reference, description);
  };

  const onReferenceChange = index => event => {
    let param = props.trialParameters[index];
    let name = param.name;
    let description = param.description;
    props.changeTrialParameters(index, name, event.target.value, description);
  };

  const onDescriptionChange = index => event => {
    let param = props.trialParameters[index];
    let name = param.name;
    let reference = param.reference;
    props.changeTrialParameters(index, name, reference, event.target.value);
  };

  return (
    <div className={classes.div}>
      {props.trialParameters.length > 0 ? (
        <Grid container alignItems={'center'}>
          <Grid item xs={3}>
            <Typography variant={'subtitle1'}>
              <Tooltip
                title={
                  'Parameters for the Trial Template. Name - parameter that must be replaced in Template, Reference - parameter from Suggestion assignments'
                }
              >
                <HelpOutlineIcon className={classes.help} color={'primary'} />
              </Tooltip>
              Trial Parameters
            </Typography>
          </Grid>
          <Grid item xs={9}>
            {props.trialParameters.map((trialParam, pIndex) => {
              return (
                <Grid container key={pIndex} className={classes.parametersGrid}>
                  <Grid item xs={3}>
                    <TextField
                      className={classes.textField}
                      value={trialParam.name}
                      placeholder="Name"
                      label="Name"
                      onChange={onNameChange(pIndex)}
                    />
                  </Grid>
                  <Grid item xs={3}>
                    <TextField
                      placeholder="Reference"
                      label="Reference"
                      className={classes.textField}
                      value={trialParam.reference}
                      onChange={onReferenceChange(pIndex)}
                    />
                  </Grid>
                  <Grid item xs={3}>
                    <TextField
                      className={classes.textField}
                      value={trialParam.description}
                      multiline
                      rows={4}
                      variant="outlined"
                      placeholder="Description"
                      label="Description"
                      onChange={onDescriptionChange(pIndex)}
                    />
                  </Grid>
                </Grid>
              );
            })}
          </Grid>
        </Grid>
      ) : (
        <Typography variant="h6">Unable to get parameters from Trial Spec</Typography>
      )}
    </div>
  );
};

const mapStateToProps = state => {
  return {
    trialParameters: state[GENERAL_MODULE].trialParameters,
  };
};

export default connect(mapStateToProps, { changeTrialParameters })(TrialParameters);
