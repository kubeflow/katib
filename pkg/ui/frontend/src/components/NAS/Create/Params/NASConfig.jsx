import React from 'react';
import { connect } from 'react-redux';
import makeStyles from '@material-ui/styles/makeStyles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import IconButton from '@material-ui/core/IconButton';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import Fab from '@material-ui/core/Fab';

import { addSize, editSize, deleteSize } from '../../../../actions/nasCreateActions';

const module = "nasCreate";


const useStyles = makeStyles({
    textField: {
        marginLeft: 4,
        marginRight: 4,
        width: '80%'
    },
    help: {
        padding: 4 / 2,
        verticalAlign: "middle",
        marginRight: 5,
    },
    parameter: {
        padding: 2,
    },
    formControl: {
        margin: 4,
        width: '100%',
    },
    selectEmpty: {
        marginTop: 10,
    },
    group: {
        flexDirection: 'row',
        justifyContent: 'space-around',
    },
    divider: {
        margin: 5,
    },
    addButton: {
        margin: 10,
    },
    fab: {
        margin: 2,
    },
})

const NASConfig = (props) => {
    
    const classes = useStyles();

    const addSize = (type) => (event) => {
        props.addSize(type);
    }

    const editSize = (index, type) => (event) => {
        props.editSize(type, index, event.target.value);
    }

    const deleteSize = (index, type) => (event) => {
        props.deleteSize(type, index);
    }

    return (
        <div>
            {/* INPUT SIZE */}
            <div className={classes.parameter}>
                <Grid container alignItems={"center"}>
                    <Grid item xs={12} sm={3}>
                        <Typography variant={"subheading"}>
                            <Tooltip title={"Input Size"}>
                                <HelpOutlineIcon className={classes.help} color={"primary"}/>
                            </Tooltip>
                            {"Input Size"}
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={6}>
                        {props.inputSize.map((size, index) => {
                            return (
                                <div>
                                    <TextField
                                        className={classes.textField}
                                        value={size}
                                        onChange={editSize(index, "inputSize")}
                                    />
                                    <IconButton
                                        key="close"
                                        aria-label="Close"
                                        color={"primary"}
                                        className={classes.fab}
                                        onClick={deleteSize(index, "inputSize")}
                                        >
                                            <DeleteIcon />
                                    </IconButton>
                                </div>
                            )
                        })}
                    </Grid>
                    <Grid item xs={12} sm={2}>
                        <Fab color={"primary"} className={classes.fab} onClick={addSize("inputSize")}>
                            <AddIcon />
                         </Fab>
                    </Grid>
                </Grid>
            </div>
            {/* OUTOUT SIZE */}
            <div className={classes.parameter}>
                <Grid container alignItems={"center"}>
                    <Grid item xs={12} sm={3}>
                        <Typography variant={"subheading"}>
                            <Tooltip title={"Output Size"}>
                                <HelpOutlineIcon className={classes.help} color={"primary"}/>
                            </Tooltip>
                            {"Output Size"}
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={6}>
                        {props.outputSize.map((size, index) => {
                            return (
                                <div>
                                    <TextField
                                        className={classes.textField}
                                        value={size}
                                        onChange={editSize(index, "outputSize")}
                                    />
                                    <IconButton
                                        key="close"
                                        aria-label="Close"
                                        color={"primary"}
                                        className={classes.fab}
                                        onClick={deleteSize(index, "outputSize")}
                                        >
                                            <DeleteIcon />
                                    </IconButton>
                                </div>
                            )
                        })}
                    </Grid>
                    <Grid item xs={12} sm={2}>
                        <Fab color={"primary"} className={classes.fab} onClick={addSize("outputSize")} >
                            <AddIcon />
                         </Fab>
                    </Grid>
                </Grid>
            </div>
            <Button variant={"contained"} color={"primary"} className={classes.addButton}>
                    Add parameter
            </Button>
        </div>
    )
}


const mapStateToProps = state => {
    return {
        inputSize: state[module].inputSize,
        outputSize: state[module].outputSize,
    }
}

export default connect(mapStateToProps, { addSize, editSize, deleteSize })(NASConfig);