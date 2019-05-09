import React  from 'react';

import Header from './Menu/Header';
import Snack from './Menu/Snack';

import { makeStyles } from '@material-ui/styles';

import  { Route } from 'react-router-dom';
import Main from './Menu/Main';
import HP from './HP/Create/HP';
import HPMonitor from './HP/Monitor/Monitor';
import HPJobInfo from './HP/Monitor/JobInfo';
import NAS from './NAS/Create/NAS';
import NASMonitor from './NAS/Monitor/Monitor';
import NASJobInfo from './NAS/Monitor/JobInfo';
import Trial from './Templates/Trial';
import Collector from './Templates/Collector';


const useStyles = makeStyles({
    root: {
        width: '90%',
        margin: '0 auto',
        paddingTop: 20,
    }
});

const App = (props) => { 
    const classes = useStyles();
    return (
        <div className={classes.root}>
            <Header />
            <Route exact path="/" component={Main} />
            <Route path="/katib/hp" component={HP} />
            <Route exact path="/katib/hp_monitor" component={HPMonitor} />
            <Route path="/katib/hp_monitor/:id" component={HPJobInfo} />
            <Route path="/katib/nas" component={NAS} />
            <Route exact path="/katib/nas_monitor" component={NASMonitor} />
            <Route path="/katib/nas_monitor/:id" component={NASJobInfo} />
            <Route path="/katib/trial" component={Trial} />
            <Route path="/katib/collector" component={Collector} />
            <Snack />
        </div>
    )
};

export default App;