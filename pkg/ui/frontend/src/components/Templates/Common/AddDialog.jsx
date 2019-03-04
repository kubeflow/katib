import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import TextField from '@material-ui/core/TextField';
import DialogTitle from '@material-ui/core/DialogTitle';
import Slide from '@material-ui/core/Slide';

import { connect } from 'react-redux';
import { closeDialog } from '../../../actions/templateActions';

const module = "template";

const useStyles = makeStyles({
    textField: {
        marginLeft: 4,
        marginRight: 4,
        width: 400,
    },
});

function Transition(props) {
    return <Slide direction={"up"} {...props} />
}

// FIX DIALOG TEXTFIELD SIZE

const AddDialog = (props) => {
    const classes = useStyles();

    return (
        <div>
            <Dialog
                open={props.addOpen}
                TransitionComponent={Transition}
                keepMounted
                onClose={props.closeDialog}
            >
                <DialogTitle >
                    {"Adding a template"}
                </DialogTitle>
                <DialogContent>
                    <TextField 
                        multiline
                        className={classes.textField}
                        variant={"outlined"}
                        rows={"10"}
                        />
                </DialogContent>
                <DialogActions>
                    <Button onClick={props.closeDialog} color={"primary"}>
                        Save
                    </Button>
                    <Button onClick={props.closeDialog} color={"secondary"}>
                        Discard
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    )
}

const mapStateToProps = (state) => {
    return {
        addOpen: state[module].addOpen,
    };
};

export default connect(mapStateToProps, { closeDialog })(AddDialog);