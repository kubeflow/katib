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

import { openDeleteDialog } from '../../../actions/generalActions';
import DeleteDialog from '../../Menu/DeleteDialog';

const module = "nasMonitor";


const styles = theme => ({
    running: {
        color: '#8b8ffb',
    },
    failed: {
        color: '#f26363',
    },
    finished: {
        color: '#63f291',
    },
});


const JobList = (props) => {

    const { classes } = props;

    const onDelete = (id) => (event) => {
        props.openDeleteDialog(id);
    }

    return (
        <div>
            <List component="nav">
                {props.filteredJobsList.map((job, i) => {
                    let icon;
                    if (job.status === 'Running') {
                        icon = (<ScheduleIcon className={classes.running}/>)
                    } else if (job.status === 'Failed') {
                        icon = (<HighlightOffIcon className={classes.failed}/>)
                    } else {
                        icon = (<DoneIcon className={classes.finished}/>)
                    }
                    return (
                        <ListItem button key={i} component={Link} to={`/katib/nas_monitor/${job.id}`}>
                            <ListItemIcon>
                                {icon}
                            </ListItemIcon>
                            <ListItemText inset primary={job.name} />
                            <ListItemSecondaryAction>
                                <IconButton aria-label={"Delete"} onClick={onDelete(job.id)}>
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

export default connect(mapStateToProps, { openDeleteDialog })(withStyles(styles)(JobList));