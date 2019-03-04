import React, { useState } from 'react';
import { connect } from 'react-redux';
import makeStyles from '@material-ui/styles/makeStyles';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormHelperText from '@material-ui/core/FormHelperText';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';

const module = "hpCreate";


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
    }
})

const ParameterConfig = (props) => {
    
    const classes = useStyles();

    return (
        <div>
            {props.parameterConfig.map((param, i) => {
                return (
                    <div className={classes.parameter} key={i}>
                        <Grid container alignItems={"center"}>
                            <Grid item xs={1}>
                                <TextField
                                    label={"Name"}
                                    className={classes.textField}
                                    value={param.name}
                                    />
                            </Grid>
                            <Grid item xs={2}>
                                <FormControl variant="outlined" className={classes.formControl}>
                                    <InputLabel>
                                        Parameter Type
                                    </InputLabel>
                                    <Select
                                        value={param.type}
                                        // onChange={handleChange}
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
                                    >
                                    <FormControlLabel value="feasible" control={<Radio />} label="Feasible" />
                                    <FormControlLabel value="list" control={<Radio />} label="List" />
                                </RadioGroup>
                            </Grid>
                            <Grid item xs={5}>
                                <TextField
                                    label={"Value"}
                                    className={classes.textField}
                                    value={param.name}
                                    />
                            </Grid>
                            <Grid item xs={1} >
                                <IconButton
                                        key="close"
                                        aria-label="Close"
                                        color={"primary"}
                                        className={classes.icon}
                                    >
                                        <DeleteIcon />
                                </IconButton>
                            </Grid>
                        </Grid>
                    </div>
                )
            })}
        </div>
    )
}


const mapStateToProps = state => {
    return {
        parameterConfig: state[module].parameterConfig,
        paramTypes: state[module].paramTypes,
    }
}

export default connect(mapStateToProps, null)(ParameterConfig);