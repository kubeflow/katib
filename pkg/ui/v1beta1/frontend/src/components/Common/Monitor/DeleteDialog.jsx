import React from 'react';
import { connect } from 'react-redux';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

import { deleteExperiment, closeDeleteExperimentDialog } from '../../../actions/generalActions';

import { GENERAL_MODULE } from '../../../constants/constants';

const DeleteDialog = props => {
  const onDelete = () => {
    props.deleteExperiment(props.deleteExperimentName, props.deleteExperimentNamespace);
  };

  return (
    <Dialog
      open={props.open}
      onClose={props.closeDeleteExperimentDialog}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{'Delete Experiment?'}</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are you sure you want to delete this experiment?
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.closeDeleteExperimentDialog} color="primary">
          Disagree
        </Button>
        <Button onClick={onDelete} color="primary" autoFocus>
          Agree
        </Button>
      </DialogActions>
    </Dialog>
  );
};

const mapStateToProps = state => ({
  open: state[GENERAL_MODULE].deleteDialog,
  deleteExperimentName: state[GENERAL_MODULE].deleteExperimentName,
  deleteExperimentNamespace: state[GENERAL_MODULE].deleteExperimentNamespace,
});

export default connect(mapStateToProps, { closeDeleteExperimentDialog, deleteExperiment })(
  DeleteDialog,
);
