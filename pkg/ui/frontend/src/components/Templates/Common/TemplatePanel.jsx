import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';

import { connect } from 'react-redux';

import DeleteIcon from '@material-ui/icons/Delete';
import CreateIcon from '@material-ui/icons/Create';

import { openDialog } from '../../../actions/templateActions';

const module = "template";

const useStyles = makeStyles({
    root: {
        flexGrow: 1,
    },
    grid: {
        marginTop: 30,
        textAlign: 'right',
    },
    icon: {
        margin: 4,
    }
});

const TemplatePanel = (props) => {

    const classes = useStyles();

    const openEditDialog = (index) => (event) => {
        props.openDialog("edit", index, props.type);
    };

    const openDeleteDialog = (index) => (event) => {
        props.openDialog("delete", index);
    };

    return (
        <div className={classes.root}>
            {props.text}
            <br />
            <Grid container spacing={24} className={classes.grid}>
                <Grid item xs={10}>
                    <Button variant={"contained"} color={"primary"} onClick={openEditDialog(props.index)}>
                        <CreateIcon color={"secondary"} className={classes.icon} />
                            Edit
                    </Button>
                </Grid>
                <Grid item xs={1}>
                    <Button variant={"contained"} color={"primary"} onClick={openDeleteDialog(props.index)}>
                        <DeleteIcon color={"secondary"} className={classes.icon} />
                            Delete
                    </Button>
                </Grid>
            </Grid>
        </div>
    )
}

const mapStateToProps = (state) => {
    return {

    };
};


export default connect(mapStateToProps, { openDialog })(TemplatePanel);