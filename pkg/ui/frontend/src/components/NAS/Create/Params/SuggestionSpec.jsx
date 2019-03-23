import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import TextField from '@material-ui/core/TextField';
import IconButton from '@material-ui/core/IconButton';

import DeleteIcon from '@material-ui/icons/Delete';

import { connect } from 'react-redux';
import { changeAlgorithm, addSuggestionParameter, changeSuggestionParameter, deleteSuggestionParameter, changeRequestNumber } from '../../../../actions/nasCreateActions';

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
    icon: {
        padding: 4,
        margin: '0 auto',
        verticalAlign: "middle !important",
    },
    formControl: {
        margin: 4,
        width: '100%',
    },
    addButton: {
        margin: 10,
    }
})

const SuggestionSpec = (props) => {
    
    const classes = useStyles();

    const onAlgorithmChange = (event) => {
        props.changeAlgorithm(event.target.value);
    }

    const onAddParameter = (event) => {
        props.addSuggestionParameter();
    }

    const onChangeParameter = (name, index) => (event) => {
        props.changeSuggestionParameter(index, name, event.target.value);
    }

    const onDeleteParameter = (index) => (event) => {
        props.deleteSuggestionParameter(index);
    }

    const onRequestNumberChange = (event) => {
        props.changeRequestNumber(event.target.value);
    }

    return (
        <div>
            <Button variant={"contained"} color={"primary"} className={classes.addButton} onClick={onAddParameter}>
                    Add parameter
            </Button>
            <div className={classes.parameter}> 
                <Grid container alignItems={"center"}>
                    <Grid item xs={12} sm={3}>
                        <Typography>
                            <Tooltip title={"Suggestion algorithm"}>
                                <HelpOutlineIcon className={classes.help} color={"primary"}/>
                            </Tooltip>
                            {"Suggestion Algorithm"}
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={8}>
                        <FormControl variant="outlined" className={classes.formControl}>
                            <InputLabel>
                                Suggestion Algorithm
                            </InputLabel>
                            <Select
                                value={props.suggestionAlgorithm}
                                onChange={onAlgorithmChange}
                                input={
                                    <OutlinedInput name={"SuggestionAlgorithm"} labelWidth={160}/>
                                }
                                className={classes.select}
                                >
                                    {props.suggestionAlgorithms.map((algorithm, i) => {
                                        return (
                                                <MenuItem value={algorithm} key={i}>{algorithm}</MenuItem>
                                            )
                                    })}
                            </Select>
                        </FormControl>
                    </Grid>
                </Grid>
            </div>
            <div className={classes.parameter}> 
                <Grid container alignItems={"center"}>
                    <Grid item xs={12} sm={3}>
                        <Typography>
                            <Tooltip title={"Number of trials in parallel"}>
                                <HelpOutlineIcon className={classes.help} color={"primary"}/>
                            </Tooltip>
                            {"RequestNumber"}
                        </Typography>
                    </Grid>
                    <Grid item xs={12} sm={8}>
                        <FormControl variant="outlined" className={classes.formControl}>
                            <TextField
                                label={"Request Number"}
                                className={classes.textField}
                                value={props.requestNumber}
                                onChange={onRequestNumberChange}
                                />
                        </FormControl>
                    </Grid>
                </Grid>
            </div>
            <br />
            {props.suggestionParameters.map((param, i) => {
                return (
                    <div key={i} className={classes.parameter}>
                        <Grid container alignItems={"center"}>
                            <Grid item xs={3} />
                            <Grid item xs={4}>
                                <TextField
                                    label={"Name"}
                                    className={classes.textField}
                                    value={param.name}
                                    onChange={onChangeParameter("name", i)}
                                    />
                            </Grid>
                            <Grid item xs={4}>
                                <TextField
                                    label={"Value"}
                                    className={classes.textField}
                                    value={param.value}
                                    onChange={onChangeParameter("value", i)}
                                    />
                            </Grid>
                            <Grid item xs={1} >
                                <IconButton
                                        key="close"
                                        aria-label="Close"
                                        color={"primary"}
                                        className={classes.icon}
                                        onClick={onDeleteParameter(i)}
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
        suggestionAlgorithm: state[module].suggestionAlgorithm,
        suggestionAlgorithms: state[module].suggestionAlgorithms,
        suggestionParameters: state[module].suggestionParameters,
    }
}

export default connect(mapStateToProps, { changeAlgorithm, addSuggestionParameter, changeSuggestionParameter, deleteSuggestionParameter, changeRequestNumber })(SuggestionSpec);