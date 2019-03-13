import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/core/styles';

import TextField from '@material-ui/core/TextField';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import Button from '@material-ui/core/Button';

import { filterJobs, changeType } from '../../../actions/hpMonitorActions';


const module = "hpMonitor";


const styles = theme => ({
    textField: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
    },
    filter: {
        margin: '0 auto',
        textAlign: 'center',
    },
});

const FilterPanel = (props) => {

    const { classes } = props;

    const handleType = (name) => (event) => {
        props.changeType(name, event.target.checked);
    }
    
    return (
        <div className={classes.filter}>
            <FormGroup row>
                <TextField
                    id="outlined-name"
                    label="Name"
                    className={classes.textField}
                    value={props.filter}
                    onChange={(event) => props.filterJobs(event.target.value)}
                    margin="normal"
                    variant="outlined"
                />
                {
                    Object.keys(props.filterType).map((filter, i) => {
                        return(
                            <FormControlLabel
                                key={i}
                                control={
                                    <Switch
                                        checked={props.filterType[filter]}
                                        onChange={handleType(filter)}
                                        value={filter}
                                        color={"primary "}
                                        />
                                    }
                                label={filter}
                            />
                        );
                    })
                }
            </FormGroup>
            <Button color={"secondary"} variant={"raised"}>
                Update
            </Button>
        </div>   
    )
}

const mapStateToProps = state => {
    return {
        filter: state[module].filter,
        filterType: state[module].filterType,
    }
}

export default connect(mapStateToProps, { filterJobs, changeType })(withStyles(styles)(FilterPanel));