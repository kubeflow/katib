import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-xcode';
import 'ace-builds/src-noconflict/mode-json';

import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';

import { closeDialogExperiment } from '../../actions/generalActions';

import { GENERAL_MODULE } from '../../constants/constants';

const styles = theme => ({
  header: {
    textAlign: 'center',
    width: 900,
  },
});

const ExperimentInfoDialog = props => {
  const { classes } = props;

  return (
    <Dialog
      open={props.dialogExperimentOpen}
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
        <AceEditor
          mode="json"
          theme="xcode"
          value={JSON.stringify(props.experiment, null, 2)}
          tabSize={2}
          fontSize={13}
          width={'100%'}
          showPrintMargin={false}
          autoScrollEditorIntoView={true}
          maxLines={40}
          minLines={10}
          readOnly={true}
          setOptions={{ useWorker: false }}
        />
      </DialogContent>
    </Dialog>
  );
};

const mapStateToProps = state => {
  return {
    dialogExperimentOpen: state[GENERAL_MODULE].dialogExperimentOpen,
    experiment: state[GENERAL_MODULE].experiment,
  };
};

export default connect(mapStateToProps, { closeDialogExperiment })(
  withStyles(styles)(ExperimentInfoDialog),
);
