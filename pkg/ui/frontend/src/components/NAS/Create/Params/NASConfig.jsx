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
import Divider from '@material-ui/core/Divider';
import Select from '@material-ui/core/Select';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';

import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';


import { addSize, editSize, deleteSize, addOperation, deleteOperation, changeOperation, addParameter, changeParameter, deleteParameter, addListParameter, editListParameter, deleteListParameter } from '../../../../actions/nasCreateActions';

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
        marginBottom: 10,
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

const SectionInTypography = (name, classes, variant) => {
    return (
        <div className={classes.section}>
            <Grid container>
                <Grid item xs={12} sm={12}>
                    <Typography variant={variant}>
                        {name}
                    </Typography>
                <hr />
                </Grid>
            </Grid>
        </div>
    )
}


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

    const deleteOperation = (index) => (event) => {
        props.deleteOperation(index);
    }

    const changeOperation = (index) => (event) => {
        props.changeOperation(index, event.target.value);
    }

    const addParameter = (opIndex) => (event) => {
        props.addParameter(opIndex);
    }

    const changeParameter = (opIndex, paramIndex, name) => (event) => {
        props.changeParameter(opIndex, paramIndex, name, event.target.value);
    }

    const deleteParameter = (opIndex, paramIndex) => (event) => {
        props.deleteParameter(opIndex, paramIndex);
    }

    const addListParameter = (opIndex, paramIndex) => (event) => {
        props.addListParameter(opIndex, paramIndex);
    }

    const deleteListParameter = (opIndex, paramIndex, listIndex) => (event) => {
        props.deleteListParameter(opIndex, paramIndex, listIndex);
    }

    const editListParameter = (opIndex, paramIndex, listIndex) => (event) => {
        props.editListParameter(opIndex, paramIndex, listIndex, event.target.value);
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
                                <div key={index}>
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
            {/* OUTPUT SIZE */}
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
                                <div key={index}>
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
            {/* OPERATIONS */}
            {SectionInTypography("Operations", classes, "h6")} 
            <div>
                <Button variant={"contained"} color={"primary"} className={classes.addButton} onClick={props.addOperation}>
                        Add operation
                </Button>
            </div>
            {
                props.operations.map((operation, opIndex) => {
                    return (
                        <div key={opIndex}>
                            <div className={classes.section}>
                                <Grid container spacing={24}>
                                    <Grid item xs={4}>
                                        <Typography variant={"h5"}>
                                            OperationType
                                        </Typography>
                                    </Grid>
                                    <Grid item xs={7}>
                                        <TextField 
                                            value={operation.operationType}
                                            classes={classes.textField}
                                            onChange={changeOperation(opIndex)}
                                            />
                                    </Grid>
                                    <Grid item xs={1}>
                                        <IconButton
                                            key="close"
                                            aria-label="Close"
                                            color={"primary"}
                                            className={classes.fab}
                                            onClick={deleteOperation(opIndex)}
                                        >
                                            <DeleteIcon />
                                        </IconButton>
                                    </Grid>
                                    <hr />
                                </Grid>
                                <div>
                                    <Button variant={"contained"} color={"primary"} className={classes.addButton} onClick={addParameter(opIndex)}>
                                        Add parameter
                                    </Button>
                                </div>
                            </div>
                            {operation.parameterconfigs.map((param, paramIndex) => {
                                return (
                                    <div className={classes.parameter} key={paramIndex}>
                                        <Grid container alignItems={"center"}>
                                            <Grid item xs={1}>
                                                <TextField
                                                    label={"Name"}
                                                    className={classes.textField}
                                                    value={param.name}
                                                    onChange={changeParameter(opIndex, paramIndex, "name")}
                                                    />
                                            </Grid>
                                            <Grid item xs={2}>
                                                <FormControl variant="outlined" className={classes.formControl}>
                                                    <InputLabel>
                                                        Parameter Type
                                                    </InputLabel>
                                                    <Select
                                                        onChange={changeParameter(opIndex, paramIndex, "parameterType")}
                                                        value={param.parameterType}
                                                        input={
                                                            <OutlinedInput name={"paramType"} labelWidth={120}/>
                                                        }
                                                        className={classes.select}
                                                        >
                                                            {props.paramTypes.map((type, i) => {
                                                                return (
                                                                        <MenuItem value={type} key={i}>{type}</MenuItem>
                                                                    )
                                                            })}
                                                    </Select>
                                                </FormControl>
                                            </Grid>
                                            <Grid item xs={3}>
                                                <RadioGroup
                                                        aria-label="Gender"
                                                        name="gender1"
                                                        className={classes.group}
                                                        value={param.feasible}

                                                        onChange={changeParameter(opIndex, paramIndex, "feasible")}
                                                    >
                                                    <FormControlLabel value="feasible" control={<Radio color={"primary"} />} label="Feasible" />
                                                    <FormControlLabel value="list" control={<Radio color={"primary"} />} label="List" />
                                                </RadioGroup>
                                            </Grid>
                                            <Grid item xs={4}>
                                                {param.feasible === "list" && 
                                                    (param.list.map((element, elIndex) => {
                                                        return (
                                                            <div key={elIndex}>
                                                                <TextField
                                                                    className={classes.textField}
                                                                    value={element.value}
                                                                    onChange={editListParameter(opIndex, paramIndex, elIndex)}
                                                                />
                                                                <IconButton
                                                                    key="close"
                                                                    aria-label="Close"
                                                                    color={"primary"}
                                                                    className={classes.icon}
                                                                    onClick={deleteListParameter(opIndex, paramIndex, elIndex)}
                                                                    >
                                                                        <DeleteIcon />
                                                                </IconButton>
                                                            </div>
                                                        )
                                                    }))
                                                    
                                                }
                                                {param.feasible === "feasible" && 
                                                    <div>
                                                        <TextField
                                                            label={"Min"}
                                                            className={classes.textField}
                                                            value={param.min}

                                                            onChange={changeParameter(opIndex, paramIndex, "min")}
                                                        />
                                                        <TextField
                                                            label={"Max"}
                                                            className={classes.textField}
                                                            value={param.max}
                                                            onChange={changeParameter(opIndex, paramIndex, "max")}
                                                        />
                                                        <TextField
                                                            label={"Step"}
                                                            className={classes.textField}
                                                            value={param.step}
                                                            onChange={changeParameter(opIndex, paramIndex, "step")}
                                                        />
                                                    </div>
                                                }
                                            </Grid>
                                            <Grid item xs={1}>
                                                {param.feasible === "list" && 
                                                    <Fab color={"primary"} className={classes.fab} onClick={addListParameter(opIndex, paramIndex)}>
                                                        <AddIcon />
                                                    </Fab>
                                                }
                                            </Grid>
                                            <Grid item xs={1} >
                                                <IconButton
                                                        key="close"
                                                        aria-label="Close"
                                                        color={"primary"}
                                                        className={classes.fab}
                                                        onClick={deleteParameter(opIndex, paramIndex)}
                                                    >
                                                        <DeleteIcon />
                                                </IconButton>
                                            </Grid>
                                        </Grid>
                                    </div>
                                )
                            })}
                            <Divider />
                        </div>
                    )
                })
            }
        </div>
    )
}


const mapStateToProps = state => {
    return {
        inputSize: state[module].inputSize,
        outputSize: state[module].outputSize,
        operations: state[module].operations,
        paramTypes: state[module].paramTypes,
    }
}

export default connect(mapStateToProps, { addSize, editSize, deleteSize, addOperation, deleteOperation, changeOperation, addParameter, changeParameter, deleteParameter, addListParameter, editListParameter, deleteListParameter })(NASConfig);