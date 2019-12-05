import React from 'react';
import withStyles from '@material-ui/styles/withStyles';
import Grid from '@material-ui/core/Grid';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Typography from '@material-ui/core/Typography';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import { connect } from 'react-redux';
import { changeTrial, changeTrialNamespace } from '../../../../actions/nasCreateActions';
import { fetchTrialTemplates } from '../../../../actions/templateActions';

const module = "nasCreate";
const templateModule = "template";

const styles = theme => ({
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
        marginBottom: 10,
    },
    formControl: {
        margin: 4,
        width: '100%',
    },
    selectEmpty: {
        marginTop: 10,
    },
})


class TrialSpecParam extends React.Component {

    componentDidMount() {
        this.props.fetchTrialTemplates(this.props.trialNamespace);
        this.props.changeTrialNamespace(this.props.trialNamespace);
    }

    onTrialNamespaceChange = (event) => {
        this.props.fetchTrialTemplates(event.target.value);
        this.props.changeTrialNamespace(event.target.value);
    }

    onTrialChange = (event) => {
        this.props.changeTrial(event.target.value);
    }

    render() {
        const names = this.props.templates.map((template, i) => template.name)

        const { classes } = this.props
        return (
            <div>
                <div className={classes.parameter}>
                    <Grid container alignItems={"center"}>
                        <Grid item xs={12} sm={3}>
                            <Typography variant={"subheading"}>
                                <Tooltip title={"Trial namespace"}>
                                    <HelpOutlineIcon className={classes.help} color={"primary"} />
                                </Tooltip>
                                {"Namespace"}
                            </Typography>
                        </Grid>
                        <Grid item xs={12} sm={8}>
                            <TextField
                                className={"Trial Namespace"}
                                value={this.props.trialNamespace}
                                onChange={this.onTrialNamespaceChange}
                            />
                        </Grid>
                    </Grid>
                </div>
                <div className={classes.parameter}>
                    <Grid container alignItems={"center"}>
                        <Grid item xs={12} sm={3}>
                            <Typography variant={"subheading"}>
                                <Tooltip title={"Trial spec template"}>
                                    <HelpOutlineIcon className={classes.help} color={"primary"} />
                                </Tooltip>
                                {"TrialSpec"}
                            </Typography>
                        </Grid>
                        <Grid item xs={12} sm={8}>
                            <FormControl variant="outlined" className={classes.formControl}>
                                <InputLabel>
                                    Trial Spec
                                </InputLabel>
                                <Select
                                    value={this.props.trial}
                                    onChange={this.onTrialChange}
                                    input={
                                        <OutlinedInput name={"TrialSpec"} labelWidth={100} />
                                    }
                                    className={classes.select}
                                >
                                    {names.map((spec, i) => {
                                        return (
                                            <MenuItem value={spec} key={i}>{spec}</MenuItem>
                                        )
                                    })}
                                </Select>
                            </FormControl>
                        </Grid>
                    </Grid>
                </div>
            </div>
        )
    }
}


const mapStateToProps = state => {
    return {
        trial: state[module].trial,
        templates: state[templateModule].trialTemplates,
        trialNamespace: state[module].trialNamespace,
    }
}

export default connect(mapStateToProps, { changeTrialNamespace, changeTrial, fetchTrialTemplates })(withStyles(styles)(TrialSpecParam));
