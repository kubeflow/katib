import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Slide from '@material-ui/core/Slide';

import { connect } from 'react-redux';
import { closeDialog, deleteTemplate } from '../../../actions/templateActions';

const module = "template";

const useStyles = makeStyles({

});

function Transition(props) {
    return <Slide direction={"up"} {...props} />
}

const DeleteDialog = (props) => {
    const classes = useStyles();

    const deleteTemplate = (type) => (event) => {
        props.deleteTemplate(props.currentTemplateName, type);
    }

    return (
        <div>
            <Dialog
                open={props.deleteOpen}
                TransitionComponent={Transition}
                keepMounted
                onClose={props.closeDialog}
                aria-labelledby="alert-dialog-slide-title"
                aria-describedby="alert-dialog-slide-description"
            >
                <DialogTitle >
                    {"Are you sure?"}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Let Google help apps determine location. This means sending anonymous location data to
                        Google, even when no apps are running.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={props.closeDialog} color={"primary"}>
                        Disagee
                    </Button>
                    <Button onClick={deleteTemplate(props.type)} color={"secondary"}>
                        Agree
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    )
}

const mapStateToProps = (state) => {
    return {
        deleteOpen: state[module].deleteOpen,
        currentTemplateIndex: state[module].currentTemplateIndex,
        currentTemplateName: state[module].currentTemplateName,
    };
};

export default connect(mapStateToProps, { closeDialog, deleteTemplate })(DeleteDialog);