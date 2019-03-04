import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';

import { connect } from 'react-redux';
import { changeWorker } from '../../../../actions/nasCreateActions';


const module = "nasCreate";

const useStyles = makeStyles({
    help: {
        padding: 4 / 2,
        verticalAlign: "middle",
        marginRight: 5,
    },
    section: {
        padding: 4,
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
})

const WorkerSpecParam = (props) => {

    const classes = useStyles();

    const onWorkerChange = (event) => {
        props.changeWorker(event.target.value);
    }

    return (
        <div className={classes.parameter}> 
            <Grid container alignItems={"center"}>
                <Grid item xs={12} sm={3}>
                    <Typography variant={"subheading"}>
                        <Tooltip title={"Worker spec template"}>
                            <HelpOutlineIcon className={classes.help} color={"primary"}/>
                        </Tooltip>
                        {"WorkerSpec"}
                    </Typography>
                </Grid>
                <Grid item xs={12} sm={8}>
                    <FormControl variant="outlined" className={classes.formControl}>
                        <InputLabel>
                            Worker Spec
                        </InputLabel>
                        <Select
                            value={props.worker}
                            onChange={onWorkerChange}
                            input={
                                <OutlinedInput name={"workerSpec"} labelWidth={100}/>
                            }
                            className={classes.select}
                            >
                                {props.workerSpec.map((spec, i) => {
                                    return (
                                            <MenuItem value={spec} key={i}>{spec}</MenuItem>
                                        )
                                })}
                        </Select>
                    </FormControl>
                </Grid>
            </Grid>
        </div>
    )
}
const mapStateToProps = state => {
    return {
        workerSpec: state[module].workerSpec,
        worker: state[module].worker,
    }
}

export default connect(mapStateToProps, { changeWorker })(WorkerSpecParam);