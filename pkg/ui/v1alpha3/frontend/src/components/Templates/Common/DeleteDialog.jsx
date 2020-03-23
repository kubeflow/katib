import React from 'react';
import { connect } from 'react-redux';
import withStyles from '@material-ui/styles/withStyles';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

import { closeDialog, deleteTemplate } from '../../../actions/templateActions';

const module = 'template';

const DeleteDialog = props => {
  const submitDeleteTemplate = () => {
    props.deleteTemplate(
      props.edittedTemplateNamespace,
      props.edittedTemplateConfigMapName,
      props.edittedTemplateName,
    );
  };

  return (
    <div>
      <Dialog open={props.deleteOpen} onClose={props.closeDialog}>
        <DialogTitle id="alert-dialog-title">{'Are you sure?'}</DialogTitle>
        <DialogContent>
          <DialogContentText>Are you sure you want to delete this template?</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={props.closeDialog} color={'primary'}>
            Disagee
          </Button>
          <Button onClick={submitDeleteTemplate} color={'primary'}>
            Agree
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

const mapStateToProps = state => {
  return {
    deleteOpen: state[module].deleteOpen,
    edittedTemplateNamespace: state[module].edittedTemplateNamespace,
    edittedTemplateConfigMapName: state[module].edittedTemplateConfigMapName,
    edittedTemplateName: state[module].edittedTemplateName,
  };
};

export default connect(mapStateToProps, { closeDialog, deleteTemplate })(DeleteDialog);
