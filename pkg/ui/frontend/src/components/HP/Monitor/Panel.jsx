import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/core/styles';

import TextField from '@material-ui/core/TextField';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import Button from '@material-ui/core/Button';
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

    const {classes} = props;
    return (
        <div className={classes.filter}>
            <FormGroup row>
                <TextField
                    id="outlined-name"
                    label="Name"
                    className={classes.textField}
                    value={props.filter}
                    // onChange={}
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
                                        // onChange={this.handleType(filter)}
                                        value={filter}
                                        color={"secondary"}
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

export default connect(mapStateToProps, null)(withStyles(styles)(FilterPanel));