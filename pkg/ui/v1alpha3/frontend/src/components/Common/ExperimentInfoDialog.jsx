import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import MonacoEditor from 'react-monaco-editor';

import { closeDialogExperiment } from '../../actions/generalActions';

const module = 'general';

const styles = theme => ({
  header: {
    textAlign: 'center',
  },
});

const ExperimentInfoDialog = props => {
  const { classes } = props;

  return (
    <Dialog
      open={props.open}
      onClose={props.closeDialogExperiment}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
      maxWidth={'xl'}
    >
      <DialogTitle id="alert-dialog-title" className={classes.header}>
        {props.experiment.metadata && props.experiment.metadata.name
          ? 'Experiment ' + JSON.stringify(props.experiment.metadata.name, null, 2)
          : ''}
      </DialogTitle>
      <DialogContent>
        <MonacoEditor
          value={JSON.stringify(props.experiment, null, 2)}
          width="900"
          height="650"
          language="json"
          options={{
            readOnly: true,
          }}
        />
      </DialogContent>
    </Dialog>
  );
};

const mapStateToProps = state => {
  return {
    open: state[module].dialogExperimentOpen,
    experiment: state[module].experiment,
  };
};

export default connect(mapStateToProps, { closeDialogExperiment })(
  withStyles(styles)(ExperimentInfoDialog),
);
