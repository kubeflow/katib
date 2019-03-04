import React from 'react';
import withStyles from '@material-ui/styles/withStyles';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import TextField from '@material-ui/core/TextField';
import DialogTitle from '@material-ui/core/DialogTitle';
import Slide from '@material-ui/core/Slide';

import { connect } from 'react-redux';
import { closeDialog, changeTemplate } from '../../../actions/templateActions';

const module = "template";

const styles = theme => ({
    textField: {
        margin: 10,
        width: 400,
    },
});

function Transition(props) {
    return <Slide direction={"up"} {...props} />
}

// FIX DIALOG TEXTFIELD SIZE

class EditDialog extends React.Component {

    state = {
        name: '',
        yaml: '',
    }

    componentDidMount() {
        this.setState({
            name: this.props.edittedTemplate.name,
            yaml: this.props.edittedTemplate.yaml
        })
    }

    componentWillReceiveProps(newProps) {
        this.setState({
            name: newProps.edittedTemplate.name,
            yaml: newProps.edittedTemplate.yaml,
        })
    }
    render () {
        const { classes } = this.props;
        return (
            <div>
                <Dialog
                    open={this.props.editOpen}
                    TransitionComponent={Transition}
                    keepMounted
                    onClose={this.props.closeDialog}
                >
                    <DialogTitle >
                        {"Editing a template"}
                    </DialogTitle>
                    <DialogContent>
                        <TextField
                            className={classes.textField}
                            value={this.state.name}
                            onChange={(event) => this.setState({
                                name: event.target.value
                            })}
                        /> 
                        <br />
                        <TextField 
                            multiline
                            className={classes.textField}
                            variant={"outlined"}
                            rows={"10"}
                            value={this.state.yaml}
                            onChange={(event) => this.setState({
                                yaml: event.target.value
                            })}
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={this.props.closeDialog} color={"primary"}>
                            Save
                        </Button>
                        <Button onClick={this.props.closeDialog} color={"secondary"}>
                            Discard
                        </Button>
                    </DialogActions>
                </Dialog>
            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        editOpen: state[module].editOpen,
        currentTemplateIndex: state[module].currentTemplateIndex,
        edittedTemplate: state[module].edittedTemplate,
    };
};

export default connect(mapStateToProps, { closeDialog, changeTemplate })(withStyles(styles)(EditDialog));