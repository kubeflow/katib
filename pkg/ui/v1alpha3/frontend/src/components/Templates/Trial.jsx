import React from 'react';
import withStyles from '@material-ui/styles/withStyles';

import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';

import TemplateList from './Common/TemplateList';

import { connect } from 'react-redux';
import { openDialog, fetchTrialTemplates } from '../../actions/templateActions';
import AddDialog from './Common/AddDialog';


const styles = theme => ({
    root: {
        flexGrow: 1,
        marginTop: 40,
    },
});

class Trial extends React.Component {

    componentDidMount() {
        // TODO: Add possibility to change namespace in Trial Manifest form
        // Right now we get templates only from kubeflow namespace
        this.props.fetchTrialTemplates("");
    }

    openAddDialog = () => {
        this.props.openDialog("add");
    }

    render () {
        const { classes } = this.props;

        const type = "trial";
        
        return (
            <div className={classes.root}>
                <Typography variant={"headline"} color={"primary"}>
                    Trial Manifests
                </Typography>
                <Button variant={"contained"} color={"primary"} onClick={this.openAddDialog}>
                    Add
                </Button>

                <TemplateList type={type} />
                <AddDialog type={type}/>
                
            </div>
        )
    }
}
export default connect(null, { openDialog, fetchTrialTemplates })(withStyles(styles)(Trial));
