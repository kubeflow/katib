import React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import makeStyles from '@material-ui/styles/makeStyles';

import { connect } from 'react-redux';
import { deleteExperiment, closeDeleteExperimentDialog } from '../../actions/generalActions';


const module = "general";

const useStyles = makeStyles({
    root: {
    }
  })

const DeleteDialog = (props) => {
    const classes = useStyles();

    const onDelete = () => {
        props.deleteExperiment(props.deleteExperimentName);
    }

    return (
        <Dialog
          open={props.open}
          onClose={props.closeDeleteExperimentDialog}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
            <DialogTitle id="alert-dialog-title">{"Delete Experiment?"}</DialogTitle>
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
    )
}

const mapStateToProps = (state) => ({
    open: state[module].deleteDialog,
    deleteExperimentName: state[module].deleteExperimentName,
})

export default connect(mapStateToProps, { closeDeleteExperimentDialog, deleteExperiment })(DeleteDialog);
