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
import { closeDialog, addTemplate } from '../../../actions/templateActions';

const module = "template";

const styles = theme => ({
    textField: {
        marginLeft: 4,
        marginRight: 4,
        width: 400,
        marginBottom: 10,
    },
});

function Transition(props) {
    return <Slide direction={"up"} {...props} />
}

// FIX DIALOG TEXTFIELD SIZE

class AddDialog extends React.Component {

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

    onChange = (name) => (event) => {
        this.setState({
            [name]: event.target.value
        })
    }

    addTemplate = () => {
        // this.props.addTemplate(this.state.name, this.state.yaml, this.props.type);
    }

    render () {
        const { classes } = this.props;

    // const addTemplate = () => {
    //     props.addTemplate(props.type);
    // }

        return (
            <div>
                <Dialog
                    open={this.props.addOpen}
                    TransitionComponent={Transition}
                    keepMounted
                    onClose={this.props.closeDialog}
                >
                    <DialogTitle >
                        {"Adding a template"}
                    </DialogTitle>
                    <DialogContent>
                        <TextField
                            className={classes.textField}
                            value={this.state.name}
                            onChange={this.onChange("name")}
                            />
                        <br />
                        <TextField 
                            multiline
                            className={classes.textField}
                            variant={"outlined"}
                            rows={"10"}
                            value={this.state.yaml}
                            onChange={this.onChange("yaml")}
                            />
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={this.addTemplate} color={"primary"}>
                            Save
                        </Button>
                        <Button onClick={this.props.closeDialog} color={"primary"}>
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
        addOpen: state[module].addOpen,
        edittedTemplate: state[module].edittedTemplate,
    };
};

export default connect(mapStateToProps, { closeDialog, addTemplate })(withStyles(styles)(AddDialog));