import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';

import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';

import TemplateList from './Common/TemplateList';

import { connect } from 'react-redux';
import { openDialog } from '../../actions/templateActions';
import AddDialog from './Common/AddDialog';

const useStyles = makeStyles({
    root: {
        flexGrow: 1,
        marginTop: 40,
    },
});

const Worker = (props) => {

    const classes = useStyles();

    const type = props.match.path.replace("\/", "");

    const openAddDialog = () => {
        props.openDialog("add");
    }
    return (
        <div className={classes.root}>
            <Typography variant={"headline"} color={"secondary"}>
                Worker Manifests
            </Typography>
            <Button variant={"contained"} color={"primary"} onClick={openAddDialog}>
                Add
            </Button>

            <TemplateList type={type} />
            <AddDialog />
            
        </div>
    )
}
export default connect(null, { openDialog })(Worker)