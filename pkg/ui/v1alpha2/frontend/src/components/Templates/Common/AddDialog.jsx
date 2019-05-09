import React from 'react';
import withStyles from '@material-ui/styles/withStyles';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import TextField from '@material-ui/core/TextField';
import DialogTitle from '@material-ui/core/DialogTitle';
import Slide from '@material-ui/core/Slide';
import AceEditor from 'react-ace';

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

    onChangeYaml = (name) => (value) => {
        this.setState({
            [name]: value,
        })
    }

    addTemplate = () => {
        this.props.addTemplate(this.state.name, this.state.yaml, this.props.type, "add");
    }

    render () {
        const { classes } = this.props;

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
                        <AceEditor
                            mode="text"
                            theme="tomorrow"
                            value={this.state.yaml}
                            onChange={this.onChangeYaml("yaml")}
                            name="UNIQUE_ID_OF_DIV"
                            editorProps={{$blockScrolling: true}}
                            tabSize={2}
                            enableLiveAutocompletion={true}
                            fontSize={14}
                            width={480}
                            height={640}
                            />
                        {/* <TextField 
                            multiline
                            className={classes.textField}
                            variant={"outlined"}
                            rows={"10"}
                            value={this.state.yaml}
                            onChange={this.onChange("yaml")}
                            /> */}
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