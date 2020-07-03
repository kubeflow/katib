import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import { ListItemSecondaryAction, IconButton } from '@material-ui/core';
import ScheduleIcon from '@material-ui/icons/Schedule';
import RestoreIcon from '@material-ui/icons/Restore';
import HighlightOffIcon from '@material-ui/icons/HighlightOff';
import DoneIcon from '@material-ui/icons/Done';
import DeleteIcon from '@material-ui/icons/Delete';
import HourglassFullIcon from '@material-ui/icons/HourglassFull';

import DeleteDialog from './DeleteDialog';

import { openDeleteExperimentDialog } from '../../../actions/generalActions';
import {
  GENERAL_MODULE,
  EXPERIMENT_TYPE_HP,
  LINK_HP_MONITOR,
  LINK_NAS_MONITOR,
  EXPERIMENT_TYPE_NAS,
} from '../../../constants/constants';

const styles = theme => ({
  created: {
    color: theme.colors.created,
  },
  running: {
    color: theme.colors.running,
  },
  restarting: {
    color: theme.colors.restarting,
  },
  succeeded: {
    color: theme.colors.succeeded,
  },
  failed: {
    color: theme.colors.failed,
  },
});

const ExperimentList = props => {
  const { classes } = props;

  const onDeleteExperiment = (name, namespace) => event => {
    props.openDeleteExperimentDialog(name, namespace);
  };

  return (
    <div>
      <List component="nav">
        {props.filteredExperiments
          .filter(experiment =>
            props.experimentType === EXPERIMENT_TYPE_HP
              ? experiment.type === EXPERIMENT_TYPE_HP
              : experiment.type === EXPERIMENT_TYPE_NAS,
          )
          .map((experiment, i) => {
            let icon;
            if (experiment.status === 'Created') {
              icon = <HourglassFullIcon className={classes.created} />;
            } else if (experiment.status === 'Running') {
              icon = <ScheduleIcon className={classes.running} />;
            } else if (experiment.status === 'Restarting') {
              icon = <RestoreIcon className={classes.restarting} />;
            } else if (experiment.status === 'Succeeded') {
              icon = <DoneIcon className={classes.succeeded} />;
            } else if (experiment.status === 'Failed') {
              icon = <HighlightOffIcon className={classes.failed} />;
            }
            return (
              <ListItem
                button
                key={i}
                component={Link}
                to={
                  props.experimentType === EXPERIMENT_TYPE_HP
                    ? LINK_HP_MONITOR + '/' + experiment.namespace + '/' + experiment.name
                    : LINK_NAS_MONITOR + '/' + experiment.namespace + '/' + experiment.name
                }
              >
                <ListItemIcon>{icon}</ListItemIcon>
                <ListItemText
                  inset
                  primary={`${experiment.name}`}
                  secondary={experiment.namespace}
                />
                <ListItemSecondaryAction>
                  <IconButton
                    aria-label={'Delete'}
                    onClick={onDeleteExperiment(experiment.name, experiment.namespace)}
                  >
                    <DeleteIcon />
                  </IconButton>
                </ListItemSecondaryAction>
              </ListItem>
            );
          })}
      </List>
      <DeleteDialog />
    </div>
  );
};

const mapStateToProps = state => {
  return {
    filteredExperiments: state[GENERAL_MODULE].filteredExperiments,
  };
};

export default connect(mapStateToProps, { openDeleteExperimentDialog })(
  withStyles(styles)(ExperimentList),
);
