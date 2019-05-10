import React from 'react';
import {connect} from 'react-redux';
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import ScheduleIcon from '@material-ui/icons/Schedule';
import HighlightOffIcon from '@material-ui/icons/HighlightOff';
import DoneIcon from '@material-ui/icons/Done';
import { Link } from 'react-router-dom';
import { ListItemSecondaryAction, IconButton } from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';

import { openDeleteJobDialog } from '../../../actions/generalActions';
import DeleteDialog from '../../Menu/DeleteDialog';

const module = "hpMonitor";


const styles = theme => ({
    running: {
        color: theme.colors.running,
    },
    failed: {
        color: theme.colors.failed,
    },
    finished: {
        color: theme.colors.finished,
    },
});


const HPJobList = (props) => {

    const { classes } = props;

    const onDeleteJob = (experimentName) => (event) => {
        props.openDeleteJobDialog(experimentName);
    }

    return (
        <div>
            <List component="nav">
                {props.filteredJobsList.map((job, i) => {
                    let icon;
                    //TODO: Change status name
                    if (job.status === 'Running') {
                        icon = (<ScheduleIcon className={classes.running}/>)
                    } else if (job.status === 'Failed') {
                        icon = (<HighlightOffIcon className={classes.failed}/>)
                    } else {
                        icon = (<DoneIcon className={classes.finished}/>)
                    }
                    return (
                        <ListItem button key={i} component={Link} to={`/katib/hp_monitor/${job.name}`}>
                            <ListItemIcon>
                                {icon}
                            </ListItemIcon>
                            <ListItemText inset primary={job.name} />
                            <ListItemSecondaryAction>
                                <IconButton aria-label={"Delete"} onClick={onDeleteJob(job.name)}>
                                    <DeleteIcon />
                                </IconButton>
                            </ListItemSecondaryAction>
                        </ListItem>
                    );
                 })}
            </List>     
            <DeleteDialog />  
        </div>
    )
}

const mapStateToProps = (state) => {
    return {
        filteredJobsList: state[module].filteredJobsList,
    }
}

export default connect(mapStateToProps, { openDeleteJobDialog })(withStyles(styles)(HPJobList));