import React from 'react';
import { Route } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';

import Main from './Menu/Main';
import HPJobInfo from './HP/Monitor/HPJobInfo';
import NASJobInfo from './NAS/Monitor/NASJobInfo';
import Trial from './Templates/Trial';
import Header from './Menu/Header';
import Snack from './Menu/Snack';
import TabPanel from './Common/Create/TabPanel';
import ExperimentMonitor from './Common/Monitor/ExperimentMonitor';

import {
  LINK_HP_CREATE,
  LINK_NAS_CREATE,
  LINK_HP_MONITOR,
  LINK_NAS_MONITOR,
  LINK_TRIAL_TEMPLATE,
} from '../constants/constants';

const useStyles = makeStyles({
  root: {
    width: '90%',
    margin: '0 auto',
    paddingTop: 20,
  },
});

const App = props => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Header />
      <Route exact path="/" component={Main} />
      <Route path={LINK_HP_CREATE} component={TabPanel} />
      <Route exact path={LINK_HP_MONITOR} component={ExperimentMonitor} />
      <Route path={LINK_HP_MONITOR + '/:namespace/:name'} component={HPJobInfo} />
      <Route path={LINK_NAS_CREATE} component={TabPanel} />
      <Route exact path={LINK_NAS_MONITOR} component={ExperimentMonitor} />
      <Route path={LINK_NAS_MONITOR + '/:namespace/:name'} component={NASJobInfo} />
      <Route path={LINK_TRIAL_TEMPLATE} component={Trial} />
      <Snack />
    </div>
  );
};

export default App;
