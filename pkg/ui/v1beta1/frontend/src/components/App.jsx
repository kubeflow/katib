import React from 'react';
import { Route } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';

import Main from './Menu/Main';
import HPJobMonitor from './HP/Monitor/HPJobMonitor';
import HPJobInfo from './HP/Monitor/HPJobInfo';
import NASJobMonitor from './NAS/Monitor/NASJobMonitor';
import NASJobInfo from './NAS/Monitor/NASJobInfo';
import Trial from './Templates/Trial';
import Header from './Menu/Header';
import Snack from './Menu/Snack';
import TabPanel from './Common/Create/TabPanel';

import { LINK_HP_CREATE, LINK_NAS_CREATE } from '../constants/constants';

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
      <Route exact path="/katib/hp_monitor" component={HPJobMonitor} />
      <Route path="/katib/hp_monitor/:namespace/:name" component={HPJobInfo} />
      <Route path={LINK_NAS_CREATE} component={TabPanel} />
      <Route exact path="/katib/nas_monitor" component={NASJobMonitor} />
      <Route path="/katib/nas_monitor/:namespace/:name" component={NASJobInfo} />
      <Route path="/katib/trial" component={Trial} />
      <Snack />
    </div>
  );
};

export default App;
